package utils

import (
	"bytes"
	"crypto/tls"
	"html/template"
	"net/mail"
	"os"
	"path/filepath"

	"github.com/devusSs/crosshairs/config"
	"github.com/devusSs/crosshairs/database"
	"github.com/k3a/html2text"
	"gopkg.in/gomail.v2"
)

const (
	MailVerfication          = "mailVerification"
	MailVerificationPassword = "passwordVerification"
)

type EmailData struct {
	URL     string
	Subject string
}

func ParseTemplateDir(dir string) (*template.Template, error) {
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

func SendEmail(user *database.UserAccount, data *EmailData, cfg *config.Config, mailType string) error {
	from := cfg.EmailFrom
	smtpPass := cfg.SMTPPass
	smtpUser := cfg.SMTPUser
	to := user.EMail
	smtpHost := cfg.SMTPHost
	smtpPort := cfg.SMTPPort

	var body bytes.Buffer

	template, err := ParseTemplateDir("templates")
	if err != nil {
		return err
	}

	if mailType == MailVerfication {
		if err := template.ExecuteTemplate(&body, "verificationCode.html", &data); err != nil {
			return err
		}
	}

	if mailType == MailVerificationPassword {
		if err := template.ExecuteTemplate(&body, "passwordVerification.html", &data); err != nil {
			return err
		}
	}

	m := gomail.NewMessage()

	m.SetHeader("From", from)
	m.SetHeader("To", to)
	m.SetHeader("Subject", data.Subject)
	m.SetBody("text/html", body.String())
	m.AddAlternative("text/plain", html2text.HTML2Text(body.String()))

	d := gomail.NewDialer(smtpHost, smtpPort, smtpUser, smtpPass)
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	if err := d.DialAndSend(m); err != nil {
		return err
	}

	return nil
}

func IsEmailValid(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}
