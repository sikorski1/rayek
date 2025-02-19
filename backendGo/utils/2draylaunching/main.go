package main
import (
	. "backendGo/types"

	"fmt"

	"math"
	
	"time"
	
)

type RayLaunching struct {
	Step                 float64
	TransmitterPos       Point
	TransmitterPower     float64
	TransmitterFreq      float64
	WaveLength           float64
	ReflectionFactor     float64
	Map                  [][]float64
}

func NewRayLaunching(matrixDimensions Point, tPos Point, tPower float64, tFreq float64, rFactor float64, wallPos []Vector) *RayLaunching {
	step := 0.1
	rows := int(matrixDimensions.Y*(1/step))+1
	cols := int(matrixDimensions.X*(1/step))+1
	Map := make([][]float64, rows)
	for i := range Map {
		Map[i] = make([]float64, cols)
	}
	
	return &RayLaunching{
		Step:             step,
		TransmitterPos:   tPos,
		TransmitterPower: tPower,
		TransmitterFreq:  tFreq,
		WaveLength:       299792458 / (tFreq * math.Pow(10, 9)),
		ReflectionFactor: rFactor,
		Map: Map,
		
	}
}




func main() {
	start := time.Now()
	matrixDimensions := Point{X:15, Y:20}
	transmitterPos := Point{X:5, Y:20}
	transmitterPower := 5.0 // mW
	transmitterFreq := 3.4   // GHz
	reflectionFactor := 0.7
	walls := []Vector{{A:Point{X:6.6,Y:3}, B:Point{X:6.6,Y:12}}, {A:Point{X:4,Y:16}, B:Point{X:13,Y:16}}}
	raylaunching := NewRayLaunching(matrixDimensions, transmitterPos, transmitterPower, transmitterFreq, reflectionFactor, walls)
	fmt.Printf("%v \n", raylaunching.Map[200][150])
	stop := time.Since(start)
	fmt.Printf("Computation time: %v \n", stop)
}