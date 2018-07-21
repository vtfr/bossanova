package service

import (
	"io"

	"github.com/spf13/afero"
)

// FilesystemStorage stores medias in the filesystem
type FilesystemStorage struct {
	fs afero.Fs
}

// NewFilesystemStorage creates a new filesystem storage
func NewFilesystemStorage(path string) *FilesystemStorage {
	return &FilesystemStorage{
		afero.NewBasePathFs(afero.NewOsFs(), path),
	}
}

// Upload uploads a file to the filesystem and return it's id
func (s *FilesystemStorage) Upload(r io.Reader) (string, error) {
	name := GenerateStorageID()

	// create a new file
	file, err := s.fs.Create(name)
	if err != nil {
		return "", err
	}
	defer file.Close()

	// copy contents to the new file
	if _, err := io.Copy(file, r); err != nil {
		s.fs.Remove(name)
		return "", err
	}

	return name, nil
}

// Delete deletes a filesystem file by it's id
func (s *FilesystemStorage) Delete(id string) error {
	return s.fs.Remove(id)
}
