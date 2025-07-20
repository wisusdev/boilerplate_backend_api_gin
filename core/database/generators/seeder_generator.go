package generators

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"
)

// SeederGenerator genera seeders automáticamente desde structs
type SeederGenerator struct {
	OutputDir string
}

// NewSeederGenerator crea una nueva instancia del generador
func NewSeederGenerator(outputDir string) *SeederGenerator {
	return &SeederGenerator{
		OutputDir: outputDir,
	}
}

// GenerateFromStruct genera un seeder a partir de una struct
func (sg *SeederGenerator) GenerateFromStruct(structType interface{}, tableName string, dependencies []string) error {
	// Obtener información de la struct usando reflection
	typeOf := reflect.TypeOf(structType)
	if typeOf.Kind() == reflect.Ptr {
		typeOf = typeOf.Elem()
	}

	// Generar nombre de archivo
	fileName := fmt.Sprintf("%s_seeder.go", tableName)
	filePath := filepath.Join(sg.OutputDir, fileName)

	// Generar contenido del seeder
	content := sg.generateSeederContent(typeOf, tableName, dependencies)

	// Escribir archivo
	err := os.WriteFile(filePath, []byte(content), 0644)
	if err != nil {
		return fmt.Errorf("error writing seeder file: %v", err)
	}

	fmt.Printf("Seeder generated: %s\n", filePath)
	return nil
}

// generateSeederContent genera el contenido del seeder
func (sg *SeederGenerator) generateSeederContent(structType reflect.Type, tableName string, dependencies []string) string {
	seederName := fmt.Sprintf("%sSeeder", toPascalCase(tableName))

	// Generar datos de ejemplo
	sampleData := sg.generateSampleData(structType, tableName)

	// Generar lista de dependencias
	dependenciesStr := sg.generateDependenciesArray(dependencies)

	// Generar lista de tablas (con relaciones)
	tables := sg.generateTablesArray(tableName)

	template := `package seeders

import (
	"log"
	"semita/app/core/database"
	"semita/config"
)

// %s seeder para %s
type %s struct {
	database.BaseSeeder
}

// New%s crea una nueva instancia del seeder
func New%s() *%s {
	return &%s{
		BaseSeeder: database.BaseSeeder{
			DB:   database_connections.DatabaseConnectSQL(),
			Name: "%s_seeder",
		},
	}
}

// GetName retorna el nombre del seeder
func (s *%s) GetName() string {
	return s.Name
}

// GetDependencies retorna las dependencias del seeder
func (s *%s) GetDependencies() []string {
	return %s
}

// GetTables retorna las tablas que maneja este seeder
func (s *%s) GetTables() []string {
	return %s
}

// Seed ejecuta el seeding de %s
func (s *%s) Seed() error {
	log.Println("Seeding %s...")

%s

	log.Println("%s seeding completed successfully!")
	return nil
}
`

	return fmt.Sprintf(template,
		seederName,              // Type comment
		tableName,               // Comment description
		seederName,              // Type declaration
		seederName,              // Constructor name
		seederName,              // Constructor name 2
		seederName,              // Constructor return type
		seederName,              // Constructor initialization
		tableName,               // Name field
		seederName,              // GetName receiver
		seederName,              // GetDependencies receiver
		dependenciesStr,         // Dependencies array
		seederName,              // GetTables receiver
		tables,                  // Tables array
		tableName,               // Seed comment
		seederName,              // Seed receiver
		tableName,               // Seed log
		sampleData,              // Sample data and insertion
		toPascalCase(tableName), // Success log
	)
}

// generateSampleData genera datos de ejemplo para el seeder
func (sg *SeederGenerator) generateSampleData(structType reflect.Type, tableName string) string {
	fields, placeholders := sg.getInsertFieldsAndPlaceholders(structType)
	dataItems := sg.getSampleDataItems(structType, fields)
	if len(fields) == 0 {
		return "\t// No fields to seed"
	}

	insertCode := sg.buildInsertCode(structType, tableName, fields, placeholders, dataItems)
	return insertCode
}

