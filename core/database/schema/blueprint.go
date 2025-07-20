package schema

import (
	"fmt"
	"strings"
)

// Blueprint representa la definición de una tabla
type Blueprint struct {
	tableName  string
	columns    []*Column
	indexes    []Index
	foreign    []ForeignKey
	primaryKey []string
}

// Column representa una columna de la tabla
type Column struct {
	Name          string
	Type          string
	Length        int
	Precision     int
	Scale         int
	IsUnsigned    bool
	IsNullable    bool
	DefaultValue  interface{}
	AutoIncrement bool
	IsPrimary     bool
	IsUnique      bool
	HasIndex      bool
	CommentText   string
	EnumValues    []string // Para ENUM y SET
	OnUpdate      string   // Para ON UPDATE CURRENT_TIMESTAMP
}

// Index representa un índice
type Index struct {
	Name    string
	Columns []string
	Type    string // 'index', 'unique', 'primary'
}

// ForeignKey representa una clave foránea
type ForeignKey struct {
	Column           string
	ReferencedTable  string
	ReferencedColumn string
	OnDelete         string
	OnUpdate         string
}

// Schema es el builder principal
type Schema struct{}

// NewSchema crea una nueva instancia de Schema
func NewSchema() *Schema {
	return &Schema{}
}

// Create crea una nueva tabla
func (s *Schema) Create(tableName string, callback func(*Blueprint)) string {
	blueprint := &Blueprint{
		tableName: tableName,
		columns:   []*Column{},
		indexes:   []Index{},
		foreign:   []ForeignKey{},
	}

	callback(blueprint)

	return blueprint.ToSQL()
}

// Drop elimina una tabla
func (s *Schema) Drop(tableName string) string {
	return fmt.Sprintf("DROP TABLE IF EXISTS %s", tableName)
}

// Increments crea una columna AUTO_INCREMENT PRIMARY KEY
func (b *Blueprint) Increments(name string) *Column {
	col := &Column{
		Name:          name,
		Type:          "INT",
		IsUnsigned:    true,
		AutoIncrement: true,
		IsPrimary:     true,
		IsNullable:    false,
	}
	b.columns = append(b.columns, col)
	return col
}

// Integer crea una columna INTEGER
func (b *Blueprint) Integer(name string) *Column {
	col := &Column{
		Name:       name,
		Type:       "INT",
		IsNullable: false,
	}
	b.columns = append(b.columns, col)
	return col
}

// UnsignedInteger crea una columna INTEGER UNSIGNED
func (b *Blueprint) UnsignedInteger(name string) *Column {
	col := &Column{
		Name:       name,
		Type:       "INT",
		IsUnsigned: true,
		IsNullable: false,
	}
	b.columns = append(b.columns, col)
	return col
}

// BigInteger crea una columna BIGINT
func (b *Blueprint) BigInteger(name string) *Column {
	col := &Column{
		Name:       name,
		Type:       "BIGINT",
		IsNullable: false,
	}
	b.columns = append(b.columns, col)
	return col
}

// String crea una columna VARCHAR
func (b *Blueprint) String(name string, length ...int) *Column {
	l := 255 // default length
	if len(length) > 0 {
		l = length[0]
	}

	col := &Column{
		Name:       name,
		Type:       "VARCHAR",
		Length:     l,
		IsNullable: false,
	}
	b.columns = append(b.columns, col)
	return col
}

// Text crea una columna TEXT
func (b *Blueprint) Text(name string) *Column {
	col := &Column{
		Name:       name,
		Type:       "TEXT",
		IsNullable: false,
	}
	b.columns = append(b.columns, col)
	return col
}

// Boolean crea una columna BOOLEAN
func (b *Blueprint) Boolean(name string) *Column {
	col := &Column{
		Name:       name,
		Type:       "BOOLEAN",
		IsNullable: false,
	}
	b.columns = append(b.columns, col)
	return col
}

// DateTime crea una columna DATETIME
func (b *Blueprint) DateTime(name string) *Column {
	col := &Column{
		Name:       name,
		Type:       "DATETIME",
		IsNullable: false,
	}
	b.columns = append(b.columns, col)
	return col
}

