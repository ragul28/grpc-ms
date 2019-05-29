FROM alpine:latest

RUN mkdir /app
WORKDIR /app
COPY consignment.json .
COPY client-cli .

ENTRYPOINT ["./client-cli"]