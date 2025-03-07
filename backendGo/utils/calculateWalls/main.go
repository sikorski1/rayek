package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
)

// GeoJSON structures with flexible property types
type FeatureCollection struct {
	Type     string    `json:"type"`
	Features []Feature `json:"features"`
}

type Feature struct {
	Type       string          `json:"type"`
	Properties map[string]any  `json:"properties"`
	Geometry   Geometry        `json:"geometry"`
	ID         string          `json:"id,omitempty"`
}

type Geometry struct {
	Type        string        `json:"type"`
	Coordinates [][][]float64 `json:"coordinates"`
}

// Point3D represents a point in 3D space
type Point3D struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
	Z float64 `json:"z"`
}

// Wall represents a plane between two 3D points
type Wall struct {
	Start Point3D `json:"start"`
	End   Point3D `json:"end"`
}

// BuildingOutput is the structure for each building
type BuildingOutput struct {
	Name   string `json:"name"`
	Height float64 `json:"height"`
	Walls  []Wall `json:"walls"`
}

func main() {
	// Read the GeoJSON file
	data, err := os.ReadFile("buildingsData.json")
	if err != nil {
		log.Fatalf("Failed to read data file: %v", err)
	}

	// Parse the GeoJSON
	var featureCollection FeatureCollection
	err = json.Unmarshal(data, &featureCollection)
	if err != nil {
		log.Fatalf("Error parsing JSON: %v", err)
	}

	// Create a map to hold all buildings
	buildings := make(map[string]BuildingOutput)

	// Process each building (feature)
	for i, feature := range featureCollection.Features {
		buildingIndex := i + 1
		
		// Extract building name
		buildingName := fmt.Sprintf("Building %d", buildingIndex)
		if name, ok := feature.Properties["addr:housename"]; ok {
			buildingName = fmt.Sprintf("%v", name)
		}
		
		// Extract building height (levels)
		var heightInLevels float64 = 1 // Default height if not specified
		if levels, ok := feature.Properties["building:levels"]; ok {
			// Convert the interface{} to float64
			switch v := levels.(type) {
			case float64:
				heightInLevels = v
			case int:
				heightInLevels = float64(v)
			case string:
				fmt.Sscanf(v, "%f", &heightInLevels)
			}
		}
		
		// Assume each level is 3 meters high
		heightInMeters := heightInLevels * 3.0
		
		// Create output structure
		buildingOutput := BuildingOutput{
			Name:   buildingName,
			Height: heightInMeters,
			Walls:  []Wall{},
		}

		// Process geometry if it's a polygon (building walls)
		if feature.Geometry.Type == "Polygon" {
			// Process exterior wall (first ring)
			if len(feature.Geometry.Coordinates) > 0 {
				ring := feature.Geometry.Coordinates[0]
				
				// Create walls between consecutive points
				for i := 0; i < len(ring); i++ {
					// Current point
					current := Point3D{
						X: ring[i][0],          // Longitude
						Y: ring[i][1],          // Latitude
						Z: heightInMeters,      // Height in meters
					}
					
					// Next point (wrapping around to the first point if we're at the end)
					nextIdx := (i + 1) % len(ring)
					next := Point3D{
						X: ring[nextIdx][0],    // Longitude
						Y: ring[nextIdx][1],    // Latitude
						Z: heightInMeters,      // Height in meters
					}
					
					// Create a wall between these two points
					wall := Wall{Start: current, End: next}
					buildingOutput.Walls = append(buildingOutput.Walls, wall)
				}
			}
		}
		
		// Add building to the map with key "building1", "building2", etc.
		buildingKey := "building" + strconv.Itoa(buildingIndex)
		buildings[buildingKey] = buildingOutput
	}
	
	// Convert map to JSON
	outputJSON, err := json.MarshalIndent(buildings, "", "  ")
	if err != nil {
		log.Fatalf("Error creating JSON: %v", err)
	}
	
	// Write to file
	outputFilename := "buildings.json"
	err = os.WriteFile(outputFilename, outputJSON, 0644)
	if err != nil {
		log.Fatalf("Error writing file %s: %v", outputFilename, err)
	}
	
	fmt.Printf("Saved all buildings to %s\n", outputFilename)
	fmt.Println("Processing complete")
}