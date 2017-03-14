//go:generate mockgen -destination mock_providers/mock_filesystem.go github.com/alioygur/imgbucket/providers Filesystem

package providers

import (
	"io"
	"os"
	"path"

	"github.com/alioygur/imgbucket/service"
	"github.com/pkg/errors"
)

type (
	fileSystem struct {
		basePath string
	}

	// Filesystem just form mocking
	Filesystem service.FileSystem
)

// NewFileSystem instances a FileSystem
func NewFileSystem(basePath string) service.FileSystem {
	return &fileSystem{basePath}
}

// Create a file on filesystem then returns it
func (fs *fileSystem) Create(file string, data io.Reader) (*os.File, error) {
	fp := path.Join(fs.basePath, path.Dir(file))
	if err := os.MkdirAll(fp, os.ModePerm); err != nil {
		return nil, errors.WithStack(err)
	}
	f, err := os.Create(path.Join(fp, path.Base(file)))
	if err != nil {
		return nil, errors.WithStack(err)
	}
	defer f.Close()

	_, err = io.Copy(f, data)
	return f, errors.WithStack(err)
}

// Open file from filesystem
func (fs *fileSystem) Open(file string) (*os.File, error) {
	f, err := os.Open(path.Join(fs.basePath, file))
	return f, errors.WithStack(err)
}

// Remove file from filesystem
func (fs *fileSystem) Remove(file string) error {
	err := os.Remove(path.Join(fs.basePath, file))
	return errors.WithStack(err)
}
