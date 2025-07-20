package middleware

import (
	"net/http"
	"semita/core/helpers"
	"semita/core/roles_and_permissions/models_roles_and_permissions"

	"github.com/gin-gonic/gin"
)

// getUserFromSession obtiene el usuario autenticado de la sesión
func getUserFromSession(c *gin.Context) (int, bool) {
	user, authenticated := helpers.GetAuthenticatedUser(c.Request)
	if !authenticated {
		return 0, false
	}
	return user.ID, true
}

// RequireRole middleware que verifica si el usuario tiene un rol específico
func RequireRole(roleName string, guardName ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Verificar si el usuario está autenticado
		if !helpers.IsUserAuthenticated(c.Request) {
			helpers.CreateFlashNotification(c.Writer, c.Request, "error", "You must be logged in to access this page.")
			c.Redirect(http.StatusSeeOther, "/auth/login")
			c.Abort()
			return
		}

		// Obtener el usuario autenticado de la sesión
		userID, authenticated := getUserFromSession(c)
		if !authenticated {
			helpers.CreateFlashNotification(c.Writer, c.Request, "error", "Session expired. Please log in again.")
			c.Redirect(http.StatusSeeOther, "/auth/login")
			c.Abort()
			return
		}

		// Determinar el guard name
		guard := "web"
		if len(guardName) > 0 && guardName[0] != "" {
			guard = guardName[0]
		}

		// Verificar si el usuario tiene el rol
		hasRole, err := models_roles_and_permissions.UserHasRoleByName(userID, roleName, guard)
		if err != nil {
			helpers.CreateFlashNotification(c.Writer, c.Request, "error", "Error checking user permissions.")
			c.Redirect(http.StatusSeeOther, "/")
			c.Abort()
			return
		}

		if !hasRole {
			helpers.CreateFlashNotification(c.Writer, c.Request, "error", "You don't have permission to access this page.")
			c.Redirect(http.StatusSeeOther, "/")
			c.Abort()
			return
		}

		c.Next()
	}
}

// RequireAnyRole middleware que verifica si el usuario tiene al menos uno de los roles especificados
func RequireAnyRole(roleNames []string, guardName ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Verificar si el usuario está autenticado
		if !helpers.IsUserAuthenticated(c.Request) {
			helpers.CreateFlashNotification(c.Writer, c.Request, "error", "You must be logged in to access this page.")
			c.Redirect(http.StatusSeeOther, "/auth/login")
			c.Abort()
			return
		}

		// Obtener el usuario autenticado de la sesión
		userID, authenticated := getUserFromSession(c)
		if !authenticated {
			helpers.CreateFlashNotification(c.Writer, c.Request, "error", "Session expired. Please log in again.")
			c.Redirect(http.StatusSeeOther, "/auth/login")
			c.Abort()
			return
		}

		// Determinar el guard name
		guard := "web"
		if len(guardName) > 0 && guardName[0] != "" {
			guard = guardName[0]
		}

		// Verificar si el usuario tiene al menos uno de los roles
		hasAnyRole, err := models_roles_and_permissions.UserHasAnyRole(userID, roleNames, guard)
		if err != nil {
			helpers.CreateFlashNotification(c.Writer, c.Request, "error", "Error checking user permissions.")
			c.Redirect(http.StatusSeeOther, "/")
			c.Abort()
			return
		}

		if !hasAnyRole {
			helpers.CreateFlashNotification(c.Writer, c.Request, "error", "You don't have permission to access this page.")
			c.Redirect(http.StatusSeeOther, "/")
			c.Abort()
			return
		}

		c.Next()
	}
}

