package main

import (
	"context"
	"log"

	pb "github.com/ragul28/grpc-ms/consignment-service/proto/consignment"
	"go.mongodb.org/mongo-driver/bson"
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
	cur, err := repository.collection.Find(context.TODO(), bson.D{{}}, nil)
	if err != nil {
		log.Println(err)
	}

	var consignments []*pb.Consignment

	for cur.Next(context.TODO()) {
		var consignment *pb.Consignment
		err := cur.Decode(&consignment)
		if err != nil {
			log.Panicln(err)
			return nil, err
		}
		consignments = append(consignments, consignment)
	}
	return consignments, err
}
