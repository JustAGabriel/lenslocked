package models

import (
	"fmt"
	"log"

	"github.com/justagabriel/lenslocked/util"
	"github.com/wneessen/go-mail"
)

const (
	DefaultSender = "support@lenslocked.com"
)

type Email struct {
	From, To, Subject, Plaintext, HTML string
}

type SMTPConfig struct {
	Host     string
	Port     int
	Username string
	Password string
}

type EmailService struct {
	DefaultSender string
	streamLogger  *util.StreamLogger
	dailer        *mail.Client
}

func NewEmailService(config SMTPConfig, defaultSenderEmail string) (*EmailService, error) {
	smtpClient, err := mail.NewClient(config.Host, mail.WithPort(config.Port), mail.WithSMTPAuth(mail.SMTPAuthPlain),
		mail.WithUsername(config.Username), mail.WithPassword(config.Password))
	if err != nil {
		log.Fatalf("failed to create mail client: %s", err)
		return nil, err
	}

	es := EmailService{
		DefaultSender: defaultSenderEmail,
		streamLogger:  &util.StreamLogger{Logger: log.Default()},
		dailer:        smtpClient,
	}

	return &es, nil
}

func (es *EmailService) Send(email Email) error {
	msg := mail.NewMsg()

	if err := msg.To(email.To); err != nil {
		log.Print(err)
	}

	if email.From == "" {
		email.From = es.DefaultSender
	}
	if err := msg.From(email.From); err != nil {
		log.Print(err)
	}

	switch {
	case email.Plaintext != "" && email.HTML != "":
		msg.SetBodyString(mail.TypeTextPlain, email.Plaintext)
		msg.SetBodyString(mail.TypeTextHTML, email.HTML)
	case email.Plaintext != "":
		msg.SetBodyString(mail.TypeTextPlain, email.Plaintext)
	case email.HTML != "":
		msg.SetBodyString(mail.TypeTextPlain, email.HTML)
	}

	msg.Subject(email.Subject)
	msg.WriteTo(es.streamLogger)

	if err := es.dailer.DialAndSend(msg); err != nil {
		log.Fatalf("failed to send mail: %s", err)
		return err
	}
	return nil
}

func (es *EmailService) SendForgotPasswordEmail(to, resetURL string) error {
	email := Email{
		To:        to,
		Subject:   "lenslock, password reset",
		Plaintext: "To reset your password, please click the following link: " + resetURL,
		HTML:      `<p>To reset your password, please visit the following <a href="` + resetURL + `"/>link</a>.`,
	}

	err := es.Send(email)
	if err != nil {
		return fmt.Errorf("error while sending 'forgot pw' eamil:\n%w", err)
	}

	return nil
}
