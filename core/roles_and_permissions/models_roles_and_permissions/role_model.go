package models_roles_and_permissions

import (
	"fmt"
	"semita/app/data/structs"
	"semita/common/nulltypes"
	"semita/core/database/database_connections"
	"strings"
)

var rolesTable = "roles"
var userRolesTable = "user_roles"

// GetAllRoles obtiene todos los roles
func GetAllRoles() ([]structs.RoleStruct, error) {
	database := database_connections.DatabaseConnectSQL()
	defer database.Close()

	query := `SELECT id, name, guard_name, description, created_at, updated_at FROM ` + rolesTable + ` ORDER BY name`
	rows, err := database.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var roles []structs.RoleStruct
	for rows.Next() {
		var role structs.RoleStruct
		var description nulltypes.NullString
		err = rows.Scan(&role.ID, &role.Name, &role.GuardName, &description, &role.CreatedAt, &role.UpdatedAt)
		if err != nil {
			return nil, err
		}
		role.Description = description.String
		roles = append(roles, role)
	}

	return roles, nil
}

// GetRoleByID obtiene un rol por su ID
func GetRoleByID(id int) (*structs.RoleStruct, error) {
	database := database_connections.DatabaseConnectSQL()
	defer database.Close()

	query := `SELECT id, name, guard_name, description, created_at, updated_at FROM ` + rolesTable + ` WHERE id = ?`
	row := database.QueryRow(query, id)

	var role structs.RoleStruct
	var description nulltypes.NullString
	err := row.Scan(&role.ID, &role.Name, &role.GuardName, &description, &role.CreatedAt, &role.UpdatedAt)
	if err != nil {
		return nil, err
	}
	role.Description = description.String

	return &role, nil
}

// GetRoleByName obtiene un rol por su nombre
func GetRoleByName(name string, guardName string) (*structs.RoleStruct, error) {
	database := database_connections.DatabaseConnectSQL()
	defer database.Close()

	query := `SELECT id, name, guard_name, description, created_at, updated_at FROM ` + rolesTable + ` WHERE name = ? AND guard_name = ?`
	row := database.QueryRow(query, name, guardName)

	var role structs.RoleStruct
	var description nulltypes.NullString
	err := row.Scan(&role.ID, &role.Name, &role.GuardName, &description, &role.CreatedAt, &role.UpdatedAt)
	if err != nil {
		return nil, err
	}
	role.Description = description.String

	return &role, nil
}

// CreateRole crea un nuevo rol
func CreateRole(roleData structs.CreateRoleStruct) (*structs.RoleStruct, error) {
	database := database_connections.DatabaseConnectSQL()
	defer database.Close()

	if roleData.GuardName == "" {
		roleData.GuardName = "web"
	}

	query := `INSERT INTO ` + rolesTable + ` (name, guard_name, description) VALUES (?, ?, ?)`
	result, err := database.Exec(query, roleData.Name, roleData.GuardName, roleData.Description)
	if err != nil {
		return nil, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	return GetRoleByID(int(id))
}

// UpdateRole actualiza un rol existente
func UpdateRole(id int, roleData structs.CreateRoleStruct) (*structs.RoleStruct, error) {
	database := database_connections.DatabaseConnectSQL()
	defer database.Close()

	query := `UPDATE ` + rolesTable + ` SET name = ?, description = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?`
	_, err := database.Exec(query, roleData.Name, roleData.Description, id)
	if err != nil {
		return nil, err
	}

	return GetRoleByID(id)
}

// DeleteRole elimina un rol
func DeleteRole(id int) error {
	database := database_connections.DatabaseConnectSQL()
	defer database.Close()

	query := `DELETE FROM ` + rolesTable + ` WHERE id = ?`
	_, err := database.Exec(query, id)
	return err
}

// GetUserRoles obtiene todos los roles de un usuario
func GetUserRoles(userID int) ([]structs.RoleStruct, error) {
	database := database_connections.DatabaseConnectSQL()
	defer database.Close()

	query := `
		SELECT r.id, r.name, r.guard_name, r.description, r.created_at, r.updated_at 
		FROM ` + rolesTable + ` r
		INNER JOIN ` + userRolesTable + ` ur ON r.id = ur.role_id
		WHERE ur.user_id = ?
		ORDER BY r.name
	`
	rows, err := database.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var roles []structs.RoleStruct
	for rows.Next() {
		var role structs.RoleStruct
		var description nulltypes.NullString
		err = rows.Scan(&role.ID, &role.Name, &role.GuardName, &description, &role.CreatedAt, &role.UpdatedAt)
		if err != nil {
			return nil, err
		}
		role.Description = description.String
		roles = append(roles, role)
	}

	return roles, nil
}

// AssignRoleToUser asigna un rol a un usuario
func AssignRoleToUser(userID int, roleID int) error {
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

	query := `INSERT INTO ` + userRolesTable + ` (user_id, role_id) VALUES (?, ?)`
	_, err = database.Exec(query, userID, roleID)
	return err
}

// RevokeRoleFromUser revoca un rol de un usuario
func RevokeRoleFromUser(userID int, roleID int) error {
	database := database_connections.DatabaseConnectSQL()
	defer database.Close()

	query := `DELETE FROM ` + userRolesTable + ` WHERE user_id = ? AND role_id = ?`
	_, err := database.Exec(query, userID, roleID)
	return err
}

// UserHasRole verifica si un usuario tiene un rol específico
func UserHasRole(userID int, roleID int) (bool, error) {
	database := database_connections.DatabaseConnectSQL()
	defer database.Close()

	query := `SELECT COUNT(*) FROM ` + userRolesTable + ` WHERE user_id = ? AND role_id = ?`
	var count int
	err := database.QueryRow(query, userID, roleID).Scan(&count)
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

// UserHasRoleByName verifica si un usuario tiene un rol por nombre
func UserHasRoleByName(userID int, roleName string, guardName string) (bool, error) {
	database := database_connections.DatabaseConnectSQL()
	defer database.Close()

	if guardName == "" {
		guardName = "web"
	}

	query := `
		SELECT COUNT(*) 
		FROM ` + userRolesTable + ` ur
		INNER JOIN ` + rolesTable + ` r ON ur.role_id = r.id
		WHERE ur.user_id = ? AND r.name = ? AND r.guard_name = ?
	`
	var count int
	err := database.QueryRow(query, userID, roleName, guardName).Scan(&count)
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

// UserHasAnyRole verifica si un usuario tiene al menos uno de los roles especificados
func UserHasAnyRole(userID int, roleNames []string, guardName string) (bool, error) {
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
		INNER JOIN %s r ON ur.role_id = r.id
		WHERE ur.user_id = ? AND r.name IN (%s) AND r.guard_name = ?
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

// UserHasAllRoles verifica si un usuario tiene todos los roles especificados
func UserHasAllRoles(userID int, roleNames []string, guardName string) (bool, error) {
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
		INNER JOIN %s r ON ur.role_id = r.id
		WHERE ur.user_id = ? AND r.name IN (%s) AND r.guard_name = ?
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
