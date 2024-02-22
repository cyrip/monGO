package controllers

import (
	"net/http"

	"github.com/cyrip/monGO/driver/mongo"
	"github.com/gin-gonic/gin"
)

var mongoCars mongo.MongoCars

func Init() {
	mongoCars = mongo.MongoCars{}
	mongoCars.Init()
}

func PostCar(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "OK",
	})
}

func GetCar(c *gin.Context) {
	mongoCars.Find0()
	c.JSON(http.StatusOK, gin.H{
		"status": "OK",
	})
}
