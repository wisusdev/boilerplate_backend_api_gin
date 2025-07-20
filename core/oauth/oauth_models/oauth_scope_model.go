package oauth_models

import "semita/core/database/database_connections"

type OAuthScope struct {
	ID          int64  `db:"id"`
	Name        string `db:"name"`
	Description string `db:"description"`
	CreatedAt   string `db:"created_at"`
	UpdatedAt   string `db:"updated_at"`
}

// Tabla de scopes OAuth
const oauthScopeTable = "oauth_scopes"

// GetScopeByName obtiene un scope por su nombre
func GetScopeByName(name string) (*OAuthScope, error) {
	db := database_connections.DatabaseConnectSQL()
	defer db.Close()

	query := `SELECT id, name, description, created_at, updated_at 
              FROM ` + oauthScopeTable + ` WHERE name = ?`

	var scope OAuthScope
	err := db.QueryRow(query, name).Scan(
		&scope.ID, &scope.Name, &scope.Description,
		&scope.CreatedAt, &scope.UpdatedAt)

	if err != nil {
		return nil, err
	}

	return &scope, nil
}

// GetAllScopes obtiene todos los scopes
func GetAllScopes() ([]OAuthScope, error) {
	db := database_connections.DatabaseConnectSQL()
	defer db.Close()

	query := `SELECT id, name, description, created_at, updated_at FROM ` + oauthScopeTable

	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var scopes []OAuthScope

	for rows.Next() {
		var scope OAuthScope
		err := rows.Scan(
			&scope.ID, &scope.Name, &scope.Description,
			&scope.CreatedAt, &scope.UpdatedAt)
		if err != nil {
			return nil, err
		}
		scopes = append(scopes, scope)
	}

	return scopes, nil
}

// CreateScope crea un nuevo scope
func CreateScope(name, description string) (*OAuthScope, error) {
	db := database_connections.DatabaseConnectSQL()
	defer db.Close()

	query := `INSERT INTO ` + oauthScopeTable + ` (name, description) VALUES (?, ?)`

	result, err := db.Exec(query, name, description)
	if err != nil {
		return nil, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	// Recuperar el scope creado
	return GetScopeByID(id)
}

// UpdateScope actualiza un scope existente
func UpdateScope(id int64, name, description string) (*OAuthScope, error) {
	db := database_connections.DatabaseConnectSQL()
	defer db.Close()

	query := `UPDATE ` + oauthScopeTable + ` SET name = ?, description = ? WHERE id = ?`

	_, err := db.Exec(query, name, description, id)
	if err != nil {
		return nil, err
	}

	// Recuperar el scope actualizado
	return GetScopeByID(id)
}

// DeleteScope elimina un scope
func DeleteScope(id int64) error {
	db := database_connections.DatabaseConnectSQL()
	defer db.Close()

	_, err := db.Exec("DELETE FROM "+oauthScopeTable+" WHERE id = ?", id)
	return err
}

// GetScopeByID obtiene un scope por su ID
func GetScopeByID(id int64) (*OAuthScope, error) {
	db := database_connections.DatabaseConnectSQL()
	defer db.Close()

	query := `SELECT id, name, description, created_at, updated_at 
              FROM ` + oauthScopeTable + ` WHERE id = ?`

	var scope OAuthScope
	err := db.QueryRow(query, id).Scan(
		&scope.ID, &scope.Name, &scope.Description,
		&scope.CreatedAt, &scope.UpdatedAt)

	if err != nil {
		return nil, err
	}

	return &scope, nil
}

// ValidateScopes verifica que todos los scopes proporcionados existan
func ValidateScopes(scopes []string) (bool, error) {
	if len(scopes) == 0 {
		return true, nil
	}

	db := database_connections.DatabaseConnectSQL()
	defer db.Close()

	for _, scope := range scopes {
		var count int
		err := db.QueryRow("SELECT COUNT(*) FROM "+oauthScopeTable+" WHERE name = ?", scope).Scan(&count)
		if err != nil {
			return false, err
		}
		if count == 0 {
			return false, nil
		}
	}

	return true, nil
}
