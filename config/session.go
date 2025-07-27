package config

type Session struct {
	Lifetime int    `json:"lifetime"`  // Duración de la sesión en minutos
	Driver   string `json:"driver"`    // Controlador de sesión (ej. "file", "redis", etc.)
	Cookie   string `json:"cookie"`    // Nombre de la cookie de sesión
	Domain   string `json:"domain"`    // Dominio para la cookie de sesión
	Path     string `json:"path"`      // Ruta para la cookie de sesión
	Secure   bool   `json:"secure"`    // Si la cookie es segura (HTTPS)
	HttpOnly bool   `json:"http_only"` // Si la cookie es accesible solo por HTTP
	SameSite string `json:"same_site"` // Política SameSite para la cookie ("Strict", "Lax", "None")
}

func SessionConfig() *Session {
	return &Session{
		Lifetime: GetEnvInt("SESSION_LIFETIME", 120), // 120 minutos por defecto
		Driver:   GetEnv("SESSION_DRIVER", "file"),
		Cookie:   GetEnv("SESSION_COOKIE", "session_id"),
		Domain:   GetEnv("SESSION_DOMAIN", ""),
		Path:     GetEnv("SESSION_PATH", "/"),
		Secure:   GetEnvBool("SESSION_SECURE", false),
		HttpOnly: GetEnvBool("SESSION_HTTP_ONLY", true),
		SameSite: GetEnv("SESSION_SAME_SITE", "Lax"),
	}
}
