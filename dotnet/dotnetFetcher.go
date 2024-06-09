package dotnet

import (
	"github.com/chadsmith12/dotsec/passbolt"
)

type DotNetFetcher struct {
	project string
}

func NewFetcher(project string) DotNetFetcher {
	return DotNetFetcher{ project: project }
}

func (fetcher DotNetFetcher) FetchSecrets() ([]passbolt.SecretData, error) {
	stdOut, err := ListSecrets(fetcher.project)
	if err != nil {
		return []passbolt.SecretData{}, err
	}
	values, err := ParseSecrets(stdOut)
	if err != nil {
		return []passbolt.SecretData{}, err
	}

	secretsData := passbolt.SecretDataFromSlice(values)
	return secretsData, nil
}
