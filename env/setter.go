package env

import "github.com/chadsmith12/dotsec/secrets"

type EnvSetter struct {
	envFile string
}

func NewSetter(envFile string) EnvSetter {
	return EnvSetter{envFile: envFile}
}

func (setter EnvSetter) SetSecrets(secrets []secrets.SecretData) error {
	err := setSecrets(setter.envFile, secrets)

	return err
}
