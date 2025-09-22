package examples

import (
	"boilerplate_backend_api_gin/app/data/models"
	"boilerplate_backend_api_gin/core/database/orm"
	"time"
)

// TransferExample muestra cómo usar transacciones para operaciones que involucran múltiples modelos
func TransferExample() error {
	// Usar una transacción para operaciones atómicas
	return orm.WithTransaction(func(tx *orm.Transaction) error {
		// Crear repositorios usando la misma transacción
		userRepo := orm.NewRepositoryWithDB("users", tx.GetAdapter())

		// Ejemplo: actualizar múltiples usuarios en una transacción
		var user1 models.UserStruct
		var user2 models.UserStruct

		// Obtener usuario 1
		row1, _ := userRepo.Query().Where("id = ?", 1).First()
		err := orm.RowScanner(row1, &user1)
		if err != nil {
			return err // Esto provocará un rollback
		}

		// Obtener usuario 2
		row2, _ := userRepo.Query().Where("id = ?", 2).First()
		err = orm.RowScanner(row2, &user2)
		if err != nil {
			return err // Esto provocará un rollback
		}

		// Actualizar usuario 1
		user1.FirstName = "Nuevo Nombre"
		user1.UpdatedAt = time.Now()
		_, err = userRepo.Update(user1.ID, user1)
		if err != nil {
			return err // Esto provocará un rollback
		}

		// Actualizar usuario 2
		user2.FirstName = "Otro Nombre"
		user2.UpdatedAt = time.Now()
		_, err = userRepo.Update(user2.ID, user2)
		if err != nil {
			return err // Esto provocará un rollback
		}

		// Si llegamos aquí sin errores, la transacción se confirma automáticamente
		return nil
	})
}

// UserRegistrationExample demuestra cómo usar transacciones para registrar un usuario y crear recursos relacionados
func UserRegistrationExample(newUser models.UserStruct) error {
	return orm.WithTransaction(func(tx *orm.Transaction) error {
		// Crear repositorio de usuarios con la transacción
		userRepo := orm.NewRepositoryWithDB("users", tx.GetAdapter())

		// Crear el usuario
		result, err := userRepo.Create(newUser)
		if err != nil {
			return err
		}

		// Obtener el ID del usuario recién creado
		_, err = result.LastInsertId()
		if err != nil {
			return err
		}

		// Aquí podrías crear otros registros relacionados con el usuario
		// Por ejemplo, un perfil, configuraciones, etc.
		// Todos estos cambios se realizarían en la misma transacción

		return nil // La transacción se confirma
	})
}
