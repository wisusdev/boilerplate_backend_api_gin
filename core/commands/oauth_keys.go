package commands

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

var OauthKeysCmd = &cobra.Command{
	Use:   "oauth:keys",
	Short: "Genera las llaves oauth-private.key y oauth-public.key en el directorio storage",
	Run: func(cmd *cobra.Command, args []string) {
		storageDir := "storage/oauth"
		privateKeyPath := filepath.Join(storageDir, "oauth-private.key")
		publicKeyPath := filepath.Join(storageDir, "oauth-public.key")

		if _, err := os.Stat(storageDir); os.IsNotExist(err) {
			os.MkdirAll(storageDir, 0755)
		}

		privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
		if err != nil {
			fmt.Println("Error generando la llave privada:", err)
			return
		}

		privFile, err := os.Create(privateKeyPath)
		if err != nil {
			fmt.Println("No se pudo crear el archivo de llave privada:", err)
			return
		}
		defer privFile.Close()

		privBytes := x509.MarshalPKCS1PrivateKey(privateKey)
		pem.Encode(privFile, &pem.Block{Type: "RSA PRIVATE KEY", Bytes: privBytes})

		pubFile, err := os.Create(publicKeyPath)
		if err != nil {
			fmt.Println("No se pudo crear el archivo de llave pública:", err)
			return
		}
		defer pubFile.Close()

		pubASN1, err := x509.MarshalPKIXPublicKey(&privateKey.PublicKey)

		if err != nil {
			fmt.Println("Error serializando la llave pública:", err)
			return
		}

		pem.Encode(pubFile, &pem.Block{Type: "PUBLIC KEY", Bytes: pubASN1})

		fmt.Println("Llaves OAuth generadas en el directorio storage:")
		fmt.Println("-", privateKeyPath)
		fmt.Println("-", publicKeyPath)
	},
}
