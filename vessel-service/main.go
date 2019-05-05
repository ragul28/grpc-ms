package main

import (
	"context"
	"fmt"
	"log"
	"os"

	pb "github.com/grpc-ms/vessel-service/proto/vessel"
	"github.com/micro/go-micro"
)

const (
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
	srv := micro.NewService(
		micro.Name("vessel.service"),
	)

	srv.Init()

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

	pb.RegisterVesselServiceHandler(srv.Server(), &handler{repository})

	if err := srv.Run(); err != nil {
		fmt.Println(err)
	}
}
