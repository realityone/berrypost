.DEFAULT_GOAL := all

gosrc = $(GOPATH)/src
gopbsrc = $(gosrc)/github.com/golang/protobuf
current_dir = $(shell pwd)
api_dir = $(current_dir)/api

proto:
	protoc --go_out=$(api_dir) --go-grpc_out=$(api_dir) --proto_path=/usr/local/include:$(gosrc):$(gopbsrc):$(api_dir) $(api_dir)/*.proto

dist:
	pushd ./statics && \
	npm run build && \
	popd

berrypost:
	go build -o cli/berrypost cli/main.go

all: clean proto dist berrypost

clean:
	rm ./cli/berrypost ./cli/main
	rm -rf ./api/*.pb.go
	rm -rf ./statics/dist/*