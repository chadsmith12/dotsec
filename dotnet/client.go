package dotnet

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func InitSecrets(projectPath string) error {
	cmd := exec.Command("dotnet", "user-secrets", "init")
	if projectPath != "" {
		cmd.Args = append(cmd.Args, "--project")
		cmd.Args = append(cmd.Args, projectPath)
	}

	return logAndRunCommand(cmd)
}

func SetSecret(projectPath, key, value string) error {
	cmd := exec.Command("dotnet", "user-secrets", "set", key, value)
	if projectPath != "" {
		cmd.Args = append(cmd.Args, "--project")
		cmd.Args = append(cmd.Args, projectPath)
	}

	return logAndRunCommand(cmd)
}

func ListSecrets(projectPath string) (bytes.Buffer, error) {
	cmd := exec.Command("dotnet", "user-secrets", "list")
	if projectPath != "" {
		cmd.Args = append(cmd.Args, "--project")
		cmd.Args = append(cmd.Args, projectPath)
	}

	stdOut, errOut, err := runCmd(cmd)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error running %s %s - %v\n", cmd.Args[0], cmd.Args[1], errOut.String())
		return stdOut, fmt.Errorf("%s %s error: %w", cmd.Args[0], cmd.Args[1], err)
	}

	return stdOut, nil

}

func ParseSecrets(bufr bytes.Buffer) ([]string, error) {
	scanner := bufio.NewScanner(&bufr)
	if !scanner.Scan() {
		if err := scanner.Err(); err != nil {
			return []string{}, err
		}

		return []string{}, nil
	}

	firstline := scanner.Text()
	if strings.Contains(firstline, "No secrets configured") {
		return []string{}, nil
	}

	secrets := []string{firstline}
	for scanner.Scan() {
		secrets = append(secrets, scanner.Text())
	}

	return secrets, nil
}

func runCmd(cmd *exec.Cmd) (bytes.Buffer, bytes.Buffer, error) {
	var stdOut bytes.Buffer
	var errOut bytes.Buffer
	cmd.Stdout = &stdOut
	cmd.Stderr = &errOut

	err := cmd.Run()

	return stdOut, errOut, err
}

func logAndRunCommand(cmd *exec.Cmd) error {
	stdOut, errOut, err := runCmd(cmd)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error running %s %s - %v\n", cmd.Args[0], cmd.Args[1], errOut.String())
		return fmt.Errorf("%s %s error: %w", cmd.Args[0], cmd.Args[1], err)
	}

	fmt.Fprintf(os.Stderr, "%v", stdOut.String())
	return nil
}