// Timestamp crea una columna TIMESTAMP
func (b *Blueprint) Timestamp(name string) *Column {
	col := &Column{
		Name:       name,
		Type:       "TIMESTAMP",
		IsNullable: false,
	}
	b.columns = append(b.columns, col)
	return col
}

// Decimal crea una columna DECIMAL
func (b *Blueprint) Decimal(name string, precision, scale int) *Column {
	col := &Column{
		Name:       name,
		Type:       "DECIMAL",
		Precision:  precision,
		Scale:      scale,
		IsNullable: false,
	}
	b.columns = append(b.columns, col)
	return col
}

// Uuid crea una columna UUID (BINARY(16) o CHAR(36))
func (b *Blueprint) Uuid(name string) *Column {
	col := &Column{
		Name:       name,
		Type:       "CHAR",
		Length:     36,
		IsNullable: false,
	}
	b.columns = append(b.columns, col)
	return col
}

// UuidPrimary crea una columna UUID PRIMARY KEY
func (b *Blueprint) UuidPrimary(name string) *Column {
	col := &Column{
		Name:       name,
		Type:       "CHAR",
		Length:     36,
		IsPrimary:  true,
		IsNullable: false,
	}
	b.columns = append(b.columns, col)
	return col
}

// BinaryUuid crea una columna UUID en formato binario (más eficiente)
func (b *Blueprint) BinaryUuid(name string) *Column {
	col := &Column{
		Name:       name,
		Type:       "BINARY",
		Length:     16,
		IsNullable: false,
	}
	b.columns = append(b.columns, col)
	return col
}

// Enum crea una columna ENUM
func (b *Blueprint) Enum(name string, values []string) *Column {
	col := &Column{
		Name:       name,
		Type:       "ENUM",
		EnumValues: values,
		IsNullable: false,
	}
	b.columns = append(b.columns, col)
	return col
}

// Set crea una columna SET
func (b *Blueprint) Set(name string, values []string) *Column {
	col := &Column{
		Name:       name,
		Type:       "SET",
		EnumValues: values,
		IsNullable: false,
	}
	b.columns = append(b.columns, col)
	return col
}

// Float crea una columna FLOAT
func (b *Blueprint) Float(name string, precision ...int) *Column {
	p := 8 // default precision
	if len(precision) > 0 {
		p = precision[0]
	}

	col := &Column{
		Name:       name,
		Type:       "FLOAT",
		Precision:  p,
		IsNullable: false,
	}
	b.columns = append(b.columns, col)
	return col
}

// Double crea una columna DOUBLE
func (b *Blueprint) Double(name string, precision ...int) *Column {
	p := 15 // default precision
	if len(precision) > 0 {
		p = precision[0]
	}

	col := &Column{
		Name:       name,
		Type:       "DOUBLE",
		Precision:  p,
		IsNullable: false,
	}
	b.columns = append(b.columns, col)
	return col
}

// Binary crea una columna BINARY
func (b *Blueprint) Binary(name string, length int) *Column {
	col := &Column{
		Name:       name,
		Type:       "BINARY",
		Length:     length,
		IsNullable: false,
	}
	b.columns = append(b.columns, col)
	return col
}

// VarBinary crea una columna VARBINARY
func (b *Blueprint) VarBinary(name string, length int) *Column {
	col := &Column{
		Name:       name,
		Type:       "VARBINARY",
		Length:     length,
		IsNullable: false,
	}
	b.columns = append(b.columns, col)
	return col
}

// Json crea una columna JSON
func (b *Blueprint) Json(name string) *Column {
	col := &Column{
		Name:       name,
		Type:       "JSON",
		IsNullable: false,
	}
	b.columns = append(b.columns, col)
	return col
}

// Nullable hace que la columna permita NULL
func (c *Column) Nullable() *Column {
	c.IsNullable = true
	return c
}

// Default establece un valor por defecto
func (c *Column) Default(value interface{}) *Column {
	c.DefaultValue = value
	return c
}

// Unsigned hace que la columna sea unsigned (solo para enteros)
func (c *Column) Unsigned() *Column {
	c.IsUnsigned = true
	return c
}

// Unique hace que la columna sea única
func (c *Column) Unique() *Column {
	c.IsUnique = true
	return c
}

// Index crea un índice para la columna
func (c *Column) Index() *Column {
	c.HasIndex = true
	return c
}

