package commands

import (
	crand "crypto/rand"
	"encoding/hex"
	"fmt"
	mrand "math/rand"
	"os"
	"semita/core/cli"
	"semita/core/oauth/oauth_models"
)

var OauthClientCmd = cli.Command{
	Name:        "oauth:client",
	Description: "Crea un cliente OAuth en la base de datos",
	Execute: func(args []string) {
		name := "Default Client"
		if len(args) > 0 {
			name = args[0]
		}
		clientID := int64(mrand.Int31() + 1)
		clientSecret := randomHex(32)

		err := oauth_models.CreateOAuthClient(clientID, name, clientSecret)
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
	_, err := crand.Read(b)
	if err != nil {
		panic(err)
	}
	return hex.EncodeToString(b)
}
