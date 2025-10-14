package raylaunching

import (
	. "backendGo/types"
	"fmt"

	"math"
	"math/cmplx"
)

type RayLaunching3DConfig struct {
	NumOfRaysAzim, NumOfRaysElev, NumOfInteractions, WallMapNumber, BuldingInteriorNumber, RoofMapNumber, CornerMapNumber, RoofCornerMapNumber, DiffractionRayNumber int
	SizeX, SizeY, SizeZ, Step, ReflFactor, TransmitterPower, MinimalRayPower, TransmitterFreq, WaveLength                                                            float64
	TransmitterPos                                                                                                                                                   Point3D
	SingleRays                                                                                                                                                       []SingleRay
}

type RayPoint struct {
	X     float64 `json:"x"`
	Y     float64 `json:"y"`
	Z     float64 `json:"z"`
	Power float64 `json:"power"`
}
type RayLaunching3D struct {
	PowerMap       [][][]float64
	WallNormals    []Normal3D
	Config         RayLaunching3DConfig
	RayPaths       [][]RayPoint
	PowerMapLegend map[int]PowerMapLegendEntry
}

type PowerMapLegendEntry struct {
	Ptotal   float64 `json:"total"`
	P0_20    float64 `json:"< 0dbm"`
	P20_40   float64 `json:"< -20dbm"`
	P40_60   float64 `json:"< -40dbm"`
	P60_80   float64 `json:"< -60dbm"`
	P80_100  float64 `json:"< -80dbm"`
	P100_120 float64 `json:"< -100dbm"`
	P120_140 float64 `json:"< -120dbm"`
	P140plus float64 `json:"< -140dbm"`
}

type RayState struct {
	x, y, z                     float64
	dx, dy, dz                  float64
	currInteractions            int
	currPower                   float64
	currWallIndex               int
	currStartLengthPos          Point3D
	currRayLength               float64
	currSumRayLength            float64
	currReflectionFactor        float64
	diffLossLdB                 float64
	targetRayIndex              int
	toDiffractionPointRayLength float64
	diffTheta                   float64
	diffRayIndex                int
}

func NewRayLaunching3D(matrix [][][]float64, wallNormals []Normal3D, config RayLaunching3DConfig) *RayLaunching3D {
	return &RayLaunching3D{
		PowerMap:    matrix,
		WallNormals: wallNormals,
		Config:      config,
		RayPaths:    make([][]RayPoint, len(config.SingleRays)),
	}
}

func (rl *RayLaunching3D) calculateRayDirection(i, j int) (float64, float64, float64) {
	theta := (2 * math.Pi) / float64(rl.Config.NumOfRaysAzim) * float64(i)

	var phi, dx, dy, dz float64

	if rl.Config.TransmitterPos.Z == 0 {
		// half sphere – from 0 to π/2
		phi = ((math.Pi / 2) / float64(rl.Config.NumOfRaysElev)) * float64(j)
		dx = math.Cos(theta) * math.Cos(phi) * rl.Config.Step
		dy = math.Sin(theta) * math.Cos(phi) * rl.Config.Step
		dz = math.Sin(phi) * rl.Config.Step
	} else {
		//full sphere – from 0 to π
		phi = math.Pi*float64(j)/float64(rl.Config.NumOfRaysElev) - math.Pi/2
		dx = math.Cos(theta) * math.Cos(phi) * rl.Config.Step
		dy = math.Sin(theta) * math.Cos(phi) * rl.Config.Step
		dz = math.Sin(phi) * rl.Config.Step
	}

	dx = math.Round(dx*1e15) / 1e15
	dy = math.Round(dy*1e15) / 1e15
	dz = math.Round(dz*1e15) / 1e15

	return dx, dy, dz
}

func (rl *RayLaunching3D) getMapIndices(x, y, z float64) (int, int, int) {
	xIdx := int(math.Round(x / rl.Config.Step))
	yIdx := int(math.Round(y / rl.Config.Step))
	zIdx := int(math.Round(z / rl.Config.Step))
	return xIdx, yIdx, zIdx
}