// After especifica que la columna debe ser colocada después de otra columna
func (c *Column) After(columnName string) *Column {
	// En una implementación completa, esto se manejaría en el SQL ALTER TABLE
	c.CommentText = fmt.Sprintf("AFTER %s", columnName)
	return c
}

// First especifica que la columna debe ser la primera
func (c *Column) First() *Column {
	// En una implementación completa, esto se manejaría en el SQL ALTER TABLE
	c.CommentText = "FIRST"
	return c
}

// Charset especifica el charset para columnas de texto
func (c *Column) Charset(charset string) *Column {
	// Esto se podría implementar como parte del tipo de datos
	if c.CommentText == "" {
		c.CommentText = fmt.Sprintf("CHARSET %s", charset)
	} else {
		c.CommentText += fmt.Sprintf(" CHARSET %s", charset)
	}
	return c
}

// Collation especifica la collation para columnas de texto
func (c *Column) Collation(collation string) *Column {
	if c.CommentText == "" {
		c.CommentText = fmt.Sprintf("COLLATE %s", collation)
	} else {
		c.CommentText += fmt.Sprintf(" COLLATE %s", collation)
	}
	return c
}

// UseCurrent establece DEFAULT CURRENT_TIMESTAMP para columnas timestamp
func (c *Column) UseCurrent() *Column {
	if c.Type == "TIMESTAMP" || c.Type == "DATETIME" {
		c.DefaultValue = "CURRENT_TIMESTAMP"
	}
	return c
}

// OnUpdateCurrent establece ON UPDATE CURRENT_TIMESTAMP
func (c *Column) OnUpdateCurrent() *Column {
	if c.Type == "TIMESTAMP" || c.Type == "DATETIME" {
		c.OnUpdate = "CURRENT_TIMESTAMP"
	}
	return c
}

// Timestamps añade created_at y updated_at
func (b *Blueprint) Timestamps() {
	b.Timestamp("created_at")
	b.Timestamp("updated_at")
}

// SoftDeletes añade deleted_at
func (b *Blueprint) SoftDeletes() {
	b.Timestamp("deleted_at").Nullable()
}

// UniqueIndex crea un índice único
func (b *Blueprint) UniqueIndex(columns []string, name string) {
	index := Index{
		Name:    name,
		Columns: columns,
		Type:    "unique",
	}
	b.indexes = append(b.indexes, index)
}

// Foreign crea una clave foránea
func (b *Blueprint) Foreign(column string) *ForeignKeyBuilder {
	return &ForeignKeyBuilder{
		blueprint: b,
		column:    column,
	}
}

type ForeignKeyBuilder struct {
	blueprint        *Blueprint
	column           string
	referencedTable  string
	referencedColumn string
	onDelete         string
	onUpdate         string
}

func (fkb *ForeignKeyBuilder) References(column string) *ForeignKeyBuilder {
	fkb.referencedColumn = column
	return fkb
}

func (fkb *ForeignKeyBuilder) On(table string) *ForeignKeyBuilder {
	fkb.referencedTable = table
	return fkb
}

func (fkb *ForeignKeyBuilder) OnDelete(action string) *ForeignKeyBuilder {
	fkb.onDelete = action

	// Añadir la clave foránea al blueprint cuando se especifica OnDelete
	fk := ForeignKey{
		Column:           fkb.column,
		ReferencedTable:  fkb.referencedTable,
		ReferencedColumn: fkb.referencedColumn,
		OnDelete:         fkb.onDelete,
		OnUpdate:         fkb.onUpdate,
	}
	fkb.blueprint.foreign = append(fkb.blueprint.foreign, fk)
	return fkb
}

func (fkb *ForeignKeyBuilder) OnUpdate(action string) *ForeignKeyBuilder {
	fkb.onUpdate = action
	return fkb
}

