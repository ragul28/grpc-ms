build:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build 
	docker build -t user-service .

run:
	go build && ./user-service
	
init:
	GO111MODULE=on go mod init github.com/ragul28/grpc-ms/user-service
	GO111MODULE=on go get -u

grpcgw:
	protoc -I/usr/local/include -I. \
	-I${GOPATH}/src -I${GOPATH}/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis \
	--go_out=plugins=grpc:. proto/user/user.proto

	protoc -I/usr/local/include -I. \
	-I${GOPATH}/src -I${GOPATH}/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis \
	--grpc-gateway_out=logtostderr=true:. proto/user/user.proto

	protoc -I/usr/local/include -I. \
	-I${GOPATH}/src -I${GOPATH}/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis \
	--swagger_out=logtostderr=true:. proto/user/user.proto

# grpc:
# 	protoc -I. --go_out=plugins=grpc:. \
# 	  proto/user/user.proto
