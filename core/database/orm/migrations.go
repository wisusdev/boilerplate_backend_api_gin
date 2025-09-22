package orm

import (
	"boilerplate_backend_api_gin/core/database/database_connections"
	"fmt"
	"reflect"
	"strings"
	"time"
)

// Migration representa una migración de base de datos
type Migration struct {
	Name      string
	Timestamp time.Time
	Up        func(database_connections.SQLAdapter) error
	Down      func(database_connections.SQLAdapter) error
}

// MigrationManager maneja las migraciones de la base de datos
type MigrationManager struct {
	db            database_connections.SQLAdapter
	migrationsDir string
	migrations    []Migration
}

// NewMigrationManager crea un nuevo gestor de migraciones
func NewMigrationManager() *MigrationManager {
	return &MigrationManager{
		db:            database_connections.DatabaseConnectSQL(),
		migrationsDir: "database/migrations",
		migrations:    []Migration{},
	}
}

// RegisterMigration registra una nueva migración
func (m *MigrationManager) RegisterMigration(migration Migration) {
	m.migrations = append(m.migrations, migration)
}

// CreateMigrationsTable crea la tabla de migraciones si no existe
func (m *MigrationManager) CreateMigrationsTable() error {
	query := `
	CREATE TABLE IF NOT EXISTS migrations (
		id INT AUTO_INCREMENT PRIMARY KEY,
		migration VARCHAR(255) NOT NULL,
		batch INT NOT NULL,
		applied_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);
	`
	_, err := m.db.Exec(query)
	return err
}

// GetAppliedMigrations obtiene las migraciones ya aplicadas
func (m *MigrationManager) GetAppliedMigrations() (map[string]int, error) {
	query := "SELECT migration, batch FROM migrations"
	rows, err := m.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	applied := make(map[string]int)
	for rows.Next() {
		var migration string
		var batch int
		err := rows.Scan(&migration, &batch)
		if err != nil {
			return nil, err
		}
		applied[migration] = batch
	}

	return applied, nil
}

// GetLastBatch obtiene el último lote de migraciones aplicado
func (m *MigrationManager) GetLastBatch() (int, error) {
	query := "SELECT MAX(batch) FROM migrations"
	var batch int
	err := m.db.QueryRow(query).Scan(&batch)
	if err != nil {
		// Si no hay migraciones, devolver 0
		return 0, nil
	}
	return batch, nil
}

// Migrate ejecuta las migraciones pendientes
func (m *MigrationManager) Migrate() error {
	// Crear tabla de migraciones si no existe
	err := m.CreateMigrationsTable()
	if err != nil {
		return err
	}

	// Obtener migraciones ya aplicadas
	applied, err := m.GetAppliedMigrations()
	if err != nil {
		return err
	}

	// Obtener último lote
	lastBatch, err := m.GetLastBatch()
	if err != nil {
		return err
	}

	// Ejecutar migraciones pendientes
	newBatch := lastBatch + 1
	var migrationsRun int

	for _, migration := range m.migrations {
		// Verificar si la migración ya fue aplicada
		if _, exists := applied[migration.Name]; exists {
			continue
		}

		// Ejecutar migración
		fmt.Printf("Migrando: %s\n", migration.Name)
		err := migration.Up(m.db)
		if err != nil {
			return fmt.Errorf("error en migración %s: %v", migration.Name, err)
		}

		// Registrar migración como aplicada
		_, err = m.db.Exec("INSERT INTO migrations (migration, batch) VALUES (?, ?)", migration.Name, newBatch)
		if err != nil {
			return err
		}

		migrationsRun++
	}

	fmt.Printf("Migraciones completadas: %d\n", migrationsRun)
	return nil
}

// Rollback revierte el último lote de migraciones
func (m *MigrationManager) Rollback() error {
	// Obtener migraciones ya aplicadas
	applied, err := m.GetAppliedMigrations()
	if err != nil {
		return err
	}

	// Obtener último lote
	lastBatch, err := m.GetLastBatch()
	if err != nil {
		return err
	}

	// Revertir migraciones del último lote
	var migrationsRolledBack int

	// Iterar en orden inverso para revertir correctamente
	for i := len(m.migrations) - 1; i >= 0; i-- {
		migration := m.migrations[i]

		// Verificar si la migración está en el último lote
		batch, exists := applied[migration.Name]
		if !exists || batch != lastBatch {
			continue
		}

		// Revertir migración
		fmt.Printf("Revirtiendo: %s\n", migration.Name)
		err := migration.Down(m.db)
		if err != nil {
			return fmt.Errorf("error al revertir migración %s: %v", migration.Name, err)
		}

		// Eliminar registro de la migración
		_, err = m.db.Exec("DELETE FROM migrations WHERE migration = ?", migration.Name)
		if err != nil {
			return err
		}

		migrationsRolledBack++
	}

	fmt.Printf("Migraciones revertidas: %d\n", migrationsRolledBack)
	return nil
}

// CreateTableMigration crea una migración para crear una tabla
func CreateTableMigration(name string, model interface{}) Migration {
	return Migration{
		Name:      name,
		Timestamp: time.Now(),
		Up: func(db database_connections.SQLAdapter) error {
			// Generar SQL para crear la tabla
			sql := GenerateCreateTableSQL(model)
			_, err := db.Exec(sql)
			return err
		},
		Down: func(db database_connections.SQLAdapter) error {
			// Obtener nombre de la tabla
			var tableName string
			if m, ok := model.(Model); ok {
				tableName = m.TableName()
			} else {
				v := reflect.ValueOf(model)
				if v.Kind() == reflect.Ptr {
					v = v.Elem()
				}
				tableName = strings.ToLower(v.Type().Name())
			}

			// Drop table
			_, err := db.Exec(fmt.Sprintf("DROP TABLE IF EXISTS %s", tableName))
			return err
		},
	}
}
