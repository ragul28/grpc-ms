build:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build 
	docker build -t consignment-service .

run:
	go build && ./consignment-service

init:
	GO111MODULE=on go mod init github.com/grpc-ms/consignment-service
	GO111MODULE=on go get -u
	#go get -u github.com/golang/protobuf/{proto,protoc-gen-go}

grpcgw:
	protoc -I/usr/local/include -I. \
	-I${GOPATH}/src -I${GOPATH}/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis \
	--go_out=plugins=grpc:. proto/consignment/consignment.proto

	protoc -I/usr/local/include -I. \
	-I${GOPATH}/src -I${GOPATH}/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis \
	--grpc-gateway_out=logtostderr=true:. proto/consignment/consignment.proto

	protoc -I/usr/local/include -I. \
	-I${GOPATH}/src -I${GOPATH}/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis \
	--swagger_out=logtostderr=true:. proto/consignment/consignment.proto

# grpc:
# 	protoc -I. --go_out=plugins=grpc:. \
# 	  proto/consignment/consignment.proto
