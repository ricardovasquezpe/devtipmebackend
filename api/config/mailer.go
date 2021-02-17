package config

import (
	"bytes"
	"fmt"
	"log"
	"net/smtp"
	"os"
	"text/template"
)

type Mailer struct {
	Port     string
	Server   string
	Email    string
	Password string
	Auth     smtp.Auth
}

/*
type Request struct {
	from    string
	to      []string
	subject string
	body    string
}*/

func NewMailer(port, server, email, password string) Mailer {
	m := Mailer{port, server, email, password, nil}
	return m
}

func (m *Mailer) SetUpMailer() {
	m.Auth = smtp.PlainAuth("", m.Email, m.Password, m.Server)
}

func (m Mailer) SendEmail(to []string, subject string, templateName string, items interface{}) error {
	mime := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	bodyTemplate, err := parseTemplate(templateName, items)
	if err != nil {
		return err
	}
	body := "To: " + to[0] + "\r\nSubject: " + subject + "\r\n" + mime + "\r\n" + bodyTemplate
	SMTP := m.Server + ":" + m.Port
	if err := smtp.SendMail(SMTP, smtp.PlainAuth("", m.Email, m.Password, m.Server), m.Email, to, []byte(body)); err != nil {
		return err
	}
	return nil
}

func parseTemplate(fileName string, data interface{}) (string, error) {
	wd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(wd + fileName)
	t, err := template.ParseFiles(wd + "/api/" + fileName)
	if err != nil {
		return "", err
	}
	buffer := new(bytes.Buffer)
	if err = t.Execute(buffer, data); err != nil {
		return "", err
	}
	return buffer.String(), nil
}
