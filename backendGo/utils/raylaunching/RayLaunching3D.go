package raylaunching

import (
	."backendGo/types"
)

type RayLaunching3DConfig struct {
	sizeX, sizeY, numOfRaysAzim, numOfRaysElev, numOfInteractions int
	step, reflFactor, stationPower, stationHeight, stationFreq, waveLength float64
}

func RayLaunching3D(matrix [][][]float64, wallNormals []Normal3D, config []RayLaunching3DConfig) [][][]float64 {
	return matrix
}