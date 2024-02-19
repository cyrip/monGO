package application

import (
	"github.com/cyrip/monGO/controllers"
	"github.com/gin-gonic/gin"
)

var (
	router = gin.Default()
)

func initHTTPServer(port string) {
	bindUrls()
	router.Run(":" + port)
}

func bindUrls() {
	v1 := router.Group("/v1")
	{
		v1.GET("/healthcheck", controllers.HealthCheck)
	}

	router.POST("/jarmuvek", controllers.PostCar)
	router.GET("/jarmuvek", controllers.GetCar)
}
