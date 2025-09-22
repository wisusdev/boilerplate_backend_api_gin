package providers

import (
	"boilerplate_backend_api_gin/app/data/repositories"
)

// CreatePasswordReset creates a new password reset token
func CreatePasswordReset(email, token string) error {
	return repositories.CreatePasswordReset(email, token)
}

// GetPasswordResetByToken retrieves a password reset by token
func GetPasswordResetByToken(token string) (repositories.PasswordReset, error) {
	return repositories.GetPasswordResetByToken(token)
}

// DeletePasswordReset deletes a password reset token
func DeletePasswordReset(token string) error {
	return repositories.DeletePasswordReset(token)
}
