package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// LanguageMiddleware detecta el idioma desde cookie, query o header y lo guarda en el contexto
func LanguageMiddleware() gin.HandlerFunc {
	return func(context *gin.Context) {
		lang := context.Query("lang")

		if lang == "" {
			langCookie, err := context.Cookie("lang")
			if err == nil {
				lang = langCookie
			}
		}

		if lang == "" {
			lang = context.GetHeader("Accept-Language")
			if len(lang) > 2 {
				lang = lang[:2]
			}
		}

		if lang != "es" && lang != "en" {
			lang = "es"
		}

		context.Set("lang", lang)
		context.Next()
	}
}

// SetLangHandler permite cambiar el idioma desde el dropdown
func SetLangHandler(context *gin.Context) {
	lang := context.PostForm("lang")

	if lang != "es" && lang != "en" {
		lang = "es"
	}

	http.SetCookie(context.Writer, &http.Cookie{
		Name:   "lang",
		Value:  lang,
		Path:   "/",
		MaxAge: 60 * 60 * 24 * 365, // 1 año
	})

	context.Redirect(http.StatusSeeOther, context.Request.Referer())
}
