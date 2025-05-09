export PATH := $(GOPATH)/bin:$(PATH)
export GO111MODULE=on
LDFLAGS := -s -w

all: fmt build

build: LogBeetle

LogBeetle:
	env CGO_ENABLED=0 go build -trimpath -ldflags "$(LDFLAGS)" -o bin/LogBeetle ./cmd/LogBeetle


# compile assets into binary file
file:
	rm -rf ./assets/frps/static/*
	rm -rf ./assets/frpc/static/*
	cp -rf ./web/frps/dist/* ./assets/frps/static
	cp -rf ./web/frpc/dist/* ./assets/frpc/static

fmt:
	go fmt ./...

package:
	./package.sh

vet:
	go vet ./...

LogBeetle:
	env CGO_ENABLED=0 go build -trimpath -ldflags "$(LDFLAGS)" ./cmd/LogBeetle

test: gotest

gotest:
	go test -v --cover ./cmd/...
	go test -v --cover ./pkg/...

clean:
	rm -f ./bin/*