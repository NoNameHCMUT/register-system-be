package email

import (
	"fmt"
	"net/smtp"
	"strings"
)

type EmailSender interface {
	Send(to, subject, body string) error
}

type smtpSender struct {
	host     string
	port     string
	user     string
	password string
	from     string
}

func NewSMTPSender(host, port, user, password, from string) EmailSender {
	return &smtpSender{host: host, port: port, user: user, password: password, from: from}
}

type noopSender struct{}

func NewNoopSender() EmailSender {
	return &noopSender{}
}

func (s *noopSender) Send(to, subject, body string) error {
	fmt.Printf("[EMAIL] to=%s subject=%s\n", to, subject)
	return nil
}

func (s *smtpSender) Send(to, subject, body string) error {
	auth := smtp.PlainAuth("", s.user, s.password, s.host)
	addr := s.host + ":" + s.port

	var msg strings.Builder
	msg.WriteString("From: " + s.from + "\r\n")
	msg.WriteString("To: " + to + "\r\n")
	msg.WriteString("Subject: " + subject + "\r\n")
	msg.WriteString("MIME-Version: 1.0\r\n")
	msg.WriteString("Content-Type: text/plain; charset=\"utf-8\"\r\n")
	msg.WriteString("\r\n")
	msg.WriteString(body)

	return smtp.SendMail(addr, auth, s.user, []string{to}, []byte(msg.String()))
}

func NotifyApplicationStatus(sender EmailSender, to string, status string, projectName string) error {
	subject := "Application Update: " + projectName
	body := fmt.Sprintf("Your application for project \"%s\" has been updated to status: %s", projectName, status)
	return sender.Send(to, subject, body)
}

func NotifyAccountStatus(sender EmailSender, to string, approved bool) error {
	var subject, body string
	if approved {
		subject = "Account Approved"
		body = "Your account has been approved. You can now log in to the system."
	} else {
		subject = "Account Rejected"
		body = "Your account registration has been rejected."
	}
	return sender.Send(to, subject, body)
}
