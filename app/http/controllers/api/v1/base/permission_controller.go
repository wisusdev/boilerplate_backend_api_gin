package base

import (
	"net/http"
	"semita/core/helpers"
	"semita/core/roles_and_permissions/models"
	"semita/core/roles_and_permissions/providers"
	"strconv"

	"github.com/gin-gonic/gin"
)

// PermissionController maneja las operaciones CRUD de permisos
type PermissionController struct{}

// Index muestra todos los permisos
func (pc *PermissionController) Index(c *gin.Context) {
	permissions, err := providers.GetAllPermissions()
	if err != nil {
		helpers.CreateFlashNotification(c.Writer, c.Request, "error", "Error retrieving permissions: "+err.Error())
		c.Redirect(http.StatusSeeOther, "/")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   permissions,
	})
}

// Show muestra un permiso específico
func (pc *PermissionController) Show(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Invalid permission ID",
		})
		return
	}

	permission, err := providers.GetPermissionByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  "error",
			"message": "Permission not found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   permission,
	})
}

// Store crea un nuevo permiso
func (pc *PermissionController) Store(c *gin.Context) {
	var permissionData models.CreatePermissionStruct
	if err := c.ShouldBindJSON(&permissionData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Invalid input data",
			"errors":  err.Error(),
		})
		return
	}

	permission, err := providers.CreatePermission(permissionData)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Error creating permission: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"status":  "success",
		"message": "Permission created successfully",
		"data":    permission,
	})
}

// Update actualiza un permiso existente
func (pc *PermissionController) Update(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Invalid permission ID",
		})
		return
	}

	var permissionData models.CreatePermissionStruct
	if err := c.ShouldBindJSON(&permissionData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Invalid input data",
			"errors":  err.Error(),
		})
		return
	}

	permission, err := providers.UpdatePermission(id, permissionData)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Error updating permission: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Permission updated successfully",
		"data":    permission,
	})
}

// Delete elimina un permiso
func (pc *PermissionController) Delete(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Invalid permission ID",
		})
		return
	}

	err = providers.DeletePermission(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Error deleting permission: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Permission deleted successfully",
	})
}

// AssignToUser asigna un permiso directamente a un usuario
func (pc *PermissionController) AssignToUser(c *gin.Context) {
	var request models.AssignPermissionRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Invalid input data",
			"errors":  err.Error(),
		})
		return
	}

	if request.UserID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "User ID is required",
		})
		return
	}

	err := providers.AssignPermissionToUser(request.UserID, request.PermissionID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Error assigning permission to user: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Permission assigned to user successfully",
	})
}

// AssignToRole asigna un permiso a un rol
func (pc *PermissionController) AssignToRole(c *gin.Context) {
	var request models.AssignPermissionRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Invalid input data",
			"errors":  err.Error(),
		})
		return
	}

	if request.RoleID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Role ID is required",
		})
		return
	}

	err := providers.AssignPermissionToRole(request.RoleID, request.PermissionID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Error assigning permission to role: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Permission assigned to role successfully",
	})
}

// RevokeFromUser revoca un permiso directo de un usuario
func (pc *PermissionController) RevokeFromUser(c *gin.Context) {
	var request models.AssignPermissionRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Invalid input data",
			"errors":  err.Error(),
		})
		return
	}

	if request.UserID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "User ID is required",
		})
		return
	}

	err := providers.RevokePermissionFromUser(request.UserID, request.PermissionID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Error revoking permission from user: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Permission revoked from user successfully",
	})
}

// RevokeFromRole revoca un permiso de un rol
func (pc *PermissionController) RevokeFromRole(c *gin.Context) {
	var request models.AssignPermissionRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Invalid input data",
			"errors":  err.Error(),
		})
		return
	}

	if request.RoleID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Role ID is required",
		})
		return
	}

	err := providers.RevokePermissionFromRole(request.RoleID, request.PermissionID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Error revoking permission from role: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Permission revoked from role successfully",
	})
}

// GetUserPermissions obtiene todos los permisos de un usuario (directos + heredados)
func (pc *PermissionController) GetUserPermissions(c *gin.Context) {
	userID := c.Param("user_id")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Invalid user ID",
		})
		return
	}

	directPermissions, err := providers.GetUserDirectPermissions(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Error retrieving user direct permissions: " + err.Error(),
		})
		return
	}

	allPermissions, err := providers.GetUserAllPermissions(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Error retrieving user permissions: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data": gin.H{
			"direct_permissions": directPermissions,
			"all_permissions":    allPermissions,
		},
	})
}

// GetRolePermissions obtiene todos los permisos de un rol
func (pc *PermissionController) GetRolePermissions(c *gin.Context) {
	roleIDParam := c.Param("role_id")
	roleID, err := strconv.Atoi(roleIDParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Invalid role ID",
		})
		return
	}

	permissions, err := providers.GetRolePermissions(roleID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Error retrieving role permissions: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   permissions,
	})
}
