package application

import (
	"log"

	"github.com/cyrip/monGO/config"
)

func StartServer() {
	log.Printf("monGO starting webserver %d", config.MONGO_PORT)
	initHTTPServer(config.MONGO_PORT)
}
