package validators

import (
	"encoding/json"
	"fmt"
	"semita/core/database/database_connections"
	"semita/core/helpers"
	"strings"
	"sync"

	"github.com/gin-gonic/gin"
)

// ValidationError representa un error de validación individual
type ValidationError struct {
	Field   string      `json:"field"`
	Rule    string      `json:"rule"`
	Value   interface{} `json:"value"`
	Message string      `json:"message"`
}

// ValidationResult contiene el resultado de una validación
type ValidationResult struct {
	Valid  bool               `json:"valid"`
	Errors []*ValidationError `json:"errors"`
}

// Validator es el validador principal
type Validator struct {
	rules    map[string]*FieldValidator
	messages map[string]string
	db       database_connections.SQLAdapter
	language string // Idioma para traduciones
	mu       sync.RWMutex
}

// FieldValidator maneja las reglas para un campo específico
type FieldValidator struct {
	fieldName string
	rules     []Rule
	sometimes bool
	validator *Validator
}

// Rule representa una regla de validación
type Rule interface {
	Validate(value interface{}, data map[string]interface{}) error
	Name() string
}

// Validatable interface que deben implementar los requests
type Validatable interface {
	Rules() *Validator
	Messages() map[string]string
}

// ValidationResponse formato de respuesta para errores de validación
type ValidationResponse struct {
	Errors []ValidationErrorResponse `json:"errors"`
}

// ValidationErrorResponse formato JSON API para errores
type ValidationErrorResponse struct {
	Status string                `json:"status"`
	Title  string                `json:"title"`
	Detail string                `json:"detail"`
	Source ValidationErrorSource `json:"source"`
	Meta   ValidationErrorMeta   `json:"meta"`
}

type ValidationErrorSource struct {
	Pointer string `json:"pointer"`
}

type ValidationErrorMeta struct {
	Field   string `json:"field"`
	Rule    string `json:"rule"`
	Message string `json:"message"`
}

var (
	defaultValidator *Validator
	once             sync.Once
)

// New crea un nuevo validador
func New() *Validator {
	return &Validator{
		rules:    make(map[string]*FieldValidator),
		messages: make(map[string]string),
	}
}

// Default retorna el validador por defecto (singleton)
func Default() *Validator {
	once.Do(func() {
		defaultValidator = New()
	})
	return defaultValidator
}

// SetDatabase establece la conexión a la base de datos para validaciones
func (v *Validator) SetDatabase(db database_connections.SQLAdapter) *Validator {
	v.mu.Lock()
	defer v.mu.Unlock()
	v.db = db
	return v
}

// SetLanguage establece el idioma para las traducciones
func (v *Validator) SetLanguage(lang string) *Validator {
	v.mu.Lock()
	defer v.mu.Unlock()
	v.language = lang
	return v
}

// Field define las reglas para un campo específico
func (v *Validator) Field(name string) *FieldValidator {
	v.mu.Lock()
	defer v.mu.Unlock()

	if validator, exists := v.rules[name]; exists {
		return validator
	}

	fieldValidator := &FieldValidator{
		fieldName: name,
		rules:     make([]Rule, 0),
		validator: v,
	}
	v.rules[name] = fieldValidator
	return fieldValidator
}

// Messages establece mensajes personalizados
func (v *Validator) Messages(messages map[string]string) *Validator {
	v.mu.Lock()
	defer v.mu.Unlock()
	for key, message := range messages {
		v.messages[key] = message
	}
	return v
}

// Validate valida los datos contra las reglas definidas
func (v *Validator) Validate(data interface{}) *ValidationResult {
	result := &ValidationResult{
		Valid:  true,
		Errors: make([]*ValidationError, 0),
	}

	// Convertir data a map para facilitar el acceso
	dataMap := v.toMap(data)

	v.mu.RLock()
	defer v.mu.RUnlock()

	for fieldName, fieldValidator := range v.rules {
		// Obtener valor usando notación de puntos
		value, exists := v.getNestedValue(dataMap, fieldName)

		// Si el campo tiene la regla "sometimes" y no existe, omitir validación
		if fieldValidator.sometimes && !exists {
			continue
		}

		// Validar reglas en orden, deteniéndose en el primer error
		for _, rule := range fieldValidator.rules {
			if err := rule.Validate(value, dataMap); err != nil {
				validationError := &ValidationError{
					Field:   fieldName,
					Rule:    rule.Name(),
					Value:   value,
					Message: v.getMessage(fieldName, rule.Name(), err.Error()),
				}
				result.Errors = append(result.Errors, validationError)
				result.Valid = false
				// Detener validaciones para este campo al encontrar el primer error
				break
			}
		}
	}

	return result
}

