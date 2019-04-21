package main

import (
	"fmt"

	// Import the generated protobuf code
	pb "github.com/grpc-ms/consignment-service/proto/consignment"

	micro "github.com/micro/go-micro"
	"golang.org/x/net/context"
)

const (
	port = ":50051"
)

type IRepository interface {
	Create(*pb.Consignment) (*pb.Consignment, error)
	GetAll() []*pb.Consignment
}

// Repository - Dummy datastore
type Repository struct {
	consignments []*pb.Consignment
}

func (repo *Repository) Create(consignment *pb.Consignment) (*pb.Consignment, error) {
	updated := append(repo.consignments, consignment)
	repo.consignments = updated
	return consignment, nil
}

func (repo *Repository) GetAll() []*pb.Consignment {
	return repo.consignments
}

// Service - methods defined in protobuf
type service struct {
	repo IRepository
}

// CreateConsignment - takes context & request as argument
func (s *service) CreateConsignment(ctx context.Context, req *pb.Consignment, res *pb.Response) error {

	// Save our consignment
	consignment, err := s.repo.Create(req)
	if err != nil {
		return err
	}

	// Return matching the `Response` message by protobuf definition.
	res.Created = true
	res.Consignment = consignment
	return nil
}

func (s *service) GetConsignments(ctx context.Context, req *pb.GetRequest, res *pb.Response) error {
	consignments := s.repo.GetAll()
	res.Consignments = consignments
	return nil
}

func main() {

	repo := &Repository{}

	//create new srv
	srv := micro.NewService(
		micro.Name("consignment.service"),
	)

	srv.Init()

	// Register handler
	pb.RegisterShippingServiceHandler(srv.Server(), &service{repo})

	//Run server
	if err := srv.Run(); err != nil {
		fmt.Println(err)
	}
}
