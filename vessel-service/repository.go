package main

import (
	"context"

	pb "github.com/ragul28/grpc-ms/vessel-service/proto/vessel"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type repository interface {
	FindAvailable(spec *pb.Specification) (*pb.Vessel, error)
	Create(vessel *pb.Vessel) error
}

type VesselRepository struct {
	collection *mongo.Collection
}

//FindAvailable - checks weight spec
// capacity (vessel) >= spec.capacity (tot.consignment cap)
func (repository *VesselRepository) FindAvailable(spec *pb.Specification) (*pb.Vessel, error) {
	filter := bson.D{
		{"capacity", bson.D{{"$gte", spec.Capacity}}},
		{"maxweight", bson.D{{"$gte", spec.MaxWeight}}},
	}

	//filter := bson.D{{"capacity", bson.D{{"$gte", spec.Capacity}}}}
	//filter := bson.M{"capacity": bson.M{"$lte": 5}}
	// filter := bson.D{{"capacity", 500}}
	var vessel *pb.Vessel
	//db.vessel.find( { capacity: { $gt: 1 }, maxweight: { $gt: 55000 } } )
	if err := repository.collection.FindOne(context.TODO(), filter).Decode(&vessel); err != nil {
		return nil, err
	}
	return vessel, nil
}

// Create a new vessel
func (repository *VesselRepository) Create(vessel *pb.Vessel) error {
	_, err := repository.collection.InsertOne(context.TODO(), vessel)
	return err
}
