package helpers

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"semita/config"
	_ "strconv"

	"gopkg.in/gomail.v2"
)

type Notifier interface {
	Send(to, subject, body string) error
}

type EmailNotifier struct{}

func (emailNotifier EmailNotifier) Send(to, subject, body string) error {
	var mailConfig = config.MailConfig()

	var from = mailConfig.From.Address
	var port = mailConfig.SMTP.Port
	var host = mailConfig.SMTP.Host
	var username = mailConfig.SMTP.Username
	var password = mailConfig.SMTP.Password

	var goMail = gomail.NewMessage()

	// Set email headers
	goMail.SetHeader("From", from)
	goMail.SetHeader("To", to)
	goMail.SetHeader("Subject", subject)

	// Set email body
	goMail.SetBody("text/html", body)

	// Set up SMTP server configuration
	var dialer = gomail.NewDialer(host, port, username, password)

	var sendErr = dialer.DialAndSend(goMail)
	if sendErr != nil {
		return fmt.Errorf("error sending email: %w", sendErr)
	}

	return nil
}

// GenerateResetToken genera un hash simple para recuperación
func GenerateResetToken(email string) string {
	h := sha256.New()
	h.Write([]byte(email + ":reset"))
	return hex.EncodeToString(h.Sum(nil))
}
