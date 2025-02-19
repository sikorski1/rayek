package calculations

import (
	. "backendGo/types"
	"image"
	"image/color"
	"math"
	"math/cmplx"
)

func CalculateTransmittance(p1, p2 Point, waveLength, reflectionRef float64) complex128 {
	r := CalculateDist(p1, p2)
	if r > 0 {
		i := complex(0, 1) 
		H := complex(reflectionRef, 0) * complex(waveLength/(4*math.Pi*r), 0) *
			cmplx.Exp(complex(-2*math.Pi*r/waveLength,0)*i) 
		return H
	} else {
		return 0
	}
}
func CalculateCrossPoint (A, B, C, D Point) Point {
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
func CalculateDist(p1, p2 Point) float64 {
	dist := math.Sqrt(math.Pow(p1.X  - p2.X, 2) +  math.Pow(p1.Y - p2.Y, 2))
	return dist
}
func TwoVectors(A, B, C, D Point) int8 {
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
