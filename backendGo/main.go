package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"backendGo/routes"
)

func main() {
	router := gin.Default()

	config := cors.Config{
		AllowOrigins:    []string{"http://localhost:5173"},
		AllowMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders: []string{"Origin", "Content-Length", "Content-Type", "Authorization"},
		AllowCredentials: true,
	}
	router.Use(cors.New(config))
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	routes.SetupRayCheckRoutes(router)

	router.Run(":3000")
}