func (rl *RayLaunching3D) isValidPosition(x, y, z float64) bool {
	return (x >= 0 && x <= rl.Config.SizeX) && (y >= 0 && y < rl.Config.SizeY) && (z <= rl.Config.SizeZ)
}

func (rl *RayLaunching3D) shouldContinueRay(state *RayState) bool {
	return rl.isValidPosition(state.x, state.y, state.z) &&
		state.currInteractions < rl.Config.NumOfInteractions &&
		state.currPower >= rl.Config.MinimalRayPower
}

func (rl *RayLaunching3D) shouldBreakRayPropagation(state *RayState, index int) bool {
	if index == rl.Config.BuldingInteriorNumber {
		return true
	} else {
		return false
	}
}

func (rl *RayLaunching3D) handleGroundReflection(state *RayState) {
	// reflection from the ground when z is below 0
	if state.z < 0 && state.currWallIndex != rl.Config.RoofMapNumber {
		state.currWallIndex = rl.Config.RoofMapNumber
		state.currInteractions++
		state.currSumRayLength += calculateDistance(state.currStartLengthPos, Point3D{X: state.x, Y: state.y, Z: state.z})
		state.currStartLengthPos = Point3D{X: state.x, Y: state.y, Z: state.z}
		nx, ny, nz := 0.0, 0.0, 1.0
		// calculate angle of incidence
		cosTheta := -(state.dx*nx + state.dy*ny + state.dz*nz)
		cosTheta = rl.clampCosTheta(cosTheta)
		theta := math.Acos(cosTheta)
		state.currReflectionFactor *= calculateReflectionFactor(theta, "medium-dry-ground")
		state.z = 0
	}

	if state.z < 0 {
		if state.dz < 0 {
			state.dz = -state.dz
		}
		state.z += state.dz
	}
}

func (rl *RayLaunching3D) handleRoofReflection(state *RayState, index int) bool {
	if (index == rl.Config.RoofMapNumber) && state.currWallIndex != rl.Config.RoofMapNumber {
		state.dz = -state.dz
		state.currWallIndex = rl.Config.RoofMapNumber
		state.currInteractions++
		state.currSumRayLength += calculateDistance(state.currStartLengthPos, Point3D{X: state.x, Y: state.y, Z: state.z})
		state.currStartLengthPos = Point3D{X: state.x, Y: state.y, Z: state.z}

		nx, ny, nz := 0.0, 0.0, 1.0
		cosTheta := -(state.dx*nx + state.dy*ny + state.dz*nz)
		cosTheta = rl.clampCosTheta(cosTheta)
		theta := math.Acos(cosTheta)
		state.currReflectionFactor *= calculateReflectionFactor(theta, "concrete")
		return true
	}
	return false
}

func (rl *RayLaunching3D) clampCosTheta(cosTheta float64) float64 {
	if cosTheta > 1 {
		return 1
	}
	if cosTheta < -1 {
		return -1
	}
	return cosTheta
}

