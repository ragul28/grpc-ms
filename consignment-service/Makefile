build:
	protoc -I. --go_out=plugins=grpc:. \
	  proto/consignment/consignment.proto

	GO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -o consignment-service
		
init:
	GO111MODULE=on go mod init github.com/grpc-ms/consignment-service
	GO111MODULE=on go get -u

run:
	./consignment-service
