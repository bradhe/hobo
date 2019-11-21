.PHONY: build

VERSION := 1.0
GIT_COMMIT := $(shell git rev-list -1 HEAD)

DOCKER_IMAGE := bradhe/location-search
DOCKER_TAG := latest

clean:
	rm -rf ./bin
	mkdir ./bin

build: clean
	go build -ldflags="-X main.gitCommit=$(GIT_COMMIT) -X main.version=$(VERSION)" -o ./bin/location-search ./cmd/location-search

image: clean
	GOOS=linux GOARCH=amd64 go build -ldflags="-X main.gitCommit=$(GIT_COMMIT) -X main.version=$(VERSION)" -o ./bin/location-search ./cmd/location-search
	docker build -t $(DOCKER_IMAGE):$(DOCKER_TAG) .

test: clean
	go test -v ./...
