package service

import (
	"io"
	"os"

	"github.com/alioygur/imgbucket"
)

type (
	// FileSystem interface
	FileSystem interface {
		Create(string, io.Reader) error
		Open(string) (*os.File, error)
		Remove(string) error
	}

	// Service interface
	Service interface {
		Upload(*UploadRequest) (*imgbucket.Image, error)
		// GenProxyURL(*http.Request) (string, error)
	}

	// service is implementation of Service interface
	service struct {
		fs   FileSystem
		repo Repository
	}

	// Repository interface
	Repository interface {
		ImageExistsByUserIDAndBucketIDAndName(uid int64, bid int64, name string) (bool, error)
		ImageExistsByName(name string) (bool, error)
		DefaultBucketByUserID(uid int64) (*imgbucket.Bucket, error)
		BucketByUserIDAndName(uid int64, name string) (*imgbucket.Bucket, error)
		IsNotFoundErr(err error) bool
		AddBucket(bkt *imgbucket.Bucket) error
		AddImage(img *imgbucket.Image) error
	}
)

// NewService instances a impl of ImageBucketServer
func NewService(fs FileSystem, repo Repository) Service {
	return &service{fs, repo}
}
