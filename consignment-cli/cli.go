package main

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"os"

	pb "github.com/grpc-ms/consignment-service/proto/consignment"
	micro "github.com/micro/go-micro"
	"github.com/micro/go-micro/metadata"
)

const (
	address         = "localhost:50051"
	defaultFilename = "consignment.json"
	defaultToken    = "jwttokenhere"
)

func parseFile(file string) (*pb.Consignment, error) {
	var consignment *pb.Consignment
	data, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}
	json.Unmarshal(data, &consignment)
	return consignment, err
}

func main() {
	service := micro.NewService(micro.Name("go.micro.consignment.cli"))
	service.Init()

	client := pb.NewShippingServiceClient("go.micro.srv.consignment", service.Client())

	// Contact the server and print out its response.
	file := defaultFilename
	token := defaultToken
	if len(os.Args) > 1 {
		file = os.Args[1]
		token = os.Args[2]
	}

	consignment, err := parseFile(file)

	if err != nil {
		log.Fatalf("Could not parse file: %v", err)
	}

	ctx := metadata.NewContext(context.Background(), map[string]string{
		"token": token,
	})

	r, err := client.CreateConsignment(ctx, consignment)
	if err != nil {
		log.Fatalf("Could not greet: %v", err)
	}
	log.Printf("Created: %t", r.Created)

	getAll, err := client.GetConsignments(ctx, &pb.GetRequest{})
	if err != nil {
		log.Fatalf("Could not list consignments: %v", err)
	}
	for _, v := range getAll.Consignments {
		log.Println(v)
	}
}
