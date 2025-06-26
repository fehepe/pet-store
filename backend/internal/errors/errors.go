package errors

import (
	"fmt"

	"github.com/google/uuid"
)

// Custom error types for better error handling

// NotFoundError represents when a resource is not found
type NotFoundError struct {
	Resource string
	ID       string
}

func (e NotFoundError) Error() string {
	return fmt.Sprintf("%s with ID %s not found", e.Resource, e.ID)
}

// ValidationError represents validation failures
type ValidationError struct {
	Field   string
	Message string
}

func (e ValidationError) Error() string {
	return fmt.Sprintf("validation error for %s: %s", e.Field, e.Message)
}

// ConflictError represents resource conflicts
type ConflictError struct {
	Resource string
	Message  string
}

func (e ConflictError) Error() string {
	return fmt.Sprintf("conflict with %s: %s", e.Resource, e.Message)
}

// PetNotFoundError specific error for pet operations
type PetNotFoundError struct {
	PetID uuid.UUID
}

func (e PetNotFoundError) Error() string {
	return fmt.Sprintf("pet with ID %s not found", e.PetID)
}

// StoreNotFoundError specific error for store operations
type StoreNotFoundError struct {
	StoreID uuid.UUID
}

func (e StoreNotFoundError) Error() string {
	return fmt.Sprintf("store with ID %s not found", e.StoreID)
}

// Helper functions for creating common errors

func NewPetNotFound(petID uuid.UUID) error {
	return PetNotFoundError{PetID: petID}
}

func NewStoreNotFound(storeID uuid.UUID) error {
	return StoreNotFoundError{StoreID: storeID}
}

func NewValidationError(field, message string) error {
	return ValidationError{Field: field, Message: message}
}

// BusinessRuleError represents business logic violations
type BusinessRuleError struct {
	Message string
}

func (e BusinessRuleError) Error() string {
	return fmt.Sprintf("business rule violation: %s", e.Message)
}

// OrderNotFoundError specific error for order operations
type OrderNotFoundError struct {
	OrderID uuid.UUID
}

func (e OrderNotFoundError) Error() string {
	return fmt.Sprintf("order with ID %s not found", e.OrderID)
}

func NewBusinessRuleError(message string) error {
	return BusinessRuleError{Message: message}
}

func NewOrderNotFound(orderID uuid.UUID) error {
	return OrderNotFoundError{OrderID: orderID}
}
