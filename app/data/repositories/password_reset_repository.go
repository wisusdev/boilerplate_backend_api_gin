package repositories

import (
	"semita/core/database/database_connections"
	"time"
)

type PasswordReset struct {
	Email     string
	Token     string
	CreatedAt time.Time
}

func CreatePasswordReset(email, token string) error {
	db := database_connections.DatabaseConnectSQL()
	defer db.Close()
	_, err := db.Exec("INSERT INTO password_resets (email, token, created_at) VALUES (?, ?, ?)", email, token, time.Now().Format("2006-01-02 15:04:05"))
	return err
}

func GetPasswordResetByToken(token string) (PasswordReset, error) {
	db := database_connections.DatabaseConnectSQL()
	defer db.Close()

	var pr PasswordReset
	var createdAtStr string

	err := db.QueryRow("SELECT email, token, created_at FROM password_resets WHERE token = ?", token).Scan(&pr.Email, &pr.Token, &createdAtStr)
	if err != nil {
		return pr, err
	}

	loc := time.Local
	pr.CreatedAt, err = time.ParseInLocation("2006-01-02 15:04:05", createdAtStr, loc)
	if err != nil {
		return pr, err
	}

	return pr, nil
}

func DeletePasswordReset(token string) error {
	db := database_connections.DatabaseConnectSQL()
	defer db.Close()
	_, err := db.Exec("DELETE FROM password_resets WHERE token = ?", token)
	return err
}
