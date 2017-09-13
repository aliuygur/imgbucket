package providers

import (
	"bytes"
	"io/ioutil"
	"os"
	"reflect"
	"strings"
	"testing"

	"github.com/pkg/errors"
)

func TestNewFileSystem(t *testing.T) {
	type args struct {
		basePath string
	}
	tests := []struct {
		name string
		args args
		want *fileSystem
	}{
		{"with valid", args{basePath: "/tmp/images"}, &fileSystem{"/tmp/images"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewFileSystem(tt.args.basePath); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewFileSystem() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_fileSystem_Create(t *testing.T) {
	t.Parallel()
	t.Run("with empty filename", func(t *testing.T) {
		fs := &fileSystem{}
		if err := fs.Create("", nil); errors.Cause(err) != errInvalidFileName {
			t.Errorf("not equal, got %v, want %v", err, errInvalidFileName)
			return
		}
	})
	t.Run("pass nil instead of reader as file data", func(t *testing.T) {
		fs := &fileSystem{}

		if err := fs.Create("file", nil); errors.Cause(err) != errInvalidFileData {
			t.Errorf("not equal, got %v, want %v", err, errInvalidFileData)
			return
		}
	})
	t.Run("with valid args", func(t *testing.T) {
		fs := &fileSystem{"/tmp"}
		if err := fs.Create("file1", strings.NewReader("filedata1")); err != nil {
			t.Fatal(err)
		}

		f, err := fs.Open("file1")
		if err != nil {
			t.Fatal(err)
		}
		defer f.Close()

		fd, err := ioutil.ReadAll(f)
		if err != nil {
			t.Fatal(err)
		}
		if bytes.Compare(fd, []byte("filedata1")) != 0 {
			t.Errorf("not equal, got %v want %v", fd, []byte("filedata1"))
		}
	})
}

func Test_fileSystem_Open(t *testing.T) {
	type fields struct {
		basePath string
	}
	type args struct {
		file string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *os.File
		wantErr bool
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fs := &fileSystem{
				basePath: tt.fields.basePath,
			}
			got, err := fs.Open(tt.args.file)
			if (err != nil) != tt.wantErr {
				t.Errorf("fileSystem.Open() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("fileSystem.Open() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_fileSystem_Remove(t *testing.T) {
	type fields struct {
		basePath string
	}
	type args struct {
		file string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fs := &fileSystem{
				basePath: tt.fields.basePath,
			}
			if err := fs.Remove(tt.args.file); (err != nil) != tt.wantErr {
				t.Errorf("fileSystem.Remove() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
