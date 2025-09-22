package resources

import (
	userModels "boilerplate_backend_api_gin/app/data/models"
	rpModels "boilerplate_backend_api_gin/core/roles_and_permissions/models"
)

// AuthResource estructura para la respuesta de autenticación
type AuthResource struct {
	ID   string         `json:"id"`
	User UserAttributes `json:"user"`
}

// UserAttributes estructura para los atributos del usuario
type UserAttributes struct {
	FirstName string      `json:"first_name"`
	LastName  string      `json:"last_name"`
	Username  string      `json:"username"`
	Email     string      `json:"email"`
	Avatar    interface{} `json:"avatar"`
	Language  string      `json:"language"`
}

// NewAuthResource construye la respuesta de autenticación
func NewAuthResource(user userModels.UserStruct) AuthResource {
	var avatar interface{}
	if user.Avatar.Valid {
		avatar = user.Avatar.String
	} else {
		avatar = nil
	}

	return AuthResource{
		ID: user.ID,
		User: UserAttributes{
			FirstName: user.FirstName,
			LastName:  user.LastName,
			Username:  user.Username,
			Email:     user.Email,
			Avatar:    avatar,
			Language:  user.Language,
		},
	}
}

// AuthLoginResponse estructura para la respuesta completa de login
type AuthLoginResponse struct {
	Data AuthLoginData `json:"data"`
}

type AuthLoginData struct {
	Type          string                 `json:"type"`
	ID            string                 `json:"id"`
	Attributes    AuthLoginAttrs         `json:"attributes"`
	Relationships AuthLoginRelationships `json:"relationships"`
}

type AuthLoginAttrs struct {
	User UserAttributes `json:"user"`
}

type AuthLoginRelationships struct {
	Roles       []string        `json:"roles"`
	Permissions []string        `json:"permissions"`
	Access      AuthLoginAccess `json:"access"`
}

type AuthLoginAccess struct {
	Token     string `json:"token"`
	TokenType string `json:"token_type"`
	ExpiresAt string `json:"expires_at"`
}

// NewAuthLoginResponse construye la respuesta completa de login
func NewAuthLoginResponse(user userModels.UserStruct, roles []rpModels.RoleStruct, permissions []rpModels.PermissionStruct, token string, expiresAt string) AuthLoginResponse {
	// Convertir roles a slice de strings
	var roleNames []string
	for _, role := range roles {
		roleNames = append(roleNames, role.Name)
	}

	// Convertir permisos a slice de strings
	var permissionNames []string
	for _, permission := range permissions {
		permissionNames = append(permissionNames, permission.Name)
	}

	resource := NewAuthResource(user)

	return AuthLoginResponse{
		Data: AuthLoginData{
			Type: "users",
			ID:   user.ID,
			Attributes: AuthLoginAttrs{
				User: resource.User,
			},
			Relationships: AuthLoginRelationships{
				Roles:       roleNames,
				Permissions: permissionNames,
				Access: AuthLoginAccess{
					Token:     token,
					TokenType: "Bearer",
					ExpiresAt: expiresAt,
				},
			},
		},
	}
}
