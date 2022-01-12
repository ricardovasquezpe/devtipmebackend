package config

import (
	"fmt"

	"devtipmebackend/api/models"

	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

type SendGridMailer struct {
	Email     string
	EmailName string
	ApiKey    string
	Url       string
	Api       string
}

func NewSendGridMailer(Email, EmailName, ApiKey, Url, Api string) SendGridMailer {
	m := SendGridMailer{Email, EmailName, ApiKey, Url, Api}
	return m
}

func (sm SendGridMailer) SendEmail(to []string, subject string, templateId string, items []models.TemplateData) error {
	m := mail.NewV3Mail()
	m.Subject = subject

	address := sm.Email
	name := sm.EmailName
	e := mail.NewEmail(name, address)
	m.SetFrom(e)

	m.SetTemplateID(templateId)

	p := mail.NewPersonalization()

	tos := []*mail.Email{}
	for _, element := range to {
		toElement := mail.NewEmail("User", element)
		tos = append(tos, toElement)
	}
	p.AddTos(tos...)

	for _, element := range items {
		p.SetDynamicTemplateData(element.Key, element.Value)
	}

	m.AddPersonalizations(p)

	fmt.Println(sm.Url)
	fmt.Println(sm.Api)
	request := sendgrid.GetRequest(sm.ApiKey, sm.Api, sm.Url)
	request.Method = "POST"
	var Body = mail.GetRequestBody(m)
	request.Body = Body
	response, err := sendgrid.API(request)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(response.StatusCode)
		fmt.Println(response.Body)
		fmt.Println(response.Headers)
	}

	return nil
}