func (rl *RayLaunching3D) calculateWallReflection(state *RayState, wallIndex int, i, j int) {
	currWallIndex := wallIndex - rl.Config.WallMapNumber

	//get wall normal
	nx, ny, nz := rl.WallNormals[currWallIndex].Nx, rl.WallNormals[currWallIndex].Ny, rl.WallNormals[currWallIndex].Nz
	//!!! MAP IS MIRRORED BY Y SO ALL Y NORMALS SHOULD BE MIRRORED !!!
	ny = -ny
	dot := state.dx*nx + state.dy*ny + state.dz*nz
	if dot > 0 {
		nx, ny, nz = -nx, -ny, -nz
		dot = -dot
	}
	dot *= 2
	// fmt.Printf("Promien %d, %d Nx=%.3f, Ny=%.3f, Nz=%.3f dot:%.3f \n", i, j, nx, ny, nz, dot)
	// println("stateDx", state.dx, "stateDy", state.dy, "stateDz", state.dz)

	cosTheta := -(state.dx*nx + state.dy*ny + state.dz*nz)
	cosTheta = rl.clampCosTheta(cosTheta)
	theta := math.Acos(cosTheta)
	// println("theta", theta)
	state.currReflectionFactor *= calculateReflectionFactor(theta, "concrete")
	state.dx = state.dx - dot*nx
	state.dy = state.dy - dot*ny
	state.dz = state.dz - dot*nz
	state.currInteractions++
	state.currWallIndex = wallIndex

	// sum distance and set new start position
	state.currSumRayLength += calculateDistance(state.currStartLengthPos, Point3D{X: state.x, Y: state.y, Z: state.z})
	state.currStartLengthPos = Point3D{X: state.x, Y: state.y, Z: state.z}
	// fmt.Println("Reflection Factor: %.3f", state.currReflectionFactor)
}

func (rl *RayLaunching3D) updatePowerMap(state *RayState, xIdx, yIdx, zIdx int) {
	state.currRayLength = calculateDistance(state.currStartLengthPos, Point3D{X: state.x, Y: state.y, Z: state.z}) + state.currSumRayLength

	d1 := state.toDiffractionPointRayLength
	d2 := state.currRayLength - state.toDiffractionPointRayLength

	lambda := rl.Config.WaveLength
	rel := float64(state.diffRayIndex) / float64(rl.Config.DiffractionRayNumber-1)
	alpha := rel * state.diffTheta
	baseLoss := bergDiffractionLoss(d1, d2, lambda, alpha)

	state.diffLossLdB = baseLoss
	H := calculateTransmittance(state.currRayLength, rl.Config.WaveLength, state.currReflectionFactor)
	absH := cmplx.Abs(H)
	if absH <= 0 {
		absH = 1e-15
	}
	txPower := rl.Config.TransmitterPower
	if txPower <= 0 {
		txPower = 1e-15
	}
	state.currPower = 10*math.Log10(txPower) + 20*math.Log10(absH) - state.diffLossLdB
	// println("baseLoss: ", baseLoss, "rayIndex: ", state.diffRayIndex)
	// println("currPorwe: ", state.currPower)

	// update power map if power is higher than previous one
	if rl.PowerMap[zIdx][yIdx][xIdx] == -160 || rl.PowerMap[zIdx][yIdx][xIdx] < state.currPower {
		rl.PowerMap[zIdx][yIdx][xIdx] = state.currPower
	}
}

func (rl *RayLaunching3D) addToRayPath(targetRayIndex int, state *RayState) {
	if targetRayIndex >= 0 {
		rl.RayPaths[targetRayIndex] = append(rl.RayPaths[targetRayIndex], RayPoint{
			X:     float64(math.Round(state.x / rl.Config.Step)),
			Y:     float64(math.Round(state.y / rl.Config.Step)),
			Z:     float64(math.Round(state.z / rl.Config.Step)),
			Power: state.currPower,
		})
	}
}

