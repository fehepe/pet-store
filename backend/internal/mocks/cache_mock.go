package mocks

import (
	"context"
	"sync"
	"time"

	"github.com/stretchr/testify/mock"
)

// MockCache is a mock implementation of cache.CacheInterface
type MockCache struct {
	mock.Mock
	mu sync.RWMutex
}

func (m *MockCache) Get(ctx context.Context, key string, dest interface{}) error {
	m.mu.RLock()
	defer m.mu.RUnlock()
	ret := m.Called(ctx, key, dest)
	return ret.Error(0)
}

func (m *MockCache) Set(ctx context.Context, key string, value interface{}, ttl ...time.Duration) error {
	m.mu.RLock()
	defer m.mu.RUnlock()
	args := []interface{}{ctx, key, value}
	for _, t := range ttl {
		args = append(args, t)
	}
	ret := m.Called(args...)
	return ret.Error(0)
}

func (m *MockCache) Delete(ctx context.Context, keys ...string) error {
	m.mu.RLock()
	defer m.mu.RUnlock()
	args := []interface{}{ctx}
	for _, k := range keys {
		args = append(args, k)
	}
	ret := m.Called(args...)
	return ret.Error(0)
}

func (m *MockCache) InvalidatePattern(ctx context.Context, pattern string) error {
	m.mu.RLock()
	defer m.mu.RUnlock()
	ret := m.Called(ctx, pattern)
	return ret.Error(0)
}

func (m *MockCache) Close() error {
	m.mu.RLock()
	defer m.mu.RUnlock()
	ret := m.Called()
	return ret.Error(0)
}

func (m *MockCache) Ping() error {
	m.mu.RLock()
	defer m.mu.RUnlock()
	ret := m.Called()
	return ret.Error(0)
}
