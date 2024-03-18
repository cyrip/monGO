package main

import (
	"fmt"
	"log"
	"os"
	"runtime"

	"github.com/cyrip/monGO/application"
	"github.com/cyrip/monGO/config"
	"github.com/cyrip/monGO/controllers"
	"github.com/cyrip/monGO/driver"
	"github.com/cyrip/monGO/driver/elastic"
	"github.com/cyrip/monGO/driver/mongo"
	"github.com/cyrip/monGO/utils"
	"github.com/joho/godotenv"
	"github.com/pbnjay/memory"
	"github.com/pborman/getopt/v2"
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
	mode := getopt.StringLong("mode", 'm', "", "[-m|--mode] [server|migrate]")
	backend := getopt.StringLong("backend", 'b', "", "[-b|--backend] [mongo|elastic|sql]")

	getopt.Parse()

	config.InitEnv()

	log.Println(utils.IsDateValue("1999-12-25"))
	log.Println(utils.IsDateValue("2024-03-15"))

	switch *mode {
	case "server":
		switch *backend {
		case "mongo":
			backend := mongo.MongoCars{}
			controllers.Init(&backend)
			application.StartServer()
		case "elastic":
			backend := elastic.Elastic{}
			backend.Init()
			backend.CreateIndex()
			//backend.DeleteIndex()
			//os.Exit(1)
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
