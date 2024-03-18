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

	formValues := c.PostFormMap("")
	log.Println(formValues)

	if err := c.ShouldBindJSON(&postData); err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	log.Println(postData)
	if !(utils.IsDateValue(postData.ValidUntil)) {
		log.Println("forgalmi_ervenyes must be date")
		c.JSON(http.StatusBadRequest, gin.H{"error": "forgalmi_ervenyes must be date"})
		return
	}

	inserted := backend.InsertOne(postData)

	if inserted == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "already exists"})
		return
	}

	c.Header("Location", "/jarmuvek/"+inserted.UUID)

	c.JSON(http.StatusCreated, inserted)
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
	if !exists || query == "" {
		c.String(http.StatusBadRequest, "")
	}

	response := backend.Search3(query)
	//	log.Println(response)
	c.JSON(http.StatusOK, response)
}

func CountCars(c *gin.Context) {
	count := backend.CountDocuments()
	c.String(http.StatusOK, strconv.FormatInt(count, 10))
}
