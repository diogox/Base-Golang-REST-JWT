package email

import (
	"github.com/diogox/REST-JWT/server/pkg/models"
	"gopkg.in/gomail.v2"
	"io/ioutil"
	"strings"
)

func NewEmailClient(from string, opts EmailClientOptions) *EmailClient {
	return &EmailClient{
		From:     from,
		Host:     opts.Host,
		Port:     opts.Port,
		Username: opts.Username,
		Password: opts.Password,
	}
}

type EmailClient struct {
	From     string
	Host     string
	Port     int
	Username string
	Password string
}

func (ec *EmailClient) SendEmail(user models.User, opts models.NewEmail) error {

	// Create new message
	m := gomail.NewMessage()
	m.SetHeader("From", ec.From)
	m.SetHeader("To", user.Email)
	m.SetHeader("Subject", opts.Subject)

	// Set email body
	emailBody, err := getBody()
	if err != nil {
		return err
	}
	m.SetBody("text/html", strings.Replace(emailBody, "%m", opts.Message, -1))

	// Connect to email server
	d := gomail.NewDialer(ec.Host, ec.Port, ec.Username, ec.Password)

	// Send the email to Bob, Cora and Dan.
	if err := d.DialAndSend(m); err != nil {
		return err
	}

	return nil
}

func getBody() (string, error) {
	contents, err := ioutil.ReadFile("./server/pkg/email/body.html")
	if err != nil {
		return "", err
	}

	return string(contents), nil
}
