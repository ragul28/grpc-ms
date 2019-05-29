package main

import (
	"context"
	"encoding/json"
	"flag"
	"io/ioutil"
	"log"
	"os"

	pb "github.com/grpc-ms/consignment-service/proto/consignment"
	userpb "github.com/grpc-ms/user-service/proto/user"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

const (
	consignmentAddress = "localhost:50051"
	userAddress        = "localhost:50053"
)

var filename = flag.String("filename", "consignment.json", "the consignment json file")
var name = flag.String("name", "Tony Stark", "name of person")
var email = flag.String("email", "Tony.stark@email.com", "email of person")
var password = flag.String("password", "iorn@man", "password")
var company = flag.String("company", "Stark Industries", "company")

func main() {

	ConsignmentAddress := getEnv("CONSIGMENT_HOST", consignmentAddress)
	UserAddress := getEnv("USER_HOST", userAddress)

	flag.Parse()

	// setup connection to userservice
	conn, err := grpc.Dial(UserAddress, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Did not connect: %v", err)
	}
	defer conn.Close()
	userClient := userpb.NewUserServiceClient(conn)

	token := userAuthCli(userClient)

	// setup connection to consignmentservice
	conn, err = grpc.Dial(ConsignmentAddress, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Did not connect: %v", err)
	}
	defer conn.Close()
	consmtClient := pb.NewShippingServiceClient(conn)

	err = consmtSrvCli(consmtClient, token)
	if err != nil {
		log.Fatalf("consignments client failed: %v", err)
	}
}

func userAuthCli(client userpb.UserServiceClient) (token string) {

	log.Println(*name, *email, *password)

	r, err := client.Create(context.TODO(), &userpb.User{
		Name:     *name,
		Email:    *email,
		Password: *password,
		Company:  *company,
	})
	if err != nil {
		log.Fatalf("Could not create: %v", err)
	}
	log.Printf("Created: %s", r.User.Id)

	getAll, err := client.GetAll(context.Background(), &userpb.Request{})
	if err != nil {
		log.Fatalf("Could not list users: %v", err)
	}
	for _, v := range getAll.Users {
		log.Println(v)
	}

	authRes, err := client.Auth(context.TODO(), &userpb.User{
		Email:    *email,
		Password: *password,
	})

	if err != nil {
		log.Fatalf("Could't authenticate user: %s error: %v\n", *email, err)
	}

	log.Printf("Your access token: %s\n", authRes.Token)
	return authRes.Token
}

func consmtSrvCli(client pb.ShippingServiceClient, token string) (err error) {
	consignment, err := parseFile(*filename)
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
		log.Printf("Could not list consignments: %v", err)
		return err
	}
	for _, v := range getAll.Consignments {
		log.Println(v)
	}
	return nil
}

// Getenv Helper func
func getEnv(key, fallback string) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		return fallback
	}
	return value
}

// parseFile json
func parseFile(file string) (*pb.Consignment, error) {
	var consignment *pb.Consignment
	data, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}
	json.Unmarshal(data, &consignment)
	return consignment, err
}
