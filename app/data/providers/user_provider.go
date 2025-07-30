package providers

import (
	"semita/app/data/models"
	"semita/app/data/repositories"
)

var tableName = "users"

// GetAllUsers obtiene todos los usuarios a través del repositorio
func GetAllUsers() ([]models.UserStruct, error) {
	repository := repositories.NewUserRepository()

	var users []models.UserStruct
	err := repository.GetAll(&users)
	if err != nil {
		return nil, err
	}
	return users, nil
}

// StoreUser guarda un nuevo usuario a través del repositorio
func StoreUser(storeUser models.UserStruct) (user models.UserStruct, err error) {
	return repositories.StoreUser(storeUser)
}

// GetUserByID obtiene un usuario por ID a través del repositorio
func GetUserByID(id string) (models.UserStruct, error) {
	return repositories.GetUserByID(id)
}

// GetUserByEmail obtiene un usuario por email a través del repositorio
func GetUserByEmail(email string) (models.UserStruct, error) {
	return repositories.GetUserByEmail(email)
}

// UpdateUser actualiza un usuario a través del repositorio
func UpdateUser(user models.UserStruct) error {
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
