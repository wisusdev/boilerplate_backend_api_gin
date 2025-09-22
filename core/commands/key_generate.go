package commands

import (
	"boilerplate_backend_api_gin/core/cli"
	"boilerplate_backend_api_gin/core/helpers"
	"crypto/rand"
	"encoding/base64"
)

var KeyGenerateCmd = cli.Command{
	Name:        "key:generate",
	Description: "Genera una nueva clave JWT y la guarda en el archivo .env",
	Execute: func(args []string) {
		key := generateRandomKey(32)

		// Intenta actualizar el archivo .env
		helpers.UpdateEnvFile("APP_KEY", key)
	},
}

func generateRandomKey(length int) string {
	b := make([]byte, length)
	_, err := rand.Read(b)
	if err != nil {
		panic(err)
	}
	return base64.StdEncoding.EncodeToString(b)
}
