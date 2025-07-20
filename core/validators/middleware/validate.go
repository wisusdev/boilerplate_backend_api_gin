package middleware

import (
	"fmt"
	"reflect"
	"semita/core/validators"

	"github.com/gin-gonic/gin"
)

// ValidateRequest middleware que valida automáticamente los requests
func ValidateRequest(requestType interface{}) gin.HandlerFunc {
	return gin.HandlerFunc(func(c *gin.Context) {
		// Crear una nueva instancia del tipo de request
		requestValue := reflect.New(reflect.TypeOf(requestType).Elem())
		request := requestValue.Interface()

		// Verificar si implementa la interfaz Validatable
		if validatable, ok := request.(validators.Validatable); ok {
			if err := validators.Validate(c, validatable); err != nil {
				c.Abort()
				return
			}

			// Almacenar el request validado en el contexto
			c.Set("validated_request", request)
		} else {
			// Si no implementa Validatable, usar validación básica
			if err := c.ShouldBind(request); err != nil {
				c.JSON(422, validators.ValidationResponse{
					Errors: []validators.ValidationErrorResponse{{
						Status: "422",
						Title:  "Validation Error",
						Detail: "Invalid request format",
						Source: validators.ValidationErrorSource{Pointer: "/data"},
						Meta: validators.ValidationErrorMeta{
							Field:   "request",
							Rule:    "format",
							Message: err.Error(),
						},
					}},
				})
				c.Abort()
				return
			}

			c.Set("validated_request", request)
		}

		c.Next()
	})
}

// GetValidatedRequest helper para obtener el request validado del contexto
func GetValidatedRequest(c *gin.Context) (interface{}, bool) {
	return c.Get("validated_request")
}

// ValidateJSON middleware específico para validar JSON requests
func ValidateJSON(requestType interface{}) gin.HandlerFunc {
	return gin.HandlerFunc(func(c *gin.Context) {
		// Verificar que el Content-Type sea JSON
		contentType := c.GetHeader("Content-Type")
		if contentType != "application/json" {
			c.JSON(400, validators.ValidationResponse{
				Errors: []validators.ValidationErrorResponse{{
					Status: "400",
					Title:  "Bad Request",
					Detail: "Content-Type must be application/json",
					Source: validators.ValidationErrorSource{Pointer: "/headers/content-type"},
					Meta: validators.ValidationErrorMeta{
						Field:   "content-type",
						Rule:    "required",
						Message: "Content-Type header is required and must be application/json",
					},
				}},
			})
			c.Abort()
			return
		}

		// Usar el middleware de validación estándar
		ValidateRequest(requestType)(c)
	})
}

// ValidateForm middleware específico para validar form requests
func ValidateForm(requestType interface{}) gin.HandlerFunc {
	return gin.HandlerFunc(func(c *gin.Context) {
		// Verificar que el Content-Type sea form
		contentType := c.GetHeader("Content-Type")
		if contentType != "application/x-www-form-urlencoded" && contentType != "multipart/form-data" {
			c.JSON(400, validators.ValidationResponse{
				Errors: []validators.ValidationErrorResponse{{
					Status: "400",
					Title:  "Bad Request",
					Detail: "Content-Type must be application/x-www-form-urlencoded or multipart/form-data",
					Source: validators.ValidationErrorSource{Pointer: "/headers/content-type"},
					Meta: validators.ValidationErrorMeta{
						Field:   "content-type",
						Rule:    "required",
						Message: "Invalid Content-Type for form request",
					},
				}},
			})
			c.Abort()
			return
		}

		// Usar el middleware de validación estándar
		ValidateRequest(requestType)(c)
	})
}

// ConditionalValidation middleware que aplica validación basada en condiciones
func ConditionalValidation(condition func(c *gin.Context) bool, requestType interface{}) gin.HandlerFunc {
	return gin.HandlerFunc(func(c *gin.Context) {
		if condition(c) {
			ValidateRequest(requestType)(c)
		} else {
			c.Next()
		}
	})
}

// AuthenticatedValidation middleware que solo valida si el usuario está autenticado
func AuthenticatedValidation(requestType interface{}) gin.HandlerFunc {
	return ConditionalValidation(func(c *gin.Context) bool {
		// Verificar si hay un token de autorización o sesión activa
		token := c.GetHeader("Authorization")
		_, sessionExists := c.Get("user")
		return token != "" || sessionExists
	}, requestType)
}

// ValidateWithCallback middleware que permite un callback después de la validación
func ValidateWithCallback(requestType interface{}, callback func(c *gin.Context, request interface{})) gin.HandlerFunc {
	return gin.HandlerFunc(func(c *gin.Context) {
		// Crear una nueva instancia del tipo de request
		requestValue := reflect.New(reflect.TypeOf(requestType).Elem())
		request := requestValue.Interface()

		// Verificar si implementa la interfaz Validatable
		if validatable, ok := request.(validators.Validatable); ok {
			if err := validators.Validate(c, validatable); err != nil {
				c.Abort()
				return
			}

			// Ejecutar callback
			if callback != nil {
				callback(c, request)
			}

			// Almacenar el request validado en el contexto
			c.Set("validated_request", request)
		}

		c.Next()
	})
}

// MultiValidation middleware que permite validar múltiples tipos de request
func MultiValidation(requestTypes map[string]interface{}) gin.HandlerFunc {
	return gin.HandlerFunc(func(c *gin.Context) {
		method := c.Request.Method
		requestType, exists := requestTypes[method]

		if !exists {
			c.JSON(405, validators.ValidationResponse{
				Errors: []validators.ValidationErrorResponse{{
					Status: "405",
					Title:  "Method Not Allowed",
					Detail: fmt.Sprintf("Method %s is not allowed for this endpoint", method),
					Source: validators.ValidationErrorSource{Pointer: "/method"},
					Meta: validators.ValidationErrorMeta{
						Field:   "method",
						Rule:    "allowed",
						Message: fmt.Sprintf("Method %s is not supported", method),
					},
				}},
			})
			c.Abort()
			return
		}

		ValidateRequest(requestType)(c)
	})
}
