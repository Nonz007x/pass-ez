package utils

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/crypto/pbkdf2"
)

const (
	headerLenField  = 1
	saltLenField    = 1
	keyLenField     = 1
	iterationSize   = 4
	minHeaderLen    = 55
	maxHeaderLen    = 87
	minSaltLen      = 32
	maxSaltLen      = 64
	magicNumbersLen = 4
)

func magicNumbers() []byte {
	return []byte{69, 115, 97, 110}
}

// sensitive
func DeriveKey(password []byte, salt []byte, iteration int, keyLength int) ([]byte, error) {
	// will be configurable
	keyByte := pbkdf2.Key(password, salt, iteration, keyLength, sha256.New)

	return keyByte, nil
}

// The length must be between 32 and 64 bytes.
func GenerateSalt(length int) ([]byte, error) {
	if length < 32 || length > 64 {
		return nil, errors.New("length must be between 32 and 64 bytes")
	}

	salt := make([]byte, length)

	if _, err := rand.Read(salt); err != nil {
		return nil, fmt.Errorf("failed to generate salt: %w", err)
	}

	return salt, nil
}

func GenerateIV() ([]byte, error) {
	iv := make([]byte, aes.BlockSize)
	if _, err := rand.Read(iv); err != nil {
		return nil, err
	}
	return iv, nil
}

// sensitive
func EncryptFile(key []byte, iv []byte, header []byte, fullFileName string) error {

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

	encryptedData, err := EncryptAES(key, iv, header, string(data))
	if err != nil {
		return fmt.Errorf("error encrypting data: %v", err)
	}

	err = file.Truncate(0)
	if err != nil {
		return fmt.Errorf("error truncating file \"%s\": %v", fullFileName, err)
	}

	_, err = file.Seek(0, io.SeekStart)
	if err != nil {
		return fmt.Errorf("error seeking to start of file \"%s\": %v", fullFileName, err)
	}

	_, err = file.Write(encryptedData)
	if err != nil {
		return fmt.Errorf("error writing data to file \"%s\": %v", fullFileName, err)
	}

	file.Close()

	newFileName := fullFileName + ENC_FILE_EXT
	err = os.Rename(fullFileName, newFileName)
	if err != nil {
		return fmt.Errorf("error renaming file \"%s\" to \"%s\": %v", fullFileName, newFileName, err)
	} else {
		fmt.Println(newFileName)
	}
	return nil
}

// sensitive
func DecryptFile(key []byte, iv []byte, fullFileName string) error {
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

	encryptedData, err := RemoveHeader(data)
	if err != nil {
		return err
	}

	decryptedData, err := DecryptAES(key, iv, encryptedData)
	if err != nil {
		return err
	}
	defer WipeData(decryptedData)

	err = file.Truncate(0)
	if err != nil {
		return fmt.Errorf("error truncating file \"%s\": %v", fullFileName, err)
	}

	_, err = file.Seek(0, io.SeekStart)
	if err != nil {
		return fmt.Errorf("error seeking to start of file \"%s\": %v", fullFileName, err)
	}

	_, err = file.Write(decryptedData)
	if err != nil {
		return fmt.Errorf("error writing data to file \"%s\": %v", fullFileName, err)
	}

	file.Close()

	decryptedFileName := strings.TrimSuffix(fullFileName, filepath.Ext(fullFileName))
	err = os.Rename(fullFileName, decryptedFileName)
	if err != nil {
		return fmt.Errorf("error renaming file \"%s\" to \"%s\": %v", fullFileName, decryptedFileName, err)
	}

	return nil
}

// sensitive
func EncryptAES(key []byte, iv []byte, header []byte, plaintext string) ([]byte, error) {
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

	combined := append(header, cipherText...)
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

	if len(iv) != block.BlockSize() {
		return nil, fmt.Errorf("invalid IV length. Expected %d, got %d", block.BlockSize(), len(iv))
	}

	// defer func() {
	// 	if r := recover(); r != nil {
	// 		fmt.Printf("Recovered from panic: %v\n", r)

	// 	}
	// }()
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
		return nil, errors.New("incorrect password")
	}

	for i := len(data) - paddingLength; i < len(data); i++ {
		if data[i] != byte(paddingLength) {
			return nil, fmt.Errorf("incorrect password")
		}
	}

	return data[:len(data)-paddingLength], nil
}

func CreateHeader(salt []byte, iteration int, keyLength int, iv []byte) ([]byte, error) {

	saltLength := len(salt)
	ivLength := len(iv)
	headerLength := headerLenField + saltLenField + keyLenField +
		iterationSize + saltLength + ivLength

	if headerLength > 255 {
		return nil, errors.New("header length exceeds maximum allowed size")
	}

	header := make([]byte, headerLength)

	header[0] = byte(headerLength)
	header[1] = byte(saltLength)
	header[2] = byte(keyLength)

	binary.BigEndian.PutUint32(header[3:7], uint32(iteration))

	offsetSalt := headerLenField + saltLenField + keyLenField + iterationSize
	offsetIV := offsetSalt + saltLength

	copy(header[offsetSalt:offsetSalt+saltLength], salt)
	copy(header[offsetIV:offsetIV+ivLength], iv)

	magicNumbers := magicNumbers()
	return append(magicNumbers, header...), nil
}

func ParseHeader(cipherText []byte) (int, int, []byte, []byte, []byte, error) {

	if !IsEncrypted(cipherText) {
		return 0, 0, nil, nil, nil, errors.New("file is unencrypted")
	}

	if len(cipherText) < magicNumbersLen {
		return 0, 0, nil, nil, nil, errors.New("invalid header length")
	}

	headerLength := int(cipherText[magicNumbersLen])

	if headerLength < minHeaderLen || headerLength > maxHeaderLen {
		return 0, 0, nil, nil, nil, errors.New("invalid header length")
	}

	header := cipherText[magicNumbersLen : magicNumbersLen+headerLength]

	saltLength := int(header[1])
	keyLength := int(header[2])

	if saltLength < minSaltLen || saltLength > maxSaltLen {
		return 0, 0, nil, nil, nil, errors.New("invalid salt length")
	}

	iteration := int(binary.BigEndian.Uint32(header[3:7]))

	offsetSalt := headerLenField + saltLenField + keyLenField + iterationSize
	offsetIV := offsetSalt + saltLength

	salt := header[offsetSalt : offsetSalt+saltLength]
	iv := header[offsetIV : offsetIV+(headerLength-offsetIV)]

	encryptedData := cipherText[magicNumbersLen+headerLength:]

	return iteration, keyLength, salt, iv, encryptedData, nil
}

func IsEncrypted(data []byte) bool {
	return bytes.Equal(data[:magicNumbersLen], magicNumbers())
}

func RemoveHeader(cipherText []byte) ([]byte, error) {
	if len(cipherText) <= magicNumbersLen {
		return nil, fmt.Errorf("invalid cipherText: too short to contain header and magic numbers")
	}

	headerLength := int(cipherText[magicNumbersLen])

	if len(cipherText) <= magicNumbersLen+headerLength {
		return nil, fmt.Errorf("invalid header length: exceeds cipherText length")
	}

	encryptedData := cipherText[magicNumbersLen+headerLength:]

	return encryptedData, nil
}