// Extrae los campos y placeholders válidos para la inserción
func (sg *SeederGenerator) getInsertFieldsAndPlaceholders(structType reflect.Type) ([]string, []string) {
	var fields []string
	var placeholders []string
	for i := 0; i < structType.NumField(); i++ {
		field := structType.Field(i)
		fieldName := sg.getFieldName(field)
		if sg.isSkippableField(fieldName) {
			continue
		}
		fields = append(fields, fieldName)
		placeholders = append(placeholders, "?")
	}
	return fields, placeholders
}

// Genera los datos de ejemplo para los items
func (sg *SeederGenerator) getSampleDataItems(structType reflect.Type, fields []string) []string {
	var dataItems []string
	for i := 1; i <= 3; i++ {
		values := sg.getSampleValuesForItem(structType, fields, i)
		if len(values) > 0 {
			dataItem := fmt.Sprintf("\t\t{%s}", strings.Join(values, ", "))
			dataItems = append(dataItems, dataItem)
		}
	}
	return dataItems
}

// Obtiene los valores de ejemplo para un item
func (sg *SeederGenerator) getSampleValuesForItem(structType reflect.Type, fields []string, index int) []string {
	var values []string
	for j := 0; j < structType.NumField(); j++ {
		field := structType.Field(j)
		fieldName := sg.getFieldName(field)
		if sg.isSkippableField(fieldName) {
			continue
		}
		sampleValue := sg.generateSampleValue(field, fieldName, index)
		values = append(values, sampleValue)
	}
	return values
}

// Determina si un campo debe ser omitido
func (sg *SeederGenerator) isSkippableField(fieldName string) bool {
	return fieldName == "id" || fieldName == "created_at" || fieldName == "updated_at" || fieldName == "" || fieldName == "-"
}

// Construye el código de inserción
func (sg *SeederGenerator) buildInsertCode(structType reflect.Type, tableName string, fields, placeholders, dataItems []string) string {
	return fmt.Sprintf(`	data := []struct {
%s
	}{
%s,
	}

	for _, item := range data {
		insertQuery := `+"`"+`
			INSERT INTO %s (%s) 
			VALUES (%s)
		`+"`"+`

		_, err := s.DB.Exec(insertQuery%s)
		if err != nil {
			log.Printf("Error creating %s: %%v", err)
			continue
		}

		log.Printf("Created %s item")
	}`,
		sg.generateStructFields(structType),
		strings.Join(dataItems, ",\n"),
		tableName,
		strings.Join(fields, ", "),
		strings.Join(placeholders, ", "),
		sg.generateExecParams(structType),
		tableName,
		tableName,
	)
}

// generateStructFields genera los campos de la struct para el seeder
func (sg *SeederGenerator) generateStructFields(structType reflect.Type) string {
	var fields []string

	for i := 0; i < structType.NumField(); i++ {
		field := structType.Field(i)
		fieldName := sg.getFieldName(field)

		// Saltar campos especiales
		if fieldName == "id" || fieldName == "created_at" || fieldName == "updated_at" || fieldName == "" || fieldName == "-" {
			continue
		}

		goType := sg.getGoTypeString(field.Type)
		structFieldName := toPascalCase(fieldName)
		fields = append(fields, fmt.Sprintf("\t\t%s %s", structFieldName, goType))
	}

	return strings.Join(fields, "\n")
}

// generateExecParams genera los parámetros para la función Exec
func (sg *SeederGenerator) generateExecParams(structType reflect.Type) string {
	var params []string

	for i := 0; i < structType.NumField(); i++ {
		field := structType.Field(i)
		fieldName := sg.getFieldName(field)

		// Saltar campos especiales
		if fieldName == "id" || fieldName == "created_at" || fieldName == "updated_at" || fieldName == "" || fieldName == "-" {
			continue
		}

		structFieldName := toPascalCase(fieldName)
		params = append(params, fmt.Sprintf("item.%s", structFieldName))
	}

	if len(params) > 0 {
		return ", " + strings.Join(params, ", ")
	}

	return ""
}

