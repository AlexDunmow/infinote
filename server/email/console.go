package email

import "fmt"

type Console struct {
}

// ReceivedSignup from the user
func (c *Console) ReceivedSignup(email string) error {
	fmt.Println("Send ReceivedSignup: ", "email:", email)
	return nil
}

// ForgotPassword request from the user
func (c *Console) ForgotPassword(email string) error {
	fmt.Println("Send ForgotPassword: ", "email:", email)
	return nil
}
