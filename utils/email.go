package utils

import (
	"fmt"
	"net/smtp"
	"os"
)

func SendEmail(to string, subject string, body string) error {
	from := os.Getenv("FROM_EMAIL")
	user := os.Getenv("SMTP_USER")
	password := os.Getenv("SMTP_PASS")
	host := os.Getenv("SMTP_HOST")
	port := os.Getenv("SMTP_PORT")

	if host == "" || port == "" || user == "" || password == "" || from == "" {
		return fmt.Errorf("smtp config incomplete")
	}

	addr := fmt.Sprintf("%s:%s", host, port)

	msg := []byte("To: " + to + "\r\n" +
		"From: " + from + "\r\n" +
		"Subject: " + subject + "\r\n" +
		"MIME-Version: 1.0\r\n" +
		"Content-Type: text/html; charset=UTF-8\r\n" +
		"\r\n" + body)

	auth := smtp.PlainAuth("", user, password, host)
	return smtp.SendMail(addr, auth, from, []string{to}, msg)
}
