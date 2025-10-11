package main

import (
	. "backendGo/types"
	"backendGo/utils/calculations"
	"encoding/json"
	"flag"
	"fmt"
	"image/png"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	genHeatmaps := flag.Bool("gen-heatmaps", false, "Generate heatmaps to /imgs/raw folder")
	flag.Parse()

	if flag.NArg() < 1 {
		fmt.Println("Usage: go run main.go [-gen-heatmaps] ../../data/<MAP_NAME>")
		os.Exit(1)
	}
	mapFolderPath := flag.Arg(0)

	configFilePath := filepath.Join(mapFolderPath, "mapConfig.json")
	if _, err := os.Stat(configFilePath); os.IsNotExist(err) {
		log.Fatalf("Config file does not exist: %s", configFilePath)
	}

	file, err := os.ReadFile(configFilePath)
	if err != nil {
		log.Fatalf("Error reading file %s: %v", configFilePath, err)
	}
	var config MapConfig
	if err := json.Unmarshal(file, &config); err != nil {
		log.Fatalf("Error parsing JSON %v", err)
	}

	calculations.CalculateWallsMatrix3D(mapFolderPath, config)

	wallsMatrixPath := filepath.Join(mapFolderPath, "wallsMatrix3D_floor.bin")
	wallNormalsPath := filepath.Join(mapFolderPath, "wallNormals3D.bin")

	var matrixInt [][][]int16
	var wallNormals []Normal3D
	if err := calculations.LoadMatrixBinary(wallsMatrixPath, &matrixInt); err != nil {
		log.Fatalf("Error loading matrix %v", err)
	}

	matrix := calculations.ConvertInt16MatrixToFloat64(matrixInt)

	if err := calculations.LoadMatrixBinary(wallNormalsPath, &wallNormals); err != nil {
		log.Fatalf("Error loading normals %v", err)
	}

	fmt.Println(len(matrix[0]), len(matrix[0][0]), len(matrix))
	mapName := filepath.Base(mapFolderPath)

	if *genHeatmaps {
		imgFolder := filepath.Join(mapFolderPath, "imgs", "raw")
		os.MkdirAll(imgFolder, os.ModePerm)
		for height := 0; height < 30; height++ {
			outputFileName := fmt.Sprintf("%s-%dm.png", strings.ReplaceAll(mapName, " ", "_"), height)
			outputFilePath := filepath.Join(imgFolder, outputFileName)
			heatmap := calculations.GenerateHeatmap(matrix[height])
			f, err := os.Create(outputFilePath)
			if err != nil {
				fmt.Printf("Failed to create file for height %d: %v\n", height, err)
				continue
			}
			if err := png.Encode(f, heatmap); err != nil {
				fmt.Printf("Failed to encode PNG for height %d: %v\n", height, err)
				f.Close()
				continue
			}
			f.Close()
			fmt.Printf("Heatmap saved to %s\n", outputFilePath)
		}
	}
}
