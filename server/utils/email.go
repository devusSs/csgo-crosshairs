package utils

import (
	"bytes"
	"errors"
	"html/template"
	"net/mail"
	"os"
	"path/filepath"

	"github.com/devusSs/crosshairs/config"
	"github.com/devusSs/crosshairs/database"
	"gopkg.in/gomail.v2"
)

var (
	mailServer *gomail.Dialer
	mailSender string
)

type EmailData struct {
	URL     string
	Subject string
}

type EmailDataAdmin struct {
	Subject string
	Data    interface{}
}

func InitMail(cfg *config.Config) {
	mailServer = gomail.NewDialer(cfg.SMTPHost, cfg.SMTPPort, cfg.SMTPUser, cfg.SMTPPass)

	mailSender = cfg.EmailFrom
}

func parseTemplateDir(dir string) (*template.Template, error) {
	var paths []string
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			paths = append(paths, path)
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return template.ParseFiles(paths...)
}

func SendVerificationMail(user *database.UserAccount, data *EmailData) error {
	var body bytes.Buffer

	template, err := parseTemplateDir("templates")
	if err != nil {
		return err
	}

	if err := template.ExecuteTemplate(&body, "verificationCode.html", &data); err != nil {
		return err
	}

	m := gomail.NewMessage()

	m.SetHeader("From", mailSender)
	m.SetHeader("To", user.EMail)
	m.SetHeader("Subject", data.Subject)
	m.SetBody("text/html", body.String())

	if err := mailServer.DialAndSend(m); err != nil {
		return err
	}

	return nil
}

func SendAdminMail(user *database.UserAccount, data *EmailDataAdmin) error {
	if user.Role != "admin" {
		return errors.New("user is not an admin")
	}

	var body bytes.Buffer

	template, err := parseTemplateDir("templates")
	if err != nil {
		return err
	}

	if err := template.ExecuteTemplate(&body, "adminEvent.html", &data); err != nil {
		return err
	}

	m := gomail.NewMessage()

	m.SetHeader("From", mailSender)
	m.SetHeader("To", user.EMail)
	m.SetHeader("Subject", data.Subject)
	m.SetBody("text/html", body.String())

	if err := mailServer.DialAndSend(m); err != nil {
		return err
	}

	return nil
}

func IsEmailValid(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}
