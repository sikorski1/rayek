package main

import (
	. "backendGo/types"
	"backendGo/utils/calculateWalls/tools"
	"backendGo/utils/calculations"
	"encoding/json"
	"image/png"
	"log"
	"os"
)
func main() {
	tools.CalculateWalls("AGHFragment_buildings_raw.json")
	data, err := os.ReadFile("AGHFragment_buildings.json")
	if err != nil {
		log.Fatalf("Failed to read data file: %v", err)
	}

	// Parse the GeoJSON
	var buildings []Building
	err = json.Unmarshal(data, &buildings)
	if err != nil {
		log.Fatalf("Error parsing JSON: %v", err)
	}
	matrix := tools.GenerateBuildingMatrix(buildings, 50.065311, 50.067556, 19.914029, 19.917527, 250, 25)
	heatmap := calculations.GenerateHeatmap(matrix[4])
	 f, _ := os.Create("heatmap.png")
	 defer f.Close()
	 png.Encode(f, heatmap)
}