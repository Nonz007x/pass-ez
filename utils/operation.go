package utils

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"
)

func EasyMode() {
	fmt.Printf("%s Ez mode %s\n", NAME, VER)

	for {
		fmt.Println()
		fmt.Println("(1) ")
		fmt.Println("(2) add a new password")
		fmt.Println("(3) encrypt a file")
		fmt.Println("(4) decrypt a file")
		fmt.Println("(5) get password")
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

		case "2":
			choice_AddPassword()

		case "3":
			choice_EncryptFile()

		case "4":
			choice_DecryptFile()

		case "5":
			choice_GetPassword()

		case "q":
			fmt.Println("Exiting program...")
			return
		default:
			fmt.Println("Invalid choice. Please try again.")
		}
	}
}

func choice_EncryptFile() {
	reader := bufio.NewReader(os.Stdin)

	fileName, err := getFileName(reader)
	if err != nil {
		fmt.Println(err)
		return
	}
	if fileName == "" {
		return
	}

	password, err := getSecret(true)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer WipeData(password)

	salt, err := GenerateSalt(DEFAULT_SALT_LEN)
	if err != nil {
		fmt.Println(err)
		return
	}

	key, err := DeriveKey(
		password,
		salt,
		DEFAULT_ITERATION,
		DEFAULT_KEY_LEN,
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

	header, err := CreateHeader(salt, DEFAULT_ITERATION, len(key), iv)
	if err != nil {
		fmt.Println(err)
		return
	}

	if err := EncryptFile(key, iv, header, fileName); err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("New file \"%s\" created successfully\n", fileName)
}

func choice_DecryptFile() {
	reader := bufio.NewReader(os.Stdin)

	fileName, err := getFileName(reader)
	if err != nil {
		fmt.Println(err)
		return
	}
	if fileName == "" {
		return
	}

	password, err := getSecret(false)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer WipeData(password)

	file, err := os.OpenFile(fileName, os.O_RDWR, 0644)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		fmt.Printf("error reading file \"%s\": %v\n", fileName, err)
		return
	}

	iteration, keyLength, salt, iv, _, err := ParseHeader(data)
	if err != nil {
		fmt.Printf("Error decrypting file: %v\n", err)
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

	file.Close()
	if err := DecryptFile(key, iv, fileName); err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("File \"%s\" decrypted successfully\n", fileName)
}

func choice_AddPassword() {
	reader := bufio.NewReader(os.Stdin)

	fileName, err := getFileName(reader)
	if err != nil {
		fmt.Printf("Error getting file name: %v\n", err)
		return
	}

	file, err := os.OpenFile(fileName, os.O_RDWR, 0644)
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Printf("File \"%s\" does not exist\n", fileName)
			return
		} else {
			fmt.Printf("Error opening file \"%s\": %v\n", fileName, err)
			return
		}
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		fmt.Printf("Error reading file \"%s\": %v\n", fileName, err)
		return
	}
	defer WipeData(data)

	var (
		key []byte
		iv []byte
		salt []byte
		iteration int
		keyLength int
		encryptedData []byte
		decryptedData []byte
	)

	encyptedFlag := IsEncrypted(data[:magicNumbersLen])
	
	if encyptedFlag {
		fmt.Println("File is encrypted. Please enter password.")
		
		password, err := getSecret(false)
		if err != nil {
			fmt.Println(err)
			return
		}
		defer WipeData(password)
		
		iteration, keyLength, salt, iv, encryptedData, err = ParseHeader(data)
		if err != nil {
			fmt.Printf("Error decrypting file: %v\n", err)
			return
		}
		
		key, err = DeriveKey(
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

		decryptedData, err = DecryptAES(key, iv, encryptedData)
		if err != nil {
			fmt.Println(err)
			return
		}
		defer WipeData(decryptedData)
	} else {
		decryptedData = data
		fmt.Println("WARNING: File is unencrypted.")
	}

	serviceName, err := getInput(reader, "Enter service: ")
	if err != nil {
		fmt.Printf("Error getting service name: %v\n", err)
		return
	}

	identifier, err := getInput(reader, "Enter identifier: ")
	if err != nil {
		fmt.Printf("Error getting identifier: %v\n", err)
		return
	}

	password, err := getSecret(true)
	if err != nil {
		fmt.Printf("Error getting password: %v\n", err)
		return
	}
	defer WipeData(password)

	var services []Service
	if len(decryptedData) > 0 {
		err = json.Unmarshal(decryptedData, &services)
		if err != nil {
			fmt.Printf("Error unmarshaling data from file \"%s\": %v\n", fileName, err)
			return
		}
	}

	var foundService *Service
	for i, service := range services {
		if service.ServiceName == serviceName {
			foundService = &services[i]
			break
		}
	}

	newCredential := Credential{
		Identifier: identifier,
		Password:   password,
	}
	
	if foundService != nil {
		foundService.Credentials = append(foundService.Credentials, newCredential)
		} else {
			newService := Service{
				ServiceName: serviceName,
				Credentials: []Credential{newCredential},
		}
		services = append(services, newService)
	}
	
	updatedData, err := json.Marshal(services)
	if err != nil {
		fmt.Printf("Error marshaling data: %v\n", err)
		return
	}

	if encyptedFlag {
		header, err := CreateHeader(salt, iteration, keyLength, iv)
		if err != nil {
			fmt.Println(err)
			return
		}
		updatedData, err = EncryptAES(key, iv, header, string(updatedData))
		if err != nil {
			fmt.Println(err)
			return
		}
	}
	
	err = file.Truncate(0)
	if err != nil {
		fmt.Printf("Error truncating file \"%s\": %v\n", fileName, err)
		return
	}

	_, err = file.Seek(0, io.SeekStart)
	if err != nil {
		fmt.Printf("Error seeking start of file \"%s\": %v\n", fileName, err)
		return
	}

	_, err = file.Write(updatedData)
	if err != nil {
		fmt.Printf("Error writing data to file \"%s\": %v\n", fileName, err)
		return
	}
	
	
	fmt.Printf("New credential added successfully to file \"%s\".\n", fileName)
}

