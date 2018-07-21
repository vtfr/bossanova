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
