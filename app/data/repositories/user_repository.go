package repositories

import (
	"semita/app/data/structs"
	"semita/core/database/database_connections"
	"semita/core/helpers"
	"time"
)

var userTable = "users"

type UserRepository struct {
	DB database_connections.SQLAdapter
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

// scanUserRow escanea una fila de usuario y maneja la conversión de fechas
func scanUserRow(scanner interface {
	Scan(dest ...interface{}) error
}) (structs.UserStruct, error) {
	var user structs.UserStruct
	var createdAtStr, updatedAtStr string
	var avatarPtr *string

	err := scanner.Scan(
		&user.ID, &user.FirstName, &user.LastName, &user.Username,
		&avatarPtr, &user.Language, &user.Email, &user.Password,
		&createdAtStr, &updatedAtStr,
	)

	if avatarPtr != nil {
		user.Avatar = *avatarPtr
	} else {
		user.Avatar = ""
	}

	if err != nil {
		return structs.UserStruct{}, err
	}

	// Convertir las fechas usando la función helper
	user.CreatedAt = parseDateTime(createdAtStr)
	user.UpdatedAt = parseDateTime(updatedAtStr)

	return user, nil
}

func (r *UserRepository) Where(field string, value interface{}) ([]structs.UserStruct, error) {
	query := "SELECT id, first_name, last_name, username, avatar, language, email, password, created_at, updated_at FROM users WHERE " + field + " = ?"
	rows, err := r.DB.Query(query, value)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []structs.UserStruct
	for rows.Next() {
		user, err := scanUserRow(rows)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, nil
}

func GetAllUsers() ([]structs.UserStruct, error) {
	// Instanciamos la conexión a la base de datos
	var database = database_connections.DatabaseConnectSQL()

	// Aseguramos que la conexión se cierre al final de la función
	defer database.Close()

	// Preparamos la consulta para obtener todos los usuarios
	var query = "SELECT id, first_name, last_name, username, avatar, language, email, password, created_at, updated_at FROM " + userTable

	// Ejecutamos la consulta y obtenemos los resultados
	rows, err := database.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Creamos un slice para almacenar los usuarios
	var users []structs.UserStruct

	// Iteramos sobre los resultados y los agregamos al slice
	for rows.Next() {
		user, err := scanUserRow(rows)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}

func StoreUser(user structs.StoreUserStruct) (err error) {
	// Instanciamos la conexion a la base de datos
	var database = database_connections.DatabaseConnectSQL()

	// Aseguramos que la conexion se cierre al final de la funcion
	defer database.Close()

	// Preparamos la consulta para insertar un nuevo usuario
	var query = "INSERT INTO " + userTable + " (first_name, last_name, username, email, password, language) VALUES (?, '', ?, ?, ?, 'es')"

	// Ejecutamos la consulta con los datos del usuario
	_, err = database.Exec(query, user.Name, user.Email, user.Email, user.Password)

	// Si hubo un error al ejecutar la consulta, retornamos el error
	if err != nil {
		return err
	}

	return nil
}

func GetUserByID(id string) (user structs.UserStruct, err error) {
	// Instanciamos la conexión a la base de datos
	var database = database_connections.DatabaseConnectSQL()

	// Aseguramos que la conexión se cierre al final de la función
	defer database.Close()

	// Preparamos la consulta para obtener un usuario por su ID
	var query = "SELECT id, first_name, last_name, username, avatar, language, email, password, created_at, updated_at FROM " + userTable + " WHERE id = ?"

	// Ejecutamos la consulta y obtenemos los resultados usando la función helper
	user, err = scanUserRow(database.QueryRow(query, id))
	if err != nil {
		return structs.UserStruct{}, err
	}

	return user, nil
}

func GetUserByEmail(email string) (user structs.UserStruct, err error) {
	// Instanciamos la conexión a la base de datos
	var database = database_connections.DatabaseConnectSQL()

	// Aseguramos que la conexión se cierre al final de la función
	defer database.Close()

	// Preparamos la consulta para obtener un usuario por su email
	var query = "SELECT id, first_name, last_name, username, avatar, language, email, password, created_at, updated_at FROM " + userTable + " WHERE email = ?"

	// Ejecutamos la consulta y obtenemos los resultados usando la función helper
	user, err = scanUserRow(database.QueryRow(query, email))
	if err != nil {
		helpers.Logs("ERROR", "Error al obtener el usuario por email: "+err.Error())
		return structs.UserStruct{}, err
	}

	return user, nil
}

func UpdateUser(user structs.UpdateUserStruct) (err error) {
	// Instanciamos la conexión a la base de datos
	var database = database_connections.DatabaseConnectSQL()

	// Aseguramos que la conexión se cierre al final de la función
	defer database.Close()

	// Preparamos la consulta para actualizar un usuario por su ID
	var query = "UPDATE " + userTable + " SET first_name = ?, email = ?, password = ? WHERE id = ?"

	// Ejecutamos la consulta con los datos del usuario
	_, err = database.Exec(query, user.Name, user.Email, user.Password, user.ID)

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
func MarkEmailVerified(userID int) error {
	db := database_connections.DatabaseConnectSQL()
	defer db.Close()
	_, err := db.Exec("UPDATE "+userTable+" SET email_verified_at = ? WHERE id = ?", time.Now().Format("2006-01-02 15:04:05"), userID)
	return err
}
