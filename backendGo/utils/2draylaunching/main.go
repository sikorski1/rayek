package main

import (
	. "backendGo/types"
	. "backendGo/utils/calculations"
	"encoding/json"
	"fmt"
	"image/png"
	"math"
	"math/cmplx"
	"os"
	"time"
)

type RayLaunching struct {
	Step                 float64
	NumberOfRays		 int
	NumberOfInteracitons int
	MinimalPower         float64
	TransmitterPos       Point
	TransmitterPower     float64
	TransmitterFreq      float64
	WaveLength           float64
	ReflectionFactor     float64
	Map                  [][]float64
	WallNormals  []Normal
	WallMapNumber int
}

func NewRayLaunching(matrixDimensions Point, tPos Point, tPower float64, tFreq float64, rFactor float64, wallPos []Vector) *RayLaunching {
	step := 0.1
	numberOfRays := 36000
	numberOfInteractions := 4
	minimalPower := -150.0
	wallMapNumber := 1000
	rows := int(matrixDimensions.Y*(1/step))+1
	cols := int(matrixDimensions.X*(1/step))+1
	Normal := make([]Normal, len(wallPos))
	Map := make([][]float64, rows)
	for i := range Map {
    	Map[i] = make([]float64, cols)
    	for j := range Map[i] {
       	 	Map[i][j] = -150
    	}
	}

	
	setWallsIn2DMap(&Map, &Normal, wallPos, step, wallMapNumber)
	return &RayLaunching{
		Step:             step,
		NumberOfRays:     numberOfRays,
		NumberOfInteracitons:numberOfInteractions,
		MinimalPower:  minimalPower, 
		TransmitterPos:   tPos,
		TransmitterPower: tPower,
		TransmitterFreq:  tFreq,
		WaveLength:       299792458 / (tFreq * math.Pow(10, 9)),
		ReflectionFactor: rFactor,
		Map: Map,
		WallNormals: Normal,
		WallMapNumber:wallMapNumber,
	}
}

func (rl *RayLaunching) calculateRayLaunching() {
	maxSizeX := (float64(len(rl.Map[0]))-1)*rl.Step
	maxSizeY := (float64(len(rl.Map))-1)*rl.Step
	for i := range rl.NumberOfRays {
		currInteractions := 0
		currPower := 0.0
		dRadians := (2*math.Pi)/float64(rl.NumberOfRays)*float64(i)  // calculating angle change for every ray
		dx, dy := math.Cos(dRadians)*rl.Step, math.Sin(dRadians)*rl.Step // calculating x, y change for every step
		dx, dy = math.Round(dx*1e15)/1e15, math.Round(dy*1e15)/1e15 // floating point numbers correction
		x, y := rl.TransmitterPos.X + dx, rl.TransmitterPos.Y + dy
		currWallIndex := 0
		for (x >= 0 && x <= maxSizeX) && (y >= 0 && y <= maxSizeY) && currInteractions < rl.NumberOfInteracitons && currPower >= rl.MinimalPower {
			xIdx := int(math.Round(x / rl.Step))
			yIdx := int(math.Round(y / rl.Step))
			index := int(rl.Map[yIdx][xIdx])
			// check if there is wall and if its diffrent from previous
			if index >= rl.WallMapNumber && index != currWallIndex+rl.WallMapNumber {
				currWallIndex = index-rl.WallMapNumber
				nx, ny := rl.WallNormals[currWallIndex].Nx, rl.WallNormals[currWallIndex].Ny
				dot := 2 * (dx*nx + dy*ny)
				dx = dx - dot*nx
				dy = dy - dot*ny
				currInteractions++
				// fmt.Printf("ray%v: %v \n", i,rl.Map[yIdx][xIdx])
			} else {
				H := CalculateTransmittance(Point{X:rl.TransmitterPos.X,Y:rl.TransmitterPos.Y},Point{X:x,Y:y},rl.WaveLength,math.Pow(rl.ReflectionFactor, float64(currInteractions)))
				currPower = 10*math.Log10(rl.TransmitterPower) + 20*math.Log10(cmplx.Abs(H))
				if rl.Map[yIdx][xIdx] != -150 {
					existingPowerLin := math.Pow(10, rl.Map[yIdx][xIdx]/10)
					currPowerLin := math.Pow(10, currPower/10)
					newPowerDb := 10 * math.Log10(existingPowerLin + currPowerLin)
					rl.Map[yIdx][xIdx] = newPowerDb
					
				} else {
					rl.Map[yIdx][xIdx] = currPower
				}
			}	
			x += dx
			y += dy
			
		}
	}
}


