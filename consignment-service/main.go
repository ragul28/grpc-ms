package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"

	// Import the generated protobuf code
	pb "github.com/grpc-ms/consignment-service/proto/consignment"
	userProto "github.com/grpc-ms/user-service/proto/user"
	vesselProto "github.com/grpc-ms/vessel-service/proto/vessel"
	"github.com/micro/go-micro"
	"github.com/micro/go-micro/client"
	"github.com/micro/go-micro/metadata"
	"github.com/micro/go-micro/server"
)

const (
	defaultHost = "localhost:27017"
)

func main() {
	// setup go-micro
	srv := micro.NewService(
		micro.Name("go.micro.srv.consignment"),
		micro.Version("latest"),
		//auth middleware
		micro.WrapHandler(AuthWrapper),
	)

	srv.Init()

	//Get database host
	uri := os.Getenv("DB_HOST")
	if uri == "" {
		uri = "mongodb://" + defaultHost
	}

	//create mongodb client session
	session, err := createClient(uri)
	if err != nil {
		log.Panic(err)
	}

	log.Println("DB connected at", uri)
	defer session.Disconnect(context.TODO())

	consignmentCollection := session.Database("grpc-ms").Collection("consignments")

	repository := &MongoRepository{consignmentCollection}
	vesselClient := vesselProto.NewVesselServiceClient("go.micro.srv.vessel", srv.Client())
	h := &handler{repository, vesselClient}

	// Register handler
	pb.RegisterShippingServiceHandler(srv.Server(), h)

	//Run server
	if err := srv.Run(); err != nil {
		fmt.Println(err)
	}
}

//Auth middleware to validate token in consignment svc api
func AuthWrapper(fn server.HandlerFunc) server.HandlerFunc {
	return func(ctx context.Context, req server.Request, resp interface{}) error {
		//To skip auth for dev
		DISABLE_AUTH := false
		if os.Getenv("DISABLE_AUTH") == "true" || DISABLE_AUTH {
			return fn(ctx, req, resp)
		}

		meta, ok := metadata.FromContext(ctx)
		if !ok {
			return errors.New("no auth meta-data found in request")
		}

		token := meta["Token"]
		log.Println("Authenticating with token: ", token)

		authClient := userProto.NewUserServiceClient("go.micro.srv.user", client.DefaultClient)
		authResp, err := authClient.ValidateToken(ctx, &userProto.Token{
			Token: token,
		})
		log.Println("Auth resp:", authResp)

		if err != nil {
			log.Println("Err:", err)
			return err
		}

		err = fn(ctx, req, resp)
		return err
	}
}
