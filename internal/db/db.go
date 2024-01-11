package db

import (
  "context"
  "os"
  "log"

  "go.mongodb.org/mongo-driver/bson"
  "go.mongodb.org/mongo-driver/mongo"
  "go.mongodb.org/mongo-driver/mongo/options"
)

var Client *mongo.Client

func InitDB() {
  serverAPI := options.ServerAPI(options.ServerAPIVersion1)
  opts := options.Client().ApplyURI(os.Getenv("DB_URL")).SetServerAPIOptions(serverAPI)

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
      Keys: bson.D{{Key: "username", Value: 1}},
      Options: options.Index().SetUnique(true),
    },
  ); err != nil {
    log.Println(err)
  }
}

func GetCollection(collectionName string) *mongo.Collection {
  return Client.Database(os.Getenv("DB_NAME")).Collection(collectionName)
}