// ToSQL genera el SQL DDL para la tabla
func (b *Blueprint) ToSQL() string {
	var sql strings.Builder

	sql.WriteString(fmt.Sprintf("CREATE TABLE %s (\n", b.tableName))

	// Columnas
	columnDefs := make([]string, 0, len(b.columns))
	for _, col := range b.columns {
		columnDefs = append(columnDefs, b.columnToSQL(col))
	}

	// Añadir índices únicos
	for _, col := range b.columns {
		if col.IsUnique && !col.IsPrimary {
			columnDefs = append(columnDefs, fmt.Sprintf("UNIQUE KEY (%s)", col.Name))
		}
		if col.HasIndex && !col.IsPrimary && !col.IsUnique {
			columnDefs = append(columnDefs, fmt.Sprintf("KEY (%s)", col.Name))
		}
	}

	// Añadir índices personalizados
	for _, index := range b.indexes {
		indexSQL := b.indexToSQL(index)
		if indexSQL != "" {
			columnDefs = append(columnDefs, indexSQL)
		}
	}

	// Añadir claves foráneas
	for _, fk := range b.foreign {
		fkSQL := b.foreignKeyToSQL(fk)
		if fkSQL != "" {
			columnDefs = append(columnDefs, fkSQL)
		}
	}

	// Añadir clave primaria compuesta si existe
	if len(b.primaryKey) > 0 {
		primaryKeyColumns := strings.Join(b.primaryKey, ", ")
		columnDefs = append(columnDefs, fmt.Sprintf("PRIMARY KEY (%s)", primaryKeyColumns))
	}

	sql.WriteString("\t")
	sql.WriteString(strings.Join(columnDefs, ",\n\t"))
	sql.WriteString("\n);")

	return sql.String()
}

// columnToSQL convierte una columna a SQL
func (b *Blueprint) columnToSQL(col *Column) string {
	var parts []string

	// Nombre de la columna
	parts = append(parts, col.Name)

	// Tipo de datos
	dataType := b.getDataType(col)
	parts = append(parts, dataType)

	// Unsigned
	if col.IsUnsigned && (col.Type == "INT" || col.Type == "BIGINT") {
		parts = append(parts, "UNSIGNED")
	}

	// NOT NULL / NULL
	if col.IsNullable {
		// Para campos nullable no ponemos nada o ponemos NULL explícitamente
	} else {
		parts = append(parts, "NOT NULL")
	}

	// Default value
	if col.DefaultValue != nil {
		switch v := col.DefaultValue.(type) {
		case string:
			switch v {
			case "NULL":
				parts = append(parts, "DEFAULT NULL")
			case "CURRENT_TIMESTAMP":
				parts = append(parts, "DEFAULT CURRENT_TIMESTAMP")
			default:
				parts = append(parts, fmt.Sprintf("DEFAULT '%s'", v))
			}
		case bool:
			if v {
				parts = append(parts, "DEFAULT true")
			} else {
				parts = append(parts, "DEFAULT false")
			}
		default:
			parts = append(parts, fmt.Sprintf("DEFAULT %v", v))
		}
	}

	// Auto increment
	if col.AutoIncrement {
		parts = append(parts, "AUTO_INCREMENT")
	}

	// Primary key
	if col.IsPrimary {
		parts = append(parts, "PRIMARY KEY")
	}

	// ON UPDATE modifier
	if col.OnUpdate != "" {
		parts = append(parts, fmt.Sprintf("ON UPDATE %s", col.OnUpdate))
	}

	// Modificadores avanzados
	if col.CommentText != "" {
		parts = append(parts, col.CommentText)
	}

	return strings.Join(parts, " ")
}

// getDataType obtiene el tipo de datos SQL
func (b *Blueprint) getDataType(col *Column) string {
	switch col.Type {
	case "VARCHAR":
		return b.getVarcharType(col)
	case "CHAR":
		return b.getCharType(col)
	case "BINARY":
		return b.getBinaryType(col)
	case "VARBINARY":
		return b.getVarBinaryType(col)
	case "DECIMAL":
		return b.getDecimalType(col)
	case "FLOAT":
		return b.getFloatType(col)
	case "DOUBLE":
		return b.getDoubleType(col)
	case "ENUM":
		return b.getEnumOrSetType("ENUM", col.EnumValues)
	case "SET":
		return b.getEnumOrSetType("SET", col.EnumValues)
	default:
		return col.Type
	}
}

func (b *Blueprint) getVarcharType(col *Column) string {
	return fmt.Sprintf("VARCHAR(%d)", col.Length)
}

func (b *Blueprint) getCharType(col *Column) string {
	if col.Length > 0 {
		return fmt.Sprintf("CHAR(%d)", col.Length)
	}
	return "CHAR(255)"
}

