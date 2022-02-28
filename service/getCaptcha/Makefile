
GOPATH:=$(shell go env GOPATH)
.PHONY: init
init:
	go get -u github.com/golang/protobuf/proto
	go get -u github.com/golang/protobuf/protoc-gen-go
	go get github.com/micro/micro/v3/cmd/protoc-gen-micro
	go get github.com/micro/micro/v3/cmd/protoc-gen-openapi

.PHONY: api
api:
	protoc --openapi_out=. --proto_path=. proto/getCaptcha.proto

.PHONY: proto
proto:
	protoc --proto_path=. --micro_out=. --gofast_out=:. proto/getCaptcha.proto
	
.PHONY: build
build:
	go build -o getCaptcha *.go

.PHONY: test
test:
	go test -v ./... -cover

.PHONY: docker
docker:
	docker build . -t getCaptcha:latest
