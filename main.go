package main

import (
	"fmt"
	"log"
	"os"
	"runtime"
	"time"

	"github.com/cyrip/monGO/application"
	"github.com/cyrip/monGO/config"
	"github.com/cyrip/monGO/controllers"
	"github.com/cyrip/monGO/driver"
	"github.com/cyrip/monGO/driver/elastic"
	"github.com/cyrip/monGO/driver/mongo"
	"github.com/joho/godotenv"
	"github.com/pbnjay/memory"
)

func init() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Println(err)
	}
}

// https://www.mongodb.com/blog/post/mongodb-go-driver-tutorial
func main() {
	// runtime.GOMAXPROCS(1)

	log.Printf("monGO started! CPUs: %d, total/free memory %dMB/%dMB", runtime.NumCPU(), memory.TotalMemory()/1024, memory.FreeMemory()/1024)

	config.InitEnv()

	mode := config.SERVER_MODE
	backend := config.SERVER_BACKEND

	log.Printf("Mode: %s, backend: %s", config.SERVER_MODE, config.SERVER_BACKEND)

	switch mode {
	case "server":
		switch backend {
		case "mongo":
			backend := mongo.MongoCars{}
			controllers.Init(&backend)
			application.StartServer()
		case "elastic":
			// it is an ugly hack, but on my slow machine elastic docker start very slowly, need to tune conenction timeout
			if config.ELASTIC_MODE == "single" {
				time.Sleep(time.Second * 60)
			} else {
				time.Sleep(time.Second * 60)
			}
			backend := elastic.Elastic{}
			controllers.Init(&backend)
			application.StartServer()
		case "sql":
			log.Fatalln("Not implemented yet")
		default:
			log.Fatalf("There is no such backend!")
			os.Exit(1)
		}

	case "migrate":
		log.Printf("Start migration")
		//container := container.Container{}
		//mongoCars := container.GetMongo()
		mongoCars := mongo.MongoCars{}
		mongoCars.Init()

		car := driver.Car{}
		car.Owner = "Owner1"
		car.PlateNumber = "IOP-920"
		mongoCars.InsertOne(car)
		fmt.Println(mongoCars.Search3(".*IOP-91.*"))
		fmt.Println(mongoCars.CountDocuments())
		//mongoCars.InsertFakeCars()
		//mongoCars.Find0()
		//mongoCars.Disconnect()
		//ela := elastic.Elastic{}
		//ela.Init("cars")
		//ela.Search3("A.*")

	default:
		log.Fatalf("There is no such mode!")
		os.Exit(1)
	}
}
