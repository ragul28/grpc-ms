package main

import (
	"context"
	"log"

	// Import the generated protobuf code
	pb "github.com/grpc-ms/consignment-service/proto/consignment"
	vesselProto "github.com/grpc-ms/vessel-service/proto/vessel"
)

type handler struct {
	repository
	vesselClient vesselProto.VesselServiceClient
}

// CreateConsignment - takes context & request as argument
func (s *handler) CreateConsignment(ctx context.Context, req *pb.Consignment) (*pb.Response, error) {

	// Call vessel client with consignment weight and no. containers
	vesselResponse, err := s.vesselClient.FindAvailable(ctx, &vesselProto.Specification{
		MaxWeight: req.Weight,
		Capacity:  int32(len(req.Containers)),
	})
	if err != nil {
		return nil, err
	}
	log.Printf("Found vessel: %s \n", vesselResponse.Vessel.Name)

	//set VesselID from vesselsrv
	req.VesselId = vesselResponse.Vessel.Id

	// Save our consignment
	if err = s.repository.Create(req); err != nil {
		return nil, err
	}

	// Return `Response` message by protobuf def
	return &pb.Response{Created: true, Consignment: req}, nil
}

// GetConsignments
func (s *handler) GetConsignments(ctx context.Context, req *pb.GetRequest) (*pb.Response, error) {

	consignments, err := s.repository.GetAll()
	if err != nil {
		return nil, err
	}

	return &pb.Response{Consignments: consignments}, nil
}
