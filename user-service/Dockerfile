FROM alpine:latest

RUN mkdir /app
WORKDIR /app
COPY user-service .

CMD ["./user-service"]
EXPOSE 50051