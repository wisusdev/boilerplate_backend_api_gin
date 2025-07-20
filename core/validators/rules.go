package validators

import (
	"fmt"
	"reflect"
	"regexp"
	"semita/core/database/database_connections"
	"strconv"
	"strings"
)

// RequiredRule válida que el campo esté presente y no esté vacío
type RequiredRule struct{}

func (r *RequiredRule) Name() string {
	return "required"
}

func (r *RequiredRule) Validate(value interface{}, data map[string]interface{}) error {
	if value == nil {
		return fmt.Errorf("the field is required")
	}

	// Verificar valores vacíos según el tipo
	valueOf := reflect.ValueOf(value)
	switch valueOf.Kind() {
	case reflect.String:
		if strings.TrimSpace(valueOf.String()) == "" {
			return fmt.Errorf("the field is required")
		}
	case reflect.Slice, reflect.Array, reflect.Map:
		if valueOf.Len() == 0 {
			return fmt.Errorf("the field is required")
		}
	case reflect.Ptr:
		if valueOf.IsNil() {
			return fmt.Errorf("the field is required")
		}
	default:
		panic("unhandled default case")
	}

	return nil
}

// EmailRule valida que el campo sea un email válido
type EmailRule struct{}

func (r *EmailRule) Name() string {
	return "email"
}

func (r *EmailRule) Validate(value interface{}, data map[string]interface{}) error {
	if value == nil {
		return nil // Si es nil, no validar (usar Required si es necesario)
	}

	str, ok := value.(string)
	if !ok {
		return fmt.Errorf("the field must be a string")
	}

	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(str) {
		return fmt.Errorf("the field must be a valid email address")
	}

	return nil
}

// MinRule valida el valor mínimo
type MinRule struct {
	Min int
}

func (r *MinRule) Name() string {
	return "min"
}

func (r *MinRule) Validate(value interface{}, data map[string]interface{}) error {
	if value == nil {
		return nil
	}

	switch v := value.(type) {
	case string:
		if len(v) < r.Min {
			return fmt.Errorf("the field must be at least %d characters", r.Min)
		}
	case float32, float64:
		floatVal := reflect.ValueOf(v).Float()
		if int(floatVal) < r.Min {
			return fmt.Errorf("the field must be at least %d", r.Min)
		}
	case []interface{}:
		if len(v) < r.Min {
			return fmt.Errorf("the field must have at least %d items", r.Min)
		}
	default:
		return fmt.Errorf("the min rule cannot be applied to this field type")
	}

	return nil
}

// MaxRule valida el valor máximo
type MaxRule struct {
	Max int
}

func (r *MaxRule) Name() string {
	return "max"
}

func (r *MaxRule) Validate(value interface{}, data map[string]interface{}) error {
	if value == nil {
		return nil
	}

	switch v := value.(type) {
	case string:
		if len(v) > r.Max {
			return fmt.Errorf("the field may not be greater than %d characters", r.Max)
		}
	case int, int8, int16, int32, int64:
		intVal := reflect.ValueOf(v).Int()
		if int(intVal) > r.Max {
			return fmt.Errorf("the field may not be greater than %d", r.Max)
		}
	case float32, float64:
		floatVal := reflect.ValueOf(v).Float()
		if int(floatVal) > r.Max {
			return fmt.Errorf("the field may not be greater than %d", r.Max)
		}
	case []interface{}:
		if len(v) > r.Max {
			return fmt.Errorf("the field may not have more than %d items", r.Max)
		}
	default:
		return fmt.Errorf("the max rule cannot be applied to this field type")
	}

	return nil
}

// BetweenRule valida que el valor esté entre min y max
type BetweenRule struct {
	Min int
	Max int
}

func (r *BetweenRule) Name() string {
	return "between"
}

func (r *BetweenRule) Validate(value interface{}, data map[string]interface{}) error {
	if value == nil {
		return nil
	}

	switch v := value.(type) {
	case string:
		length := len(v)
		if length < r.Min || length > r.Max {
			return fmt.Errorf("the field must be between %d and %d characters", r.Min, r.Max)
		}
	case int, int8, int16, int32, int64:
		intVal := int(reflect.ValueOf(v).Int())
		if intVal < r.Min || intVal > r.Max {
			return fmt.Errorf("the field must be between %d and %d", r.Min, r.Max)
		}
	case float32, float64:
		floatVal := int(reflect.ValueOf(v).Float())
		if floatVal < r.Min || floatVal > r.Max {
			return fmt.Errorf("the field must be between %d and %d", r.Min, r.Max)
		}
	default:
		return fmt.Errorf("the between rule cannot be applied to this field type")
	}

	return nil
}

