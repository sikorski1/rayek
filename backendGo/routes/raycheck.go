package routes

import (
	"backendGo/controllers"

	"github.com/gin-gonic/gin"
)

func SetupRayCheckRoutes (router *gin.Engine) {
	raycheckRouter := router.Group("/api/maps")
	{
		raycheckRouter.GET("/", controllers.GetMaps)
		raycheckRouter.GET("/:mapTitle", controllers.GetMapById)
		raycheckRouter.GET("/wallmatrix/:mapTitle", controllers.GetWallMatrix)
		raycheckRouter.POST("/rayLaunch/:mapTitle", controllers.Create3DRayLaunching)
		raycheckRouter.GET("/buildings/:mapTitle", controllers.GetBuildings)
		raycheckRouter.POST("/compute", controllers.ComputeRays)
	}
}