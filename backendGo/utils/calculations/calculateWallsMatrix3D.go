package calculations

import (
	. "backendGo/types"
	"encoding/binary"
	"encoding/gob"
	"encoding/json" // Potrzebne do Unmarshal
	"fmt"
	"io"
	"log" // Używasz log.Fatalf
	"math"
	"os"
	"path/filepath"
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
}

type Geometry struct {
	Type        string        `json:"type"`
	Coordinates [][][]float64 `json:"coordinates"`
}

func calculateWalls(folderPath string) string {
	rawPath := filepath.Join(folderPath, "rawBuildings.json")
	data, err := os.ReadFile(rawPath)
	if err != nil {
		log.Fatalf("Failed to read data file: %v", err)
	}
	var featureCollection FeatureCollection
	err = json.Unmarshal(data, &featureCollection)
	if err != nil {
		log.Fatalf("Error parsing JSON: %v", err)
	}
	var buildings []Building
	for i, feature := range featureCollection.Features {
		buildingIndex := i + 1
		
		buildingName := fmt.Sprintf("Building %d", buildingIndex)
		if name, ok := feature.Properties["addr:housename"]; ok {
			buildingName = fmt.Sprintf("%v", name)
		}
		
		var heightInLevels float64 = 3
		if levels, ok := feature.Properties["building:levels"]; ok {
			switch v := levels.(type) {
			case float64:
				heightInLevels = v
			case int:
				heightInLevels = float64(v)
			case string:
				fmt.Sscanf(v, "%f", &heightInLevels)
			}
		}
		
		heightInMeters := heightInLevels * 3.0
		buildingOutput := Building{
			Name:   buildingName,
			Height: heightInMeters,
			Walls:  []Wall{},
		}
		if feature.Geometry.Type == "Polygon" {
			if len(feature.Geometry.Coordinates) > 0 {
				ring := feature.Geometry.Coordinates[0]
				for i := 0; i < len(ring); i++ {
				
					current := Point3D{
						X: ring[i][0],         
						Y: ring[i][1],          
						Z: heightInMeters,     
					}
					nextIdx := (i + 1) % len(ring)
					next := Point3D{
						X: ring[nextIdx][0],  
						Y: ring[nextIdx][1],   
						Z: heightInMeters,     
					}
					wall := Wall{Start: current, End: next}
					buildingOutput.Walls = append(buildingOutput.Walls, wall)
				}
			}
		}
		
		buildings = append(buildings, buildingOutput)
	}
	outputJSON, err := json.MarshalIndent(buildings, "", "  ")
	if err != nil {
		log.Fatalf("Error creating JSON: %v", err)
	}
	outputFilePath := filepath.Join(folderPath, "buildings.json")
	err = os.WriteFile(outputFilePath, outputJSON, 0644)
	if err != nil {
		log.Fatalf("Error writing file %s: %v", outputFilePath, err)
	}
	fmt.Printf("Saved all buildings to %s\n", outputFilePath)
	fmt.Println("Processing complete")
	return outputFilePath
}
func geoToMatrixIndex(lat, lon, latMin, latMax, lonMin, lonMax float64, size int) (int, int) {
	y := (lat - latMin) / (latMax - latMin) * float64(size-1)
	x := (lon - lonMin) / (lonMax - lonMin) * float64(size-1)
	i := int(math.Round(x))
	j := int(math.Round(y))
	return i, j
}

func AngleBetweenNormals(a, b Normal3D) float64 {
	dot := a.Nx*b.Nx + a.Ny*b.Ny + a.Nz*b.Nz
	lenA := math.Sqrt(a.Nx*a.Nx + a.Ny*a.Ny + a.Nz*a.Nz)
	lenB := math.Sqrt(b.Nx*b.Nx + b.Ny*b.Ny + b.Nz*b.Nz)

	if lenA == 0 || lenB == 0 {
		return 0
	}

	cosTheta := dot / (lenA * lenB)

	cosTheta = math.Max(-1, math.Min(1, cosTheta))

	theta := math.Acos(cosTheta)
	return theta * 180 / math.Pi 
}

