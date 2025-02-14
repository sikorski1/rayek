package main

import (
	"fmt"
	"image/color"
	"math"
	"math/cmplx"

	"gonum.org/v1/plot"
    "gonum.org/v1/plot/plotter"
    "gonum.org/v1/plot/vg"
    "gonum.org/v1/plot/palette"
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
				H += calculateTransmittance(receiverPos, rt.TransmitterPos, rt.WaveLength, 1.0)
			}
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
		if twoVectors(receiverPos, transmitterPos,wall.A, wall.B) >= 0 {
			return false
		}
	}
	return true
}

func calculateTransmittance(p1, p2 Point, waveLength, reflectionRef float64) complex128 {
	r := calculateDist(p1, p2)
	if r > 0 {
		i := complex(0, 1) 
		H := complex(reflectionRef, 0) * complex(waveLength/(4*math.Pi*r), 0) *
			cmplx.Exp(complex(-2*math.Pi*r/waveLength,0)*i) 
		return H
	} else {
		return 0
	}
}
func calculateDist(p1, p2 Point) float64 {
	dist := math.Sqrt(math.Pow(p1.X  - p2.X, 2) +  math.Pow(p1.Y - p2.Y, 2))
	return dist
}
func twoVectors(A, B, C, D Point) int8 {
	result := (((C.X - A.X)*(B.Y - A.Y) - (B.X - A.X)*(C.Y - A.Y)) * ((D.X - A.X) * (B.Y - A.Y) - (B.X - A.X) * (D.Y - A.Y)))  
        if result > 0{
            return -1
		} else {
            result2 := (((A.X - C.X)*(D.Y - C.Y) - (D.X - C.X)*(A.Y - C.Y)) * ((B.X - C.X) * (D.Y - C.Y) - (D.X - C.X) * (B.Y - C.Y)))
            if result2 > 0 {
                return -1
			} else if result < 0 && result2 < 0 {
                return 1
			} else if result == 0 && result2 < 0 {
                return 0
			} else if result < 0 &&  result2 == 0 {
                return 0
			} else if A.X < C.X && A.X < D.X && B.X < C.X && B.X < D.X {
                return -1
			} else if A.Y < C.Y && A.Y < D.Y && B.Y < C.Y && B.Y < D.Y {
                return -1
			} else if A.X > C.X && A.X > D.X && B.X > C.X && B.X > D.X {
                return -1
			} else if A.Y > C.Y && A.Y > D.Y && B.Y > C.Y && B.Y > D.Y {
                return -1
			}  else {
				return 0
			}
		}
}

// Interpolacja kolorów między niebieskim a czerwonym
func interpolateColor(value, min, max float64) color.RGBA {
    norm := (value - min) / (max - min) // Normalizacja do zakresu [0,1]
    r := uint8(255 * norm)
    g := uint8(0)
    b := uint8(255 * (1 - norm))
    return color.RGBA{R: r, G: g, B: b, A: 255}
}

func (rt *RayTracing) PlotVisualization(filename string) error {
    p := plot.New()
    p.Title.Text = "RayTracing Heatmap"
    p.X.Label.Text = "X"
    p.Y.Label.Text = "Y"

    // Określenie zakresu wykresu
    xMin, xMax := 0.0, rt.Matrix[0][len(rt.Matrix[0])-1].X
    yMin, yMax := 0.0, rt.Matrix[len(rt.Matrix)-1][0].Y

    // Normalizacja osi
    xRange := xMax - xMin
    yRange := yMax - yMin
    if xRange > yRange {
        diff := (xRange - yRange) / 2
        yMin -= diff
        yMax += diff
    } else {
        diff := (yRange - xRange) / 2
        xMin -= diff
        xMax += diff
    }
    p.X.Min, p.X.Max = xMin, xMax
    p.Y.Min, p.Y.Max = yMin, yMax

    // Znalezienie minimalnej i maksymalnej wartości mocy
    minPower, maxPower := math.Inf(1), math.Inf(-1)
    for i := range rt.PowerMap {
        for j := range rt.PowerMap[i] {
            if rt.PowerMap[i][j] < minPower {
                minPower = rt.PowerMap[i][j]
            }
            if rt.PowerMap[i][j] > maxPower {
                maxPower = rt.PowerMap[i][j]
            }
        }
    }

    // Rysowanie mapy cieplnej
    heatmap := make(plotter.XYs, 0)
    colors := make([]color.RGBA, 0)
    for i := range rt.Matrix {
        for j := range rt.Matrix[i] {
            heatmap = append(heatmap, plotter.XY{
                X: rt.Matrix[i][j].X,
                Y: rt.Matrix[i][j].Y,
            })
            colors = append(colors, interpolateColor(rt.PowerMap[i][j], minPower, maxPower))
        }
    }
    scatter, err := plotter.NewScatter(heatmap)
    if err != nil {
        return err
    }
    scatter.GlyphStyle.Radius = vg.Points(2)
    scatter.GlyphStyle.Color = color.Transparent // Ustawienie koloru przez własną mapę
    scatter.ColorMap = plotter.ColorMapFunc(func(x, y float64) color.Color {
        for k, pt := range heatmap {
            if pt.X == x && pt.Y == y {
                return colors[k]
            }
        }
        return color.Black
    })
    p.Add(scatter)

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
        l.Color = color.Black
        l.Width = vg.Points(2)
        p.Add(l)
    }

    // Rysowanie transmitera
    transmitter := plotter.XYs{{X: rt.TransmitterPos.X, Y: rt.TransmitterPos.Y}}
    tScatter, err := plotter.NewScatter(transmitter)
    if err != nil {
        return err
    }
    tScatter.Color = color.RGBA{R: 255, G: 0, B: 0, A: 255}
    tScatter.Radius = vg.Points(5)
    p.Add(tScatter)

    return p.Save(10*vg.Inch, 10*vg.Inch, filename)
}

func main() {
	matrixDimensions := Point{X:20, Y:30}
	transmitterPos := Point{X:17, Y:7}
	transmitterPower := 10.0 // mW
	transmitterFreq := 2.4   // GHz
	reflectionFactor := 0.8
	walls := []Vector{{A:Point{X:0,Y:3}, B:Point{X:3,Y:6}}, {A:Point{X:1,Y:3}, B:Point{X:6,Y:3}}, {A:Point{X:6,Y:10}, B:Point{X:12,Y:12}}, {A:Point{X:15,Y:15}, B:Point{X:15,Y:3}}}

	raytracing := NewRayTracing(matrixDimensions, transmitterPos, transmitterPower, transmitterFreq, reflectionFactor, walls)
	fmt.Printf("%v", raytracing.MirroredTransmitters)
	raytracing.calculateRayTracing()
	err := raytracing.PlotVisualization("raytracing.png")
    if err != nil {
        fmt.Printf("Error creating visualization: %v\n", err)
        return
    }
    fmt.Println("Visualization saved to raytracing.png")
}
