package mocks

import (
	"context"
	"database/sql"
	"sync"

	"github.com/stretchr/testify/mock"
)

// MockRepository is a mock implementation of database.Repository
type MockRepository struct {
	mock.Mock
	mu sync.RWMutex
}

func (m *MockRepository) QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row {
	m.mu.RLock()
	defer m.mu.RUnlock()
	argsList := make([]interface{}, len(args)+2)
	argsList[0] = ctx
	argsList[1] = query
	copy(argsList[2:], args)
	ret := m.Called(argsList...)
	return ret.Get(0).(*sql.Row)
}

func (m *MockRepository) QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	argsList := make([]interface{}, len(args)+2)
	argsList[0] = ctx
	argsList[1] = query
	copy(argsList[2:], args)
	ret := m.Called(argsList...)
	return ret.Get(0).(*sql.Rows), ret.Error(1)
}

func (m *MockRepository) ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	argsList := make([]interface{}, len(args)+2)
	argsList[0] = ctx
	argsList[1] = query
	copy(argsList[2:], args)
	ret := m.Called(argsList...)
	return ret.Get(0).(sql.Result), ret.Error(1)
}

func (m *MockRepository) Transaction(fn func(*sql.Tx) error) error {
	m.mu.RLock()
	defer m.mu.RUnlock()
	ret := m.Called(fn)
	return ret.Error(0)
}

func (m *MockRepository) Close() error {
	m.mu.RLock()
	defer m.mu.RUnlock()
	ret := m.Called()
	return ret.Error(0)
}

func (m *MockRepository) Ping() error {
	m.mu.RLock()
	defer m.mu.RUnlock()
	ret := m.Called()
	return ret.Error(0)
}

func (m *MockRepository) Migrate() error {
	m.mu.RLock()
	defer m.mu.RUnlock()
	ret := m.Called()
	return ret.Error(0)
}

// MockResult is a mock implementation of sql.Result
type MockResult struct {
	mock.Mock
}

func (m *MockResult) LastInsertId() (int64, error) {
	ret := m.Called()
	return ret.Get(0).(int64), ret.Error(1)
}

func (m *MockResult) RowsAffected() (int64, error) {
	ret := m.Called()
	return ret.Get(0).(int64), ret.Error(1)
}
