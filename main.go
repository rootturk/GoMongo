package main

import (
	"time"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log"
	"context"
	)

type MongoConnection struct {
	Client *mongo.Client
	Db *mongo.Database
}

var mg MongoConnection

const dbStorageName = "animalarchive"
const mongoURL = "mongodb://localhost:27017" + dbStorageName

type Animals struct{
	ID string  `json:"id, omitempty" bson:"_id, omitempty"`
	Name string `json:"name"`
	Description string `json:"description"`
	CreatedOn time.Time `json:"createdOn"`
 	Location string `json:"location"`
 	ImageUrl string `json:"description"`
 	Kind string `json:"description"`
}

func Connect() error {
	client, err := mongo.NewClient(options.Client().ApplyURI(mongoURL))
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)

	defer cancel()

	err = client.Connect(ctx)
	db := client.Database(dbStorageName)

	if err!=nil {
		return err
	}

	mg = MongoConnection {
		Client: client,
		Db: db,
	}

	return nil
}

func main(){

	if err := Connect(); err != nil{
		log.Fatal(err)
	}

	app := fiber.New()

	app.Get("/animal", func(c *fiber.Ctx) error {
		query := bson.D{{}}
		cursor, err := mg.Db.Collection("animals").Find(c.Context(), query)

		if err!=nil {
			return c.Status(500).SendString(err.Error())
		}

		var animals []Animals = make([]Animals, 0)

		if err := cursor.All(c.Context(), &animals); err !=nil {
			return c.Status(500).SendString(err.Error())
		}

		return c.JSON(animals)
	})


	app.Post("/animal", func(c *fiber.Ctx) error {
		collection := mg.Db.Collection("animals")

		animals := new (Animals)

		if err := c.BodyParser(animals); err!=nil {
			return c.Status(400).SendString(err.Error())
		}

		animals.ID = ""

		insertResult, err := collection.InsertOne(c.Context(), animals)

		if err !=nil {
			return c.Status(500).SendString(err.Error())
		}

		filter := bson.D{{Key:"_id", Value: insertResult.InsertedID}}

		createdRecords := collection.FindOne(c.Context(), filter)

		createdAnimal:= &Animals{}

		createdRecords.Decode(createdRecords)

		return c.Status(201).JSON(createdAnimal)
	})

	app.Put("/animal/:id", func(c *fiber.Ctx) error {
		Id := c.Params("id")

		animalID, err := primitive.ObjectIDFromHex(Id)

		if err!=nil{
			return c.SendStatus(400)
		}

		animals := new (Animals)

		if err := c.BodyParser(animals); err != nil{
			return c.Status(400).SendString(err.Error())
		}

		query := bson.D{{Key:"_id", Value: animalID}}
		update :=
		bson.D{
				{Key:"$set",
				Value: bson.D{
					{Key:"name", Value: animals.Name},
					{Key:"description", Value: animals.Description},
					{Key:"location", Value: animals.Location},
					{Key:"imageurl", Value: animals.ImageUrl},
				},
			},
		}

		err = mg.Db.Collection("animals").FindOneAndUpdate(c.Context(), query, update).Err()

		if err != nil {
			if err == mongo.ErrNoDocuments{
				return c.SendStatus(400)
			}

			return c.SendStatus(500)
		}

		animals.ID = Id

		return c.Status(200).JSON(animals)

	})
	
	app.Delete("/animal/:id", func(c* fiber.Ctx) error {
		animalID, err := primitive.ObjectIDFromHex(c.Params("id"))

		if err != nil{
			return c.SendStatus(400)
		}

		query := bson.D{{Key:"_id", Value: animalID}}
		result, err := mg.Db.Collection("animals").DeleteOne(c.Context(), &query)

		if err !=nil {
			return c.SendStatus(500)
		}

		if result.DeletedCount < 1 {
			return c.SendStatus(404)
		}

		return c.Status(200).JSON("record deleted")
	})


}