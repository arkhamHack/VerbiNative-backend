package configs

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/go-redis/redis"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	rdb *redis.Client
)

func ConnectDB() *mongo.Client {
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatal("Error loading .env file")
	}
	client, err := mongo.NewClient(options.Client().ApplyURI(EnvMongoURI()))
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

// func RedisConn() *redis.Client {
// 	redisURL := os.Getenv("REDIS_URL")
// 	if redisURL == "" {
// 		log.Fatal("REDIS_URL environment variable not set")
// 	}
// 	opt, err := redis.ParseURL(redisURL)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	rdb := redis.NewClient(opt)

// 	_, err = rdb.Ping().Result()
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	return rdb
// }
