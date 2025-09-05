package env_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/chadsmith12/dotsec/env"
	"github.com/chadsmith12/dotsec/secrets"
)

func createTempEnvFile(t *testing.T, content string) string {
	t.Helper()
	tempDir := t.TempDir()
	envFile := filepath.Join(tempDir, ".env")

	err := os.WriteFile(envFile, []byte(content), 0600)
	if err != nil {
		t.Fatalf("Failed to create temp env file: %v", err)
	}

	return envFile
}

func readEnvFile(t *testing.T, path string) string {
	t.Helper()
	content, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("Failed to read env file: %v", err)
	}
	return string(content)
}

func TestSetSecrets_EmptyFile(t *testing.T) {
	envFile := createTempEnvFile(t, "")

	secretsData := []secrets.SecretData{
		{Key: "API_KEY", Value: "secret123"},
		{Key: "DB_PASSWORD", Value: "password456"},
	}

	setter := env.NewSetter(envFile)
	err := setter.SetSecrets(secretsData)
	if err != nil {
		t.Fatalf("SetSecrets failed: %v", err)
	}

	content := readEnvFile(t, envFile)

	if !strings.Contains(content, `API_KEY="secret123"`) {
		t.Error("Expected API_KEY to be set")
	}
	if !strings.Contains(content, `DB_PASSWORD="password456"`) {
		t.Error("Expected DB_PASSWORD to be set")
	}
}

func TestSetSecrets_UpdateExistingValue(t *testing.T) {
	initialContent := `API_KEY="oldvalue"
DB_HOST=localhost
`
	envFile := createTempEnvFile(t, initialContent)

	secretsData := []secrets.SecretData{
		{Key: "API_KEY", Value: "newvalue"},
	}

	setter := env.NewSetter(envFile)
	err := setter.SetSecrets(secretsData)
	if err != nil {
		t.Fatalf("SetSecrets failed: %v", err)
	}

	content := readEnvFile(t, envFile)

	if !strings.Contains(content, `API_KEY="newvalue"`) {
		t.Error("Expected API_KEY to be updated to newvalue")
	}
	if strings.Contains(content, "oldvalue") {
		t.Error("Old value should be replaced")
	}
	if !strings.Contains(content, "DB_HOST=localhost") {
		t.Error("Existing unchanged values should be preserved")
	}
}

func TestSetSecrets_PreserveUnchangedValues(t *testing.T) {
	initialContent := `API_KEY="keep_this"
DB_HOST=localhost
DB_PORT=5432
SECRET_TOKEN="also_keep_this"
`
	envFile := createTempEnvFile(t, initialContent)

	secretsData := []secrets.SecretData{
		{Key: "DB_HOST", Value: "newhost"},
	}

	setter := env.NewSetter(envFile)
	err := setter.SetSecrets(secretsData)
	if err != nil {
		t.Fatalf("SetSecrets failed: %v", err)
	}

	content := readEnvFile(t, envFile)

	if !strings.Contains(content, `API_KEY="keep_this"`) {
		t.Error("API_KEY should be preserved")
	}
	if !strings.Contains(content, `DB_HOST="newhost"`) {
		t.Error("DB_HOST should be updated")
	}
	if !strings.Contains(content, "DB_PORT=5432") {
		t.Error("DB_PORT should be preserved")
	}
	if !strings.Contains(content, `SECRET_TOKEN="also_keep_this"`) {
		t.Error("SECRET_TOKEN should be preserved")
	}
}

func TestSetSecrets_HandleQuotedValues(t *testing.T) {
	initialContent := `API_KEY="quoted_value"
DB_PASSWORD='single_quoted'
PLAIN_VALUE=no_quotes
`
	envFile := createTempEnvFile(t, initialContent)

	secretsData := []secrets.SecretData{
		{Key: "API_KEY", Value: "new_quoted"},
		{Key: "DB_PASSWORD", Value: "new_single"},
		{Key: "PLAIN_VALUE", Value: "new_plain"},
	}

	setter := env.NewSetter(envFile)
	err := setter.SetSecrets(secretsData)
	if err != nil {
		t.Fatalf("SetSecrets failed: %v", err)
	}

	content := readEnvFile(t, envFile)

	if !strings.Contains(content, `API_KEY="new_quoted"`) {
		t.Error("Double quoted value should be updated")
	}
	if !strings.Contains(content, `DB_PASSWORD="new_single"`) {
		t.Error("Single quoted value should be updated")
	}
	if !strings.Contains(content, `PLAIN_VALUE="new_plain"`) {
		t.Error("Plain value should be updated")
	}
}

