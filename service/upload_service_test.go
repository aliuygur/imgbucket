package service

import (
	"testing"

	"github.com/alioygur/imgbucket"
	"github.com/pkg/errors"

	"github.com/alioygur/imgbucket/providers/mock_providers"
)

func Test_service_Upload(t *testing.T) {
	t.Parallel()

	fs := struct {
		FileSystem
	}{}

	// test errors states

	t.Run("with empty bucket name also user doesn't have default bucket", func(t *testing.T) {
		errNoDefaultBucket := errors.New("no default bucket")
		repo := &mock_providers.MockRepo{}
		repo.DefaultBucketByUserIDFnc = func(int64) (*imgbucket.Bucket, error) {
			return nil, errNoDefaultBucket
		}

		s := NewService(fs, repo)

		_, err := s.Upload(&UploadRequest{})
		if err != errNoDefaultBucket {
			t.Errorf("want err: %v, got: %v", errNoDefaultBucket, err)
		}
	})

	t.Run("with not exists bucket name", func(t *testing.T) {
		repo := &mock_providers.MockRepo{}
		repo.BucketByUserIDAndNameFnc = func(int64, string) (*imgbucket.Bucket, error) {
			return nil, errors.New("not found")
		}
		repo.IsNotFoundErrFnc = func(error) bool {
			return true
		}

		s := NewService(fs, repo)

		_, err := s.Upload(&UploadRequest{BucketName: "bucket"})
		if err != ErrBucketNotFound {
			t.Errorf("want err: %v, got: %v", ErrBucketNotFound, err)
		}
	})

	t.Run("with invalid image name", func(t *testing.T) {
		repo := &mock_providers.MockRepo{}
		repo.DefaultBucketByUserIDFnc = func(int64) (*imgbucket.Bucket, error) {
			return &imgbucket.Bucket{}, nil
		}
		s := NewService(fs, repo)

		_, err := s.Upload(&UploadRequest{ImageName: "/ "})
		if err != ErrInvalidImageName {
			t.Errorf("want err: %v, got: %v", ErrInvalidImageName, err)
		}
	})

	t.Run("with image name that already exists", func(t *testing.T) {
		repo := &mock_providers.MockRepo{}
		repo.DefaultBucketByUserIDFnc = func(int64) (*imgbucket.Bucket, error) {
			return &imgbucket.Bucket{ID: 1}, nil
		}
		repo.ImageExistsByUserIDAndBucketIDAndNameFnc = func(int64, int64, string) (bool, error) {
			return true, nil
		}
		s := NewService(fs, repo)

		_, err := s.Upload(&UploadRequest{ImageName: "image"})
		if err != ErrImageAlreadyExists {
			t.Errorf("want err: %v, got: %v", ErrImageAlreadyExists, err)
		}
	})
}
