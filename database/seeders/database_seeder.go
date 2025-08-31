package seeders

import (
	"log"
	"semita/core/database/database_connections"
	"semita/core/database/generate_seeders"
)

// DatabaseSeeder seeder principal que ejecuta todos los seeders
type DatabaseSeeder struct {
	generate_seeders.BaseSeeder
}

// NewDatabaseSeeder crea una nueva instancia del seeder principal
func NewDatabaseSeeder() *DatabaseSeeder {
	return &DatabaseSeeder{
		BaseSeeder: generate_seeders.BaseSeeder{
			DB:   database_connections.DatabaseConnectSQL(),
			Name: "database_seeder",
		},
	}
}

// GetName retorna el nombre del seeder
func (ds *DatabaseSeeder) GetName() string {
	return ds.BaseSeeder.Name
}

// GetDependencies retorna las dependencias del seeder
func (ds *DatabaseSeeder) GetDependencies() []string {
	return []string{} // No tiene dependencias
}

// GetTables retorna las tablas que maneja este seeder
func (ds *DatabaseSeeder) GetTables() []string {
	return []string{} // No maneja tablas directamente
}

// Seed ejecuta todos los seeders en el orden correcto, exactamente como Laravel
func (ds *DatabaseSeeder) Seed() error {
	log.Println("Starting DatabaseSeeder...")

	// Crear y ejecutar seeders en el orden exacto de Laravel
	seeders := []generate_seeders.Seeder{
		NewPermissionSeeder(),
		NewRoleSeeder(),
		NewUserSeeder(),
	}

	for _, seeder := range seeders {
		log.Printf("Running %s...", seeder.GetName())
		err := seeder.Seed()
		if err != nil {
			log.Printf("Error running %s: %v", seeder.GetName(), err)
			return err
		}
		log.Printf("Completed %s successfully", seeder.GetName())
	}

	log.Println("DatabaseSeeder completed successfully!")
	return nil
}

// CreateSeederManager crea y configura el manager de seeders
func CreateSeederManager() *generate_seeders.SeederManager {
	manager := generate_seeders.NewSeederManager()

	// Solo registrar el DatabaseSeeder principal, como en Laravel
	// Este internamente ejecutará PermissionSeeder, RoleSeeder y UserSeeder
	manager.RegisterSeeder(NewDatabaseSeeder())

	return manager
}
