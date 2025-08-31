package seeders

import (
	"log"
	"semita/core/database/database_connections"
	"semita/core/database/generate_seeders"
)

// PermissionSeeder seeder para permisos
type PermissionSeeder struct {
	generate_seeders.BaseSeeder
}

// NewPermissionSeeder crea una nueva instancia del seeder
func NewPermissionSeeder() *PermissionSeeder {
	return &PermissionSeeder{
		BaseSeeder: generate_seeders.BaseSeeder{
			DB:   database_connections.DatabaseConnectSQL(),
			Name: "permission_seeder",
		},
	}
}

// GetName retorna el nombre del seeder
func (ps *PermissionSeeder) GetName() string {
	return ps.BaseSeeder.Name
}

// GetDependencies retorna las dependencias del seeder
func (ps *PermissionSeeder) GetDependencies() []string {
	return []string{} // No tiene dependencias
}

// GetTables retorna las tablas que maneja este seeder
func (ps *PermissionSeeder) GetTables() []string {
	return []string{"permissions"}
}

// Seed ejecuta el seeding de permisos exactamente como en Laravel
func (ps *PermissionSeeder) Seed() error {
	permissions := []string{
		// role
		"roles:index",
		"roles:create",
		"roles:store",
		"roles:show",
		"roles:edit",
		"roles:update",
		"roles:delete",

		// permission
		"permissions:index",
		"permissions:by-role",

		// user
		"users:index",
		"users:create",
		"users:store",
		"users:show",
		"users:edit",
		"users:update",
		"users:delete",
	}

	for _, permissionName := range permissions {
		// Insertar directamente usando SQL para coincidir exactamente con Laravel
		query := `INSERT INTO permissions (uuid, name, guard_name, created_at, updated_at) 
				  VALUES (UUID(), ?, ?, NOW(), NOW())`

		_, err := ps.BaseSeeder.DB.Exec(query, permissionName, "api")
		if err != nil {
			log.Printf("Error creating permission %s: %v", permissionName, err)
			return err
		}
	}

	log.Printf("Successfully created %d permissions", len(permissions))
	return nil
}