func drawLine(matrix [][][]float64, wallNormals []Normal3D, x1, y1, z1, x2, y2, z2, heightLevels, wallIndex, sizeX, sizeY int) {
	dx := x2 - x1
	dy := y2 - y1
	minCornerAngle := 40.0
	maxCornerAngle := 120.0
	if z1 >= heightLevels {
		z1 = heightLevels - 1
	}
	if z2 >= heightLevels {
		z2 = heightLevels - 1
	}
	if x1 == x2 && x1 >= 0 && x1 < sizeX{
		if y1 > y2 {
			y1, y2 = y2, y1 
		}
		for y := y1; y <= y2; y++ {
			if y >= 0 && y < sizeY{
				for z := 0; z <= z1; z++ {
					if matrix[z][y][x1] >= 1000  && matrix[z][y][x1] < 5000 && matrix[z][y][x1] != float64(1000 + wallIndex) {
						wallA := wallNormals[int(matrix[z][y][x1]) - 1000]
						wallB := wallNormals[wallIndex]
						angle := AngleBetweenNormals(wallA, wallB)
						if angle > minCornerAngle && angle < maxCornerAngle{
							matrix[z][y][x1] = 10000
						} else {
							matrix[z][y][x1] =  float64(1000 + wallIndex)
						}
					} else {
						matrix[z][y][x1] = float64(1000 + wallIndex)
					}
				}
			}
		}
	} else if y1 == y2  && y1 >= 0 && y1 < sizeY {
		if x1 > x2 {
			x1, x2 = x2, x1
		}
		for x := x1; x <= x2; x++ {
			if x >= 0 && x < sizeX {
				for z := 0; z <= z1; z++ {
					if matrix[z][y1][x] >= 1000 && matrix[z][y1][x] < 5000 && matrix[z][y1][x] != float64(1000 + wallIndex) {
						wallA := wallNormals[int(matrix[z][y1][x]) - 1000]
						wallB := wallNormals[wallIndex]
						angle := AngleBetweenNormals(wallA, wallB)
						if angle > minCornerAngle && angle < maxCornerAngle{
							matrix[z][y1][x] = 10000
						} else {
							matrix[z][y1][x] =  float64(1000 + wallIndex)
						}
					} else {
						matrix[z][y1][x] = float64(1000 + wallIndex)
					}
				}
			}
		}
	} else {
		steps := int(math.Max(math.Abs(float64(dx)), math.Abs(float64(dy))))
		prevXIdx := x1
		prevYIdx := y1
		for j := 0; j <= steps; j++ {
			x := x1 + dx*j/steps
			y := y1 + dy*j/steps
			xIdx := x
			yIdx := y
			if prevXIdx < xIdx && prevYIdx < yIdx || prevXIdx < xIdx && prevYIdx > yIdx  {
				if yIdx >= 0 && yIdx < sizeY && prevXIdx >= 0 && prevXIdx < sizeX {
					for z := 0; z <= z1; z++ {
						if matrix[z][yIdx][prevXIdx] >= 1000 && matrix[z][yIdx][prevXIdx] < 5000 && matrix[z][yIdx][prevXIdx] != float64(1000 + wallIndex) {
							wallA := wallNormals[int(matrix[z][yIdx][prevXIdx]) - 1000]
							wallB := wallNormals[wallIndex]
							angle := AngleBetweenNormals(wallA, wallB)
							if angle > minCornerAngle && angle < maxCornerAngle{
								matrix[z][yIdx][prevXIdx] = 10000
							} else {
								matrix[z][yIdx][prevXIdx] =  float64(1000 + wallIndex)
							}
						} else {
							matrix[z][yIdx][prevXIdx] = float64(1000 + wallIndex)
						}
					}
				}
			}
			if prevXIdx > xIdx && prevYIdx < yIdx  || prevXIdx > xIdx && prevYIdx > yIdx  {
				if xIdx >=0 && xIdx < sizeX && prevYIdx >= 0 && prevYIdx < sizeY {
					for z := 0; z <= z1; z++ {
						if matrix[z][prevYIdx][xIdx] >= 1000 && matrix[z][prevYIdx][xIdx] < 5000 && matrix[z][prevYIdx][xIdx] != float64(1000 + wallIndex) {
							wallA := wallNormals[int(matrix[z][prevYIdx][xIdx]) - 1000]
							wallB := wallNormals[wallIndex]
							angle := AngleBetweenNormals(wallA, wallB)
							if angle > minCornerAngle && angle < maxCornerAngle {
								matrix[z][prevYIdx][xIdx] = 10000
							} else {
								matrix[z][prevYIdx][xIdx] =  float64(1000 + wallIndex)
							}
						} else {
							matrix[z][prevYIdx][xIdx] = float64(1000 + wallIndex)
						}
					}
				}
			} // walls continuity
			if xIdx >= 0 && xIdx < sizeX && yIdx >= 0 && yIdx < sizeY {
				for z := 0; z <= z1; z++ {
					if matrix[z][yIdx][xIdx] >= 1000 && matrix[z][yIdx][xIdx] < 5000 && matrix[z][yIdx][xIdx] != float64(1000 + wallIndex) {
						wallA := wallNormals[int(matrix[z][yIdx][xIdx]) - 1000]
						wallB := wallNormals[wallIndex]
						angle := AngleBetweenNormals(wallA, wallB)
						println(angle)
						if angle > minCornerAngle && angle < maxCornerAngle {
							matrix[z][yIdx][xIdx] = 10000
						} else {
							matrix[z][yIdx][xIdx] =  float64(1000 + wallIndex)
						}
					} else {
						matrix[z][yIdx][xIdx] = float64(1000 + wallIndex)
					}
				}
			}
			prevXIdx = xIdx
			prevYIdx = yIdx
		} 
	}
}

