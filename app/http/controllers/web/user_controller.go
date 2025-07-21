package web

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"semita/app/data/models"
	"semita/app/data/structs"
	"semita/core/helpers"
	"strconv"
)

func UserIndex(context *gin.Context) {
	var users, errorUsers = models.GetAllUsers()

	if errorUsers != nil {
		helpers.Logs("ERROR", fmt.Sprintf("Error al obtener los usuarios: %v", errorUsers))
		helpers.CreateFlashNotification(context.Writer, context.Request, "error", "Error al obtener los usuarios")
		http.Error(context.Writer, "Error al obtener los usuarios desde la base de datos", http.StatusInternalServerError)
		return
	}

	helpers.View(context, "users/index", "User Index", users)
}

func UserCreate(context *gin.Context) {
	helpers.View(context, "users/create", "User Create", nil)
}

func UserStore(context *gin.Context) {
	var user = structs.StoreUserStruct{
		Name:     context.PostForm("name"),
		Email:    context.PostForm("email"),
		Password: context.PostForm("password"),
	}

	var errorStore = models.StoreUser(user)
	if errorStore != nil {
		http.Error(context.Writer, "Error al guardar el usuario en la base de datos", http.StatusInternalServerError)
		return
	}

	context.Redirect(http.StatusSeeOther, "/users")
	context.Abort()
}

func UserShow(context *gin.Context) {
	var id = context.Param("id")

	var user, errorUser = models.GetUserByID(id)
	if errorUser != nil {
		http.Error(context.Writer, "Error al obtener el usuario desde la base de datos", http.StatusInternalServerError)
		return
	}

	helpers.View(context, "users/show", "User Show", user)
}

func UserEdit(context *gin.Context) {
	var id = context.Param("id")

	var user, errorUser = models.GetUserByID(id)
	if errorUser != nil {
		http.Error(context.Writer, "Error al obtener el usuario desde la base de datos", http.StatusInternalServerError)
		return
	}

	helpers.View(context, "users/edit", "User Edit", user)
}

func UserUpdate(context *gin.Context) {
	var id = context.Param("id")

	var intID, errorParse = strconv.ParseInt(id, 10, 64)
	if errorParse != nil {
		http.Error(context.Writer, "ID de usuario inválido", http.StatusBadRequest)
		return
	}

	var user = structs.UpdateUserStruct{
		ID:       int(intID),
		Name:     context.PostForm("name"),
		Email:    context.PostForm("email"),
		Password: context.PostForm("password"),
	}

	var errorUpdate = models.UpdateUser(user)
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

	var errorDelete = models.DeleteUser(strconv.FormatInt(intID, 10))
	if errorDelete != nil {
		http.Error(context.Writer, "Error al eliminar el usuario desde la base de datos", http.StatusInternalServerError)
		return
	}

	context.Redirect(http.StatusSeeOther, "/users")
	context.Abort()
}
