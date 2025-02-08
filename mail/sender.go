package mail

import (
	"fmt"
	"net/smtp"

	"github.com/jordan-wright/email"
)

const (
	smtpAuthAddress = "smtp.gmail.com"
	smtpServerAdr   = "smtp.gmail.com:587"
)

type EmailSender interface {
	SendEmail(subject string, content string, to []string, cc []string, bcc []string, attachFiles []string) error
}

type GmailSender struct {
	name              string
	fromEmailAdr      string
	fromEmailPassword string
}

func NewGmailSender(name string,
	fromEmailAdr string,
	fromEmailPassword string) EmailSender {
	return &GmailSender{
		name,
		fromEmailAdr,
		fromEmailPassword,
	}
}

func (sender *GmailSender) SendEmail(subject string, content string, to []string, cc []string, bcc []string, attachFiles []string) error {
	e := email.NewEmail()
	e.From = fmt.Sprintf("%s <%s>", sender.name, sender.fromEmailAdr)
	e.Subject = subject
	e.HTML = []byte(content)
	e.To = to
	e.Bcc = bcc
	e.Cc = cc

	for _, a := range attachFiles {
		_, err := e.AttachFile(a)
		if err != nil {
			return fmt.Errorf("failed to attach file %s: %w", a, err)
		}
	}

	auth := smtp.PlainAuth("", sender.fromEmailAdr, sender.fromEmailPassword, smtpAuthAddress)
	return e.Send(smtpServerAdr, auth)
}
