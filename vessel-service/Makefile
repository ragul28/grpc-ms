build:
	protoc -I. --go_out=plugins=grpc:. \
	  proto/vessel/vessel.proto

	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build 
	docker build -t vessel-service .

run:
	docker run -p 50052:50051 \
	-e MICRO_SERVER_ADDRESS=:50051 \
	vessel-service

init:
	GO111MODULE=on go mod init github.com/grpc-ms/vessel-service
	GO111MODULE=on go get -u
