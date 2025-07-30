package examples

import (
	"database/sql"
	"semita/app/data/models"
	"semita/core/database/orm"
	"strconv"
)

// Este archivo muestra cómo migrar del código existente al nuevo ORM

// Ejemplo de migración de GetAllUsers
func GetAllUsersWithORM() ([]models.UserStruct, error) {
	// Usando el nuevo ORM
	repo := NewUserRepository()

	// Obtener todos los usuarios
	users, err := repo.GetAllUsers()
	if err != nil {
		return nil, err
	}

	return users, nil
}

// Ejemplo de migración de GetUserByID
func GetUserByIDWithORM(id string) (models.UserStruct, error) {
	// Usando el nuevo ORM
	repo := NewUserRepository()

	// Convertir ID a int si es necesario
	idInt, _ := strconv.Atoi(id)

	// Obtener usuario por ID
	user, err := repo.GetUserByID(idInt)
	if err != nil {
		return models.UserStruct{}, err
	}

	return user, nil
}

// Ejemplo de migración de GetUserByEmail
func GetUserByEmailWithORM(email string) (models.UserStruct, error) {
	// Usando el nuevo ORM
	repo := NewUserRepository()

	// Obtener usuario por email
	user, err := repo.GetUserByEmail(email)
	if err != nil {
		return models.UserStruct{}, err
	}

	return user, nil
}

// Ejemplo de migración de StoreUser
func StoreUserWithORM(storeUser models.UserStruct) (models.UserStruct, error) {
	// Usando el nuevo ORM
	repo := NewUserRepository()

	// Guardar el usuario
	user, err := repo.StoreUser(storeUser)
	if err != nil {
		return models.UserStruct{}, err
	}

	return user, nil
}

// Ejemplo de migración de UpdateUser
func UpdateUserWithORM(user models.UserStruct) error {
	// Usando el nuevo ORM
	repo := NewUserRepository()

	// Actualizar el usuario
	return repo.UpdateUser(user)
}

// Ejemplo de migración de DeleteUser
func DeleteUserWithORM(id string) (sql.Result, error) {
	// Usando el nuevo ORM
	repo := NewUserRepository()

	// Convertir ID a int si es necesario
	idInt, _ := strconv.Atoi(id)

	// Eliminar el usuario
	return repo.DeleteUser(idInt)
}

// Ejemplo de migración de consultas complejas usando QueryBuilder
func FindUsersWithFiltersORM(minAge int, country string, active bool) ([]models.UserStruct, error) {
	// Usando el nuevo ORM
	repo := NewUserRepository()

	// Construir la consulta
	query := repo.Query()

	if minAge > 0 {
		query.Where("age >= ?", minAge)
	}

	if country != "" {
		query.Where("country = ?", country)
	}

	if active {
		query.Where("email_verified_at IS NOT NULL")
	}

	// Ordenar por nombre
	query.OrderBy("last_name ASC, first_name ASC")

	// Ejecutar la consulta
	rows, err := query.Execute()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Escanear los resultados
	var users []models.UserStruct
	for rows.Next() {
		var user models.UserStruct
		err := orm.RowScanner(rows, &user)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}
