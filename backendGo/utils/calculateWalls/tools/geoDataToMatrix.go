package tools

import (
	. "backendGo/types"
	"math"
)
func GeoToMatrixIndex(lat, lon, latMin, latMax, lonMin, lonMax float64, size int) (int, int) {
	if lat < latMin || lat > latMax || lon < lonMin || lon > lonMax {
		return -1, -1
	}
	y := (lat - latMin) / (latMax - latMin) * float64(size-1)
	x := (lon - lonMin) / (lonMax - lonMin) * float64(size-1)
	i := int(math.Round(x))
	j := int(math.Round(y))
	return i, j
}

func DrawLine(matrix [][][]float64, x1, y1, z1, x2, y2, z2, heightLevels int) {
	dx := x2 - x1
	dy := y2 - y1
	if z1 >= heightLevels {
		z1 = heightLevels - 1
	}
	if z2 >= heightLevels {
		z2 = heightLevels - 1
	}
	if x1 == x2 {
		if y1 > y2 {
			y1, y2 = y2, y1 
		}
		for y := y1; y <= y2; y++ {
			for z := 0; z <= z1; z++ {
				matrix[z][y][x1] = 1
			}
		}
	} else if y1 == y2 {
		if x1 > x2 {
			x1, x2 = x2, x1
		}
		for x := x1; x <= x2; x++ {
			for z := 0; z <= z1; z++ {
				matrix[z][y1][x] = 1
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
				for z := 0; z <= z1; z++ {
					matrix[z][yIdx][prevXIdx] = 1
				}
			}
			if prevXIdx > xIdx && prevYIdx < yIdx  || prevXIdx > xIdx && prevYIdx > yIdx  {
				for z := 0; z <= z1; z++ {
					matrix[z][prevYIdx][xIdx] = 1
				}
			} // walls continuity
			for z := 0; z <= z1; z++ {
				matrix[z][yIdx][xIdx] = 1
			}
			prevXIdx = xIdx
			prevYIdx = yIdx
		} 
	}
}

func GenerateBuildingMatrix(buildings []Building, latMin, latMax, lonMin, lonMax float64, size, heightLevels int) [][][]float64{
	matrix := make([][][]float64, heightLevels)
	for z := range matrix {
		matrix[z] = make([][]float64, size)
		for y := range matrix[z] {
			matrix[z][y] = make([]float64, size)
		}
	}

	for _, building := range buildings {
		for _, wall := range building.Walls {
			i1, j1 := GeoToMatrixIndex(wall.Start.Y, wall.Start.X, latMin, latMax, lonMin, lonMax, size)
			i2, j2 := GeoToMatrixIndex(wall.End.Y, wall.End.X, latMin, latMax, lonMin, lonMax, size)
			if i1 == -1 || i2 == -1 || j1 == -1 || j2 == -1 {
				continue
			}
			z1 := int(math.Round(wall.Start.Z))
			z2 := int(math.Round(wall.End.Z))
			DrawLine(matrix, i1, j1, z1, i2, j2, z2, heightLevels)
		}
	}
	return matrix
}