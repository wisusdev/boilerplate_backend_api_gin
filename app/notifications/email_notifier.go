package notifications

import (
	"boilerplate_backend_api_gin/core/helpers"
	"fmt"
)

var DefaultNotifier helpers.Notifier = helpers.EmailNotifier{}

func SendEmailVerification(to, url string) error {
	subject := "Verifica tu correo electrónico"
	body := fmt.Sprintf("<p>Haz clic en el siguiente enlace para verificar tu email:</p><p><a href=\"%s\">Verificar Email</a></p>", url)
	return DefaultNotifier.Send(to, subject, body)
}

func SendPasswordReset(to, url string) error {
	subject := "Restablece tu contraseña"
	body := fmt.Sprintf("<p>Haz clic en el siguiente enlace para restablecer tu contraseña:</p><p><a href=\"%s\">Restablecer contraseña</a></p>", url)
	return DefaultNotifier.Send(to, subject, body)
}
