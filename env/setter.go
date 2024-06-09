package env

import "github.com/chadsmith12/dotsec/passbolt"

type EnvSetter struct {
	envFile string
}

func NewSetter(envFile string) EnvSetter {
	return EnvSetter{ envFile: envFile }
}

func (setter EnvSetter) SetSecrets(secrets []passbolt.SecretData) error {
	err := setSecrets(setter.envFile, secrets)

	return err
}
