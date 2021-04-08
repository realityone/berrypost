gosrc = $(GOPATH)/src
gogosrc = $(gosrc)/github.com/gogo/protobuf
current_dir = $(shell pwd)
api_dir = $(current_dir)/api

proto:
	protoc --gogo_out=plugins=grpc:$(api_dir) --proto_path=/usr/local/include:$(gosrc):$(gogosrc):$(api_dir) $(api_dir)/*.proto
