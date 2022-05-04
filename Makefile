GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
BINARY_NAME=ordinary-kriging-cli

all: build test
build:
	$(GOBUILD) -o $(BINARY_NAME) -v ./cmd/$(BINARY_NAME)
test:
	$(GOTEST) -v ./...
gentest: build
	$(GOTEST) -v ./... -gen_golden_files
clean:
	$(GOCLEAN)
	rm -f $(BINARY_NAME)