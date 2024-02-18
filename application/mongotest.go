package application

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/brianvoe/gofakeit"
	gofakeit "github.com/brianvoe/gofakeit/v7"
	uuid "github.com/satori/go.uuid"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type FakeCar struct {
	UUID5       string
	PlateNumber string   `fake:"{regex:[A-Z]{3}}-{regex:[0-9]{3}}"`
	ValidUntil  string   `fake:"{year}-{month}-{day}" format:"2006-01-02"`
	Owner       string   `fake:"{name}"`
	Data        []string `fakesize:"3"`
}

type Cars struct {
	Id          primitive.ObjectID `bson:"_id,omitempty"`
	UUID5       string             `bson:"UUID5,omitempty"`
	PlateNumber string             `bson:"plateNumber,omitempty"`
	Owner       string             `bson:"owner,omitempty"`
	ValidUntil  string             `bson:"validUntil,omitempty"`
	Data        []string           `bson:"data,omitempty"`
}

// var collection *mongo.Collection
func MongoTest() {

	getFakeCar()
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")

	// Connect to MongoDB
	client, err := mongo.Connect(ctx, clientOptions)

	if err != nil {
		log.Fatal(err)
	}

	defer client.Disconnect(ctx)
	// Check the connection
	err = client.Ping(context.TODO(), nil)

	if err != nil {
		log.Fatal(err)
	}

	database := client.Database("darth")
	carsCollection := database.Collection("Cars")

	namespaceUUID := uuid.NewV4()
	name := "unique_name_for_car1"
	uuidV5 := uuid.NewV5(namespaceUUID, name)

	car1 := Cars{
		UUID5:       uuidV5.String(),
		PlateNumber: "ROBOTD-04",
		Owner:       "Bill Gates",
		ValidUntil:  "2026-02-28",
		Data:        []string{"zöld", "VIN: WP0ZZZ99ZTS392124"},
	}

	fmt.Print(car1)
	response, _ := carsCollection.InsertOne(ctx, car1)
	fmt.Print(response)

	fmt.Println("Connected to MongoDB!")
}

func getFakeCar() FakeCar {
	var fakeCar FakeCar
	err := gofakeit.Struct(&fakeCar)
	if err != nil {
		log.Fatal(err)
	}

	namespaceUUID := uuid.NewV4()
	name := fakeCar.Owner
	uuidV5 := uuid.NewV5(namespaceUUID, name)

	fakeCar.UUID5 = uuidV5.String()
	fakeCar.Data[0] = gofakeit.Color()
	log.Fatal(fakeCar)
	return fakeCar
}
