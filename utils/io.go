package utils

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"os"
	"strings"

	"golang.org/x/term"
)

const INT_SIZE int = 4

func getFileName(reader *bufio.Reader) (string, error) {
	for {
		fmt.Print("Enter file name witout extension (type 'cancel' to exit): ")

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

func SavePBKDF2Params(fullFileName string, salt []byte, iteration int, keyLength int) error {
	file, err := os.OpenFile(fullFileName, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("error creating file \"%s\": %v", fullFileName, err)
	}
	defer file.Close()

	buf := make([]byte, len(salt)+INT_SIZE*2)

	copy(buf, salt)

	binary.BigEndian.PutUint32(buf[len(salt):len(salt)+INT_SIZE], uint32(iteration))

	binary.BigEndian.PutUint32(buf[len(salt)+INT_SIZE:], uint32(keyLength))

	_, err = file.Write(buf)
	if err != nil {
		return fmt.Errorf("error writing data to file \"%s\": %v", fullFileName, err)
	}
	return nil
}

func RetrievePBKDF2Params(fullFileName string) ([]byte, int, int, error) {
	file, err := os.OpenFile(fullFileName, os.O_RDONLY, 0444)
	if err != nil {
		return nil, 0, 0, fmt.Errorf("error opening file \"%s\": %v", fullFileName, err)
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		return nil, 0, 0, fmt.Errorf("error reading file \"%s\": %v", fullFileName, err)
	}

	if len(data) < INT_SIZE * 2 {
		return nil, 0, 0, fmt.Errorf("invalid file data length")
	}

	saltLength := len(data) - INT_SIZE * 2
	salt := data[:saltLength]

	iterations := binary.BigEndian.Uint32(data[saltLength : saltLength + INT_SIZE])

	keyLength := binary.BigEndian.Uint32(data[saltLength + INT_SIZE:])

	return salt, int(iterations), int(keyLength), nil
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
