package repositories

import (
	"boilerplate_backend_api_gin/core/database/database_connections"
	"boilerplate_backend_api_gin/core/helpers"
	"time"
)

type PasswordReset struct {
	Email     string
	Token     string
	CreatedAt time.Time
}

var PasswordResetTable = "password_reset_tokens"

func CreatePasswordReset(email, token string) error {
	db := database_connections.DatabaseConnectSQL()
	defer db.Close()
	_, err := db.Exec("INSERT INTO "+PasswordResetTable+" (email, token, created_at) VALUES (?, ?, ?)", email, token, time.Now().Format("2006-01-02 15:04:05"))
	return err
}

func GetPasswordResetByToken(token string) (PasswordReset, error) {
	db := database_connections.DatabaseConnectSQL()
	defer db.Close()

	var pr PasswordReset
	var createdAtStr string

	err := db.QueryRow("SELECT email, token, created_at FROM "+PasswordResetTable+" WHERE token = ?", token).Scan(&pr.Email, &pr.Token, &createdAtStr)
	if err != nil {
		helpers.Logs("ERROR", "Error retrieving password reset: "+err.Error())
		return pr, err
	}

	loc := time.Local
	pr.CreatedAt, err = time.ParseInLocation("2006-01-02 15:04:05", createdAtStr, loc)
	if err != nil {
		helpers.Logs("ERROR", "Error parsing password reset created at: "+err.Error())
		return pr, err
	}

	return pr, nil
}

func DeletePasswordReset(token string) error {
	db := database_connections.DatabaseConnectSQL()
	defer db.Close()
	_, err := db.Exec("DELETE FROM "+PasswordResetTable+" WHERE token = ?", token)
	return err
}
