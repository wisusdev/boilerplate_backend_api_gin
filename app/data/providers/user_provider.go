package providers

import (
	"boilerplate_backend_api_gin/app/data/models"
	"boilerplate_backend_api_gin/app/data/repositories"
)

var tableName = "users"

// GetAllUsers obtiene todos los usuarios a través del repositorio
func GetAllUsers() ([]models.UserStruct, error) {
	repository := repositories.NewUserRepository()
	return repository.GetAllUsers()
}

// StoreUser guarda un nuevo usuario a través del repositorio
func StoreUser(storeUser models.UserStruct) (user models.UserStruct, err error) {
	return repositories.StoreUser(storeUser)
}

func UpdateUser(user models.UserStruct) (err error) {
	return repositories.UpdateUser(user)
}

// GetUserByID obtiene un usuario por ID a través del repositorio
func GetUserByID(id string) (models.UserStruct, error) {
	return repositories.GetUserByID(id)
}

// GetUserByEmail obtiene un usuario por email a través del repositorio
func GetUserByEmail(email string) (models.UserStruct, error) {
	return repositories.GetUserByEmail(email)
}

// UpdateUserPassword actualiza un usuario a través del repositorio
func UpdateUserPassword(user models.UserStruct) error {
	return repositories.UpdateUserPassword(user)
}

// DeleteUser elimina un usuario a través del repositorio
func DeleteUser(id string) error {
	return repositories.DeleteUser(id)
}

// MarkEmailVerified marca el email como verificado a través del repositorio
func MarkEmailVerified(userID string) error {
	return repositories.MarkEmailVerified(userID)
}
