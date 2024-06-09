package secretsfetcher 

import "github.com/chadsmith12/dotsec/passbolt"

type SecretsFetcher interface {
	FetchSecrets() ([]passbolt.SecretData, error)
}
