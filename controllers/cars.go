package controllers

import (
	"net/http"

	"github.com/cyrip/monGO/driver"
	"github.com/gin-gonic/gin"
)

var backend driver.Backend

func Init(backend driver.Backend) {
	backend.Init()
}

func PostCar(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "OK",
	})
}

func GetCar(c *gin.Context) {
	backend.Search3("A.*")
	c.JSON(http.StatusOK, gin.H{
		"status": "OK",
	})
}
