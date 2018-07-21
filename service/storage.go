package service

import (
	"io"

	"github.com/satori/go.uuid"
)

// GenerateMediaID generates a new media ID
func GenerateStorageID() string {
	return uuid.Must(uuid.NewV4()).String()
}

// Storage stores files somewhere
type Storage interface {
	// Upload uploads a new file
	Upload(io.Reader) (string, error)
	// Delete deletes a file
	Delete(id string) error
}
