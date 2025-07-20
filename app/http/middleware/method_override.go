package middleware

import (
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
)

// Middleware compatible con Gin
func MethodOverride() gin.HandlerFunc {
	return func(context *gin.Context) {
		if context.Request.Method == "POST" {
			// Parsear el formulario
			if err := context.Request.ParseForm(); err != nil {
				fmt.Printf("Error parsing form: %v\n", err)
				context.AbortWithStatusJSON(400, gin.H{"error": "Error parsing form data"})
				return
			}

			method := context.Request.FormValue("_method")
			if method != "" {
				method = strings.ToUpper(method)
				if method == "PUT" || method == "PATCH" || method == "DELETE" {
					context.Request.Method = method
				}
			}
		}

		context.Next()
	}
}
