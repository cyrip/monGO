package application

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"sync"
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
	UUID5       string   `bson:"uuid,omitempty" fake:"{uuid}"`
	PlateNumber string   `bson:"rendszam,omitempty" fake:"{regex:[A-Z]{3}}-{regex:[0-9]{3}}"`
	ValidUntil  string   `bson:"forgalmi_ervenyes,omitempty" fake:"{year}-{month}-{day}" format:"2006-01-02"`
	Owner       string   `bson:"tulajdonos,omitempty" fake:"{name}"`
	Data        []string `bson:"adatok,omitempty" fakesize:"3"`
}

type Cars struct {
	Id          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UUID5       string             `bson:"uuid,omitempty" json:"uuid"`
	PlateNumber string             `bson:"rendszam,omitempty" json:"rendszam"`
	Owner       string             `bson:"tulajdonos,omitempty" json:"tulajdonos"`
	ValidUntil  string             `bson:"forgalmi_ervenyes,omitempty" json:"forgalmi_ervenyes"`
	Data        []string           `bson:"adatok,omitempty" json:"adatok"`
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

	for id := 0; id < MONGO_SHARDS; id++ {
		createConnection(id)
		defer connections[id].Disconnect(contexts[id])
		createCollection(id, contexts[id])
		addIndex(id)
	}

	//insertFakeCars()
	find0()
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
	for i := 0; i < MONGO_SHARDS; i++ {
		log.Printf("Inserted docs shard%d %d", i, inserted[i])
	}
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
		Keys:    bson.D{{Key: "rendszam", Value: 1}},
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
	return int(firstChar) % MONGO_SHARDS
}

func find3(id int, regex string, wg *sync.WaitGroup) {

	defer wg.Done()

	filter := bson.M{
		"$or": []interface{}{
			bson.M{"rendszam": bson.M{"$regex": regex, "$options": "i"}},
			bson.M{"tulajdonos": bson.M{"$regex": regex, "$options": "i"}},
			bson.M{"adatok": bson.M{"$regex": regex, "$options": "i"}},
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

func findAsync(regex string) {
	var wg sync.WaitGroup
	var mutex sync.Mutex

	var response []Cars
	response = make([]Cars, 0)

	for i := 0; i < MONGO_SHARDS; i++ {
		wg.Add(1)
		go func(id int, regex string, wg *sync.WaitGroup) {
			defer wg.Done()

			filter := bson.M{
				"$or": []interface{}{
					bson.M{"rendszam": bson.M{"$regex": regex, "$options": "i"}},
					bson.M{"tulajdonos": bson.M{"$regex": regex, "$options": "i"}},
					bson.M{"adatok": bson.M{"$regex": regex, "$options": "i"}},
				},
			}

			cursor, err := collections[id].Find(contexts[id], filter)
			if err != nil {
				log.Fatal(err)
			}
			defer cursor.Close(contexts[id])

			mutex.Lock()
			found := 0

			var car Cars

			for cursor.Next(contexts[id]) {
				var result bson.M
				if err := cursor.Decode(&result); err != nil {
					log.Fatal(err)
				}
				// jsonData, _ := json.Marshal(result)
				// response = append(response, string(jsonData))
				bsonBytes, _ := bson.Marshal(result)
				bson.Unmarshal(bsonBytes, &car)
				response = append(response, car)
				// fmt.Println(string(jsonData))
				found = found + 1
			}

			if err := cursor.Err(); err != nil {
				log.Fatal(err)
			}

			log.Printf("Found in shard %d %d", id, found)
			mutex.Unlock()
		}(i, regex, &wg)
	}

	wg.Wait()
	//log.Println(response)
	log.Println(len(response))
	for i, e := range response {
		log.Printf("%d %s", i, e.PlateNumber)
	}
}

func find0() {
	findAsync("FW.*")
}
