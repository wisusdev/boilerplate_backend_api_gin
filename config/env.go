package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

// GoDotEnv carga el archivo .env una sola vez al inicio de la aplicación.
func GoDotEnv() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("❌ No se pudo cargar el archivo .env, usando variables de entorno del sistema")
	}

	fmt.Println("✅ Archivo .env cargado correctamente")
}

// GetEnv obtiene una variable de entorno, o retorna un valor por defecto si no existe.
func GetEnv(key string, defaultKey string) string {
	var val = os.Getenv(key)

	if val == "" {
		if defaultKey == "" {
			fmt.Printf("❌ La variable de entorno '%s' no está definida\n", key)
			return ""
		}

		val = defaultKey
	}

	return val
}

func GetEnvBool(key string, defaultKey bool) bool {
	val := os.Getenv(key)
	if val == "" {
		return defaultKey
	}

	if val == "true" || val == "1" {
		return true
	} else if val == "false" || val == "0" {
		return false
	}

	fmt.Printf("❌ La variable de entorno '%s' no es un booleano válido, usando valor por defecto: %v\n", key, defaultKey)
	return defaultKey
}

func GetEnvInt(key string, defaultKey int) int {
	val := os.Getenv(key)
	if val == "" {
		return defaultKey
	}

	intVal, err := strconv.Atoi(val)
	if err != nil {
		fmt.Printf("❌ La variable de entorno '%s' no es un entero válido, usando valor por defecto: %d\n", key, defaultKey)
		return defaultKey
	}

	return intVal
}

func UpdateEnvFile(key, value string) {
	file := ".env"
	input, err := os.ReadFile(file)
	if err != nil {
		fmt.Printf("No se pudo leer %s, solo se mostrará la clave.\n", file)
		return
	}
	lines := strings.Split(string(input), "\n")
	found := false
	for i, line := range lines {
		if strings.HasPrefix(line, key+"=") {
			lines[i] = key + "=" + value
			found = true
			break
		}
	}
	if !found {
		lines = append(lines, key+"="+value)
	}
	output := strings.Join(lines, "\n")
	err = os.WriteFile(file, []byte(output), 0644)

	if err != nil {
		fmt.Printf("❌ Error al actualizar %s: %v\n", file, err)
		return
	}

	fmt.Printf("%s actualizado en %s\n", key, file)
}
