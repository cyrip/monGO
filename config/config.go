package config

import (
	"github.com/cyrip/monGO/utils"
)

const (
	MONGO_PORT_DEFAULT = "8080"
)

var (
	MONGO_PORT string
)

func InitEnv() {
	MONGO_PORT = utils.GetEnvFallback("MONGO_PORT", MONGO_PORT_DEFAULT)
}
