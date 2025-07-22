package generate_migrations

import (
	"fmt"
	"semita/core/database/database_connections"
	"semita/core/helpers"
	"sort"
)

type Migrator struct {
	db         database_connections.SQLAdapter
	migrations []Migration
}

func NewMigrator(db database_connections.SQLAdapter) *Migrator {
	return &Migrator{
		db:         db,
		migrations: make([]Migration, 0),
	}
}

// Register registra una nueva migración
func (m *Migrator) Register(migration Migration) {
	m.migrations = append(m.migrations, migration)
}

// CreateMigrationsTable crea la tabla de migraciones si no existe
func (m *Migrator) CreateMigrationsTable() error {
	query := `
		CREATE TABLE IF NOT EXISTS generate_migrations (
			id INT PRIMARY KEY AUTO_INCREMENT,
			migration VARCHAR(255) NOT NULL,
			batch INT NOT NULL,
			executed_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)
	`
	_, err := m.db.Exec(query)
	if err != nil {
		helpers.Logs("info", fmt.Sprintf("❌ Error ejecutando query de creación de tabla: %v\n", err))
		fmt.Printf("📝 Query ejecutada: %s\n", query)
		return fmt.Errorf("error creating generate_migrations table: %v", err)
	}
	return nil
}

// Migrate ejecuta todas las migraciones pendientes
func (m *Migrator) Migrate() error {

	if err := m.CreateMigrationsTable(); err != nil {
		fmt.Printf("❌ Error creando tabla de migraciones: %v\n", err)
		return fmt.Errorf("error creating database table: %v", err)
	}

	executed, err := m.getExecutedMigrations()
	if err != nil {
		fmt.Printf("❌ Error obteniendo migraciones ejecutadas: %v\n", err)
		return fmt.Errorf("error fetching executed generate_migrations: %v", err)
	}

	// Ordenar migraciones por timestamp
	sort.Slice(m.migrations, func(i, j int) bool {
		return m.migrations[i].GetTimestamp() < m.migrations[j].GetTimestamp()
	})

	batch, err := m.getNextBatch()
	if err != nil {
		fmt.Printf("❌ Error obteniendo siguiente lote: %v\n", err)
		return fmt.Errorf("error getting next batch: %v", err)
	}

	executedCount := 0
	for _, migration := range m.migrations {
		migrationName := fmt.Sprintf("%s_%s", migration.GetTimestamp(), migration.GetName())

		if _, exists := executed[migrationName]; !exists {
			if err := migration.Up(m.db); err != nil {
				fmt.Printf("❌ Error ejecutando migración %s: %v\n", migrationName, err)
				return fmt.Errorf("error executing migration %s: %v", migrationName, err)
			}

			if err := m.recordMigration(migrationName, batch); err != nil {
				fmt.Printf("❌ Error registrando migración %s: %v\n", migrationName, err)
				return fmt.Errorf("error recording migration %s: %v", migrationName, err)
			}

			fmt.Printf("✅ Migrated: %s\n", migrationName)
			executedCount++
		} else {
			fmt.Printf("⏭️  Saltando (ya ejecutada): %s\n", migrationName)
		}
	}

	if executedCount == 0 {
		helpers.Logs("info", "ℹ️  No hay migraciones pendientes")
	} else {
		fmt.Printf("🎉 Se ejecutaron %d migraciones exitosamente\n", executedCount)
	}

	return nil
}

func (m *Migrator) Fresh() error {
	if err := dropAllMigrationsTable(m.db); err != nil {
		fmt.Printf("❌ Error eliminando tablas: %v\n", err)
		return fmt.Errorf("error dropping generate_migrations table: %v", err)
	}

	if err := m.CreateMigrationsTable(); err != nil {
		fmt.Printf("❌ Error recreando tabla de migraciones: %v\n", err)
		return fmt.Errorf("error recreating generate_migrations table: %v", err)
	}

	// Ejecutar todas las migraciones nuevamente
	if err := m.Migrate(); err != nil {
		fmt.Printf("❌ Error ejecutando migraciones después de fresh: %v\n", err)
		return fmt.Errorf("error running generate_migrations after fresh: %v", err)
	}

	return nil
}

// Rollback revierte el último lote de migraciones
func (m *Migrator) Rollback() error {
	lastBatch, err := m.getLastBatch()
	if err != nil {
		return err
	}

	if lastBatch == 0 {
		helpers.Logs("info", "Nothing to rollback")
		return nil
	}

	migrations, err := m.getMigrationsByBatch(lastBatch)
	if err != nil {
		return err
	}

	// Ejecutar rollback en orden inverso
	for i := len(migrations) - 1; i >= 0; i-- {
		migrationName := migrations[i]
		migration := m.findMigrationByName(migrationName)

		if migration == nil {
			return fmt.Errorf("migration %s not found in registered database", migrationName)
		}

		fmt.Printf("Rolling back: %s\n", migrationName)

		if err := migration.Down(m.db); err != nil {
			return fmt.Errorf("error rolling back migration %s: %v", migrationName, err)
		}

		if err := m.deleteMigrationRecord(migrationName); err != nil {
			return err
		}

		fmt.Printf("Rolled back: %s\n", migrationName)
	}

	return nil
}

