package main

import (
	"context"
	"log"
	"os"

	microclient "github.com/micro/go-micro/client"
	"github.com/micro/go-micro/cmd"

	pb "github.com/grpc-ms/user-service/proto/user"
)

func main() {

	cmd.Init()

	client := pb.NewUserServiceClient("gomicro.user.service", microclient.DefaultClient)

	name := "Tony Stark"
	email := "Tony.stark@email.com"
	password := "iorn@man"
	company := "Stark Industries"

	log.Println(name, email, password)

	r, err := client.Create(context.TODO(), &pb.User{
		Name:     name,
		Email:    email,
		Password: password,
		Company:  company,
	})
	if err != nil {
		log.Fatalf("Could not create: %v", err)
	}
	log.Printf("Created: %s", r.User.Id)

	getAll, err := client.GetAll(context.Background(), &pb.Request{})
	if err != nil {
		log.Fatalf("Could not list users: %v", err)
	}
	for _, v := range getAll.Users {
		log.Println(v)
	}

	os.Exit(0)
}
