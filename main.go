package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/go-redis/redis/v8"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	mongoClient *mongo.Client
	redisClient *redis.Client
)

func main() {
	// Initialize MongoDB client
	mongoURI := "mongodb://mongodb:27017"
	mongoClient = connectMongoDB(mongoURI)
	defer mongoClient.Disconnect(context.Background())

	// Initialize Redis client
	redisClient = redis.NewClient(&redis.Options{
		Addr:     "redis:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	defer redisClient.Close()

	// Test MongoDB connection
	testMongoDBConnection()

	// Test Redis connection
	testRedisConnection()

	// Set up HTTP server
	http.HandleFunc("/", handler)
	fmt.Println("Server listening on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, World!")
}

func connectMongoDB(uri string) *mongo.Client {
	client, err := mongo.NewClient(options.Client().ApplyURI(uri))
	if err != nil {
		log.Fatalf("Error creating MongoDB client: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = client.Connect(ctx)
	if err != nil {
		log.Fatalf("Error connecting to MongoDB: %v", err)
	}

	fmt.Println("Connected to MongoDB")
	return client
}

func testMongoDBConnection() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := mongoClient.Ping(ctx, nil)
	if err != nil {
		log.Fatalf("Error pinging MongoDB: %v", err)
	}

	fmt.Println("MongoDB connection successful")
}

func testRedisConnection() {
	ctx := context.Background()

	pong, err := redisClient.Ping(ctx).Result()
	if err != nil {
		log.Fatalf("Error pinging Redis: %v", err)
	}

	fmt.Println("Redis Ping Response:", pong)
}
