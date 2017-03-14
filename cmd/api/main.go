package main

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"time"

	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"

	_ "golang.org/x/image/webp"
	mgo "gopkg.in/mgo.v2"

	"github.com/alioygur/goutil"
	"github.com/alioygur/imgbucket/providers"
	"github.com/alioygur/imgbucket/service"
)

func main() {
	if err := ensureEnv(); err != nil {
		log.Fatal(err)
	}

	if err := waitForServices(); err != nil {
		log.Fatal(err)
	}

	// deps
	// proxy server
	targetURL, err := url.Parse(os.Getenv("PROXY_URL"))
	if err != nil {
		log.Fatal(err)
	}
	ps := httputil.NewSingleHostReverseProxy(targetURL)

	// filesystem
	fs := providers.NewFileSystem(os.Getenv("UPLOAD_PATH"))
	// repo
	msess, err := mgo.Dial(os.Getenv("MONGO_HOST") + "/" + os.Getenv("MONGO_DB"))
	if err != nil {
		log.Fatal(err)
	}
	repo := providers.NewMongoDBRepository(msess)
	// service
	s := service.NewService(fs, repo)

	h := NewHandler(s, ps)

	log.Printf("application starting on port: %s", os.Getenv("PORT"))
	if err := http.ListenAndServe(":"+os.Getenv("PORT"), h); err != nil {
		log.Fatal(err)
	}
}

func ensureEnv() error {
	envs := map[string]string{
		"MONGO_HOST":  "mongo",
		"MONGO_DB":    "imgbucket",
		"PORT":        "5000",
		"UPLOAD_PATH": "/images",
		"PROXY_URL":   "",
	}
	for k, v := range envs {
		if os.Getenv(k) == "" && v == "" {
			return fmt.Errorf("missing required env:  %s", k)
		}
	}
	return nil
}

func waitForServices() error {
	img, err := url.Parse(os.Getenv("PROXY_URL"))
	if err != nil {
		log.Fatal(err)
	}

	services := []url.URL{
		url.URL{Scheme: "tcp", Host: img.Host},
	}

	if err := goutil.WaitForServices(services, 15*time.Second); err != nil {
		return err
	}
	log.Println("services are ready")
	return nil
}
