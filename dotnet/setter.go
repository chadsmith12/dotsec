package dotnet

import "github.com/chadsmith12/dotsec/secrets"

type DotNetSetter struct {
	project string
}

func NewSetter(project string) DotNetSetter {
	return DotNetSetter{ project: project }
}


func (setter DotNetSetter) SetSecrets(secrets []secrets.SecretData) error {
	if err := InitSecrets(setter.project); err != nil {
		return err
	}

	for _, secret := range secrets {
		SetSecret(setter.project, secret.Key, secret.Value)
	}

	return nil
}
