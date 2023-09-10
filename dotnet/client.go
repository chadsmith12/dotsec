package dotnet

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
)

func InitSecrets(projectPath string) {
	cmd := exec.Command("dotnet", "user-secrets", "init")
	if projectPath != "" {
		cmd.Args = append(cmd.Args, "--project")
		cmd.Args = append(cmd.Args, projectPath)
	}
	var stdOut bytes.Buffer
	var errOut bytes.Buffer
	cmd.Stdout = &stdOut
	cmd.Stderr = &errOut

	err := cmd.Run()
	
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error user-secrets init: %v\n", errOut.String())
		os.Exit(1)
	}

	fmt.Fprintf(os.Stderr, "%v\n", stdOut.String())
}
