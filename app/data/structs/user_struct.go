package structs

import (
	"time"
)

type UserBase struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

type UserStruct struct {
	ID        int       `json:"id" db:"INT PRIMARY KEY AUTO_INCREMENT"`
	FirstName string    `json:"first_name" db:"VARCHAR(255)" nullable:"false"`
	LastName  string    `json:"last_name" db:"VARCHAR(255)" nullable:"false"`
	Username  string    `json:"username" db:"VARCHAR(255)" unique:"true" nullable:"false"`
	Avatar    string    `json:"avatar" db:"VARCHAR(255)" nullable:"true" default:"NULL"`
	Language  string    `json:"language" db:"VARCHAR(10)" default:"'en'" nullable:"true"`
	Email     string    `json:"email" db:"VARCHAR(255)" unique:"true" nullable:"false"`
	Password  string    `json:"password" db:"VARCHAR(255)" nullable:"false"`
	CreatedAt time.Time `json:"created_at" db:"DATETIME" default:"CURRENT_TIMESTAMP"`
	UpdatedAt time.Time `json:"updated_at" db:"DATETIME" default:"CURRENT_TIMESTAMP"`
}

type Users []UserStruct

type RegisterUserStruct struct {
	Name            string `json:"name"`
	Email           string `json:"email"`
	Password        string `json:"password"`
	ConfirmPassword string `json:"confirm_password"`
}

type StoreUserStruct struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	Password  string `json:"password"`
}

type UpdateUserStruct struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginUserStruct struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
