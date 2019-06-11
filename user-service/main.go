package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

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

	reflection.Register(s)

	go runHttp(fmt.Sprintf("localhost%s", Port))

	log.Println("Running on port:", Port)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func runHttp(clientAddr string) {

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithInsecure()}
	if err := pb.RegisterUserServiceHandlerFromEndpoint(ctx, mux, clientAddr, opts); err != nil {
		log.Fatalf("Failed to start HTTP server: %v", err)
	}

	http.ListenAndServe(":6000", mux)
}
