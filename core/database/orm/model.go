package orm

import (
	"boilerplate_backend_api_gin/app/data/models"
	"boilerplate_backend_api_gin/core/database/database_connections"
	"boilerplate_backend_api_gin/core/helpers"
	"database/sql"
	"fmt"
	"reflect"
	"strings"
	"time"
)

// ModelRepository es la estructura para trabajar con modelos específicos
type ModelRepository struct {
	Repository
	Model     Model
	TableName string
}

// NewModelRepository crea un nuevo repositorio para un modelo específico
func NewModelRepository(model *models.UserStruct) *ModelRepository {
	tableName := model.TableName()
	return &ModelRepository{
		Repository: Repository{
			DB:        database_connections.DatabaseConnectSQL(),
			TableName: tableName,
		},
		Model:     model,
		TableName: tableName,
	}
}

// GetAll recupera todos los registros del modelo
func (repository *ModelRepository) GetAll(result interface{}) error {
	query := fmt.Sprintf("SELECT * FROM %s", repository.TableName)
	rows, err := repository.DB.Query(query)
	if err != nil {
		helpers.Logs("ERROR", fmt.Sprintf("Error al ejecutar la consulta: %v", err))
		return err
	}
	defer rows.Close()

	return scanAll(rows, result)
}

// FindById busca un registro por su ID
func (repository *ModelRepository) FindById(id interface{}, result interface{}) error {
	query := fmt.Sprintf("SELECT * FROM %s WHERE id = ?", repository.TableName)
	row := repository.DB.QueryRow(query, id)

	return RowScanner(row, result)
}

// FindBy busca registros por un campo específico
func (repository *ModelRepository) FindBy(field string, value interface{}, result interface{}) error {
	query := fmt.Sprintf("SELECT * FROM %s WHERE %s = ?", repository.TableName, field)
	rows, err := repository.DB.Query(query, value)
	if err != nil {
		return err
	}
	defer rows.Close()

	return scanAll(rows, result)
}

// Where busca registros que cumplan una condición
func (repository *ModelRepository) Where(condition string, args []interface{}, result interface{}) error {
	query := fmt.Sprintf("SELECT * FROM %s WHERE %s", repository.TableName, condition)
	rows, err := repository.DB.Query(query, args...)
	if err != nil {
		return err
	}
	defer rows.Close()

	return scanAll(rows, result)
}

// First devuelve el primer registro que cumple una condición
func (repository *ModelRepository) First(condition string, args []interface{}, result interface{}) error {
	query := fmt.Sprintf("SELECT * FROM %s WHERE %s LIMIT 1", repository.TableName, condition)
	row := repository.DB.QueryRow(query, args...)

	return RowScanner(row, result)
}

// Save guarda un modelo (crea o actualiza)
func (repository *ModelRepository) Save(model interface{}) error {
	v := reflect.ValueOf(model)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	// Verificar si es un objeto nuevo o existente
	idField := v.FieldByName("ID")
	if !idField.IsValid() {
		return fmt.Errorf("model must have an ID field")
	}

	var result sql.Result
	var err error

	// Si ID es cero, insertar, de lo contrario actualizar
	if idField.Int() == 0 {
		result, err = repository.Create(model)
		if err != nil {
			return err
		}

		// Actualizar el ID del modelo con el generado
		lastID, err := result.LastInsertId()
		if err != nil {
			return err
		}

		idField.SetInt(lastID)
	} else {
		_, err = repository.Update(idField.Interface(), model)
		if err != nil {
			return err
		}
	}

	// Actualizar timestamps
	updatedAtField := v.FieldByName("UpdatedAt")
	if updatedAtField.IsValid() && updatedAtField.Type() == reflect.TypeOf(time.Time{}) {
		updatedAtField.Set(reflect.ValueOf(time.Now()))
	}

	return nil
}

// scanAll escanea múltiples filas en una slice de un modelo
func scanAll(rows *sql.Rows, result interface{}) error {
	v := reflect.ValueOf(result)
	if v.Kind() != reflect.Ptr || v.IsNil() {
		return fmt.Errorf("result must be a non-nil pointer to a slice")
	}

	sliceVal := v.Elem()
	if sliceVal.Kind() != reflect.Slice {
		return fmt.Errorf("result must be a pointer to a slice")
	}

	sliceType := sliceVal.Type()
	elementType := sliceType.Elem()

	var items []reflect.Value

	for rows.Next() {
		// Crear una nueva instancia del tipo de elemento
		newItem := reflect.New(elementType).Elem()
		newItemPtr := reflect.New(elementType)

		// Escanear la fila en el elemento
		err := RowScanner(rows, newItemPtr.Interface())
		if err != nil {
			return err
		}

		newItem.Set(newItemPtr.Elem())
		items = append(items, newItem)
	}

	// Verificar errores después de la iteración
	if err := rows.Err(); err != nil {
		return err
	}

	// Crear una nueva slice con los elementos escaneados
	newSlice := reflect.MakeSlice(sliceType, len(items), len(items))
	for i, item := range items {
		newSlice.Index(i).Set(item)
	}

	// Asignar la nueva slice al resultado
	sliceVal.Set(newSlice)

	return nil
}

// GenerateCreateTableSQL genera el SQL para crear una tabla a partir de un modelo
func GenerateCreateTableSQL(model interface{}) string {
	v := reflect.ValueOf(model)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	t := v.Type()

	// Obtener el nombre de la tabla
	var tableName string
	if m, ok := model.(Model); ok {
		tableName = m.TableName()
	} else {
		tableName = strings.ToLower(t.Name())
	}

	var columns []string
	var primaryKey string

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)

		// Obtener el nombre de columna desde la etiqueta json
		colName := field.Tag.Get("json")
		if colName == "" {
			colName = strings.ToLower(field.Name)
		} else {
			// Si tiene opciones como omitempty, las quitamos
			colName = strings.Split(colName, ",")[0]
		}

		// Obtener el tipo de columna
		dbType := field.Tag.Get("db")
		if dbType == "" {
			// Asignar un tipo por defecto según el tipo Go
			switch field.Type.Kind() {
			case reflect.String:
				dbType = "VARCHAR(255)"
			case reflect.Int, reflect.Int64:
				dbType = "INT"
			case reflect.Float64:
				dbType = "DOUBLE"
			case reflect.Bool:
				dbType = "BOOLEAN"
			case reflect.Struct:
				if field.Type == reflect.TypeOf(time.Time{}) {
					dbType = "DATETIME"
				}
			case reflect.Ptr:
				if field.Type == reflect.TypeOf((*time.Time)(nil)) {
					dbType = "DATETIME"
				}
			}
		}

		// Comprobar si es clave primaria
		if strings.Contains(dbType, "PRIMARY KEY") {
			primaryKey = colName
		}

		// Comprobar si es nullable
		nullable := field.Tag.Get("nullable")
		if nullable == "false" {
			dbType += " NOT NULL"
		}

		// Comprobar si es único
		unique := field.Tag.Get("unique")
		if unique == "true" {
			dbType += " UNIQUE"
		}

		// Comprobar si tiene un valor por defecto
		defaultVal := field.Tag.Get("default")
		if defaultVal != "" {
			dbType += " DEFAULT " + defaultVal
		}

		columns = append(columns, colName+" "+dbType)
	}

	// Si no se encontró una clave primaria, agregar una automáticamente
	if primaryKey == "" {
		columns = append([]string{"id INT PRIMARY KEY AUTO_INCREMENT"}, columns...)
	}

	// Generar la sentencia SQL
	return fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s (\n  %s\n);",
		tableName,
		strings.Join(columns, ",\n  "),
	)
}
