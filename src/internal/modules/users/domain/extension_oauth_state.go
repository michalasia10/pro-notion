package domain

// ExtensionOAuthStateService handles generating and validating OAuth state for extension flows.
type ExtensionOAuthStateService interface {
	Generate(redirectURI string) (string, error)
	Validate(state string) (string, error)
}
