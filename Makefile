# Go parameters
GOCMD=go
GOTEST=$(GOCMD) test
DEPCMD=dep ensure

all: test
test:
		$(GOTEST) -v ./...
deps:
		$(DEPCMD)
