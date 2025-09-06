package helpers

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/hex"
	"encoding/pem"
	"errors"
	"fmt"
	"os"
	"slices"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// OAuthTokenClaims define la estructura de los claims del token JWT compatible con Laravel Passport
type OAuthTokenClaims struct {
	Audience  string   `json:"aud"`
	JTI       string   `json:"jti"`
	IssuedAt  float64  `json:"iat"`
	NotBefore float64  `json:"nbf"`
	ExpiresAt float64  `json:"exp"`
	Subject   string   `json:"sub"`
	Scopes    []string `json:"scopes"`
}

// Valid implements jwt.Claims
func (c OAuthTokenClaims) Valid() error {
	now := float64(time.Now().Unix())

	if c.ExpiresAt < now {
		return errors.New("token is expired")
	}

	if c.NotBefore > now {
		return errors.New("token used before valid")
	}

	return nil
}

// GetExpirationTime implements jwt.Claims
func (c OAuthTokenClaims) GetExpirationTime() (*jwt.NumericDate, error) {
	return jwt.NewNumericDate(time.Unix(int64(c.ExpiresAt), 0)), nil
}

// GetNotBefore implements jwt.Claims
func (c OAuthTokenClaims) GetNotBefore() (*jwt.NumericDate, error) {
	return jwt.NewNumericDate(time.Unix(int64(c.NotBefore), 0)), nil
}

// GetIssuedAt implements jwt.Claims
func (c OAuthTokenClaims) GetIssuedAt() (*jwt.NumericDate, error) {
	return jwt.NewNumericDate(time.Unix(int64(c.IssuedAt), 0)), nil
}

// GetAudience implements jwt.Claims
func (c OAuthTokenClaims) GetAudience() (jwt.ClaimStrings, error) {
	return jwt.ClaimStrings{c.Audience}, nil
}

// GetIssuer implements jwt.Claims
func (c OAuthTokenClaims) GetIssuer() (string, error) {
	return "", nil
}

// GetSubject implements jwt.Claims
func (c OAuthTokenClaims) GetSubject() (string, error) {
	return c.Subject, nil
}

// parseRSAPrivateKey parsea una clave privada RSA desde PEM
func parseRSAPrivateKey(keyData string) (*rsa.PrivateKey, error) {
	block, _ := pem.Decode([]byte(keyData))
	if block == nil {
		return nil, errors.New("failed to parse PEM block containing the key")
	}

	// Try PKCS1 first
	if key, err := x509.ParsePKCS1PrivateKey(block.Bytes); err == nil {
		return key, nil
	}

	// Try PKCS8
	if key, err := x509.ParsePKCS8PrivateKey(block.Bytes); err == nil {
		if rsaKey, ok := key.(*rsa.PrivateKey); ok {
			return rsaKey, nil
		}
		return nil, errors.New("key is not RSA private key")
	}

	return nil, errors.New("failed to parse RSA private key")
}

// parseRSAPublicKey parsea una clave pública RSA desde PEM
func parseRSAPublicKey(keyData string) (*rsa.PublicKey, error) {
	block, _ := pem.Decode([]byte(keyData))
	if block == nil {
		return nil, errors.New("failed to parse PEM block containing the key")
	}

	// Try parsing as public key
	if key, err := x509.ParsePKIXPublicKey(block.Bytes); err == nil {
		if rsaKey, ok := key.(*rsa.PublicKey); ok {
			return rsaKey, nil
		}
		return nil, errors.New("key is not RSA public key")
	}

	return nil, errors.New("failed to parse RSA public key")
}

