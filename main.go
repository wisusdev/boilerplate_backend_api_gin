package main

import (
	"boilerplate_backend_api_gin/app/http/controllers/web"
	"boilerplate_backend_api_gin/bootstrap"
	"boilerplate_backend_api_gin/config"
	"boilerplate_backend_api_gin/core/helpers"
	"boilerplate_backend_api_gin/core/internationalization"
	"boilerplate_backend_api_gin/routes"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {

	config.GoDotEnv()                       // Load environment variables
	internationalization.LoadTranslations() // Load internationalization

	if os.Getenv("AIR") != "" || len(os.Args) == 1 {
		StartServer()
		return
	} else if len(os.Args) > 1 {
		args := os.Args[1:]
		bootstrap.Execute(args) // Ejecutar comandos de Cobra cuando hay argumentos
	}

}

func StartServer() {
	// Cargar variables de entorno
	var appUrl = config.AppConfig().Url

	router := gin.Default()

	// Rutas web
	routes.Web(router.Group("/"))

	// Rutas API
	routes.Api(router.Group("/api/v1"))

	// Archivos estáticos
	router.Static("/public", "./public")

	// Ruta 404 personalizada
	router.NoRoute(web.Error404)

	// Ejecución del servidor
	server := &http.Server{
		Addr:         appUrl,
		Handler:      router,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	fmt.Printf("✅ Servidor corriendo en %v\n", helpers.ColorGreen("http://"+appUrl))
	log.Fatal(server.ListenAndServe())
}
