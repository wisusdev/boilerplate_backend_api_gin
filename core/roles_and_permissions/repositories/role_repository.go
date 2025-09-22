package repositories

import (
	"fmt"
	"semita/core/database/database_connections"
	"semita/core/roles_and_permissions/models"
	"strings"
)

var rolesTable = "roles"
var userRolesTable = "model_has_roles"

func GetAllRoles() ([]models.RoleStruct, error) {
	database := database_connections.DatabaseConnectSQL()
	defer database.Close()

	query := `SELECT id, name, guard_name, created_at, updated_at FROM ` + rolesTable + ` ORDER BY name`
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

func GetRoleByID(id int) (*models.RoleStruct, error) {
	database := database_connections.DatabaseConnectSQL()
	defer database.Close()

	query := `SELECT id, name, guard_name, created_at, updated_at FROM ` + rolesTable + ` WHERE id = ?`
	row := database.QueryRow(query, id)

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
		roleData.GuardName = "web"
	}

	query := `INSERT INTO ` + rolesTable + ` (name, guard_name) VALUES (?, ?, ?)`
	result, err := database.Exec(query, roleData.Name, roleData.GuardName)
	if err != nil {
		return nil, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	return GetRoleByID(int(id))
}

func UpdateRole(id int, roleData models.CreateRoleStruct) (*models.RoleStruct, error) {
	database := database_connections.DatabaseConnectSQL()
	defer database.Close()

	query := `UPDATE ` + rolesTable + ` SET name = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?`
	_, err := database.Exec(query, roleData.Name, id)
	if err != nil {
		return nil, err
	}

	return GetRoleByID(id)
}

func DeleteRole(id int) error {
	database := database_connections.DatabaseConnectSQL()
	defer database.Close()

	query := `DELETE FROM ` + rolesTable + ` WHERE id = ?`
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