// UniqueRule valida que el valor sea único en la base de datos
type UniqueRule struct {
	Table  string
	Column string
	Except []interface{} // Valores a excluir (útil para updates)
	DB     database_connections.SQLAdapter
}

func (r *UniqueRule) Name() string {
	return "unique"
}

func (r *UniqueRule) Validate(value interface{}, data map[string]interface{}) error {
	if value == nil {
		return nil
	}

	if r.DB == nil {
		return fmt.Errorf("database connection not provided for unique validation")
	}

	// Construir la consulta base
	query := fmt.Sprintf("SELECT COUNT(*) FROM %s WHERE %s = ?", r.Table, r.Column)
	args := []interface{}{value}

	// Agregar condiciones de exclusión si existen
	if len(r.Except) > 0 {
		query += fmt.Sprintf(" AND %s NOT IN (?", r.Column)
		for i := 1; i < len(r.Except); i++ {
			query += ",?"
		}
		query += ")"
		args = append(args, r.Except...)
	}

	var count int64
	err := r.DB.QueryRow(query, args...).Scan(&count)
	if err != nil {
		return fmt.Errorf("database error: %v", err)
	}

	if count > 0 {
		return fmt.Errorf("the field has already been taken")
	}

	return nil
}

// ExistsRule valida que el valor exista en la base de datos
type ExistsRule struct {
	Table  string
	Column string
	DB     database_connections.SQLAdapter
}

func (r *ExistsRule) Name() string {
	return "exists"
}

func (r *ExistsRule) Validate(value interface{}, data map[string]interface{}) error {
	if value == nil {
		return nil
	}

	if r.DB == nil {
		return fmt.Errorf("database connection not provided for exists validation")
	}

	query := fmt.Sprintf("SELECT COUNT(*) FROM %s WHERE %s = ?", r.Table, r.Column)
	var count int64
	err := r.DB.QueryRow(query, value).Scan(&count)
	if err != nil {
		return fmt.Errorf("database error: %v", err)
	}

	if count == 0 {
		return fmt.Errorf("the selected field is invalid")
	}

	return nil
}

// ConfirmedRule valida que el campo coincida con su confirmación
type ConfirmedRule struct {
	Field string
}

func (r *ConfirmedRule) Name() string {
	return "confirmed"
}

func (r *ConfirmedRule) Validate(value interface{}, data map[string]interface{}) error {
	if value == nil {
		return fmt.Errorf("the field confirmation is required")
	}

	confirmationField := r.Field + "_confirmation"
	confirmationValue, exists := data[confirmationField]
	if !exists {
		return fmt.Errorf("the field confirmation is required")
	}

	if !reflect.DeepEqual(value, confirmationValue) {
		return fmt.Errorf("the field confirmation does not match")
	}

	return nil
}

// RegexRule valida contra una expresión regular
type RegexRule struct {
	Pattern string
}

func (r *RegexRule) Name() string {
	return "regex"
}

func (r *RegexRule) Validate(value interface{}, data map[string]interface{}) error {
	if value == nil {
		return nil
	}

	str, ok := value.(string)
	if !ok {
		return fmt.Errorf("the field must be a string for regex validation")
	}

	regex, err := regexp.Compile(r.Pattern)
	if err != nil {
		return fmt.Errorf("invalid regex pattern: %s", err.Error())
	}

	if !regex.MatchString(str) {
		return fmt.Errorf("the field format is invalid")
	}

	return nil
}

// AlphaRule valida que solo contenga letras
type AlphaRule struct{}

func (r *AlphaRule) Name() string {
	return "alpha"
}

func (r *AlphaRule) Validate(value interface{}, data map[string]interface{}) error {
	if value == nil {
		return nil
	}

	str, ok := value.(string)
	if !ok {
		return fmt.Errorf("the field must be a string")
	}

	alphaRegex := regexp.MustCompile(`^[a-zA-ZÀ-ÿĀ-žА-я\u4e00-\u9fff]+$`)
	if !alphaRegex.MatchString(str) {
		return fmt.Errorf("the field may only contain letters")
	}

	return nil
}

// AlphaNumRule valida que solo contenga letras y números
type AlphaNumRule struct{}

func (r *AlphaNumRule) Name() string {
	return "alpha_num"
}

func (r *AlphaNumRule) Validate(value interface{}, data map[string]interface{}) error {
	if value == nil {
		return nil
	}

	str, ok := value.(string)
	if !ok {
		return fmt.Errorf("the field must be a string")
	}

	alphaNumRegex := regexp.MustCompile(`^[a-zA-Z0-9À-ÿĀ-žА-я\u4e00-\u9fff]+$`)
	if !alphaNumRegex.MatchString(str) {
		return fmt.Errorf("the field may only contain letters and numbers")
	}

	return nil
}

