package providers

import (
	"io"
	"os"
	"path"

	"github.com/alioygur/is"

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

// Errors
var (
	errInvalidFileName = errors.New("invalid file name")
	errInvalidFileData = errors.New("invalid file data")
	errInvalidBasePath = errors.New("invalid base path")
)

// NewFileSystem instances a FileSystem
func NewFileSystem(basePath string) service.FileSystem {
	return &fileSystem{basePath}
}

// Create a file on filesystem then returns it
func (fs *fileSystem) Create(file string, data io.Reader) error {
	if file == "" || !is.Alphanumeric(file) {
		return errors.WithStack(errInvalidFileName)
	}
	if data == nil {
		return errors.WithStack(errInvalidFileData)
	}
	fp := path.Join(fs.basePath, path.Dir(file))
	if err := os.MkdirAll(fp, os.ModePerm); err != nil {
		return errors.WithStack(err)
	}
	f, err := os.Create(path.Join(fp, path.Base(file)))
	if err != nil {
		return errors.WithStack(err)
	}
	defer f.Close()

	_, err = io.Copy(f, data)
	return errors.WithStack(err)
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
