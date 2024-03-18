package application

import (
	"log"

	"github.com/cyrip/monGO/config"
)

func StartServer() {
	log.Printf("monGO starting webserver %d", config.SERVER_PORT)
	initHTTPServer(config.SERVER_PORT)
}
