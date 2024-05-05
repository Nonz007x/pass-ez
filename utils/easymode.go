package utils

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"encoding/json"

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
			addPassword()

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

	iv, err := GenerateIV()
	if err != nil {
		fmt.Println(err)
		return
	}

	header, err := CreateHeader(salt, consts.DEFAULT_ITERATION, len(key), iv)
	if err != nil {
		fmt.Println(err)
		return
	}

	EncryptFile(key, iv, header, fileName+consts.ENC_FILE_EXT)
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

	file, err := os.OpenFile(fileName+consts.ENC_FILE_EXT, os.O_RDWR, 0644)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		fmt.Printf("error reading file \"%s\": %v\n", fileName+consts.ENC_FILE_EXT, err)
		return
	}

	iteration, keyLength, salt, iv, _, err := ParseHeader(data)
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

	if err := DecryptFile(key, iv, fileName+consts.ENC_FILE_EXT); err != nil {
		fmt.Println(err)
		return
	}
}

func addPassword() {
	reader := bufio.NewReader(os.Stdin)

	fileName, err := getFileName(reader)
	if err != nil {
		fmt.Println(err)
		return
	}

	serviceName, err := getInput(reader, "Enter service: ")
	if err != nil {
		fmt.Println(err)
		return
	}
	identifier, err := getInput(reader, "Enter identifier: ")
	if err != nil {
		fmt.Println(err)
		return
	}
	password, err := getPassword()
	if err != nil {
		fmt.Println(err)
		return
	}
	defer WipeData(password)
	
	fullFileName := fileName + ".json"
	file, err := OpenExistingFile(fullFileName)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		fmt.Printf("error reading file \"%s\": %v", fullFileName, err)
		return
	}

	var services []consts.Service
	if len(data) > 0 {
		err = json.Unmarshal(data, &services)
		if err != nil {
			fmt.Printf("error unmarshaling data from file \"%s\": %v", fullFileName, err)
			return
		}
	}

	var foundService *consts.Service
	for i, service := range services {
		if service.ServiceName == serviceName {
			foundService = &services[i]
			break
		}
	}

	if foundService != nil {
		newCredential := consts.Credential{
			Identifier: identifier,
			Password:   password,
		}
		foundService.Credentials = append(foundService.Credentials, newCredential)
	} else {

		newService := consts.Service{
			ServiceName: serviceName,
			Credentials: []consts.Credential{
				{
					Identifier: identifier,
					Password:   password,
				},
			},
		}
		services = append(services, newService)
	}

	updatedData, err := json.Marshal(services)
	if err != nil {
		fmt.Printf("error marshaling data: %v", err)
		return 
	}

	file, err = os.OpenFile(fullFileName, os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		fmt.Printf("error opening file \"%s\" for writing: %v", fullFileName, err)
		return 
	}
	defer file.Close()

	_, err = file.Write(updatedData)
	if err != nil {
		fmt.Printf("error writing data to file \"%s\": %v", fullFileName, err)
		return 
	}

	fmt.Printf("New credential added successfully to file \"%s\".\n", fullFileName)
}
