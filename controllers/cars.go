package controllers

import (
	"log"
	"net/http"

	"strconv"

	"github.com/cyrip/monGO/driver"
	"github.com/cyrip/monGO/utils"
	"github.com/gin-gonic/gin"
)

var backend driver.Backend

func Init(bck driver.Backend) {
	backend = bck
	backend.Init()
}

func PostCar(c *gin.Context) {
	var postData driver.Car

	log.Println("New post request")

	if err := c.ShouldBindJSON(&postData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if !(utils.IsDateValue(postData.ValidUntil)) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "forgalmi_ervenyes must be date"})
		return
	}

	log.Println(postData)
	inserted := backend.InsertOne(postData)

	if inserted == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "already exists"})
		return
	}

	c.Header("Location", "/jarmuvek/"+inserted.UUID)

	c.JSON(http.StatusCreated, gin.H{
		"uuid":              inserted.UUID,
		"rendszam":          postData.PlateNumber,
		"tulajdonos":        postData.Owner,
		"forgalmi_ervenyes": postData.ValidUntil,
		"adatok":            postData.Data,
	})
}

func GetCar(c *gin.Context) {
	// fmt.Println(backend.Search3("IOP.*"))
	uuid := c.Param("uuid")
	car := backend.GetByUUID(uuid)
	if car == nil {
		log.Print("not found by UUID")
		c.JSON(http.StatusNotFound, gin.H{"error": "not found " + uuid})
		return
	}

	c.JSON(http.StatusOK, car)
}

func Search(c *gin.Context) {
	query, exists := c.GetQuery("q")
	if !exists {
		c.String(http.StatusNotFound, "")
		return
	}

	response := backend.Search3(query)
	//	log.Println(response)
	c.JSON(http.StatusOK, response)
}

func CountCars(c *gin.Context) {
	count := backend.CountDocuments()
	c.String(http.StatusOK, strconv.FormatInt(count, 10))
}