func setWallsIn2DMap(Map *[][]float64, WallNormals  *[]Normal, walls []Vector, step float64, wallMapNumber int) {
	for i, wall := range walls {
		x1, y1 := wall.A.X, wall.A.Y
		x2, y2 := wall.B.X, wall.B.Y
		x1Idx := int(math.Round(x1 / step))
		y1Idx := int(math.Round(y1 / step))
		x2Idx := int(math.Round(x2 / step))
		y2Idx := int(math.Round(y2 / step))
		
		dx := x2 - x1
		dy := y2 - y1
		length := math.Hypot(dx, dy)
		if length != 0  {
			nx := -dy/length
			ny := dx/length
			(*WallNormals)[i] = Normal{Nx:nx, Ny:ny}
		}
		if x1 == x2 {
			if y1 > y2 {
				y1Idx, y2Idx = y2Idx, y1Idx 
			}
			for y := y1Idx; y <= y2Idx; y++ {
				(*Map)[y][x1Idx] = float64(wallMapNumber + i)
			}
		} else if y1 == y2 {
			if x1 > x2 {
				x1Idx, x2Idx = x2Idx, x1Idx
			}
			for x := x1Idx; x <= x2Idx; x++ {
				(*Map)[y1Idx][x] = float64(wallMapNumber + i)
			}
		} else {
			steps := int(math.Max(math.Abs(dx/step), math.Abs(dy/step)))
			prevXIdx := int(math.Round(x1 / step))
			prevYIdx := int(math.Round(y1 / step))
			for j := 0; j <= steps; j++ {
				x := x1 + (dx*float64(j))/float64(steps)
				y := y1 + (dy*float64(j))/float64(steps)
				xIdx := int(math.Round(x / step))
				yIdx := int(math.Round(y / step))
				if prevXIdx < xIdx && prevYIdx < yIdx || prevXIdx < xIdx && prevYIdx > yIdx  {
					(*Map)[yIdx][prevXIdx] = float64(wallMapNumber + i)
				}
				if prevXIdx > xIdx && prevYIdx < yIdx  || prevXIdx > xIdx && prevYIdx > yIdx  {
					(*Map)[prevYIdx][xIdx] = float64(wallMapNumber + i)
				} // walls continuity
				(*Map)[yIdx][xIdx] = float64(wallMapNumber + i)
				prevXIdx = xIdx
				prevYIdx = yIdx
			} 
		}
	}
}
func SaveMapToCSV(Map [][]float64, filename string) error {
    file, err := os.Create(filename)
    if err != nil {
        return err
    }
    defer file.Close()

    for _, row := range Map {
        for i, value := range row {
            if i > 0 {
                file.WriteString(",")
            }
            file.WriteString(fmt.Sprintf("%f", value))
        }
        file.WriteString("\n")
    }
    return nil
}

func main() {
	start := time.Now()
	matrixDimensions := Point{X:40, Y:40}
	transmitterPos := Point{X:20, Y:10}
	transmitterPower := 5.0 
	transmitterFreq := 2.4
	reflectionFactor := 0.8
	walls := []Vector{{A:Point{X:0,Y:3}, B:Point{X:3,Y:6}}, {A:Point{X:1,Y:3}, B:Point{X:6,Y:3}}, {A:Point{X:12,Y:12}, B:Point{X:6,Y:10}},{A:Point{X:25,Y:10}, B:Point{X:25,Y:15}},{A:Point{X:10,Y:36}, B:Point{X:5,Y:30}},{A:Point{X:23,Y:36}, B:Point{X:25,Y:39}},{A:Point{X:1,Y:24}, B:Point{X:1,Y:26}},{A:Point{X:1,Y:28}, B:Point{X:1,Y:30}},{A:Point{X:1,Y:37}, B:Point{X:1,Y:40}},{A:Point{X:35,Y:36}, B:Point{X:30,Y:28}},{A:Point{X:40,Y:1}, B:Point{X:36,Y:2}},{A:Point{X:24,Y:3}, B:Point{X:25,Y:6}},{A:Point{X:16,Y:21}, B:Point{X:18,Y:22}},{A:Point{X:12,Y:18}, B:Point{X:12,Y:20}},{A:Point{X:18,Y:36}, B:Point{X:12,Y:36}}}
	raylaunching := NewRayLaunching(matrixDimensions, transmitterPos, transmitterPower, transmitterFreq, reflectionFactor, walls)
	raylaunching.calculateRayLaunching()
	stop := time.Since(start)
	fmt.Printf("Computation time: %v \n", stop)
	 // Zapisz mapę do pliku CSV
	 err := SaveMapToCSV(raylaunching.Map, "ray_map.csv")
	 if err != nil {
		 fmt.Printf("Error saving map: %v\n", err)
	 }
	 
	 // Zapisz też pozycję nadajnika i ściany
	 saveConfigToJSON(transmitterPos, walls, "ray_config.json")
	 
	 // Możesz też wygenerować obrazek jak dotychczas
	 heatmap := GenerateHeatmap(raylaunching.Map)
	 f, _ := os.Create("heatmap.png")
	 defer f.Close()
	 png.Encode(f, heatmap)
 }
 
 // Funkcja do zapisywania konfiguracji
 func saveConfigToJSON(transmitter Point, walls []Vector, filename string) error {
	 type Config struct {
		 TransmitterX float64   `json:"tx"`
		 TransmitterY float64   `json:"ty"`
		 Walls        [][]float64 `json:"walls"`
	 }
	 
	 wallsData := make([][]float64, len(walls))
	 for i, wall := range walls {
		 wallsData[i] = []float64{
			 wall.A.X, wall.A.Y, wall.B.X, wall.B.Y,
		 }
	 }
	 
	 config := Config{
		 TransmitterX: transmitter.X,
		 TransmitterY: transmitter.Y,
		 Walls:        wallsData,
	 }
	 
	 data, err := json.Marshal(config)
	 if err != nil {
		 return err
	 }
	 
	 return os.WriteFile(filename, data, 0644)
 }