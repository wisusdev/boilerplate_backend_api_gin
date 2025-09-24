package repositories

import (
	"boilerplate_backend_api_gin/core/database/database_connections"
	"boilerplate_backend_api_gin/core/roles_and_permissions/models"
	"fmt"
	"strings"
)

var rolesTable = "roles"
var userRolesTable = "model_has_roles"

func GetAllRoles() ([]models.RoleStruct, error) {
	database := database_connections.DatabaseConnectSQL()
	defer database.Close()

	query := `SELECT uuid, name, guard_name, created_at, updated_at FROM ` + rolesTable + ` ORDER BY name`
	rows, err := database.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var roles []models.RoleStruct
	for rows.Next() {
		var role models.RoleStruct
		err = rows.Scan(&role.ID, &role.Name, &role.GuardName, &role.CreatedAt, &role.UpdatedAt)
		if err != nil {
			return nil, err
		}
		roles = append(roles, role)
	}

	return roles, nil
}

func GetRoleByID(id string) (*models.RoleStruct, error) {
	database := database_connections.DatabaseConnectSQL()
	defer database.Close()

	query := `SELECT uuid, name, guard_name, created_at, updated_at FROM ` + rolesTable + ` WHERE uuid = ?`
	row := database.QueryRow(query, id)

	var role models.RoleStruct
	err := row.Scan(&role.ID, &role.Name, &role.GuardName, &role.CreatedAt, &role.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return &role, nil
}

func GetRoleByUUID(uuid string) (*models.RoleStruct, error) {
	database := database_connections.DatabaseConnectSQL()
	defer database.Close()

	query := `SELECT uuid, name, guard_name, created_at, updated_at FROM ` + rolesTable + ` WHERE uuid = ?`
	row := database.QueryRow(query, uuid)

	var role models.RoleStruct
	err := row.Scan(&role.ID, &role.Name, &role.GuardName, &role.CreatedAt, &role.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return &role, nil
}

func GetRoleByName(name string, guardName string) (*models.RoleStruct, error) {
	database := database_connections.DatabaseConnectSQL()
	defer database.Close()

	query := `SELECT id, name, guard_name, created_at, updated_at FROM ` + rolesTable + ` WHERE name = ? AND guard_name = ?`
	row := database.QueryRow(query, name, guardName)

	var role models.RoleStruct
	err := row.Scan(&role.ID, &role.Name, &role.GuardName, &role.CreatedAt, &role.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return &role, nil
}

func CreateRole(roleData models.CreateRoleStruct) (*models.RoleStruct, error) {
	database := database_connections.DatabaseConnectSQL()
	defer database.Close()

	if roleData.GuardName == "" {
		roleData.GuardName = "api"
	}

	// Validar permisos si se proporcionaron
	if len(roleData.Permissions) > 0 {
		validPermissions, invalidPermissions, err := ValidatePermissionsExist(roleData.Permissions, roleData.GuardName)
		if err != nil {
			return nil, err
		}
		if len(invalidPermissions) > 0 {
			return nil, fmt.Errorf("los siguientes permisos no existen: %s", strings.Join(invalidPermissions, ", "))
		}

		// Crear el rol
		query := `INSERT INTO ` + rolesTable + ` (uuid, name, guard_name, created_at, updated_at) VALUES (UUID(), ?, ?, NOW(), NOW())`
		_, err = database.Exec(query, roleData.Name, roleData.GuardName)
		if err != nil {
			return nil, err
		}

		// Obtener el UUID del rol recién creado
		var roleUUID string
		err = database.QueryRow("SELECT uuid FROM "+rolesTable+" WHERE name = ? AND guard_name = ? ORDER BY created_at DESC LIMIT 1", roleData.Name, roleData.GuardName).Scan(&roleUUID)
		if err != nil {
			return nil, err
		}

		// Asignar permisos al rol
		for _, permission := range validPermissions {
			assignQuery := `INSERT INTO ` + rolePermissionsTable + ` (role_id, permission_id) VALUES (?, ?)`
			_, err := database.Exec(assignQuery, roleUUID, permission.ID)
			if err != nil {
				return nil, fmt.Errorf("error al asignar permiso %s al rol: %v", permission.Name, err)
			}
		}

		// Retornar el rol creado
		return GetRoleByUUID(roleUUID)
	} else {
		// Crear rol sin permisos
		query := `INSERT INTO ` + rolesTable + ` (uuid, name, guard_name, created_at, updated_at) VALUES (UUID(), ?, ?, NOW(), NOW())`
		_, err := database.Exec(query, roleData.Name, roleData.GuardName)
		if err != nil {
			return nil, err
		}

		// Obtener el rol recién creado
		var roleUUID string
		err = database.QueryRow("SELECT uuid FROM "+rolesTable+" WHERE name = ? AND guard_name = ? ORDER BY created_at DESC LIMIT 1", roleData.Name, roleData.GuardName).Scan(&roleUUID)
		if err != nil {
			return nil, err
		}

		return GetRoleByUUID(roleUUID)
	}
}

func UpdateRole(id string, roleData models.CreateRoleStruct) (*models.RoleStruct, error) {
	database := database_connections.DatabaseConnectSQL()
	defer database.Close()

	query := `UPDATE ` + rolesTable + ` SET name = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?`
	_, err := database.Exec(query, roleData.Name, id)
	if err != nil {
		return nil, err
	}

	return GetRoleByID(id)
}

func UpdateRoleByUUID(uuid string, roleData models.CreateRoleStruct) (*models.RoleStruct, error) {
	database := database_connections.DatabaseConnectSQL()
	defer database.Close()

	if roleData.GuardName == "" {
		roleData.GuardName = "api"
	}

	// Validar permisos si se proporcionaron
	if len(roleData.Permissions) > 0 {
		validPermissions, invalidPermissions, err := ValidatePermissionsExist(roleData.Permissions, roleData.GuardName)
		if err != nil {
			return nil, err
		}
		if len(invalidPermissions) > 0 {
			return nil, fmt.Errorf("los siguientes permisos no existen: %s", strings.Join(invalidPermissions, ", "))
		}

		// Actualizar el rol
		query := `UPDATE ` + rolesTable + ` SET name = ?, updated_at = NOW() WHERE uuid = ?`
		_, err = database.Exec(query, roleData.Name, uuid)
		if err != nil {
			return nil, err
		}

		// Actualizar permisos del rol (eliminar los antiguos y agregar los nuevos)
		// Primero eliminar permisos actuales
		deleteQuery := `DELETE FROM ` + rolePermissionsTable + ` WHERE role_id = ?`
		_, err = database.Exec(deleteQuery, uuid)
		if err != nil {
			return nil, err
		}

		// Luego asignar los nuevos permisos
		for _, permission := range validPermissions {
			assignQuery := `INSERT INTO ` + rolePermissionsTable + ` (role_id, permission_id) VALUES (?, ?)`
			_, err = database.Exec(assignQuery, uuid, permission.ID)
			if err != nil {
				return nil, fmt.Errorf("error al asignar permiso %s al rol: %v", permission.Name, err)
			}
		}
	} else {
		// Solo actualizar el rol sin modificar permisos
		query := `UPDATE ` + rolesTable + ` SET name = ?, updated_at = NOW() WHERE uuid = ?`
		_, err := database.Exec(query, roleData.Name, uuid)
		if err != nil {
			return nil, err
		}
	}

	return GetRoleByUUID(uuid)
}

func DeleteRole(id string) error {
	database := database_connections.DatabaseConnectSQL()
	defer database.Close()

	query := `DELETE FROM ` + rolesTable + ` WHERE uuid = ?`
	_, err := database.Exec(query, id)
	return err
}

func GetUserRoles(userID string) ([]models.RoleStruct, error) {
	database := database_connections.DatabaseConnectSQL()
	defer database.Close()

	query := `
		SELECT r.uuid, r.name, r.guard_name, r.created_at, r.updated_at
		FROM ` + rolesTable + ` r
		INNER JOIN ` + userRolesTable + ` ur ON r.uuid = ur.role_id
		WHERE ur.model_id = ? AND ur.model_type = 'users'
		ORDER BY r.name
	`
	rows, err := database.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var roles []models.RoleStruct
	for rows.Next() {
		var role models.RoleStruct
		err = rows.Scan(&role.ID, &role.Name, &role.GuardName, &role.CreatedAt, &role.UpdatedAt)
		if err != nil {
			return nil, err
		}
		roles = append(roles, role)
	}

	return roles, nil
}

func AssignRoleToUser(userID string, roleID string) error {
	database := database_connections.DatabaseConnectSQL()
	defer database.Close()

	// Verificar si el usuario ya tiene el rol
	exists, err := UserHasRole(userID, roleID)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("user already has this role")
	}

	query := `INSERT INTO ` + userRolesTable + ` (model_id, role_id, model_type) VALUES (?, ?, 'users')`
	_, err = database.Exec(query, userID, roleID)
	return err
}

func RevokeRoleFromUser(userID string, roleID string) error {
	database := database_connections.DatabaseConnectSQL()
	defer database.Close()

	query := `DELETE FROM ` + userRolesTable + ` WHERE model_id = ? AND role_id = ? AND model_type = 'users'`
	_, err := database.Exec(query, userID, roleID)
	return err
}

func UserHasRole(userID string, roleID string) (bool, error) {
	database := database_connections.DatabaseConnectSQL()
	defer database.Close()

	query := `SELECT COUNT(*) FROM ` + userRolesTable + ` WHERE model_id = ? AND role_id = ? AND model_type = 'users'`
	var count int
	err := database.QueryRow(query, userID, roleID).Scan(&count)
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func UserHasRoleByName(userID string, roleName string, guardName string) (bool, error) {
	database := database_connections.DatabaseConnectSQL()
	defer database.Close()

	if guardName == "" {
		guardName = "web"
	}

	query := `
		SELECT COUNT(*) 
		FROM ` + userRolesTable + ` ur
		INNER JOIN ` + rolesTable + ` r ON ur.role_id = r.uuid
		WHERE ur.model_id = ? AND ur.model_type = 'users' AND r.name = ? AND r.guard_name = ?
	`
	var count int
	err := database.QueryRow(query, userID, roleName, guardName).Scan(&count)
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func UserHasAnyRole(userID string, roleNames []string, guardName string) (bool, error) {
	if len(roleNames) == 0 {
		return false, nil
	}

	database := database_connections.DatabaseConnectSQL()
	defer database.Close()

	if guardName == "" {
		guardName = "web"
	}

	placeholders := strings.Repeat("?,", len(roleNames))
	placeholders = placeholders[:len(placeholders)-1] // Remover la última coma

	query := fmt.Sprintf(`
		SELECT COUNT(*) 
		FROM %s ur
		INNER JOIN %s r ON ur.role_id = r.uuid
		WHERE ur.model_id = ? AND ur.model_type = 'users' AND r.name IN (%s) AND r.guard_name = ?
	`, userRolesTable, rolesTable, placeholders)

	args := make([]interface{}, 0, len(roleNames)+2)
	args = append(args, userID)
	for _, name := range roleNames {
		args = append(args, name)
	}
	args = append(args, guardName)

	var count int
	err := database.QueryRow(query, args...).Scan(&count)
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func UserHasAllRoles(userID string, roleNames []string, guardName string) (bool, error) {
	if len(roleNames) == 0 {
		return true, nil
	}

	database := database_connections.DatabaseConnectSQL()
	defer database.Close()

	if guardName == "" {
		guardName = "web"
	}

	placeholders := strings.Repeat("?,", len(roleNames))
	placeholders = placeholders[:len(placeholders)-1] // Remover la última coma

	query := fmt.Sprintf(`
		SELECT COUNT(*) 
		FROM %s ur
		INNER JOIN %s r ON ur.role_id = r.uuid
		WHERE ur.model_id = ? AND ur.model_type = 'users' AND r.name IN (%s) AND r.guard_name = ?
	`, userRolesTable, rolesTable, placeholders)

	args := make([]interface{}, 0, len(roleNames)+2)
	args = append(args, userID)
	for _, name := range roleNames {
		args = append(args, name)
	}
	args = append(args, guardName)

	var count int
	err := database.QueryRow(query, args...).Scan(&count)
	if err != nil {
		return false, err
	}

	return count == len(roleNames), nil
}
