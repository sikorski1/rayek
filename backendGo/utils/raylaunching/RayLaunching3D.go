package raylaunching

import (
	."backendGo/types"
	"math"
)

type RayLaunching3DConfig struct {
	numOfRaysAzim, numOfRaysElev, numOfInteractions int
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
		for j := 0; j < rl.Config.numOfRaysElev; j++ { // loop over vertical dim
			theta := (2*math.Pi)/float64(rl.Config.numOfRaysAzim)*float64(i) // from -π to π
			phi := (-math.Pi/2) + (math.Pi / float64(rl.Config.numOfRaysElev-1)) *  float64(j) // from -π/2 to π/2

			// spherical coordinates
			dx := math.Cos(theta) * math.Cos(phi) * rl.Config.step
			dy := math.Sin(theta) * math.Cos(phi) * rl.Config.step
			dz := math.Sin(phi) *  rl.Config.step 
			dx, dy, dz = math.Round(dx*1e15)/1e15, math.Round(dy*1e15)/1e15, math.Round(dz*1e15)/1e15
			
			/* getting past to next step,
			 omitting calculations for transmitter */

			x := rl.TransmitterPos.X + dx
			y := rl.TransmitterPos.Y + dy
			z := rl.TransmitterPos.Z + dz

			// initial counters
			currInteractions := 0
			currPower := 0.0

			// main loop
			for (x >= 0 && x <= rl.Config.sizeX) && (y >= 0 && y <= rl.Config.sizeY) && (z >= 0 && z <= rl.Config.sizeZ) && currInteractions < rl.Config.numOfInteractions && currPower >= rl.Config.minimalRayPower {
				continue
			}
		}
	}
}