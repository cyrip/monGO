package controllers

import (
	"fmt"
	"log"
	"net/http"

	"github.com/cyrip/monGO/driver"
	"github.com/gin-gonic/gin"
)

var backend driver.Backend

func Init(bck driver.Backend) {
	backend = bck
	backend.Init()
}

func PostCar(c *gin.Context) {
	var postData driver.Car

	if err := c.ShouldBindJSON(&postData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	log.Println(postData)
	c.JSON(http.StatusOK, gin.H{
		"uuid":              postData.UUID,
		"rendszam":          postData.PlateNumber,
		"tulajdonos":        postData.Owner,
		"forgalmi_ervenyes": postData.ValidUntil,
		"adatok":            postData.Data,
	})
}

func GetCar(c *gin.Context) {
	fmt.Println(backend.Search3("A.*"))

	uuid := c.Param("uuid")
	c.JSON(http.StatusOK, gin.H{
		"uuid": uuid,
	})
}
