PROJECT_NAME := "bingo"
PKG := "$(PROJECT_NAME)"
PKG_LIST := $(shell go list ${PKG}/... | grep -v /vendor/)

BUILD_PATH := "cmd/server.go"
OUTPUT_FILE := "build/server"
DOCKER_FILE_RELEASE := "build/package/Dockerfile.release"
DOCKER_FILE_DEV := "build/package/Dockerfile.dev"
GO_FILES := $(shell find . -name '*.go' | grep -v /vendor/ | grep -v _test.go)

build: release

docker:
	@docker build -f $(DOCKER_FILE_RELEASE) -t bingo:latest .

docker-dev: release
	@docker build -f $(DOCKER_FILE_DEV) -t bingo-dev:latest .

debug: fix
	@CGO_ENABLED=0 go build -v -o $(OUTPUT_FILE) $(BUILD_PATH)

release: fix
	@CGO_ENABLED=0 go build -v -o $(OUTPUT_FILE) -ldflags '-s -w' $(BUILD_PATH)

fix: tidy verify format vet

vet:
	@go vet ${PKG_LIST}

lint:
	@golint -set_exit_status ${PKG_LIST}

test:
	@go test -cover -race -short ${PKG_LIST}

tidy:
	@go mod tidy

verify:
	@go mod verify

format:
	@gofmt -w -s -d .

clean:
	@rm -f $(OUTPUT_FILE)
