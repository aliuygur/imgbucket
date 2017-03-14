// +build ignore

package imgbucket

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
)

func (s *service) GenProxyURL(r *http.Request) (string, error) {
	op, w, h, file, err := s.parseURL(r.URL.Path)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("/%s?width=%s&height=%s&file=%s", op, w, h, file), nil
}

func (s *service) parseURL(urlPath string) (string, string, string, string, error) {
	// /images/{op}/{sizes}/{file}
	parts := strings.Split(urlPath, "/")[2:]

	if len(parts) != 3 {
		return "", "", "", "", errors.New("parts aren't equal 3")
	}

	op, size, file := parts[0], parts[1], parts[2]

	i := strings.Index(size, "x")
	if i == -1 {
		return "", "", "", "", fmt.Errorf("invalid size %s", size)
	}
	w, h := size[0:i], size[i+1:]

	return op, w, h, file, nil
}
