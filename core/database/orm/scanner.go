package orm

import (
	"errors"
	"reflect"
	"semita/core/common/nulltypes"
	"semita/core/helpers"
	"strings"
	"time"
)

// Scanner es una interfaz para escanear filas de la base de datos
type Scanner interface {
	Scan(dest ...interface{}) error
}

// RowScanner escanea una fila de la base de datos en una estructura
func RowScanner(row Scanner, dest interface{}) error {
	v := reflect.ValueOf(dest)
	if v.Kind() != reflect.Ptr || v.IsNil() {
		return errors.New("dest must be a non-nil pointer")
	}

	v = v.Elem()
	if v.Kind() != reflect.Struct {
		return errors.New("dest must be a pointer to a struct")
	}

	t := v.Type()

	// Primero recopilamos todos los campos que necesitamos escanear
	var fieldPtrs []interface{}
	for i := 0; i < t.NumField(); i++ {
		fieldValue := v.Field(i)

		// Ignorar campos no exportados
		if !fieldValue.CanSet() {
			continue
		}

		// Manejar diferentes tipos de campos
		switch fieldValue.Kind() {
		case reflect.String, reflect.Int, reflect.Int64, reflect.Bool, reflect.Float64:
			fieldPtrs = append(fieldPtrs, fieldValue.Addr().Interface())
		case reflect.Struct:
			// Verificar si es nulltypes.NullString
			if fieldValue.Type() == reflect.TypeOf(nulltypes.NullString{}) {
				fieldPtrs = append(fieldPtrs, fieldValue.Addr().Interface())
			} else if fieldValue.Type() == reflect.TypeOf(time.Time{}) {
				// Manejo especial para time.Time
				var timeStr string
				fieldPtrs = append(fieldPtrs, &timeStr)
			} else {
				fieldPtrs = append(fieldPtrs, fieldValue.Addr().Interface())
			}
		case reflect.Ptr:
			// Manejo especial para *time.Time
			if fieldValue.Type() == reflect.TypeOf((*time.Time)(nil)) {
				var timeStr nulltypes.NullString
				fieldPtrs = append(fieldPtrs, &timeStr)
			} else {
				// Para otros punteros, usar un puntero nulo del tipo adecuado
				ptrType := fieldValue.Type().Elem()
				newPtr := reflect.New(ptrType)
				fieldPtrs = append(fieldPtrs, newPtr.Interface())
			}
		default:
			fieldPtrs = append(fieldPtrs, fieldValue.Addr().Interface())
		}
	}

	// Escanear la fila en los punteros recopilados
	err := row.Scan(fieldPtrs...)
	if err != nil {
		return err
	}

	// Procesar los resultados escaneados
	fieldPtrIndex := 0
	for i := 0; i < t.NumField(); i++ {
		fieldValue := v.Field(i)

		// Ignorar campos no exportados
		if !fieldValue.CanSet() {
			continue
		}

		// Manejar diferentes tipos de campos
		switch fieldValue.Kind() {
		case reflect.String, reflect.Int, reflect.Int64, reflect.Bool, reflect.Float64:
			// Estos tipos ya están escaneados directamente
		case reflect.Struct:
			// Verificar si es nulltypes.NullString
			if fieldValue.Type() == reflect.TypeOf(nulltypes.NullString{}) {
				// Ya está escaneado directamente
			} else if fieldValue.Type() == reflect.TypeOf(time.Time{}) {
				// Convertir string a time.Time
				timeStr := *fieldPtrs[fieldPtrIndex].(*string)
				parsedTime := parseDateTime(timeStr)
				fieldValue.Set(reflect.ValueOf(parsedTime))
			}
		case reflect.Ptr:
			// Convertir nulltypes.NullString a *time.Time
			if fieldValue.Type() == reflect.TypeOf((*time.Time)(nil)) {
				timeStr := *fieldPtrs[fieldPtrIndex].(*nulltypes.NullString)
				if timeStr.Valid && timeStr.String != "" {
					parsedTime := parseDateTimePtr(timeStr.String)
					fieldValue.Set(reflect.ValueOf(parsedTime))
				} else {
					fieldValue.Set(reflect.Zero(fieldValue.Type()))
				}
			} else {
				// Para otros punteros, establecer el valor escaneado
				ptr := fieldPtrs[fieldPtrIndex]
				ptrValue := reflect.ValueOf(ptr).Elem()
				if !ptrValue.IsZero() {
					fieldValue.Set(ptrValue.Addr())
				}
			}
		default:
			// Para cualquier otro tipo, intentar escanearlo directamente
			helpers.Logs("WARNING", "Tipo no manejado en RowScanner: "+fieldValue.Type().String())
		}

		fieldPtrIndex++
	}

	return nil
}

// parseDateTime convierte una fecha string a time.Time
func parseDateTime(dateStr string) time.Time {
	if dateStr == "" {
		return time.Now()
	}

	parsedTime, err := time.Parse("2006-01-02 15:04:05", dateStr)
	if err != nil {
		helpers.Logs("ERROR", "Error al parsear fecha: "+err.Error())
		return time.Now()
	}

	return parsedTime
}

// parseDateTimePtr convierte una fecha string a *time.Time o nil si está vacía
func parseDateTimePtr(dateStr string) *time.Time {
	if dateStr == "" {
		return nil
	}
	parsedTime, err := time.Parse("2006-01-02 15:04:05", dateStr)
	if err != nil {
		helpers.Logs("ERROR", "Error al parsear fecha: "+err.Error())
		return nil
	}
	return &parsedTime
}

// GetColumnNames extrae los nombres de columna de una estructura
func GetColumnNames(structType reflect.Type) []string {
	var columns []string

	for i := 0; i < structType.NumField(); i++ {
		field := structType.Field(i)

		// Obtener el nombre de columna desde la etiqueta json
		colName := field.Tag.Get("json")
		if colName == "" {
			colName = strings.ToLower(field.Name)
		} else {
			// Si tiene opciones como omitempty, las quitamos
			colName = strings.Split(colName, ",")[0]
		}

		columns = append(columns, colName)
	}

	return columns
}
