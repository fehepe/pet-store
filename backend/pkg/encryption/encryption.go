package encryption

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"io"
)

// EncryptorInterface defines the interface for encryption operations
type EncryptorInterface interface {
	// Encrypt encrypts plaintext and returns ciphertext
	Encrypt(plaintext string) (string, error)

	// Decrypt decrypts ciphertext and returns plaintext
	Decrypt(ciphertext string) (string, error)
}

// Ensure Encryptor implements EncryptorInterface
var _ EncryptorInterface = (*Encryptor)(nil)

type Encryptor struct {
	key []byte
}

func NewEncryptor(key string) (*Encryptor, error) {
	keyBytes := []byte(key)
	if len(keyBytes) != 32 {
		return nil, errors.New("encryption key must be 32 bytes long")
	}
	return &Encryptor{key: keyBytes}, nil
}

func (e *Encryptor) Encrypt(plaintext string) (string, error) {
	block, err := aes.NewCipher(e.key)
	if err != nil {
		return "", err
	}

	plaintextBytes := []byte(plaintext)
	ciphertext := make([]byte, aes.BlockSize+len(plaintextBytes))
	iv := ciphertext[:aes.BlockSize]

	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return "", err
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], plaintextBytes)

	return base64.URLEncoding.EncodeToString(ciphertext), nil
}

func (e *Encryptor) Decrypt(ciphertext string) (string, error) {
	ciphertextBytes, err := base64.URLEncoding.DecodeString(ciphertext)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(e.key)
	if err != nil {
		return "", err
	}

	if len(ciphertextBytes) < aes.BlockSize {
		return "", errors.New("ciphertext too short")
	}

	iv := ciphertextBytes[:aes.BlockSize]
	ciphertextBytes = ciphertextBytes[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(ciphertextBytes, ciphertextBytes)

	return string(ciphertextBytes), nil
}
