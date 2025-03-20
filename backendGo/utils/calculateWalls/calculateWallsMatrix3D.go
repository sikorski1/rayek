package calculateWalls

import (
	. "backendGo/types"
	"backendGo/utils/calculateWalls/tools"
	"encoding/json"
	"log"
	"os"
)

type MapConfig struct {
	LatMin, LatMax, LonMin, LonMax float64
	Size, HeightMaxLevels int
}

func CalculateWallsMatrix3D(filePath string, mapConfig MapConfig) [][][]float64 {
	buildingsFilePath := tools.CalculateWalls(filePath)
	data, err := os.ReadFile(buildingsFilePath)
	if err != nil {
		log.Fatalf("Failed to read data file: %v", err)
	}
	var buildings []Building
	err = json.Unmarshal(data, &buildings)
	if err != nil {
		log.Fatalf("Error parsing JSON: %v", err)
	}
	matrix := tools.GenerateBuildingMatrix(buildings, mapConfig.LatMin, mapConfig.LatMax, mapConfig.LonMin, mapConfig.LonMax, mapConfig.Size, mapConfig.HeightMaxLevels)
	return matrix
}