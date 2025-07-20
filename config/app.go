package config

type App struct {
	Name     string `json:"name"`
	Env      string `json:"env"`
	Key      string `json:"key"`
	Debug    bool   `json:"debug"`
	Timezone string `json:"timezone"`
	Url      string `json:"url"`
	Lang     string `json:"lang"`
	Layout   string `json:"layout"`
}

func AppConfig() *App {
	return &App{
		Name:     GetEnv("APP_NAME", ""),
		Env:      GetEnv("APP_ENV", "production"),
		Key:      GetEnv("APP_KEY", ""),
		Debug:    GetEnvBool("APP_DEBUG", false),
		Timezone: GetEnv("APP_TIMEZONE", "UTC"),
		Url:      GetEnv("APP_URL", "http://localhost"),
		Lang:     GetEnv("APP_LANG", "en"),
		Layout:   "resources/views/layouts/app.html",
	}
}
