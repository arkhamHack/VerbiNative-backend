package configs

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func ConnectDB() *mongo.Client {
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatal("Error loading .env file")
	}
	client, err := mongo.NewClient(options.Client().ApplyURI(os.Getenv("MONGOURI")))
	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connected to MongoDB")
	return client
}

var Mongo_DB *mongo.Client = ConnectDB()

func GetCollec(client *mongo.Client, collec_name string) *mongo.Collection {
	fmt.Println("Check 1")
	collection := client.Database("golangApi").Collection(collec_name)
	return collection
}

// func GetChats(client *mongo.Client, collec_name string) *mongo.Collection {
// 	fmt.Println("Check 1")
// 	collection := client.Database("chats").Collection(collec_name)
// 	return collection
// }
