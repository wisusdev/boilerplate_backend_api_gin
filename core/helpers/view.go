package helpers

import (
	"fmt"
	"net/http"
	"path/filepath"
	"semita/config"
	"text/template"

	"github.com/gin-gonic/gin"
)

// View renderiza una vista con el layout principal y contexto de sesión usando gin.Context
func View(context *gin.Context, viewPath string, title string, data interface{}) {
	authData := AuthSessionService(context.Writer, context.Request, title, data)

	fullViewPath := filepath.Join("resources", "views", viewPath)
	tmpl := template.Must(template.ParseFiles(fullViewPath, config.AppConfig().Layout))

	err := tmpl.Execute(context.Writer, authData)

	if err != nil {
		Logs("ERROR", fmt.Sprintf("Error al renderizar la vista: %v", err))
		http.Error(context.Writer, "Error interno al renderizar la vista", http.StatusInternalServerError)
	}
}
