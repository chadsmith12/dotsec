package dotnet

import "github.com/chadsmith12/dotsec/passbolt"

type DotNetSetter struct {
	project string
}

func NewSetter(project string) DotNetSetter {
	return DotNetSetter{ project: project }
}


func (setter DotNetSetter) SetSecrets(secrets []passbolt.SecretData) error {
	if err := InitSecrets(setter.project); err != nil {
		return err
	}

	for _, secret := range secrets {
		SetSecret(setter.project, secret.Key, secret.Value)
	}

	return nil
}
