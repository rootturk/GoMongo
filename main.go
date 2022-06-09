package main

import (
	"time"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	)

type MongoConnection struct {
	Client
	Db
}

var mg MongoConnection

const db = "animalarchive"
const mongoURL = "mongodb://localhost:27017" + db

type Animals struct{
	ID string
 	Name string
 	Location string
 	ImageUrl string
 	Kind string
 	CreatedOn time.Time
}

func Connect() error {
	mongo.NewClient
}

func main(){
	if err := Connect(); err != nil{
		log.Fatal(err)
	}

	app := fiber.New()

	app.Get("/animal")
	app.Put("/animal/:id")
	app.Delete("/animal/:id")
}