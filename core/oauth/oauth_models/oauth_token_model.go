package oauth_models

import (
	"errors"
	"semita/core/database/database_connections"
	"semita/core/helpers"
	"strings"
	"time"
)

type OAuthToken struct {
	ID           int64  `db:"id"`
	UserID       int64  `db:"user_id"`
	ClientID     int64  `db:"client_id"`
	AccessToken  string `db:"access_token"`
	RefreshToken string `db:"refresh_token"`
	Scopes       string `db:"scopes"` // Coma separada
	Revoked      bool   `db:"revoked"`
	ExpiresAt    string `db:"expires_at"`
	CreatedAt    string `db:"created_at"`
	UpdatedAt    string `db:"updated_at"`
}

// Tabla de tokens OAuth
const oauthTokenTable = "oauth_tokens"

// GetTokenByAccessToken obtiene un token por su access_token
func GetTokenByAccessToken(accessToken string) (*OAuthToken, error) {
	database := database_connections.DatabaseConnectSQL()
	defer database.Close()

	query := `SELECT id, user_id, client_id, access_token, refresh_token, 
              scopes, revoked, expires_at, created_at, updated_at 
              FROM ` + oauthTokenTable + ` 
              WHERE access_token = ? AND revoked = 0`

	var token OAuthToken
	err := database.QueryRow(query, accessToken).Scan(
		&token.ID, &token.UserID, &token.ClientID,
		&token.AccessToken, &token.RefreshToken, &token.Scopes,
		&token.Revoked, &token.ExpiresAt, &token.CreatedAt, &token.UpdatedAt)

	if err != nil {
		return nil, err
	}

	return &token, nil
}

// GetTokenByRefreshToken obtiene un token por su refresh_token
func GetTokenByRefreshToken(refreshToken string) (*OAuthToken, error) {
	database := database_connections.DatabaseConnectSQL()
	defer database.Close()

	query := `SELECT id, user_id, client_id, access_token, refresh_token, 
              scopes, revoked, expires_at, created_at, updated_at 
              FROM ` + oauthTokenTable + ` 
              WHERE refresh_token = ? AND revoked = 0`

	var token OAuthToken
	err := database.QueryRow(query, refreshToken).Scan(
		&token.ID, &token.UserID, &token.ClientID,
		&token.AccessToken, &token.RefreshToken, &token.Scopes,
		&token.Revoked, &token.ExpiresAt, &token.CreatedAt, &token.UpdatedAt)

	if err != nil {
		return nil, err
	}

	return &token, nil
}

// CreateToken crea un nuevo token de acceso
func CreateToken(userID int64, clientID int64, scopes string) (*OAuthToken, error) {
	database := database_connections.DatabaseConnectSQL()
	defer database.Close()

	// Obtener el cliente para el ID
	client, err := GetClientByID(clientID)
	if err != nil {
		return nil, err
	}

	// Generar identificadores únicos para los tokens
	accessTokenId, err := helpers.GenerateRandomToken(16)
	if err != nil {
		return nil, err
	}

	refreshTokenId, err := helpers.GenerateRandomToken(16)
	if err != nil {
		return nil, err
	}

	// Convertir scopes de string a slice
	scopesSlice := []string{}
	if scopes != "" {
		scopesSlice = strings.Split(scopes, ",")
	}

	// Generar token de acceso JWT
	accessTokenString, expiresAt, err := helpers.GenerateJWTToken(userID, client.ClientID, accessTokenId, scopesSlice, false)
	if err != nil {
		return nil, err
	}

	// Generar token de refresco JWT
	refreshTokenString, _, err := helpers.GenerateJWTToken(userID, client.ClientID, refreshTokenId, scopesSlice, true)
	if err != nil {
		return nil, err
	}

	// Insertar token en la base de datos
	query := `INSERT INTO ` + oauthTokenTable + ` 
              (user_id, client_id, access_token, refresh_token, scopes, revoked, expires_at) 
              VALUES (?, ?, ?, ?, ?, 0, ?)`

	result, err := database.Exec(query, userID, clientID, accessTokenString, refreshTokenString, scopes, expiresAt.Format("2006-01-02 15:04:05"))
	if err != nil {
		return nil, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	// Recuperar el token creado
	return getTokenByID(database, id)
}

// RefreshToken renueva un token usando el refresh_token
func RefreshToken(refreshToken string) (*OAuthToken, error) {
	// Validar el refresh token
	_, err := helpers.ValidateJWTToken(refreshToken)
	if err != nil {
		return nil, err
	}

	database := database_connections.DatabaseConnectSQL()
	defer database.Close()

	// Buscar el token original
	existingToken, err := GetTokenByRefreshToken(refreshToken)
	if err != nil {
		return nil, err
	}

	// Verificar que no haya sido revocado
	if existingToken.Revoked {
		return nil, errors.New("el token ha sido revocado")
	}

	// Revocar el token antiguo
	_, err = database.Exec("UPDATE "+oauthTokenTable+" SET revoked = 1 WHERE id = ?", existingToken.ID)
	if err != nil {
		return nil, err
	}

	// Crear un nuevo token
	return CreateToken(existingToken.UserID, existingToken.ClientID, existingToken.Scopes)
}

// RevokeToken revoca un token específico
func RevokeToken(accessToken string) error {
	database := database_connections.DatabaseConnectSQL()
	defer database.Close()

	_, err := database.Exec("UPDATE "+oauthTokenTable+" SET revoked = 1 WHERE access_token = ?", accessToken)
	return err
}

// RevokeAllUserTokens revoca todos los tokens de un usuario
func RevokeAllUserTokens(userID int64) error {
	database := database_connections.DatabaseConnectSQL()
	defer database.Close()

	_, err := database.Exec("UPDATE "+oauthTokenTable+" SET revoked = 1 WHERE user_id = ?", userID)
	return err
}

// IsTokenValid verifica si un token es válido (no expirado y no revocado)
func IsTokenValid(accessToken string) (bool, error) {
	token, err := GetTokenByAccessToken(accessToken)
	if err != nil {
		return false, err
	}

	// Verificar que no haya sido revocado
	if token.Revoked {
		return false, nil
	}

	// Verificar que no haya expirado
	expiresAt, err := time.Parse("2006-01-02 15:04:05", token.ExpiresAt)
	if err != nil {
		return false, err
	}

	return time.Now().Before(expiresAt), nil
}

// GetScopesArray devuelve los scopes como un array
func (t *OAuthToken) GetScopesArray() []string {
	if t.Scopes == "" {
		return []string{}
	}
	return strings.Split(t.Scopes, ",")
}

// HasScope verifica si el token tiene un scope específico
func (t *OAuthToken) HasScope(requiredScope string) bool {
	scopes := t.GetScopesArray()
	for _, scope := range scopes {
		if scope == requiredScope {
			return true
		}
	}
	return false
}

// Función auxiliar para obtener un token por ID
func getTokenByID(database database_connections.SQLAdapter, id int64) (*OAuthToken, error) {
	query := `SELECT id, user_id, client_id, access_token, refresh_token, 
              scopes, revoked, expires_at, created_at, updated_at 
              FROM ` + oauthTokenTable + ` WHERE id = ?`

	var token OAuthToken
	err := database.QueryRow(query, id).Scan(
		&token.ID, &token.UserID, &token.ClientID,
		&token.AccessToken, &token.RefreshToken, &token.Scopes,
		&token.Revoked, &token.ExpiresAt, &token.CreatedAt, &token.UpdatedAt)

	if err != nil {
		return nil, err
	}

	return &token, nil
}
