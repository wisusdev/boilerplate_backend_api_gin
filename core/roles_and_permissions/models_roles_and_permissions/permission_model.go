package models_roles_and_permissions

import (
	"fmt"
	"semita/app/data/structs"
	"semita/common/nulltypes"
	"semita/core/database/database_connections"
	"strings"
)

var permissionsTable = "permissions"
var rolePermissionsTable = "role_permissions"
var userPermissionsTable = "user_permissions"

// GetAllPermissions obtiene todos los permisos
func GetAllPermissions() ([]structs.PermissionStruct, error) {
	database := database_connections.DatabaseConnectSQL()
	defer database.Close()

	query := `SELECT id, name, guard_name, description, created_at, updated_at FROM ` + permissionsTable + ` ORDER BY name`
	rows, err := database.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var permissions []structs.PermissionStruct
	for rows.Next() {
		var permission structs.PermissionStruct
		var description nulltypes.NullString
		err = rows.Scan(&permission.ID, &permission.Name, &permission.GuardName, &description, &permission.CreatedAt, &permission.UpdatedAt)
		if err != nil {
			return nil, err
		}
		permission.Description = description.String
		permissions = append(permissions, permission)
	}

	return permissions, nil
}

// GetPermissionByID obtiene un permiso por su ID
func GetPermissionByID(id int) (*structs.PermissionStruct, error) {
	database := database_connections.DatabaseConnectSQL()
	defer database.Close()

	query := `SELECT id, name, guard_name, description, created_at, updated_at FROM ` + permissionsTable + ` WHERE id = ?`
	row := database.QueryRow(query, id)

	var permission structs.PermissionStruct
	var description nulltypes.NullString
	err := row.Scan(&permission.ID, &permission.Name, &permission.GuardName, &description, &permission.CreatedAt, &permission.UpdatedAt)
	if err != nil {
		return nil, err
	}
	permission.Description = description.String

	return &permission, nil
}

// GetPermissionByName obtiene un permiso por su nombre
func GetPermissionByName(name string, guardName string) (*structs.PermissionStruct, error) {
	database := database_connections.DatabaseConnectSQL()
	defer database.Close()

	query := `SELECT id, name, guard_name, description, created_at, updated_at FROM ` + permissionsTable + ` WHERE name = ? AND guard_name = ?`
	row := database.QueryRow(query, name, guardName)

	var permission structs.PermissionStruct
	var description nulltypes.NullString
	err := row.Scan(&permission.ID, &permission.Name, &permission.GuardName, &description, &permission.CreatedAt, &permission.UpdatedAt)
	if err != nil {
		return nil, err
	}
	permission.Description = description.String

	return &permission, nil
}

