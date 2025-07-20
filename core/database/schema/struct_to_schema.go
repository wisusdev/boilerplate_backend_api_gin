package schema

import (
	"fmt"
	"reflect"
	"regexp"
	"strings"
)

// StructToSchemaGenerator convierte structs a código de Schema Builder
type StructToSchemaGenerator struct{}

// NewStructToSchemaGenerator crea una nueva instancia del generador
func NewStructToSchemaGenerator() *StructToSchemaGenerator {
	return &StructToSchemaGenerator{}
}

// GenerateSchemaCode genera código de Schema Builder a partir de una struct
func (schemaGenerator *StructToSchemaGenerator) GenerateSchemaCode(structType interface{}, tableName string) string {
	typeOf := reflect.TypeOf(structType)
	if typeOf.Kind() == reflect.Ptr {
		typeOf = typeOf.Elem()
	}

	var lines []string

	// Agregar la declaración inicial
	lines = append(lines, fmt.Sprintf(`schema := NewSchema()
sql := schema.Create("%s", func(table *Blueprint) {`, tableName))

	// Procesar cada campo de la struct
	for i := 0; i < typeOf.NumField(); i++ {
		field := typeOf.Field(i)
		line := schemaGenerator.generateFieldCode(field)
		if line != "" {
			lines = append(lines, "\t"+line)
		}
	}

	lines = append(lines, "})")
	lines = append(lines, "fmt.Println(sql)")

	return strings.Join(lines, "\n")
}

// generateFieldCode genera el código para un campo específico
func (schemaGenerator *StructToSchemaGenerator) generateFieldCode(field reflect.StructField) string {
	fieldName := schemaGenerator.getFieldName(field)
	if fieldName == "" || fieldName == "-" {
		return ""
	}

	// Casos especiales para campos comunes
	if special := schemaGenerator.handleSpecialFields(field, fieldName); special != "" {
		return special
	}

	baseLine := schemaGenerator.getBaseLine(field, fieldName)
	modifiers := schemaGenerator.generateModifiers(field)
	if len(modifiers) > 0 {
		return baseLine + modifiers
	}
	return baseLine
}

// handleSpecialFields maneja los campos especiales como id, created_at, etc.
func (schemaGenerator *StructToSchemaGenerator) handleSpecialFields(field reflect.StructField, fieldName string) string {
	dbTag := field.Tag.Get("db")
	goType := field.Type
	if goType.Kind() == reflect.Ptr {
		goType = goType.Elem()
	}

	switch fieldName {
	case "id":
		if schemaGenerator.isUUIDField(dbTag, goType, field.Name) {
			return `table.UuidPrimary("id")`
		}
		return `table.Increments("id")`
	case "created_at":
		return `table.Timestamps()`
	case "updated_at":
		return ""
	case "deleted_at":
		return `table.SoftDeletes()`
	}
	return ""
}

// isUUIDField determina si el campo es un UUID
func (schemaGenerator *StructToSchemaGenerator) isUUIDField(dbTag string, goType reflect.Type, name string) bool {
	return strings.Contains(strings.ToUpper(dbTag), "UUID") ||
		strings.Contains(strings.ToUpper(dbTag), "CHAR(36)") ||
		(goType.Kind() == reflect.String && strings.Contains(strings.ToLower(name), "uuid"))
}

// getBaseLine obtiene la línea base del tipo de campo
func (schemaGenerator *StructToSchemaGenerator) getBaseLine(field reflect.StructField, fieldName string) string {
	dbTag := field.Tag.Get("db")
	goType := field.Type
	if goType.Kind() == reflect.Ptr {
		goType = goType.Elem()
	}

	switch {
	case goType.Kind() == reflect.Int || goType.Kind() == reflect.Int32:
		return fmt.Sprintf(`table.Integer("%s")`, fieldName)
	case goType.Kind() == reflect.Int64:
		return fmt.Sprintf(`table.BigInteger("%s")`, fieldName)
	case goType.Kind() == reflect.String:
		return schemaGenerator.stringFieldLine(fieldName, dbTag)
	case goType.Kind() == reflect.Bool:
		return fmt.Sprintf(`table.Boolean("%s")`, fieldName)
	case goType.Kind() == reflect.Float32 || goType.Kind() == reflect.Float64:
		return fmt.Sprintf(`table.Decimal("%s", 10, 2)`, fieldName)
	case goType.String() == "time.Time":
		return fmt.Sprintf(`table.DateTime("%s")`, fieldName)
	default:
		return fmt.Sprintf(`table.Text("%s")`, fieldName)
	}
}

// stringFieldLine maneja los campos string con heurística de tags
func (schemaGenerator *StructToSchemaGenerator) stringFieldLine(fieldName, dbTag string) string {
	switch {
	case strings.Contains(strings.ToUpper(dbTag), "UUID"),
		strings.Contains(strings.ToUpper(dbTag), "CHAR(36)"),
		strings.Contains(strings.ToLower(fieldName), "uuid"):
		return fmt.Sprintf(`table.Uuid("%s")`, fieldName)
	case strings.Contains(dbTag, "VARCHAR"):
		re := regexp.MustCompile(`VARCHAR\((\d+)\)`)
		if matches := re.FindStringSubmatch(dbTag); len(matches) > 1 {
			return fmt.Sprintf(`table.String("%s", %s)`, fieldName, matches[1])
		}
		return fmt.Sprintf(`table.String("%s")`, fieldName)
	case strings.Contains(strings.ToLower(dbTag), "text"):
		return fmt.Sprintf(`table.Text("%s")`, fieldName)
	case strings.Contains(fieldName, "email"):
		return fmt.Sprintf(`table.String("%s")`, fieldName)
	case strings.Contains(fieldName, "description") || strings.Contains(fieldName, "content"):
		return fmt.Sprintf(`table.Text("%s")`, fieldName)
	default:
		return fmt.Sprintf(`table.String("%s")`, fieldName)
	}
}

