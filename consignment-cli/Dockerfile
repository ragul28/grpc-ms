FROM alpine:latest

RUN mkdir /app
WORKDIR /app
COPY consignment.json .
COPY consignment-cli .

CMD ["./consignment-cli"]