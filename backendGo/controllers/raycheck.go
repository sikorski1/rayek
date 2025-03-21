package controllers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

type MapConfiguration struct {
	Title string `json:"title"`
	Coordinates [][][]float64 `json:"coordinates"`
	Center [2]float64 `json:"center"`
	Bounds [2][2]float64 `json:"bounds"`
}

type Features struct {
	Type string `json:"type"`
	Properties interface{} `json:"properties"`
	Geometry interface{} `json:"geometry"`
	Id string `json:"id"`
}
type BuildingsConfiguration struct {
	Type string `json:"type"`
	Features []Features `json:"features"`
}

func GetMapConfiguration(context *gin.Context) {
	cwd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	mapTitle := context.Param("mapTitle")
	fmt.Println(filepath.Join(cwd, "data", mapTitle,"mapData.json"))
	data, err := os.ReadFile(filepath.Join(cwd, "data", mapTitle,"mapData.json")) 
	if err != nil {
		log.Print("Failed to read data file")
		context.JSON(http.StatusInternalServerError, gin.H{"error":"Failed to read data file"})
		return
	}
	var mapData MapConfiguration
	err = json.Unmarshal(data, &mapData)
	if err != nil {
		log.Print("Failed to parse JSON")
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse JSON"})
		return
	}
	if mapData.Title == "" || len(mapData.Coordinates) == 0 {
		log.Print("Map configuration not found")
		context.JSON(http.StatusNotFound, gin.H{"error": "Map configuration not found"})
		return
	}

	context.JSON(http.StatusOK, mapData)
}

func GetBuildings(context *gin.Context) {
	cwd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	mapTitle := context.Param("mapTitle")
	data, err := os.ReadFile(filepath.Join(cwd, "data", mapTitle,"rawBuildings.json"))
	if err != nil {
		log.Print("Failed to read data file")
		context.JSON(http.StatusInternalServerError, gin.H{"error":"Failed to read data file"})
	}
	var buildingData BuildingsConfiguration
	err = json.Unmarshal(data, &buildingData)

	if err != nil {
		log.Print("Failed to parse JSON")
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse JSON"})
		return
	}
	if len(buildingData.Features) == 0 {
		log.Print("Buildings configuration not found")
		context.JSON(http.StatusNotFound, gin.H{"error": "Buildings configuration not found"})
		return
	}
	context.JSON(http.StatusOK, buildingData)
}

func ComputeRays(context *gin.Context) {
	southwest :=  [2]float64{19.914029, 50.065311}
	southeast := [2]float64{19.917527, 50.065311}
	northeast := [2]float64{19.917527, 50.067556}
	numPoints := 500

	xStep := (southeast[0] - southwest[0]) / (float64(numPoints - 1));
	yStep := (northeast[1] - southeast[1]) / (float64(numPoints - 1));
	points := make([][2]float64, 0, numPoints*numPoints)
	for i := 0; i < numPoints; i++ {
		for j := 0; j < numPoints; j++ {
			x := southwest[0] + float64(i)*xStep
			y := southwest[1] + float64(j)*yStep
			points = append(points, [2]float64{x, y})
		}
	}
	context.JSON(http.StatusOK, points)

}

type RayLaunchRequest struct {
	MapNumber      int     `json:"mapNumber" binding:"required"`
	StationHeight  float64 `json:"stationHeight" binding:"required"`
	Frequency      float64 `json:"freq" binding:"required"`
}

func Create3DRayLaunching(context *gin.Context) {
	var request RayLaunchRequest

	if err := context.ShouldBindJSON(&request); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	matrixPath := fmt.Sprintf("../data/%v/wallsMatrix3D.bin", request.MapNumber)
	fmt.Printf(matrixPath)
	println("Launching ray with params:")
	println("Map Number:", request.MapNumber)
	println("Station Height:", request.StationHeight)
	println("Frequency:", request.Frequency)

	context.JSON(http.StatusOK, gin.H{
		"message":       "Request received successfully",
		"mapNumber":     request.MapNumber,
		"stationHeight": request.StationHeight,
		"frequency":     request.Frequency,
	})
}