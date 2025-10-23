package controllers

import (
	. "backendGo/types"
	"backendGo/utils/calculations"
	"backendGo/utils/raylaunching"
	"encoding/json"
	"fmt"

	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
)

type Map struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Img         string `json:"img"`
	Size        string `json:"size"`
}
type MapAndBuildingsResponse struct {
	MapData       MapConfiguration       `json:"mapData"`
	BuildingsData BuildingsConfiguration `json:"buildingsData"`
}

type MapConfiguration struct {
	Title       string        `json:"title"`
	Coordinates [][][]float64 `json:"coordinates"`
	Center      [2]float64    `json:"center"`
	Bounds      [2][2]float64 `json:"bounds"`
	Size        int           `json:"size"`
}

type Features struct {
	Type       string      `json:"type"`
	Properties interface{} `json:"properties"`
	Geometry   interface{} `json:"geometry"`
}
type BuildingsConfiguration struct {
	Type     string     `json:"type"`
	Features []Features `json:"features"`
}

func GetMaps(context *gin.Context) {
	cwd, err := os.Getwd()
	if err != nil {
		log.Println(err)
	}
	dataPath := filepath.Join(cwd, "data", "maps.json")
	fmt.Println(dataPath)
	data, err := os.ReadFile(dataPath)
	if err != nil {
		log.Println("Failed to read data file")
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read data file"})
		return
	}
	var mapsData []Map
	err = json.Unmarshal(data, &mapsData)
	if err != nil {
		log.Println("Failed to parse JSON")
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse JSON"})
		return
	}
	print(mapsData)
	context.JSON(http.StatusOK, mapsData)
}

func GetMapById(context *gin.Context) {
	cwd, err := os.Getwd()
	if err != nil {
		log.Println(err)
	}
	mapTitle := context.Param("mapTitle")
	fmt.Println(filepath.Join(cwd, "data", mapTitle, "mapData.json"))
	data, err := os.ReadFile(filepath.Join(cwd, "data", mapTitle, "mapData.json"))
	if err != nil {
		log.Println("Failed to read data file")
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read data file"})
		return
	}
	var mapData MapConfiguration
	err = json.Unmarshal(data, &mapData)
	if err != nil {
		log.Println("Failed to parse JSON")
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse JSON"})
		return
	}
	if mapData.Title == "" || len(mapData.Coordinates) == 0 {
		log.Println("Map configuration not found")
		context.JSON(http.StatusNotFound, gin.H{"error": "Map configuration not found"})
		return
	}

	data, err = os.ReadFile(filepath.Join(cwd, "data", mapTitle, "rawBuildings.json"))
	if err != nil {
		log.Println("Failed to read data file")
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read data file"})
	}
	var buildingsData BuildingsConfiguration
	err = json.Unmarshal(data, &buildingsData)

	if err != nil {
		log.Println("Failed to parse JSON")
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse JSON"})
		return
	}
	if len(buildingsData.Features) == 0 {
		log.Println("Buildings configuration not found")
		context.JSON(http.StatusNotFound, gin.H{"error": "Buildings configuration not found"})
		return
	}
	response := MapAndBuildingsResponse{
		MapData:       mapData,
		BuildingsData: buildingsData,
	}
	context.JSON(http.StatusOK, response)
}

type RayLaunchRequest struct {
	NumberOfRaysAzimuth   int         `json:"numberOfRaysAzimuth" binding:"required,min=1,max=2879"`
	NumberOfRaysElevation int         `json:"numberOfRaysElevation" binding:"required,min=1,max=2879"`
	NumberOfInteractions  int         `json:"numberOfInteractions" binding:"required,min=1,max=10"`
	ReflectionFactor      float64     `json:"reflectionFactor" binding:"required,gte=0,lte=1"`
	StationPower          float64     `json:"stationPower" binding:"required,gte=0.01,lte=100"`
	MinimalRayPower       float64     `json:"minimalRayPower" binding:"required,gte=-160,lte=-60"`
	Frequency             float64     `json:"frequency" binding:"required,gte=0.1,lte=100"`
	Size                  int         `json:"size" binding:"required,oneof=250 400 500"`
	StationPos            Point3D     `json:"stationPos" binding:"required"`
	SingleRays            []SingleRay `json:"singleRays" binding:"omitempty,dive,required"`
	DiffractionRayNumber  int         `json:"diffractionRayNumber" binding:"omitempty,gte=0,lte=120"`
}

