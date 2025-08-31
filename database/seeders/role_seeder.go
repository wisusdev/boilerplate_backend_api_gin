package seeders

import (
	"log"
	"semita/core/database/database_connections"
	"semita/core/database/generate_seeders"
)

// RoleSeeder seeder para roles
type RoleSeeder struct {
	generate_seeders.BaseSeeder
}

// NewRoleSeeder crea una nueva instancia del seeder
func NewRoleSeeder() *RoleSeeder {
	return &RoleSeeder{
		BaseSeeder: generate_seeders.BaseSeeder{
			DB:   database_connections.DatabaseConnectSQL(),
			Name: "role_seeder",
		},
	}
}

// GetName retorna el nombre del seeder
func (rs *RoleSeeder) GetName() string {
	return rs.BaseSeeder.Name
}

// GetDependencies retorna las dependencias del seeder
func (rs *RoleSeeder) GetDependencies() []string {
	return []string{"permission_seeder"} // Depende del PermissionSeeder
}

// GetTables retorna las tablas que maneja este seeder
func (rs *RoleSeeder) GetTables() []string {
	return []string{"roles", "role_has_permissions"}
}

// Seed ejecuta el seeding de roles exactamente como en Laravel
func (rs *RoleSeeder) Seed() error {
	roles := []string{
		"user",
		"admin",
	}

	for _, roleName := range roles {
		// Insertar rol directamente usando SQL
		roleQuery := `INSERT INTO roles (uuid, name, guard_name, created_at, updated_at) 
					  VALUES (UUID(), ?, ?, NOW(), NOW())`

		_, err := rs.BaseSeeder.DB.Exec(roleQuery, roleName, "api")
		if err != nil {
			log.Printf("Error creating role %s: %v", roleName, err)
			return err
		}

		// Si es admin, darle todos los permisos
		if roleName == "admin" {
			// Obtener UUID del rol recién creado
			var roleUUID string
			err = rs.BaseSeeder.DB.QueryRow("SELECT uuid FROM roles WHERE name = ? AND guard_name = ?", roleName, "api").Scan(&roleUUID)
			if err != nil {
				log.Printf("Error getting role UUID for %s: %v", roleName, err)
				return err
			}

			// Obtener todos los permisos y asignarlos al rol admin
			permissionQuery := `SELECT uuid FROM permissions`
			rows, err := rs.BaseSeeder.DB.Query(permissionQuery)
			if err != nil {
				log.Printf("Error getting permissions: %v", err)
				return err
			}
			defer rows.Close()

			for rows.Next() {
				var permissionUUID string
				err = rows.Scan(&permissionUUID)
				if err != nil {
					log.Printf("Error scanning permission UUID: %v", err)
					continue
				}

				// Asignar permiso al rol
				assignQuery := `INSERT INTO role_has_permissions (permission_id, role_id) VALUES (?, ?)`
				_, err = rs.BaseSeeder.DB.Exec(assignQuery, permissionUUID, roleUUID)
				if err != nil {
					log.Printf("Error assigning permission %s to role %s: %v", permissionUUID, roleName, err)
				}
			}
		}
	}

	log.Printf("Successfully created %d roles", len(roles))
	return nil
}
