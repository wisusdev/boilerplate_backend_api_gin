package commands

import (
	"crypto/rand"
	"encoding/base64"
	"semita/core/helpers"

	"github.com/spf13/cobra"
)

var KeyGenerateCmd = &cobra.Command{
	Use:   "key:generate",
	Short: "Genera una nueva clave JWT y la guarda en el archivo .env",
	Run: func(cmd *cobra.Command, args []string) {
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
