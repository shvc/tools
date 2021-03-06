BINARY=sshtail
BUILDDATEID=$(shell date +'%Y-%m-%dT%H:%M:%SZ')@$(shell git rev-parse --short HEAD)
VERSION=1.0.$(shell git rev-list HEAD --count)
LONGVER=${VERSION}@${BUILDDATEID}
REGISTRY='myshare.io:5000'
LDFLAGS=-ldflags "-X main.version=${LONGVER}"

.DEFAULT_GOAL:=default

test:
	go test ./...

vet:
	go vet ./...

default:
	@echo "Building ${BINARY}-${VERSION}"
	go build ${LDFLAGS} -o ${BINARY}

image: clean
	@echo "Build docker image: ${REGISTRY}/${BINARY}:${VERSION}"
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a ${LDFLAGS} -o ./docker/${BINARY}
	docker build -t ${REGISTRY}/${BINARY}:${VERSION} -f docker/Dockerfile docker/

pkg: image
	docker save ${REGISTRY}/${BINARY}:${VERSION} -o ${BINARY}-${VERSION}.tar

clean:
	go clean
	@rm -f ${BINARY}

.PHONY: test vet default image pkg clean
