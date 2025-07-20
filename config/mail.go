package config

type Mailers struct {
	SMTP     Smtp
	Sendmail Sendmail
	From     From
}

type Smtp struct {
	Transport  string `json:"transport"` // "smtp" or "sendmail"
	Host       string `json:"host"`
	Port       int    `json:"port"`
	Username   string `json:"username"`
	Password   string `json:"password"`
	Encryption string `json:"encryption"` // "tls" or "ssl"
}

type Sendmail struct {
	Transport string `json:"transport"` // "sendmail"
	Path      string `json:"path"`      // Path to the sendmail binary
}

type From struct {
	Address string `json:"address"` // Email address to use as the sender
	Name    string `json:"name"`    // Name to use as the sender
}

func (m *Mailers) GetSMTP() Smtp {
	return m.SMTP
}

func (m *Mailers) GetSendmail() Sendmail {
	return m.Sendmail
}

func (m *Mailers) GetFrom() From {
	return m.From
}

func MailConfig() *Mailers {
	return &Mailers{
		SMTP: Smtp{
			Transport:  GetEnv("MAIL_MAILER", ""),
			Host:       GetEnv("MAIL_HOST", ""),
			Port:       GetEnvInt("MAIL_PORT", 587),
			Username:   GetEnv("MAIL_USERNAME", ""),
			Password:   GetEnv("MAIL_PASSWORD", ""),
			Encryption: GetEnv("MAIL_ENCRYPTION", "tls"),
		},
		Sendmail: Sendmail{
			Transport: GetEnv("MAIL_MAILER", ""),
			Path:      GetEnv("MAIL_SENDMAIL_PATH", "/usr/sbin/sendmail"),
		},
		From: From{
			Address: GetEnv("MAIL_FROM_ADDRESS", ""),
			Name:    GetEnv("MAIL_FROM_NAME", "Semita"),
		},
	}
}