func (rl *RayLaunching3D) processCornerDiffraction(state *RayState, xIdx, yIdx, zIdx int, i, j int, diffractionRayNumber int, index int) {
	state.currWallIndex = index
	state.currSumRayLength += calculateDistance(state.currStartLengthPos, Point3D{X: state.x, Y: state.y, Z: state.z})
	state.toDiffractionPointRayLength = state.currSumRayLength
	state.currStartLengthPos = Point3D{X: state.x, Y: state.y, Z: state.z}
	normals := getNeighborWallNormals(xIdx, yIdx, zIdx, rl)
	// fmt.Printf("xIdx: %v yIdx: %v zIdx: %v dx: %v dy: %v dz: %v rayLength: %.3f \n", xIdx, yIdx, zIdx, state.dx, state.dy, state.dz, state.currSumRayLength)

	if len(normals) == 0 {
		return
	}

	bestDot := -math.MaxFloat64
	var bestNormal Normal3D
	if index == rl.Config.CornerMapNumber {
		for _, n := range normals {
			n.Nx, n.Ny, n.Nz = -n.Nx, n.Ny, n.Nz
			dot := state.dx*n.Nx + state.dy*n.Ny + state.dz*n.Nz

			if dot > bestDot {
				bestDot = dot
				bestNormal = n

			}
			// fmt.Printf("Promien %d, %d Normalna %d: Nx=%.3f, Ny=%.3f, Nz=%.3f Dot=%.3f\n", i, j, k, n.Nx, n.Ny, n.Nz, dot)
		}
		// fmt.Printf("Best Dot: %.3f, bestNormal: %v \n", bestDot, bestNormal)
		cosTheta := state.dx*bestNormal.Nx + state.dy*bestNormal.Ny + state.dz*bestNormal.Nz
		cosTheta = rl.clampCosTheta(cosTheta)
		theta := math.Acos(cosTheta)
		cross := state.dx*(-bestNormal.Ny) - state.dy*bestNormal.Nx
		finalTheta := (math.Pi/2 - theta)
		if cross < 0 {
			finalTheta = -finalTheta
		}
		stepResolution := diffractionRayNumber

		oneStep := computeOneStep(bestNormal.Nx, bestNormal.Ny, state.dx, state.dy, finalTheta, stepResolution)

		// fmt.Printf("cosTheta: %.3f theta: %.3f, finalTheta: %.3f cross: %.3f oneStep: %.3f\n", cosTheta, theta, finalTheta, cross, oneStep)

		state.diffTheta = math.Abs(finalTheta)
		rl.processDiffractionSteps(state, oneStep, stepResolution, i, j, normals, index)
	} else if index == rl.Config.RoofCornerMapNumber {
		stepResolution := diffractionRayNumber
		var startDz, endDz, theta, finalTheta float64
		if state.dz < 0 {
			cosTheta := state.dx*0 + state.dy*0 + state.dz*-1
			cosTheta = rl.clampCosTheta(cosTheta)
			theta = math.Acos(cosTheta)
			finalTheta = (math.Pi/2 - theta)
			startDz = state.dz
			endDz = -1.0
		} else {
			cosTheta := state.dx*0 + state.dy*0 + state.dz*1
			cosTheta = rl.clampCosTheta(cosTheta)
			theta = math.Acos(cosTheta)
			finalTheta = (math.Pi/2 - theta)
			startDz = state.dz
			endDz = 0.0
		}

		oneStep := (endDz - startDz) / float64(stepResolution)

		// fmt.Printf(
		// 	"Roof: dz=%.3f -> range [%.3f .. %.3f], oneStep=%.3f theta=%.3f finalTheta=%.3f\n",
		// 	state.dz, startDz, endDz, oneStep, theta, finalTheta,
		// )
		state.diffTheta = finalTheta
		rl.processDiffractionSteps(state, oneStep, stepResolution, i, j, normals, index)
	}

}

func (rl *RayLaunching3D) processDiffractionSteps(state *RayState, oneStep float64, stepResolution int, i, j int, normalsAround []Normal3D, index int) {
	for step := 0; step <= stepResolution; step++ {
		x := state.x
		y := state.y
		z := state.z
		state.diffRayIndex = step
		var newDx, newDy, newDz float64
		if index == rl.Config.RoofCornerMapNumber {
			if state.dz > 0 {
				newDz = state.dz + float64(step)*oneStep
				newDx = state.dx
				newDy = state.dy
			} else if state.dz <= 0 {
				if index == rl.Config.RoofCornerMapNumber {
					t := float64(step) / float64(stepResolution+1)
					newDz = state.dz + float64(step)*oneStep
					factor := 1.0 - t
					newDx = state.dx * factor
					newDy = state.dy * factor
				}
			}
		} else {
			angle := float64(step) * oneStep
			newDx = state.dx*math.Cos(angle) - state.dy*math.Sin(angle)
			newDy = state.dx*math.Sin(angle) + state.dy*math.Cos(angle)
			newDz = state.dz
		}

		// fmt.Printf("Step %d: dx=%.3f dy=%.3f dz=%.3f\n", step, newDx, newDy, newDz)

		x += newDx
		y += newDy
		z += newDz

		// fmt.Printf("x: %v y: %v z: %v \n", x, y, z)
		rl.processDiffractionRayPath(x, y, z, newDx, newDy, newDz, *state, i, j, normalsAround, index)
	}
}

