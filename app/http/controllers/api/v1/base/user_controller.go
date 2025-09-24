package base

import (
	"boilerplate_backend_api_gin/app/data/models"
	"boilerplate_backend_api_gin/app/data/providers"
	"boilerplate_backend_api_gin/core/helpers"
	"net/http"

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

	var userData models.UserStruct
	if err := c.ShouldBindJSON(&userData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Invalid request data: " + err.Error(),
		})
		return
	}

	newUser, err := providers.StoreUser(userData)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Error creating user: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"status": "success",
		"data":   newUser,
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

	id := c.Param("id")
	var userData models.UserStruct
	if err := c.ShouldBindJSON(&userData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Invalid request data: " + err.Error(),
		})
		return
	}

	userData.ID = id // Asegurarse de que el ID esté establecido
	err := providers.UpdateUser(userData)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Error updating user: " + err.Error(),
		})

		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "User updated successfully",
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
