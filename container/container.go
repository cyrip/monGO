package container

import (
	"github.com/cyrip/monGO/driver"
	"github.com/cyrip/monGO/driver/mongo"
)

type Container struct {
	Backend *driver.Backend
}

func (this *Container) GetMongo() *mongo.MongoCars {
	if this.MongoCars == nil {
		mongo := mongo.MongoCars{}
		this.MongoCars = &mongo
	}
	return this.MongoCars
}
