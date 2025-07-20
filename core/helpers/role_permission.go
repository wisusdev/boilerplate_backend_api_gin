package helpers

import (
	"net/http"
	"semita/core/roles_and_permissions/models_roles_and_permissions"
)

// RolePermissionHelper proporciona métodos auxiliares para verificar roles y permisos
type RolePermissionHelper struct{}

// NewRolePermissionHelper crea una nueva instancia del helper
func NewRolePermissionHelper() *RolePermissionHelper {
	return &RolePermissionHelper{}
}

// getUserIDFromRequest obtiene el ID del usuario autenticado desde la request
func (rph *RolePermissionHelper) getUserIDFromRequest(request *http.Request) (int, bool) {
	user, authenticated := GetAuthenticatedUser(request)
	if !authenticated {
		return 0, false
	}
	return user.ID, true
}

// HasRole verifica si el usuario autenticado tiene un rol específico
func (rph *RolePermissionHelper) HasRole(request *http.Request, roleName string, guardName ...string) bool {
	userID, authenticated := rph.getUserIDFromRequest(request)
	if !authenticated {
		return false
	}

	guard := "web"
	if len(guardName) > 0 && guardName[0] != "" {
		guard = guardName[0]
	}

	hasRole, err := models_roles_and_permissions.UserHasRoleByName(userID, roleName, guard)
	if err != nil {
		return false
	}

	return hasRole
}

// HasAnyRole verifica si el usuario autenticado tiene al menos uno de los roles especificados
func (rph *RolePermissionHelper) HasAnyRole(request *http.Request, roleNames []string, guardName ...string) bool {
	userID, authenticated := rph.getUserIDFromRequest(request)
	if !authenticated {
		return false
	}

	guard := "web"
	if len(guardName) > 0 && guardName[0] != "" {
		guard = guardName[0]
	}

	hasAnyRole, err := models_roles_and_permissions.UserHasAnyRole(userID, roleNames, guard)
	if err != nil {
		return false
	}

	return hasAnyRole
}

// HasAllRoles verifica si el usuario autenticado tiene todos los roles especificados
func (rph *RolePermissionHelper) HasAllRoles(request *http.Request, roleNames []string, guardName ...string) bool {
	userID, authenticated := rph.getUserIDFromRequest(request)
	if !authenticated {
		return false
	}

	guard := "web"
	if len(guardName) > 0 && guardName[0] != "" {
		guard = guardName[0]
	}

	hasAllRoles, err := models_roles_and_permissions.UserHasAllRoles(userID, roleNames, guard)
	if err != nil {
		return false
	}

	return hasAllRoles
}

// HasPermission verifica si el usuario autenticado tiene un permiso específico
func (rph *RolePermissionHelper) HasPermission(request *http.Request, permissionName string, guardName ...string) bool {
	userID, authenticated := rph.getUserIDFromRequest(request)
	if !authenticated {
		return false
	}

	guard := "web"
	if len(guardName) > 0 && guardName[0] != "" {
		guard = guardName[0]
	}

	hasPermission, err := models_roles_and_permissions.UserHasPermission(userID, permissionName, guard)
	if err != nil {
		return false
	}

	return hasPermission
}

// HasAnyPermission verifica si el usuario autenticado tiene al menos uno de los permisos especificados
func (rph *RolePermissionHelper) HasAnyPermission(request *http.Request, permissionNames []string, guardName ...string) bool {
	userID, authenticated := rph.getUserIDFromRequest(request)
	if !authenticated {
		return false
	}

	guard := "web"
	if len(guardName) > 0 && guardName[0] != "" {
		guard = guardName[0]
	}

	hasAnyPermission, err := models_roles_and_permissions.UserHasAnyPermission(userID, permissionNames, guard)
	if err != nil {
		return false
	}

	return hasAnyPermission
}

// HasAllPermissions verifica si el usuario autenticado tiene todos los permisos especificados
func (rph *RolePermissionHelper) HasAllPermissions(request *http.Request, permissionNames []string, guardName ...string) bool {
	userID, authenticated := rph.getUserIDFromRequest(request)
	if !authenticated {
		return false
	}

	guard := "web"
	if len(guardName) > 0 && guardName[0] != "" {
		guard = guardName[0]
	}

	hasAllPermissions, err := models_roles_and_permissions.UserHasAllPermissions(userID, permissionNames, guard)
	if err != nil {
		return false
	}

	return hasAllPermissions
}

