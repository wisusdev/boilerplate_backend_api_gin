package base

import (
	"net/http"
	"semita/app/data/models"
	"semita/app/data/structs"
	"semita/core/helpers"
	"semita/core/roles_and_permissions/models_roles_and_permissions"
	"strconv"

	"github.com/gin-gonic/gin"
)

// UserPermissionController maneja verificaciones y gestión de permisos de usuarios
type UserPermissionController struct{}

// CheckUserPermissions verifica los permisos de un usuario específico
func (upc *UserPermissionController) CheckUserPermissions(c *gin.Context) {
	userIDParam := c.Param("user_id")
	userID, err := strconv.Atoi(userIDParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Invalid user ID",
		})
		return
	}

	// Obtener roles del usuario
	roles, err := models_roles_and_permissions.GetUserRoles(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Error retrieving user roles: " + err.Error(),
		})
		return
	}

	// Obtener permisos directos del usuario
	directPermissions, err := models_roles_and_permissions.GetUserDirectPermissions(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Error retrieving user direct permissions: " + err.Error(),
		})
		return
	}

	// Obtener todos los permisos del usuario (directos + heredados)
	allPermissions, err := models_roles_and_permissions.GetUserAllPermissions(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Error retrieving user permissions: " + err.Error(),
		})
		return
	}

	// Obtener información del usuario
	user, err := models.GetUserByID(strconv.Itoa(userID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  "error",
			"message": "User not found",
		})
		return
	}

	userWithPerms := structs.UserWithRolesAndPermissions{
		UserStruct:        user,
		Roles:             roles,
		DirectPermissions: directPermissions,
		AllPermissions:    allPermissions,
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   userWithPerms,
	})
}

// CheckCurrentUserPermissions verifica los permisos del usuario actualmente logueado
func (upc *UserPermissionController) CheckCurrentUserPermissions(c *gin.Context) {
	// Obtener el usuario autenticado de la sesión
	user, authenticated := helpers.GetAuthenticatedUser(c.Request)
	if !authenticated {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status":  "error",
			"message": "User not authenticated",
		})
		return
	}

	userID := user.ID

	// Obtener roles del usuario
	roles, err := models_roles_and_permissions.GetUserRoles(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Error retrieving user roles: " + err.Error(),
		})
		return
	}

	// Obtener permisos directos del usuario
	directPermissions, err := models_roles_and_permissions.GetUserDirectPermissions(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Error retrieving user direct permissions: " + err.Error(),
		})
		return
	}

	// Obtener todos los permisos del usuario (directos + heredados)
	allPermissions, err := models_roles_and_permissions.GetUserAllPermissions(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Error retrieving user permissions: " + err.Error(),
		})
		return
	}

	var userData = structs.UserStruct{
		ID:        user.ID,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Username:  user.Username,
		Avatar:    helpers.NullStringToString(user.Avatar),
		Language:  helpers.NullStringToString(user.Language),
	}

	userWithPerms := structs.UserWithRolesAndPermissions{
		UserStruct:        userData,
		Roles:             roles,
		DirectPermissions: directPermissions,
		AllPermissions:    allPermissions,
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   userWithPerms,
	})
}

// CheckRole verifica si un usuario tiene un rol específico
func (upc *UserPermissionController) CheckRole(c *gin.Context) {
	userIDParam := c.Param("user_id")
	userID, err := strconv.Atoi(userIDParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Invalid user ID",
		})
		return
	}

	roleName := c.Query("role")
	guardName := c.Query("guard")
	if guardName == "" {
		guardName = "web"
	}

	if roleName == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Role name is required",
		})
		return
	}

	hasRole, err := models_roles_and_permissions.UserHasRoleByName(userID, roleName, guardName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Error checking user role: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data": gin.H{
			"user_id":  userID,
			"role":     roleName,
			"guard":    guardName,
			"has_role": hasRole,
		},
	})
}

// CheckPermission verifica si un usuario tiene un permiso específico
func (upc *UserPermissionController) CheckPermission(c *gin.Context) {
	userIDParam := c.Param("user_id")
	userID, err := strconv.Atoi(userIDParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Invalid user ID",
		})
		return
	}

	permissionName := c.Query("permission")
	guardName := c.Query("guard")
	if guardName == "" {
		guardName = "web"
	}

	if permissionName == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Permission name is required",
		})
		return
	}

	hasPermission, err := models_roles_and_permissions.UserHasPermission(userID, permissionName, guardName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Error checking user permission: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data": gin.H{
			"user_id":        userID,
			"permission":     permissionName,
			"guard":          guardName,
			"has_permission": hasPermission,
		},
	})
}

// CheckCurrentUserRole verifica si el usuario logueado tiene un rol específico
func (upc *UserPermissionController) CheckCurrentUserRole(c *gin.Context) {
	// Obtener el usuario autenticado de la sesión
	user, authenticated := helpers.GetAuthenticatedUser(c.Request)
	if !authenticated {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status":  "error",
			"message": "User not authenticated",
		})
		return
	}

	roleName := c.Query("role")
	guardName := c.Query("guard")
	if guardName == "" {
		guardName = "web"
	}

	if roleName == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Role name is required",
		})
		return
	}

	hasRole, err := models_roles_and_permissions.UserHasRoleByName(user.ID, roleName, guardName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Error checking user role: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data": gin.H{
			"user_id":  user.ID,
			"role":     roleName,
			"guard":    guardName,
			"has_role": hasRole,
		},
	})
}

// CheckCurrentUserPermission verifica si el usuario logueado tiene un permiso específico
func (upc *UserPermissionController) CheckCurrentUserPermission(c *gin.Context) {
	// Obtener el usuario autenticado de la sesión
	user, authenticated := helpers.GetAuthenticatedUser(c.Request)
	if !authenticated {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status":  "error",
			"message": "User not authenticated",
		})
		return
	}

	permissionName := c.Query("permission")
	guardName := c.Query("guard")
	if guardName == "" {
		guardName = "web"
	}

	if permissionName == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Permission name is required",
		})
		return
	}

	hasPermission, err := models_roles_and_permissions.UserHasPermission(user.ID, permissionName, guardName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Error checking user permission: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data": gin.H{
			"user_id":        user.ID,
			"permission":     permissionName,
			"guard":          guardName,
			"has_permission": hasPermission,
		},
	})
}
