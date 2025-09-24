package base

import (
	"boilerplate_backend_api_gin/app/data/models"
	"boilerplate_backend_api_gin/app/data/providers"
	"boilerplate_backend_api_gin/app/http/requests"
	"boilerplate_backend_api_gin/core/helpers"
	roleProviders "boilerplate_backend_api_gin/core/roles_and_permissions/providers"
	"boilerplate_backend_api_gin/core/validators"
	"net/http"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"github.com/gin-gonic/gin"
)

type UserController struct{}

func (uc *UserController) Index(c *gin.Context) {

	if !helpers.HasPermissionGin(c, "users:index") {
		c.JSON(http.StatusForbidden, gin.H{
			"status":  "error",
			"message": "You don't have permission to access this resource",
		})
		return
	}

	users, err := providers.GetAllUsers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Error retrieving users: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   users,
	})
}

func (uc *UserController) Show(c *gin.Context) {

	if !helpers.HasPermissionGin(c, "users:show") {
		c.JSON(http.StatusForbidden, gin.H{
			"status":  "error",
			"message": "You don't have permission to access this resource",
		})
		return
	}

	id := c.Param("id")
	user, err := providers.GetUserByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  "error",
			"message": "User not found: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   user,
	})
}

func (uc *UserController) Store(c *gin.Context) {

	if !helpers.HasPermissionGin(c, "users:create") {
		c.JSON(http.StatusForbidden, gin.H{
			"status":  "error",
			"message": "You don't have permission to access this resource",
		})
		return
	}

	var request requests.UserRequest

	if err := validators.Validate(c, &request); err != nil {
		return
	}

	// Hash de la contraseña
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(request.Data.Attributes.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Error processing password",
		})
		return
	}

	// Convertir a la estructura interna
	userData := models.UserStruct{
		Username:  request.Data.Attributes.Username,
		FirstName: request.Data.Attributes.FirstName,
		LastName:  request.Data.Attributes.LastName,
		Email:     request.Data.Attributes.Email,
		Password:  string(hashedPassword),
		Language:  "es", // valor por defecto
	}

	newUser, err := providers.StoreUser(userData)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Error creating user: " + err.Error(),
		})
		return
	}

	// Asignar roles si se especificaron
	if len(request.Data.Attributes.Roles) > 0 {
		for _, roleName := range request.Data.Attributes.Roles {
			role, err := roleProviders.GetRoleByName(roleName, "api")
			if err == nil {
				roleProviders.AssignRoleToUser(newUser.ID, role.ID)
			}
		}
	}

	c.JSON(http.StatusCreated, gin.H{
		"status":  "success",
		"message": "User created successfully",
		"data":    newUser,
	})
}

func (uc *UserController) Update(c *gin.Context) {

	if !helpers.HasPermissionGin(c, "users:update") {
		c.JSON(http.StatusForbidden, gin.H{
			"status":  "error",
			"message": "You don't have permission to access this resource",
		})
		return
	}

	userID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Invalid user ID format",
		})
		return
	}

	var request requests.UserRequest

	if err := validators.Validate(c, &request); err != nil {
		return
	}

	// Convertir a la estructura interna
	userData := models.UserStruct{
		ID:        userID.String(),
		Username:  request.Data.Attributes.Username,
		FirstName: request.Data.Attributes.FirstName,
		LastName:  request.Data.Attributes.LastName,
		Email:     request.Data.Attributes.Email,
	}

	// Hash de la contraseña solo si se proporcionó
	if request.Data.Attributes.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(request.Data.Attributes.Password), bcrypt.DefaultCost)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  "error",
				"message": "Error processing password",
			})
			return
		}
		userData.Password = string(hashedPassword)
	}

	err = providers.UpdateUser(userData)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Error updating user: " + err.Error(),
		})
		return
	}

	// Obtener el usuario actualizado
	updatedUser, err := providers.GetUserByID(userID.String())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Error retrieving updated user: " + err.Error(),
		})
		return
	}

	// Actualizar roles si se especificaron
	if len(request.Data.Attributes.Roles) > 0 {
		// Obtener roles actuales del usuario para revocarlos
		currentRoles, err := roleProviders.GetUserRoles(userID.String())
		if err == nil {
			for _, role := range currentRoles {
				roleProviders.RevokeRoleFromUser(userID.String(), role.ID)
			}
		}

		// Luego asignar los nuevos roles
		for _, roleName := range request.Data.Attributes.Roles {
			role, err := roleProviders.GetRoleByName(roleName, "api")
			if err == nil {
				roleProviders.AssignRoleToUser(userID.String(), role.ID)
			}
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "User updated successfully",
		"data":    updatedUser,
	})
}

func (uc *UserController) Delete(c *gin.Context) {

	if !helpers.HasPermissionGin(c, "users:delete") {
		c.JSON(http.StatusForbidden, gin.H{
			"status":  "error",
			"message": "You don't have permission to access this resource",
		})
		return
	}

	id := c.Param("id")
	err := providers.DeleteUser(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Error deleting user: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "User deleted successfully",
	})
}
