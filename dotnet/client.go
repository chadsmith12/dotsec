package dotnet

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
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

func logAndRunCommand(cmd *exec.Cmd) error {
	var stdOut bytes.Buffer
	var errOut bytes.Buffer
	cmd.Stdout = &stdOut
	cmd.Stderr = &errOut

	err := cmd.Run()

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error running %s %s - %v\n", cmd.Args[0], cmd.Args[1], errOut.String())
		return fmt.Errorf("%s %s error: %w", cmd.Args[0], cmd.Args[1], err)
	}

	fmt.Fprintf(os.Stderr, "%v", stdOut.String())
	return nil
}
