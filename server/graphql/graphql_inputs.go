package graphql

import "github.com/volatiletech/null"

// OnboardingInput used for mutations during onboarding process
type OnboardingInput struct {
	Email string      `json:"email"`
	Name  null.String `json:"name"`
}
