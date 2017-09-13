package mock_providers

import (
	"io"
	"os"
)

type MockFileSystem struct {
	CreateFnc func(string, io.Reader) error
	OpenFnc   func(string) (*os.File, error)
	RemoveFnc func(string) error
}

func (fs *MockFileSystem) Create(name string, data io.Reader) error {
	if fs.CreateFnc != nil {
		return fs.CreateFnc(name, data)
	}
	panic("not implemented")
}

func (fs *MockFileSystem) Open(name string) (*os.File, error) {
	if fs.OpenFnc != nil {
		return fs.OpenFnc(name)
	}
	panic("not implemented")
}

func (fs *MockFileSystem) Remove(name string) error {
	if fs.RemoveFnc != nil {
		return fs.RemoveFnc(name)
	}
	panic("not implemented")
}
