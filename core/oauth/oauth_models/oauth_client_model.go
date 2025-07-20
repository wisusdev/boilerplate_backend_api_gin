package oauth_models

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"semita/core/database/database_connections"
	"strings"
)

type OAuthClient struct {
	ID           int64  `db:"id"`
	Name         string `db:"name"`
	ClientID     string `db:"client_id"`
	ClientSecret string `db:"client_secret"`
	RedirectURI  string `db:"redirect_uri"`
	GrantTypes   string `db:"grant_types"` // Coma separada
	Scopes       string `db:"scopes"`      // Coma separada
	CreatedAt    string `db:"created_at"`
	UpdatedAt    string `db:"updated_at"`
}

// Tabla de clientes OAuth
const oauthClientTable = "oauth_clients"

// GetClientByID obtiene un cliente OAuth por su ID
func GetClientByID(id int64) (*OAuthClient, error) {
	db := database_connections.DatabaseConnectSQL()
	defer db.Close()

	query := `SELECT id, name, client_id, client_secret, redirect_uri, grant_types, scopes, 
              created_at, updated_at FROM ` + oauthClientTable + ` WHERE id = ?`

	var client OAuthClient
	err := db.QueryRow(query, id).Scan(
		&client.ID, &client.Name, &client.ClientID, &client.ClientSecret,
		&client.RedirectURI, &client.GrantTypes, &client.Scopes,
		&client.CreatedAt, &client.UpdatedAt)

	if err != nil {
		return nil, err
	}

	return &client, nil
}

// GetClientByClientID obtiene un cliente OAuth por su client_id
func GetClientByClientID(clientID string) (*OAuthClient, error) {
	db := database_connections.DatabaseConnectSQL()
	defer db.Close()

	query := `SELECT id, name, client_id, client_secret, redirect_uri, grant_types, scopes, 
              created_at, updated_at FROM ` + oauthClientTable + ` WHERE client_id = ?`

	var client OAuthClient
	err := db.QueryRow(query, clientID).Scan(
		&client.ID, &client.Name, &client.ClientID, &client.ClientSecret,
		&client.RedirectURI, &client.GrantTypes, &client.Scopes,
		&client.CreatedAt, &client.UpdatedAt)

	if err != nil {
		return nil, err
	}

	return &client, nil
}

// GetAllClients obtiene todos los clientes OAuth
func GetAllClients() ([]OAuthClient, error) {
	db := database_connections.DatabaseConnectSQL()
	defer db.Close()

	query := `SELECT id, name, client_id, client_secret, redirect_uri, grant_types, scopes, 
              created_at, updated_at FROM ` + oauthClientTable

	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var clients []OAuthClient

	for rows.Next() {
		var client OAuthClient
		err := rows.Scan(
			&client.ID, &client.Name, &client.ClientID, &client.ClientSecret,
			&client.RedirectURI, &client.GrantTypes, &client.Scopes,
			&client.CreatedAt, &client.UpdatedAt)
		if err != nil {
			return nil, err
		}
		clients = append(clients, client)
	}

	return clients, nil
}

// CreateClient crea un nuevo cliente OAuth
func CreateClient(name, redirectURI, grantTypes, scopes string) (*OAuthClient, error) {
	// Generar client_id y client_secret aleatorios
	clientID, err := generateSecureToken(16)
	if err != nil {
		return nil, err
	}

	clientSecret, err := generateSecureToken(32)
	if err != nil {
		return nil, err
	}

	db := database_connections.DatabaseConnectSQL()
	defer db.Close()

	query := `INSERT INTO ` + oauthClientTable + ` 
              (name, client_id, client_secret, redirect_uri, grant_types, scopes) 
              VALUES (?, ?, ?, ?, ?, ?)`

	result, err := db.Exec(query, name, clientID, clientSecret, redirectURI, grantTypes, scopes)
	if err != nil {
		return nil, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	return GetClientByID(id)
}

// CreateOAuthClient crea un cliente OAuth con client_id y client_secret personalizados
func CreateOAuthClient(name, clientID, clientSecret string) error {
	db := database_connections.DatabaseConnectSQL()
	defer db.Close()

	query := `INSERT INTO ` + oauthClientTable + ` (name, client_id, client_secret, redirect_uri, grant_types, scopes) VALUES (?, ?, ?, '', 'password,refresh_token', '*')`
	_, err := db.Exec(query, name, clientID, clientSecret)
	return err
}

// UpdateClient actualiza un cliente OAuth existente
func UpdateClient(id int64, name, redirectURI, grantTypes, scopes string) (*OAuthClient, error) {
	db := database_connections.DatabaseConnectSQL()
	defer db.Close()

	query := `UPDATE ` + oauthClientTable + ` 
              SET name = ?, redirect_uri = ?, grant_types = ?, scopes = ? 
              WHERE id = ?`

	_, err := db.Exec(query, name, redirectURI, grantTypes, scopes, id)
	if err != nil {
		return nil, err
	}

	return GetClientByID(id)
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
func ValidateClientCredentials(clientID, clientSecret string) (*OAuthClient, error) {
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
