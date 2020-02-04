PROJECT_NAME := "bingo"
PKG := "$(PROJECT_NAME)"
PKG_LIST := $(shell go list ${PKG}/... | grep -v /vendor/)

BUILD_PATH := "cmd/server.go"
OUTPUT_FILE := "build/server"
DOCKER_FILE := "build/package/Dockerfile"
GO_FILES := $(shell find . -name '*.go' | grep -v /vendor/ | grep -v _test.go)

all: build docker

build: release

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

debug: vet tidy verify format
	@CGO_ENABLED=0 go build -v -o $(OUTPUT_FILE) $(BUILD_PATH)

release: vet tidy verify format
	@CGO_ENABLED=0 go build -v -o $(OUTPUT_FILE) -ldflags '-s -w' $(BUILD_PATH)

docker:
	@sudo docker build -f $(DOCKER_FILE) -t bingo:latest .

clean:
	@rm -f $(OUTPUT_FILE)