func (rl *RayLaunching3D) processDiffractionRayPath(x, y, z, newDx, newDy, newDz float64, state RayState, i, j int, normalsAround []Normal3D, startIndex int) {
	state.dx, state.dy, state.dz = newDx, newDy, newDz
	state.x, state.y, state.z = x, y, z
	for rl.shouldContinueRay(&state) {
		rl.handleGroundReflection(&state)
		xIdx, yIdx, zIdx := rl.getMapIndices(state.x, state.y, state.z)
		index := int(rl.PowerMap[zIdx][yIdx][xIdx])
		// fmt.Println("xIdx: ", xIdx, "yIdx: ", yIdx, "zIdx: ", zIdx, "index: ", index, "currWallIndex: ", state.currWallIndex)
		// fmt.Println("x: ", state.x, "y: ", state.y, "z: ", state.z)
		if rl.shouldBreakRayPropagation(&state, index) || ((state.currWallIndex == rl.Config.RoofMapNumber) && newDz > 0) {
			break
		}

		if !(startIndex == rl.Config.RoofCornerMapNumber && newDz > 0) && rl.handleRoofReflection(&state, index) {
			continue
		}

		if index >= rl.Config.WallMapNumber && index < rl.Config.RoofMapNumber && index != state.currWallIndex+rl.Config.WallMapNumber {
			normalIndex := index - rl.Config.WallMapNumber

			if !containsNormal(normalsAround, rl.WallNormals[normalIndex]) {
				rl.calculateWallReflection(&state, index, i, j)
			}
			rl.updatePowerMap(&state, xIdx, yIdx, zIdx)
			rl.addToRayPath(state.targetRayIndex, &state)
		} else {
			rl.updatePowerMap(&state, xIdx, yIdx, zIdx)
			rl.addToRayPath(state.targetRayIndex, &state)
		}
		state.x += state.dx
		state.y += state.dy
		state.z += state.dz
	}
}

func (rl *RayLaunching3D) CreatePowerMapLegend() {
	legend := make(map[int]PowerMapLegendEntry)

	for z := 0; z < len(rl.PowerMap); z++ {
		entry := PowerMapLegendEntry{}
		var coveredPoints, totalPoints float64

		for y := 0; y < len(rl.PowerMap[z]); y++ {
			for x := 0; x < len(rl.PowerMap[z][y]); x++ {
				power := rl.PowerMap[z][y][x]
				if power < 0 {
					totalPoints++
					if power > -159.9 {
						coveredPoints++
						loss := math.Abs(power)
						switch {
						case loss < 20:
							entry.P0_20++
						case loss < 40:
							entry.P20_40++
						case loss < 60:
							entry.P40_60++
						case loss < 80:
							entry.P60_80++
						case loss < 100:
							entry.P80_100++
						case loss < 120:
							entry.P100_120++
						case loss < 140:
							entry.P120_140++
						default:
							entry.P140plus++
						}
					}
				}
			}
		}

		if coveredPoints > 0 {
			scale := 100.0 / coveredPoints
			entry.P0_20 *= scale
			entry.P20_40 *= scale
			entry.P40_60 *= scale
			entry.P60_80 *= scale
			entry.P80_100 *= scale
			entry.P100_120 *= scale
			entry.P120_140 *= scale
			entry.P140plus *= scale
		}

		entry.Ptotal = (coveredPoints / totalPoints) * 100.0

		legend[z] = entry
	}

	rl.PowerMapLegend = legend
}

