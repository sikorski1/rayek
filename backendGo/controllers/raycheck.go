package controllers

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

type MapConfiguration struct {
	Title string `json:"title"`
	Coordinates [][][]float64 `json:"coordinates"`
	Center [2]float64 `json:"center"`
	Bounds [2][2]float64 `json:"bounds"`
}


func GetMapConfiguration(context *gin.Context) {
	cwd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	mapTitle := context.Param("mapTitle")
	data, err := os.ReadFile(cwd + "/data/mapData.json") 
	if err != nil {
		log.Print("Failed to read data file")
		context.JSON(http.StatusInternalServerError, gin.H{"error":"Failed to read data file"})
		return
	}
	var mapData map[string]MapConfiguration
	err = json.Unmarshal(data, &mapData)
	if err != nil {
		log.Print("Failed to parse JSON")
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse JSON"})
		return
	}
	if config, exists := mapData[mapTitle]; exists {
		context.JSON(http.StatusOK, config)
	} else {
		log.Print("Map configuration not found")
		context.JSON(http.StatusNotFound, gin.H{"error": "Map configuration not found"})
	}
}

func GetBuildings(context *gin.Context) {
	return
}

func ComputeRays(context *gin.Context) {
	return
}