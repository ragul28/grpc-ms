package main

import (
	"context"
	"log"
	"os"

	pb "github.com/grpc-ms/user-service/proto/user"
	"google.golang.org/grpc"
)

const (
	userAddress = "localhost:50053"
)

func main() {

	UserAddress := os.Getenv("USER_HOST")
	if UserAddress == "" {
		UserAddress = userAddress
	}
	conn, err := grpc.Dial(UserAddress, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Did not connect: %v", err)
	}
	defer conn.Close()
	client := pb.NewUserServiceClient(conn)

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

	authRes, err := client.Auth(context.TODO(), &pb.User{
		Email:    email,
		Password: password,
	})

	if err != nil {
		log.Fatalf("Could't authenticate user: %s error: %v\n", email, err)
	}

	log.Printf("Your access token: %s\n", authRes.Token)

	os.Exit(0)
}