// generateSampleValue genera un valor de ejemplo basado en el tipo y nombre del campo
func (sg *SeederGenerator) generateSampleValue(field reflect.StructField, fieldName string, index int) string {
	goType := field.Type
	if goType.Kind() == reflect.Ptr {
		goType = goType.Elem()
	}

	switch goType.Kind() {
	case reflect.String:
		// Generar valores específicos basados en el nombre del campo
		switch fieldName {
		case "email":
			return fmt.Sprintf(`"user%d@example.com"`, index)
		case "username":
			return fmt.Sprintf(`"user%d"`, index)
		case "first_name":
			names := []string{"John", "Jane", "Mike"}
			return fmt.Sprintf(`"%s"`, names[(index-1)%len(names)])
		case "last_name":
			surnames := []string{"Doe", "Smith", "Johnson"}
			return fmt.Sprintf(`"%s"`, surnames[(index-1)%len(surnames)])
		case "name":
			names := []string{"Item One", "Item Two", "Item Three"}
			return fmt.Sprintf(`"%s"`, names[(index-1)%len(names)])
		case "description":
			return fmt.Sprintf(`"Description for item %d"`, index)
		case "slug":
			return fmt.Sprintf(`"item-%d"`, index)
		default:
			return fmt.Sprintf(`"Sample %s %d"`, strings.ReplaceAll(fieldName, "_", " "), index)
		}
	case reflect.Int, reflect.Int32, reflect.Int64:
		return fmt.Sprintf("%d", index)
	case reflect.Bool:
		return fmt.Sprintf("%t", index%2 == 0)
	case reflect.Float32, reflect.Float64:
		return fmt.Sprintf("%.2f", float64(index)*10.5)
	default:
		return fmt.Sprintf(`"value_%d"`, index)
	}
}

// getGoTypeString convierte un reflect.Type a string
func (sg *SeederGenerator) getGoTypeString(goType reflect.Type) string {
	if goType.Kind() == reflect.Ptr {
		return "*" + sg.getGoTypeString(goType.Elem())
	}

	switch goType.Kind() {
	case reflect.String:
		return "string"
	case reflect.Int:
		return "int"
	case reflect.Int32:
		return "int32"
	case reflect.Int64:
		return "int64"
	case reflect.Bool:
		return "bool"
	case reflect.Float32:
		return "float32"
	case reflect.Float64:
		return "float64"
	default:
		return goType.String()
	}
}

// getFieldName obtiene el nombre del campo desde el tag json
func (sg *SeederGenerator) getFieldName(field reflect.StructField) string {
	jsonTag := field.Tag.Get("json")
	if jsonTag == "" {
		return toSnakeCase(field.Name)
	}

	// Parsear el tag json (ej: "field_name,omitempty")
	parts := strings.Split(jsonTag, ",")
	return parts[0]
}

// generateDependenciesArray genera el array de dependencias
func (sg *SeederGenerator) generateDependenciesArray(dependencies []string) string {
	if len(dependencies) == 0 {
		return "[]string{} // No dependencies"
	}

	var deps []string
	for _, dep := range dependencies {
		deps = append(deps, fmt.Sprintf(`"%s"`, dep))
	}

	return fmt.Sprintf("[]string{%s}", strings.Join(deps, ", "))
}

// generateTablesArray genera el array de tablas
func (sg *SeederGenerator) generateTablesArray(tableName string) string {
	// Por ahora solo incluye la tabla principal
	// Se puede extender para incluir tablas relacionales
	return fmt.Sprintf(`[]string{"%s"}`, tableName)
}