// GetUserRoles obtiene todos los roles del usuario autenticado
func (rph *RolePermissionHelper) GetUserRoles(request *http.Request) ([]string, bool) {
	userID, authenticated := rph.getUserIDFromRequest(request)
	if !authenticated {
		return nil, false
	}

	roles, err := models_roles_and_permissions.GetUserRoles(userID)
	if err != nil {
		return nil, false
	}

	roleNames := make([]string, len(roles))
	for i, role := range roles {
		roleNames[i] = role.Name
	}

	return roleNames, true
}

// GetUserPermissions obtiene todos los permisos del usuario autenticado
func (rph *RolePermissionHelper) GetUserPermissions(request *http.Request) ([]string, bool) {
	userID, authenticated := rph.getUserIDFromRequest(request)
	if !authenticated {
		return nil, false
	}

	permissions, err := models_roles_and_permissions.GetUserAllPermissions(userID)
	if err != nil {
		return nil, false
	}

	permissionNames := make([]string, len(permissions))
	for i, permission := range permissions {
		permissionNames[i] = permission.Name
	}

	return permissionNames, true
}

// IsUserAdmin verifica si el usuario es administrador (tiene rol admin o super-admin)
func (rph *RolePermissionHelper) IsUserAdmin(request *http.Request) bool {
	return rph.HasAnyRole(request, []string{"admin", "super-admin"})
}

// IsUserSuperAdmin verifica si el usuario es super administrador
func (rph *RolePermissionHelper) IsUserSuperAdmin(request *http.Request) bool {
	return rph.HasRole(request, "super-admin")
}

// CanManageUsers verifica si el usuario puede gestionar usuarios
func (rph *RolePermissionHelper) CanManageUsers(request *http.Request) bool {
	return rph.HasAnyPermission(request, []string{"create-users", "edit-users", "delete-users"})
}

// CanManageRoles verifica si el usuario puede gestionar roles
func (rph *RolePermissionHelper) CanManageRoles(request *http.Request) bool {
	return rph.HasAnyPermission(request, []string{"create-roles", "edit-roles", "delete-roles", "assign-roles"})
}

// CanManagePermissions verifica si el usuario puede gestionar permisos
func (rph *RolePermissionHelper) CanManagePermissions(request *http.Request) bool {
	return rph.HasAnyPermission(request, []string{"create-permissions", "edit-permissions", "delete-permissions", "assign-permissions"})
}

// CanAccessDashboard verifica si el usuario puede acceder al dashboard
func (rph *RolePermissionHelper) CanAccessDashboard(request *http.Request) bool {
	return rph.HasPermission(request, "view-dashboard")
}

// Instancia global del helper para uso fácil
var RolePermissionHelperInstance = NewRolePermissionHelper()

// Funciones auxiliares globales para uso directo
func HasRole(request *http.Request, roleName string, guardName ...string) bool {
	return RolePermissionHelperInstance.HasRole(request, roleName, guardName...)
}

func HasAnyRole(request *http.Request, roleNames []string, guardName ...string) bool {
	return RolePermissionHelperInstance.HasAnyRole(request, roleNames, guardName...)
}

func HasAllRoles(request *http.Request, roleNames []string, guardName ...string) bool {
	return RolePermissionHelperInstance.HasAllRoles(request, roleNames, guardName...)
}

func HasPermission(request *http.Request, permissionName string, guardName ...string) bool {
	return RolePermissionHelperInstance.HasPermission(request, permissionName, guardName...)
}

func HasAnyPermission(request *http.Request, permissionNames []string, guardName ...string) bool {
	return RolePermissionHelperInstance.HasAnyPermission(request, permissionNames, guardName...)
}

func HasAllPermissions(request *http.Request, permissionNames []string, guardName ...string) bool {
	return RolePermissionHelperInstance.HasAllPermissions(request, permissionNames, guardName...)
}

func GetUserRoles(request *http.Request) ([]string, bool) {
	return RolePermissionHelperInstance.GetUserRoles(request)
}

func GetUserPermissions(request *http.Request) ([]string, bool) {
	return RolePermissionHelperInstance.GetUserPermissions(request)
}

func IsUserAdmin(request *http.Request) bool {
	return RolePermissionHelperInstance.IsUserAdmin(request)
}

func IsUserSuperAdmin(request *http.Request) bool {
	return RolePermissionHelperInstance.IsUserSuperAdmin(request)
}

func CanManageUsers(request *http.Request) bool {
	return RolePermissionHelperInstance.CanManageUsers(request)
}

func CanManageRoles(request *http.Request) bool {
	return RolePermissionHelperInstance.CanManageRoles(request)
}

func CanManagePermissions(request *http.Request) bool {
	return RolePermissionHelperInstance.CanManagePermissions(request)
}

func CanAccessDashboard(request *http.Request) bool {
	return RolePermissionHelperInstance.CanAccessDashboard(request)
}
