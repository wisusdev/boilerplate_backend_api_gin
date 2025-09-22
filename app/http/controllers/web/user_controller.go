package web

import (
	"boilerplate_backend_api_gin/app/data/models"
	"boilerplate_backend_api_gin/app/data/providers"
	"boilerplate_backend_api_gin/core/helpers"
	roleAndPermissionModels "boilerplate_backend_api_gin/core/roles_and_permissions/providers"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func UserIndex(context *gin.Context) {
	var users, errorUsers = providers.GetAllUsers()

	if errorUsers != nil {
		helpers.Logs("ERROR", fmt.Sprintf("Error al obtener los usuarios: %v", errorUsers))
		helpers.CreateFlashNotification(context.Writer, context.Request, "error", "Error al obtener los usuarios")
		context.Redirect(http.StatusSeeOther, "/")
		context.Abort()
		return
	}

	helpers.View(context, "users/index", "User Index", users)
}

func UserCreate(context *gin.Context) {
	var roles, errorRoles = roleAndPermissionModels.GetAllRoles()

	if errorRoles != nil {
		helpers.Logs("ERROR", fmt.Sprintf("Error al obtener los roles: %v", errorRoles))
		helpers.CreateFlashNotification(context.Writer, context.Request, "error", "Error al obtener los roles")
		context.Redirect(http.StatusSeeOther, "/")
		context.Abort()
		return
	}

	helpers.View(context, "users/create", "User Create", roles)
}

func UserStore(context *gin.Context) {

	firstName := context.PostForm("first_name")
	lastName := context.PostForm("last_name")
	username := context.PostForm("username")
	email := context.PostForm("email")
	password := context.PostForm("password")
	confirmPassword := context.PostForm("confirm_password")
	role := context.PostForm("role")

	if firstName == "" || lastName == "" || username == "" || email == "" || password == "" || confirmPassword == "" || role == "" {
		helpers.CreateFlashNotification(context.Writer, context.Request, "error", "Todos los campos son obligatorios")
		context.Redirect(http.StatusSeeOther, "/users/create")
		context.Abort()
		return
	}

	if len(password) < 6 {
		helpers.CreateFlashNotification(context.Writer, context.Request, "error", "La contraseña debe tener al menos 6 caracteres")
		context.Redirect(http.StatusSeeOther, "/users/create")
		context.Abort()
		return
	}

	if password != confirmPassword {
		helpers.CreateFlashNotification(context.Writer, context.Request, "error", "Las contraseñas no coinciden")
		context.Redirect(http.StatusSeeOther, "/users/create")
		context.Abort()
		return
	}

	// Verificar si el nombre de usuario ya existe

	var user = models.UserStruct{
		FirstName: context.PostForm("first_name"),
		LastName:  context.PostForm("last_name"),
		Username:  context.PostForm("username"),
		Email:     context.PostForm("email"),
		Password:  context.PostForm("password"),
	}

	var userStore, errorStore = providers.StoreUser(user)
	if errorStore != nil {
		http.Error(context.Writer, "Error al guardar el usuario en la base de datos", http.StatusInternalServerError)
		return
	}

	// Asignar rol al usuario
	var intRoleID, errorParse = strconv.ParseInt(role, 10, 64)
	if errorParse != nil {
		http.Error(context.Writer, "ID de rol inválido", http.StatusBadRequest)
		return
	}

	var errorAssignRole = roleAndPermissionModels.AssignRoleToUser(userStore.ID, string(intRoleID))
	if errorAssignRole != nil {
		http.Error(context.Writer, "Error al asignar el rol al usuario", http.StatusInternalServerError)
		return
	}

	context.Redirect(http.StatusSeeOther, "/users")
	context.Abort()
}

func UserShow(context *gin.Context) {
	var id = context.Param("id")

	var user, errorUser = providers.GetUserByID(id)
	if errorUser != nil {
		http.Error(context.Writer, "Error al obtener el usuario desde la base de datos", http.StatusInternalServerError)
		return
	}

	helpers.View(context, "users/show", "User Show", user)
}

func UserEdit(context *gin.Context) {
	var id = context.Param("id")

	var user, errorUser = providers.GetUserByID(id)
	if errorUser != nil {
		http.Error(context.Writer, "Error al obtener el usuario desde la base de datos", http.StatusInternalServerError)
		return
	}

	helpers.View(context, "users/edit", "User Edit", user)
}

func UserUpdate(context *gin.Context) {
	var id = context.Param("id")

	var user, errorUser = providers.GetUserByID(id)
	if errorUser != nil {
		http.Error(context.Writer, "Error al obtener el usuario desde la base de datos", http.StatusInternalServerError)
		return
	}

	user.FirstName = context.PostForm("first_name")
	user.LastName = context.PostForm("last_name")
	user.Email = context.PostForm("email")
	user.Password = context.PostForm("password")

	var errorUpdate = providers.UpdateUser(user)
	if errorUpdate != nil {
		http.Error(context.Writer, "Error al actualizar el usuario en la base de datos", http.StatusInternalServerError)
		return
	}

	context.Redirect(http.StatusSeeOther, "/users")
	context.Abort()
}

func UserDelete(context *gin.Context) {
	var id = context.Param("id")

	var intID, errorParse = strconv.ParseInt(id, 10, 64)
	if errorParse != nil {
		http.Error(context.Writer, "ID de usuario inválido", http.StatusBadRequest)
		return
	}

	var errorDelete = providers.DeleteUser(strconv.FormatInt(intID, 10))
	if errorDelete != nil {
		http.Error(context.Writer, "Error al eliminar el usuario desde la base de datos", http.StatusInternalServerError)
		return
	}

	context.Redirect(http.StatusSeeOther, "/users")
	context.Abort()
}