// getNestedValue obtiene un valor de un mapa usando notación de puntos
func (v *Validator) getNestedValue(data map[string]interface{}, fieldPath string) (interface{}, bool) {
	keys := strings.Split(fieldPath, ".")
	current := data

	for i, key := range keys {
		value, exists := current[key]
		if !exists {
			return nil, false
		}

		// Si es el último key, retornar el valor
		if i == len(keys)-1 {
			return value, true
		}

		// Si no es el último key, necesitamos que sea un mapa para continuar
		if nextMap, ok := value.(map[string]interface{}); ok {
			current = nextMap
		} else {
			return nil, false
		}
	}

	return nil, false
}

// getMessage obtiene el mensaje de error personalizado usando el sistema de traducción
func (v *Validator) getMessage(field, rule, defaultMessage string) string {
	// Buscar mensaje específico para campo.regla en mensajes personalizados
	key := fmt.Sprintf("%s.%s", field, rule)
	if message, exists := v.messages[key]; exists {
		return strings.ReplaceAll(message, ":field", field)
	}

	// Buscar mensaje genérico para la regla en mensajes personalizados
	if message, exists := v.messages[rule]; exists {
		return strings.ReplaceAll(message, ":field", field)
	}

	// Usar el sistema de traducción existente
	translationKey := fmt.Sprintf("validation_%s", rule)
	translatedMessage := helpers.Translate(translationKey, v.language)

	// Si encontró una traducción (no devolvió la clave), usarla
	if translatedMessage != translationKey {
		return strings.ReplaceAll(translatedMessage, ":field", field)
	}

	// Fallback al mensaje por defecto
	return defaultMessage
}

// toMap convierte una estructura a un mapa
func (v *Validator) toMap(data interface{}) map[string]interface{} {
	result := make(map[string]interface{})

	// Si ya es un mapa, devolverlo
	if m, ok := data.(map[string]interface{}); ok {
		return m
	}

	// Convertir usando JSON como intermedio
	jsonData, err := json.Marshal(data)
	if err != nil {
		return result
	}

	json.Unmarshal(jsonData, &result)
	return result
}

// Métodos del FieldValidator para sintaxis fluida

// Required marca el campo como requerido
func (f *FieldValidator) Required() *FieldValidator {
	f.rules = append(f.rules, &RequiredRule{})
	return f
}

// Email valida que el campo sea un email válido
func (f *FieldValidator) Email() *FieldValidator {
	f.rules = append(f.rules, &EmailRule{})
	return f
}

// Min establece el valor mínimo
func (f *FieldValidator) Min(min int) *FieldValidator {
	f.rules = append(f.rules, &MinRule{Min: min})
	return f
}

// Max establece el valor máximo
func (f *FieldValidator) Max(max int) *FieldValidator {
	f.rules = append(f.rules, &MaxRule{Max: max})
	return f
}

// Between valida que el valor esté entre min y max
func (f *FieldValidator) Between(min, max int) *FieldValidator {
	f.rules = append(f.rules, &BetweenRule{Min: min, Max: max})
	return f
}

// Unique valida que el valor sea único en la base de datos
func (f *FieldValidator) Unique(table, column string, except ...interface{}) *FieldValidator {
	f.rules = append(f.rules, &UniqueRule{
		Table:  table,
		Column: column,
		Except: except,
		DB:     f.validator.db,
	})
	return f
}

// Exists valida que el valor exista en la base de datos
func (f *FieldValidator) Exists(table, column string) *FieldValidator {
	f.rules = append(f.rules, &ExistsRule{
		Table:  table,
		Column: column,
		DB:     f.validator.db,
	})
	return f
}

