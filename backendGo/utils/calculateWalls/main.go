package main

import (
	"fmt"
	"encoding/json"
	"log"
	"os"
)

type Coordinates [][]float64 // Tablica tablic float64 (para współrzędnych: longitude, latitude)

type Geometry struct {
	Type        string     `json:"type"`
	Coordinates Coordinates `json:"coordinates"`
}

type Properties struct {
	ID        string `json:"@id"`
	City      string `json:"addr:city"`
	Street    string `json:"addr:street"`
	Postcode  string `json:"addr:postcode"`
	Building  string `json:"building"`
	Levels    int    `json:"building:levels"`
	Wheelchair string `json:"wheelchair"`
}

type Feature struct {
	Type       string     `json:"type"`
	Properties Properties `json:"properties"`
	Geometry   Geometry  `json:"geometry"`
	ID         string     `json:"id"`
}

type FeatureCollection struct {
	Type     string    `json:"type"`
	Features []Feature `json:"features"`
}

func main() {
	cwd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	// Read raw data from the JSON file
	data, err := os.ReadFile(cwd + "/buildingsData.json")
	if err != nil {
		log.Fatal("Failed to read data file", err)
	}

	// Declare a variable to store the parsed data
	var featureCollection FeatureCollection

	// Unmarshal the JSON data into featureCollection
	err = json.Unmarshal(data, &featureCollection)
	if err != nil {
		log.Fatal(err)
	}

	// Print out the coordinates for each feature
	for _, feature := range featureCollection.Features {
		fmt.Println("Coordinates for feature ID:", feature.ID)
		for _, coord := range feature.Geometry.Coordinates {
			fmt.Printf("Longitude: %f, Latitude: %f\n", coord[0], coord[1])
		}
	}
}