func (rl *RayLaunching3D) PrintPowerMapLegend() {
	fmt.Println("===== Power Map Legend per Floor =====")
	for z := 0; z < len(rl.PowerMap); z++ {
		entry := rl.PowerMapLegend[z]
		fmt.Printf("Floor z = %d:\n", z)
		fmt.Printf("  Total coverage: %.2f%%\n", entry.Ptotal)
		fmt.Printf("  0-20 dB   : %.2f%%\n", entry.P0_20)
		fmt.Printf("  20-40 dB  : %.2f%%\n", entry.P20_40)
		fmt.Printf("  40-60 dB  : %.2f%%\n", entry.P40_60)
		fmt.Printf("  60-80 dB  : %.2f%%\n", entry.P60_80)
		fmt.Printf("  80-100 dB : %.2f%%\n", entry.P80_100)
		fmt.Printf("  100-120 dB: %.2f%%\n", entry.P100_120)
		fmt.Printf("  120-140 dB: %.2f%%\n", entry.P120_140)
		fmt.Printf("  >140 dB   : %.2f%%\n", entry.P140plus)
		fmt.Println("--------------------------------------")
	}
}

func (rl *RayLaunching3D) CalculateRayLaunching3D() {
	for z := 0; z < int(rl.Config.TransmitterPos.Z); z++ {
		rl.PowerMap[z][int(rl.Config.TransmitterPos.Y)][int(rl.Config.TransmitterPos.X)] = 0
	}

	for i := 0; i < rl.Config.NumOfRaysAzim; i++ { // loop over horizontal dim
		for j := 0; j < rl.Config.NumOfRaysElev; j++ { // loop over vertical dim
			dx, dy, dz := rl.calculateRayDirection(i, j)
			// main loop
			// if !(i == 15 && j == 13) {
			// 	continue
			// }
			targetRayIndex := rl.isTargetRay(i, j)
			state := &RayState{
				x:  rl.Config.TransmitterPos.X + dx,
				y:  rl.Config.TransmitterPos.Y + dy,
				z:  rl.Config.TransmitterPos.Z + dz,
				dx: dx, dy: dy, dz: dz,
				currInteractions:            0,
				currPower:                   0.0,
				currWallIndex:               0,
				currStartLengthPos:          Point3D{X: rl.Config.TransmitterPos.X, Y: rl.Config.TransmitterPos.Y, Z: rl.Config.TransmitterPos.Z},
				currRayLength:               0.0,
				currSumRayLength:            0.0,
				currReflectionFactor:        1.0,
				diffLossLdB:                 0.0,
				targetRayIndex:              targetRayIndex,
				toDiffractionPointRayLength: 0.0,
				diffTheta:                   0.0,
				diffRayIndex:                0,
			}

			for rl.shouldContinueRay(state) {
				// reflection from the ground when z is below 0
				rl.handleGroundReflection(state)

				xIdx, yIdx, zIdx := rl.getMapIndices(state.x, state.y, state.z)
				index := int(rl.PowerMap[zIdx][yIdx][xIdx])
				// fmt.Println("xIdx: ", xIdx, "yIdx: ", yIdx, "zIdx: ", zIdx, "index: ", index, "currWallIndex: ", state.currWallIndex)
				if rl.shouldBreakRayPropagation(state, index) || (index == rl.Config.RoofCornerMapNumber && state.dz == 0) {
					break
				}
				// reflection from the building roof
				if rl.handleRoofReflection(state, index) {
					continue
				}

				if (index == rl.Config.CornerMapNumber && state.currWallIndex != rl.Config.CornerMapNumber) || (index == rl.Config.RoofCornerMapNumber && state.currWallIndex != rl.Config.CornerMapNumber && !(state.currWallIndex >= rl.Config.WallMapNumber && index < rl.Config.RoofMapNumber)) {
					if rl.Config.DiffractionRayNumber < 2 {
						break
					}
					nextXIdx, nextYIdx, nextZIdx := rl.getMapIndices(state.x+state.dx, state.y+state.dy, state.z+state.dz)
					// fmt.Println("NextXIdx: ", nextXIdx, "NextYIdx: ", nextYIdx, "NextZIdx: ", nextZIdx)
					if nextZIdx < 0 || !rl.isValidPosition(float64(nextXIdx), float64(nextYIdx), float64(nextZIdx)) {
						break
					}

					nextIndex := int(rl.PowerMap[nextZIdx][nextYIdx][nextXIdx])
					if nextIndex == rl.Config.RoofMapNumber {
						break
					}

					if !(state.currWallIndex >= rl.Config.WallMapNumber && state.currWallIndex < rl.Config.RoofMapNumber) {
						rl.processCornerDiffraction(state, xIdx, yIdx, zIdx, i, j, rl.Config.DiffractionRayNumber-1, index)
						break
					}
				}

				if index >= rl.Config.WallMapNumber && index < rl.Config.RoofMapNumber && index != state.currWallIndex {
					rl.calculateWallReflection(state, index, i, j)
				} else {
					rl.updatePowerMap(state, xIdx, yIdx, zIdx)
					rl.addToRayPath(targetRayIndex, state)
				}

				// update position
				state.x += state.dx
				state.y += state.dy
				state.z += state.dz
			}
		}
	}

	rl.CreatePowerMapLegend()
}

