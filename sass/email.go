package sass

import "github.com/go-mail/mail"

// EmailSenderInterface defines
type EmailSenderInterface interface {
	SendEmail(Configuration, *email.Message) error
}

// LocalHostEmailSender Mail sender
type LocalHostEmailSender struct{}

// SendEmail bind a function
func (s *LocalHostEmailSender) SendEmail(config Configuration, m *mail.Message) error {

	return nil
}

// EmailSender defines
var EmailSender EmailSenderInterface = &LocalHostEmailSender{}