// RequireAllRoles middleware que verifica si el usuario tiene todos los roles especificados
func RequireAllRoles(roleNames []string, guardName ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Verificar si el usuario está autenticado
		if !helpers.IsUserAuthenticated(c.Request) {
			helpers.CreateFlashNotification(c.Writer, c.Request, "error", "You must be logged in to access this page.")
			c.Redirect(http.StatusSeeOther, "/auth/login")
			c.Abort()
			return
		}

		// Obtener el usuario autenticado de la sesión
		userID, authenticated := getUserFromSession(c)
		if !authenticated {
			helpers.CreateFlashNotification(c.Writer, c.Request, "error", "Session expired. Please log in again.")
			c.Redirect(http.StatusSeeOther, "/auth/login")
			c.Abort()
			return
		}

		// Determinar el guard name
		guard := "web"
		if len(guardName) > 0 && guardName[0] != "" {
			guard = guardName[0]
		}

		// Verificar si el usuario tiene todos los roles
		hasAllRoles, err := models_roles_and_permissions.UserHasAllRoles(userID, roleNames, guard)
		if err != nil {
			helpers.CreateFlashNotification(c.Writer, c.Request, "error", "Error checking user permissions.")
			c.Redirect(http.StatusSeeOther, "/")
			c.Abort()
			return
		}

		if !hasAllRoles {
			helpers.CreateFlashNotification(c.Writer, c.Request, "error", "You don't have permission to access this page.")
			c.Redirect(http.StatusSeeOther, "/")
			c.Abort()
			return
		}

		c.Next()
	}
}

// RequirePermission middleware que verifica si el usuario tiene un permiso específico
func RequirePermission(permissionName string, guardName ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Verificar si el usuario está autenticado
		if !helpers.IsUserAuthenticated(c.Request) {
			helpers.CreateFlashNotification(c.Writer, c.Request, "error", "You must be logged in to access this page.")
			c.Redirect(http.StatusSeeOther, "/auth/login")
			c.Abort()
			return
		}

		// Obtener el usuario autenticado de la sesión
		userID, authenticated := getUserFromSession(c)
		if !authenticated {
			helpers.CreateFlashNotification(c.Writer, c.Request, "error", "Session expired. Please log in again.")
			c.Redirect(http.StatusSeeOther, "/auth/login")
			c.Abort()
			return
		}

		// Determinar el guard name
		guard := "web"
		if len(guardName) > 0 && guardName[0] != "" {
			guard = guardName[0]
		}

		// Verificar si el usuario tiene el permiso
		hasPermission, err := models_roles_and_permissions.UserHasPermission(userID, permissionName, guard)
		if err != nil {
			helpers.CreateFlashNotification(c.Writer, c.Request, "error", "Error checking user permissions.")
			c.Redirect(http.StatusSeeOther, "/")
			c.Abort()
			return
		}

		if !hasPermission {
			helpers.CreateFlashNotification(c.Writer, c.Request, "error", "You don't have permission to access this page.")
			c.Redirect(http.StatusSeeOther, "/")
			c.Abort()
			return
		}

		c.Next()
	}
}

// RequireAnyPermission middleware que verifica si el usuario tiene al menos uno de los permisos especificados
func RequireAnyPermission(permissionNames []string, guardName ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Verificar si el usuario está autenticado
		if !helpers.IsUserAuthenticated(c.Request) {
			helpers.CreateFlashNotification(c.Writer, c.Request, "error", "You must be logged in to access this page.")
			c.Redirect(http.StatusSeeOther, "/auth/login")
			c.Abort()
			return
		}

		// Obtener el usuario autenticado de la sesión
		userID, authenticated := getUserFromSession(c)
		if !authenticated {
			helpers.CreateFlashNotification(c.Writer, c.Request, "error", "Session expired. Please log in again.")
			c.Redirect(http.StatusSeeOther, "/auth/login")
			c.Abort()
			return
		}

		// Determinar el guard name
		guard := "web"
		if len(guardName) > 0 && guardName[0] != "" {
			guard = guardName[0]
		}

		// Verificar si el usuario tiene al menos uno de los permisos
		hasAnyPermission, err := models_roles_and_permissions.UserHasAnyPermission(userID, permissionNames, guard)
		if err != nil {
			helpers.CreateFlashNotification(c.Writer, c.Request, "error", "Error checking user permissions.")
			c.Redirect(http.StatusSeeOther, "/")
			c.Abort()
			return
		}

		if !hasAnyPermission {
			helpers.CreateFlashNotification(c.Writer, c.Request, "error", "You don't have permission to access this page.")
			c.Redirect(http.StatusSeeOther, "/")
			c.Abort()
			return
		}

		c.Next()
	}
}

