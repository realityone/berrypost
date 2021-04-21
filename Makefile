.DEFAULT_GOAL := all

gosrc = $(GOPATH)/src
gopbsrc = $(gosrc)/github.com/golang/protobuf
current_dir = $(shell pwd)
api_dir = $(current_dir)/api

proto:
	protoc --go_out=$(api_dir) --go-grpc_out=$(api_dir) --proto_path=/usr/local/include:$(gosrc):$(gopbsrc):$(api_dir) $(api_dir)/*.proto

dist:
	pushd ./statics && \
	npm install && \
	npm run build && \
	popd

berrypost-main:
	go build -o berrypost/berrypost berrypost/main.go

run: dist
	go run berrypost/main.go

all: clean proto dist berrypost-main

clean:
	rm -f ./berrypost/berrypost
	rm -f ./statics/dist/*
