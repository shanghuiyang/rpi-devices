package base

import (
	"fmt"
	"log"
	"net/smtp"
)

const (
	logTagEmail = "email"
)

var (
	chEmail   = make(chan *EmailInfo, 4)
	emailList []string
)

// Init ...
func Init(cfg *Config) {
	emailList = cfg.EmailTo.List
	e := NewEmail(cfg.Email)
	go e.Start()
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

// Start ...
func (e *Email) Start() {
	log.Printf("[%v]start working", logTagEmail)
	for info := range chEmail {
		if err := e.Send(info); err != nil {
			log.Printf("[%v]faied to send email, error: %v", logTagEmail, err)
		}
	}
	return
}

// Send ...
func (e *Email) Send(info *EmailInfo) error {
	msg := fmt.Sprintf("From: %v\nTo: %v\nSubject: %v\n\n%v", e.Address, info.To, info.Subject, info.Body)
	err := smtp.SendMail(fmt.Sprintf("%v:%v", e.SMTP, e.SMTPPort),
		smtp.PlainAuth("", e.Address, e.Password, e.SMTP),
		e.Address, info.To, []byte(msg))
	return err
}
