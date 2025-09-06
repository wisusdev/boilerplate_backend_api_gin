package repositories

import (
	"semita/app/data/models"
	"semita/core/common/nulltypes"
	"semita/core/database/database_connections"
	"semita/core/database/orm"
	"semita/core/helpers"
	"time"

	"github.com/google/uuid"
)

var userTable = "users"

type UserRepository struct {
	*orm.ModelRepository
}

func NewUserRepository() *UserRepository {
	return &UserRepository{
		ModelRepository: orm.NewModelRepository(&models.UserStruct{}),
	}
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

// scanUserRow escanea una fila de usuario y maneja la conversión de fechas
func scanUserRow(scanner interface {
	Scan(dest ...interface{}) error
}) (models.UserStruct, error) {
	var user models.UserStruct
	var createdAtStr, updatedAtStr string
	var emailVerifiedAtStr nulltypes.NullString
	var avatarPtr *string

	err := scanner.Scan(
		&user.ID, &user.FirstName, &user.LastName, &user.Username,
		&avatarPtr, &user.Language, &user.Email, &emailVerifiedAtStr, &user.Password,
		&createdAtStr, &updatedAtStr,
	)

	if err != nil {
		return models.UserStruct{}, err
	}

	// Convertir las fechas usando la función helper
	user.CreatedAt = parseDateTime(createdAtStr)
	user.UpdatedAt = parseDateTime(updatedAtStr)

	if emailVerifiedAtStr.Valid && emailVerifiedAtStr.String != "" {
		user.EmailVerifiedAt = parseDateTimePtr(emailVerifiedAtStr.String)
	} else {
		user.EmailVerifiedAt = nil // nil si no está verificado
	}

	return user, nil
}

func (r *UserRepository) Where(field string, value interface{}) ([]models.UserStruct, error) {
	query := "SELECT id, first_name, last_name, username, avatar, language, email, email_verified_at, password, created_at, updated_at FROM users WHERE " + field + " = ?"
	rows, err := r.DB.Query(query, value)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []models.UserStruct
	for rows.Next() {
		user, err := scanUserRow(rows)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, nil
}

func (r *UserRepository) GetAllUsers() ([]models.UserStruct, error) {
	var users models.Users

	errorGetAll := r.GetAll(&users)
	if errorGetAll != nil {
		return nil, errorGetAll
	}

	return users, nil
}

func StoreUser(storeUser models.UserStruct) (user models.UserStruct, err error) {
	// Instanciamos la conexion a la base de datos
	var database = database_connections.DatabaseConnectSQL()
	defer database.Close()

	// Generar UUID para el ID
	id := uuid.New().String()
	storeUser.ID = id

	// Preparamos la consulta para insertar un nuevo usuario
	var query = "INSERT INTO " + userTable + " (id, first_name, last_name, username, email, password, language, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, 'es', NOW(), NOW())"

	// Ejecutamos la consulta con los datos del usuario
	_, errorErr := database.Exec(query, storeUser.ID, storeUser.FirstName, storeUser.LastName, storeUser.Username, storeUser.Email, storeUser.Password)
	if errorErr != nil {
		helpers.Logs("ERROR", "Error al guardar el usuario: "+errorErr.Error())
		return models.UserStruct{}, errorErr
	}

	user, err = GetUserByID(storeUser.ID)
	if err != nil {
		helpers.Logs("ERROR", "Error al obtener el usuario recién insertado: "+err.Error())
		return models.UserStruct{}, err
	}

	return user, nil
}

func GetUserByID(id string) (user models.UserStruct, err error) {
	// Instanciamos la conexión a la base de datos
	var database = database_connections.DatabaseConnectSQL()

	// Aseguramos que la conexión se cierre al final de la función
	defer database.Close()

	// Preparamos la consulta para obtener un usuario por su ID
	var query = "SELECT id, first_name, last_name, username, avatar, language, email, email_verified_at, password, created_at, updated_at FROM " + userTable + " WHERE id = ?"

	// Ejecutamos la consulta y obtenemos los resultados usando la función helper
	user, err = scanUserRow(database.QueryRow(query, id))
	if err != nil {
		return models.UserStruct{}, err
	}

	return user, nil
}

func GetUserByEmail(email string) (user models.UserStruct, err error) {
	// Instanciamos la conexión a la base de datos
	var database = database_connections.DatabaseConnectSQL()

	// Aseguramos que la conexión se cierre al final de la función
	defer database.Close()

	// Preparamos la consulta para obtener un usuario por su email
	var query = "SELECT id, first_name, last_name, username, avatar, language, email, email_verified_at, password, created_at, updated_at FROM " + userTable + " WHERE email = ?"

	// Ejecutamos la consulta y obtenemos los resultados usando la función helper
	user, err = scanUserRow(database.QueryRow(query, email))
	if err != nil {
		helpers.Logs("ERROR", "Error al obtener el usuario por email: "+err.Error())
		return models.UserStruct{}, err
	}

	return user, nil
}

func UpdateUserPassword(user models.UserStruct) (err error) {
	// Instanciamos la conexión a la base de datos
	var database = database_connections.DatabaseConnectSQL()

	// Aseguramos que la conexión se cierre al final de la función
	defer database.Close()

	// Preparamos la consulta para actualizar un usuario por su ID
	var query = "UPDATE " + userTable + " SET password = ? WHERE id = ?"

	// Ejecutamos la consulta con los datos del usuario
	_, err = database.Exec(query, user.Password, user.ID)

	// Si hubo un error al ejecutar la consulta, retornamos el error
	if err != nil {
		return err
	}

	return nil
}

func DeleteUser(id string) (err error) {
	// Instanciamos la conexión a la base de datos
	var database = database_connections.DatabaseConnectSQL()

	// Aseguramos que la conexión se cierre al final de la función
	defer database.Close()

	// Preparamos la consulta para eliminar un usuario por su ID
	var query = "DELETE FROM " + userTable + " WHERE id = ?"

	// Ejecutamos la consulta con el ID del usuario
	_, err = database.Exec(query, id)

	// Si hubo un error al ejecutar la consulta, retornamos el error
	if err != nil {
		return err
	}

	return nil
}

// MarkEmailVerified actualiza el campo email_verified_at del usuario
func MarkEmailVerified(userID string) error {
	db := database_connections.DatabaseConnectSQL()
	defer db.Close()
	_, err := db.Exec("UPDATE "+userTable+" SET email_verified_at = ? WHERE id = ?", time.Now().Format("2006-01-02 15:04:05"), userID)
	return err
}
