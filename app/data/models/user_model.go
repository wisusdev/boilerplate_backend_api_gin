package models

import (
	"semita/core/common/nulltypes"
	"time"
)

type UserStruct struct {
	ID              string               `json:"id" db:"VARCHAR(36) PRIMARY KEY"`
	FirstName       string               `json:"first_name" db:"VARCHAR(255)" nullable:"false"`
	LastName        string               `json:"last_name" db:"VARCHAR(255)" nullable:"false"`
	Username        string               `json:"username" db:"VARCHAR(255)" unique:"true" nullable:"false"`
	Avatar          nulltypes.NullString `json:"avatar" db:"VARCHAR(255)" nullable:"true" default:"NULL"`
	Language        string               `json:"language" db:"VARCHAR(10)" default:"'en'" nullable:"true"`
	Email           string               `json:"email" db:"VARCHAR(255)" unique:"true" nullable:"false"`
	EmailVerifiedAt *time.Time           `json:"email_verified_at" db:"DATETIME" default:"NULL"`
	RememberToken   nulltypes.NullString `json:"remember_token" db:"VARCHAR(100)" nullable:"true"`
	Password        string               `json:"password" db:"VARCHAR(255)" nullable:"false"`
	CreatedAt       time.Time            `json:"created_at" db:"DATETIME" default:"CURRENT_TIMESTAMP"`
	UpdatedAt       time.Time            `json:"updated_at" db:"DATETIME" default:"CURRENT_TIMESTAMP"`
}

type Users []UserStruct

func (user UserStruct) TableName() string {
	return "users"
}
