export GO111MODULE=on
export GOSUMDB=off
export GOPROXY=direct

SHELL=/bin/bash
IMAGE_TAG := $(shell git rev-parse HEAD)
DOCKER_REPO=hub.docker.com

.PHONY: all
all: gen deps deps_check lint unit_test build

.PHONY: ci
ci: lint unit_test

.PHONY: deps
deps:
	go mod tidy
	go mod download
	go mod vendor

.PHONY: deps_check
deps_check:
	go mod verify

.PHONY: build
build:
	CGO_ENABLED=0 GOOS=linux go build -mod=vendor -a -o artifacts/svc .

.PHONY: gen
gen:
	mockgen -package mock -source api/service.go Service > api/mock/service.go
	mockgen -package mock -source storage/decorator/unique_lines.go > storage/decorator/mock/unique_lines.go

.PHONY: unit_test
unit_test:
	go test -mod=vendor -count=1 -v -cover ./...

.PHONY: lint
lint:
	golangci-lint run

.PHONY: dockerise
dockerise:
	docker build -t "${DOCKER_REPO}/file-worker:${IMAGE_TAG}" .