package email

import (
	"fmt"
	"net/smtp"
)

type SMTPConfig struct {
	Host     string
	Port     int
	Username string
	Password string
	From     string
}

type EmailService struct {
	config *SMTPConfig
}

func NewEmailService(config *SMTPConfig) *EmailService {
	return &EmailService{
		config: config,
	}
}

func (s *EmailService) Send(to, subject, message string) error {
	auth := smtp.PlainAuth("", s.config.Username, s.config.Password, s.config.Host)
	msg := []byte(
		"From: " + s.config.From + "\r\n" +
			"To: " + to + "\r\n" +
			"Subject: " + subject + "\r\n\r\n" +
			message + "\r\n",
	)
	addr := fmt.Sprintf("%s:%d", s.config.Host, s.config.Port)
	return smtp.SendMail(
		addr,
		auth,
		s.config.From,
		[]string{to},
		msg,
	)
}
