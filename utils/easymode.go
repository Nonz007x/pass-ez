package utils

import (
	"bufio"
	"fmt"
	"os"

	consts "github.com/spf13/myapp/constants"
)

func EasyMode() {
	fmt.Printf("%s Ez mode %s\n", consts.NAME, consts.VER)

	for {
		fmt.Println()
		fmt.Println("(1) create new file")
		fmt.Println("(2) add a new password")
		fmt.Println("(3) delete an existing password")
		fmt.Println("(q) exit program")
		var choice string
		fmt.Print(": ")

		_, err := fmt.Scanln(&choice)
		if err != nil {
			fmt.Println("Error reading input. Please try again.")
			continue
		}

		switch choice {
		case "1":
			createFile()

		case "2":
			reader := bufio.NewReader(os.Stdin)

			fileName, err := getFileName(reader)
			if err != nil {
				fmt.Println(err)
				continue
			}

			service, err := getInput(reader, "Enter service: ")
			if err != nil {
				fmt.Println(err)
				return
			}
			username, err := getInput(reader, "Enter username: ")
			if err != nil {
				fmt.Println(err)
				return
			}
			password, err := getPassword()
			if err != nil {
				fmt.Println(err)
				return
			}

			if err := AddPassword(fileName, service, username, password); err != nil {
				fmt.Println(err)
				continue
			}

			WipeData(password)

			fmt.Printf("New credential added successfully to file \"%s\".\n", fileName)

		case "3":


		case "4":
			decryptFile()

		case "q":
			fmt.Println("Exiting program...")
			return
		default:
			fmt.Println("Invalid choice. Please try again.")
		}
	}
}

func createFile() {
	reader := bufio.NewReader(os.Stdin)

	fileName, err := getFileName(reader)
	if err != nil {
		fmt.Println(err)
		return
	}
	if fileName == "" {
		return
	}

	password, err := getPassword()
	if err != nil {
		fmt.Println(err)
		return
	}
	defer WipeData(password)

	salt, err := GenerateSalt(consts.DEFAULT_SALT_LEN)
	if err != nil {
		fmt.Println(err)
		return
	}

	key, err := DeriveKey(
		password,
		salt,
		consts.DEFAULT_ITERATION,
		consts.DEFAULT_KEY_LEN,
	)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer WipeData(key)

	if err := CreateFileAndEncrypt(key, fileName); err != nil {
		fmt.Println(err)
		return
	}
	SavePBKDF2Params(fileName+consts.PARAMS_FILE_EXT, salt, consts.DEFAULT_ITERATION, consts.DEFAULT_KEY_LEN)

	fmt.Printf("New file \"%s\" created successfully\n", fileName)
}

func decryptFile() {
	reader := bufio.NewReader(os.Stdin)

	fileName, err := getFileName(reader)
	if err != nil {
		fmt.Println(err)
		return
	}
	if fileName == "" {
		return
	}

	password, err := getPassword()
	if err != nil {
		fmt.Println(err)
		return
	}
	defer WipeData(password)

	salt, iteration, keyLength, err := RetrievePBKDF2Params(fileName + consts.PARAMS_FILE_EXT)
	if err != nil {
		fmt.Println(err)
		return
	}

	key, err := DeriveKey(
		password,
		salt,
		iteration,
		keyLength,
	)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer WipeData(key)

	if err := DecryptFile(key, fileName+consts.ENC_FILE_EXT); err != nil {
		fmt.Println(err)
		return
	}
}
