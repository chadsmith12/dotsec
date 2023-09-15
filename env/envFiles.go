package env

import (
	"fmt"
	"os"
	"path"

	"github.com/chadsmith12/dotsec/passbolt"
	"github.com/hashicorp/go-envparse"
)

func SetSecrets(project, env string, secretsData []passbolt.SecretData) error {
	filePath := path.Join(project, env)
	file, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE, 0600)
	if err != nil {
		return fmt.Errorf("SetSecrets - failed to open file. %w", err)
	}
	defer file.Close()

	envLines := make([]string, len(secretsData))
	secretsMap := createSecretsMap(secretsData)
	secrets, err := envparse.Parse(file)
	
	if err != nil {
		return fmt.Errorf("SetSecrets - failed to parse env file %s. %w", filePath, err)
	}
	file.Close()
	index := 0

	for key  := range secrets {
		if val, ok := secretsMap[key]; ok {
			envLines[index] = fmt.Sprintf("%s=\"%s\"", key, val)
			delete(secretsMap, key)
			index++
		} 
	}
	
	for key, value := range secretsMap {
		envLines[index] = fmt.Sprintf("%s=\"%s\"", key, value)
		index++
	}

	file, err = os.Create(filePath)
	if err != nil {
		return fmt.Errorf("SetSecrets - failed to create env file %s. %w", filePath, err)
	}
	defer file.Close()
	
	for _, line := range envLines {
		fmt.Fprintln(file, line)
	}

	return nil
}

func createSecretsMap(secretsData []passbolt.SecretData) map[string]string {
	secretMap := make(map[string]string, len(secretsData))

	for i := range secretsData {
		secret := secretsData[i]
		secretMap[secret.Key] = secret.Value
	}

	return secretMap
}