func choice_GetPassword() {
	reader := bufio.NewReader(os.Stdin)
	
	fileName, err := getFileName(reader)
	if err != nil {
		fmt.Printf("Error getting file name: %v\n", err)
		return
	}

	file, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Printf("File \"%s\" does not exist, creating a new file...\n", fileName)
		} else {
			fmt.Printf("Error opening file \"%s\": %v\n", fileName, err)
			return
		}
	}
	defer file.Close()

	serviceName, err := getInput(reader, "Enter service: ")
	if err != nil {
		fmt.Printf("Error getting service name: %v\n", err)
		return
	}

	identifier, err := getInput(reader, "Enter identifier: ")
	if err != nil {
		fmt.Printf("Error getting identifier: %v\n", err)
		return
	}

	password, err := getSecret(false)
	if err != nil {
		fmt.Printf("Error getting password: %v\n", err)
		return
	}
	defer WipeData(password)

	data, err := io.ReadAll(file)
	if err != nil {
		fmt.Printf("Error reading file \"%s\": %v\n", fileName, err)
		return
	}
	defer WipeData(data)

	var services []Service
	if len(data) > 0 {
		err = json.Unmarshal(data, &services)
		if err != nil {
			fmt.Printf("Error unmarshaling data from file \"%s\": %v\n", fileName, err)
			return
		}
	}
	defer clearServiceData(services)

	var foundService *Service
	for i, service := range services {
		if service.ServiceName == serviceName {
			foundService = &services[i]
			break
		}
	}
	defer func() {
		foundService = nil
	}()

	if foundService == nil {
		fmt.Printf("Service \"%s\" not found. Please try again.\n", serviceName)
		return
	}

	var foundCredential *Credential
	for i, credential := range foundService.Credentials {
		if credential.Identifier == identifier {
			foundCredential = &foundService.Credentials[i]
			break
		}
	}
	defer func() {
		foundCredential = nil
	}()

	if foundCredential != nil {
		fmt.Println(string(foundCredential.Password))
	} else {
		fmt.Printf("Identifier \"%s\" not found. Please try again.\n", identifier)
	}
}

func clearServiceData(services []Service) {

	for i := range services {
		// Clear each service's credentials
		for j := range services[i].Credentials {
			// Clear each credential's password
			for k := range services[i].Credentials[j].Password {
				services[i].Credentials[j].Password[k] = 0
			}
			// Clear the credential's identifier (if it contains sensitive data)
			services[i].Credentials[j].Identifier = ""
			// Reset the credential to its zero value
			services[i].Credentials[j] = Credential{}
		}
		// Clear the service's name (if it contains sensitive data)
		services[i].ServiceName = ""
		// Reset the service to its zero value
		services[i] = Service{}
	}

	// Clear the entire slice
	services = nil
}