// Confirmed valida que el campo coincida con su confirmación
func (f *FieldValidator) Confirmed() *FieldValidator {
	f.rules = append(f.rules, &ConfirmedRule{Field: f.fieldName})
	return f
}

// Regex valida contra una expresión regular
func (f *FieldValidator) Regex(pattern string) *FieldValidator {
	f.rules = append(f.rules, &RegexRule{Pattern: pattern})
	return f
}

// Alpha valida que solo contenga letras
func (f *FieldValidator) Alpha() *FieldValidator {
	f.rules = append(f.rules, &AlphaRule{})
	return f
}

// AlphaNum valida que solo contenga letras y números
func (f *FieldValidator) AlphaNum() *FieldValidator {
	f.rules = append(f.rules, &AlphaNumRule{})
	return f
}

// Numeric valida que sea numérico
func (f *FieldValidator) Numeric() *FieldValidator {
	f.rules = append(f.rules, &NumericRule{})
	return f
}

// In valida que el valor esté en la lista de valores permitidos
func (f *FieldValidator) In(values ...interface{}) *FieldValidator {
	f.rules = append(f.rules, &InRule{Values: values})
	return f
}

// NotIn valida que el valor NO esté en la lista de valores
func (f *FieldValidator) NotIn(values ...interface{}) *FieldValidator {
	f.rules = append(f.rules, &NotInRule{Values: values})
	return f
}

// Sometimes marca el campo como opcional (solo validar si está presente)
func (f *FieldValidator) Sometimes() *FieldValidator {
	f.sometimes = true
	return f
}

// RequiredIf valida que sea requerido si otro campo tiene un valor específico
func (f *FieldValidator) RequiredIf(field string, value interface{}) *FieldValidator {
	f.rules = append(f.rules, &RequiredIfRule{Field: field, Value: value})
	return f
}

// RequiredUnless valida que sea requerido a menos que otro campo tenga un valor específico
func (f *FieldValidator) RequiredUnless(field string, value interface{}) *FieldValidator {
	f.rules = append(f.rules, &RequiredUnlessRule{Field: field, Value: value})
	return f
}

// Custom añade una regla personalizada
func (f *FieldValidator) Custom(rule Rule) *FieldValidator {
	f.rules = append(f.rules, rule)
	return f
}

// Array valida que el campo sea un array/objeto
func (f *FieldValidator) Array() *FieldValidator {
	f.rules = append(f.rules, &ArrayRule{})
	return f
}

// String valida que el campo sea una cadena
func (f *FieldValidator) String() *FieldValidator {
	f.rules = append(f.rules, &StringRule{})
	return f
}

// Field permite continuar con otro campo (para sintaxis fluida)
func (f *FieldValidator) Field(name string) *FieldValidator {
	return f.validator.Field(name)
}

// Función helper para validar requests en controladores
func Validate(context *gin.Context, request Validatable) error {
	// Obtener las reglas del request
	validator := request.Rules()

	// Obtener mensajes personalizados
	if messages := request.Messages(); messages != nil {
		validator.Messages(messages)
	}

	// Bind de los datos del request
	if err := context.ShouldBind(request); err != nil {
		context.JSON(422, ValidationResponse{
			Errors: []ValidationErrorResponse{{
				Status: "422",
				Title:  "Validation Error",
				Detail: "Invalid request format",
				Source: ValidationErrorSource{Pointer: "/data"},
				Meta: ValidationErrorMeta{
					Field:   "request",
					Rule:    "format",
					Message: err.Error(),
				},
			}},
		})
		return err
	}

	// Validar usando las reglas
	result := validator.Validate(request)

	if !result.Valid {
		errors := make([]ValidationErrorResponse, len(result.Errors))
		for i, err := range result.Errors {
			errors[i] = ValidationErrorResponse{
				Status: "422",
				Title:  "Validation Error",
				Detail: "The given data was invalid",
				Source: ValidationErrorSource{
					Pointer: fmt.Sprintf("/data/attributes/%s", err.Field),
				},
				Meta: ValidationErrorMeta{
					Field:   err.Field,
					Rule:    err.Rule,
					Message: err.Message,
				},
			}

			break // Detener al primer error encontrado
		}

		context.JSON(422, ValidationResponse{Errors: errors})
		return fmt.Errorf("validation failed")
	}

	return nil
}
