package main

import (
	"log"
	"net"
	"os"

	"google.golang.org/grpc"

	pb "github.com/grpc-ms/user-service/proto/user"
)

const (
	defaultPort = ":50053"
)

func main() {

	Port := os.Getenv("GRPC_PORT")
	if Port == "" {
		Port = defaultPort
	}
	lis, err := net.Listen("tcp", Port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()

	db, err := CreateConnection()
	defer db.Close()

	if err != nil {
		log.Fatalf("Database not connected: %v", err)
	} else {
		log.Println("Connected to Postgras DB!")
	}

	// auto migrate user struct to db
	db.AutoMigrate(&pb.User{})

	repo := &UserRepository{db}
	tokenService := &TokenService{repo}

	pb.RegisterUserServiceServer(s, &handler{repo, tokenService})

	log.Println("Running on port:", Port)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
