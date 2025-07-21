package seeders

import (
	"fmt"
	"log"
	"semita/core/database/database_connections"
	"semita/core/database/generate_seeders"
	"semita/core/helpers"

	"golang.org/x/crypto/bcrypt"
)

// UsersSeeder seeder para usuarios de prueba
type UsersSeeder struct {
	generate_seeders.BaseSeeder
}

// NewUsersSeeder crea una nueva instancia del seeder
func NewUsersSeeder() *UsersSeeder {
	return &UsersSeeder{
		BaseSeeder: generate_seeders.BaseSeeder{
			DB:   database_connections.DatabaseConnectSQL(),
			Name: "users_seeder",
		},
	}
}

func (us *UsersSeeder) GetName() string {
	return us.BaseSeeder.Name
}

// GetDependencies retorna las dependencias del seeder
func (us *UsersSeeder) GetDependencies() []string {
	return []string{"roles_permissions_seeder"} // Depende de roles y permisos
}

// GetTables retorna las tablas que maneja este seeder (en orden de limpieza)
func (us *UsersSeeder) GetTables() []string {
	return []string{"user_permissions", "user_roles", "users"}
}

// Seed ejecuta el seeding de usuarios
func (us *UsersSeeder) Seed() error {
	log.Println("Seeding users...")

	// 12345678aA
	var passwordHash, _ = bcrypt.GenerateFromPassword([]byte("12345678aA"), bcrypt.DefaultCost)

	users := []struct {
		FirstName string
		LastName  string
		Username  string
		Email     string
		Password  string
		Role      string
	}{
		{
			FirstName: "Super",
			LastName:  "Admin",
			Username:  "superadmin",
			Email:     "superadmin@example.com",
			Password:  string(passwordHash),
			Role:      "super-admin",
		},
		{
			FirstName: "Admin",
			LastName:  "User",
			Username:  "admin",
			Email:     "admin@example.com",
			Password:  string(passwordHash),
			Role:      "admin",
		},
		{
			FirstName: "Editor",
			LastName:  "Principal",
			Username:  "editor",
			Email:     "editor@example.com",
			Password:  string(passwordHash),
			Role:      "editor",
		},
		{
			FirstName: "Moderador",
			LastName:  "00",
			Username:  "moderator",
			Email:     "moderator@example.com",
			Password:  string(passwordHash),
			Role:      "moderator",
		},
		{
			FirstName: "Usuario",
			LastName:  "General",
			Username:  "user",
			Email:     "user@example.com",
			Password:  string(passwordHash),
			Role:      "user",
		},
		{
			FirstName: "Maria",
			LastName:  "Garcia",
			Username:  "maria.garcia",
			Email:     "maria.garcia@example.com",
			Password:  string(passwordHash),
			Role:      "user",
		},
		{
			FirstName: "Carlos",
			LastName:  "Lopez",
			Username:  "carlos.lopez",
			Email:     "carlos.lopez@example.com",
			Password:  string(passwordHash),
			Role:      "user",
		},
		{
			FirstName: "Ana",
			LastName:  "Martinez",
			Username:  "ana.martinez",
			Email:     "ana.martinez@example.com",
			Password:  string(passwordHash),
			Role:      "editor",
		},
		{
			FirstName: "Jesús",
			LastName:  "Avelar",
			Username:  "user00",
			Email:     "user00@wisus.dev",
			Password:  string(passwordHash),
			Role:      "user",
		},
	}

	for _, user := range users {
		// Crear nuevo usuario directamente (ya se limpiaron los datos)
		insertQuery := `
			INSERT INTO users (first_name, last_name, username, email, password, email_verified_at, created_at, updated_at)
			VALUES (?, ?, ?, ?, ?, NULL, NOW(), NOW())
			`

		result, err := us.BaseSeeder.DB.Exec(insertQuery, user.FirstName, user.LastName, user.Username, user.Email, user.Password)
		if err != nil {
			helpers.Logs("ERROR", fmt.Sprintf("Error creating user '%s': %v", user.Email, err))
			continue
		}

		userID, _ := result.LastInsertId()
		log.Printf("Created user: %s (ID: %d)", user.Email, userID)

		// Asignar rol al usuario
		err = us.assignRoleToUser(int(userID), user.Role)
		if err != nil {
			log.Printf("Error assigning role '%s' to user '%s': %v", user.Role, user.Email, err)
		} else {
			log.Printf("Assigned role '%s' to user '%s'", user.Role, user.Email)
		}
	}

	log.Println("Users seeding completed successfully!")
	return nil
}

// assignRoleToUser asigna un rol a un usuario
func (us *UsersSeeder) assignRoleToUser(userID int, roleName string) error {
	// Obtener el ID del rol
	var roleID int
	roleQuery := `SELECT id FROM roles WHERE name = ? AND guard_name = 'web'`
	err := us.BaseSeeder.DB.QueryRow(roleQuery, roleName).Scan(&roleID)
	if err != nil {
		return err
	}

	// Crear la relación directamente (ya se limpiaron los datos)
	insertQuery := `INSERT INTO user_roles (user_id, role_id) VALUES (?, ?)`
	_, err = us.BaseSeeder.DB.Exec(insertQuery, userID, roleID)
	return err
}
