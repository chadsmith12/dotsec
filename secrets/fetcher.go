package secrets 

import "github.com/chadsmith12/dotsec/passbolt"

// A SecretsFetcher is an interface you implement for different ways to fetch secrets from their underlying sources.
// Right now we support env file and dotnet user-secrets
type SecretsFetcher interface {
	FetchSecrets() ([]passbolt.SecretData, error)
}

// A SecretsSetter is an interface you implement for different ways to set secrets retrieved from passbolt.
// This will set the secrets into their underlying source.
type SecretsSetter interface {
	SetSecrets([]passbolt.SecretData) error
}
