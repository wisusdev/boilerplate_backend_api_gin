package models

import (
	"semita/app/data/repositories"
	"semita/app/data/structs"
)

var tableName = "users"

// GetAllUsers obtiene todos los usuarios a través del repositorio
func GetAllUsers() ([]structs.UserStruct, error) {
	return repositories.GetAllUsers()
}

// StoreUser guarda un nuevo usuario a través del repositorio
func StoreUser(user structs.StoreUserStruct) error {
	return repositories.StoreUser(user)
}

// GetUserByID obtiene un usuario por ID a través del repositorio
func GetUserByID(id string) (structs.UserStruct, error) {
	return repositories.GetUserByID(id)
}

// GetUserByEmail obtiene un usuario por email a través del repositorio
func GetUserByEmail(email string) (structs.UserStruct, error) {
	return repositories.GetUserByEmail(email)
}

// UpdateUser actualiza un usuario a través del repositorio
func UpdateUser(user structs.UpdateUserStruct) error {
	return repositories.UpdateUser(user)
}

// DeleteUser elimina un usuario a través del repositorio
func DeleteUser(id string) error {
	return repositories.DeleteUser(id)
}

// MarkEmailVerified marca el email como verificado a través del repositorio
func MarkEmailVerified(userID int) error {
	return repositories.MarkEmailVerified(userID)
}
