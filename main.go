package main

import (
	"log"
	"os"
	"runtime"

	"github.com/cyrip/monGO/application"
	"github.com/cyrip/monGO/config"
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
	log.Printf("monGO started! CPUs: %d, total/free memory %dMB/%dMB", runtime.NumCPU(), memory.TotalMemory()/1024, memory.FreeMemory()/1024)
	mode := getopt.StringLong("mode", 'm', "", "[-m|--mode] [server|migrate]")
	getopt.Parse()

	config.InitEnv()

	switch *mode {
	case "server":
		application.StartServer()
	case "migrate":
		log.Printf("Start migration")
		application.MongoTest()
	default:
		log.Fatalf("There is no such mode!")
		os.Exit(1)
	}
}
