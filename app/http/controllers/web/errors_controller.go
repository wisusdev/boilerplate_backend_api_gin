package web

import (
	"net/http"
	"semita/config"
	"text/template"

	"github.com/gin-gonic/gin"
)

func Error404(context *gin.Context) {
	templateError404 := template.Must(template.ParseFiles("resources/error/404.html", config.AppConfig().Layout))
	err := templateError404.Execute(context.Writer, nil)
	if err != nil {
		context.String(http.StatusInternalServerError, "Error al cargar la plantilla 404")
		return
	}
	context.Status(http.StatusNotFound)
}
