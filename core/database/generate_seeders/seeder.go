package generate_seeders

import (
	"fmt"
	"log"
	"semita/core/database/database_connections"
	"semita/core/helpers"
)

// Seeder interface que deben implementar todos los generate_seeders
type Seeder interface {
	Seed() error
	GetName() string
	GetDependencies() []string
	GetTables() []string // Nuevo método para especificar qué tablas maneja el seeder
}

// BaseSeeder estructura base para todos los generate_seeders
type BaseSeeder struct {
	DB   database_connections.SQLAdapter
	Name string
}

// SeederManager gestiona la ejecución de generate_seeders
type SeederManager struct {
	DB      database_connections.SQLAdapter
	seeders map[string]Seeder
}

// NewSeederManager crea una nueva instancia del manager
func NewSeederManager() *SeederManager {
	db := database_connections.DatabaseConnectSQL()
	return &SeederManager{
		DB:      db,
		seeders: make(map[string]Seeder),
	}
}

// RegisterSeeder registra un seeder en el manager
func (sm *SeederManager) RegisterSeeder(seeder Seeder) {
	sm.seeders[seeder.GetName()] = seeder
	helpers.Logs("INFO", fmt.Sprintf("Seeder '%s' registered successfully", seeder.GetName()))
}

// GetAllSeeders retorna todos los generate_seeders registrados
func (sm *SeederManager) GetAllSeeders() map[string]Seeder {
	return sm.seeders
}

// GetSeeder retorna un seeder específico por nombre
func (sm *SeederManager) GetSeeder(name string) (Seeder, error) {
	seeder, exists := sm.seeders[name]
	if !exists {
		helpers.Logs("ERROR", fmt.Sprintf("El seeder '%s' no existe", name))
		return nil, fmt.Errorf("seeder '%s' not found", name)
	}
	return seeder, nil
}

// RunSeeder ejecuta un seeder específico
func (sm *SeederManager) RunSeeder(name string) error {

	seeder, err := sm.GetSeeder(name)
	if err != nil {
		helpers.Logs("ERROR", fmt.Sprintf("Seeder '%s' not found: %v", name, err))
		return err
	}

	// Ejecutar dependencias primero
	dependencies := seeder.GetDependencies()
	for _, dep := range dependencies {
		log.Printf("Running dependency seeder: %s", dep)
		err := sm.RunSeeder(dep)
		if err != nil {
			helpers.Logs("ERROR", fmt.Sprintf("Error running dependency '%s': %v", dep, err))
			return fmt.Errorf("error running dependency '%s': %v", dep, err)
		}
	}

	// Ejecutar el seeder (sin transacción por ahora debido a limitaciones del adapter)
	log.Printf("Running seeder: %s", name)

	// Limpiar datos automáticamente antes del seeding
	err = sm.cleanSeederData(seeder)
	if err != nil {
		helpers.Logs("ERROR", fmt.Sprintf("Error cleaning data for seeder '%s': %v", name, err))
		return fmt.Errorf("error cleaning data for seeder '%s': %v", name, err)
	}

	// Ejecutar el seeding
	err = seeder.Seed()
	if err != nil {
		helpers.Logs("ERROR", fmt.Sprintf("Error running seeder '%s': %v", name, err))
		return fmt.Errorf("error running seeder '%s': %v", name, err)
	}

	log.Printf("Seeder '%s' executed successfully", name)
	return nil
}

// RunAllSeeders ejecuta todos los generate_seeders registrados
func (sm *SeederManager) RunAllSeeders() error {

	// Crear un grafo de dependencias y ejecutar en orden
	executed := make(map[string]bool)

	var runSeederWithDeps func(string) error
	runSeederWithDeps = func(name string) error {
		if executed[name] {
			return nil
		}

		seeder, exists := sm.seeders[name]
		if !exists {
			return fmt.Errorf("seeder '%s' not found", name)
		}

		// Ejecutar dependencias primero
		for _, dep := range seeder.GetDependencies() {
			err := runSeederWithDeps(dep)
			if err != nil {
				return err
			}
		}

		// Ejecutar el seeder (que incluye limpieza automática)
		err := sm.RunSeeder(name)
		if err != nil {
			fmt.Printf("Seeder '%s' failed: %v", name, err)
			return err
		}

		executed[name] = true
		fmt.Printf("✅ Seeder: %s\n", name)
		return nil
	}

	// Ejecutar todos los generate_seeders
	for name := range sm.seeders {
		err := runSeederWithDeps(name)
		if err != nil {
			return err
		}
	}

	return nil
}

// ResetSeeder ejecuta un seeder (ahora equivalente a RunSeeder ya que incluye limpieza automática)
func (sm *SeederManager) ResetSeeder(name string) error {
	log.Printf("Resetting seeder: %s", name)

	// Con el nuevo comportamiento, RunSeeder ya hace cleanup automáticamente
	err := sm.RunSeeder(name)
	if err != nil {
		return err
	}

	log.Printf("Seeder '%s' reset successfully", name)
	return nil
}

// cleanSeederData limpia automáticamente las tablas especificadas por un seeder
func (sm *SeederManager) cleanSeederData(seeder Seeder) error {
	tables := seeder.GetTables()
	if len(tables) == 0 {
		log.Printf("No tables specified for seeder '%s', skipping cleanup", seeder.GetName())
		return nil
	}

	log.Printf("Cleaning tables for seeder '%s': %v", seeder.GetName(), tables)

	// Deshabilitar temporalmente las verificaciones de claves foráneas
	_, err := sm.DB.Exec("SET FOREIGN_KEY_CHECKS = 0")
	if err != nil {
		log.Printf("Warning: could not disable foreign key checks: %v", err)
	}

	// Limpiar cada tabla en el orden especificado
	for _, table := range tables {
		query := fmt.Sprintf("DELETE FROM %s", table)
		_, err := sm.DB.Exec(query)
		if err != nil {
			log.Printf("Error cleaning table '%s': %v", table, err)
			// Continuar con las otras tablas
		} else {
			log.Printf("Cleaned table: %s", table)
		}
	}

	// Rehabilitar las verificaciones de claves foráneas
	_, err = sm.DB.Exec("SET FOREIGN_KEY_CHECKS = 1")
	if err != nil {
		log.Printf("Warning: could not re-enable foreign key checks: %v", err)
	}

	log.Printf("Data cleanup completed for seeder '%s'", seeder.GetName())
	return nil
}