// generateModifiers genera los modificadores para un campo
func (schemaGenerator *StructToSchemaGenerator) generateModifiers(field reflect.StructField) string {
	var modifiers []string

	// Nullable
	if field.Tag.Get("nullable") == "true" || field.Type.Kind() == reflect.Ptr {
		modifiers = append(modifiers, ".Nullable()")
	}

	// Unique
	if field.Tag.Get("unique") == "true" {
		modifiers = append(modifiers, ".Unique()")
	}

	// Index
	if schemaGenerator.shouldIndex(field) {
		modifiers = append(modifiers, ".Index()")
	}

	// Default
	if defaultVal := field.Tag.Get("default"); defaultVal != "" {
		switch defaultVal {
		case "CURRENT_TIMESTAMP", "NULL":
			modifiers = append(modifiers, fmt.Sprintf(`.Default("%s")`, defaultVal))
		case "TRUE", "FALSE":
			modifiers = append(modifiers, fmt.Sprintf(".Default(%s)", strings.ToLower(defaultVal)))
		default:
			// Remover comillas si ya las tiene
			cleanDefault := strings.Trim(defaultVal, "'\"")
			modifiers = append(modifiers, fmt.Sprintf(`.Default("%s")`, cleanDefault))
		}
	}

	// Comment (si hay descripción en los tags)
	if comment := field.Tag.Get("comment"); comment != "" {
		modifiers = append(modifiers, fmt.Sprintf(`.Comment("%s")`, comment))
	}

	// Unsigned para números
	if field.Tag.Get("unsigned") == "true" || schemaGenerator.isUnsignedType(field) {
		modifiers = append(modifiers, ".Unsigned()")
	}

	return strings.Join(modifiers, "")
}

// shouldIndex determina si un campo debería tener índice
func (schemaGenerator *StructToSchemaGenerator) shouldIndex(field reflect.StructField) bool {
	fieldName := schemaGenerator.getFieldName(field)

	// Campos que típicamente necesitan índices
	indexFields := []string{
		"user_id", "business_id", "account_id", "created_by", "updated_by",
		"category_id", "status", "type", "slug",
	}

	for _, indexField := range indexFields {
		if fieldName == indexField {
			return true
		}
	}

	// Si termina en _id, probablemente es FK
	if strings.HasSuffix(fieldName, "_id") {
		return true
	}

	return false
}

// isUnsignedType determina si un tipo numérico debería ser unsigned
func (schemaGenerator *StructToSchemaGenerator) isUnsignedType(field reflect.StructField) bool {
	fieldName := schemaGenerator.getFieldName(field)

	// IDs son típicamente unsigned
	if fieldName == "id" || strings.HasSuffix(fieldName, "_id") {
		return true
	}

	return false
}

// getFieldName obtiene el nombre del campo desde el tag json
func (schemaGenerator *StructToSchemaGenerator) getFieldName(field reflect.StructField) string {
	jsonTag := field.Tag.Get("json")
	if jsonTag == "" {
		return toSnakeCase(field.Name)
	}

	// Parsear el tag json (ej: "field_name,omitempty")
	parts := strings.Split(jsonTag, ",")
	return parts[0]
}

// toSnakeCase convierte PascalCase a snake_case
func toSnakeCase(s string) string {
	re := regexp.MustCompile("([a-z0-9])([A-Z])")
	snake := re.ReplaceAllString(s, "${1}_${2}")
	return strings.ToLower(snake)
}

// GenerateFromUserStruct ejemplo específico para el UserStruct
func (schemaGenerator *StructToSchemaGenerator) GenerateFromUserStruct() string {
	return `schema := NewSchema()
sql := schema.Create("users", func(table *Blueprint) {
	table.Increments("id")
	table.String("first_name")
	table.String("last_name")
	table.String("username").Unique()
	table.String("avatar").Nullable()
	table.String("language", 10).Default("en").Nullable()
	table.String("email").Unique()
	table.String("password")
	table.Timestamps()
})
fmt.Println(sql)`
}

// GenerateAccountsExample genera el ejemplo de cuentas como en Laravel
func (schemaGenerator *StructToSchemaGenerator) GenerateAccountsExample() string {
	return `schema := NewSchema()
sql := schema.Create("accounts", func(table *Blueprint) {
	table.Increments("id")
	table.Integer("business_id").Index()
	table.String("name", 191)
	table.String("account_number", 191)
	table.Text("account_details").Nullable()
	table.Integer("account_type_id").Nullable().Index()
	table.Text("note").Nullable()
	table.Integer("created_by").Index()
	table.Boolean("is_closed").Default(false)
	table.SoftDeletes()
	table.Timestamps()
})
fmt.Println(sql)`
}
