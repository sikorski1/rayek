package main

import (
	. "backendGo/types"
	"backendGo/utils/calculations"
	"encoding/json"
	"fmt"
	"image/png"
	"log"
	"os"
	"path/filepath"
	"strings"
)




func main() {
	if len(os.Args) < 2 {
		fmt.Println("Use: go run main.go ../../data/<MAP_NAME>")
		os.Exit(1)
	}
	mapFolderPath := os.Args[1]
	configFilePath := filepath.Join(mapFolderPath, "mapConfig.json")
	if _, err := os.Stat(configFilePath); os.IsNotExist(err) {
		log.Fatalf("Config file does not exist: %s", configFilePath)
	}
	//config load
	file, err := os.ReadFile(configFilePath)
	if err != nil {
		log.Fatalf("Error reading file %s: %v",configFilePath, err)
	}
	var config MapConfig
	err = json.Unmarshal(file, &config)
	if err != nil {
		log.Fatalf("Error parsing JSON %v", err)
	}
	//calculating walls matrix
	calculations.CalculateWallsMatrix3D(mapFolderPath, config)
	wallsMatrixPath := filepath.Join(mapFolderPath, "wallsMatrix3D.bin")
	wallNormalsPath := filepath.Join(mapFolderPath, "wallNormals3D.bin")
	//check if matrix calculated properly
	var matrix [][][]float64
	var wallNormals []Normal3D
	err = calculations.LoadMatrixBinary(wallsMatrixPath, &matrix)
	if err != nil {
		log.Fatalf("Error loading matrix %v", err)
	}
	err = calculations.LoadMatrixBinary(wallNormalsPath, &wallNormals)
	if err != nil {
		log.Fatalf("Error loading matrix %v", err)
	}
	mapName := filepath.Base(mapFolderPath)	
	HEIGHT := 0
	outputFileName := fmt.Sprintf("%s-%dm.png", strings.ReplaceAll(mapName, " ", "_"), HEIGHT)
	outputFilePath := filepath.Join(mapFolderPath, outputFileName)
	heatmap := calculations.GenerateHeatmap(matrix[HEIGHT])
    f, _ := os.Create(outputFilePath)
    defer f.Close()
    png.Encode(f, heatmap)
	fmt.Printf("Heatmap saved to %s\n", outputFilePath)
}