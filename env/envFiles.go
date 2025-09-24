package env

import (
	"bufio"
	"fmt"
	"os"
	"strings"

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

func stripQuotes(value string) string {
	if len(value) <= 2 {
		return value
	}
	if value[0] == '"' && value[len(value)-1] == '"' {
		return value[1 : len(value)-1]
	}
	if value[0] == '\'' && value[len(value)-1] == '\'' {
		return value[1 : len(value)-1]
	}

	return value
}

func setSecrets(envFile string, secretsData []secrets.SecretData) error {
	currEnvFile, err := os.OpenFile(envFile, os.O_RDWR|os.O_CREATE, 0600)
	if err != nil {
		return fmt.Errorf("SetSecrets - failed to open file. %w", err)
	}
	

	tempEnvFile, err := os.CreateTemp("", ".env.temp")
	if err != nil {
		currEnvFile.Close()
		return fmt.Errorf("SetSecrets - failed to create temporary file. %w", err)
	}
	defer func() {
		tempEnvFile.Close()
		if err != nil {
			os.Remove(tempEnvFile.Name())
		}
	}()

	secretsMap := createSecretsMap(secretsData)
	scanner := bufio.NewScanner(currEnvFile)
	writer := bufio.NewWriter(tempEnvFile)

	for scanner.Scan() {
		line := scanner.Text()
		if err := processExistingLine(line, writer, secretsMap); err != nil {
			return fmt.Errorf("SetSecrets - failed to process line: %w", err)
		}
	}

	if err := scanner.Err(); err != nil {
		currEnvFile.Close()
		return fmt.Errorf("SetSecrets - failed to scan file: %w", err)
	}

	if err := writeRemainingSecrets(writer, secretsMap); err != nil {
		currEnvFile.Close()
		return fmt.Errorf("SetSecrets - failed to write remaining secrets: %w", err)
	}

	if err := writer.Flush(); err != nil {
		currEnvFile.Close()
		return fmt.Errorf("SetSecrets - failed to flush writer: %w", err)
	}

	currEnvFile.Close()
	tempEnvFile.Close()

	if err := os.Remove(envFile); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("SetSecrets - failed to remove old file: %w", err)
	}

	if err := os.Rename(tempEnvFile.Name(), currEnvFile.Name()); err != nil {
		return fmt.Errorf("SetSecrets - failed to rename temp file: %w", err)
	}

	return nil
}

func processExistingLine(line string, writer *bufio.Writer, secretsMap map[string]string) error {
	key, currentValue, hasValue := parseEnvLine(line)

	if key == "" {
		return writeLineToFile(writer, line)
	}

	secretValue, found := secretsMap[key]
	if !found {
		return writeLineToFile(writer, line)
	}

	if shouldUpdateValue(currentValue, secretValue, hasValue) {
		updatedLine := formatEnvLine(key, secretValue)
		if err := writeLineToFile(writer, updatedLine); err != nil {
			return err
		}
	} else {
		if err := writeLineToFile(writer, line); err != nil {
			return err
		}
	}

	delete(secretsMap, key)
	return nil
}

func writeRemainingSecrets(writer *bufio.Writer, secretsMap map[string]string) error {
	for key, value := range secretsMap {
		line := formatEnvLine(key, value)
		if err := writeLineToFile(writer, line); err != nil {
			return err
		}
	}
	return nil
}

func parseEnvLine(line string) (key, value string, hasValue bool) {
	line = strings.TrimSpace(line)

	if line == "" || strings.HasPrefix(line, "#") {
		return "", "", false
	}

	parts := strings.SplitN(line, "=", 2)
	if len(parts) == 1 {
		return strings.TrimSpace(parts[0]), "", false
	}

	return strings.TrimSpace(parts[0]), strings.TrimSpace(parts[1]), true
}

func shouldUpdateValue(currentValue, secretValue string, hasValue bool) bool {
	if !hasValue {
		return true
	}
	return stripQuotes(currentValue) != secretValue
}

func formatEnvLine(key, value string) string {
	return fmt.Sprintf("%s=\"%s\"", key, value)
}

func writeLineToFile(writer *bufio.Writer, line string) error {
	if _, err := writer.WriteString(line); err != nil {
		return err
	}
	return writer.WriteByte('\n')
}

func createSecretsMap(secretsData []secrets.SecretData) map[string]string {
	secretMap := make(map[string]string, len(secretsData))

	for i := range secretsData {
		secret := secretsData[i]
		secretMap[secret.Key] = secret.Value
	}

	return secretMap
}
