package email

import "errors"

// Mailer Note email implementation of notifier
type Mailer struct {
}

func (m *Mailer) ReceivedSignup(email string) error {
	return errors.New("not implemented")
}
func (m *Mailer) ForgotPassword(email string) error {
	return errors.New("not implemented")
}
