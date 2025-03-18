package tools

import (
	. "backendGo/types"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
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

func CalculateWalls(filePath string) {
	// Read the GeoJSON file
	data, err := os.ReadFile(filePath)
	if err != nil {
		log.Fatalf("Failed to read data file: %v", err)
	}

	// Parse the GeoJSON
	var featureCollection FeatureCollection
	err = json.Unmarshal(data, &featureCollection)
	if err != nil {
		log.Fatalf("Error parsing JSON: %v", err)
	}

	// Create a slice to hold all buildings
	var buildings []Building

	// Process each building (feature)
	for i, feature := range featureCollection.Features {
		buildingIndex := i + 1
		
		// Extract building name
		buildingName := fmt.Sprintf("Building %d", buildingIndex)
		if name, ok := feature.Properties["addr:housename"]; ok {
			buildingName = fmt.Sprintf("%v", name)
		}
		
		// Extract building height (levels)
		var heightInLevels float64 = 3 // Default height if not specified
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
		buildingOutput := Building{
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
		
		// Add building to the list
		buildings = append(buildings, buildingOutput)
	}
	
	// Convert slice to JSON
	outputJSON, err := json.MarshalIndent(buildings, "", "  ")
	if err != nil {
		log.Fatalf("Error creating JSON: %v", err)
	}
	
	// Write to file
	outputFilename := strings.Join(strings.Split(filePath, "_raw"), "")
	err = os.WriteFile(outputFilename, outputJSON, 0644)
	if err != nil {
		log.Fatalf("Error writing file %s: %v", outputFilename, err)
	}
	
	fmt.Printf("Saved all buildings to %s\n", outputFilename)
	fmt.Println("Processing complete")
}
