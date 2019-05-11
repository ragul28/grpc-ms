build:
	protoc -I. --go_out=plugins=micro:. \
	  proto/user/user.proto

	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build 
	docker build -t user-service .

run:
	docker run -p 50053:50051 \
	-e MICRO_SERVER_ADDRESS=:50051 \
	vessel-service

init:
	GO111MODULE=on go mod init github.com/grpc-ms/user-service
	GO111MODULE=on go get -u
