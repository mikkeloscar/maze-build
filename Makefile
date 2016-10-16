.PHONY: clean deps docker build

EXECUTABLE ?= maze-build
IMAGE      ?= mikkeloscar/$(EXECUTABLE)
SOURCES    = $(shell find . -name '*.go')

all: build

clean:
	go clean -i ./..

deps:
	go get -t

docker: build
	docker build --rm -t $(IMAGE) .

docker-test:
	docker build --rm -t $(IMAGE)-test -f Dockerfile.test .
	docker run --rm $(IMAGE)-test

$(EXECUTABLE): $(SOURCES)
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -ldflags "-s"

build: $(EXECUTABLE)
