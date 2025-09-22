package orm

import (
	"boilerplate_backend_api_gin/core/database/database_connections"
	"boilerplate_backend_api_gin/core/helpers"
	"database/sql"
	"errors"
	"fmt"
	"reflect"
	"strings"
)

// Model es la interfaz que define un modelo de datos para el ORM
type Model interface {
	TableName() string
}

// QueryBuilder ayuda a construir consultas SQL
type QueryBuilder struct {
	tableName  string
	selections []string
	conditions []string
	parameters []interface{}
	orders     []string
	limit      int
	offset     int
	joins      []string
	groupBy    []string
	distinct   bool
}

// Repository es la estructura base para trabajar con el ORM
type Repository struct {
	DB        database_connections.SQLAdapter
	TableName string
}

// NewRepository crea un nuevo repositorio con una conexión a la base de datos
func NewRepository(tableName string) *Repository {
	return &Repository{
		DB:        database_connections.DatabaseConnectSQL(),
		TableName: tableName,
	}
}

// NewRepositoryWithDB crea un nuevo repositorio con una conexión a la base de datos existente
func NewRepositoryWithDB(tableName string, db database_connections.SQLAdapter) *Repository {
	return &Repository{
		DB:        db,
		TableName: tableName,
	}
}

// Query inicia un nuevo constructor de consultas
func (r *Repository) Query() *QueryBuilder {
	return &QueryBuilder{
		tableName:  r.TableName,
		selections: []string{"*"},
		limit:      0,
		offset:     0,
	}
}

// Select establece los campos a seleccionar
func (qb *QueryBuilder) Select(columns ...string) *QueryBuilder {
	qb.selections = columns
	return qb
}

// Where agrega una condición a la consulta
func (qb *QueryBuilder) Where(condition string, args ...interface{}) *QueryBuilder {
	qb.conditions = append(qb.conditions, condition)
	qb.parameters = append(qb.parameters, args...)
	return qb
}

// OrderBy agrega un ordenamiento a la consulta
func (qb *QueryBuilder) OrderBy(order string) *QueryBuilder {
	qb.orders = append(qb.orders, order)
	return qb
}

// Limit establece el límite de resultados
func (qb *QueryBuilder) Limit(limit int) *QueryBuilder {
	qb.limit = limit
	return qb
}

// Offset establece el desplazamiento de resultados
func (qb *QueryBuilder) Offset(offset int) *QueryBuilder {
	qb.offset = offset
	return qb
}

// Join agrega una unión a la consulta
func (qb *QueryBuilder) Join(join string) *QueryBuilder {
	qb.joins = append(qb.joins, join)
	return qb
}

// GroupBy agrega un agrupamiento a la consulta
func (qb *QueryBuilder) GroupBy(column string) *QueryBuilder {
	qb.groupBy = append(qb.groupBy, column)
	return qb
}

// Distinct hace que la consulta devuelva resultados únicos
func (qb *QueryBuilder) Distinct() *QueryBuilder {
	qb.distinct = true
	return qb
}

// BuildQuery construye la consulta SQL
func (qb *QueryBuilder) BuildQuery() (string, []interface{}) {
	var query strings.Builder
	query.WriteString("SELECT ")

	if qb.distinct {
		query.WriteString("DISTINCT ")
	}

	query.WriteString(strings.Join(qb.selections, ", "))
	query.WriteString(" FROM ")
	query.WriteString(qb.tableName)

	for _, join := range qb.joins {
		query.WriteString(" ")
		query.WriteString(join)
	}

	if len(qb.conditions) > 0 {
		query.WriteString(" WHERE ")
		query.WriteString(strings.Join(qb.conditions, " AND "))
	}

	if len(qb.groupBy) > 0 {
		query.WriteString(" GROUP BY ")
		query.WriteString(strings.Join(qb.groupBy, ", "))
	}

	if len(qb.orders) > 0 {
		query.WriteString(" ORDER BY ")
		query.WriteString(strings.Join(qb.orders, ", "))
	}

	if qb.limit > 0 {
		query.WriteString(fmt.Sprintf(" LIMIT %d", qb.limit))
	}

	if qb.offset > 0 {
		query.WriteString(fmt.Sprintf(" OFFSET %d", qb.offset))
	}

	return query.String(), qb.parameters
}

// Execute ejecuta la consulta y devuelve las filas
func (qb *QueryBuilder) Execute() (*sql.Rows, error) {
	query, params := qb.BuildQuery()
	rows, err := qb.getDB().Query(query, params...)
	if err != nil {
		helpers.Logs("ERROR", "Error ejecutando consulta: "+err.Error())
		return nil, err
	}
	return rows, nil
}

