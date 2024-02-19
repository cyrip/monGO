package application

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"time"

	gofakeit "github.com/brianvoe/gofakeit/v7"
	uuid "github.com/satori/go.uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type FakeCar struct {
	UUID5       string   `fake:"{uuid}"`
	PlateNumber string   `fake:"{regex:[A-Z]{3}}-{regex:[0-9]{3}}"`
	ValidUntil  string   `fake:"{year}-{month}-{day}" format:"2006-01-02"`
	Owner       string   `fake:"{name}"`
	Data        []string `fakesize:"3"`
}

type Cars struct {
	Id          primitive.ObjectID `bson:"_id,omitempty"`
	UUID5       string             `bson:"uuid5,omitempty"`
	PlateNumber string             `bson:"platenumber,omitempty"`
	Owner       string             `bson:"owner,omitempty"`
	ValidUntil  string             `bson:"validuntil,omitempty"`
	Data        []string           `bson:"data,omitempty"`
}

var carsConnection0 *mongo.Client
var carsCollection0 *mongo.Collection

var ctx context.Context

// var collection *mongo.Collection
func MongoTest() {

	// getFakeCar()
	ctx, _ = context.WithTimeout(context.Background(), 10*time.Second)
	clientOptions := options.Client().ApplyURI(fmt.Sprintf("mongodb://localhost:%d", 27018))

	var err error
	// Connect to MongoDB
	carsConnection0, err = mongo.Connect(ctx, clientOptions)

	if err != nil {
		log.Fatal(err)
	}

	defer carsConnection0.Disconnect(ctx)
	// Check the connection
	err = carsConnection0.Ping(context.TODO(), nil)

	if err != nil {
		log.Fatal(err)
	}

	database := carsConnection0.Database("darth")
	carsCollection0 = database.Collection("Cars")

	//	namespaceUUID := uuid.NewV4()
	//name := "unique_name_for_car1"
	//uuidV5 := uuid.NewV5(namespaceUUID, name)

	//car1 := Cars{
	//UUID5:       uuidV5.String(),
	//PlateNumber: "ROBOTD-04",
	//Owner:       "Bill Gates",
	//ValidUntil:  "2026-02-28",
	//Data:        []string{"zöld", "VIN: WP0ZZZ99ZTS392124"},
	//}

	//fmt.Print(car1)
	//response, _ := carsCollection.InsertOne(ctx, car1)
	//fmt.Print(response)
	insertFakeCars()
	addIndex()
	// insertOne("KQA-901")

	fmt.Println("Connected to MongoDB!")
}

func insertOne(plateNumber string) {
	namespaceUUID := uuid.NewV4()
	uuidV5 := uuid.NewV5(namespaceUUID, plateNumber)

	car1 := Cars{
		UUID5:       uuidV5.String(),
		PlateNumber: plateNumber,
		Owner:       "Bill Gates",
		ValidUntil:  "2026-02-28",
		Data:        []string{"zöld", "VIN: WP0ZZZ99ZTS392124"},
	}

	fmt.Print(car1)
	response, err := carsCollection0.InsertOne(ctx, car1)
	if err != nil {
		fmt.Println(response)
		fmt.Println(err)
	}
}

func insertFakeCars() {
	gofakeit.Seed(rand.Intn(10))
	inserted := 0
	for i := 0; i < 10000; i++ {
		car := getFakeCar()
		response, err := carsCollection0.InsertOne(ctx, car)
		if err != nil {
			fmt.Println(err)
			fmt.Println(response)
		} else {
			inserted = inserted + 1
		}
		// fmt.Print(response)
	}
	log.Printf("Inserted docs %d", inserted)
}

func getFakeCar() FakeCar {
	var fakeCar FakeCar
	err := gofakeit.Struct(&fakeCar)
	if err != nil {
		log.Fatal(err)
	}

	fakeCar.Data[0] = gofakeit.Color()
	fakeCar.Data[1] = gofakeit.City()
	fakeCar.Data[2] = gofakeit.BeerName()
	return fakeCar
}

func addIndex() {
	indexModel := mongo.IndexModel{
		Keys:    bson.D{{Key: "platenumber", Value: 1}},
		Options: options.Index().SetUnique(true),
	}

	// Create the index
	indexName, err := carsCollection0.Indexes().CreateOne(ctx, indexModel)
	if err != nil {
		log.Println("Index already exists")
	} else {
		log.Printf("Index added %s", indexName)
	}
}

func getMongoShard(plateNumber string) int {
	firstChar := plateNumber[0]
	return int(firstChar) % 2
}
