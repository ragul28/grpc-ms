build:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build
	docker build -t client-cli .

run:
	docker run client-cli