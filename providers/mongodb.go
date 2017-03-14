//go:generate mockgen -destination mock_providers/mock_mongodb.go github.com/alioygur/imgbucket/providers Repository

package providers

import (
	"github.com/alioygur/imgbucket"
	"github.com/alioygur/imgbucket/service"
	"github.com/pkg/errors"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type (
	mongodb struct {
		sess *mgo.Session
	}
	nextID struct {
		Next int64 `bson:"n"`
	}

	// Repository just for mock
	Repository service.Repository
)

// tables
const (
	idTbl      = `ids`
	usersTbl   = `users`
	bucketsTbl = `buckets`
	imagesTbl  = `images`
)

// errors
var (
	ErrNotFound = errors.New("not found")
)

// NewMongoDBRepository instances new repository
func NewMongoDBRepository(sess *mgo.Session) service.Repository {
	return &mongodb{sess: sess}
}

// id returns next id.
// if sess nil then it uses default session
func (r *mongodb) id(c string) int64 {
	ids := r.c(idTbl)
	change := mgo.Change{
		Update:    bson.M{"$inc": bson.M{"n": 1}},
		Upsert:    true,
		ReturnNew: true,
	}
	id := new(nextID)
	ids.Find(bson.M{"_id": c}).Apply(change, id)
	return id.Next
}

// c returns mgo collection.
// if sess nil then it uses default session
func (r *mongodb) c(c string) *mgo.Collection {
	return r.sess.DB("").C(c)
}

func (r *mongodb) oneBy(c *mgo.Collection, q interface{}, result interface{}) error {
	return c.Find(q).One(result)
}

func (r *mongodb) oneByID(c *mgo.Collection, id uint64, result interface{}) error {
	return r.oneBy(c, bson.M{"id": id}, result)
}

func (r *mongodb) existsBy(c *mgo.Collection, q interface{}) (bool, error) {
	n, err := c.Find(q).Limit(1).Count()
	return n > 0, err
}

func wrapErr(err error) error {
	return err
}

func (r *mongodb) ImageExistsByUserIDAndBucketIDAndName(uid int64, bid int64, name string) (bool, error) {
	return r.existsBy(r.c(imagesTbl), bson.M{"user_id": uid, "bucket_id": bid, "name": name})
}

func (r *mongodb) ImageExistsByName(name string) (bool, error) {
	return r.existsBy(r.c(imagesTbl), bson.M{"name": name})
}

func (r *mongodb) DefaultBucketByUserID(uid int64) (*imgbucket.Bucket, error) {
	var bucket imgbucket.Bucket
	return &bucket, r.oneBy(r.c(bucketsTbl), bson.M{"user_id": uid, "default": true}, &bucket)
}

func (r *mongodb) BucketByUserIDAndName(int64, string) (*imgbucket.Bucket, error) {
	panic("not implemented")
}

func (r *mongodb) IsNotFoundErr(err error) bool {
	return errors.Cause(err) == ErrNotFound
}

func (r *mongodb) AddBucket(b *imgbucket.Bucket) error {
	b.ID = r.id(bucketsTbl)
	err := r.c(bucketsTbl).Insert(b)
	return errors.WithStack(err)
}

func (r *mongodb) AddImage(i *imgbucket.Image) error {
	i.ID = r.id(imagesTbl)
	err := r.c(imagesTbl).Insert(i)
	return errors.WithStack(err)
}
