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

const MONGO_SHARDS int = 2

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

var connections [MONGO_SHARDS]*mongo.Client
var collections [MONGO_SHARDS]*mongo.Collection
var contexts [MONGO_SHARDS]context.Context

func createConnection(id int) *mongo.Client {
	contexts[id], _ = context.WithTimeout(context.Background(), 10*time.Second)
	clientOptions := options.Client().ApplyURI(fmt.Sprintf("mongodb://localhost:%d", 27017+id))
	var err error
	connections[id], err = mongo.Connect(contexts[id], clientOptions)

	if err != nil {
		log.Fatal(err)
	}

	err = connections[id].Ping(context.TODO(), nil)

	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Connected to mongodb on port %d", 27017+id)
	return connections[id]
}

func createCollection(id int, ctx context.Context) *mongo.Collection {
	database := connections[id].Database("darth")
	collections[id] = database.Collection("Cars")
	return collections[id]
}

// var collection *mongo.Collection
func MongoTest() {

	for id := 0; id <= 1; id++ {
		createConnection(id)
		defer connections[id].Disconnect(contexts[id])
		createCollection(id, contexts[id])
		addIndex(id)
	}

	// insertFakeCars()
	find3(0, "PW.*")
	find3(1, "PW.*")
}

func insertCar(plateNumber string) {
	namespaceUUID := uuid.NewV4()
	uuidV5 := uuid.NewV5(namespaceUUID, plateNumber)

	car := Cars{
		UUID5:       uuidV5.String(),
		PlateNumber: plateNumber,
		Owner:       "Bill Gates",
		ValidUntil:  "2026-02-28",
		Data:        []string{"zöld", "VIN: WP0ZZZ99ZTS392124"},
	}
	insertOne(car)
}

func insertOne(car Cars) {

	shard := getMongoShard(car.PlateNumber)

	response, err := collections[shard].InsertOne(contexts[shard], car)
	if err != nil {
		fmt.Println(response)
		fmt.Println(err)
	}
}

func insertFakeCars() {
	gofakeit.Seed(rand.Intn(10))
	var inserted [MONGO_SHARDS]int
	for i := 0; i < 10000; i++ {
		car := getFakeCar()
		shard := getMongoShard(car.PlateNumber)
		response, err := collections[shard].InsertOne(contexts[shard], car)
		if err != nil {
			fmt.Println(err)
			fmt.Println(response)
		} else {
			inserted[shard] = inserted[shard] + 1
		}
	}
	log.Printf("Inserted docs shard0 %d shard1 %d", inserted[0], inserted[1])
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

// @TODO: index on mongo 1?
func addIndex(id int) {
	indexModel := mongo.IndexModel{
		Keys:    bson.D{{Key: "platenumber", Value: 1}},
		Options: options.Index().SetUnique(true),
	}

	// Create the index
	indexName, err := collections[id].Indexes().CreateOne(contexts[id], indexModel)
	if err != nil {
		log.Println("Index already exists")
		log.Println(err)
	} else {
		log.Printf("Index added %s", indexName)
	}
}

func getMongoShard(plateNumber string) int {
	firstChar := plateNumber[0]
	return int(firstChar) % 2
}

func find3(id int, regex string) {

	filter := bson.M{
		"$or": []interface{}{
			bson.M{"platenumber": bson.M{"$regex": regex, "$options": "i"}},
			bson.M{"owner": bson.M{"$regex": regex, "$options": "i"}},
			bson.M{"data": bson.M{"$regex": regex, "$options": "i"}},
		},
	}

	// Find documents
	cursor, err := collections[id].Find(contexts[id], filter)
	if err != nil {
		log.Fatal(err)
	}
	defer cursor.Close(contexts[id])

	for cursor.Next(contexts[id]) {
		var result bson.M
		if err := cursor.Decode(&result); err != nil {
			log.Fatal(err)
		}
		fmt.Println(result)
	}

	if err := cursor.Err(); err != nil {
		log.Fatal(err)
	}
}
