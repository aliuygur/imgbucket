package service

import (
	"bytes"
	"fmt"
	"image"
	"io"
	"io/ioutil"

	"github.com/pkg/errors"

	"github.com/ventu-io/go-shortid"

	"github.com/alioygur/imgbucket"
	"github.com/alioygur/is"
)

type (
	// UploadRequest ...
	UploadRequest struct {
		ImageName   string
		FileContent io.Reader
		BucketName  string
		UserID      int64
	}
)

// errors
var (
	ErrBucketNotFound     = errors.New("bucket not found")
	ErrInvalidBucketName  = errors.New("invalid bucket name")
	ErrInvalidImageName   = errors.New("invalid image name")
	ErrImageAlreadyExists = errors.New("image already exists")
	ErrInvalidImage       = errors.New("invalid image")
)

// Upload uploads an image to FileSystem
func (s *service) Upload(r *UploadRequest) (*imgbucket.Image, error) {
	// if bucket name has been provided find it else get user's default bucket.
	var bucket *imgbucket.Bucket
	if r.BucketName == "" {
		b, err := s.repo.DefaultBucketByUserID(r.UserID)
		if err != nil {
			return nil, err
		}
		bucket = b
	} else {
		b, err := s.repo.BucketByUserIDAndName(r.UserID, r.BucketName)
		if err != nil {
			if s.repo.IsNotFoundErr(err) {
				return nil, ErrBucketNotFound
			}
			return nil, err
		}
		bucket = b
	}

	// if image name hasn't been provided then generate a name, else validate it
	if r.ImageName == "" {
		uniqid, err := shortid.Generate()
		if err != nil {
			return nil, errors.WithStack(err)
		}
		r.ImageName = uniqid
	} else {
		if !is.Alphanumeric(r.ImageName) {
			return nil, ErrInvalidImageName
		}
	}

	var img imgbucket.Image
	img.BucketID = bucket.ID
	img.Name = r.ImageName

	exists, err := s.repo.ImageExistsByUserIDAndBucketIDAndName(r.UserID, bucket.ID, r.ImageName)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, ErrImageAlreadyExists
	}

	imgdata, err := ioutil.ReadAll(r.FileContent)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	cnf, format, err := image.DecodeConfig(bytes.NewReader(imgdata))
	if err != nil {
		return nil, ErrInvalidImage
	}

	img.Width = cnf.Width
	img.Height = cnf.Height
	img.Size = len(imgdata)
	img.Format = format

	_, err = s.fs.Create(fmt.Sprintf("%d/%d", bucket.UserID, img.ID), bytes.NewReader(imgdata))
	if err != nil {
		return nil, err
	}

	return &img, nil
}
