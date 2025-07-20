package helpers

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func SetupLogger() *os.File {
	// Crear el directorio de logs si no existe
	var logDir = "storage/logs"
	if errorCreateDir := os.MkdirAll(logDir, 0775); errorCreateDir != nil {
		log.Fatal("No se pudo crear el directorio de logs:", errorCreateDir)
	}

	// Generamos el nombre del archivo de log con la fecha actual
	var today = time.Now().Format("2006-01-02")
	var logFileName = fmt.Sprintf("app-%s.log", today)
	var logPath = filepath.Join(logDir, logFileName)

	// Crea o abre el archivo de log del día
	var file, err = os.OpenFile(logPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal("No se pudo abrir el archivo de log:", err)
	}

	log.SetOutput(file)

	return file
}

func Logs(logType string, logDescription string) {
	// Configuramos el logger y obtenemos el archivo
	var file = SetupLogger()

	// Aseguramos cerrar el archivo al finalizar
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			log.Fatal("No se pudo abrir el archivo de log:", err)
		}
	}(file)

	log.Printf("%s: %s", strings.ToUpper(logType), logDescription)
}
