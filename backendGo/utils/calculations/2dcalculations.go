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
func CalculateTransmittanceWithLength(r, waveLength, reflectionRef float64) complex128 {
	if r > 0 {
		H := complex(reflectionRef, 0) * complex(waveLength/(4*math.Pi*r), 0) *
			cmplx.Exp(complex(0,-2*math.Pi*r/waveLength)) 
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
func GenerateHeatmap(powerMap [][]float64) *image.RGBA {
    height := len(powerMap)
    width := len(powerMap[0])
    
    // Fixed range values
    minVal := -150.0  
    maxVal := 0.0
    
    img := image.NewRGBA(image.Rect(0, 0, width, height))
    for y := 0; y < height; y++ {
        yFlipped := height - 1 - y
        for x := 0; x < width; x++ {
            val := powerMap[y][x]
            
            if val >= 10000 {
                // Corner points are pink
                img.Set(x, yFlipped, color.RGBA{255, 0, 255, 255}) // Bright pink
            } else if val == 5000 {
				  img.Set(x, yFlipped, color.RGBA{192, 192, 192, 255})
			} else if val >= 1000 {
                // Regular walls are black
                img.Set(x, yFlipped, color.RGBA{0, 0, 0, 255})
            } else if val == 0 {
                // Empty space is white
                img.Set(x, yFlipped, color.RGBA{255, 255, 255, 255})
            } else {
                // Clamp values to range
                if val < minVal {
                    val = minVal
                }
                if val > maxVal {
                    val = maxVal
                }
                
                // Normalize to 0-1 range
                normalizedVal := (val - minVal) / (maxVal - minVal)
                r, g, b := getHeatmapColor(normalizedVal)
                img.Set(x, yFlipped, color.RGBA{r, g, b, 255})
            }
        }
    }

    return img
}

func getHeatmapColor(value float64) (uint8, uint8, uint8) {
    switch {
    case value < 0.25:
        return 0, uint8(255 * value * 4), 255
    case value < 0.5:
        return 0, 255, uint8(255 * (2 - value*4))
    case value < 0.75:
        return uint8(255 * (value*4 - 2)), 255, 0
    default:
        return 255, uint8(255 * (4 - value*4)), 0
    }
}
