package main

package main

import (
	"backendGo/config"
	"backendGo/routes"
	"github.com/gin-gonic/gin"
)

func main() {
	config.LoadEnv()
	config.ConnectDB()

	r := gin.Default()

	
	routes.SetupRoutes(r)


	r.Run(":8080")
}

