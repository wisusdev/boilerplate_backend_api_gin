package repositories

import (
	"boilerplate_backend_api_gin/core/database/database_connections"
	"boilerplate_backend_api_gin/core/roles_and_permissions/models"
	"fmt"
	"strings"
)

var permissionsTable = "permissions"
var rolePermissionsTable = "role_has_permissions"
var userPermissionsTable = "model_has_permissions"

func GetAllPermissions() ([]models.PermissionStruct, error) {
	database := database_connections.DatabaseConnectSQL()
	defer database.Close()

	query := `SELECT uuid, name, guard_name, created_at, updated_at FROM ` + permissionsTable + ` ORDER BY name`
	rows, err := database.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var permissions []models.PermissionStruct
	for rows.Next() {
		var permission models.PermissionStruct
		err = rows.Scan(&permission.ID, &permission.Name, &permission.GuardName, &permission.CreatedAt, &permission.UpdatedAt)
		if err != nil {
			return nil, err
		}
		permissions = append(permissions, permission)
	}

	return permissions, nil
}

func GetPermissionByID(id int) (*models.PermissionStruct, error) {
	database := database_connections.DatabaseConnectSQL()
	defer database.Close()

	query := `SELECT id, name, guard_name, created_at, updated_at FROM ` + permissionsTable + ` WHERE id = ?`
	row := database.QueryRow(query, id)

	var permission models.PermissionStruct
	err := row.Scan(&permission.ID, &permission.Name, &permission.GuardName, &permission.CreatedAt, &permission.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return &permission, nil
}

func GetPermissionByName(name string, guardName string) (*models.PermissionStruct, error) {
	database := database_connections.DatabaseConnectSQL()
	defer database.Close()

	query := `SELECT id, name, guard_name, created_at, updated_at FROM ` + permissionsTable + ` WHERE name = ? AND guard_name = ?`
	row := database.QueryRow(query, name, guardName)

	var permission models.PermissionStruct
	err := row.Scan(&permission.ID, &permission.Name, &permission.GuardName, &permission.CreatedAt, &permission.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return &permission, nil
}

func CreatePermission(permissionData models.CreatePermissionStruct) (*models.PermissionStruct, error) {
	database := database_connections.DatabaseConnectSQL()
	defer database.Close()

	if permissionData.GuardName == "" {
		permissionData.GuardName = "web"
	}

	query := `INSERT INTO ` + permissionsTable + ` (name, guard_name) VALUES (?, ?, ?)`
	result, err := database.Exec(query, permissionData.Name, permissionData.GuardName)
	if err != nil {
		return nil, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	return GetPermissionByID(int(id))
}

func UpdatePermission(id int, permissionData models.CreatePermissionStruct) (*models.PermissionStruct, error) {
	database := database_connections.DatabaseConnectSQL()
	defer database.Close()

	query := `UPDATE ` + permissionsTable + ` SET name = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?`
	_, err := database.Exec(query, permissionData.Name, id)
	if err != nil {
		return nil, err
	}

	return GetPermissionByID(id)
}

func DeletePermission(id int) error {
	database := database_connections.DatabaseConnectSQL()
	defer database.Close()

	query := `DELETE FROM ` + permissionsTable + ` WHERE id = ?`
	_, err := database.Exec(query, id)
	return err
}

func GetRolePermissions(roleID string) ([]models.PermissionStruct, error) {
	database := database_connections.DatabaseConnectSQL()
	defer database.Close()

	query := `
		SELECT p.uuid, p.name, p.guard_name, p.created_at, p.updated_at 
		FROM ` + permissionsTable + ` p
		INNER JOIN ` + rolePermissionsTable + ` rp ON p.uuid = rp.permission_id
		WHERE rp.role_id = ?
		ORDER BY p.name
	`
	rows, err := database.Query(query, roleID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var permissions []models.PermissionStruct
	for rows.Next() {
		var permission models.PermissionStruct
		err = rows.Scan(&permission.ID, &permission.Name, &permission.GuardName, &permission.CreatedAt, &permission.UpdatedAt)
		if err != nil {
			return nil, err
		}
		permissions = append(permissions, permission)
	}

	return permissions, nil
}

func GetUserDirectPermissions(userID string) ([]models.PermissionStruct, error) {
	database := database_connections.DatabaseConnectSQL()
	defer database.Close()

	query := `
		SELECT p.uuid, p.name, p.guard_name, p.created_at, p.updated_at 
		FROM ` + permissionsTable + ` p
		INNER JOIN ` + userPermissionsTable + ` up ON p.uuid = up.permission_id
		WHERE up.model_id = ? AND up.model_type = 'users'
		ORDER BY p.name
	`
	rows, err := database.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var permissions []models.PermissionStruct
	for rows.Next() {
		var permission models.PermissionStruct
		err = rows.Scan(&permission.ID, &permission.Name, &permission.GuardName, &permission.CreatedAt, &permission.UpdatedAt)
		if err != nil {
			return nil, err
		}
		permissions = append(permissions, permission)
	}

	return permissions, nil
}

func GetUserAllPermissions(userID string) ([]models.PermissionStruct, error) {
	database := database_connections.DatabaseConnectSQL()
	defer database.Close()

	query := `
		(
			SELECT DISTINCT p.uuid, p.name, p.guard_name, p.created_at, p.updated_at 
			FROM ` + permissionsTable + ` p
			INNER JOIN ` + userPermissionsTable + ` up ON p.uuid = up.permission_id
			WHERE up.model_id = ? AND up.model_type = 'users'
		)
		UNION
		(
			SELECT DISTINCT p.uuid, p.name, p.guard_name, p.created_at, p.updated_at 
			FROM ` + permissionsTable + ` p
			INNER JOIN ` + rolePermissionsTable + ` rp ON p.uuid = rp.permission_id
			INNER JOIN ` + userRolesTable + ` ur ON rp.role_id = ur.role_id
			WHERE ur.model_id = ? AND ur.model_type = 'users'
		)
		ORDER BY name
	`
	rows, err := database.Query(query, userID, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var permissions []models.PermissionStruct
	for rows.Next() {
		var permission models.PermissionStruct
		err = rows.Scan(&permission.ID, &permission.Name, &permission.GuardName, &permission.CreatedAt, &permission.UpdatedAt)
		if err != nil {
			return nil, err
		}
		permissions = append(permissions, permission)
	}

	return permissions, nil
}

func AssignPermissionToRole(roleID string, permissionID string) error {
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

func RevokePermissionFromRole(roleID string, permissionID string) error {
	database := database_connections.DatabaseConnectSQL()
	defer database.Close()

	query := `DELETE FROM ` + rolePermissionsTable + ` WHERE role_id = ? AND permission_id = ?`
	_, err := database.Exec(query, roleID, permissionID)
	return err
}

func AssignPermissionToUser(userID string, permissionID string) error {
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

	query := `INSERT INTO ` + userPermissionsTable + ` (model_id, permission_id, model_type) VALUES (?, ?, 'users')`
	_, err = database.Exec(query, userID, permissionID)
	return err
}

func RevokePermissionFromUser(userID string, permissionID string) error {
	database := database_connections.DatabaseConnectSQL()
	defer database.Close()

	query := `DELETE FROM ` + userPermissionsTable + ` WHERE model_id = ? AND permission_id = ? AND model_type = 'users'`
	_, err := database.Exec(query, userID, permissionID)
	return err
}

func RoleHasPermission(roleID string, permissionID string) (bool, error) {
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

func UserHasDirectPermission(userID string, permissionID string) (bool, error) {
	database := database_connections.DatabaseConnectSQL()
	defer database.Close()

	query := `SELECT COUNT(*) FROM ` + userPermissionsTable + ` WHERE model_id = ? AND permission_id = ? AND model_type = 'users'`
	var count int
	err := database.QueryRow(query, userID, permissionID).Scan(&count)
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func UserHasPermission(userID string, permissionName string, guardName string) (bool, error) {
	database := database_connections.DatabaseConnectSQL()
	defer database.Close()

	if guardName == "" {
		guardName = "web"
	}

	query := `
		SELECT COUNT(*) FROM (
			SELECT 1 
			FROM ` + userPermissionsTable + ` up
			INNER JOIN ` + permissionsTable + ` p ON up.permission_id = p.uuid
			WHERE up.model_id = ? AND up.model_type = 'users' AND p.name = ? AND p.guard_name = ?
			UNION
			SELECT 1 
			FROM ` + userRolesTable + ` ur
			INNER JOIN ` + rolePermissionsTable + ` rp ON ur.role_id = rp.role_id
			INNER JOIN ` + permissionsTable + ` p ON rp.permission_id = p.uuid
			WHERE ur.model_id = ? AND ur.model_type = 'users' AND p.name = ? AND p.guard_name = ?
		) AS combined_permissions
	`

	var count int
	err := database.QueryRow(query, userID, permissionName, guardName, userID, permissionName, guardName).Scan(&count)
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func UserHasAnyPermission(userID string, permissionNames []string, guardName string) (bool, error) {
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
				INNER JOIN %s p ON up.permission_id = p.uuid
				WHERE up.model_id = ? AND up.model_type = 'users' AND p.name IN (%s) AND p.guard_name = ?
			)
			UNION
			(
				SELECT 1 
				FROM %s ur
				INNER JOIN %s rp ON ur.role_id = rp.role_id
				INNER JOIN %s p ON rp.permission_id = p.uuid
				WHERE ur.model_id = ? AND ur.model_type = 'users' AND p.name IN (%s) AND p.guard_name = ?
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

func UserHasAllPermissions(userID string, permissionNames []string, guardName string) (bool, error) {
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

// ValidatePermissionsExist valida que todos los permisos en la lista existan en la base de datos
// Retorna los permisos válidos y una lista de permisos inexistentes
func ValidatePermissionsExist(permissionNames []string, guardName string) ([]models.PermissionStruct, []string, error) {
	if len(permissionNames) == 0 {
		return []models.PermissionStruct{}, []string{}, nil
	}

	database := database_connections.DatabaseConnectSQL()
	defer database.Close()

	if guardName == "" {
		guardName = "api"
	}

	validPermissions := []models.PermissionStruct{}
	invalidPermissions := []string{}

	for _, permissionName := range permissionNames {
		var permission models.PermissionStruct
		query := `SELECT uuid, name, guard_name, created_at, updated_at FROM ` + permissionsTable + ` WHERE name = ? AND guard_name = ?`
		err := database.QueryRow(query, permissionName, guardName).Scan(&permission.ID, &permission.Name, &permission.GuardName, &permission.CreatedAt, &permission.UpdatedAt)
		if err != nil {
			invalidPermissions = append(invalidPermissions, permissionName)
		} else {
			validPermissions = append(validPermissions, permission)
		}
	}

	return validPermissions, invalidPermissions, nil
}

// AssignMultiplePermissionsToRole asigna múltiples permisos a un rol
func AssignMultiplePermissionsToRole(roleID string, permissionIDs []string) error {
	if len(permissionIDs) == 0 {
		return nil
	}

	database := database_connections.DatabaseConnectSQL()
	defer database.Close()

	// Primero, limpiar los permisos actuales del rol
	deleteQuery := `DELETE FROM ` + rolePermissionsTable + ` WHERE role_id = ?`
	_, err := database.Exec(deleteQuery, roleID)
	if err != nil {
		return err
	}

	// Luego, insertar los nuevos permisos
	for _, permissionID := range permissionIDs {
		query := `INSERT INTO ` + rolePermissionsTable + ` (role_id, permission_id) VALUES (?, ?)`
		_, err := database.Exec(query, roleID, permissionID)
		if err != nil {
			return err
		}
	}

	return nil
}
