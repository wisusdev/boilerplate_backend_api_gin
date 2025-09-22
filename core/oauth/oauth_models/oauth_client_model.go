package oauth_models

import (
	"boilerplate_backend_api_gin/core/database/database_connections"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"strings"
)

type OAuthClient struct {
	ClientID     int64  `db:"id"`
	Name         string `db:"name"`
	ClientSecret string `db:"secret"`
	RedirectURI  string `db:"redirect"`
	GrantTypes   string `db:"grant_types"`
	Scopes       string `db:"scopes"`
	CreatedAt    string `db:"created_at"`
	UpdatedAt    string `db:"updated_at"`
}

// Tabla de clientes OAuth
const oauthClientTable = "oauth_clients"

// GetClientByID obtiene un cliente OAuth por su ID
func GetClientByID(id int64) (*OAuthClient, error) {
	db := database_connections.DatabaseConnectSQL()
	defer db.Close()

	query := `SELECT id, name, secret, redirect, created_at, updated_at FROM ` + oauthClientTable + ` WHERE id = ?`

	var client OAuthClient
	err := db.QueryRow(query, id).Scan(
		&client.ClientID, &client.Name, &client.ClientSecret,
		&client.RedirectURI, &client.CreatedAt, &client.UpdatedAt)

	if err != nil {
		return nil, err
	}

	return &client, nil
}

// GetClientByClientID obtiene un cliente OAuth por su client_id
func GetClientByClientID(clientID int64) (*OAuthClient, error) {
	db := database_connections.DatabaseConnectSQL()
	defer db.Close()

	query := `SELECT id, name, secret, redirect, created_at, updated_at FROM ` + oauthClientTable + ` WHERE id = ?`

	var client OAuthClient
	err := db.QueryRow(query, clientID).Scan(
		&client.ClientID, &client.Name, &client.ClientSecret,
		&client.RedirectURI, &client.CreatedAt, &client.UpdatedAt)

	if err != nil {
		return nil, err
	}

	return &client, nil
}

// GetAllClients obtiene todos los clientes OAuth
func GetAllClients() ([]OAuthClient, error) {
	db := database_connections.DatabaseConnectSQL()
	defer db.Close()

	query := `SELECT id, name, secret, redirect, created_at, updated_at FROM ` + oauthClientTable

	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var clients []OAuthClient

	for rows.Next() {
		var client OAuthClient
		err := rows.Scan(
			&client.ClientID, &client.Name, &client.ClientSecret,
			&client.RedirectURI, &client.CreatedAt, &client.UpdatedAt)
		if err != nil {
			return nil, err
		}
		clients = append(clients, client)
	}

	return clients, nil
}

// CreateClient crea un nuevo cliente OAuth
func CreateClient(name, redirectURI string) (*OAuthClient, error) {
	// Generar client_secret aleatorio
	clientSecret, err := generateSecureToken(32)
	if err != nil {
		return nil, err
	}

	db := database_connections.DatabaseConnectSQL()
	defer db.Close()

	query := `INSERT INTO ` + oauthClientTable + ` 
              (name, secret, redirect, personal_access_client, password_client, revoked, created_at, updated_at) 
              VALUES (?, ?, ?, 0, 1, 0, NOW(), NOW())`

	result, err := db.Exec(query, name, clientSecret, redirectURI)
	if err != nil {
		return nil, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	return GetClientByClientID(id)
}

// CreateOAuthClient crea un cliente OAuth con client_id y client_secret personalizados
func CreateOAuthClient(clientID int64, name, clientSecret string) error {
	db := database_connections.DatabaseConnectSQL()
	defer db.Close()

	query := `INSERT INTO ` + oauthClientTable + ` (id, name, secret, redirect, personal_access_client, password_client, revoked, created_at, updated_at) VALUES (?, ?, ?, '', 0, 1, 0, NOW(), NOW())`
	_, err := db.Exec(query, clientID, name, clientSecret)
	return err
}

// UpdateClient actualiza un cliente OAuth existente
func UpdateClient(clientID int64, name, redirectURI string) (*OAuthClient, error) {
	db := database_connections.DatabaseConnectSQL()
	defer db.Close()

	query := `UPDATE ` + oauthClientTable + ` 
              SET name = ?, redirect = ? 
              WHERE id = ?`

	_, err := db.Exec(query, name, redirectURI, clientID)
	if err != nil {
		return nil, err
	}

	return GetClientByClientID(clientID)
}

// DeleteClient elimina un cliente OAuth
func DeleteClient(id int64) error {
	db := database_connections.DatabaseConnectSQL()
	defer db.Close()

	// Primero eliminamos los tokens asociados a este cliente
	_, err := db.Exec("DELETE FROM oauth_tokens WHERE client_id = ?", id)
	if err != nil {
		return err
	}

	// Luego eliminamos el cliente
	_, err = db.Exec("DELETE FROM "+oauthClientTable+" WHERE id = ?", id)
	return err
}

// ValidateClientCredentials valida las credenciales de un cliente
func ValidateClientCredentials(clientID int64, clientSecret string) (*OAuthClient, error) {
	client, err := GetClientByClientID(clientID)
	if err != nil {
		return nil, errors.New("cliente no encontrado")
	}

	if client.ClientSecret != clientSecret {
		return nil, errors.New("credenciales de cliente inválidas")
	}

	return client, nil
}

// SupportsGrantType verifica si un cliente soporta un tipo de grant específico
func (c *OAuthClient) SupportsGrantType(grantType string) bool {
	grantTypes := strings.Split(c.GrantTypes, ",")
	for _, gt := range grantTypes {
		if strings.TrimSpace(gt) == grantType {
			return true
		}
	}
	return false
}

// GetScopesArray devuelve los scopes como un array
func (c *OAuthClient) GetScopesArray() []string {
	if c.Scopes == "" {
		return []string{}
	}
	return strings.Split(c.Scopes, ",")
}

// generateSecureToken genera un token aleatorio seguro
func generateSecureToken(length int) (string, error) {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}