// CreatePermission crea un nuevo permiso
func CreatePermission(permissionData structs.CreatePermissionStruct) (*structs.PermissionStruct, error) {
	database := database_connections.DatabaseConnectSQL()
	defer database.Close()

	if permissionData.GuardName == "" {
		permissionData.GuardName = "web"
	}

	query := `INSERT INTO ` + permissionsTable + ` (name, guard_name, description) VALUES (?, ?, ?)`
	result, err := database.Exec(query, permissionData.Name, permissionData.GuardName, permissionData.Description)
	if err != nil {
		return nil, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	return GetPermissionByID(int(id))
}

// UpdatePermission actualiza un permiso existente
func UpdatePermission(id int, permissionData structs.CreatePermissionStruct) (*structs.PermissionStruct, error) {
	database := database_connections.DatabaseConnectSQL()
	defer database.Close()

	query := `UPDATE ` + permissionsTable + ` SET name = ?, description = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?`
	_, err := database.Exec(query, permissionData.Name, permissionData.Description, id)
	if err != nil {
		return nil, err
	}

	return GetPermissionByID(id)
}

// DeletePermission elimina un permiso
func DeletePermission(id int) error {
	database := database_connections.DatabaseConnectSQL()
	defer database.Close()

	query := `DELETE FROM ` + permissionsTable + ` WHERE id = ?`
	_, err := database.Exec(query, id)
	return err
}

// GetRolePermissions obtiene todos los permisos de un rol
func GetRolePermissions(roleID int) ([]structs.PermissionStruct, error) {
	database := database_connections.DatabaseConnectSQL()
	defer database.Close()

	query := `
		SELECT p.id, p.name, p.guard_name, p.description, p.created_at, p.updated_at 
		FROM ` + permissionsTable + ` p
		INNER JOIN ` + rolePermissionsTable + ` rp ON p.id = rp.permission_id
		WHERE rp.role_id = ?
		ORDER BY p.name
	`
	rows, err := database.Query(query, roleID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var permissions []structs.PermissionStruct
	for rows.Next() {
		var permission structs.PermissionStruct
		var description nulltypes.NullString
		err = rows.Scan(&permission.ID, &permission.Name, &permission.GuardName, &description, &permission.CreatedAt, &permission.UpdatedAt)
		if err != nil {
			return nil, err
		}
		permission.Description = description.String
		permissions = append(permissions, permission)
	}

	return permissions, nil
}

// GetUserDirectPermissions obtiene los permisos directos de un usuario (no heredados de roles)
func GetUserDirectPermissions(userID int) ([]structs.PermissionStruct, error) {
	database := database_connections.DatabaseConnectSQL()
	defer database.Close()

	query := `
		SELECT p.id, p.name, p.guard_name, p.description, p.created_at, p.updated_at 
		FROM ` + permissionsTable + ` p
		INNER JOIN ` + userPermissionsTable + ` up ON p.id = up.permission_id
		WHERE up.user_id = ?
		ORDER BY p.name
	`
	rows, err := database.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var permissions []structs.PermissionStruct
	for rows.Next() {
		var permission structs.PermissionStruct
		var description nulltypes.NullString
		err = rows.Scan(&permission.ID, &permission.Name, &permission.GuardName, &description, &permission.CreatedAt, &permission.UpdatedAt)
		if err != nil {
			return nil, err
		}
		permission.Description = description.String
		permissions = append(permissions, permission)
	}

	return permissions, nil
}

// GetUserAllPermissions obtiene todos los permisos de un usuario (directos + heredados de roles)
func GetUserAllPermissions(userID int) ([]structs.PermissionStruct, error) {
	database := database_connections.DatabaseConnectSQL()
	defer database.Close()

	query := `
		(
			SELECT DISTINCT p.id, p.name, p.guard_name, p.description, p.created_at, p.updated_at 
			FROM ` + permissionsTable + ` p
			INNER JOIN ` + userPermissionsTable + ` up ON p.id = up.permission_id
			WHERE up.user_id = ?
		)
		UNION
		(
			SELECT DISTINCT p.id, p.name, p.guard_name, p.description, p.created_at, p.updated_at 
			FROM ` + permissionsTable + ` p
			INNER JOIN ` + rolePermissionsTable + ` rp ON p.id = rp.permission_id
			INNER JOIN ` + userRolesTable + ` ur ON rp.role_id = ur.role_id
			WHERE ur.user_id = ?
		)
		ORDER BY name
	`
	rows, err := database.Query(query, userID, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var permissions []structs.PermissionStruct
	for rows.Next() {
		var permission structs.PermissionStruct
		var description nulltypes.NullString
		err = rows.Scan(&permission.ID, &permission.Name, &permission.GuardName, &description, &permission.CreatedAt, &permission.UpdatedAt)
		if err != nil {
			return nil, err
		}
		permission.Description = description.String
		permissions = append(permissions, permission)
	}

	return permissions, nil
}

// AssignPermissionToRole asigna un permiso a un rol
func AssignPermissionToRole(roleID int, permissionID int) error {
	database := database_connections.DatabaseConnectSQL()
	defer database.Close()

	// Verificar si el rol ya tiene el permiso
	exists, err := RoleHasPermission(roleID, permissionID)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("role already has this permission")
	}

	query := `INSERT INTO ` + rolePermissionsTable + ` (role_id, permission_id) VALUES (?, ?)`
	_, err = database.Exec(query, roleID, permissionID)
	return err
}

// RevokePermissionFromRole revoca un permiso de un rol
func RevokePermissionFromRole(roleID int, permissionID int) error {
	database := database_connections.DatabaseConnectSQL()
	defer database.Close()

	query := `DELETE FROM ` + rolePermissionsTable + ` WHERE role_id = ? AND permission_id = ?`
	_, err := database.Exec(query, roleID, permissionID)
	return err
}

// AssignPermissionToUser asigna un permiso directamente a un usuario
func AssignPermissionToUser(userID int, permissionID int) error {
	database := database_connections.DatabaseConnectSQL()
	defer database.Close()

	// Verificar si el usuario ya tiene el permiso directamente
	exists, err := UserHasDirectPermission(userID, permissionID)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("user already has this direct permission")
	}

	query := `INSERT INTO ` + userPermissionsTable + ` (user_id, permission_id) VALUES (?, ?)`
	_, err = database.Exec(query, userID, permissionID)
	return err
}

// RevokePermissionFromUser revoca un permiso directo de un usuario
func RevokePermissionFromUser(userID int, permissionID int) error {
	database := database_connections.DatabaseConnectSQL()
	defer database.Close()

	query := `DELETE FROM ` + userPermissionsTable + ` WHERE user_id = ? AND permission_id = ?`
	_, err := database.Exec(query, userID, permissionID)
	return err
}

// RoleHasPermission verifica si un rol tiene un permiso específico
func RoleHasPermission(roleID int, permissionID int) (bool, error) {
	database := database_connections.DatabaseConnectSQL()
	defer database.Close()

	query := `SELECT COUNT(*) FROM ` + rolePermissionsTable + ` WHERE role_id = ? AND permission_id = ?`
	var count int
	err := database.QueryRow(query, roleID, permissionID).Scan(&count)
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

// UserHasDirectPermission verifica si un usuario tiene un permiso directo
func UserHasDirectPermission(userID int, permissionID int) (bool, error) {
	database := database_connections.DatabaseConnectSQL()
	defer database.Close()

	query := `SELECT COUNT(*) FROM ` + userPermissionsTable + ` WHERE user_id = ? AND permission_id = ?`
	var count int
	err := database.QueryRow(query, userID, permissionID).Scan(&count)
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

// UserHasPermission verifica si un usuario tiene un permiso (directo o heredado)
func UserHasPermission(userID int, permissionName string, guardName string) (bool, error) {
	database := database_connections.DatabaseConnectSQL()
	defer database.Close()

	if guardName == "" {
		guardName = "web"
	}

	query := `
		SELECT COUNT(*) FROM (
			(
				SELECT 1 
				FROM ` + userPermissionsTable + ` up
				INNER JOIN ` + permissionsTable + ` p ON up.permission_id = p.id
				WHERE up.user_id = ? AND p.name = ? AND p.guard_name = ?
			)
			UNION
			(
				SELECT 1 
				FROM ` + userRolesTable + ` ur
				INNER JOIN ` + rolePermissionsTable + ` rp ON ur.role_id = rp.role_id
				INNER JOIN ` + permissionsTable + ` p ON rp.permission_id = p.id
				WHERE ur.user_id = ? AND p.name = ? AND p.guard_name = ?
			)
		) AS combined_permissions
	`
	var count int
	err := database.QueryRow(query, userID, permissionName, guardName, userID, permissionName, guardName).Scan(&count)
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

// UserHasAnyPermission verifica si un usuario tiene al menos uno de los permisos especificados
func UserHasAnyPermission(userID int, permissionNames []string, guardName string) (bool, error) {
	if len(permissionNames) == 0 {
		return false, nil
	}

	database := database_connections.DatabaseConnectSQL()
	defer database.Close()

	if guardName == "" {
		guardName = "web"
	}

	placeholders := strings.Repeat("?,", len(permissionNames))
	placeholders = placeholders[:len(placeholders)-1] // Remover la última coma

	query := fmt.Sprintf(`
		SELECT COUNT(*) FROM (
			(
				SELECT 1 
				FROM %s up
				INNER JOIN %s p ON up.permission_id = p.id
				WHERE up.user_id = ? AND p.name IN (%s) AND p.guard_name = ?
			)
			UNION
			(
				SELECT 1 
				FROM %s ur
				INNER JOIN %s rp ON ur.role_id = rp.role_id
				INNER JOIN %s p ON rp.permission_id = p.id
				WHERE ur.user_id = ? AND p.name IN (%s) AND p.guard_name = ?
			)
		) AS combined_permissions
	`, userPermissionsTable, permissionsTable, placeholders, userRolesTable, rolePermissionsTable, permissionsTable, placeholders)

	args := make([]interface{}, 0, len(permissionNames)*2+4)
	args = append(args, userID)
	for _, name := range permissionNames {
		args = append(args, name)
	}
	args = append(args, guardName, userID)
	for _, name := range permissionNames {
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

// UserHasAllPermissions verifica si un usuario tiene todos los permisos especificados
func UserHasAllPermissions(userID int, permissionNames []string, guardName string) (bool, error) {
	if len(permissionNames) == 0 {
		return true, nil
	}

	// Verificar cada permiso individualmente
	for _, permissionName := range permissionNames {
		hasPermission, err := UserHasPermission(userID, permissionName, guardName)
		if err != nil {
			return false, err
		}
		if !hasPermission {
			return false, nil
		}
	}

	return true, nil
}
