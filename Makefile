.PHONY: clean deps docker build

EXECUTABLE ?= drone-pkgbuild
IMAGE ?= mikkeloscar/$(EXECUTABLE)

all: build

clean:
	go clean -i ./..

deps:
	go get -t ./..

docker-run: docker
	./test_plugin.sh

docker:
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -ldflags "-s"
	docker build --rm -t $(IMAGE) .

docker-test: docker
	docker build --rm -t $(IMAGE)-test -f Dockerfile.test .
	docker run --rm $(IMAGE)-test

$(EXECUTABLE): $(wildcard *.go)
	go build

build: $(EXECUTABLE)
