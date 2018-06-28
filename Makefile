GOOS ?= $(uname -s)
GOARCH ?= amd64
export GOOS
export GOARCH

GO_PKG = gitlab.jetstack.net/jetstack/vault-unsealer

REGISTRY := quay.io/jetstack
IMAGE_NAME := vault-unsealer
BUILD_TAG := build
IMAGE_TAGS := canary

BUILD_IMAGE_NAME := golang:1.9.2

GOPATH ?= /tmp/go

CI_COMMIT_TAG ?= unknown
CI_COMMIT_SHA ?= unknown

help:
	# all 		- runs verify, build and docker_build targets
	# test 		- runs go_test target
	# build 	- runs go_build target
	# verify 	- verifies generated files & scripts

# Util targets
##############
.PHONY: all build verify

all: verify build docker_build

verify: go_verify

.builder_image:
	docker pull ${BUILD_IMAGE_NAME}

test: go_test
# Go targets
#################
go_verify: go_fmt go_vet go_test

clean:
	go clean

get:
	go get -t .

## Build a statically linked binary using a Docker container
BUILD_APP_PATH = /gopath/src/github.com/starlingbank/$(shell basename $(shell pwd))
build: clean get
	docker run --rm -t -v "$(GOPATH)":/gopath -v "$(shell pwd)":"$(BUILD_APP_PATH)" -e "GOPATH=/gopath" -w $(BUILD_APP_PATH) golang:1.9.2-alpine3.7 sh -c 'CGO_ENABLED=0 go build -a -tags -netgo --installsuffix cgo --ldflags="-s -w" -o vault-unsealer'

go_test:
	go test $$(go list ./... | grep -v '/vendor/')

docker: build
	docker build -t quay.io/starlingbank/vault-unsealer:$(BUILD_NUMBER) .
	docker build -t quay.io/starlingbank/vault-unsealer:latest .

ifneq ($(findstring SNAPSHOT, $(BUILD_NUMBER)), SNAPSHOT)
  ifeq ($(shell git rev-parse --abbrev-ref HEAD), master)
publish: docker
	docker push quay.io/starlingbank/talebearer:latest
	docker push quay.io/starlingbank/talebearer:$(BUILD_NUMBER)
  else
publish:
	$(warning skipping target "publish", not on master branch)
  endif
else
publish:
	$(error the target "publish" requires that BUILD_NUMBER be set)
endif

.PHONY: install, build, docker, test, publish
