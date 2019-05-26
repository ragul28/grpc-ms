package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	pb "github.com/grpc-ms/vessel-service/proto/vessel"
)

const (
	Port   = ":50052"
	DBHost = "mongodb://localhost:27017"
)

func main() {

	Port := getEnv("GRPC_PORT", Port)
	DBHost := getEnv("DB_HOST", DBHost)

	lis, err := net.Listen("tcp", Port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()

	//create mongodb client
	client, err := CreateClient(DBHost)
	if err != nil {
		log.Panic(err)
	}

	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("DB connected at", DBHost)
	defer client.Disconnect(context.TODO())

	vesselCollection := client.Database("grpc-ms").Collection("vessel")
	repository := &VesselRepository{
		vesselCollection,
	}

	createDummyData(repository)
	// checkFindAv(repository)
	pb.RegisterVesselServiceServer(s, &handler{repository})

	reflection.Register(s)

	log.Println("Running on port:", Port)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

//Ceate dummy vessl data
func createDummyData(repo repository) {
	vessels := []*pb.Vessel{
		{Id: "vessel001", Name: "Kane's Salty Secret", MaxWeight: 200000, Capacity: 500},
	}
	for _, v := range vessels {
		repo.Create(v)
	}
}

// Debug function for Findav Vessel
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

// Getenv Helper func
func getEnv(key, fallback string) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		return fallback
	}
	return value
}
