package config

import (
	"github.com/cyrip/monGO/utils"
)

const (
	SERVER_PORT_DEFAULT    = "8080"
	SERVER_MODE_DEFAULT    = "server"
	SERVER_BACKEND_DEFAULT = "mongo"
	ELASTIC_MODE_DEFAULT   = "single"
)

var (
	SERVER_PORT    string
	SERVER_MODE    string
	SERVER_BACKEND string
	ELASTIC_MODE   = "single"
)

func InitEnv() {
	SERVER_PORT = utils.GetEnvFallback("SERVER_PORT", SERVER_PORT_DEFAULT)
	SERVER_MODE = utils.GetEnvFallback("SERVER_MODE", SERVER_MODE_DEFAULT)
	SERVER_BACKEND = utils.GetEnvFallback("SERVER_BACKEND", SERVER_BACKEND_DEFAULT)
	ELASTIC_MODE = utils.GetEnvFallback("ELASTIC_MODE", ELASTIC_MODE_DEFAULT)
}
