package imgbucket

import (
	"io/ioutil"
	"log"
	"os"
	"testing"
)

func TestAliko(t *testing.T) {
	err := ioutil.WriteFile("/tmp/ali", []byte("aliko"), os.ModePerm)

	log.Println(err)
}
