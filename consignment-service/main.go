package main

import (
	"context"
	"fmt"
	"log"
	"os"

	// Import the generated protobuf code
	pb "github.com/grpc-ms/consignment-service/proto/consignment"
	vesselProto "github.com/grpc-ms/vessel-service/proto/vessel"
	micro "github.com/micro/go-micro"
)

const (
	port        = ":50051"
	defaultHost = "localhost:27017"
)

func main() {
	// setup go-micro
	srv := micro.NewService(
		micro.Name("consignment.service"),
	)

	srv.Init()

	uri := os.Getenv("DB_HOST")
	if uri == "" {
		uri = defaultHost
	}

	client, err := createClient(uri)
	if err != nil {
		log.Panic(err)
	}
	defer client.Disconnect(context.TODO())

	consignmentCollection := client.Database("grpc-ms").Collection("consignments")

	repository := &MongoRepository{consignmentCollection}
	vesselClient := vesselProto.NewVesselServiceClient("vessel.service", srv.Client())
	h := &handler{repository, vesselClient}

	// Register handler
	pb.RegisterShippingServiceHandler(srv.Server(), h)

	//Run server
	if err := srv.Run(); err != nil {
		fmt.Println(err)
	}
}
