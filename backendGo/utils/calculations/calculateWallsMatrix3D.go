package calculations

import (
	. "backendGo/types"
	"bytes"
	"encoding/binary"
	"encoding/gob"
	"encoding/json" // Potrzebne do Unmarshal
	"fmt"
	"io"
	"log" // Używasz log.Fatalf
	"math"
	"os"
	"os/exec"
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
	ID         string          `json:"id,omitempty"`
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

func drawLine(matrix [][][]float64, x1, y1, z1, x2, y2, z2, heightLevels, wallIndex, sizeX, sizeY int) {
	dx := x2 - x1
	dy := y2 - y1
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
					if matrix[z][y][x1] >= 1000 && matrix[z][y][x1] != float64(1000 + wallIndex) {
						matrix[z][y][x1] = 10000 // Mark as corner
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
					if matrix[z][y1][x] >= 1000 && matrix[z][y1][x] != float64(1000 + wallIndex) {
						matrix[z][y1][x] = 10000 // Mark as corner
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
						if matrix[z][yIdx][prevXIdx] >= 1000 && matrix[z][yIdx][prevXIdx] != float64(1000 + wallIndex) {
							matrix[z][yIdx][prevXIdx] = 10000 // Mark as corner
						} else {
							matrix[z][yIdx][prevXIdx] = float64(1000 + wallIndex)
						}
					}
				}
			}
			if prevXIdx > xIdx && prevYIdx < yIdx  || prevXIdx > xIdx && prevYIdx > yIdx  {
				if xIdx >=0 && xIdx < sizeX && prevYIdx >= 0 && prevYIdx < sizeY {
					for z := 0; z <= z1; z++ {
						if matrix[z][prevYIdx][xIdx] >= 1000 && matrix[z][prevYIdx][xIdx] != float64(1000 + wallIndex) {
							matrix[z][prevYIdx][xIdx] = 10000 // Mark as corner
						} else {
							matrix[z][prevYIdx][xIdx] = float64(1000 + wallIndex)
						}
					}
				}
			} // walls continuity
			if xIdx >= 0 && xIdx < sizeX && yIdx >= 0 && yIdx < sizeY {
				for z := 0; z <= z1; z++ {
					if matrix[z][yIdx][xIdx] >= 1000 && matrix[z][yIdx][xIdx] != float64(1000 + wallIndex) {
						matrix[z][yIdx][xIdx] = 10000 // Mark as corner
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
			drawLine(matrix, i1, j1, z1, i2, j2, z2, heightLevels, wallsMapIndex, 250, 250)
			wallNormals = append(wallNormals, normal)
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

func addFloor(matrix [][][]float64, folderPath, filename string) error {
	path := filepath.Join(folderPath, filename)
    file, err := os.Create(path)
    if err != nil {
        return err
    }
    defer file.Close()

    for _, slice := range matrix {
        for _, row := range slice {
            for _, val := range row {
                err := binary.Write(file, binary.LittleEndian, val)
                if err != nil {
                    return err
                }
            }
        }
    }
    return nil
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

	// Sprawdzenie rozmiaru pliku (opcjonalne, ale dobre)
	fileInfo, err := file.Stat()
	if err != nil {
		return nil, fmt.Errorf("błąd odczytu informacji o pliku '%s': %w", path, err)
	}
	expectedSizeBytes := int64(z * y * x * 8) // 8 bajtów na float64
	if fileInfo.Size() != expectedSizeBytes {
		fmt.Printf("Ostrzeżenie: Rozmiar pliku '%s' (%d bajtów) różni się od oczekiwanego (%d bajtów dla %dx%dx%d float64).\n",
			path, fileInfo.Size(), expectedSizeBytes, z, y, x)
        // Możesz zdecydować, czy to jest błąd krytyczny
        // if fileInfo.Size() < expectedSizeBytes {
        //     return nil, fmt.Errorf("plik '%s' jest za mały", path)
        // }
	}


	data := make([][][]float64, z)
	for i := range data {
		data[i] = make([][]float64, y)
		for j := range data[i] {
			data[i][j] = make([]float64, x)
		}
	}

	byteOrder := binary.LittleEndian // Musi pasować do zapisu w Pythonie

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

// --- ZMODYFIKOWANA Funkcja do zapisu formatu GOB ---
// Teraz przyjmuje folderPath i filename osobno
func saveGobBinary(data interface{}, folderPath, filename string) error {
	// Utwórz folder, jeśli nie istnieje
	err := os.MkdirAll(folderPath, os.ModePerm)
	if err != nil {
		return fmt.Errorf("błąd tworzenia folderu wyjściowego '%s': %w", folderPath, err)
	}

	// Połącz ścieżkę folderu i nazwę pliku
	finalPath := filepath.Join(folderPath, filename)

	file, err := os.Create(finalPath)
	if err != nil {
		return fmt.Errorf("błąd tworzenia pliku Gob '%s': %w", finalPath, err)
	}
	defer file.Close()

	encoder := gob.NewEncoder(file)
	fmt.Printf("Zapisywanie danych do pliku Gob: %s\n", finalPath)
	err = encoder.Encode(data)
	if err == nil {
		fmt.Println("Zakończono zapis pliku Gob.")
	} else {
        return fmt.Errorf("błąd kodowania danych do formatu Gob w '%s': %w", finalPath, err)
    }
	return nil
}

	
func CalculateWallsMatrix3D(folderPath string, mapConfig MapConfig) {
	fmt.Println("Rozpoczynanie CalculateWallsMatrix3D...")

	// === Krok 1: Generowanie danych początkowych w Go ===
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

	// Generuj macierz i normalne
	matrix, wallNormals := generateBuildingMatrix(buildings, mapConfig.LatMin, mapConfig.LatMax, mapConfig.LonMin, mapConfig.LonMax, mapConfig.Size, mapConfig.HeightMaxLevels)

	// Opcjonalnie: Zapisz oryginalną macierz w formacie Gob (jeśli jest potrzebna)
	// saveBinary(matrix, folderPath, "wallsMatrix3D.bin")

	// === Krok 2: Zapisz macierz do formatu surowego dla Pythona ===
	rawInputForPythonFilename := "wallsMatrix3D_raw.bin"
	rawInputForPythonPath := filepath.Join(folderPath, rawInputForPythonFilename)
	err = saveRawBinary3D(matrix, folderPath, rawInputForPythonFilename)
	if err != nil {
		fmt.Printf("BŁĄD KRYTYCZNY: Nie udało się zapisać surowych danych 3D dla Pythona: %v\n", err)
		return // Przerwij, bo Python nie będzie miał danych
	}

	// === Krok 3: Wykonaj skrypt Pythona ===
	// Konfiguracja ścieżek i komend
	pythonScriptPath := "/home/kacper/workspace/react/raycheck/pythonscripts/learning/process_map.py" // Użyj absolutnej ścieżki lub skonfiguruj inaczej
	processedRawOutputFromPythonFilename := "wallsMatrix3D_processed.bin"
	processedRawOutputFromPythonPath := filepath.Join(folderPath, processedRawOutputFromPythonFilename)
    pythonCmd := "/home/kacper/workspace/react/raycheck/pythonscripts/venv/bin/python" // Lub "python", lub skonfiguruj

	// Sprawdź, czy skrypt istnieje
	if _, err := os.Stat(pythonScriptPath); os.IsNotExist(err) {
		fmt.Printf("BŁĄD KRYTYCZNY: Skrypt Pythona nie istnieje w: %s\n", pythonScriptPath)
		return
	}
    // Sprawdź komendę Pythona
    if _, err := exec.LookPath(pythonCmd); err != nil {
         fmt.Printf("BŁĄD KRYTYCZNY: Komenda '%s' nie znaleziona w PATH.\n", pythonCmd)
         return
    }

	fmt.Printf("Uruchamianie skryptu Pythona: %s\n", pythonScriptPath)
	fmt.Printf("  Input:  %s\n", rawInputForPythonPath)
	fmt.Printf("  Output: %s\n", processedRawOutputFromPythonPath)
	fmt.Printf("  Dims:   %d,%d,%d\n", mapConfig.HeightMaxLevels, mapConfig.Size, mapConfig.Size)

	// Budowanie i wykonanie komendy
	cmdArgs := []string{
		pythonScriptPath,
		"--input", rawInputForPythonPath,
		"--output", processedRawOutputFromPythonPath, // Skrypt Pythona musi tu zapisać!
		"--dims", fmt.Sprintf("%d,%d,%d", mapConfig.HeightMaxLevels, mapConfig.Size, mapConfig.Size),
	}
	cmd := exec.Command(pythonCmd, cmdArgs...)

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err = cmd.Run() // Uruchom i poczekaj
	stdoutStr := stdout.String()
	stderrStr := stderr.String()

	if len(stdoutStr) > 0 {
		fmt.Println("--- Python stdout ---")
		fmt.Print(stdoutStr)
		fmt.Println("---------------------")
	}
	if len(stderrStr) > 0 {
		fmt.Println("--- Python stderr ---")
		fmt.Print(stderrStr)
		fmt.Println("---------------------")
	}

	if err != nil {
		fmt.Printf("BŁĄD: Wykonanie skryptu Pythona nie powiodło się: %v\n", err)
		// Możesz zdecydować, czy kontynuować mimo błędu Pythona
		// return
	} else {
		fmt.Println("Skrypt Pythona zakończony (kod wyjścia 0).")
	}

	// === Krok 4: Sprawdź, czy Python wygenerował plik wyjściowy ===
	if _, err := os.Stat(processedRawOutputFromPythonPath); os.IsNotExist(err) {
		fmt.Printf("\nBŁĄD KRYTYCZNY: Skrypt Pythona NIE utworzył oczekiwanego pliku wyjściowego: '%s'\n", processedRawOutputFromPythonPath)
		fmt.Println("Sprawdź output Pythona (powyżej) w poszukiwaniu błędów.")
		return // Przerwij, bo nie ma przetworzonych danych
	} else if err != nil {
         fmt.Printf("\nBŁĄD KRYTYCZNY: Nie można sprawdzić statusu pliku wyjściowego Pythona '%s': %v\n", processedRawOutputFromPythonPath, err)
         return
    } else {
        fmt.Println("Plik wyjściowy Pythona znaleziony.")
    }

	// === Krok 5: Odczytaj przetworzone dane z pliku Pythona ===
	processedMatrix3D, err := loadRawBinary3D(processedRawOutputFromPythonPath, mapConfig.HeightMaxLevels, mapConfig.Size, mapConfig.Size)
	if err != nil {
		fmt.Printf("\nBŁĄD KRYTYCZNY: Nie udało się odczytać przetworzonych danych z '%s': %v\n", processedRawOutputFromPythonPath, err)
		return
	}

	// === Krok 6: Zapisz przetworzoną macierz do formatu Gob ===
	finalGobFilename := "wallsMatrix3D_floor.bin"
	err = saveBinary(processedMatrix3D, folderPath, finalGobFilename) // Używamy Twojej saveBinary
	if err != nil {
		fmt.Printf("\nBŁĄD KRYTYCZNY: Nie udało się zapisać finalnego pliku Gob '%s': %v\n", finalGobFilename, err)
		return
	}

	// === Krok 7: Zapisz wallNormals (jeśli nadal potrzebne) ===
	err = saveBinary(wallNormals, folderPath, "wallNormals3D.bin")
    if err != nil {
        fmt.Printf("\nBŁĄD: Nie udało się zapisać pliku wallNormals3D.bin: %v\n", err)
        // Zwykle nie przerywamy z tego powodu, ale logujemy błąd
    }


	fmt.Println("\nCalculateWallsMatrix3D zakończone sukcesem!")
}