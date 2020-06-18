# Structure from https://sohlich.github.io/post/go_makefile/

GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
MOCKGEN=$(GOPATH)/bin/mockgen
BINARY_NAME = shortener

all: test build
build: 
				$(GOBUILD) -o ./bin/$(BINARY_NAME) -v ./cmd/shortener
test:
				$(MOCKGEN) -source=./pkg/shortener/shortener.go -destination=./mocks/mock_shortener.go
				$(GOTEST) -v ./pkg/shortener
clean: 
				$(GOCLEAN) ./cmd/shortener
				rm -f ./bin/$(BINARY_NAME)