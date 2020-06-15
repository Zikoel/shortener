GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
BINARY_NAME = shortener

all: test build
build: 
				$(GOBUILD) -o bin/$(BINARY_NAME) -v ./cmd/shortener
test: 
				$(GOTEST) -v ./...