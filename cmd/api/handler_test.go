// +build ignore

package main

import (
	"bytes"
	"crypto/md5"
	"encoding/json"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"

	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"

	"github.com/alioygur/imgbucket"

	_ "golang.org/x/image/webp"
)

type (
	mockProxyServer struct{}
	mockRepository  struct{}
)

func newMockProxyServer() http.Handler {
	return &mockProxyServer{}
}

func newMockRepository() imgbucket.Repository {
	return &mockRepository{}
}

func (h *mockProxyServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Println("proxy passed")
}

func (r *mockRepository) ImageExistsByName(string) (bool, error) {
	return false, nil
}

func TestUpload(t *testing.T) {
	fs := imgbucket.NewFileSystem("/tmp")
	repo := newMockRepository()
	s := imgbucket.NewService(fs, repo)
	ps := newMockProxyServer()
	h := NewHandler(s, ps)

	w := httptest.NewRecorder()
	r, err := newfileUploadRequest("/images", nil, "file", "./testdata/valid.png")
	if err != nil {
		t.Fatal(err)
	}
	h.ServeHTTP(w, r)

	// assert http status code
	assert.Equal(t, http.StatusCreated, w.Code, "http status code")

	// assert files checksums
	var img imgbucket.Image
	if err := json.NewDecoder(w.Body).Decode(&img); err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, fileSum("./testdata/valid.png"), fileSum("/tmp/"+img.FilePath()), "files checksums")
}

// Creates a new file upload http request with optional extra params
func newfileUploadRequest(uri string, params map[string]string, paramName, path string) (*http.Request, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile(paramName, filepath.Base(path))
	if err != nil {
		return nil, err
	}
	_, err = io.Copy(part, file)

	for key, val := range params {
		_ = writer.WriteField(key, val)
	}
	err = writer.Close()
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", uri, body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	return req, err
}

func fileSum(fp string) []byte {
	f, err := os.Open(fp)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	h := md5.New()
	if _, err := io.Copy(h, f); err != nil {
		log.Fatal(err)
	}

	return h.Sum(nil)
}
