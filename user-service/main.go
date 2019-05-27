package main

import (
	"fmt"
	"log"

	pb "github.com/grpc-ms/user-service/proto/user"
	"github.com/micro/go-micro"
)

func main() {

	db, err := CreateConnection()
	defer db.Close()

	if err != nil {
		log.Fatalf("Database not connected: %v", err)
	}

	// auto migrate user struct to db
	db.AutoMigrate(&pb.User{})

	repo := &UserRepository{db}
	tokenService := &TokenService{repo}

	srv := micro.NewService(
		micro.Name("go.micro.srv.user"),
	)

	srv.Init()

	pb.RegisterUserServiceHandler(srv.Server(), &handler{repo, tokenService})

	if err := srv.Run(); err != nil {
		fmt.Println(err)
	}
}
