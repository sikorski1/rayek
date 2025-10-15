package routes

import (
	"backendGo/controllers"

	"github.com/gin-gonic/gin"
)

func SetupRayCheckRoutes(router *gin.Engine) {
	raycheckRouter := router.Group("/api/maps")
	{
		raycheckRouter.GET("/", controllers.GetMaps)
		raycheckRouter.GET("/:mapTitle", controllers.GetMapById)
		raycheckRouter.POST("/rayLaunch/:mapTitle", controllers.Create3DRayLaunching)
	}
}