func calculateDistance(p1, p2 Point3D) float64 {
	dist := math.Sqrt(math.Pow(p1.X-p2.X, 2) + math.Pow(p1.Y-p2.Y, 2) + math.Pow(p1.Z-p2.Z, 2))
	return dist
}

func calculateTransmittance(r, waveLength, reflectionRef float64) complex128 {
	const epsilon = 1e-6
	if r < epsilon {
		r = epsilon
	}
	H := complex(reflectionRef, 0) *
		complex(waveLength/(4*math.Pi*r), 0) *
		cmplx.Exp(complex(0, -2*math.Pi*r/waveLength))

	return H
}

func (rl *RayLaunching3D) isTargetRay(i, j int) int {
	for idx, singleRay := range rl.Config.SingleRays {
		if i-singleRay.Azimuth == 0 && j-singleRay.Elevation == 0 {
			return idx
		}
	}
	return -1
}

func calculateReflectionFactor(angle float64, material string) float64 {
	if angle > math.Pi/2 {
		angle = math.Pi - angle
	}
	var eta float64
	switch material {
	case "concrete":
		eta = 5.31
	case "ceiling-board":
		eta = 1.50
	case "medium-dry-ground":
		eta = 15
	}
	sinTheta := math.Sin(angle)
	cosTheta := math.Cos(angle)
	if cosTheta > 1 {
		cosTheta = 1
	}
	if cosTheta < -1 {
		cosTheta = -1
	}
	v := eta - sinTheta*sinTheta
	if v < 0 {
		v = 0
	}
	root := math.Sqrt(v)
	R_TE := (cosTheta - root) / (cosTheta + root)
	R_TM := (eta*cosTheta - root) / (eta*cosTheta + root)
	reflectionFactor := (math.Pow(R_TE, 2) + math.Pow(R_TM, 2)) / 2
	return reflectionFactor
}

func containsNormal(normals []Normal3D, target Normal3D) bool {
	for _, n := range normals {
		if n == target {
			return true
		}
	}
	return false
}

