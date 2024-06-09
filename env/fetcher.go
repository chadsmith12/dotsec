package env

import "github.com/chadsmith12/dotsec/passbolt"

type EnvFetcher struct {
	file string
}

func NewFetcher(project string) EnvFetcher {
	return EnvFetcher{ file: project }
}

func (fetcher EnvFetcher) FetchSecrets() ([]passbolt.SecretData, error) {
	values, err := GetSecrets(fetcher.file)
	if err != nil {
		return []passbolt.SecretData{}, err 
	}

	return values, nil

}
