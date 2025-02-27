package main

import (
	. "backendGo/types"
	. "backendGo/utils/calculations"
	"fmt"
	"image/color"
	"image/png"
	"math"
	"math/cmplx"
	"os"
	"time"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
)
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
	step := 0.01
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
		a1 := (wall.B.Y - wall.A.Y) / (wall.B.X -  wall.A.X)
		b1 := wall.A.Y  -  a1 * wall.A.X

		a2 := -1/a1
		b2 := tPos.Y - a2 * tPos.X

		x := (b2-b1)/(a1-a2)
		y := a1*x + b1

		mirroredTransmitters[i].X = math.Round((2*x - tPos.X) * 10)/10
		mirroredTransmitters[i].Y = math.Round((2*y - tPos.Y) * 10)/10

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

func (rt *RayTracing) calculateRayTracing() {
	for i := range(len(rt.Matrix)) {
		for j := range(len(rt.Matrix[0])) {
			H := complex(0,0)
			receiverPos := rt.Matrix[i][j]
			if checkLineOfSight(rt.TransmitterPos, receiverPos, rt.Walls) {
				H += CalculateTransmittance(receiverPos, rt.TransmitterPos, rt.WaveLength, 1.0)
			}
			H += calculateSingleWallReflection(rt.MirroredTransmitters, rt.TransmitterPos, receiverPos, rt.Walls, rt.WaveLength, rt.ReflectionFactor)
			if H == 0 {
				rt.PowerMap[i][j] = -150
			} else {
				rt.PowerMap[i][j] = 10*math.Log10(rt.TransmitterPower) + 20*math.Log10(cmplx.Abs(H))
			}
		}
	}
}

func checkLineOfSight(transmitterPos, receiverPos Point, walls []Vector) bool {
	for _, wall := range(walls) {
		if TwoVectors(receiverPos, transmitterPos,wall.A, wall.B) >= 0 {
			return false
		}
	}
	return true
}

func calculateSingleWallReflection(mirroredTransmitters []Point, transmitterPos, receiverPos Point, walls []Vector, waveLength, reflectionFactor float64) complex128 {
	H := complex(0,0)
	for i, wall := range walls {
		if TwoVectors(receiverPos, mirroredTransmitters[i], wall.A, wall.B) <= 0 {
			continue
		} 
		reflectionPoint := CalculateCrossPoint(receiverPos, mirroredTransmitters[i], wall.A, wall.B)
		collision := false
		for j := range(len(walls)-1) {
			index := (i+j+1) % len(walls)
			if TwoVectors(transmitterPos, reflectionPoint, walls[index].A, walls[index].B) >= 0 {
				collision = true
				break
			}
			if TwoVectors(reflectionPoint, receiverPos, walls[index].A, walls[index].B) >= 0{
				collision = true
				break
			}
		}
		if !collision {
			r := CalculateDist(transmitterPos,  reflectionPoint) + CalculateDist(reflectionPoint, receiverPos)
			H += CalculateTransmittanceWithLength(r, waveLength, reflectionFactor)
		}
	}
	return H
}
// Interpolacja kolorów między niebieskim a czerwonym

func (rt *RayTracing) PlotVisualization(filename string) error {
    p := plot.New()
    p.Title.Text = "RayTracing Visualization"
    p.X.Label.Text = "X"
    p.Y.Label.Text = "Y"
    // Ustalanie jednakowej skali dla osi X i Y
	xMin, xMax := 0.0, rt.Matrix[0][len(rt.Matrix[0])-1].X
	yMin, yMax := 0.0, rt.Matrix[len(rt.Matrix)-1][0].Y
	xRange := xMax - xMin
	yRange := yMax - yMin
	// Wybieramy większy zakres i dopasowujemy mniejszy
	if xRange > yRange {
		diff := (xRange - yRange) / 2
		yMin -= diff
		yMax += diff
	} else {
		diff := (yRange - xRange) / 2
		xMin -= diff
		xMax += diff
	}
	// Ustawienie zakresu wykresu po korekcie
	p.X.Min, p.X.Max = xMin, xMax
	p.Y.Min, p.Y.Max = yMin, yMax
    // Rysowanie ścian
    for _, wall := range rt.Walls {
        line := plotter.XYs{
            {X: wall.A.X, Y: wall.A.Y},
            {X: wall.B.X, Y: wall.B.Y},
        }
        l, err := plotter.NewLine(line)
        if err != nil {
            return err
        }
        l.Color = color.RGBA{R: 0, G: 0, B: 0, A: 255}
        l.Width = vg.Points(2)
        p.Add(l)
    }
    // Rysowanie transmitera
    transmitter := plotter.XYs{
        {X: rt.TransmitterPos.X, Y: rt.TransmitterPos.Y},
    }
    transmitterScatter, err := plotter.NewScatter(transmitter)
    if err != nil {
        return err
    }
    transmitterScatter.Color = color.RGBA{R: 255, G: 0, B: 0, A: 255}
    transmitterScatter.Radius = vg.Points(5)
    p.Add(transmitterScatter)
    // Rysowanie odbitych transmiterów
    for _, mt := range rt.MirroredTransmitters {
        mirroredTransmitter := plotter.XYs{
            {X: mt.X, Y: mt.Y},
        }
        mtScatter, err := plotter.NewScatter(mirroredTransmitter)
        if err != nil {
            return err
        }
        mtScatter.Color = color.RGBA{R: 0, G: 0, B: 255, A: 255}
        mtScatter.Radius = vg.Points(3)
        p.Add(mtScatter)
    }
    // Dodanie legendy
    p.Legend.Add("Transmitter", transmitterScatter)
    
    // Tworzymy linię dla legendy
    legendLine, err := plotter.NewLine(plotter.XYs{{X: 0, Y: 0}, {X: 1, Y: 1}})
    if err != nil {
        return err
    }
    p.Legend.Add("Walls", legendLine)
    
    // Tworzymy punkt dla legendy
    legendPoint, err := plotter.NewScatter(plotter.XYs{{X: 0, Y: 0}})
    if err != nil {
        return err
    }
    p.Legend.Add("Mirrored Transmitters", legendPoint)
    
    p.Legend.Top = true
    p.Legend.Left = true
    // Zapisanie wykresu do pliku
    return p.Save(10*vg.Inch, 10*vg.Inch, filename)
}


