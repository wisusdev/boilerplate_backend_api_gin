package seeders

import (
	"fmt"
	"log"
	"semita/core/database/database_connections"
	"semita/core/database/generate_seeders"
	"semita/core/helpers"

	"golang.org/x/crypto/bcrypt"
)

// UserSeeder seeder para usuarios
type UserSeeder struct {
	generate_seeders.BaseSeeder
}

// NewUserSeeder crea una nueva instancia del seeder
func NewUserSeeder() *UserSeeder {
	return &UserSeeder{
		BaseSeeder: generate_seeders.BaseSeeder{
			DB:   database_connections.DatabaseConnectSQL(),
			Name: "user_seeder",
		},
	}
}

func (us *UserSeeder) GetName() string {
	return us.BaseSeeder.Name
}

// GetDependencies retorna las dependencias del seeder
func (us *UserSeeder) GetDependencies() []string {
	return []string{"role_seeder"} // Depende del RoleSeeder
}

// GetTables retorna las tablas que maneja este seeder
func (us *UserSeeder) GetTables() []string {
	return []string{"model_has_roles", "users"}
}

// Seed ejecuta el seeding de usuarios exactamente como en Laravel
func (us *UserSeeder) Seed() error {
	// Crear 10 usuarios de prueba (equivalente a User::factory(10)->create())
	for i := 1; i <= 10; i++ {
		passwordHash, _ := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.DefaultCost)

		query := `INSERT INTO users (id, first_name, last_name, username, email, password, created_at, updated_at) 
				  VALUES (UUID(), ?, ?, ?, ?, ?, NOW(), NOW())`

		_, err := us.BaseSeeder.DB.Exec(query,
			fmt.Sprintf("User%d", i),
			fmt.Sprintf("Test%d", i),
			fmt.Sprintf("user%d", i),
			fmt.Sprintf("user%d@example.com", i),
			string(passwordHash))

		if err != nil {
			log.Printf("Error creating factory user %d: %v", i, err)
			continue
		}
	}

	// Crear usuario específico como en Laravel
	passwordHash, _ := bcrypt.GenerateFromPassword([]byte("12345678aA"), bcrypt.DefaultCost)

	userQuery := `INSERT INTO users (id, first_name, last_name, username, email, password, created_at, updated_at) 
				  VALUES (UUID(), ?, ?, ?, ?, ?, NOW(), NOW())`

	_, err := us.BaseSeeder.DB.Exec(userQuery, "Jesus", "Avelar", "user00", "user00@wisus.dev", string(passwordHash))
	if err != nil {
		helpers.Logs("ERROR", fmt.Sprintf("Error creating specific user: %v", err))
		return err
	}

	// Obtener el UUID del usuario recién creado
	var userUUID string
	err = us.BaseSeeder.DB.QueryRow("SELECT id FROM users WHERE email = ?", "user00@wisus.dev").Scan(&userUUID)
	if err != nil {
		log.Printf("Error getting user UUID: %v", err)
		return err
	}

	// Asignar rol admin al usuario específico
	err = us.assignRoleToUser(userUUID, "admin")
	if err != nil {
		log.Printf("Error assigning admin role to user: %v", err)
		return err
	}

	log.Println("UserSeeder completed successfully!")
	return nil
}

// assignRoleToUser asigna un rol a un usuario
func (us *UserSeeder) assignRoleToUser(userUUID string, roleName string) error {
	// Obtener el UUID del rol
	var roleUUID string
	err := us.BaseSeeder.DB.QueryRow("SELECT uuid FROM roles WHERE name = ? AND guard_name = ?", roleName, "api").Scan(&roleUUID)
	if err != nil {
		return err
	}

	// Asignar rol al usuario en la tabla model_has_roles
	query := `INSERT INTO model_has_roles (role_id, model_type, model_id) VALUES (?, ?, ?)`
	_, err = us.BaseSeeder.DB.Exec(query, roleUUID, "users", userUUID)
	return err
}
