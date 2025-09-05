package env

import "github.com/chadsmith12/dotsec/secrets"

type EnvFetcher struct {
	file string
}

func NewFetcher(project string) EnvFetcher {
	return EnvFetcher{file: project}
}

func (fetcher EnvFetcher) FetchSecrets() ([]secrets.SecretData, error) {
	values, err := GetSecrets(fetcher.file)
	if err != nil {
		return []secrets.SecretData{}, err
	}

	return values, nil

}
