# build
build:
	CGO_ENABLED=0 GOOS=linux go build -a -tags netgo -ldflags '-w' -o imgbucket cmd/api/main.go

docker-build: build
	docker build -t alioygur/imgbucket .

# Dev commands
dev-build:
	go build -o imgbucket github.com/alioygur/imgbucket/cmd/api

compose-up: dev-build
	docker-compose up -d

compose-down:
	docker-compose down

compose-build:
	docker-compose build

compose-restart: compose-down compose-build compose-up
	
compose-restart-web: dev-build
	docker-compose restart web