func getNeighborWallNormals(x, y, z int, rl *RayLaunching3D) []Normal3D {
	neighborNormals := make(map[int]Normal3D)

	size := 3
	for dx := -size; dx <= size; dx++ {
		for dy := -size; dy <= size; dy++ {
			for dz := -size; dz <= size; dz++ {
				xprim := x + dx
				yprim := y + dy
				zprim := z + dz

				if xprim < 0 || yprim < 0 || zprim < 0 || xprim >= int(rl.Config.SizeX) || yprim >= int(rl.Config.SizeY) || zprim >= int(rl.Config.SizeZ) {
					continue
				}

				index := int(rl.PowerMap[zprim][yprim][xprim])
				if index >= rl.Config.WallMapNumber && index < rl.Config.RoofMapNumber {
					currWallIndex := index - rl.Config.WallMapNumber
					if _, exists := neighborNormals[currWallIndex]; !exists {
						neighborNormals[currWallIndex] = rl.WallNormals[currWallIndex]
					}
				}
			}
		}
	}
	result := make([]Normal3D, 0, len(neighborNormals))
	for _, normal := range neighborNormals {
		result = append(result, normal)
	}

	return result
}

func computeOneStep(Nx, Ny, dx, dy, finalTheta float64, stepResolution int) float64 {
	if stepResolution == 0 {
		return 0
	}

	nlen := math.Hypot(Nx, Ny)
	ilen := math.Hypot(dx, dy)
	if nlen == 0 || ilen == 0 {
		return 0
	}
	Nx /= nlen
	Ny /= nlen
	dx /= ilen
	dy /= ilen

	cross := Nx*dy - Ny*dx
	dot := Nx*dx + Ny*dy

	angle := math.Atan2(cross, dot)

	eps := 1e-12
	absFinal := math.Abs(finalTheta)
	if math.Abs(angle) < eps {
		return finalTheta / float64(stepResolution)
	}

	oneStep := math.Copysign(absFinal/float64(stepResolution), angle)
	return oneStep
}

func diffractionLoss(d1, d2, lambda, alpha float64) float64 {
	println("BergDiffractionLoss d1:", d1, " d2:", d2, " lambda:", lambda, " alpha:", alpha)
	if d1 <= 0 || d2 <= 0 || lambda <= 0 || alpha <= 0 || alpha >= math.Pi {
		return 0
	}
	n := math.Pi / alpha
	if n <= 1 {
		return 0
	}

	term1 := (n / (n - 1.0)) * (math.Sin(math.Pi/n) / math.Cos(math.Pi/(2*n)))
	term2 := (d1 * d2) / ((d1 + d2) * lambda)

	Ld := 20*math.Log10(term1) + 20*math.Log10(term2)
	if Ld < 0 {
		Ld = 0
	}
	return Ld
}

func bergDiffractionLoss(d1, d2, lambda, alpha_rad float64) float64 {
	const v = 1.5
	const q_lambda = 0.1
	if d1 <= 0 || d2 <= 0 || lambda <= 0 || alpha_rad <= 0 || alpha_rad > math.Pi {
		return 0
	}

	alpha_deg := alpha_rad * 180.0 / math.Pi
	q90 := math.Sqrt(q_lambda / lambda)
	// println("alpha_deg:", alpha_deg, " q90:", q90)
	q1 := q90 * math.Pow(alpha_deg/90.0, v)

	d_real := d1 + d2
	d_final := d_real + (d1 * d2 * q1)
	Ld_diff_only := 20 * math.Log10(d_final/d_real)
	if Ld_diff_only < 0 {
		Ld_diff_only = 0
	}
	return Ld_diff_only
}
func AngularDiffractionLoss(rayIdx, nRays int, maxAngleRad, n float64) float64 {
	rel := float64(rayIdx) / float64(nRays-1)
	phi := rel * maxAngleRad

	cosPhi := math.Cos(phi)
	if cosPhi < 1e-6 {
		cosPhi = 1e-6
	}

	additionalLoss := -20 * n * math.Log10(cosPhi)
	return additionalLoss
}

func checkNaN(label string, value float64) bool {
	if math.IsNaN(value) || math.IsInf(value, 0) {
		fmt.Printf("⚠️ NaN detected at %s: %v\n", label, value)
		return true
	}
	return false
}