func calculateNormal3D( x1, y1, z1, x2, y2, z2 int) Normal3D {
	dx := x2 - x1
	dy := y2 - y1
	length := math.Hypot(float64(dx), float64(dy))
		if length == 0  {
			return Normal3D{Nx:0, Ny:0, Nz:0}
		}
		nx := -float64(dy)/length
		ny := float64(dx)/length
		return Normal3D{Nx:nx, Ny:ny, Nz:0}
}

func generateBuildingMatrix(buildings []Building, latMin, latMax, lonMin, lonMax float64, size, heightLevels int) ([][][]float64, []Normal3D){
	matrix := make([][][]float64, heightLevels)
	wallNormals := []Normal3D{}
	for z := range matrix {
		matrix[z] = make([][]float64, size)
		for y := range matrix[z] {
			matrix[z][y] = make([]float64, size)
			for x := range matrix[z][y] {
				matrix[z][y][x] = -150.0
			}
		}
	}
	wallsMapIndex := 0
	wallHeights := make(map[int]int)
	for _, building := range buildings {
		for _, wall := range building.Walls {
			i1, j1 := geoToMatrixIndex(wall.Start.Y, wall.Start.X, latMin, latMax, lonMin, lonMax, size)
			i2, j2 := geoToMatrixIndex(wall.End.Y, wall.End.X, latMin, latMax, lonMin, lonMax, size)			
			z1 := int(math.Round(wall.Start.Z))
			z2 := int(math.Round(wall.End.Z))
			normal := calculateNormal3D( i1, j1, z1, i2, j2, z2)
			if normal.Nx == 0 && normal.Ny == 0 {
				continue
			} 
			wallNormals = append(wallNormals, normal)
			drawLine(matrix, wallNormals, i1, j1, z1, i2, j2, z2, heightLevels, wallsMapIndex, size, size)
			wallHeights[wallsMapIndex] = z1
			wallsMapIndex++
		}
	}
	fmt.Printf("walls: %v \n", wallsMapIndex)
	fmt.Printf("normals: %v \n", len(wallNormals))
	return matrix, wallNormals
}

func saveBinary(data interface{}, folderPath, filename string) error {
	finalPath := filepath.Join(folderPath, filename)
	file, err := os.Create(finalPath)
	if err != nil {
		return err
	}
	defer file.Close()
	
	encoder := gob.NewEncoder(file)
	return encoder.Encode(data)
}

func LoadMatrixBinary(path string, data interface{}) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()
	decoder := gob.NewDecoder(file)
	return decoder.Decode(data)
}

func saveRawBinary3D(matrix [][][]float64, folderPath, filename string) error {
	path := filepath.Join(folderPath, filename)
	file, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("błąd tworzenia pliku %s: %w", path, err)
	}
	defer file.Close()
	
	// Upewnij się, że zapisujesz w LittleEndian (Python to odczyta)
	byteOrder := binary.LittleEndian
	
	fmt.Printf("Zapisywanie surowych danych 3D do: %s\n", path)
	for _, slice := range matrix {
		for _, row := range slice {
			for _, val := range row {
				err := binary.Write(file, byteOrder, val)
				if err != nil {
					// Zwróć bardziej szczegółowy błąd
					return fmt.Errorf("błąd zapisu wartości %f do %s: %w", val, path, err)
				}
			}
		}
	}
	fmt.Println("Zakończono zapis surowych danych 3D.")
	return nil
}