// First ejecuta la consulta y devuelve la primera fila
func (qb *QueryBuilder) First() (*sql.Row, error) {
	qb.Limit(1)
	query, params := qb.BuildQuery()
	return qb.getDB().QueryRow(query, params...), nil
}

// getDB obtiene la conexión a la base de datos
func (qb *QueryBuilder) getDB() database_connections.SQLAdapter {
	// Aquí podríamos obtener la conexión desde algún contexto global
	// Por ahora usamos la conexión por defecto
	return database_connections.DatabaseConnectSQL()
}

// FindAll encuentra todos los registros que cumplen las condiciones
func (r *Repository) FindAll(conditions map[string]interface{}, scanFunc func(*sql.Rows) (interface{}, error)) ([]interface{}, error) {
	qb := r.Query()

	for field, value := range conditions {
		qb.Where(field+" = ?", value)
	}

	rows, err := qb.Execute()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []interface{}
	for rows.Next() {
		item, err := scanFunc(rows)
		if err != nil {
			return nil, err
		}
		result = append(result, item)
	}

	return result, nil
}

// FindOne encuentra un registro que cumple las condiciones
func (r *Repository) FindOne(conditions map[string]interface{}, scanFunc func(*sql.Row) (interface{}, error)) (interface{}, error) {
	qb := r.Query()

	for field, value := range conditions {
		qb.Where(field+" = ?", value)
	}

	row, err := qb.First()
	if err != nil {
		return nil, err
	}

	return scanFunc(row)
}

// Create inserta un nuevo registro
func (r *Repository) Create(data interface{}) (sql.Result, error) {
	if r.DB == nil {
		return nil, errors.New("database connection is nil")
	}

	fields, values, err := extractFieldsAndValues(data)
	if err != nil {
		return nil, err
	}

	placeholders := make([]string, len(fields))
	for i := range fields {
		placeholders[i] = "?"
	}

	query := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)",
		r.TableName,
		strings.Join(fields, ", "),
		strings.Join(placeholders, ", "),
	)

	return r.DB.Exec(query, values...)
}

// Update actualiza un registro por su ID
func (r *Repository) Update(id interface{}, data interface{}) (sql.Result, error) {
	if r.DB == nil {
		return nil, errors.New("database connection is nil")
	}

	fields, values, err := extractFieldsAndValues(data)
	if err != nil {
		return nil, err
	}

	setClauses := make([]string, len(fields))
	for i, field := range fields {
		setClauses[i] = field + " = ?"
	}

	query := fmt.Sprintf("UPDATE %s SET %s WHERE id = ?",
		r.TableName,
		strings.Join(setClauses, ", "),
	)

	// Agregamos el ID como último parámetro
	values = append(values, id)

	return r.DB.Exec(query, values...)
}

// Delete elimina un registro por su ID
func (r *Repository) Delete(id interface{}) (sql.Result, error) {
	if r.DB == nil {
		return nil, errors.New("database connection is nil")
	}

	query := fmt.Sprintf("DELETE FROM %s WHERE id = ?", r.TableName)

	return r.DB.Exec(query, id)
}

// extractFieldsAndValues extrae los campos y valores de una estructura
func extractFieldsAndValues(data interface{}) ([]string, []interface{}, error) {
	v := reflect.ValueOf(data)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	if v.Kind() != reflect.Struct {
		return nil, nil, errors.New("data must be a struct or a pointer to a struct")
	}

	var fields []string
	var values []interface{}

	t := v.Type()
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)

		// Ignorar campos no exportados o marcados como `-`
		if !v.Field(i).CanInterface() || field.Tag.Get("json") == "-" {
			continue
		}

		// Obtenemos el nombre del campo desde la etiqueta json, o usamos el nombre del campo
		fieldName := field.Tag.Get("json")
		if fieldName == "" {
			fieldName = strings.ToLower(field.Name)
		} else {
			// Si tiene opciones como omitempty, las quitamos
			fieldName = strings.Split(fieldName, ",")[0]
		}

		// Ignorar campos como ID o timestamps que generalmente son gestionados por la BD
		if fieldName == "id" || fieldName == "created_at" || fieldName == "updated_at" {
			continue
		}

		fields = append(fields, fieldName)
		values = append(values, v.Field(i).Interface())
	}

	return fields, values, nil
}