// RequireAllPermissions middleware que verifica si el usuario tiene todos los permisos especificados
func RequireAllPermissions(permissionNames []string, guardName ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Verificar si el usuario está autenticado
		if !helpers.IsUserAuthenticated(c.Request) {
			helpers.CreateFlashNotification(c.Writer, c.Request, "error", "You must be logged in to access this page.")
			c.Redirect(http.StatusSeeOther, "/auth/login")
			c.Abort()
			return
		}

		// Obtener el usuario autenticado de la sesión
		userID, authenticated := getUserFromSession(c)
		if !authenticated {
			helpers.CreateFlashNotification(c.Writer, c.Request, "error", "Session expired. Please log in again.")
			c.Redirect(http.StatusSeeOther, "/auth/login")
			c.Abort()
			return
		}

		// Determinar el guard name
		guard := "web"
		if len(guardName) > 0 && guardName[0] != "" {
			guard = guardName[0]
		}

		// Verificar si el usuario tiene todos los permisos
		hasAllPermissions, err := models_roles_and_permissions.UserHasAllPermissions(userID, permissionNames, guard)
		if err != nil {
			helpers.CreateFlashNotification(c.Writer, c.Request, "error", "Error checking user permissions.")
			c.Redirect(http.StatusSeeOther, "/")
			c.Abort()
			return
		}

		if !hasAllPermissions {
			helpers.CreateFlashNotification(c.Writer, c.Request, "error", "You don't have permission to access this page.")
			c.Redirect(http.StatusSeeOther, "/")
			c.Abort()
			return
		}

		c.Next()
	}
}

// CheckRoleOrPermission middleware que verifica si el usuario tiene un rol O un permiso específico
func CheckRoleOrPermission(roleName string, permissionName string, guardName ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Verificar si el usuario está autenticado
		if !helpers.IsUserAuthenticated(c.Request) {
			helpers.CreateFlashNotification(c.Writer, c.Request, "error", "You must be logged in to access this page.")
			c.Redirect(http.StatusSeeOther, "/auth/login")
			c.Abort()
			return
		}

		// Obtener el usuario autenticado de la sesión
		userID, authenticated := getUserFromSession(c)
		if !authenticated {
			helpers.CreateFlashNotification(c.Writer, c.Request, "error", "Session expired. Please log in again.")
			c.Redirect(http.StatusSeeOther, "/auth/login")
			c.Abort()
			return
		}

		// Determinar el guard name
		guard := "web"
		if len(guardName) > 0 && guardName[0] != "" {
			guard = guardName[0]
		}

		// Verificar si el usuario tiene el rol o el permiso
		hasRole, err := models_roles_and_permissions.UserHasRoleByName(userID, roleName, guard)
		if err != nil {
			helpers.CreateFlashNotification(c.Writer, c.Request, "error", "Error checking user permissions.")
			c.Redirect(http.StatusSeeOther, "/")
			c.Abort()
			return
		}

		if hasRole {
			c.Next()
			return
		}

		hasPermission, err := models_roles_and_permissions.UserHasPermission(userID, permissionName, guard)
		if err != nil {
			helpers.CreateFlashNotification(c.Writer, c.Request, "error", "Error checking user permissions.")
			c.Redirect(http.StatusSeeOther, "/")
			c.Abort()
			return
		}

		if !hasPermission {
			helpers.CreateFlashNotification(c.Writer, c.Request, "error", "You don't have permission to access this page.")
			c.Redirect(http.StatusSeeOther, "/")
			c.Abort()
			return
		}

		c.Next()
	}
}
