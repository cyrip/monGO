package container

import (
	"github.com/cyrip/monGO/driver/mongo"
)

type Container struct {
	MongoCars *mongo.MongoCars
}

func (this *Container) GetMongo() *mongo.MongoCars {
	if this.MongoCars == nil {
		mongo := mongo.MongoCars{}
		this.MongoCars = &mongo
	}
	return this.MongoCars
}
