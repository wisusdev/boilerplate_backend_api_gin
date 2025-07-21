package seeders

import (
	"log"
	"semita/app/data/structs"
	"semita/core/database/database_connections"
	"semita/core/database/generate_seeders"
	"semita/core/roles_and_permissions/models_roles_and_permissions"
)

// RolesPermissionsSeeder seeder para roles y permisos
type RolesPermissionsSeeder struct {
	generate_seeders.BaseSeeder
}

// NewRolesPermissionsSeeder crea una nueva instancia del seeder
func NewRolesPermissionsSeeder() *RolesPermissionsSeeder {
	return &RolesPermissionsSeeder{
		BaseSeeder: generate_seeders.BaseSeeder{
			DB:   database_connections.DatabaseConnectSQL(),
			Name: "roles_permissions_seeder",
		},
	}
}

// GetName retorna el nombre del seeder
func (rps *RolesPermissionsSeeder) GetName() string {
	return rps.BaseSeeder.Name
}

// GetDependencies retorna las dependencias del seeder
func (rps *RolesPermissionsSeeder) GetDependencies() []string {
	return []string{} // No tiene dependencias
}

// GetTables retorna las tablas que maneja este seeder (en orden de limpieza)
func (rps *RolesPermissionsSeeder) GetTables() []string {
	return []string{"role_permissions", "user_roles", "user_permissions", "roles", "permissions"}
}

// Seed ejecuta el seeding de roles y permisos
func (rps *RolesPermissionsSeeder) Seed() error {
	log.Println("Seeding roles and permissions...")

	createdPermissions := rps.createPermissions()
	createdRoles := rps.createRoles()

	rps.assignAllPermissionsToRole("super-admin", createdRoles, createdPermissions)
	rps.assignPermissionsToRole("admin", createdRoles, createdPermissions, []string{
		"create-users", "edit-users", "delete-users", "view-users",
		"create-roles", "edit-roles", "view-roles", "assign-roles",
		"view-permissions", "assign-permissions",
		"manage-posts", "publish-posts", "edit-posts", "delete-posts",
		"view-dashboard", "manage-settings",
	})
	rps.assignPermissionsToRole("editor", createdRoles, createdPermissions, []string{
		"view-users",
		"manage-posts", "publish-posts", "edit-posts",
		"view-dashboard",
	})
	rps.assignPermissionsToRole("moderator", createdRoles, createdPermissions, []string{
		"view-users",
		"edit-posts",
		"view-dashboard",
	})

	log.Println("Roles and permissions seeding completed successfully!")
	return nil
}

func (rps *RolesPermissionsSeeder) createPermissions() map[string]*structs.PermissionStruct {
	log.Println("Creating permissions...")
	permissions := []structs.CreatePermissionStruct{
		{Name: "create-users", GuardName: "web", Description: "Crear usuarios"},
		{Name: "edit-users", GuardName: "web", Description: "Editar usuarios"},
		{Name: "delete-users", GuardName: "web", Description: "Eliminar usuarios"},
		{Name: "view-users", GuardName: "web", Description: "Ver usuarios"},
		{Name: "create-roles", GuardName: "web", Description: "Crear roles"},
		{Name: "edit-roles", GuardName: "web", Description: "Editar roles"},
		{Name: "delete-roles", GuardName: "web", Description: "Eliminar roles"},
		{Name: "view-roles", GuardName: "web", Description: "Ver roles"},
		{Name: "assign-roles", GuardName: "web", Description: "Asignar roles"},
		{Name: "create-permissions", GuardName: "web", Description: "Crear permisos"},
		{Name: "edit-permissions", GuardName: "web", Description: "Editar permisos"},
		{Name: "delete-permissions", GuardName: "web", Description: "Eliminar permisos"},
		{Name: "view-permissions", GuardName: "web", Description: "Ver permisos"},
		{Name: "assign-permissions", GuardName: "web", Description: "Asignar permisos"},
		{Name: "manage-posts", GuardName: "web", Description: "Gestionar posts"},
		{Name: "publish-posts", GuardName: "web", Description: "Publicar posts"},
		{Name: "edit-posts", GuardName: "web", Description: "Editar posts"},
		{Name: "delete-posts", GuardName: "web", Description: "Eliminar posts"},
		{Name: "view-dashboard", GuardName: "web", Description: "Ver dashboard administrativo"},
		{Name: "manage-settings", GuardName: "web", Description: "Gestionar configuración del sistema"},
	}
	createdPermissions := make(map[string]*structs.PermissionStruct)
	for _, permData := range permissions {
		// Crear el permiso directamente (ya se limpiaron los datos)
		permission, err := models_roles_and_permissions.CreatePermission(permData)
		if err != nil {
			log.Printf("Error creating permission '%s': %v", permData.Name, err)
			continue
		}
		createdPermissions[permData.Name] = permission
		log.Printf("Created permission: %s", permission.Name)
	}
	return createdPermissions
}

func (rps *RolesPermissionsSeeder) createRoles() map[string]*structs.RoleStruct {
	log.Println("Creating roles...")
	roles := []structs.CreateRoleStruct{
		{Name: "super-admin", GuardName: "web", Description: "Super administrador con todos los permisos"},
		{Name: "admin", GuardName: "web", Description: "Administrador del sistema"},
		{Name: "editor", GuardName: "web", Description: "Editor de contenido"},
		{Name: "moderator", GuardName: "web", Description: "Moderador"},
		{Name: "user", GuardName: "web", Description: "Usuario regular"},
	}
	createdRoles := make(map[string]*structs.RoleStruct)
	for _, roleData := range roles {
		// Crear el rol directamente (ya se limpiaron los datos)
		role, err := models_roles_and_permissions.CreateRole(roleData)
		if err != nil {
			log.Printf("Error creating role '%s': %v", roleData.Name, err)
			continue
		}
		createdRoles[roleData.Name] = role
		log.Printf("Created role: %s", role.Name)
	}
	return createdRoles
}

func (rps *RolesPermissionsSeeder) assignAllPermissionsToRole(roleName string, roles map[string]*structs.RoleStruct, permissions map[string]*structs.PermissionStruct) {
	if role, exists := roles[roleName]; exists {
		for _, permission := range permissions {
			err := models_roles_and_permissions.AssignPermissionToRole(role.ID, permission.ID)
			if err != nil {
				log.Printf("Error assigning permission '%s' to role '%s': %v", permission.Name, roleName, err)
			}
		}
		log.Printf("Assigned all permissions to %s role", roleName)
	}
}

func (rps *RolesPermissionsSeeder) assignPermissionsToRole(roleName string, roles map[string]*structs.RoleStruct, permissions map[string]*structs.PermissionStruct, permNames []string) {
	if role, exists := roles[roleName]; exists {
		for _, permName := range permNames {
			if permission, exists := permissions[permName]; exists {
				err := models_roles_and_permissions.AssignPermissionToRole(role.ID, permission.ID)
				if err != nil {
					log.Printf("Error assigning permission '%s' to role '%s': %v", permission.Name, roleName, err)
				}
			}
		}
		log.Printf("Assigned permissions to %s role", roleName)
	}
}