func (m *Migrator) getExecutedMigrations() (map[string]bool, error) {
	rows, err := m.db.Query("SELECT migration FROM generate_migrations")
	if err != nil {
		return nil, fmt.Errorf("error querying executed generate_migrations: %v", err)
	}
	defer rows.Close()

	executed := make(map[string]bool)
	for rows.Next() {
		var migration string
		if err := rows.Scan(&migration); err != nil {
			return nil, fmt.Errorf("error scanning executed migration: %v", err)
		}
		executed[migration] = true
	}

	return executed, nil
}

func (m *Migrator) getNextBatch() (int, error) {
	var batch int
	err := m.db.QueryRow("SELECT COALESCE(MAX(batch), 0) + 1 FROM generate_migrations").Scan(&batch)
	return batch, err
}

func (m *Migrator) getLastBatch() (int, error) {
	var batch int
	err := m.db.QueryRow("SELECT COALESCE(MAX(batch), 0) FROM generate_migrations").Scan(&batch)
	return batch, err
}

func (m *Migrator) getMigrationsByBatch(batch int) ([]string, error) {
	rows, err := m.db.Query("SELECT migration FROM generate_migrations WHERE batch = ? ORDER BY id DESC", batch)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var migrations []string
	for rows.Next() {
		var migration string
		if err := rows.Scan(&migration); err != nil {
			return nil, err
		}
		migrations = append(migrations, migration)
	}

	return migrations, nil
}

func (m *Migrator) recordMigration(name string, batch int) error {
	_, err := m.db.Exec("INSERT INTO generate_migrations (migration, batch) VALUES (?, ?)", name, batch)
	return err
}

func (m *Migrator) deleteMigrationRecord(name string) error {
	_, err := m.db.Exec("DELETE FROM generate_migrations WHERE migration = ?", name)
	return err
}

func (m *Migrator) findMigrationByName(name string) Migration {
	for _, migration := range m.migrations {
		migrationName := fmt.Sprintf("%s_%s", migration.GetTimestamp(), migration.GetName())
		if migrationName == name {
			return migration
		}
	}
	return nil
}

// ExecuteSQL ejecuta una consulta SQL arbitraria
func (m *Migrator) ExecuteSQL(sql string) error {
	_, err := m.db.Exec(sql)
	return err
}

// GetDB retorna la conexión a la base de datos
func (m *Migrator) GetDB() database_connections.SQLAdapter {
	return m.db
}

func dropAllMigrationsTable(db database_connections.SQLAdapter) error {
	// Deshabilitar claves foráneas
	_, _ = db.Exec("SET FOREIGN_KEY_CHECKS = 0;")

	var rows, errorRows = db.Query("SHOW TABLES")
	if errorRows != nil {
		return errorRows
	}
	defer rows.Close()

	var tables []string
	for rows.Next() {
		var tableName string
		if errorScan := rows.Scan(&tableName); errorScan != nil {
			return errorScan
		}
		tables = append(tables, tableName)
	}

	// Eliminar primero las tablas que contienen 'tokens' (hijas)
	for _, tableName := range tables {
		if tableName == "generate_migrations" {
			continue
		}
		if containsToken(tableName) {
			_, errorExecute := db.Exec("DROP TABLE IF EXISTS " + tableName)
			if errorExecute != nil {
				return errorExecute
			}
		}
	}
	// Luego eliminar el resto
	for _, tableName := range tables {
		if tableName == "generate_migrations" {
			continue
		}
		if !containsToken(tableName) {
			_, errorExecute := db.Exec("DROP TABLE IF EXISTS " + tableName)
			if errorExecute != nil {
				return errorExecute
			}
		}
	}
	// Finalmente, eliminar la tabla de migraciones
	_, _ = db.Exec("DROP TABLE IF EXISTS generate_migrations")

	// Volver a habilitar claves foráneas
	_, _ = db.Exec("SET FOREIGN_KEY_CHECKS = 1;")

	return nil
}

func containsToken(tableName string) bool {
	// Busca si la palabra 'token' está en cualquier parte del nombre de la tabla
	return len(tableName) >= 5 && (tableName == "tokens" || tableName == "token" || containsSubstring(tableName, "token"))
}

func containsSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
