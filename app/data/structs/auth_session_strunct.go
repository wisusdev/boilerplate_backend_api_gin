package structs

type AuthSessionStruct struct {
	User            UserStruct
	IsAuthenticated bool
	AlertId         string
	AlertMessage    string
	Title           string
	Data            any
	Lang            string              // Idioma actual
	Translate       func(string) string // Función de traducción automática
}
