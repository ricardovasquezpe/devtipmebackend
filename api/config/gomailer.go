package config

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"strconv"
	"text/template"

	mail "github.com/xhit/go-simple-mail/v2"
)

type GoMailer struct {
	Port     string
	Server   string
	Email    string
	Password string
	Client   *mail.SMTPClient
}

func NewGoMailer(port, server, email, password string) GoMailer {
	m := GoMailer{port, server, email, password, nil}
	return m
}

func (m *GoMailer) SetUpGoMailer() {
	fmt.Println(m.Server)
	fmt.Println(strconv.Atoi(m.Port))
	fmt.Println(m.Email)
	fmt.Println(m.Password)

	server := mail.NewSMTPClient()
	server.Host = m.Server
	server.Port, _ = strconv.Atoi(m.Port)
	server.Username = m.Email
	server.Password = m.Password
	server.Encryption = mail.EncryptionTLS
	/*server.KeepAlive = false
	server.ConnectTimeout = 10 * time.Second
	server.SendTimeout = 10 * time.Second
	server.TLSConfig = &tls.Config{InsecureSkipVerify: true}*/

	smtpClient, err := server.Connect()
	if err != nil {
		log.Fatal(err)
	}

	m.Client = smtpClient
}

func (m GoMailer) SendEmail(to []string, subject string, templateName string, items interface{}) error {
	email := mail.NewMSG()
	email.SetFrom("From Me <devtipmedeveloper@gmail.com>")
	email.AddTo("devtipmedeveloper@gmail.com")
	//email.AddCc("another_you@example.com")
	email.SetSubject(subject)

	bodyTemplate, err := parseTemplate(templateName, items)
	if err != nil {
		return err
	}

	email.SetBody(mail.TextHTML, bodyTemplate)

	err = email.Send(m.Client)
	if err != nil {
		return err
	}

	return nil
}

func parseTemplateGo(fileName string, data interface{}) (string, error) {
	wd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
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