// loadRawBinary2D odczytuje surowy plik binarny (stworzony przez Python) do [][][]float64
func loadRawBinary3D(path string, z, y, x int) ([][][]float64, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("błąd otwarcia pliku wejściowego '%s': %w", path, err)
	}
	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil {
		return nil, fmt.Errorf("błąd odczytu informacji o pliku '%s': %w", path, err)
	}
	expectedSizeBytes := int64(z * y * x * 8) 
	if fileInfo.Size() != expectedSizeBytes {
		fmt.Printf("Ostrzeżenie: Rozmiar pliku '%s' (%d bajtów) różni się od oczekiwanego (%d bajtów dla %dx%dx%d float64).\n",
			path, fileInfo.Size(), expectedSizeBytes, z, y, x)
      
	}


	data := make([][][]float64, z)
	for i := range data {
		data[i] = make([][]float64, y)
		for j := range data[i] {
			data[i][j] = make([]float64, x)
		}
	}

	byteOrder := binary.LittleEndian 

	fmt.Printf("Odczytywanie danych 3D z: %s (%dx%dx%d)\n", path, z, y, x)
	expectedReads := z * y * x
	reads := 0
	for zi := 0; zi < z; zi++ {
		for yi := 0; yi < y; yi++ {
			for xi := 0; xi < x; xi++ {
				var value float64
				err := binary.Read(file, byteOrder, &value)
				if err != nil {
					if err == io.EOF {
						return nil, fmt.Errorf("niespodziewany koniec pliku (EOF) w '%s' po odczytaniu %d z %d oczekiwanych liczb. Plik jest za krótki.", path, reads, expectedReads)
					}
					return nil, fmt.Errorf("błąd odczytu danych 3D dla [%d][%d][%d] z '%s': %w", zi, yi, xi, path, err)
				}
				data[zi][yi][xi] = value
				reads++
			}
		}
	}
	fmt.Println("Zakończono odczyt surowych danych 3D.")
	return data, nil
}

	
func CalculateWallsMatrix3D(folderPath string, mapConfig MapConfig) {
	fmt.Println("Rozpoczynanie CalculateWallsMatrix3D...")
	buildingsFilePath := calculateWalls(folderPath)
	data, err := os.ReadFile(buildingsFilePath)
	if err != nil {
		log.Fatalf("Nie udało się odczytać pliku danych budynków '%s': %v", buildingsFilePath, err)
	}
	var buildings []Building
	err = json.Unmarshal(data, &buildings)
	if err != nil {
		log.Fatalf("Błąd parsowania JSON z '%s': %v", buildingsFilePath, err)
	}

	matrix, wallNormals := generateBuildingMatrix(buildings, mapConfig.LatMin, mapConfig.LatMax, mapConfig.LonMin, mapConfig.LonMax, mapConfig.Size, mapConfig.HeightMaxLevels)
	
 	saveBinary(matrix, folderPath, "wallsMatrix3D.bin")
	saveBinary(wallNormals, folderPath, "wallNormals3D.bin")

 	rawInputForPythonFilename := "wallsMatrix3D_raw.bin"
 	err = saveRawBinary3D(matrix, folderPath, rawInputForPythonFilename)
 	if err != nil {
 		fmt.Printf("BŁĄD KRYTYCZNY: Nie udało się zapisać surowych danych 3D dla Pythona: %v\n", err)
 		return 
 	}
 
 	processedRawOutputFromPythonFilename := "wallsMatrix3D_processed.bin"
 	processedRawOutputFromPythonPath := filepath.Join(folderPath, processedRawOutputFromPythonFilename)
 	processedMatrix3D, err := loadRawBinary3D(processedRawOutputFromPythonPath, mapConfig.HeightMaxLevels, mapConfig.Size, mapConfig.Size)
 	if err != nil {
 		fmt.Printf("\nBŁĄD KRYTYCZNY: Nie udało się odczytać przetworzonych danych z '%s': %v\n", processedRawOutputFromPythonPath, err)
 		return
 	}
 
 	
 	finalGobFilename := "wallsMatrix3D_floor.bin"
 	err = saveBinary(processedMatrix3D, folderPath, finalGobFilename)  
 	if err != nil {
 		fmt.Printf("\nBŁĄD KRYTYCZNY: Nie udało się zapisać finalnego pliku Gob '%s': %v\n", finalGobFilename, err)
 		fmt.Printf("\nBŁĄD KRYTYCZNY: Nie udało się zapisać finalnego pliku Gob '%s': %v\n", "wallsMatrix3D.bin", err)
 		return
 	}
 	err = saveBinary(wallNormals, folderPath, "wallNormals3D.bin")
     if err != nil {
         fmt.Printf("\nBŁĄD: Nie udało się zapisać pliku wallNormals3D.bin: %v\n", err)
         
     }
	}