// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package model

import (
	"bytes"
	"fmt"
	"io"
	"strconv"
	"time"

	"github.com/google/uuid"
)

type CreatePetInput struct {
	Name         string     `json:"name"`
	Species      PetSpecies `json:"species"`
	Age          int32      `json:"age"`
	PictureURL   *string    `json:"pictureUrl,omitempty"`
	Description  *string    `json:"description,omitempty"`
	BreederName  string     `json:"breederName"`
	BreederEmail string     `json:"breederEmail"`
}

type CreateStoreInput struct {
	Name string `json:"name"`
}

type Mutation struct {
}

type Order struct {
	ID         uuid.UUID `json:"id"`
	CustomerID string    `json:"customerID"`
	Pets       []*Pet    `json:"pets"`
	TotalPets  int32     `json:"totalPets"`
	CreatedAt  time.Time `json:"createdAt"`
}

type PageInfo struct {
	HasNextPage     bool    `json:"hasNextPage"`
	HasPreviousPage bool    `json:"hasPreviousPage"`
	StartCursor     *string `json:"startCursor,omitempty"`
	EndCursor       *string `json:"endCursor,omitempty"`
}

type PaginationInput struct {
	First  *int32  `json:"first,omitempty"`
	After  *string `json:"after,omitempty"`
	Last   *int32  `json:"last,omitempty"`
	Before *string `json:"before,omitempty"`
}

type Pet struct {
	ID           uuid.UUID  `json:"id"`
	Name         string     `json:"name"`
	Species      PetSpecies `json:"species"`
	Age          int32      `json:"age"`
	PictureURL   *string    `json:"pictureUrl,omitempty"`
	Description  *string    `json:"description,omitempty"`
	BreederName  string     `json:"breederName"`
	BreederEmail string     `json:"breederEmail"`
	Status       PetStatus  `json:"status"`
	CreatedAt    time.Time  `json:"createdAt"`
}

type PetConnection struct {
	Edges      []*Pet    `json:"edges"`
	PageInfo   *PageInfo `json:"pageInfo"`
	TotalCount int32     `json:"totalCount"`
}

type PetFilterInput struct {
	Status    *PetStatus `json:"status,omitempty"`
	StartDate *time.Time `json:"startDate,omitempty"`
	EndDate   *time.Time `json:"endDate,omitempty"`
}

type Query struct {
}

type Store struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"createdAt"`
}

type PetSpecies string

const (
	PetSpeciesCat  PetSpecies = "Cat"
	PetSpeciesDog  PetSpecies = "Dog"
	PetSpeciesFrog PetSpecies = "Frog"
)

var AllPetSpecies = []PetSpecies{
	PetSpeciesCat,
	PetSpeciesDog,
	PetSpeciesFrog,
}

func (e PetSpecies) IsValid() bool {
	switch e {
	case PetSpeciesCat, PetSpeciesDog, PetSpeciesFrog:
		return true
	}
	return false
}

func (e PetSpecies) String() string {
	return string(e)
}

func (e *PetSpecies) UnmarshalGQL(v any) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = PetSpecies(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid PetSpecies", str)
	}
	return nil
}

func (e PetSpecies) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}

func (e *PetSpecies) UnmarshalJSON(b []byte) error {
	s, err := strconv.Unquote(string(b))
	if err != nil {
		return err
	}
	return e.UnmarshalGQL(s)
}

func (e PetSpecies) MarshalJSON() ([]byte, error) {
	var buf bytes.Buffer
	e.MarshalGQL(&buf)
	return buf.Bytes(), nil
}

type PetStatus string

const (
	PetStatusAvailable PetStatus = "available"
	PetStatusSold      PetStatus = "sold"
)

var AllPetStatus = []PetStatus{
	PetStatusAvailable,
	PetStatusSold,
}

func (e PetStatus) IsValid() bool {
	switch e {
	case PetStatusAvailable, PetStatusSold:
		return true
	}
	return false
}

func (e PetStatus) String() string {
	return string(e)
}

func (e *PetStatus) UnmarshalGQL(v any) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = PetStatus(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid PetStatus", str)
	}
	return nil
}

func (e PetStatus) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}

func (e *PetStatus) UnmarshalJSON(b []byte) error {
	s, err := strconv.Unquote(string(b))
	if err != nil {
		return err
	}
	return e.UnmarshalGQL(s)
}

func (e PetStatus) MarshalJSON() ([]byte, error) {
	var buf bytes.Buffer
	e.MarshalGQL(&buf)
	return buf.Bytes(), nil
}
