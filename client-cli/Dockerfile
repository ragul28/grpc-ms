FROM alpine:latest

RUN mkdir /app
WORKDIR /app

COPY consignment.json consignment.json 
COPY client-cli .

ENTRYPOINT ["./client-cli"]