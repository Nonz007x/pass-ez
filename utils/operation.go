package utils

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	consts "github.com/spf13/myapp/constants"
)

// optimization is needed
func CreateFileAndEncrypt(key []byte, fileName string) error {
	fullFileName := fileName + consts.ENC_FILE_EXT
	if err := CreateFile(fullFileName); err != nil {
		return err
	}

	if err := EncryptFile(key, fullFileName); err != nil {
		return err
	}

	return nil
}

// sensitive
func AddPassword(fileName string, serviceName string, username string, password []byte) error {
	fullFileName := fileName + consts.ENC_FILE_EXT
	file, err := OpenExistingFile(fullFileName)
	if err != nil {
		return err
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		return fmt.Errorf("error reading file \"%s\": %v", fullFileName, err)
	}

	var services []consts.Service
	if len(data) > 0 {
		err = json.Unmarshal(data, &services)
		if err != nil {
			return fmt.Errorf("error unmarshaling data from file \"%s\": %v", fullFileName, err)
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
			Username: username,
			Password: password,
		}
		foundService.Credentials = append(foundService.Credentials, newCredential)
	} else {

		newService := consts.Service{
			ServiceName: serviceName,
			Credentials: []consts.Credential{
				{
					Username: username,
					Password: password,
				},
			},
		}
		services = append(services, newService)
	}

	updatedData, err := json.Marshal(services)
	if err != nil {
		return fmt.Errorf("error marshaling data: %v", err)
	}

	file, err = os.OpenFile(fullFileName, os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return fmt.Errorf("error opening file \"%s\" for writing: %v", fullFileName, err)
	}
	defer file.Close()

	_, err = file.Write(updatedData)
	if err != nil {
		return fmt.Errorf("error writing data to file \"%s\": %v", fullFileName, err)
	}

	return nil
}
