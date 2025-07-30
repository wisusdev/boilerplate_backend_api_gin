# Simple ORM para Go

Este es un ORM simple para Go diseñado para facilitar las operaciones CRUD básicas mientras mantiene la flexibilidad para consultas más complejas.

## Características

- Operaciones CRUD básicas con modelos
- Constructor de consultas fluido
- Soporte para transacciones
- Sistema de escaneo automático de filas
- Generación de SQL para crear tablas

## Uso Básico

### Definir un Modelo

```go
package models

import (
	"time"
)

type User struct {
	ID              int        `json:"id" db:"INT PRIMARY KEY AUTO_INCREMENT"`
	FirstName       string     `json:"first_name" db:"VARCHAR(255)" nullable:"false"`
	LastName        string     `json:"last_name" db:"VARCHAR(255)" nullable:"false"`
	Email           string     `json:"email" db:"VARCHAR(255)" unique:"true" nullable:"false"`
	Password        string     `json:"password" db:"VARCHAR(255)" nullable:"false"`
	CreatedAt       time.Time  `json:"created_at" db:"DATETIME" default:"CURRENT_TIMESTAMP"`
	UpdatedAt       time.Time  `json:"updated_at" db:"DATETIME" default:"CURRENT_TIMESTAMP"`
}

// TableName implementa la interfaz Model
func (u User) TableName() string {
	return "users"
}

type Users []User
```

### Operaciones Básicas

```go
// Crear un repositorio para el modelo User
repo := orm.NewModelRepository(&models.User{})

// Obtener todos los usuarios
var users models.Users
err := repo.GetAll(&users)

// Encontrar usuario por ID
var user models.User
err = repo.FindById(1, &user)

// Encontrar usuarios por un campo
var adminUsers models.Users
err = repo.FindBy("role", "admin", &adminUsers)

// Crear un nuevo usuario
newUser := models.User{
	FirstName: "John",
	LastName:  "Doe",
	Email:     "john@example.com",
	Password:  "hashedpassword",
}
err = repo.Save(&newUser)

// Actualizar un usuario existente
user.FirstName = "Jane"
err = repo.Save(&user)

// Eliminar un usuario
err = repo.Delete(user.ID)
```

### Constructor de Consultas

```go
// Crear un repositorio para el modelo User
repo := orm.NewModelRepository(&models.User{})

// Construir una consulta personalizada
qb := repo.Query().Select("id", "first_name", "last_name", "email")
qb.Where("created_at > ?", "2023-01-01")
qb.OrderBy("last_name ASC")
qb.Limit(10).Offset(20)

// Ejecutar la consulta
rows, err := qb.Execute()
if err != nil {
	// manejar error
}
defer rows.Close()

// Procesar las filas
for rows.Next() {
	var user models.User
	err := orm.RowScanner(rows, &user)
	if err != nil {
		// manejar error
	}
	// usar user...
}
```

### Transacciones

```go
err := orm.WithTransaction(func(tx *orm.Transaction) error {
	// Crear un repositorio que use la transacción
	repo := orm.NewRepositoryWithDB("users", tx.GetAdapter())

	// Realizar operaciones dentro de la transacción
	result, err := repo.Exec("UPDATE users SET active = ? WHERE id = ?", true, 1)
	if err != nil {
		return err // Esto provocará un rollback
	}

	// Más operaciones...
	return nil // Retornar nil hará commit de la transacción
})
```

## Extensibilidad

El ORM está diseñado para ser extensible. Puedes crear repositorios específicos para cada modelo con métodos personalizados:

```go
type UserRepository struct {
	*orm.ModelRepository
}

func NewUserRepository() *UserRepository {
	return &UserRepository{
		ModelRepository: orm.NewModelRepository(&models.User{}),
	}
}

// Método personalizado para el repositorio de usuarios
func (r *UserRepository) FindByEmail(email string) (models.User, error) {
	var user models.User
	err := r.First("email = ?", []interface{}{email}, &user)
	return user, err
}
```

## Generar Esquema de Base de Datos

```go
// Generar SQL para crear la tabla de usuarios
userSQL := orm.GenerateCreateTableSQL(&models.User{})
fmt.Println(userSQL)
```
