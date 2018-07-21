package model

import (
	"time"

	"golang.org/x/crypto/bcrypt"
)

// User is a user
type User struct {
	Username       string    `json:"username" bson:"_id" validate:"required,min=4,max=32"`
	Role           string    `json:"role" bson:"role" validate:"required"`
	HashedPassword []byte    `json:"-" bson:"hashedPassword" validate:"required"`
	CreatedAt      time.Time `json:"createdAt" bson:"createdAt"`
	LastModifiedAt time.Time `json:"modifiedAt" bson:"modifiedAt"`
}

// HashPassword hashes a password using bcrypt with default cost.
func HashPassword(password string) []byte {
	hash, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return hash
}

// SetPassword sets an user password after hashing it.
func (u *User) SetPassword(password string) {
	u.HashedPassword = HashPassword(password)
}

// VerifyPassword verifies if a user has the correct password
func (u *User) VerifyPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.HashedPassword), []byte(password))
	return err == nil
}