func (b *Blueprint) getBinaryType(col *Column) string {
	if col.Length > 0 {
		return fmt.Sprintf("BINARY(%d)", col.Length)
	}
	return "BINARY(16)"
}

func (b *Blueprint) getVarBinaryType(col *Column) string {
	if col.Length > 0 {
		return fmt.Sprintf("VARBINARY(%d)", col.Length)
	}
	return "VARBINARY(255)"
}

func (b *Blueprint) getDecimalType(col *Column) string {
	return fmt.Sprintf("DECIMAL(%d,%d)", col.Precision, col.Scale)
}

func (b *Blueprint) getFloatType(col *Column) string {
	if col.Precision > 0 {
		return fmt.Sprintf("FLOAT(%d)", col.Precision)
	}
	return "FLOAT"
}

func (b *Blueprint) getDoubleType(col *Column) string {
	if col.Precision > 0 {
		return fmt.Sprintf("DOUBLE(%d)", col.Precision)
	}
	return "DOUBLE"
}

func (b *Blueprint) getEnumOrSetType(typeName string, values []string) string {
	if len(values) > 0 {
		quoted := make([]string, len(values))
		for i, v := range values {
			quoted[i] = fmt.Sprintf("'%s'", v)
		}
		return fmt.Sprintf("%s(%s)", typeName, strings.Join(quoted, ", "))
	}
	return fmt.Sprintf("%s('')", typeName)
}

// indexToSQL convierte un índice a SQL
func (b *Blueprint) indexToSQL(index Index) string {
	columns := strings.Join(index.Columns, ", ")
	switch index.Type {
	case "unique":
		return fmt.Sprintf("UNIQUE KEY %s (%s)", index.Name, columns)
	case "index":
		return fmt.Sprintf("KEY %s (%s)", index.Name, columns)
	default:
		return ""
	}
}

// foreignKeyToSQL convierte una clave foránea a SQL
func (b *Blueprint) foreignKeyToSQL(fk ForeignKey) string {
	sql := fmt.Sprintf("FOREIGN KEY (%s) REFERENCES %s (%s)",
		fk.Column, fk.ReferencedTable, fk.ReferencedColumn)

	if fk.OnDelete != "" {
		sql += fmt.Sprintf(" ON DELETE %s", fk.OnDelete)
	}

	if fk.OnUpdate != "" {
		sql += fmt.Sprintf(" ON UPDATE %s", fk.OnUpdate)
	}

	return sql
}

// UuidTimestamps añade created_at y updated_at con configuración automática (NO incluye id)
func (b *Blueprint) UuidTimestamps() {
	b.Timestamp("created_at").UseCurrent()
	b.Timestamp("updated_at").UseCurrent().OnUpdateCurrent()
}

// Id es un alias para Increments("id")
func (b *Blueprint) Id() *Column {
	return b.Increments("id")
}

// UuidId crea una columna id de tipo UUID
func (b *Blueprint) UuidId() *Column {
	return b.UuidPrimary("id")
}

// MorphsUuid crea columnas para relaciones polimórficas con UUID
func (b *Blueprint) MorphsUuid(name string) {
	b.Uuid(name + "_id").Index()
	b.String(name + "_type").Index()
}

// Morphs crea columnas para relaciones polimórficas con ID entero
func (b *Blueprint) Morphs(name string) {
	b.Integer(name + "_id").Index()
	b.String(name + "_type").Index()
}

// RememberToken añade una columna remember_token
func (b *Blueprint) RememberToken() {
	b.String("remember_token", 100).Nullable()
}

// Primary define una clave primaria compuesta
func (b *Blueprint) Primary(columns []string) {
	b.primaryKey = columns
}

// Index añade un índice simple a una columna
func (b *Blueprint) Index(column string) {
	b.indexes = append(b.indexes, Index{
		Name:    fmt.Sprintf("idx_%s_%s", b.tableName, column),
		Columns: []string{column},
		Type:    "index",
	})
}

// Unique añade un índice único compuesto a múltiples columnas
func (b *Blueprint) Unique(columns []string) {
	columnNames := strings.Join(columns, "_")
	b.indexes = append(b.indexes, Index{
		Name:    fmt.Sprintf("unique_%s_%s", b.tableName, columnNames),
		Columns: columns,
		Type:    "unique",
	})
}
