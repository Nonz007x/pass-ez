package utils

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"golang.org/x/term"
	"os"
	"strings"
)

const INT_SIZE int = 4

func getFileName(reader *bufio.Reader) (string, error) {
	for {
		fmt.Print("Enter file name (type 'cancel' to exit): ")

		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Error reading input. Please try again.")
			continue
		}

		fileName := strings.TrimSpace(input)
		if fileName == "cancel" {
			fmt.Println("Operation cancelled.")
			return "", nil
		}

		if fileName == "" {
			fmt.Println("File name cannot be empty. Please try again.")
		} else {
			return fileName, nil
		}
	}
}

func getInput(reader *bufio.Reader, prompt string) (string, error) {
	fmt.Print(prompt)

	input, err := reader.ReadString('\n')
	if err != nil {
		return "", fmt.Errorf("error reading input: %v", err)
	}
	return strings.TrimSpace(input), nil
}

// This function is sensitive because it handles sensitive data (password).
// It is recommended to securely destroy
// the password after use by overwriting it (use "WipeData" function),
// as password data should not remain in memory longer than necessary.
func getSecret(confirmation bool) ([]byte, error) {
	fmt.Print("Enter password: ")

	// Read the password securely (without echoing) from the standard input
	fd := int(os.Stdin.Fd())
	password, err := term.ReadPassword(fd)
	if err != nil {
		return nil, fmt.Errorf("error reading password: %v", err)
	}

	// Trim whitespace from the entered password
	password = bytes.TrimSpace(password)

	fmt.Println()

	if confirmation {

		fmt.Print("Confirm password: ")

		confirmPassword, err := term.ReadPassword(fd)
		if err != nil {
			return nil, fmt.Errorf("error reading confirmed password: %v", err)
		}

		confirmPassword = bytes.TrimSpace(confirmPassword)

		fmt.Println()

		if !bytes.Equal(password, confirmPassword) {
			return nil, errors.New("passwords do not match")
		}
	}

	return password, nil
}

func OpenExistingFile(fullFileName string) (*os.File, error) {
	file, err := os.OpenFile(fullFileName, os.O_RDWR, 0644)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("error: file \"%s\" does not exist", fullFileName)
		}
		return nil, fmt.Errorf("error opening file \"%s\": %v", fullFileName, err)
	}
	return file, nil
}

func CreateFile(fullFileName string) error {
	file, err := os.OpenFile(fullFileName, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0644)
	if err != nil {
		if os.IsExist(err) {
			return fmt.Errorf("file \"%s\" already exists", fullFileName)
		}

		return fmt.Errorf("error creating file \"%s\": %v", fullFileName, err)
	}

	defer file.Close()
	return nil
}
