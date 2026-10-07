package mailer

import (
	"github.com/SrJui/yap/internal/config"
	"github.com/resend/resend-go/v4"
)

type Mailer struct {
	config config.MailConfig
	client *resend.Client
}

func New(cfg config.MailConfig) *Mailer {
	return &Mailer{
		config: cfg,
		client: resend.NewClient(cfg.APIKey),
	}
}

func (m *Mailer) Send(to []string, subject, body string) error {
	params := &resend.SendEmailRequest{
		From:    m.config.From,
		To:      to,
		Subject: subject,
		Html:    body,
	}

	_, err := m.client.Emails.Send(params)
	if err != nil {
		return err
	}

	return nil
}
