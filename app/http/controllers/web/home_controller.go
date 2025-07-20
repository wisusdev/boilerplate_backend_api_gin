package web

import (
	"semita/core/helpers"

	"github.com/gin-gonic/gin"
)

type Habilidades struct {
	Nombre string
}

type Data struct {
	Nombre      string
	Edad        int
	Perfil      string
	Habilidades []Habilidades
}

func HomeIndex(c *gin.Context) {
	helpers.View(c, "home/home.html", "Inicio", nil)
}
