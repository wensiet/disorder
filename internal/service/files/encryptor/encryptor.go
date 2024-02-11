package encryptor

import (
	"crypto/aes"
	"crypto/cipher"
	crand "crypto/rand"
	"errors"
	"io"
	"os"
)

func KeyGen(keySize int) ([]byte, error) {
	key := make([]byte, keySize)
	_, err := io.ReadFull(crand.Reader, key)
	if err != nil {
		return nil, err
	}
	return key, nil
}

// SaveKeyToFile saves the given key to a file.
func SaveKeyToFile(key []byte, filename string) error {
	return os.WriteFile(filename, key, 0644)
}

// LoadKeyFromFile loads the key from a file.
func LoadKeyFromFile(filename string) ([]byte, error) {
	return os.ReadFile(filename)
}

func EncryptData(data []byte, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	// The IV needs to be unique, but doesn't have to be secret.
	iv := make([]byte, aes.BlockSize)
	if _, err := io.ReadFull(crand.Reader, iv); err != nil {
		return nil, err
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	ciphertext := make([]byte, len(data))
	stream.XORKeyStream(ciphertext, data)

	// Prepend the IV to the ciphertext.
	ciphertext = append(iv, ciphertext...)
	return ciphertext, nil
}

func DecryptData(encryptedData []byte, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	if len(encryptedData) < aes.BlockSize {
		return nil, errors.New("encrypted data is too short")
	}
	iv := encryptedData[:aes.BlockSize]
	ciphertext := encryptedData[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)

	// XORKeyStream can work in-place if the two arguments are the same.
	stream.XORKeyStream(ciphertext, ciphertext)

	return ciphertext, nil
}
