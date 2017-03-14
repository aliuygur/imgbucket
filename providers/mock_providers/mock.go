package mock_providers

import imgbucket "github.com/alioygur/imgbucket"

type MockRepo struct {
	ImageExistsByUserIDAndBucketIDAndNameFnc func(uid int64, bid int64, name string) (bool, error)
	ImageExistsByNameFnc                     func(name string) (bool, error)
	DefaultBucketByUserIDFnc                 func(uid int64) (*imgbucket.Bucket, error)
	BucketByUserIDAndNameFnc                 func(uid int64, name string) (*imgbucket.Bucket, error)
	IsNotFoundErrFnc                         func(err error) bool
	AddBucketFnc                             func(bkt *imgbucket.Bucket) error
	AddImageFnc                              func(img *imgbucket.Image) error
}

func NewMockRepo() *MockRepo {
	return &MockRepo{}
}

func (m *MockRepo) ImageExistsByUserIDAndBucketIDAndName(uid int64, bid int64, name string) (bool, error) {
	if m.ImageExistsByUserIDAndBucketIDAndNameFnc != nil {
		return m.ImageExistsByUserIDAndBucketIDAndNameFnc(uid, bid, name)
	}
	panic("not implemented")
}

func (m *MockRepo) ImageExistsByName(name string) (bool, error) {
	if m.ImageExistsByNameFnc != nil {
		return m.ImageExistsByNameFnc(name)
	}
	panic("not implemented")
}

func (m *MockRepo) DefaultBucketByUserID(uid int64) (*imgbucket.Bucket, error) {
	if m.DefaultBucketByUserIDFnc != nil {
		return m.DefaultBucketByUserIDFnc(uid)
	}
	panic("not implemented")
}

func (m *MockRepo) BucketByUserIDAndName(uid int64, name string) (*imgbucket.Bucket, error) {
	if m.BucketByUserIDAndNameFnc != nil {
		return m.BucketByUserIDAndNameFnc(uid, name)
	}
	panic("not implemented")
}

func (m *MockRepo) IsNotFoundErr(err error) bool {
	if m.IsNotFoundErrFnc != nil {
		return m.IsNotFoundErrFnc(err)
	}
	panic("not implemented")
}

func (m *MockRepo) AddBucket(bkt *imgbucket.Bucket) error {
	if m.AddBucketFnc != nil {
		return m.AddBucketFnc(bkt)
	}
	panic("not implemented")
}

func (m *MockRepo) AddImage(img *imgbucket.Image) error {
	if m.AddImageFnc != nil {
		return m.AddImageFnc(img)
	}
	panic("not implemented")
}
