package main

import (
	"context"
	"encoding/json"
	"flag"
	"io/ioutil"
	"log"

	microclient "github.com/micro/go-micro/client"
	"github.com/micro/go-micro/metadata"

	pb "github.com/grpc-ms/consignment-service/proto/consignment"
	userpb "github.com/grpc-ms/user-service/proto/user"
)

var filename = flag.String("filename", "consignment.json", "the consignment json file")
var name = flag.String("name", "Tony Stark", "name of person")
var email = flag.String("email", "Tony.stark@email.com", "email of person")
var password = flag.String("password", "iorn@man", "password")
var company = flag.String("company", "Stark Industries", "company")

func main() {

	flag.Parse()

	// setup connection to userservice
	userclient := userpb.NewUserServiceClient("go.micro.srv.user", microclient.DefaultClient)
	token := userAuthCli(userclient)
	log.Println("token: ", token)

	// setup connection to consignment srv
	consmtclient := pb.NewShippingServiceClient("go.micro.srv.consignment", microclient.DefaultClient)
	err := consmtSrvCli(consmtclient, token)
	if err != nil {
		log.Fatalf("consignment client faild: %v", err)
	}
}

// Helper userAuthCli func to get token
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

// Helper func to call consignment srv
func consmtSrvCli(client pb.ShippingServiceClient, token string) (err error) {
	consignment, err := parseFile(*filename)
	if err != nil {
		log.Fatalf("Could not parse file: %v", err)
	}

	ctx := metadata.NewContext(context.Background(), map[string]string{
		"token": token,
	})

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

func parseFile(file string) (*pb.Consignment, error) {
	var consignment *pb.Consignment
	data, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}
	json.Unmarshal(data, &consignment)
	return consignment, err
}
