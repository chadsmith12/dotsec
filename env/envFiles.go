package env

import (
	"fmt"
	"os"

	"github.com/chadsmith12/dotsec/secrets"
	"github.com/hashicorp/go-envparse"
)

func GetSecrets(envFile string) ([]secrets.SecretData, error) {
	file, err := os.Open(envFile)
	if err != nil {
		return []secrets.SecretData{}, err
	}

	parsedSecrets, err := envparse.Parse(file)
	if err != nil {
		return []secrets.SecretData{}, err
	}

	secretData := make([]secrets.SecretData, 0, len(parsedSecrets))
	for key, value := range parsedSecrets {
		secretData = append(secretData, secrets.SecretData{Key: key, Value: value})
	}

	return secretData, nil
}

func setSecrets(envFile string, secretsData []secrets.SecretData) error {
	file, err := os.OpenFile(envFile, os.O_RDWR|os.O_CREATE, 0600)
	if err != nil {
		return fmt.Errorf("SetSecrets - failed to open file. %w", err)
	}
	defer file.Close()

	// envLines := make([]string, len(secretsData))
	secretsMap := createSecretsMap(secretsData)
	currentSecrets, err := envparse.Parse(file)

	if err != nil {
		return fmt.Errorf("SetSecrets - failed to parse env file %s. %w", envFile, err)
	}
	file.Close()

	// go through our current secrets and update them if we need to
	// index := 0
	for key := range currentSecrets {
		if val, ok := secretsMap[key]; ok {
			//envLines[index] = fmt.Sprintf("%s=\"%s\"", key, val)
			//delete(secretsMap, key)
			//index++
			currentSecrets[key] = val
		}
	}

	for key, value := range secretsMap {
		currentSecrets[key] = value
		// envLines[index] = fmt.Sprintf("%s=\"%s\"", key, value)
		// index++
	}

	file, err = os.Create(envFile)
	if err != nil {
		return fmt.Errorf("SetSecrets - failed to create env file %s. %w", envFile, err)
	}
	defer file.Close()

	for key, value := range currentSecrets {
		line := fmt.Sprintf("%s=\"%s\"", key, value)
		fmt.Fprintln(file, line)
	}

	return nil
}

func createSecretsMap(secretsData []secrets.SecretData) map[string]string {
	secretMap := make(map[string]string, len(secretsData))

	for i := range secretsData {
		secret := secretsData[i]
		secretMap[secret.Key] = secret.Value
	}

	return secretMap
}
