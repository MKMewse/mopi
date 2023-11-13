package server

import (
	"context"
	"fmt"
	"log/slog"
	"reflect"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoStorer struct {
	Client *mongo.Client
}

func NewMongoStorer(uri string) (*MongoStorer, error) {
	tM := reflect.TypeOf(bson.M{})
	rb := bson.NewRegistry()
	rb.RegisterTypeMapEntry(bson.TypeEmbeddedDocument, tM)

	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri).SetRegistry(rb))
	if err != nil {
		slog.Error("failed to connect to mongo")
		return nil, err
	}
	return &MongoStorer{
		Client: client,
	}, nil
}

func (m *MongoStorer) GetAll() []Response {
	c := m.Client.Database("mopi").Collection("responses")
	cur, err := c.Find(context.TODO(), bson.D{{}})
	if err != nil {
		slog.Error(fmt.Sprintf("error getting responses from mongo: %s", err.Error()))
	}
	var responses []Response
	cur.All(context.TODO(), &responses)
	return responses
}

func (m *MongoStorer) Add(r Response) error {
	c := m.Client.Database("mopi").Collection("responses")
	_, err := c.InsertOne(context.TODO(), r)
	if err != nil {
		slog.Error(fmt.Sprintf("failed to insert new response into mongo %s", err.Error()))
		return err
	}
	return nil
}

func (m *MongoStorer) RemoveAll() error {
	c := m.Client.Database("mopi").Collection("responses")
	_, err := c.DeleteMany(context.TODO(), bson.D{{}})
	if err != nil {
		slog.Error(fmt.Sprintf("failed to delete responses in mongo: %s", err.Error()))
		return err
	}
	return nil
}