const MaxTotalRays = 2880

func Create3DRayLaunching(context *gin.Context) {
	mapTitle := context.Param("mapTitle")

	var request RayLaunchRequest
	if err := context.ShouldBindJSON(&request); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	log.Printf("Received request: %+v\n", request)
	cwd, err := os.Getwd()
	if err != nil {
		log.Println(err)
	}

	totalRays := request.NumberOfRaysAzimuth + request.NumberOfRaysElevation
	if totalRays > MaxTotalRays {
		context.JSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintf("Total number of rays (%d) exceeds maximum limit (%d). Azimuth: %d, Elevation: %d",
				totalRays, MaxTotalRays, request.NumberOfRaysAzimuth, request.NumberOfRaysElevation),
		})
		return
	}
	
	var matrixInt [][][]int16
	err = calculations.LoadMatrixBinary(filepath.Join(cwd, "data", mapTitle, "wallsMatrix3D_floor.bin"), &matrixInt)
	if err != nil {
		log.Println("Failed to load matrix:", err)
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load matrix"})
		return
	}

	matrix := calculations.ConvertInt16MatrixToFloat64(matrixInt)
	var wallNormals []Normal3D
	err = calculations.LoadMatrixBinary(filepath.Join(cwd, "data", mapTitle, "wallNormals3D.bin"), &wallNormals)
	if err != nil {
		log.Println("Failed to load matrix:", err)
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load matrix"})
		return
	}
	var diffRayNum int
	if request.DiffractionRayNumber == 0 {
		diffRayNum = 1
	} else {
		diffRayNum = request.DiffractionRayNumber
	}
	config := raylaunching.RayLaunching3DConfig{
		NumOfRaysAzim:         request.NumberOfRaysAzimuth,
		NumOfRaysElev:         request.NumberOfRaysElevation,
		NumOfInteractions:     request.NumberOfInteractions,
		WallMapNumber:         1000,
		RoofMapNumber:         5000,
		CornerMapNumber:       10000,
		RoofCornerMapNumber:   10001,
		BuldingInteriorNumber: 20000,
		SizeX:                 float64(request.Size - 1),
		SizeY:                 float64(request.Size - 1),
		SizeZ:                 30 - 1,
		Step:                  1.0,
		ReflFactor:            request.ReflectionFactor,
		TransmitterPower:      request.StationPower,
		MinimalRayPower:       request.MinimalRayPower,
		TransmitterFreq:       request.Frequency * 1e9,
		WaveLength:            0,
		TransmitterPos:        Point3D{X: request.StationPos.X, Y: request.StationPos.Y, Z: request.StationPos.Z},
		SingleRays:            request.SingleRays,
		DiffractionRayNumber:  diffRayNum,
	}
	config.WaveLength = 299792458 / (config.TransmitterFreq)
	start := time.Now()
	rayLaunching := raylaunching.NewRayLaunching3D(matrix, wallNormals, config)
	rayLaunching.CalculateRayLaunching3D()
	stop := time.Since(start)
	fmt.Printf("RayLaunching 3D calculation time: %v\n", stop)

	context.JSON(http.StatusOK, gin.H{
		"message":        "Request received successfully",
		"mapTitle":       mapTitle,
		"stationPos":     request.StationPos,
		"powerMap":       rayLaunching.PowerMap,
		"rayPaths":       rayLaunching.RayPaths,
		"powerMapLegend": rayLaunching.PowerMapLegend,
	})
}
