package helpers

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"os"
	"slices"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// OAuthTokenClaims define la estructura de los claims del token JWT
type OAuthTokenClaims struct {
	jwt.RegisteredClaims
	Scopes []string `json:"scopes,omitempty"`
}

// GenerateJWTToken genera un token JWT con los datos proporcionados
func GenerateJWTToken(userID int64, clientID string, tokenID string, scopes []string, isRefresh bool) (string, time.Time, error) {
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
			expirationSeconds = 86400 // 24 horas por defecto
		}
	}

	if expirationEnvVar != "" {
		var err error
		expirationSeconds, err = strconv.ParseInt(expirationEnvVar, 10, 64)
		if err != nil {
			return "", time.Time{}, err
		}
	}

	expirationTime := time.Now().Add(time.Second * time.Duration(expirationSeconds))

	claims := OAuthTokenClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "semita_api",
			Subject:   fmt.Sprintf("%d", userID),
			Audience:  jwt.ClaimStrings{clientID},
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			NotBefore: jwt.NewNumericDate(time.Now()),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ID:        tokenID,
		},
		Scopes: scopes,
	}

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		return "", time.Time{}, fmt.Errorf("JWT_SECRET no está configurado")
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		return "", time.Time{}, err
	}

	return tokenString, expirationTime, nil
}

// ValidateJWTToken valida un token JWT y devuelve sus claims
func ValidateJWTToken(tokenString string) (*OAuthTokenClaims, error) {
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		return nil, fmt.Errorf("JWT_SECRET no está configurado")
	}

	token, err := jwt.ParseWithClaims(tokenString, &OAuthTokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("método de firma inesperado: %v", token.Header["alg"])
		}
		return []byte(jwtSecret), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*OAuthTokenClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("token inválido")
}

// GenerateRandomToken genera un token aleatorio para usar como identificador único
func GenerateRandomToken(length int) (string, error) {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

// HasScope verifica si un conjunto de scopes incluye un scope específico
func HasScope(tokenScopes []string, requiredScope string) bool {
	return slices.Contains(tokenScopes, requiredScope)
}
