package utils

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"strings"

	"golang.org/x/term"
)

const INT_SIZE int = 4

func getFileName(reader *bufio.Reader) (string, error) {
	for {
		fmt.Print("Enter file name without extension (type 'cancel' to exit): ")

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
func getPassword() ([]byte, error) {
	fmt.Print("Enter password: ")

	fd := int(os.Stdin.Fd())

	password, err := term.ReadPassword(fd)
	if err != nil {
		return nil, fmt.Errorf("error reading password: %v", err)
	}

	password = bytes.TrimSpace(password)
	fmt.Println()

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