func main() {
	start := time.Now()
	matrixDimensions := Point{X:40, Y:40}
	transmitterPos := Point{X:16, Y:20}
	transmitterPower := 5.0 // mW
	transmitterFreq := 2.4   // GHz
	reflectionFactor := 0.8
	walls := []Vector{
		
		{A: Point{X: 2, Y: 2}, B: Point{X: 8, Y: 2}},
		{A: Point{X: 8, Y: 2}, B: Point{X: 8, Y: 8}},
		{A: Point{X: 8, Y: 8}, B: Point{X: 2, Y: 8}},
		{A: Point{X: 2, Y: 8}, B: Point{X: 2, Y: 2}},

		{A: Point{X: 12, Y: 4}, B: Point{X: 18, Y: 4}},
		{A: Point{X: 18, Y: 4}, B: Point{X: 18, Y: 10}},
		{A: Point{X: 18, Y: 10}, B: Point{X: 12, Y: 10}},
		{A: Point{X: 12, Y: 10}, B: Point{X: 12, Y: 4}},


		{A: Point{X: 22, Y: 6}, B: Point{X: 30, Y: 6}},
		{A: Point{X: 30, Y: 6}, B: Point{X: 30, Y: 14}},
		{A: Point{X: 30, Y: 14}, B: Point{X: 22, Y: 14}},
		{A: Point{X: 22, Y: 14}, B: Point{X: 22, Y: 6}},


		{A: Point{X: 4, Y: 12}, B: Point{X: 10, Y: 12}},
		{A: Point{X: 10, Y: 12}, B: Point{X: 10, Y: 18}},
		{A: Point{X: 10, Y: 18}, B: Point{X: 4, Y: 18}},
		{A: Point{X: 4, Y: 18}, B: Point{X: 4, Y: 12}},




		{A: Point{X: 26, Y: 18}, B: Point{X: 34, Y: 18}},
		{A: Point{X: 34, Y: 18}, B: Point{X: 34, Y: 26}},
		{A: Point{X: 34, Y: 26}, B: Point{X: 26, Y: 26}},
		{A: Point{X: 26, Y: 26}, B: Point{X: 26, Y: 18}},

	
		{A: Point{X: 6, Y: 22}, B: Point{X: 12, Y: 22}},
		{A: Point{X: 12, Y: 22}, B: Point{X: 12, Y: 30}},
		{A: Point{X: 12, Y: 30}, B: Point{X: 6, Y: 30}},
		{A: Point{X: 6, Y: 30}, B: Point{X: 6, Y: 22}},

		{A: Point{X: 16, Y: 26}, B: Point{X: 24, Y: 26}},
		{A: Point{X: 24, Y: 26}, B: Point{X: 24, Y: 34}},
		{A: Point{X: 24, Y: 34}, B: Point{X: 16, Y: 34}},
		{A: Point{X: 16, Y: 34}, B: Point{X: 16, Y: 26}},

		{A: Point{X: 28, Y: 30}, B: Point{X: 36, Y: 30}},
		{A: Point{X: 36, Y: 30}, B: Point{X: 36, Y: 38}},
		{A: Point{X: 36, Y: 38}, B: Point{X: 28, Y: 38}},
		{A: Point{X: 28, Y: 38}, B: Point{X: 28, Y: 30}},


		{A: Point{X: 10, Y: 34}, B: Point{X: 18, Y: 34}},
		{A: Point{X: 18, Y: 34}, B: Point{X: 18, Y: 40}},
		{A: Point{X: 18, Y: 40}, B: Point{X: 10, Y: 40}},
		{A: Point{X: 10, Y: 40}, B: Point{X: 10, Y: 34}},
	}

	raytracing := NewRayTracing(matrixDimensions, transmitterPos, transmitterPower, transmitterFreq, reflectionFactor, walls)
	fmt.Printf("%v \n", raytracing.MirroredTransmitters)
	fmt.Printf("%v \n", raytracing.Matrix[285][211])
	
	raytracing.calculateRayTracing()
	stop := time.Since(start)
	fmt.Printf("Computation time: %v \n", stop)
	err := raytracing.PlotVisualization("raytracing.png")
    if err != nil {
        fmt.Printf("Error creating visualization: %v\n", err)
        return
    }
    fmt.Println("Visualization saved to raytracing.png")
	heatmap := GenerateHeatmap(raytracing.PowerMap)
    f, _ := os.Create("heatmap.png")
    defer f.Close()
    png.Encode(f, heatmap)
}
