package main

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"os"

	"google.golang.org/grpc/metadata"

	pb "github.com/grpc-ms/consignment-service/proto/consignment"
	"google.golang.org/grpc"
)

const (
	consignmentAddress = "localhost:50051"
	defaultFilename    = "consignment.json"
	defaultToken       = "secret-token"
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

	ConsignmentAddress := os.Getenv("CONSIGMENT_HOST")
	if ConsignmentAddress == "" {
		ConsignmentAddress = consignmentAddress
	}
	// Set up a connection to the server.
	conn, err := grpc.Dial(ConsignmentAddress, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Did not connect: %v", err)
	}
	defer conn.Close()
	client := pb.NewShippingServiceClient(conn)

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

	ctx := metadata.NewOutgoingContext(context.Background(), metadata.New(map[string]string{
		"token": token,
	}))

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
