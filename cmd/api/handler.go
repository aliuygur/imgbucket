package main

import (
	"log"
	"net/http"
	"strings"

	"github.com/alioygur/gores"
	"github.com/alioygur/imgbucket/service"
	"github.com/pkg/errors"
)

type (
	handler struct {
		service.Service
		ps http.Handler // proxy server
	}
)

// NewHandler instances new handler
func NewHandler(s service.Service, ps http.Handler) http.Handler {
	return &handler{s, ps}
}

func (h *handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if strings.HasPrefix(r.URL.Path, "/images") && r.Method == http.MethodPost {
		if err := h.upload(w, r); err != nil {
			if err, ok := err.(*appErr); ok {
				gores.JSON(w, err.HTTPCode, err)
				log.Println(err.Stack())
				return
			}
			// fallback error handling
			gores.Error(w, http.StatusInternalServerError, err.Error())
		}
	}

	h.proxyPass(w, r)
}

func (h *handler) upload(w http.ResponseWriter, r *http.Request) error {
	// r.Body = http.MaxBytesReader(w, r.Body, 1024*2)
	infile, fh, err := r.FormFile("file")
	if err != nil {
		return &appErr{ImageSizeTooBigErrCode, http.StatusBadRequest, errors.WithStack(err)}
	}
	defer infile.Close()
	_ = fh

	ur := service.UploadRequest{ImageName: fh.Filename, FileContent: infile, BucketName: "default"}
	img, err := h.Upload(&ur)
	if err != nil {
		switch errors.Cause(err) {
		default:
			return &appErr{UnknownErrCode, http.StatusInternalServerError, err}
		case service.ErrImageAlreadyExists:
			return &appErr{ImageAlreadyExistsErrCode, http.StatusBadRequest, err}
		case service.ErrInvalidBucketName:
			return &appErr{InvalidBucketNameErrCode, http.StatusBadRequest, err}
		case service.ErrInvalidImage:
			return &appErr{InvalidImageErrCode, http.StatusBadRequest, err}
		}
	}
	return gores.JSON(w, http.StatusCreated, img)
}

func (h *handler) proxyPass(w http.ResponseWriter, r *http.Request) {
	// u, err := h.GenProxyURL(r)
	// if err != nil {
	// 	http.Error(w, err.Error(), http.StatusBadRequest)
	// 	return
	// }

	req, err := http.NewRequest(http.MethodGet, "", nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	h.ps.ServeHTTP(w, req)
}
