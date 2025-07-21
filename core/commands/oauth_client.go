package commands

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"os"
	"semita/core/oauth/oauth_models"

	"github.com/spf13/cobra"
)

var OauthClientCmd = &cobra.Command{
	Use:   "oauth:client",
	Short: "Crea un cliente OAuth en la base de datos",
	Run: func(cmd *cobra.Command, args []string) {
		name := "Default Client"
		if len(args) > 0 {
			name = args[0]
		}
		clientID := randomHex(16)
		clientSecret := randomHex(32)

		err := oauth_models.CreateOAuthClient(name, clientID, clientSecret)
		if err != nil {
			fmt.Println("Error creando el cliente OAuth:", err)
			os.Exit(1)
		}
		fmt.Println("Cliente OAuth creado correctamente:")
		fmt.Println("ID:", clientID)
		fmt.Println("Secret:", clientSecret)
	},
}

func randomHex(n int) string {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		panic(err)
	}
	return hex.EncodeToString(b)
}
