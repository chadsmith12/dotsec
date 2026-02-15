package input

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"golang.org/x/term"
)

// PromptUser will prompt the user for input and return the input as a string. Setting asPassword to true will not echo out the characters to the user.
func PromptUser(prompt string, asPassword bool) (string, error) {
	fmt.Print(prompt)

	if !asPassword {
		inputString, err := readInput()
		return inputString, err
	}

	inputBuf, err := term.ReadPassword(int(os.Stdin.Fd()))
	if err != nil {
		return "", err
	}

	return string(inputBuf), nil
}

// PromptYesOrNo will prompt the user as a yes or no question and will return a boolean if they indicated yes or not
func PromptYesOrNo(prompt string) (bool, error) {
	for {
		confirm, err := PromptUser(prompt, false)
		if err != nil {
			return false, err
		}
		trimmed := strings.ToLower(strings.TrimSpace(confirm))
		if trimmed == "n" {
			return false, nil
		}
		if trimmed == "y" {
			return true, nil
		}
	}
}

func readInput() (string, error) {
	reader := bufio.NewScanner(os.Stdin)
	reader.Scan()
	err := reader.Err()
	if err != nil {
		return "", nil
	}

	return reader.Text(), nil
}
