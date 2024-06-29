package secrets

import "strings"

type SecretData struct {
	Key   string
	Value string
}

func SecretDataFromSlice(values []string) []SecretData {
	secretValues := make([]SecretData, 0, len(values))

	for _, value := range values {
		key, secret, found := strings.Cut(value, "=")
		if !found {
			continue
		}
		secretValues = append(secretValues, SecretData{Key: strings.TrimSpace(key), Value: strings.TrimSpace(secret)})
	}

	return secretValues
}
// A SecretsFetcher is an interface you implement for different ways to fetch secrets from their underlying sources.
// Right now we support env file and dotnet user-secrets
type SecretsFetcher interface {
	FetchSecrets() ([]SecretData, error)
}

// A SecretsSetter is an interface you implement for different ways to set secrets retrieved from passbolt.
// This will set the secrets into their underlying source.
type SecretsSetter interface {
	SetSecrets([]SecretData) error
}

