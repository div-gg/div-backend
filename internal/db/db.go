package db

import (
	"context"
	"log"

  "github.com/divinitymn/div-backend/internal/config"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var Client *mongo.Client

func InitDB() {
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(config.Env.DbURL).SetServerAPIOptions(serverAPI)

	client, err := mongo.Connect(context.Background(), opts)
	if err != nil {
		log.Fatal("⇨ Error connecting to MongoDB: ", err)
		log.Fatal(err)
	}

	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal("⇨ Error connecting to MongoDB: ", err)
		log.Fatal(err)
	}

	log.Println("⇨ MongoDB Connected!")

	Client = client

	// Indexes
	if _, err := GetCollection("users").Indexes().CreateOne(
		context.Background(),
		mongo.IndexModel{
			Keys:    bson.D{{Key: "username", Value: 1}},
			Options: options.Index().SetUnique(true),
		},
	); err != nil {
		log.Println(err)
	}

	if _, err := GetCollection("posts").Indexes().CreateOne(
		context.Background(),
		mongo.IndexModel{
			Keys:    bson.D{{Key: "expire_at", Value: 1}},
			Options: options.Index().SetExpireAfterSeconds(0),
		},
	); err != nil {
		log.Println(err)
	}
}

func GetCollection(collectionName string) *mongo.Collection {
	return Client.Database(config.Env.DbName).Collection(collectionName)
}
