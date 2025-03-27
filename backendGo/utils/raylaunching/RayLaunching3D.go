package raylaunching

import (
	."backendGo/types"
)

type RayLaunching3DConfig struct {
	sizeX, sizeY, numOfRaysAzim, numOfRaysElev, numOfInteractions int
	step, reflFactor, stationPower, stationHeight, stationFreq, waveLength float64
}

type RayLaunching3D struct {
	PowerMap [][][]float64
	WallNormals []Normal3D
	Config []RayLaunching3DConfig
}

func NewRayLaunching3D(matrix [][][]float64, wallNormals []Normal3D, config []RayLaunching3DConfig) *RayLaunching3D {
	return &RayLaunching3D{
		PowerMap: matrix,
		WallNormals: wallNormals,
		Config: config,
	}
}