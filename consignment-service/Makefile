build:
	protoc -I. --go_out=plugins=grpc:. \
	  proto/consignment/consignment.proto

	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build 
	docker build -t consignment-service .

run:
	docker run -p 50051:50051 \
	-e MICRO_SERVER_ADDRESS=:50051 \
	consignment-service

init:
	GO111MODULE=on go mod init github.com/grpc-ms/consignment-service
	GO111MODULE=on go get -u
	#go get -u github.com/golang/protobuf/{proto,protoc-gen-go}
