package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"semita/app/http/controllers/web"
	"semita/bootstrap"
	"semita/config"
	"semita/core/helpers"
	"semita/core/internationalization"
	"semita/routes"
	"time"
)

func main() {

	config.GoDotEnv()                       // Load environment variables
	internationalization.LoadTranslations() // Load internationalization

	if os.Getenv("AIR") != "" || len(os.Args) == 1 {
		StartServer()
		return
	} else if len(os.Args) > 1 {
		// Si hay argumentos, ejecuta los comandos de Cobra
		bootstrap.Execute() // Ejecutar comandos de Cobra cuando hay argumentos
	}

}

func StartServer() {
	// Cargar variables de entorno
	var appUrl = config.AppConfig().Url

	// Inicializar el enrutador Gin
	router := routes.Web()

	// Montar rutas API
	apiGroup := router.Group("/api/v1")
	routes.Api(apiGroup)

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
