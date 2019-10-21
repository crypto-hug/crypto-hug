PROJECTNAME := $(shell basename "$(PWD)")
GOBASE := $(shell pwd)
GOPATH := $(GOBASE)/vendor:$(GOBASE)
GOBIN := $(GOBASE)/bin
GOFILES := $(wildcard *.go)

SVC_ENTRY := $(GOBASE)/cmd/chug-node
CLI_ENTRY := $(GOBASE)/cmd/chug-client
SVC_BIN_NAME := "chug-node"
CLI_BIN_NAME := "chug"



ci: go-compile go-test-cover

go-build-all:
	@echo "  ⚙️  Building binary..."
	@go build -o $(GOBIN)/$(SVC_BIN_NAME) $(SVC_ENTRY)
	@go build -o $(GOBIN)/$(CLI_BIN_NAME) $(CLI_ENTRY)


go-get:
	@echo "  🔎  Checking if there is any missing dependencies..."
	@GOPATH=$(GOPATH) GOBIN=$(GOBIN) go get $(get)

go-clean:
	@echo "  🗑  Cleaning build cache"
	@GOPATH=$(GOPATH) GOBIN=$(GOBIN) go clean

go-test-cover:
	@go test ./... -coverpkg=./... -coverprofile=coverage.txt  -timeout 30s
	@go tool cover -func=coverage.txt

go-compile: go-get go-build-all

