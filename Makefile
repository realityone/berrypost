gosrc = $(GOPATH)/src
gopbsrc = $(gosrc)/github.com/golang/protobuf
current_dir = $(shell pwd)
api_dir = $(current_dir)/api

proto:
	protoc --go_out=$(api_dir) --go-grpc_out=$(api_dir) --proto_path=/usr/local/include:$(gosrc):$(gopbsrc):$(api_dir) $(api_dir)/*.proto
