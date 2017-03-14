package providers

import (
	"fmt"
	"os"
	"reflect"
	"strconv"
	"testing"
	"time"

	"gopkg.in/mgo.v2/bson"

	"github.com/alioygur/imgbucket"
	mgo "gopkg.in/mgo.v2"
)

func randDBName() string {
	return "testdb_" + strconv.Itoa(time.Now().Nanosecond())
}

func newTestRepo(dropDB bool) (*mongodb, func(), error) {
	dbname := "imgbucket_test"
	if dropDB {
		dbname = randDBName()
	}
	mhost := fmt.Sprintf("%s/%s", os.Getenv("MONGO_HOST"), dbname)
	sess, err := mgo.Dial(mhost)
	if err != nil {
		return nil, nil, err
	}
	teardown := func() {
		if dropDB {
			sess.DB("").DropDatabase()
		}
		sess.Close()
	}

	r := mongodb{sess: sess}
	return &r, teardown, nil
}

func Test_repository_AddBucket(t *testing.T) {
	r, teardown, err := newTestRepo(false)
	if err != nil {
		t.Fatal(err)
	}
	defer teardown()

	var want imgbucket.Bucket
	want.Name = "bucket"
	want.UserID = 1

	if err := r.AddBucket(&want); err != nil {
		t.Error(err)
	}

	var got imgbucket.Bucket
	if err := r.c(bucketsTbl).Find(bson.M{"id": want.ID}).One(&got); err != nil {
		t.Error(err)
		return
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("not equal, got %+v, want %+v", got, want)
	}
}

func Test_repository_AddImage(t *testing.T) {
	r, teardown, err := newTestRepo(false)
	if err != nil {
		t.Fatal(err)
	}
	defer teardown()

	var want imgbucket.Image
	want.BucketID = 1
	want.Name = "image"
	want.Width = 400
	want.Height = 600
	want.Size = 1024
	want.Format = "image/jpg"
	want.IsDefault = true

	if err := r.AddImage(&want); err != nil {
		t.Error(err)
	}

	var got imgbucket.Image
	if err := r.c(imagesTbl).Find(bson.M{"id": want.ID}).One(&got); err != nil {
		t.Error(err)
		return
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("not equal, got %+v, want %+v", got, want)
	}
}
