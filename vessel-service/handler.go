package main

import (
	"context"

	pb "github.com/grpc-ms/vessel-service/proto/vessel"
)

type handler struct {
	repository
}

// FindAvailable vessels
func (s *handler) FindAvailable(ctx context.Context, req *pb.Specification) (*pb.Response, error) {

	// Find the next available vessel
	vessel, err := s.repository.FindAvailable(req)
	if err != nil {
		return nil, err
	}

	// Set the vessel as part of the response message type
	return &pb.Response{Vessel: vessel}, nil
}

// Create a new vessel
func (s *handler) Create(ctx context.Context, req *pb.Vessel) (*pb.Response, error) {
	if err := s.repository.Create(req); err != nil {
		return nil, err
	}

	return &pb.Response{Vessel: req, Created: true}, nil
}
