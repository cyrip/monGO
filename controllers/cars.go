package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func PostCar(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "OK",
	})
}

func GetCar(c *gin.Context) {

	c.JSON(http.StatusOK, gin.H{
		"status": "OK",
	})
}
