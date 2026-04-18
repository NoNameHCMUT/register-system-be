package email

import (
	"fmt"
	"mime/quotedprintable"
	"net/smtp"
	"strings"
	"time"
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

	textBody, htmlBody := formatEmailBody(subject, body)
	msg := buildMultipartAlternativeMessage(s.from, to, subject, textBody, htmlBody)
	return smtp.SendMail(addr, auth, s.user, []string{to}, []byte(msg))
}

func buildMultipartAlternativeMessage(from, to, subject, textBody, htmlBody string) string {
	boundary := fmt.Sprintf("boundary_%d", time.Now().UnixNano())

	var msg strings.Builder
	msg.WriteString("From: " + from + "\r\n")
	msg.WriteString("To: " + to + "\r\n")
	msg.WriteString("Subject: " + subject + "\r\n")
	msg.WriteString("MIME-Version: 1.0\r\n")
	msg.WriteString("Content-Type: multipart/alternative; boundary=\"" + boundary + "\"\r\n")
	msg.WriteString("\r\n")

	// text/plain
	msg.WriteString("--" + boundary + "\r\n")
	msg.WriteString("Content-Type: text/plain; charset=\"utf-8\"\r\n")
	msg.WriteString("Content-Transfer-Encoding: quoted-printable\r\n")
	msg.WriteString("\r\n")
	msg.WriteString(encodeQuotedPrintable(textBody))
	msg.WriteString("\r\n")

	// text/html
	msg.WriteString("--" + boundary + "\r\n")
	msg.WriteString("Content-Type: text/html; charset=\"utf-8\"\r\n")
	msg.WriteString("Content-Transfer-Encoding: quoted-printable\r\n")
	msg.WriteString("\r\n")
	msg.WriteString(encodeQuotedPrintable(htmlBody))
	msg.WriteString("\r\n")

	msg.WriteString("--" + boundary + "--\r\n")
	return msg.String()
}

func encodeQuotedPrintable(s string) string {
	var b strings.Builder
	w := quotedprintable.NewWriter(&b)
	_, _ = w.Write([]byte(s))
	_ = w.Close()
	return b.String()
}

// Keep the EmailSender interface stable: callers provide plain text; we build a
// matching HTML version for nicer rendering.
func formatEmailBody(subject, plain string) (textBody string, htmlBody string) {
	textBody = strings.TrimSpace(plain) + "\n"

	escaped := htmlEscape(textBody)
	htmlBody = "<!doctype html>" +
		"<html><head><meta charset=\"utf-8\"></head>" +
		"<body style=\"font-family:Arial,Helvetica,sans-serif;background:#f6f7fb;margin:0;padding:24px;\">" +
		"<div style=\"max-width:640px;margin:0 auto;background:#ffffff;border:1px solid #e5e7eb;border-radius:12px;overflow:hidden;\">" +
		"<div style=\"padding:20px 24px;background:#111827;color:#ffffff;\">" +
		"<div style=\"font-size:16px;font-weight:700;\">Register System</div>" +
		"<div style=\"font-size:13px;opacity:.85;margin-top:4px;\">" + htmlEscape(subject) + "</div>" +
		"</div>" +
		"<div style=\"padding:20px 24px;color:#111827;line-height:1.55;font-size:14px;\">" +
		"<div style=\"white-space:pre-wrap\">" + escaped + "</div>" +
		"</div>" +
		"<div style=\"padding:14px 24px;border-top:1px solid #e5e7eb;color:#6b7280;font-size:12px;\">" +
		"Please do not reply to this email." +
		"</div>" +
		"</div>" +
		"</body></html>"

	return textBody, htmlBody
}

func htmlEscape(s string) string {
	replacer := strings.NewReplacer(
		"&", "&amp;",
		"<", "&lt;",
		">", "&gt;",
		"\"", "&quot;",
		"'", "&#39;",
	)
	return replacer.Replace(s)
}

func NotifyApplicationStatus(sender EmailSender, to string, status string, projectName string) error {
	subject := "Update on your application - " + projectName
	body := fmt.Sprintf(
		"Hello,\n\n"+
			"We wanted to let you know that your application for the project \"%s\" has been updated.\n\n"+
			"Current status: %s\n\n"+
			"You can log in to the system to view details.\n\n"+
			"Regards,\n"+
			"Register System",
		projectName,
		status,
	)
	return sender.Send(to, subject, body)
}

func NotifyAccountStatus(sender EmailSender, to string, approved bool) error {
	var subject, body string
	if approved {
		subject = "Your account is approved - welcome!"
		body = "Hello,\n\n" +
			"Good news! Your account has been approved and is now active.\n\n" +
			"You can log in to the system using the username and password you registered with.\n\n" +
			"Regards,\n" +
			"Register System"
	} else {
		subject = "Your account request wasn't approved"
		body = "Hello,\n\n" +
			"Thank you for registering. Unfortunately, your account request wasn't approved at this time.\n\n" +
			"If you believe this is a mistake, please contact the system administrator or try registering again with correct information.\n\n" +
			"Regards,\n" +
			"Register System"
	}
	return sender.Send(to, subject, body)
}