// NumericRule valida que sea numérico
type NumericRule struct{}

func (r *NumericRule) Name() string {
	return "numeric"
}

func (r *NumericRule) Validate(value interface{}, data map[string]interface{}) error {
	if value == nil {
		return nil
	}

	switch v := value.(type) {
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, float32, float64:
		return nil
	case string:
		if _, err := strconv.ParseFloat(v, 64); err != nil {
			return fmt.Errorf("the field must be numeric")
		}
		return nil
	default:
		return fmt.Errorf("the field must be numeric")
	}
}

// InRule valida que el valor esté en la lista de valores permitidos
type InRule struct {
	Values []interface{}
}

func (r *InRule) Name() string {
	return "in"
}

func (r *InRule) Validate(value interface{}, data map[string]interface{}) error {
	if value == nil {
		return nil
	}

	for _, allowedValue := range r.Values {
		if reflect.DeepEqual(value, allowedValue) {
			return nil
		}
	}

	return fmt.Errorf("the selected field is invalid")
}

// NotInRule valida que el valor NO esté en la lista de valores
type NotInRule struct {
	Values []interface{}
}

func (r *NotInRule) Name() string {
	return "not_in"
}

func (r *NotInRule) Validate(value interface{}, data map[string]interface{}) error {
	if value == nil {
		return nil
	}

	for _, forbiddenValue := range r.Values {
		if reflect.DeepEqual(value, forbiddenValue) {
			return fmt.Errorf("the selected field is invalid")
		}
	}

	return nil
}

// RequiredIfRule valida que sea requerido si otro campo tiene un valor específico
type RequiredIfRule struct {
	Field string
	Value interface{}
}

func (r *RequiredIfRule) Name() string {
	return "required_if"
}

func (r *RequiredIfRule) Validate(value interface{}, data map[string]interface{}) error {
	otherValue, exists := data[r.Field]
	if !exists {
		return nil // Si el otro campo no existe, no se requiere este campo
	}

	if reflect.DeepEqual(otherValue, r.Value) {
		// El otro campo tiene el valor esperado, este campo es requerido
		if value == nil {
			return fmt.Errorf("the field is required when %s is %v", r.Field, r.Value)
		}
		// Verificar si es una cadena vacía
		if str, ok := value.(string); ok && strings.TrimSpace(str) == "" {
			return fmt.Errorf("the field is required when %s is %v", r.Field, r.Value)
		}
	}

	return nil
}

// RequiredUnlessRule valida que sea requerido a menos que otro campo tenga un valor específico
type RequiredUnlessRule struct {
	Field string
	Value interface{}
}

func (r *RequiredUnlessRule) Name() string {
	return "required_unless"
}

func (r *RequiredUnlessRule) Validate(value interface{}, data map[string]interface{}) error {
	otherValue, exists := data[r.Field]
	if !exists {
		// Si el otro campo no existe, este campo es requerido
		if value == nil {
			return fmt.Errorf("the field is required unless %s is %v", r.Field, r.Value)
		}
		if str, ok := value.(string); ok && strings.TrimSpace(str) == "" {
			return fmt.Errorf("the field is required unless %s is %v", r.Field, r.Value)
		}
	} else if !reflect.DeepEqual(otherValue, r.Value) {
		// El otro campo existe pero no tiene el valor esperado, este campo es requerido
		if value == nil {
			return fmt.Errorf("the field is required unless %s is %v", r.Field, r.Value)
		}
		if str, ok := value.(string); ok && strings.TrimSpace(str) == "" {
			return fmt.Errorf("the field is required unless %s is %v", r.Field, r.Value)
		}
	}

	return nil
}

// ArrayRule valida que el campo sea un array u objeto
type ArrayRule struct{}

func (r *ArrayRule) Name() string {
	return "array"
}

func (r *ArrayRule) Validate(value interface{}, data map[string]interface{}) error {
	if value == nil {
		return fmt.Errorf("the field must be an array")
	}

	v := reflect.ValueOf(value)
	switch v.Kind() {
	case reflect.Slice, reflect.Array, reflect.Map:
		return nil
	default:
		return fmt.Errorf("the field must be an array")
	}
}

// StringRule válida que el campo sea una cadena
type StringRule struct{}

func (r *StringRule) Name() string {
	return "string"
}

func (r *StringRule) Validate(value interface{}, data map[string]interface{}) error {
	if value == nil {
		return fmt.Errorf("the field must be a string")
	}

	v := reflect.ValueOf(value)
	if v.Kind() != reflect.String {
		return fmt.Errorf("the field must be a string")
	}

	return nil
}
