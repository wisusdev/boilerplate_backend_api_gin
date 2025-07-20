package requests

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

// RegisterRequest valida los datos de registro
type RegisterRequest struct {
	Name                 string `form:"name" json:"name" binding:"required,min=2"`
	Email                string `form:"email" json:"email" binding:"required,email"`
	Password             string `form:"password" json:"password" binding:"required,min=6"`
	PasswordConfirmation string `form:"password_confirmation" json:"password_confirmation" binding:"required,eqfield=Password"`
}

func (r *RegisterRequest) Validate(c *gin.Context) error {
	if err := c.ShouldBind(r); err != nil {
		return err
	}
	return validate.Struct(r)
}

// ForgotPasswordRequest valida el email para recuperar contraseña
type ForgotPasswordRequest struct {
	Email string `form:"email" json:"email" binding:"required,email"`
}

func (r *ForgotPasswordRequest) Validate(c *gin.Context) error {
	if err := c.ShouldBind(r); err != nil {
		return err
	}
	return validate.Struct(r)
}

// ResetPasswordRequest valida el reseteo de contraseña
type ResetPasswordRequest struct {
	Token                string `form:"token" json:"token" binding:"required"`
	Email                string `form:"email" json:"email" binding:"required,email"`
	Password             string `form:"password" json:"password" binding:"required,min=6"`
	PasswordConfirmation string `form:"password_confirmation" json:"password_confirmation" binding:"required,eqfield=Password"`
}

func (r *ResetPasswordRequest) Validate(c *gin.Context) error {
	if err := c.ShouldBind(r); err != nil {
		return err
	}
	return validate.Struct(r)
}
