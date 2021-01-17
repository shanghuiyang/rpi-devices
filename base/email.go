package base

import (
	"fmt"
	"net/smtp"
)

// EmailConfig ...
type EmailConfig struct {
	SMTP     string `json:"smtp"`
	SMTPPort int    `json:"smtp_port"`
	POP      string `json:"pop"`
	POPPort  int    `json:"pop_port"`
	Address  string `json:"addr"`
	Password string `json:"password"`
}

// Email ...
type Email struct {
	SMTP     string
	SMTPPort int
	POP      string
	POPPort  int
	Address  string
	Password string
}

// EmailInfo ...
type EmailInfo struct {
	To      []string
	Subject string
	Body    string
}

// NewEmail ...
func NewEmail(cfg *EmailConfig) *Email {
	return &Email{
		SMTP:     cfg.SMTP,
		SMTPPort: cfg.SMTPPort,
		POP:      cfg.POP,
		POPPort:  cfg.POPPort,
		Address:  cfg.Address,
		Password: cfg.Password,
	}
}

// Send ...
func (e *Email) Send(info *EmailInfo) error {
	msg := fmt.Sprintf("From: %v\nTo: %v\nSubject: %v\n\n%v", e.Address, info.To, info.Subject, info.Body)
	err := smtp.SendMail(fmt.Sprintf("%v:%v", e.SMTP, e.SMTPPort),
		smtp.PlainAuth("", e.Address, e.Password, e.SMTP),
		e.Address, info.To, []byte(msg))
	return err
}
