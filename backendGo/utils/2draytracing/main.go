package main

import (
	"os"
	"fmt"
	"math"
	"math/cmplx"
	"time"
	"image"
    "image/color"
    "image/png"
	"gonum.org/v1/plot"
    "gonum.org/v1/plot/plotter"
    "gonum.org/v1/plot/vg"
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
				H += calculateTransmittance(receiverPos, rt.TransmitterPos, rt.WaveLength, 1.0)
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
		if twoVectors(receiverPos, transmitterPos,wall.A, wall.B) >= 0 {
			return false
		}
	}
	return true
}

func calculateSingleWallReflection(mirroredTransmitters []Point, transmitterPos, receiverPos Point, walls []Vector, waveLength, reflectionFactor float64) complex128 {
	H := complex(0,0)
	for i, wall := range walls {
		if twoVectors(receiverPos, mirroredTransmitters[i], wall.A, wall.B) <= 0 {
			continue
		} 
		reflectionPoint := calculateCrossPoint(receiverPos, mirroredTransmitters[i], wall.A, wall.B)
		collision := false
		for j := range(len(walls)-1) {
			index := (i+j+1) % len(walls)
			if twoVectors(transmitterPos, reflectionPoint, walls[index].A, walls[index].B) >= 0 {
				collision = true
				break
			}
			if twoVectors(reflectionPoint, receiverPos, walls[index].A, walls[index].B) >= 0{
				collision = true
				break
			}
		}
		if !collision {
			H += calculateTransmittance(receiverPos, mirroredTransmitters[i], waveLength, reflectionFactor)
		}
	}
	return H
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
func calculateCrossPoint (A, B, C, D Point) Point {
	if A.X == B.X {
		x := A.X
		a2 := (D.Y - C.Y)/(D.X - C.X)
		b2 := C.Y - a2 * C.X
		y := a2 * x + b2
		return Point{X:x, Y:y}
	} else if C.X == D.X{
		x := C.X
		a1 := (B.Y-A.Y)/(B.X-A.X)
		b1 := A.Y - a1*A.X
		y := a1 * x + b1
		return Point{X:x, Y:y}
	} else {
		a1 := (B.Y-A.Y)/(B.X-A.X)
		b1 := A.Y - a1*A.X
		a2 := (D.Y - C.Y)/(D.X - C.X)
		b2 := C.Y - a2 * C.X
		x := (b2-b1)/(a1-a2)
		y := a1 * x + b1
		return Point{X:x, Y:y}
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
// GenerateHeatmap tworzy mapę ciepła z tablicy 2D wartości
func GenerateHeatmap(powerMap [][]float64) *image.RGBA {
    height := len(powerMap)
    width := len(powerMap[0])
    
    // Znajdź min i max wartości do normalizacji
    minVal, maxVal := math.Inf(1), math.Inf(-1)
    for y := 0; y < height; y++ {
        for x := 0; x < width; x++ {
            val := powerMap[y][x]
            if val < minVal {
                minVal = val
            }
            if val > maxVal {
                maxVal = val
            }
        }
    }
    
    // Utwórz nowy obraz
    img := image.NewRGBA(image.Rect(0, 0, width, height))
    
    // Wypełnij obraz kolorami
    for y := 0; y < height; y++ {
        // Odwracamy współrzędną Y
        yFlipped := height - 1 - y
        
        for x := 0; x < width; x++ {
            // Normalizuj wartość do zakresu 0-1
            normalizedVal := (powerMap[y][x] - minVal) / (maxVal - minVal)
            
            // Konwertuj na kolor (używając przejścia od niebieskiego przez zielony do czerwonego)
            r, g, b := getHeatmapColor(normalizedVal)
            // Używamy odwróconej współrzędnej Y
            img.Set(x, yFlipped, color.RGBA{r, g, b, 255})
        }
    }
    
    return img
}

// getHeatmapColor zwraca kolor RGB dla znormalizowanej wartości (0-1)
func getHeatmapColor(value float64) (uint8, uint8, uint8) {
    // Przejście kolorów: niebieski -> cyjan -> zielony -> żółty -> czerwony
    switch {
    case value < 0.25:
        // niebieski do cyjan
        return 0, uint8(255 * value * 4), 255
    case value < 0.5:
        // cyjan do zielony
        return 0, 255, uint8(255 * (2 - value*4))
    case value < 0.75:
        // zielony do żółty
        return uint8(255 * (value*4 - 2)), 255, 0
    default:
        // żółty do czerwony
        return 255, uint8(255 * (4 - value*4)), 0
    }
}

func main() {
	start := time.Now()
	matrixDimensions := Point{X:40, Y:40}
	transmitterPos := Point{X:18, Y:5}
	transmitterPower := 5.0 // mW
	transmitterFreq := 3.4   // GHz
	reflectionFactor := 0.7
	walls := []Vector{{A:Point{X:0,Y:3}, B:Point{X:3,Y:6}}, {A:Point{X:1,Y:3}, B:Point{X:6,Y:3}}, {A:Point{X:6,Y:10}, B:Point{X:12,Y:12}},{A:Point{X:25,Y:10}, B:Point{X:25,Y:30}},{A:Point{X:5,Y:30}, B:Point{X:10,Y:35}},{A:Point{X:23,Y:36}, B:Point{X:25,Y:39}}}

	raytracing := NewRayTracing(matrixDimensions, transmitterPos, transmitterPower, transmitterFreq, reflectionFactor, walls)
	fmt.Printf("%v \n", raytracing.MirroredTransmitters)
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
    
    // Zapisz do pliku
    f, _ := os.Create("heatmap.png")
    defer f.Close()
    png.Encode(f, heatmap)
}
