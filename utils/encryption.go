package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"errors"
	"fmt"
	"io"
	"os"

	"golang.org/x/crypto/pbkdf2"
)

// sensitive
func DeriveKey(password []byte, salt []byte, iteration int, keyLength int) ([]byte, error) {
	// will be configurable
	keyByte := pbkdf2.Key(password, salt, iteration, keyLength, sha256.New)

	return keyByte, nil
}

func GenerateSalt(length int) ([]byte, error) {
	salt := make([]byte, length)

	_, err := rand.Read(salt)
	if err != nil {
		return nil, err
	}

	return salt, nil
}

// sensitive
func EncryptFile(key []byte, fullFileName string) error {

	file, err := OpenExistingFile(fullFileName)
	if err != nil {
		return err
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		return fmt.Errorf("error reading file \"%s\": %v", fullFileName, err)
	}

	// Generate a random IV
	iv := make([]byte, aes.BlockSize)
	if _, err := rand.Read(iv); err != nil {
		return err
	}

	encryptedData, err := EncryptAES(key, iv, string(data))
	if err != nil {
		return err
	}

	file, err = os.OpenFile(fullFileName, os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return fmt.Errorf("error opening file \"%s\" for writing: %v", fullFileName, err)
	}
	defer file.Close()

	_, err = file.Write(encryptedData)
	if err != nil {
		return fmt.Errorf("error writing data to file \"%s\": %v", fullFileName, err)
	}

	fmt.Println(encryptedData)

	return nil
}

// sensitive
func DecryptFile(key []byte, fullFileName string) error {
	file, err := os.OpenFile(fullFileName, os.O_RDWR, 0644)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("error: file \"%s\" does not exist", fullFileName)
		}
		return fmt.Errorf("error opening file \"%s\": %v", fullFileName, err)
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		return fmt.Errorf("error reading file \"%s\": %v", fullFileName, err)
	}

	if len(data) < 16 {
		return errors.New("ciphertext is too short")
	}

	iv := data[:aes.BlockSize]
	encryptedData := data[aes.BlockSize:]

	decryptedData, err := DecryptAES(key, iv, encryptedData)
	if err != nil {
		return err
	}
	defer WipeData(decryptedData)
	fmt.Println(decryptedData)

	file, err = os.OpenFile(fullFileName, os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return fmt.Errorf("error opening file \"%s\" for writing: %v", fullFileName, err)
	}

	_, err = file.Write(decryptedData)
	if err != nil {
		return fmt.Errorf("error writing data to file \"%s\": %v", fullFileName, err)
	}

	return nil
}

// sensitive
func EncryptAES(key []byte, iv []byte, plaintext string) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	plaintextBytes := []byte(plaintext)
	plaintextPadded := PadPKCS7(plaintextBytes, aes.BlockSize)
	defer WipeData(plaintextBytes)
	defer WipeData(plaintextPadded)

	cipherText := make([]byte, len(plaintextPadded))

	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(cipherText, plaintextPadded)

	combined := append(iv, cipherText...)
	return combined, nil
}

// sensitive
func DecryptAES(key []byte, iv []byte, ciphertextBytes []byte) ([]byte, error) {

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	if len(ciphertextBytes) < aes.BlockSize {
		return nil, errors.New("ciphertext is too short")
	}

	mode := cipher.NewCBCDecrypter(block, iv)

	mode.CryptBlocks(ciphertextBytes, ciphertextBytes)

	plaintext, err := UnpadPKCS7(ciphertextBytes)
	if err != nil {
		return nil, err
	}

	return plaintext, nil
}

func PadPKCS7(data []byte, blockSize int) []byte {
	padding := blockSize - len(data)%blockSize
	padText := make([]byte, padding)
	for i := range padText {
		padText[i] = byte(padding)
	}
	return append(data, padText...)
}

func UnpadPKCS7(data []byte) ([]byte, error) {
	if len(data) == 0 {
		return nil, errors.New("input data is empty")
	}

	paddingLength := int(data[len(data)-1])
	if paddingLength <= 0 || paddingLength > len(data) {
		return nil, errors.New("invalid padding length")
	}

	for i := len(data) - paddingLength; i < len(data); i++ {
		if data[i] != byte(paddingLength) {
			return nil, fmt.Errorf("invalid padding value")
		}
	}

	return data[:len(data)-paddingLength], nil
}