func TestSetSecrets_AddNewSecrets(t *testing.T) {
	initialContent := `EXISTING_KEY="existing_value"
`
	envFile := createTempEnvFile(t, initialContent)

	secretsData := []secrets.SecretData{
		{Key: "NEW_KEY1", Value: "new_value1"},
		{Key: "NEW_KEY2", Value: "new_value2"},
	}

	setter := env.NewSetter(envFile)
	err := setter.SetSecrets(secretsData)
	if err != nil {
		t.Fatalf("SetSecrets failed: %v", err)
	}

	content := readEnvFile(t, envFile)

	if !strings.Contains(content, `EXISTING_KEY="existing_value"`) {
		t.Error("Existing key should be preserved")
	}
	if !strings.Contains(content, `NEW_KEY1="new_value1"`) {
		t.Error("NEW_KEY1 should be added")
	}
	if !strings.Contains(content, `NEW_KEY2="new_value2"`) {
		t.Error("NEW_KEY2 should be added")
	}
}

func TestSetSecrets_HandleCommentsAndEmptyLines(t *testing.T) {
	initialContent := `# This is a comment
API_KEY="value1"

# Another comment
DB_HOST=localhost

`
	envFile := createTempEnvFile(t, initialContent)

	secretsData := []secrets.SecretData{
		{Key: "API_KEY", Value: "updated_value"},
	}

	setter := env.NewSetter(envFile)
	err := setter.SetSecrets(secretsData)
	if err != nil {
		t.Fatalf("SetSecrets failed: %v", err)
	}

	content := readEnvFile(t, envFile)

	if !strings.Contains(content, "# This is a comment") {
		t.Error("First comment should be preserved")
	}
	if !strings.Contains(content, "# Another comment") {
		t.Error("Second comment should be preserved")
	}
	if !strings.Contains(content, `API_KEY="updated_value"`) {
		t.Error("API_KEY should be updated")
	}
	if !strings.Contains(content, "DB_HOST=localhost") {
		t.Error("DB_HOST should be preserved")
	}
}

func TestSetSecrets_HandleKeyWithoutValue(t *testing.T) {
	initialContent := `EMPTY_KEY=
ANOTHER_KEY
NORMAL_KEY="has_value"
`
	envFile := createTempEnvFile(t, initialContent)

	secretsData := []secrets.SecretData{
		{Key: "EMPTY_KEY", Value: "now_has_value"},
		{Key: "ANOTHER_KEY", Value: "also_has_value"},
	}

	setter := env.NewSetter(envFile)
	err := setter.SetSecrets(secretsData)
	if err != nil {
		t.Fatalf("SetSecrets failed: %v", err)
	}

	content := readEnvFile(t, envFile)

	if !strings.Contains(content, `EMPTY_KEY="now_has_value"`) {
		t.Error("EMPTY_KEY should be updated with value")
	}
	if !strings.Contains(content, `ANOTHER_KEY="also_has_value"`) {
		t.Error("ANOTHER_KEY should be updated with value")
	}
	if !strings.Contains(content, `NORMAL_KEY="has_value"`) {
		t.Error("NORMAL_KEY should be preserved")
	}
}

func TestSetSecrets_HandleValuesWithEquals(t *testing.T) {
	initialContent := `CONNECTION_STRING="server=localhost;database=test"
`
	envFile := createTempEnvFile(t, initialContent)

	secretsData := []secrets.SecretData{
		{Key: "CONNECTION_STRING", Value: "server=newhost;database=prod;user=admin"},
		{Key: "FORMULA", Value: "x=y+z"},
	}

	setter := env.NewSetter(envFile)
	err := setter.SetSecrets(secretsData)
	if err != nil {
		t.Fatalf("SetSecrets failed: %v", err)
	}

	content := readEnvFile(t, envFile)

	expected := `CONNECTION_STRING="server=newhost;database=prod;user=admin"`
	if !strings.Contains(content, expected) {
		t.Errorf("CONNECTION_STRING should handle equals in value, got: %s", content)
	}
	if !strings.Contains(content, `FORMULA="x=y+z"`) {
		t.Error("FORMULA should be added with equals in value")
	}
}

func TestSetSecrets_FilePermissionError(t *testing.T) {
	if os.Getuid() == 0 {
		t.Skip("Skipping permission test when running as root")
	}

	tempDir := t.TempDir()
	envFile := filepath.Join(tempDir, ".env")

	err := os.WriteFile(envFile, []byte("TEST=value"), 0000)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	secretsData := []secrets.SecretData{
		{Key: "API_KEY", Value: "secret"},
	}

	setter := env.NewSetter(envFile)
	err = setter.SetSecrets(secretsData)
	if err == nil {
		t.Error("Expected error due to file permissions")
	}
}

func TestSetSecrets_NonexistentDirectory(t *testing.T) {
	nonexistentPath := "/nonexistent/directory/.env"

	secretsData := []secrets.SecretData{
		{Key: "API_KEY", Value: "secret"},
	}

	setter := env.NewSetter(nonexistentPath)
	err := setter.SetSecrets(secretsData)
	if err == nil {
		t.Error("Expected error for nonexistent directory")
	}
}
