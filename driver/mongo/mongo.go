package mongo

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"sync"

	gofakeit "github.com/brianvoe/gofakeit/v7"
	"github.com/cyrip/monGO/driver"
	"github.com/cyrip/monGO/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const MONGO_SHARDS int = 1

type MongoCars struct {
	connections [MONGO_SHARDS]*mongo.Client
	collections [MONGO_SHARDS]*mongo.Collection
	contexts    [MONGO_SHARDS]context.Context
}

func (this *MongoCars) Init() {
	for id := 0; id < MONGO_SHARDS; id++ {
		this.createConnection(id)
		// defer this.connections[id].Disconnect(this.contexts[id])
		this.createCollection(id, this.contexts[id])
		this.addIndex(id)
	}
	log.Println(this.contexts)
}

func (this *MongoCars) Dispose() {
	for id := 0; id < MONGO_SHARDS; id++ {
		defer this.connections[id].Disconnect(this.contexts[id])
	}
}

func (this *MongoCars) createConnection(id int) *mongo.Client {
	this.contexts[id] = context.Background()
	clientOptions := options.Client().ApplyURI(fmt.Sprintf("mongodb://localhost:%d", 27017+id))
	var err error
	this.connections[id], err = mongo.Connect(this.contexts[id], clientOptions)

	if err != nil {
		log.Fatal(err)
	}

	err = this.connections[id].Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Connected to mongodb on port %d", 27017+id)
	return this.connections[id]
}

func (this *MongoCars) createCollection(id int, ctx context.Context) *mongo.Collection {
	database := this.connections[id].Database("darth")
	this.collections[id] = database.Collection("Cars")
	return this.collections[id]
}

func (this *MongoCars) InsertOne(car driver.Car) {

	shard := this.getMongoShard(car.PlateNumber)
	car.UUID = utils.GetUUID(car.PlateNumber)
	response, err := this.collections[shard].InsertOne(this.contexts[shard], car)
	if err != nil {
		fmt.Println(response)
		fmt.Println(err)
	}
}

func (this *MongoCars) Seed(documentNumber int) {
	gofakeit.Seed(rand.Intn(10000))
	var inserted [MONGO_SHARDS]int
	for i := 0; i < documentNumber; i++ {
		car := this.getFakeCar()
		shard := this.getMongoShard(car.PlateNumber)
		response, err := this.collections[shard].InsertOne(this.contexts[shard], car)
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

func (this *MongoCars) getFakeCar() driver.Car {
	var fakeCar driver.Car
	err := gofakeit.Struct(&fakeCar)
	if err != nil {
		log.Fatal(err)
	}

	fakeCar.Data[0] = gofakeit.Color()
	fakeCar.Data[1] = gofakeit.City()
	fakeCar.Data[2] = gofakeit.BeerName()
	return fakeCar
}

func (this *MongoCars) addIndex(id int) {
	indexModel := mongo.IndexModel{
		Keys:    bson.D{{Key: "rendszam", Value: 1}},
		Options: options.Index().SetUnique(true),
	}

	// Create the index
	indexName, err := this.collections[id].Indexes().CreateOne(this.contexts[id], indexModel)
	if err != nil {
		log.Println("Index already exists")
		log.Println(err)
	} else {
		log.Printf("Index added %s", indexName)
	}
}

func (this *MongoCars) getMongoShard(plateNumber string) int {
	firstChar := plateNumber[0]
	return int(firstChar) % MONGO_SHARDS
}

func (this *MongoCars) findAsync(regex string) []driver.Car {
	var wg sync.WaitGroup
	var mutex sync.Mutex

	var response []driver.Car
	response = make([]driver.Car, 0)

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

			cursor, err := this.collections[id].Find(this.contexts[id], filter)
			if err != nil {
				log.Fatal(err)
			}
			defer cursor.Close(this.contexts[id])

			mutex.Lock()
			found := 0

			var car driver.Car

			for cursor.Next(this.contexts[id]) {
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
	return response
	//log.Println(response)
	//log.Println(len(response))
	//for i, e := range response {
	//log.Printf("%d %s", i, e.PlateNumber)
	//}
}

func (this *MongoCars) Search3(regex string) []driver.Car {
	return this.findSync(regex, 0)
	// return this.findAsync("AA.*")
}

func (this *MongoCars) findSync(regex string, id int) []driver.Car {
	var response []driver.Car
	response = make([]driver.Car, 0)

	filter := bson.M{
		"$or": []interface{}{
			bson.M{"rendszam": bson.M{"$regex": regex, "$options": "i"}},
			bson.M{"tulajdonos": bson.M{"$regex": regex, "$options": "i"}},
			bson.M{"adatok": bson.M{"$regex": regex, "$options": "i"}},
		},
	}

	cursor, err := this.collections[id].Find(this.contexts[id], filter)
	if err != nil {
		log.Fatal(err)
	}
	defer cursor.Close(this.contexts[0])

	found := 0

	for cursor.Next(this.contexts[id]) {
		var result bson.M
		var car driver.Car
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

	log.Println(response)
	return response
}
