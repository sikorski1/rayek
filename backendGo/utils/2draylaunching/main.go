package main

import (
	. "backendGo/types"
	. "backendGo/utils/calculations"
	"fmt"
	"image/png"
	"math"
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
}

func NewRayLaunching(matrixDimensions Point, tPos Point, tPower float64, tFreq float64, rFactor float64, wallPos []Vector) *RayLaunching {
	step := 0.01
	numberOfRays := 360
	numberOfInteracitons := 3
	minimalPower := -70
	rows := int(matrixDimensions.Y*(1/step))+1
	cols := int(matrixDimensions.X*(1/step))+1
	Map := make([][]float64, rows)
	for i := range Map {
		Map[i] = make([]float64, cols)
	}
	
	setWallsIn2DMap(&Map, wallPos, step)
	return &RayLaunching{
		Step:             step,
		NumberOfRays:     numberOfRays,
		NumberOfInteracitons:numberOfInteracitons,
		MinimalPower:  minimalPower, 
		TransmitterPos:   tPos,
		TransmitterPower: tPower,
		TransmitterFreq:  tFreq,
		WaveLength:       299792458 / (tFreq * math.Pow(10, 9)),
		ReflectionFactor: rFactor,
		Map: Map,
		
	}
}

func (rl *RayLaunching) calculateRayLaunching() {
	for i := range(rl.NumberOfRays) {
		iteration := 1
		currInteractions := 0
		currPower := 0.0
		dx, dy := math.Cos(float64(i)), math.Sin(float64(i))
		for currInteractions <= rl.NumberOfInteracitons && currPower >= rl.MinimalPower {
			x:=rl.TransmitterPos.X + dx * float64(iteration) //todo
			y:=rl.TransmitterPos.Y + dy * float64(iteration) //todo
			currPoint := Point{X:x,Y:y}
		}
	}
}


func setWallsIn2DMap(Map *[][]float64, walls []Vector, step float64) {
	for _, wall := range walls {
		x1, y1 := wall.A.X, wall.A.Y
		x2, y2 := wall.B.X, wall.B.Y
		x1Idx := int(math.Round(x1 / step))
		y1Idx := int(math.Round(y1 / step))
		x2Idx := int(math.Round(x2 / step))
		y2Idx := int(math.Round(y2 / step))
		if x1 == x2 {
			if y1 > y2 {
				y1Idx, y2Idx = y2Idx, y1Idx 
			}
			for y := y1Idx; y <= y2Idx; y++ {
				(*Map)[y][x1Idx] = 1000
			}
		} else if y1 == y2 {
			if x1 > x2 {
				x1Idx, x2Idx = x2Idx, x1Idx
			}
			for x := x1Idx; x <= x2Idx; x++ {
				(*Map)[y1Idx][x] = 1000
			}
		} else {
			dx := x2 - x1
			dy := y2 - y1
			steps := int(math.Max(math.Abs(dx/step), math.Abs(dy/step)))
			for i := 0; i <= steps; i++ {
				x := x1 + (dx*float64(i))/float64(steps)
				y := y1 + (dy*float64(i))/float64(steps)
				xIdx := int(math.Round(x / step))
				yIdx := int(math.Round(y / step))
				(*Map)[yIdx][xIdx] = 1000
			} 
		}
	}
}


func main() {
	start := time.Now()
	matrixDimensions := Point{X:40, Y:40}
	transmitterPos := Point{X:5, Y:20}
	transmitterPower := 5.0 
	transmitterFreq := 3.4  
	reflectionFactor := 0.7
	walls := []Vector{{A:Point{X:0,Y:3}, B:Point{X:3,Y:6}}, {A:Point{X:1,Y:3}, B:Point{X:6,Y:3}}, {A:Point{X:6,Y:10}, B:Point{X:12,Y:12}},{A:Point{X:25,Y:10}, B:Point{X:25,Y:15}},{A:Point{X:5,Y:30}, B:Point{X:10,Y:35}},{A:Point{X:23,Y:36}, B:Point{X:25,Y:39}},{A:Point{X:1,Y:24}, B:Point{X:1,Y:26}},{A:Point{X:1,Y:28}, B:Point{X:1,Y:30}},{A:Point{X:1,Y:37}, B:Point{X:1,Y:40}},{A:Point{X:35,Y:36}, B:Point{X:30,Y:28}},{A:Point{X:40,Y:1}, B:Point{X:36,Y:2}},{A:Point{X:24,Y:3}, B:Point{X:25,Y:6}},{A:Point{X:16,Y:21}, B:Point{X:18,Y:22}},{A:Point{X:12,Y:18}, B:Point{X:12,Y:20}},{A:Point{X:18,Y:36}, B:Point{X:12,Y:36}}}
	raylaunching := NewRayLaunching(matrixDimensions, transmitterPos, transmitterPower, transmitterFreq, reflectionFactor, walls)
	fmt.Printf("%v \n", raylaunching.Map[32][66])
	stop := time.Since(start)
	fmt.Printf("Computation time: %v \n", stop)
	heatmap := GenerateHeatmap(raylaunching.Map)
    f, _ := os.Create("heatmap.png")
    defer f.Close()
    png.Encode(f, heatmap)
}