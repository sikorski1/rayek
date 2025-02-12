package routes

import (
	"github.com/gin-gonic/gin"
	"backendGo/controllers"
)

func SetupRayCheckRoutes (router *gin.Engine) {
	raycheckRouter := router.Group("/raycheck")
	{
		raycheckRouter.GET("/:mapTitle", controllers.GetMapConfiguration)
		raycheckRouter.GET("/buildings/:mapTitle", controllers.GetBuildings)
		raycheckRouter.POST("/compute", controllers.ComputeRays)
	}
}