package input

import (
	"bufio"
	"fmt"
	"os"

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

func readInput() (string, error) {
	reader := bufio.NewScanner(os.Stdin)
	reader.Scan()
	err := reader.Err()
	if err != nil {
		return "", nil
	}

	return reader.Text(), nil
}
