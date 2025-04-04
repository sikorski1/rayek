package raylaunching

import (
	."backendGo/types"
	"math"
)

type RayLaunching3DConfig struct {
	numOfRaysAzim, numOfRaysElev, numOfInteractions, wallMapNumber, cornerMapNumber int
	sizeX, sizeY, sizeZ, step, reflFactor, stationPower, minimalRayPower, stationFreq, waveLength float64
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
	for i := 0; i < rl.Config.numOfRaysAzim; i++ { // loop over horizontal dim
		theta := (2*math.Pi)/float64(rl.Config.numOfRaysAzim)*float64(i) // from -π to π
		for j := 0; j < rl.Config.numOfRaysElev; j++ { // loop over vertical dim
			var phi,dx,dy,dz float64

			// spherical coordinates
			if rl.TransmitterPos.Z == 0 {
				// half sphere – from -π/2 to π/2
				phi = (-math.Pi/2) + (math.Pi / float64(rl.Config.numOfRaysElev-1)) *  float64(j) // from -π/2 to π/2
				dx = math.Cos(theta) * math.Cos(phi) * rl.Config.step 
				dy = math.Sin(theta) * math.Cos(phi) * rl.Config.step
				dz = math.Sin(phi) * rl.Config.step
			} else {
				//full sphere – from 0 to π
				phi = math.Pi * float64(j) / float64(rl.Config.numOfRaysElev-1) // from 0 to π
				dx = math.Sin(phi) * math.Cos(theta) * rl.Config.step
				dy = math.Sin(phi) * math.Sin(theta) * rl.Config.step
				dz = math.Cos(phi) * rl.Config.step
			}

			/* getting past to next step,
			 omitting calculations for transmitter */
			x := rl.TransmitterPos.X + dx
			y := rl.TransmitterPos.Y + dy
			z := rl.TransmitterPos.Z + dz

			// initial counters
			currInteractions := 0
			currPower := 0.0
			currWallIndex := 0
			currStartLengthPos := Point3D{X:rl.TransmitterPos.X, Y:rl.TransmitterPos.Y, Z:rl.TransmitterPos.Z}
			currRayLength := 0.0
			// main loop
			for (x >= 0 && x <= rl.Config.sizeX) && (y >= 0 && y <= rl.Config.sizeY) && (z >= 0 && z <= rl.Config.sizeZ) && currInteractions < rl.Config.numOfInteractions && currPower >= rl.Config.minimalRayPower {
				xIdx := int(math.Round(x/rl.Config.step))
				yIdx := int(math.Round(y/rl.Config.step))
				zIdx := int(math.Round(z/rl.Config.step))
				index := int(rl.PowerMap[zIdx][yIdx][xIdx])

				if index == rl.Config.cornerMapNumber {
					continue
				}
				// check if there is wall and if its diffrent from previous one
				if index >= rl.Config.wallMapNumber && index != currWallIndex + rl.Config.wallMapNumber {
					currWallIndex = index - rl.Config.wallMapNumber
					nx, ny, nz := rl.WallNormals[currWallIndex].Nx, rl.WallNormals[currWallIndex].Ny, rl.WallNormals[currWallIndex].Nz

					dot := 2 * (dx*nx + dy*ny + dz*nz)
					dx = dx - dot*nx
					dy = dy - dot*ny
					dz = dz - dot*nz
					currInteractions++
					currRayLength += rl.CalculateDistance(currStartLengthPos, Point3D{X: x, Y: y, Z: z})
					currStartLengthPos = Point3D{X: x, Y: y, Z: z}
				} else {
					currRayLength = rl.CalculateDistance(currStartLengthPos, Point3D{X: x, Y: y, Z: z})
				}
			}
		}
	}
}

func (rl *RayLaunching3D) CalculateDistance(p1, p2 Point3D) float64 {
	dist := math.Sqrt(math.Pow(p1.X-p2.X,2)+math.Pow(p1.Y-p2.Y,2)+math.Pow(p1.Z-p2.Z,2))
	return dist
}