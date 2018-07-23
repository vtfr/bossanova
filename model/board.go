package model

import (
	"time"
)

// Board is a board
type Board struct {
	URI         string    `json:"uri" bson:"_id" validate:"board"`
	Name        string    `json:"title" validate:"required"`
	Description string    `json:"description,omitempty"`
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

// Valid returns an error if the board is invalid
func (board *Board) Valid() error {
	return Validate(board)
}
