package middleware

import (
	"net/http"
	"semita/core/helpers"
	"semita/core/roles_and_permissions/providers"

	"github.com/gin-gonic/gin"
)

// getUserFromSession obtiene el usuario autenticado usando el sistema unificado
// Funciona tanto para API (JWT) como para web (sesiones)
func getUserFromSession(c *gin.Context) (string, bool) {
	user, authenticated := helpers.GetAppAuthenticatedUser(c)
	if !authenticated {
		return "", false
	}
	return user.ID, true
}

// isAPIRequest determina si la petición es una API o una web request
func isAPIRequest(c *gin.Context) bool {
	// Si tiene el header Authorization con Bearer token, es API
	authHeader := c.GetHeader("Authorization")
	if authHeader != "" && len(authHeader) > 7 && authHeader[:7] == "Bearer " {
		return true
	}

	// Si el content-type o accept es application/vnd.api+json, es API
	contentType := c.GetHeader("Content-Type")
	accept := c.GetHeader("Accept")

	return contentType == "application/vnd.api+json" || accept == "application/vnd.api+json"
}

// getGuardName determina el guard correcto según el tipo de petición
func getGuardName(c *gin.Context, guardName ...string) string {
	// Si se especifica un guard explícitamente, usarlo
	if len(guardName) > 0 && guardName[0] != "" {
		return guardName[0]
	}

	// Determinar automáticamente basado en el tipo de petición
	if isAPIRequest(c) {
		return "api"
	}

	return "web"
}

// handleUnauthorized maneja la respuesta cuando el usuario no está autorizado
func handleUnauthorized(c *gin.Context, message string) {
	if isAPIRequest(c) {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": message,
		})
	} else {
		helpers.CreateFlashNotification(c.Writer, c.Request, "error", message)
		c.Redirect(http.StatusSeeOther, "/auth/login")
	}
	c.Abort()
}

// handleForbidden maneja la respuesta cuando el usuario no tiene permisos
func handleForbidden(c *gin.Context, message string) {
	if isAPIRequest(c) {
		c.JSON(http.StatusForbidden, gin.H{
			"error": message,
		})
	} else {
		helpers.CreateFlashNotification(c.Writer, c.Request, "error", message)
		c.Redirect(http.StatusSeeOther, "/")
	}
	c.Abort()
}

// RequireRole middleware que verifica si el usuario tiene un rol específico
func RequireRole(roleName string, guardName ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Obtener el usuario autenticado (funciona para API y web)
		userID, authenticated := getUserFromSession(c)
		if !authenticated {
			handleUnauthorized(c, "You must be logged in to access this resource.")
			return
		}

		// Determinar el guard name dinámicamente
		guard := getGuardName(c, guardName...)

		// Verificar si el usuario tiene el rol
		hasRole, err := providers.UserHasRoleByName(userID, roleName, guard)
		if err != nil {
			handleForbidden(c, "Error checking user permissions.")
			return
		}

		if !hasRole {
			handleForbidden(c, "You don't have the required role to access this resource.")
			return
		}

		c.Next()
	}
}

// RequireAnyRole middleware que verifica si el usuario tiene al menos uno de los roles especificados
func RequireAnyRole(roleNames []string, guardName ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Obtener el usuario autenticado (funciona para API y web)
		userID, authenticated := getUserFromSession(c)
		if !authenticated {
			handleUnauthorized(c, "You must be logged in to access this resource.")
			return
		}

		// Determinar el guard name dinámicamente
		guard := getGuardName(c, guardName...)

		// Verificar si el usuario tiene al menos uno de los roles
		hasAnyRole, err := providers.UserHasAnyRole(userID, roleNames, guard)
		if err != nil {
			handleForbidden(c, "Error checking user permissions.")
			return
		}

		if !hasAnyRole {
			handleForbidden(c, "You don't have any of the required roles to access this resource.")
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
		hasAllRoles, err := providers.UserHasAllRoles(userID, roleNames, guard)
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
		// Obtener el usuario autenticado (funciona para API y web)
		userID, authenticated := getUserFromSession(c)
		if !authenticated {
			handleUnauthorized(c, "You must be logged in to access this resource.")
			return
		}

		// Determinar el guard name dinámicamente
		guard := getGuardName(c, guardName...)

		// Verificar si el usuario tiene el permiso
		hasPermission, err := providers.UserHasPermission(userID, permissionName, guard)
		if err != nil {
			handleForbidden(c, "Error checking user permissions.")
			return
		}

		if !hasPermission {
			handleForbidden(c, "You don't have permission to access this resource.")
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
		hasAnyPermission, err := providers.UserHasAnyPermission(userID, permissionNames, guard)
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
		hasAllPermissions, err := providers.UserHasAllPermissions(userID, permissionNames, guard)
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
		hasRole, err := providers.UserHasRoleByName(userID, roleName, guard)
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

		hasPermission, err := providers.UserHasPermission(userID, permissionName, guard)
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

// === API MIDDLEWARE FUNCTIONS ===

// RequirePermissionAPI middleware que verifica permisos para endpoints API usando JWT
func RequirePermissionAPI(permissionName string, guardName ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Verificar si el usuario está autenticado vía JWT (debe haber user_id en el contexto)
		userID, exists := c.Get("user_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{
				"status":  "error",
				"message": "Authentication required",
			})
			c.Abort()
			return
		}

		userIDStr, ok := userID.(string)
		if !ok || userIDStr == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"status":  "error",
				"message": "Invalid user authentication",
			})
			c.Abort()
			return
		}

		// Determinar el guard name (por defecto "api" para endpoints API)
		guard := "api"
		if len(guardName) > 0 && guardName[0] != "" {
			guard = guardName[0]
		}

		// Verificar si el usuario tiene el permiso
		hasPermission, err := providers.UserHasPermission(userIDStr, permissionName, guard)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  "error",
				"message": "Error checking user permissions",
			})
			c.Abort()
			return
		}

		if !hasPermission {
			c.JSON(http.StatusForbidden, gin.H{
				"status":  "error",
				"message": "You dont have permission to access this resource",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}
