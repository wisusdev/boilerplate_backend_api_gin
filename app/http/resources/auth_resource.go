package resources

// AuthResource estructura para la respuesta de autenticación

type AuthResource struct {
	ID    uint   `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Token string `json:"token,omitempty"`
}

// NewAuthResource construye la respuesta de autenticación
func NewAuthResource(id uint, name, email, token string) AuthResource {
	return AuthResource{
		ID:    id,
		Name:  name,
		Email: email,
		Token: token,
	}
}

type AuthLoginResponse struct {
	Data AuthLoginData `json:"data"`
}

type AuthLoginData struct {
	Type       string         `json:"type"`
	ID         uint           `json:"id"`
	Attributes AuthLoginAttrs `json:"attributes"`
	Meta       AuthLoginMeta  `json:"meta"`
}

type AuthLoginAttrs struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

type AuthLoginMeta struct {
	Token        string      `json:"token"`
	RefreshToken string      `json:"refresh_token"`
	ExpiresIn    int         `json:"expires_in"`
	Scope        interface{} `json:"scope"`
}

func NewAuthLoginResponse(resource AuthResource, refreshToken string, expiresIn int, scope interface{}) AuthLoginResponse {
	return AuthLoginResponse{
		Data: AuthLoginData{
			Type: "users",
			ID:   resource.ID,
			Attributes: AuthLoginAttrs{
				Name:  resource.Name,
				Email: resource.Email,
			},
			Meta: AuthLoginMeta{
				Token:        resource.Token,
				RefreshToken: refreshToken,
				ExpiresIn:    expiresIn,
				Scope:        scope,
			},
		},
	}
}
