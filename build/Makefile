TARGET := http-random-stress
GOPATH := $(shell go env GOPATH)

GOFILES := $(shell find . -name "*.go")

GITCOMMIT := $(shell git rev-list -1 HEAD)
GITREPOSITORY := https://github.com/alandiegosantos/http-random-stress.git

all: $(TARGET)

clean:
	rm -rf $(TARGET)

$(TARGET): $(GOFILES)
	CGO_ENABLED=0 GOARCH=amd64 go build -a -tags netgo -ldflags "-w -X 'main.GitCommit=$(GITCOMMIT)' -X 'main.GitRepository=$(GITREPOSITORY)'" -o $(TARGET) main.go 
