package mocks

import (
	"sync"

	"github.com/stretchr/testify/mock"
)

// MockEncryptor is a mock implementation of encryption.EncryptorInterface
type MockEncryptor struct {
	mock.Mock
	mu sync.RWMutex
}

func (m *MockEncryptor) Encrypt(plaintext string) (string, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	ret := m.Called(plaintext)
	return ret.String(0), ret.Error(1)
}

func (m *MockEncryptor) Decrypt(ciphertext string) (string, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	ret := m.Called(ciphertext)
	return ret.String(0), ret.Error(1)
}
