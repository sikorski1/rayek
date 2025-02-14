package main

import (
	"fmt"
	"math"
)

type Point struct {
	X, Y float64
}

type Vector struct {
	A, B Point
}

type RayTracing struct {
	Step                 float64
	TransmitterPos       Point
	TransmitterPower     float64
	TransmitterFreq      float64
	WaveLength           float64
	ReflectionFactor     float64
	Walls                []Vector
	PowerMap             [][]float64
	Matrix               [][]Point
	MirroredTransmitters []Point
}
//create RayTracingObject
func NewRayTracing(matrixDimensions Point, tPos Point, tPower float64, tFreq float64, rFactor float64, wallPos []Vector) *RayTracing {
	step := 0.1
	rows := int(matrixDimensions.Y*(1/step))+1
	cols := int(matrixDimensions.X*(1/step))+1
	powerMap := make([][]float64, rows)
	//powermap
	for i := range powerMap {
		powerMap[i] = make([]float64, cols)
	}
	//matrix
	matrix := make([][]Point, rows)
	for i := range matrix {
		matrix[i] = make([]Point, cols)
		for j := range matrix[i] {
			matrix[i][j] = Point{X: math.Round(float64(j) * step * 10)/10, Y: math.Round(float64(i) * step * 10)/10}
		}
	}
	//mirred transmitters
	mirroredTransmitters := make([]Point, len(wallPos))
	for i, wall := range wallPos {
		if wall.A.X == wall.B.X {
			mirroredTransmitters[i].Y = tPos.Y
			distance := math.Abs(wall.A.X - tPos.X)
			if tPos.X < wall.A.X {
				mirroredTransmitters[i].X = wall.A.X + distance
			} else if tPos.X > wall.A.X {
				mirroredTransmitters[i].X = wall.A.X - distance
			} else {
				mirroredTransmitters[i].X = wall.A.X
			}
			continue
 		}
		if wall.A.Y == wall.B.Y {
			mirroredTransmitters[i].X = tPos.X
			distance := math.Abs(wall.A.Y - tPos.Y)
			if tPos.Y < wall.A.Y {
				mirroredTransmitters[i].Y = wall.A.Y + distance
			} else if tPos.Y > wall.A.Y {
				mirroredTransmitters[i].Y = wall.A.Y - distance
			} else {
				mirroredTransmitters[i].Y = wall.A.Y
			}
			continue
		}
	}
	return &RayTracing{
		Step:             step,
		TransmitterPos:   tPos,
		TransmitterPower: tPower,
		TransmitterFreq:  tFreq,
		WaveLength:       299792458 / (tFreq * math.Pow(10, 9)),
		ReflectionFactor: rFactor,
		Walls: wallPos,
		PowerMap: powerMap,
		Matrix: matrix,
		MirroredTransmitters: mirroredTransmitters,
	}
}

func main() {
	matrixDimensions := Point{X:20, Y:30}
	transmitterPos := Point{X:5, Y:5}
	transmitterPower := 10.0 // mW
	transmitterFreq := 2.4   // GHz
	reflectionFactor := 0.8
	walls := []Vector{{A:Point{X:3,Y:3}, B:Point{X:3,Y:6}}, {A:Point{X:1,Y:3}, B:Point{X:6,Y:3}}}

	raytracing := NewRayTracing(matrixDimensions, transmitterPos, transmitterPower, transmitterFreq, reflectionFactor, walls)
	fmt.Println("Raytracing instance created:", raytracing.MirroredTransmitters)
}
