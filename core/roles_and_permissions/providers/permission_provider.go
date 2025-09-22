package providers

import (
	"boilerplate_backend_api_gin/core/roles_and_permissions/models"
	"boilerplate_backend_api_gin/core/roles_and_permissions/repositories"
)

// GetAllPermissions obtiene todos los permisos
func GetAllPermissions() ([]models.PermissionStruct, error) {
	return repositories.GetAllPermissions()
}

// GetPermissionByID obtiene un permiso por su ID
func GetPermissionByID(id int) (*models.PermissionStruct, error) {
	return repositories.GetPermissionByID(id)
}

// GetPermissionByName obtiene un permiso por su nombre
func GetPermissionByName(name string, guardName string) (*models.PermissionStruct, error) {
	return repositories.GetPermissionByName(name, guardName)
}

// CreatePermission crea un nuevo permiso
func CreatePermission(permissionData models.CreatePermissionStruct) (*models.PermissionStruct, error) {
	return repositories.CreatePermission(permissionData)
}

// UpdatePermission actualiza un permiso existente
func UpdatePermission(id int, permissionData models.CreatePermissionStruct) (*models.PermissionStruct, error) {
	return repositories.UpdatePermission(id, permissionData)
}

// DeletePermission elimina un permiso
func DeletePermission(id int) error {
	return repositories.DeletePermission(id)
}

// GetRolePermissions obtiene todos los permisos de un rol
func GetRolePermissions(roleID int) ([]models.PermissionStruct, error) {
	return repositories.GetRolePermissions(roleID)
}

// GetUserDirectPermissions obtiene los permisos directos de un usuario (no heredados de roles)
func GetUserDirectPermissions(userID string) ([]models.PermissionStruct, error) {
	return repositories.GetUserDirectPermissions(userID)
}

// GetUserAllPermissions obtiene todos los permisos de un usuario (directos + heredados de roles)
func GetUserAllPermissions(userID string) ([]models.PermissionStruct, error) {
	return repositories.GetUserAllPermissions(userID)
}

// AssignPermissionToRole asigna un permiso a un rol
func AssignPermissionToRole(roleID string, permissionID string) error {
	return repositories.AssignPermissionToRole(roleID, permissionID)
}

// RevokePermissionFromRole revoca un permiso de un rol
func RevokePermissionFromRole(roleID string, permissionID string) error {
	return repositories.RevokePermissionFromRole(roleID, permissionID)
}

// AssignPermissionToUser asigna un permiso directamente a un usuario
func AssignPermissionToUser(userID string, permissionID string) error {
	return repositories.AssignPermissionToUser(userID, permissionID)
}

// RevokePermissionFromUser revoca un permiso directo de un usuario
func RevokePermissionFromUser(userID string, permissionID string) error {
	return repositories.RevokePermissionFromUser(userID, permissionID)
}

// RoleHasPermission verifica si un rol tiene un permiso específico
func RoleHasPermission(roleID string, permissionID string) (bool, error) {
	return repositories.RoleHasPermission(roleID, permissionID)
}

// UserHasDirectPermission verifica si un usuario tiene un permiso directo
func UserHasDirectPermission(userID string, permissionID string) (bool, error) {
	return repositories.UserHasDirectPermission(userID, permissionID)
}

// UserHasPermission verifica si un usuario tiene un permiso (directo o heredado)
func UserHasPermission(userID string, permissionName string, guardName string) (bool, error) {
	return repositories.UserHasPermission(userID, permissionName, guardName)
}

// UserHasAnyPermission verifica si un usuario tiene al menos uno de los permisos especificados
func UserHasAnyPermission(userID string, permissionNames []string, guardName string) (bool, error) {
	return repositories.UserHasAnyPermission(userID, permissionNames, guardName)
}

// UserHasAllPermissions verifica si un usuario tiene todos los permisos especificados
func UserHasAllPermissions(userID string, permissionNames []string, guardName string) (bool, error) {
	return repositories.UserHasAllPermissions(userID, permissionNames, guardName)
}
