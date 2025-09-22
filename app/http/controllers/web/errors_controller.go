package web

import (
	"boilerplate_backend_api_gin/core/helpers"

	"github.com/gin-gonic/gin"
)

func Error404(context *gin.Context) {
	context.Status(404)
	helpers.View(context, "error/404.html", "Página no encontrada", nil)
}
