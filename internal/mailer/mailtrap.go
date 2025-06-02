package mailer

import (
	"bytes"
	"errors"
	"fmt"
	"text/template"

	gomail "gopkg.in/mail.v2"
)

type mailtrapClient struct {
	fromEmail string
	apiKey    string
}

func NewMailTrapClient(apiKey, fromEmail string) (mailtrapClient, error) {
	if apiKey == "" {
		return mailtrapClient{}, errors.New("api key is required")
	}

	return mailtrapClient{
		fromEmail: fromEmail,
		apiKey:    apiKey,
	}, nil
}

func (m mailtrapClient) Send(templateFile, username, email string, data any, isSandbox bool) (int, error) {
	fmt.Printf("---mailtrap.go")
	// Template parsing and building
	tmpl, err := template.ParseFS(FS, "templates/"+templateFile)
	if err != nil {
		return -1, err
	}

	subject := new(bytes.Buffer)
	err = tmpl.ExecuteTemplate(subject, "subject", data)
	if err != nil {
		return -1, err
	}

	body := new(bytes.Buffer)
	err = tmpl.ExecuteTemplate(body, "body", data)
	if err != nil {
		return -1, err
	}

	message := gomail.NewMessage()
	message.SetHeader("From", m.fromEmail)
	message.SetHeader("To", email)
	message.SetHeader("Subject", subject.String())

	message.AddAlternative("text/html", body.String())
	fmt.Printf("---mailtrap.go--inside before dialer-")
	dialer := gomail.NewDialer("sandbox.smtp.mailtrap.io", 587, "ece2d425d228c0", m.apiKey)
	fmt.Printf("---mailtrap.go--inside after dialer-")
	if err := dialer.DialAndSend(message); err != nil {
		fmt.Print("---mailtrap.go--errorDialer--err:", err)
		return -1, err
	}

	return 200, nil
}
