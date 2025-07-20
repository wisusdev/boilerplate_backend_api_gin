package generators

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"strings"
	"time"
)

// MigrationGenerator genera migraciones automáticamente desde structs
type MigrationGenerator struct {
	OutputDir string
}

// NewMigrationGenerator crea una nueva instancia del generador
func NewMigrationGenerator(outputDir string) *MigrationGenerator {
	return &MigrationGenerator{
		OutputDir: outputDir,
	}
}

// GenerateFromStruct genera una migración a partir de una struct
func (mg *MigrationGenerator) GenerateFromStruct(structType interface{}, tableName string) error {
	// Obtener información de la struct usando reflection
	typeOf := reflect.TypeOf(structType)
	if typeOf.Kind() == reflect.Ptr {
		typeOf = typeOf.Elem()
	}

	// Generar nombre de archivo con timestamp
	timestamp := time.Now().Format("2006_01_02_150405")
	fileName := fmt.Sprintf("%s_create_%s_table.go", timestamp, tableName)
	filePath := filepath.Join(mg.OutputDir, fileName)

	// Generar contenido de la migración
	content := mg.generateMigrationContent(typeOf, tableName, timestamp)

	// Escribir archivo
	err := os.WriteFile(filePath, []byte(content), 0644)
	if err != nil {
		return fmt.Errorf("error writing migration file: %v", err)
	}

	fmt.Printf("Migration generated: %s\n", filePath)
	return nil
}

// generateMigrationContent genera el contenido de la migración
func (mg *MigrationGenerator) generateMigrationContent(structType reflect.Type, tableName, timestamp string) string {
	migrationName := fmt.Sprintf("Create%sTable", toPascalCase(tableName))

	// Generar campos de la tabla
	fields := mg.generateTableFields(structType)

	template := `package generate_migrations

import (
	"database/sql"
	"semita/app/core/database"
)

type %s struct {
	database.BaseMigration
}

func New%s() *%s {
	return &%s{
		BaseMigration: database.BaseMigration{
			Name:      "create_%s_table",
			Timestamp: "%s",
		},
	}
}

func (m *%s) Up(db database_connections.SQLAdapter) error {
	query := ` + "`" + `
		CREATE TABLE %s (
%s
		)
	` + "`" + `
	_, err := db.Exec(query)
	return err
}

func (m *%s) Down(db database_connections.SQLAdapter) error {
	_, err := db.Exec("DROP TABLE IF EXISTS %s")
	return err
}
`

	return fmt.Sprintf(template,
		migrationName, // Type declaration
		migrationName, // Function name
		migrationName, // Return type
		migrationName, // Constructor
		tableName,     // Name field
		timestamp,     // Timestamp field
		migrationName, // Up method receiver
		tableName,     // Table name in CREATE
		fields,        // Fields
		migrationName, // Down method receiver
		tableName,     // Table name in DROP
	)
}

// generateTableFields genera los campos de la tabla basándose en la struct
func (mg *MigrationGenerator) generateTableFields(structType reflect.Type) string {
	var fields []string

	for i := 0; i < structType.NumField(); i++ {
		field := structType.Field(i)

		// Obtener el nombre del campo desde el tag json, o usar el nombre del campo
		fieldName := mg.getFieldName(field)
		if fieldName == "" || fieldName == "-" {
			continue
		}

		// Generar definición SQL del campo
		sqlDef := mg.generateFieldSQL(field, fieldName)
		if sqlDef != "" {
			fields = append(fields, "\t\t\t"+sqlDef)
		}
	}

	return strings.Join(fields, ",\n")
}

// getFieldName obtiene el nombre del campo desde el tag json
func (mg *MigrationGenerator) getFieldName(field reflect.StructField) string {
	jsonTag := field.Tag.Get("json")
	if jsonTag == "" {
		return toSnakeCase(field.Name)
	}

	// Parsear el tag json (ej: "field_name,omitempty")
	parts := strings.Split(jsonTag, ",")
	return parts[0]
}

// generateFieldSQL genera la definición SQL para un campo
func (mg *MigrationGenerator) generateFieldSQL(field reflect.StructField, fieldName string) string {
	// Obtener tags personalizados para SQL
	dbTag := field.Tag.Get("db")
	sqlType := mg.getSQLType(field.Type, dbTag)

	var constraints []string

	// Primary key
	if fieldName == "id" {
		constraints = append(constraints, "PRIMARY KEY AUTO_INCREMENT")
	}

	// Not null (por defecto, a menos que sea pointer)
	if field.Type.Kind() != reflect.Ptr && !mg.isNullable(field) {
		constraints = append(constraints, "NOT NULL")
	}

	// Unique constraint
	if mg.hasTag(field, "unique") {
		constraints = append(constraints, "UNIQUE")
	}

	// Default values
	if defaultVal := field.Tag.Get("default"); defaultVal != "" {
		constraints = append(constraints, fmt.Sprintf("DEFAULT %s", defaultVal))
	}

	// Timestamps especiales
	if fieldName == "created_at" || fieldName == "updated_at" {
		constraints = append(constraints, "DEFAULT CURRENT_TIMESTAMP")
	}

	constraintStr := ""
	if len(constraints) > 0 {
		constraintStr = " " + strings.Join(constraints, " ")
	}

	return fmt.Sprintf("%s %s%s", fieldName, sqlType, constraintStr)
}

// getSQLType convierte tipos de Go a tipos SQL
func (mg *MigrationGenerator) getSQLType(goType reflect.Type, dbTag string) string {
	// Si hay un tag db específico, usarlo
	if dbTag != "" {
		return dbTag
	}

	// Manejar pointers
	if goType.Kind() == reflect.Ptr {
		goType = goType.Elem()
	}

	switch goType.Kind() {
	case reflect.Int, reflect.Int32:
		return "INT"
	case reflect.Int64:
		return "BIGINT"
	case reflect.String:
		return "VARCHAR(255)"
	case reflect.Bool:
		return "BOOLEAN"
	case reflect.Float32, reflect.Float64:
		return "DECIMAL(10,2)"
	default:
		// Para tipos complejos como time.Time
		switch goType.String() {
		case "time.Time":
			return "DATETIME"
		default:
			return "TEXT"
		}
	}
}

// isNullable determina si un campo puede ser NULL
func (mg *MigrationGenerator) isNullable(field reflect.StructField) bool {
	return field.Type.Kind() == reflect.Ptr ||
		field.Tag.Get("nullable") == "true" ||
		field.Tag.Get("omitempty") != ""
}

// hasTag verifica si un campo tiene un tag específico
func (mg *MigrationGenerator) hasTag(field reflect.StructField, tagName string) bool {
	return field.Tag.Get(tagName) == "true"
}

// Utility functions
func toPascalCase(s string) string {
	parts := strings.Split(s, "_")
	for i, part := range parts {
		parts[i] = strings.Title(part)
	}
	return strings.Join(parts, "")
}

func toSnakeCase(s string) string {
	// Convert PascalCase to snake_case
	re := regexp.MustCompile("([a-z0-9])([A-Z])")
	snake := re.ReplaceAllString(s, "${1}_${2}")
	return strings.ToLower(snake)
}
