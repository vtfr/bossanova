package model

import (
	"time"
)

// Board is a board
type Board struct {
	URI         string    `json:"uri" bson:"_id" validate:"board"`
	Name        string    `json:"title" validate:"required,max=16"`
	Description string    `json:"description,omitempty" validate:"max=32"`
	CreatedAt   time.Time `json:"createdAt" bson:"createdAt"`
	Hidden      bool      `json:"hidden,omitempty" bson:"hidden"`
}

// NewBoard creates a new board
func NewBoard(uri, name, description string) *Board {
	return &Board{
		URI:         uri,
		Name:        name,
		Description: description,
		CreatedAt:   time.Now().UTC().Truncate(time.Second),
	}
}
