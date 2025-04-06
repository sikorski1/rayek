package raylaunching

import (
	."backendGo/types"
	"math"
	"math/cmplx"
)

type RayLaunching3DConfig struct {
	NumOfRaysAzim, NumOfRaysElev, NumOfInteractions, WallMapNumber, CornerMapNumber int
	SizeX, SizeY, SizeZ, Step, ReflFactor, TransmitterPower, MinimalRayPower, TransmitterFreq, WaveLength float64
	TransmitterPos Point3D
}

type RayLaunching3D struct {
	TransmitterPos TransmitterPos3D
	PowerMap [][][]float64
	WallNormals []Normal3D
	Config RayLaunching3DConfig
}

func NewRayLaunching3D(matrix [][][]float64, wallNormals []Normal3D, config RayLaunching3DConfig) *RayLaunching3D {
	return &RayLaunching3D{
		PowerMap: matrix,
		WallNormals: wallNormals,
		Config: config,
	}
}

func (rl *RayLaunching3D) CalculateRayLaunching3D() {
	for z := 0; z < int(rl.Config.TransmitterPos.Z); z++ {
		rl.PowerMap[z][int(rl.Config.TransmitterPos.Y)][int(rl.Config.TransmitterPos.X)] = 0
	}
	
	for i := 0; i < rl.Config.NumOfRaysAzim; i++ { // loop over horizontal dim
		theta := (2*math.Pi)/float64(rl.Config.NumOfRaysAzim)*float64(i) // from -π to π
		for j := 0; j < rl.Config.NumOfRaysElev; j++ { // loop over vertical dim
			var phi,dx,dy,dz float64

			// spherical coordinates
			if rl.Config.TransmitterPos.Z == 0 {
				// half sphere – from 0 to π/2
				phi = ((math.Pi/2) / float64(rl.Config.NumOfRaysElev-1)) *  float64(j) // from 0 to π/2
				dx = math.Cos(theta) * math.Cos(phi) * rl.Config.Step 
				dy = math.Sin(theta) * math.Cos(phi) * rl.Config.Step
				dz = math.Sin(phi) * rl.Config.Step
			} else {
				//full sphere – from 0 to π
				phi = 2*math.Pi * float64(j) / float64(rl.Config.NumOfRaysElev-1) // from 0 to π
				dx = math.Cos(theta) * math.Cos(phi) * rl.Config.Step 
				dy = math.Sin(theta) * math.Cos(phi) * rl.Config.Step
				dz = math.Sin(phi) * rl.Config.Step
				
			}
			dx, dy, dz = math.Round(dx*1e15)/1e15, math.Round(dy*1e15)/1e15, math.Round(dz*1e15)/1e15

			/* getting past to next step,
			 omitting calculations for transmitter */

			x := rl.Config.TransmitterPos.X + dx
			y := rl.Config.TransmitterPos.Y + dy
			z := rl.Config.TransmitterPos.Z + dz

			// initial counters
			currInteractions := 0
			currPower := 0.0
			currWallIndex := 0
			currStartLengthPos := Point3D{X:rl.Config.TransmitterPos.X, Y:rl.Config.TransmitterPos.Y, Z:rl.Config.TransmitterPos.Z}
			currRayLength := 0.0
			currSumRayLength := 0.0

			// main loop
			for (x >= 0 && x <= rl.Config.SizeX) && (y >= 0 && y <= rl.Config.SizeY) && (z >= 0 && z <= rl.Config.SizeZ) && currInteractions < rl.Config.NumOfInteractions && currPower >= rl.Config.MinimalRayPower {
				xIdx := int(math.Round(x/rl.Config.Step))
				yIdx := int(math.Round(y/rl.Config.Step))
				zIdx := int(math.Round(z/rl.Config.Step))
				index := int(rl.PowerMap[zIdx][yIdx][xIdx])
				if index == rl.Config.CornerMapNumber {
					break
				}
				// check if there is wall and if its diffrent from previous one
				if index >= rl.Config.WallMapNumber && index != currWallIndex + rl.Config.WallMapNumber {
					currWallIndex = index - rl.Config.WallMapNumber

					//get wall normal
					nx, ny, nz := rl.WallNormals[currWallIndex].Nx, rl.WallNormals[currWallIndex].Ny, rl.WallNormals[currWallIndex].Nz
					dot := 2 * (dx*nx + dy*ny + dz*nz)

					// calculate new direction
					dx = dx - dot*nx
					dy = dy - dot*ny
					dz = dz - dot*nz
					currInteractions++

					// sum distance and set new start position
					currSumRayLength += calculateDistance(currStartLengthPos, Point3D{X: x, Y: y, Z: z})
					currStartLengthPos = Point3D{X: x, Y: y, Z: z}
				} else {

					// calculate distance and transmittance
					currRayLength = calculateDistance(currStartLengthPos, Point3D{X: x, Y: y, Z: z}) + currSumRayLength
					H := calculateTransmittance(currRayLength, rl.Config.WaveLength, math.Pow(rl.Config.ReflFactor, float64(currInteractions)))
					currPower = 10*math.Log10(rl.Config.TransmitterPower) + 20*math.Log10(cmplx.Abs(H))
					// update power map if power is higher than previous one
					if rl.PowerMap[zIdx][yIdx][xIdx] == -150 || rl.PowerMap[zIdx][yIdx][xIdx] < currPower {
						rl.PowerMap[zIdx][yIdx][xIdx] = currPower
					} 
				}
				// update position
				x += dx
				y += dy
				z += dz
			}
		}
	}
}

func calculateDistance(p1, p2 Point3D) float64 {
	dist := math.Sqrt(math.Pow(p1.X-p2.X,2)+math.Pow(p1.Y-p2.Y,2)+math.Pow(p1.Z-p2.Z,2))
	return dist
}

func calculateTransmittance(r, waveLength, reflectionRef float64) complex128 {
	if r > 0 {
		H := complex(reflectionRef, 0) * complex(waveLength/(4*math.Pi*r), 0) *
			cmplx.Exp(complex(0,-2*math.Pi*r/waveLength)) 
		return H
	} else {
		return 0
	}
}