package helper

import (
	"bytes"
	"errors"
	"net/smtp"
	"os"
	"text/template"

	"gopkg.in/gomail.v2"
)

const CONFIG_SMTP_HOST = "smtp.gmail.com"
const CONFIG_SMTP_PORT = 587
const CONFIG_SENDER_NAME = "Legible ID <id.legible@gmail.com>"

func SendMail(to string, subject string, message interface{}) error {
	mailer := gomail.NewMessage()
	mailer.SetHeader("From", CONFIG_SENDER_NAME)
	mailer.SetHeader("To", to)
	mailer.SetHeader("Subject", subject)
	mailer.SetBody("text/html", message.(string))

	dialer := gomail.NewDialer(
		CONFIG_SMTP_HOST,
		CONFIG_SMTP_PORT,
		os.Getenv("SMTP_EMAIL"),
		os.Getenv("SMTP_PASSWORD"),
	)

	err := dialer.DialAndSend(mailer)

	if err != nil {
		return errors.New("There was an error with mailer service")
	}

	return nil
}

//Request struct
type Request struct {
	from    string
	to      []string
	subject string
	body    string
}

func NewRequest(to []string, subject, body string) *Request {
	return &Request{
		to:      to,
		subject: subject,
		body:    body,
	}
}

func (r *Request) SendEmail() (bool, error) {
	mime := "MIME-version: 1.0;\nContent-Type: text/plain; charset=\"UTF-8\";\n\n"
	subject := "Subject: " + r.subject + "!\n"
	msg := []byte(subject + mime + "\n" + r.body)
	addr := "smtp.gmail.com:587"
	auth := smtp.PlainAuth("", os.Getenv("SMTP_EMAIL"), os.Getenv("SMTP_EMAIL"), "smtp.gmail.com")

	if err := smtp.SendMail(addr, auth, CONFIG_SENDER_NAME, r.to, msg); err != nil {
		return false, err
	}
	return true, nil
}

func (r *Request) ParseTemplate(templateFileName string, data interface{}) error {
	t, err := template.ParseFiles(templateFileName)
	if err != nil {
		return err
	}
	buf := new(bytes.Buffer)
	if err = t.Execute(buf, data); err != nil {
		return err
	}
	r.body = buf.String()
	return nil
}