// GenerateJWTToken genera un token JWT con los datos proporcionados compatible con Laravel Passport
func GenerateJWTToken(userID string, clientID string, tokenID string, scopes []string, isRefresh bool) (string, time.Time, error) {
	var expirationSeconds int64
	var expirationEnvVar string

	if isRefresh {
		expirationEnvVar = os.Getenv("OAUTH_REFRESH_TOKEN_LIFETIME")
		if expirationEnvVar == "" {
			expirationSeconds = 1209600 // 2 semanas por defecto
		}
	} else {
		expirationEnvVar = os.Getenv("OAUTH_ACCESS_TOKEN_LIFETIME")
		if expirationEnvVar == "" {
			expirationSeconds = 31536000 // 1 año por defecto (como Laravel Passport)
		}
	}

	if expirationEnvVar != "" {
		var err error
		expirationSeconds, err = strconv.ParseInt(expirationEnvVar, 10, 64)
		if err != nil {
			return "", time.Time{}, err
		}
	}

	now := time.Now()
	expirationTime := now.Add(time.Second * time.Duration(expirationSeconds))

	// Claims compatibles con Laravel Passport
	claims := OAuthTokenClaims{
		Audience:  clientID,
		JTI:       tokenID,
		IssuedAt:  float64(now.Unix()) + float64(now.Nanosecond())/1e9,
		NotBefore: float64(now.Unix()) + float64(now.Nanosecond())/1e9,
		ExpiresAt: float64(expirationTime.Unix()) + float64(expirationTime.Nanosecond())/1e9,
		Subject:   userID,
		Scopes:    scopes,
	}

	// Obtener la clave privada RSA
	privateKeyPEM := getPrivateKey()
	if privateKeyPEM == "" {
		return "", time.Time{}, fmt.Errorf("OAUTH_PRIVATE_KEY no está configurado y no se pudo leer desde archivo")
	}

	privateKey, err := parseRSAPrivateKey(privateKeyPEM)
	if err != nil {
		return "", time.Time{}, fmt.Errorf("error parsing private key: %v", err)
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)

	// Agregar header typ como JWT (Laravel Passport incluye esto)
	token.Header["typ"] = "JWT"

	tokenString, err := token.SignedString(privateKey)
	if err != nil {
		return "", time.Time{}, err
	}

	return tokenString, expirationTime, nil
}

// ValidateJWTToken valida un token JWT y devuelve sus claims compatible con Laravel Passport
func ValidateJWTToken(tokenString string) (*OAuthTokenClaims, error) {
	publicKeyPEM := getPublicKey()
	if publicKeyPEM == "" {
		return nil, fmt.Errorf("OAUTH_PUBLIC_KEY no está configurado y no se pudo leer desde archivo")
	}

	publicKey, err := parseRSAPublicKey(publicKeyPEM)
	if err != nil {
		return nil, fmt.Errorf("error parsing public key: %v", err)
	}

	token, err := jwt.ParseWithClaims(tokenString, &OAuthTokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("método de firma inesperado: %v", token.Header["alg"])
		}
		return publicKey, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*OAuthTokenClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("token inválido")
}

// getPrivateKey obtiene la clave privada desde variable de entorno o archivo
func getPrivateKey() string {
	// Intentar primero desde variable de entorno
	privateKey := os.Getenv("OAUTH_PRIVATE_KEY")
	if privateKey != "" {
		return privateKey
	}

	// Si no existe en env, intentar leer desde archivo
	keyPath := os.Getenv("OAUTH_PRIVATE_KEY_PATH")
	if keyPath == "" {
		keyPath = "storage/oauth/oauth-private.key"
	}

	keyBytes, err := os.ReadFile(keyPath)
	if err != nil {
		return ""
	}

	return string(keyBytes)
}

// getPublicKey obtiene la clave pública desde variable de entorno o archivo
func getPublicKey() string {
	// Intentar primero desde variable de entorno
	publicKey := os.Getenv("OAUTH_PUBLIC_KEY")
	if publicKey != "" {
		return publicKey
	}

	// Si no existe en env, intentar leer desde archivo
	keyPath := os.Getenv("OAUTH_PUBLIC_KEY_PATH")
	if keyPath == "" {
		keyPath = "storage/oauth/oauth-public.key"
	}

	keyBytes, err := os.ReadFile(keyPath)
	if err != nil {
		return ""
	}

	return string(keyBytes)
}

// GenerateRandomToken genera un token aleatorio para usar como identificador único (compatible con Laravel Passport)
func GenerateRandomToken(length int) (string, error) {
	// Laravel Passport usa tokens de 80 caracteres hexadecimales
	bytes := make([]byte, 40) // 40 bytes = 80 hex chars
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

// HasScope verifica si un conjunto de scopes incluye un scope específico
func HasScope(tokenScopes []string, requiredScope string) bool {
	return slices.Contains(tokenScopes, requiredScope)
}
