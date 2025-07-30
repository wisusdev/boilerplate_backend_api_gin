package examples

import (
	"database/sql"
	"semita/app/data/models"
	"semita/core/database/orm"
	"time"
)

// UserRepository es un ejemplo de repositorio específico para usuarios
type UserRepository struct {
	*orm.ModelRepository
}

// NewUserRepository crea una nueva instancia del repositorio de usuarios
func NewUserRepository() *UserRepository {
	return &UserRepository{
		ModelRepository: orm.NewModelRepository(&models.UserStruct{}),
	}
}

// GetAllUsers obtiene todos los usuarios
func (r *UserRepository) GetAllUsers() (models.Users, error) {
	var users models.Users
	err := r.GetAll(&users)
	return users, err
}

// GetUserByID obtiene un usuario por su ID
func (r *UserRepository) GetUserByID(id interface{}) (models.UserStruct, error) {
	var user models.UserStruct
	err := r.FindById(id, &user)
	return user, err
}

// GetUserByEmail obtiene un usuario por su email
func (r *UserRepository) GetUserByEmail(email string) (models.UserStruct, error) {
	var user models.UserStruct
	err := r.First("email = ?", []interface{}{email}, &user)
	return user, err
}

// StoreUser crea un nuevo usuario
func (r *UserRepository) StoreUser(user models.UserStruct) (models.UserStruct, error) {
	// Establecer valores predeterminados
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	// Guardar el usuario
	err := r.Save(&user)
	if err != nil {
		return models.UserStruct{}, err
	}

	return user, nil
}

// UpdateUser actualiza un usuario existente
func (r *UserRepository) UpdateUser(user models.UserStruct) error {
	user.UpdatedAt = time.Now()
	return r.Save(&user)
}

// DeleteUser elimina un usuario por su ID
func (r *UserRepository) DeleteUser(id interface{}) (sql.Result, error) {
	return r.Delete(id)
}

// MarkEmailVerified marca el email de un usuario como verificado
func (r *UserRepository) MarkEmailVerified(userID int) error {
	user, err := r.GetUserByID(userID)
	if err != nil {
		return err
	}

	now := time.Now()
	user.EmailVerifiedAt = &now
	return r.UpdateUser(user)
}

// FindActiveUsers encuentra usuarios activos (con email verificado)
func (r *UserRepository) FindActiveUsers() (models.Users, error) {
	var users models.Users
	err := r.Where("email_verified_at IS NOT NULL", []interface{}{}, &users)
	return users, err
}

// TableName implementa la interfaz Model para UserStruct
/*func (u models.UserStruct) TableName() string {
	return "users"
}*/
