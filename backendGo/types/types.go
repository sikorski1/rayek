package types

type Point struct {
	X, Y float64
}


type Vector struct {
	A, B Point
}

type Normal struct {
	Nx, Ny float64
}
type Normal3D struct {
	Nx, Ny, Nz float64
}
type TransmitterPos3D struct {
	X, Y, Z float64
}

type MapConfig struct {
	LatMin, LatMax, LonMin, LonMax float64
	Size, HeightMaxLevels int
}

type Point3D struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
	Z float64 `json:"z"`
}

type Wall struct {
	Start Point3D `json:"start"`
	End   Point3D `json:"end"`
}

type Building struct {
	Name   string `json:"name"`
	Height float64 `json:"height"`
	Walls  []Wall `json:"walls"`
}

type SingleRay struct {
	Azimuth int `json:"azimuth"`
	Elevation int `json:"elevation"`
}
