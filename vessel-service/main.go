package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"

	pb "github.com/grpc-ms/vessel-service/proto/vessel"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const (
	port        = ":50052"
	defaultHost = "localhost:27017"
)

func createDummyData(repo repository) {
	vessels := []*pb.Vessel{
		{Id: "vessel001", Name: "Kane's Salty Secret", MaxWeight: 200000, Capacity: 500},
	}
	for _, v := range vessels {
		repo.Create(v)
	}
}

func checkFindAv(repo repository) {

	vesselRes, err := repo.FindAvailable(&pb.Specification{
		MaxWeight: int32(55000),
		Capacity:  int32(3),
	})
	log.Printf("Found vessel: %+v\n", vesselRes)

	if err != nil {
		fmt.Println(err)
	}
}

func main() {

	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()

	uri := os.Getenv("DB_HOST")
	if uri == "" {
		uri = "mongodb://" + defaultHost
	}

	client, err := CreateClient(uri)
	if err != nil {
		log.Panic(err)
	}

	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("DB connected at", uri)
	defer client.Disconnect(context.TODO())

	vesselCollection := client.Database("grpc-ms").Collection("vessel")
	repository := &VesselRepository{
		vesselCollection,
	}

	createDummyData(repository)
	//checkFindAv(repository)
	pb.RegisterVesselServiceServer(s, &handler{repository})

	reflection.Register(s)

	log.Println("Running on port:", port)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
