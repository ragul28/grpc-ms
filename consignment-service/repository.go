package main

import (
	"context"

	pb "github.com/grpc-ms/consignment-service/proto/consignment"
	"go.mongodb.org/mongo-driver/mongo"
)

type repository interface {
	Create(consignment *pb.Consignment) error
	GetAll() ([]*pb.Consignment, error)
}

// MongoRepository imple
type MongoRepository struct {
	collection *mongo.Collection
}

func (repository *MongoRepository) Create(consignment *pb.Consignment) error {
	_, err := repository.collection.InsertOne(context.Background(), consignment)
	return err
}

func (repository *MongoRepository) GetAll() ([]*pb.Consignment, error) {
	cur, err := repository.collection.Find(context.Background(), nil, nil)

	var consignments []*pb.Consignment

	for cur.Next(context.Background()) {
		var consignment *pb.Consignment
		if err := cur.Decode(&consignment); err != nil {
			return nil, err
		}
		consignments = append(consignments, consignment)
	}
	return consignments, err
}
