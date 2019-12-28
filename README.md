# grpc-ms

Exploration gRPC project with Microservice architecture.
written in golang.

*Stack used* - grpc, protobuf, docker, go-micro, Makefile, docker-compose, mongodb, postgresql.  

### Microservices
Container management system,
-   consignments
-   users & authentication
-   vessels
-   client cli 


### Docker Build  

1. Build all the microservices using make file in respective dir 
    ```sh
    make build
    ```
2. Then Build docker with docker-compose
    ```sh
    docker-compose build
    ```
3. Run all the microservices using docker-compose 
    ```sh
    docker-compose up 
    ```

> Build Info
- To build inside gopath use GO111MODULE=on environment variable; to build outside gopath unset GO111MODULE or set to GO111MODULE=auto.
- To run on alpine linux container build go-binary with CGO_ENABLED=0 flag set.

### Reference 
- Tutorial - [ Microservices in golang by Ewan Valentine ](https://ewanvalentine.io/microservices-in-golang-part-1/)
- [Building-high-performance-apis-in-go-using-grpc-and-protocol-buffers](https://medium.com/@shijuvar/building-high-performance-apis-in-go-using-grpc-and-protocol-buffers-2eda5b80771b)
- [Understanding the context](http://p.agnihotry.com/post/understanding_the_context_package_in_golang)
- [Context cancellation in Go](https://www.sohamkamani.com/blog/golang/2018-06-17-golang-using-context-cancellation/)
- [Docker and Go Modules](https://dev.to/plutov/docker-and-go-modules-3kkn)
