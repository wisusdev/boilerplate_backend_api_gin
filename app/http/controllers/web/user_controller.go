package web

import (
	"fmt"
	"net/http"
	"semita/app/data/models"
	"semita/app/data/structs"
	"semita/config"
	"semita/core/helpers"
	"strconv"
	"text/template"

	"github.com/gin-gonic/gin"
)

func UserIndex(context *gin.Context) {
	var users, errorUsers = models.GetAllUsers()

	if errorUsers != nil {
		helpers.Logs("ERROR", fmt.Sprintf("Error al obtener los usuarios: %v", errorUsers))
		helpers.CreateFlashNotification(context.Writer, context.Request, "error", "Error al obtener los usuarios")
		http.Error(context.Writer, "Error al obtener los usuarios desde la base de datos", http.StatusInternalServerError)
		return
	}

	var viewData = helpers.AuthSessionService(context.Writer, context.Request, "User Index", users)

	var templateIndexPath = "resources/users/index.html"
	var templateIndex = template.Must(template.ParseFiles(templateIndexPath, config.AppConfig().Layout))
	var errorExecuteTemplate = templateIndex.Execute(context.Writer, viewData)
	if errorExecuteTemplate != nil {
		fmt.Println("Error al ejecutar la plantilla:", errorExecuteTemplate)
		return
	}
}

func UserCreate(context *gin.Context) {
	var viewData = helpers.AuthSessionService(context.Writer, context.Request, "User Create", nil)
	var templateCreatePath = "resources/users/create.html"
	var templateCreate = template.Must(template.ParseFiles(templateCreatePath, config.AppConfig().Layout))
	var errorExecuteTemplate = templateCreate.Execute(context.Writer, viewData)
	if errorExecuteTemplate != nil {
		fmt.Println("Error al ejecutar la plantilla:", errorExecuteTemplate)
		return
	}
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

	var viewData = helpers.AuthSessionService(context.Writer, context.Request, "User Create", user)

	var templateShowPath = "resources/users/show.html"
	var templateShow = template.Must(template.ParseFiles(templateShowPath, config.AppConfig().Layout))
	var errorExecuteTemplate = templateShow.Execute(context.Writer, viewData)
	if errorExecuteTemplate != nil {
		fmt.Println("Error al ejecutar la plantilla:", errorExecuteTemplate)
		return
	}
}

func UserEdit(context *gin.Context) {
	var id = context.Param("id")

	var user, errorUser = models.GetUserByID(id)
	if errorUser != nil {
		http.Error(context.Writer, "Error al obtener el usuario desde la base de datos", http.StatusInternalServerError)
		return
	}

	var viewData = helpers.AuthSessionService(context.Writer, context.Request, "User Edit", user)

	var templateEditPath = "resources/users/edit.html"
	var templateEdit = template.Must(template.ParseFiles(templateEditPath, config.AppConfig().Layout))
	var errorExecuteTemplate = templateEdit.Execute(context.Writer, viewData)
	if errorExecuteTemplate != nil {
		fmt.Println("Error al ejecutar la plantilla:", errorExecuteTemplate)
		return
	}
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
