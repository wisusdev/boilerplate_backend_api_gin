package middleware

import (
	"boilerplate_backend_api_gin/core/helpers"
	"bytes"
	"fmt"
	"io"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

// Structs para representar el documento JSON API
type JSONAPIDocument struct {
	Data Data `json:"data" validate:"required"`
}

type Data struct {
	Type       string                 `json:"type" validate:"required"`
	ID         string                 `json:"id"`
	Attributes map[string]interface{} `json:"attributes"`
}

func ValidateJSONAPIDocument() gin.HandlerFunc {
	validate := validator.New()

	return func(context *gin.Context) {
		method := context.Request.Method

		// Solo aplicamos validación para métodos que envían cuerpo
		if method == "POST" || method == "PATCH" || method == "PUT" {

			// Leer el cuerpo crudo primero
			body, err := context.GetRawData()
			if err != nil {
				context.AbortWithStatusJSON(400, gin.H{"error": "No se pudo leer el cuerpo"})
				return
			}

			// Reset the body for binding
			context.Request.Body = io.NopCloser(bytes.NewReader(body))

			// Leer y deserializar el cuerpo
			var doc JSONAPIDocument
			if err := context.ShouldBindJSON(&doc); err != nil {
				helpers.Logs("ERROR", "Error al parsear JSON: "+err.Error())
				context.AbortWithStatusJSON(400, gin.H{"error": "JSON inválido"})
				return
			}

			// Reset the body again for further reading in controllers
			context.Request.Body = io.NopCloser(bytes.NewReader(body))

			// Validaciones adicionales
			if method == "POST" || method == "PATCH" {
				// Validar data.type
				if err := validate.StructPartial(doc.Data, "Type"); err != nil {
					fmt.Println("Error de validación en data.type:", err)
					context.AbortWithStatusJSON(400, gin.H{"error": "data.type es requerido"})
					return
				}

				// Validar data.attributes si no es una ruta de relationships
				if !strings.Contains(context.Request.URL.Path, "relationships") {
					if doc.Data.Attributes == nil {
						context.AbortWithStatusJSON(400, gin.H{"error": "data.attributes es requerido"})
						return
					}
				}
			}

			if method == "PATCH" {
				// Validar data.id
				if doc.Data.ID == "" {
					context.AbortWithStatusJSON(400, gin.H{"error": "data.id es requerido para PATCH"})
					return
				}
			}
		}

		context.Next()
	}
}
