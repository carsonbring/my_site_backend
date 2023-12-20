package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Data model
type BlogPost struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

type Config struct {
	Mongo struct {
		ConnectionString string `json:"connectionString"`
	} `json:"Mongo"`
}

func main() {
	app := fiber.New()

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})
	app.Post("/blog", func(c *fiber.Ctx) error {
		if err := createBlogPost(c); err != nil {
			log.Println("Error in createBlogPost:", err)

		}
		return c.SendString("Blog post created successfully")
	})
	if err := app.Listen(":3000"); err != nil {
		log.Fatal(err)
	}

}
func LoadConfig() (*Config, error) {
	file, err := os.Open("config.json")
	if err != nil {
		return nil, err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	config := &Config{}
	err = decoder.Decode(config)
	if err != nil {
		return nil, err
	}

	return config, nil
}
func connectToMongoDB() (*mongo.Client, error) {
	config, err := LoadConfig()
	if err != nil {
		return nil, err
	}
	serverAPIOptions := options.ServerAPI(options.ServerAPIVersion1)
	clientOptions := options.Client().ApplyURI(config.Mongo.ConnectionString).
		SetServerAPIOptions(serverAPIOptions)
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		return nil, err
	}

	err = client.Ping(context.Background(), nil)
	if err != nil {
		return nil, err
	}

	fmt.Println("Connected to MongoDB")
	return client, nil
}

func createBlogPost(c *fiber.Ctx) error {
	//parsing request body
	var post BlogPost
	if err := c.BodyParser(&post); err != nil {
		return err
	}

	//connecting to MongoDB
	client, err := connectToMongoDB()
	if err != nil {
		log.Fatal(err)
		return err
	}
	defer client.Disconnect(context.Background())

	//if the connection successful, then insert post into db
	collection := client.Database("cbDB").Collection("blog_posts")
	_, err = collection.InsertOne(context.Background(), post)
	if err != nil {
		return err
	}

	return c.JSON(post)
}
