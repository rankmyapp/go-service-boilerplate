.PHONY: build run test test-integration swagger clean docker-up docker-down docker-build

build:
	go build -o bin/api cmd/api/main.go

run: swagger
	go run cmd/api/main.go

test:
	go test ./...

test-integration:
	go test -tags=integration -v ./...

swagger:
	swag init -g cmd/api/main.go

clean:
	rm -rf bin/

docker-build:
	docker compose build

docker-up:
	docker compose up -d

docker-down:
	docker compose down
