package base

import (
	"boilerplate_backend_api_gin/core/helpers"
	"boilerplate_backend_api_gin/core/roles_and_permissions/models"
	"boilerplate_backend_api_gin/core/roles_and_permissions/providers"
	"boilerplate_backend_api_gin/core/roles_and_permissions/repositories"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// RoleController maneja las operaciones CRUD de roles
type RoleController struct{}

// Index muestra todos los roles
func (roleController *RoleController) Index(context *gin.Context) {

	if !helpers.HasPermissionGin(context, "roles:index") {
		context.JSON(http.StatusForbidden, gin.H{
			"status":  "error",
			"message": "You don't have permission to access this resource",
		})
		return
	}

	roles, err := providers.GetAllRoles()
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Error retrieving roles: " + err.Error(),
		})
		return
	}

	context.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   roles,
	})
}

// Show muestra un rol específico con sus permisos
func (roleController *RoleController) Show(c *gin.Context) {

	if !helpers.HasPermissionGin(c, "roles:show") {
		c.JSON(http.StatusForbidden, gin.H{
			"status":  "error",
			"message": "You don't have permission to access this resource",
		})
		return
	}

	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Invalid role ID",
		})
		return
	}

	role, err := providers.GetRoleByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  "error",
			"message": "Role not found",
		})
		return
	}

	permissions, err := providers.GetRolePermissions(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Error retrieving role permissions: " + err.Error(),
		})
		return
	}

	roleWithPermissions := models.RoleWithPermissions{
		RoleStruct:  *role,
		Permissions: permissions,
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   roleWithPermissions,
	})
}

// Store crea un nuevo rol
func (roleController *RoleController) Store(c *gin.Context) {

	if !helpers.HasPermissionGin(c, "roles:store") {
		c.JSON(http.StatusForbidden, gin.H{
			"status":  "error",
			"message": "You don't have permission to access this resource",
		})
		return
	}

	var roleData models.CreateRoleStruct
	if err := c.ShouldBindJSON(&roleData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Invalid input data",
			"errors":  err.Error(),
		})
		return
	}

	role, err := providers.CreateRole(roleData)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Error creating role: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"status":  "success",
		"message": "Role created successfully",
		"data":    role,
	})
}

// Update actualiza un rol existente
func (roleController *RoleController) Update(c *gin.Context) {

	if !helpers.HasPermissionGin(c, "roles:update") {
		c.JSON(http.StatusForbidden, gin.H{
			"status":  "error",
			"message": "You don't have permission to access this resource",
		})
		return
	}

	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Invalid role ID",
		})
		return
	}

	var roleData models.CreateRoleStruct
	if err := c.ShouldBindJSON(&roleData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Invalid input data",
			"errors":  err.Error(),
		})
		return
	}

	role, err := repositories.UpdateRole(id, roleData)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Error updating role: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Role updated successfully",
		"data":    role,
	})
}

// Delete elimina un rol
func (roleController *RoleController) Delete(c *gin.Context) {

	if !helpers.HasPermissionGin(c, "roles:delete") {
		c.JSON(http.StatusForbidden, gin.H{
			"status":  "error",
			"message": "You don't have permission to access this resource",
		})
		return
	}

	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Invalid role ID",
		})
		return
	}

	err = providers.DeleteRole(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Error deleting role: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Role deleted successfully",
	})
}

// AssignToUser asigna un rol a un usuario
func (roleController *RoleController) AssignToUser(c *gin.Context) {
	var request models.AssignRoleRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Invalid input data",
			"errors":  err.Error(),
		})
		return
	}

	err := providers.AssignRoleToUser(request.UserID, request.RoleID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Error assigning role to user: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Role assigned to user successfully",
	})
}

// RevokeFromUser revoca un rol de un usuario
func (roleController *RoleController) RevokeFromUser(c *gin.Context) {
	var request models.AssignRoleRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Invalid input data",
			"errors":  err.Error(),
		})
		return
	}

	err := providers.RevokeRoleFromUser(request.UserID, request.RoleID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Error revoking role from user: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Role revoked from user successfully",
	})
}

// GetUserRoles obtiene todos los roles de un usuario
func (roleController *RoleController) GetUserRoles(c *gin.Context) {
	userID := c.Param("user_id")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Invalid user ID",
		})
		return
	}

	roles, err := providers.GetUserRoles(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Error retrieving user roles: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   roles,
	})
}
