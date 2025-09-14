package raylaunching

import (
	. "backendGo/types"
	"fmt"

	// "fmt"
	"math"
	"math/cmplx"
)

type RayLaunching3DConfig struct {
	NumOfRaysAzim, NumOfRaysElev, NumOfInteractions, WallMapNumber, BuldingInteriorNumber, RoofMapNumber, CornerMapNumber int
	SizeX, SizeY, SizeZ, Step, ReflFactor, TransmitterPower, MinimalRayPower, TransmitterFreq, WaveLength float64
	TransmitterPos Point3D
	SingleRays []SingleRay
}

type RayPoint struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
	Z float64 `json:"z"`
	Power    float64 `json:"power"`
}
type RayLaunching3D struct {
	PowerMap [][][]float64
	WallNormals []Normal3D
	Config RayLaunching3DConfig
	RayPaths [][]RayPoint  
}

type RayState struct {
	x, y, z float64
	dx, dy, dz float64
	currInteractions int
	currPower float64
	currWallIndex int
	currStartLengthPos Point3D
	currRayLength float64
	currSumRayLength float64
	currReflectionFactor float64
	diffLossLdB float64
	targetRayIndex int
	toDiffractionPointRayLength float64
}


func NewRayLaunching3D(matrix [][][]float64, wallNormals []Normal3D, config RayLaunching3DConfig) *RayLaunching3D {
	return &RayLaunching3D{
		PowerMap: matrix,
		WallNormals: wallNormals,
		Config: config,
		RayPaths: make([][]RayPoint, len(config.SingleRays)),
	}
}

func (rl *RayLaunching3D) calculateRayDirection(i, j int) (float64, float64, float64) {
	theta := (2*math.Pi)/float64(rl.Config.NumOfRaysAzim)*float64(i)
	
	var phi, dx, dy, dz float64
	
	if rl.Config.TransmitterPos.Z == 0 {
		// half sphere – from 0 to π/2
		phi = ((math.Pi/2) / float64(rl.Config.NumOfRaysElev)) * float64(j)
		dx = math.Cos(theta) * math.Cos(phi) * rl.Config.Step 
		dy = math.Sin(theta) * math.Cos(phi) * rl.Config.Step
		dz = math.Sin(phi) * rl.Config.Step
	} else {
		//full sphere – from 0 to π
		phi = math.Pi * float64(j) / float64(rl.Config.NumOfRaysElev) - math.Pi/2
		dx = math.Cos(theta) * math.Cos(phi) * rl.Config.Step 
		dy = math.Sin(theta) * math.Cos(phi) * rl.Config.Step
		dz = math.Sin(phi) * rl.Config.Step
	}
	
	dx = math.Round(dx*1e15)/1e15
	dy = math.Round(dy*1e15)/1e15
	dz = math.Round(dz*1e15)/1e15
	
	return dx, dy, dz
}

func (rl *RayLaunching3D) getMapIndices(x, y, z float64) (int, int, int) {
	xIdx := int(math.Round(x/rl.Config.Step))
	yIdx := int(math.Round(y/rl.Config.Step))
	zIdx := int(math.Round(z/rl.Config.Step))
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
	fmt.Printf("Promien %d, %d Nx=%.3f, Ny=%.3f, Nz=%.3f dot:%.3f \n", i, j, nx, ny, nz, dot)
	println("stateDx", state.dx,"stateDy", state.dy,"stateDz", state.dz)
	
	cosTheta := -(state.dx*nx + state.dy*ny + state.dz*nz)
	cosTheta = rl.clampCosTheta(cosTheta)
	theta := math.Acos(cosTheta)
	println("theta", theta)
	state.currReflectionFactor *= calculateReflectionFactor(theta, "concrete")
	state.dx = state.dx - dot*nx
	state.dy = state.dy - dot*ny
	state.dz = state.dz - dot*nz
	state.currInteractions++
	state.currWallIndex = wallIndex

	// sum distance and set new start position
	state.currSumRayLength += calculateDistance(state.currStartLengthPos, Point3D{X: state.x, Y: state.y, Z: state.z})
	state.currStartLengthPos = Point3D{X: state.x, Y: state.y, Z: state.z}
	fmt.Println("Reflection Factor: %.3f", state.currReflectionFactor)
}

func (rl *RayLaunching3D) updatePowerMap(state *RayState, xIdx, yIdx, zIdx int) {
	state.currRayLength = calculateDistance(state.currStartLengthPos, Point3D{X: state.x, Y: state.y, Z: state.z}) + state.currSumRayLength

	d1 := state.toDiffractionPointRayLength 
	d2 := state.currRayLength - state.toDiffractionPointRayLength
	lambda := rl.Config.WaveLength
	alpha := math.Pi / 2 
	Ld := BergDiffractionLoss(d1, d2, lambda, alpha)
	state.diffLossLdB = Ld

	H := calculateTransmittance(state.currRayLength, rl.Config.WaveLength, state.currReflectionFactor)
	state.currPower = 10*math.Log10(rl.Config.TransmitterPower) + 20*math.Log10(cmplx.Abs(H)) - state.diffLossLdB
	
		println("diffLoss: ",state.diffLossLdB)
		println("currPorwe: ",state.currPower)
	

	// update power map if power is higher than previous one
	if rl.PowerMap[zIdx][yIdx][xIdx] == -150 || rl.PowerMap[zIdx][yIdx][xIdx] < state.currPower {
		rl.PowerMap[zIdx][yIdx][xIdx] = state.currPower
	}
}

func (rl *RayLaunching3D) addToRayPath(targetRayIndex int, state *RayState) {
	if targetRayIndex >= 0 {
		rl.RayPaths[targetRayIndex] = append(rl.RayPaths[targetRayIndex], RayPoint{
			X: float64(math.Round(state.x/rl.Config.Step)), 
			Y: float64(math.Round(state.y/rl.Config.Step)), 
			Z: float64(math.Round(state.z/rl.Config.Step)), 
			Power: state.currPower,
		})
	}
}

func (rl *RayLaunching3D) processCornerDiffraction(state *RayState, xIdx, yIdx, zIdx int, i, j int) {
	state.currWallIndex = rl.Config.CornerMapNumber
	state.currSumRayLength += calculateDistance(state.currStartLengthPos, Point3D{X: state.x, Y: state.y, Z: state.z})
	state.toDiffractionPointRayLength = state.currSumRayLength
	state.currStartLengthPos = Point3D{X: state.x, Y: state.y, Z: state.z}
	normals := getNeighborWallNormals(xIdx, yIdx, zIdx, rl)
	fmt.Printf("xIdx: %v yIdx: %v zIdx: %v dx: %v dy: %v dz: %v rayLength: %.3f \n", xIdx, yIdx, zIdx, state.dx, state.dy, state.dz, state.currSumRayLength)
	
	if len(normals) == 0 {
		return
	}
	
	bestDot := -math.MaxFloat64
	var bestNormal Normal3D
	
	for k, n := range normals {
		n.Nx, n.Ny, n.Nz = -n.Nx, n.Ny, n.Nz
		dot := state.dx*n.Nx + state.dy*n.Ny + state.dz*n.Nz
		
		if dot > bestDot {
			bestDot = dot
			bestNormal = n
			
		}
		fmt.Printf("Promien %d, %d Normalna %d: Nx=%.3f, Ny=%.3f, Nz=%.3f Dot=%.3f\n", i, j, k, n.Nx, n.Ny, n.Nz, dot)
	}
	
	fmt.Printf("Best Dot: %.3f, bestNormal: %v \n", bestDot, bestNormal)
	cosTheta := state.dx*bestNormal.Nx + state.dy*bestNormal.Ny + state.dz*bestNormal.Nz
	cosTheta = rl.clampCosTheta(cosTheta)
	theta := math.Acos(cosTheta)
	cross := state.dx*(-bestNormal.Ny) - state.dy*bestNormal.Nx
	finalTheta := (math.Pi/2 - theta)
	if cross < 0 {
		finalTheta = -finalTheta
	}
	stepResolution := 20
	
	oneStep := computeOneStep(bestNormal.Nx, bestNormal.Ny, state.dx, state.dy, finalTheta, stepResolution)
	
	fmt.Printf("cosTheta: %.3f theta: %.3f, finalTheta: %.3f cross: %.3f oneStep: %.3f\n", cosTheta, theta, finalTheta, cross, oneStep)
	
	rl.processDiffractionSteps(state, oneStep, stepResolution, i, j, normals)
}

func (rl *RayLaunching3D) processDiffractionSteps(state *RayState, oneStep float64, stepResolution int, i, j int, normalsAround []Normal3D) {
	for step := 0; step <= stepResolution; step++ {
		x := state.x
		y := state.y
		z := state.z
		angle := float64(step) * oneStep
		newDx := state.dx * math.Cos(angle) - state.dy * math.Sin(angle)
		newDy := state.dx * math.Sin(angle) + state.dy * math.Cos(angle)
		fmt.Printf("Step %d: dx=%.3f dy=%.3f dz=%.3f\n", step, newDx, newDy, state.dz)

		x += newDx
		y += newDy
		z += state.dz

		fmt.Printf("x: %v y: %v z: %v \n", x, y, z)
		rl.processDiffractionRayPath(x, y, z, newDx, newDy, *state, i, j, normalsAround)
	}
}


func (rl *RayLaunching3D) processDiffractionRayPath(x, y, z, newDx, newDy float64, state RayState, i, j int, normalsAround []Normal3D) {
	state.dx, state.dy, state.dz = newDx, newDy, state.dz
	state.x, state.y, state.z = x, y, z
	for rl.shouldContinueRay(&state) {
		rl.handleGroundReflection(&state)
		xIdx, yIdx, zIdx := rl.getMapIndices(state.x, state.y, state.z)
		index := int(rl.PowerMap[zIdx][yIdx][xIdx])

		if (rl.shouldBreakRayPropagation(&state, index)) {
				break
		}

		if rl.handleRoofReflection(&state, index) {
			continue
		}

		if (index == rl.Config.CornerMapNumber) {
			break
		}

		if index >= rl.Config.WallMapNumber && index < rl.Config.RoofMapNumber && index != state.currWallIndex + rl.Config.WallMapNumber {
			normalIndex:= index - rl.Config.WallMapNumber
			
			if !containsNormal(normalsAround, rl.WallNormals[normalIndex]) {
				rl.calculateWallReflection(&state, index, i, j)
			} 
			rl.updatePowerMap(&state, xIdx, yIdx, zIdx,)
			rl.addToRayPath(state.targetRayIndex, &state)
		} else {
			rl.updatePowerMap(&state, xIdx, yIdx, zIdx,)
			rl.addToRayPath(state.targetRayIndex, &state)
		}
		state.x += state.dx
		state.y += state.dy
		state.z += state.dz
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
			//  if !(i == 7 && j == 5) {
        	// 	continue 
        	// }
			targetRayIndex := rl.isTargetRay(i, j)
			state := &RayState{
				x: rl.Config.TransmitterPos.X + dx,
				y: rl.Config.TransmitterPos.Y + dy,
				z: rl.Config.TransmitterPos.Z + dz,
				dx: dx, dy: dy, dz: dz,
				currInteractions: 0,
				currPower: 0.0,
				currWallIndex: 0,
				currStartLengthPos: Point3D{X: rl.Config.TransmitterPos.X, Y: rl.Config.TransmitterPos.Y, Z: rl.Config.TransmitterPos.Z},
				currRayLength: 0.0,
				currSumRayLength: 0.0,
				currReflectionFactor: 1.0,
				diffLossLdB: 0.0,
				targetRayIndex:targetRayIndex,
				toDiffractionPointRayLength: 0.0,
			}


			for rl.shouldContinueRay(state) {
				// reflection from the ground when z is below 0
				rl.handleGroundReflection(state)
				
				xIdx, yIdx, zIdx := rl.getMapIndices(state.x, state.y, state.z)
				index := int(rl.PowerMap[zIdx][yIdx][xIdx])
				fmt.Println("xIdx: ", xIdx, "yIdx: ", yIdx, "zIdx: ", zIdx, "index: ", index, "currWallIndex: ", state.currWallIndex)
				if (rl.shouldBreakRayPropagation(state, index)) {
					break
				}
				// reflection from the building roof
				if rl.handleRoofReflection(state, index) {
					continue
				}
				
				if index == rl.Config.CornerMapNumber && state.currWallIndex != rl.Config.CornerMapNumber {
					rl.processCornerDiffraction(state, xIdx, yIdx, zIdx, i, j)
				}

				if index >= rl.Config.WallMapNumber && index < rl.Config.RoofMapNumber && index != state.currWallIndex  {
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
}


// func (rl *RayLaunching3D) CalculateRayLaunching3D() {
// 	for z := 0; z < int(rl.Config.TransmitterPos.Z); z++ {
// 		rl.PowerMap[z][int(rl.Config.TransmitterPos.Y)][int(rl.Config.TransmitterPos.X)] = 0
// 	}
// 	for i := 0; i < rl.Config.NumOfRaysAzim; i++ { // loop over horizontal dim
// 		theta := (2*math.Pi)/float64(rl.Config.NumOfRaysAzim)*float64(i) // from -π to π
// 		for j := 0; j < rl.Config.NumOfRaysElev; j++ { // loop over vertical dim
				
// 			var phi,dx,dy,dz float64

// 			// spherical coordinates
// 			if rl.Config.TransmitterPos.Z == 0 {
// 				// half sphere – from 0 to π/2
// 				phi = ((math.Pi/2) / float64(rl.Config.NumOfRaysElev)) *  float64(j) // from 0 to π/2
// 				dx = math.Cos(theta) * math.Cos(phi) * rl.Config.Step 
// 				dy = math.Sin(theta) * math.Cos(phi) * rl.Config.Step
// 				dz = math.Sin(phi) * rl.Config.Step
// 			} else {
// 				//full sphere – from 0 to π
// 				phi = math.Pi * float64(j) / float64(rl.Config.NumOfRaysElev) - math.Pi/2// from 0 to π
// 				dx = math.Cos(theta) * math.Cos(phi) * rl.Config.Step 
// 				dy = math.Sin(theta) * math.Cos(phi) * rl.Config.Step
// 				dz = math.Sin(phi) * rl.Config.Step
				
// 			}
// 			dx, dy, dz = math.Round(dx*1e15)/1e15, math.Round(dy*1e15)/1e15, math.Round(dz*1e15)/1e15

// 			/* getting past to next step,
// 			 omitting calculations for transmitter */

// 			x := rl.Config.TransmitterPos.X + dx
// 			y := rl.Config.TransmitterPos.Y + dy
// 			z := rl.Config.TransmitterPos.Z + dz

// 			targetRayIndex := rl.isTargetRay(i, j)

// 			// initial counters
// 			currInteractions := 0
// 			currPower := 0.0
// 			currWallIndex := 0
// 			currStartLengthPos := Point3D{X:rl.Config.TransmitterPos.X, Y:rl.Config.TransmitterPos.Y, Z:rl.Config.TransmitterPos.Z}
// 			currRayLength := 0.0
// 			currSumRayLength := 0.0
// 			currReflectionFactor := 1.0
// 			diffLossLdB:=0.0
// 			// main loop
// 			for (x >= 0 && x <= rl.Config.SizeX) && (y >= 0 && y < rl.Config.SizeY) && (z <= rl.Config.SizeZ) && currInteractions < rl.Config.NumOfInteractions && currPower >= rl.Config.MinimalRayPower {
// 				// reflection from the ground when z is below 0
// 				if (z < 0 && currWallIndex != rl.Config.RoofMapNumber) {
// 					dz = -dz
// 					currWallIndex = rl.Config.RoofMapNumber
// 					currInteractions++
// 					currSumRayLength += calculateDistance(currStartLengthPos, Point3D{X: x, Y: y, Z: z})
// 					currStartLengthPos = Point3D{X: x, Y: y, Z: z}
// 					nx, ny, nz := 0.0, 0.0, 1.0
// 					// calculate angle of incidence
// 					cosTheta := -(dx*nx + dy*ny + dz*nz)
// 					theta := math.Acos(cosTheta)
// 					currReflectionFactor *= calculateReflectionFactor(theta, "medium-dry-ground")
// 					z = 0
// 				}
// 				if (z < 0) {
// 					if (dz < 0) {
// 						dz = -dz	
// 					}
// 					z += dz
// 				}
// 				xIdx := int(math.Round(x/rl.Config.Step))
// 				yIdx := int(math.Round(y/rl.Config.Step))
// 				zIdx := int(math.Round(z/rl.Config.Step))
// 				index := int(rl.PowerMap[zIdx][yIdx][xIdx])
// 				// if (i==0 && j == 10) {
// 					// println("i:", i, "j:", j, "x:", xIdx, "y:", yIdx, "z:", zIdx, "index:", index,"currWallIndex:", currWallIndex, "dx:", dx, "dy:", dy, "dz:", dz)
// 				// }
// 				// reflection from the building roof
// 				if (index == rl.Config.RoofMapNumber) && currWallIndex != rl.Config.RoofMapNumber {
// 					dz = -dz
// 					currWallIndex = rl.Config.RoofMapNumber
// 					currInteractions++
// 					currSumRayLength += calculateDistance(currStartLengthPos, Point3D{X: x, Y: y, Z: z})
// 					currStartLengthPos = Point3D{X: x, Y: y, Z: z}
// 					nx, ny, nz := 0.0, 0.0, 1.0
// 					cosTheta := -(dx*nx + dy*ny + dz*nz)
// 					theta := math.Acos(cosTheta)
// 					currReflectionFactor *= calculateReflectionFactor(theta, "concrete")
// 					continue
// 				} 
// 				 if index == rl.Config.CornerMapNumber && currWallIndex != rl.Config.CornerMapNumber {
// 					currWallIndex = rl.Config.CornerMapNumber
// 					currStartLengthPos = Point3D{X: x, Y: y, Z: z}
// 					normals := getNeighborWallNormals(xIdx, yIdx, zIdx, rl)
// 					fmt.Printf("xIdx: %v yIdx: %v zIdx: %v \n", xIdx, yIdx, zIdx  )
// 					if len(normals) == 0 {
// 						break
// 					}
// 					bestDot := -math.MaxFloat64
// 					var bestNormal Normal3D
					
// 					for k, n := range normals {
					
				
// 						dot := dx*n.Nx + dy*n.Ny + dz*n.Nz

				
// 						n.Nx, n.Ny, n.Nz = -n.Nx, -n.Ny, -n.Nz
// 						dot = -dot
					
					
// 						if dot > bestDot {
// 							bestDot = dot
// 							bestNormal = n
// 						}

// 						fmt.Printf("Promien %d, %d Normalna %d: Nx=%.3f, Ny=%.3f, Nz=%.3f Dot=%.3f\n", i, j, k, n.Nx, n.Ny, n.Nz, dot)
// 					}
// 					fmt.Printf("Best Dot: %.3f, bestNormal: %v \n", bestDot, bestNormal)
// 					cosTheta := dx*bestNormal.Nx + dy*bestNormal.Ny + dz*bestNormal.Nz
// 					if cosTheta > 1 {
// 						cosTheta = 1
// 					}
// 					if cosTheta < -1 {
// 						cosTheta = -1
// 					}
// 					theta := math.Acos(cosTheta)
// 					finalTheta := math.Pi/2 - theta
// 					stepResolution := 10
// 					oneStep := finalTheta/float64(stepResolution)
					
// 					fmt.Printf("cosTheta: %.3f theta: %.3f, finalTheta: %.3f \n", cosTheta, theta, finalTheta)

					
// 					for step := 0; step < stepResolution; step++ {
// 						x := x
// 						y := y
// 						z := z
// 						angle := float64(step) * oneStep
// 						newDx := dx * math.Cos(angle) - dy * math.Sin(angle)
// 						newDy := dx * math.Sin(angle) + dy * math.Cos(angle)
// 						fmt.Printf("Step %d: dx=%.3f dy=%.3f dz=%.3f\n", step, newDx, newDy, dz)
// 						x += newDx
// 						y += newDy
// 						z += dz
						
// 						xIdx := int(math.Round(x/rl.Config.Step))
// 						yIdx := int(math.Round(y/rl.Config.Step))
// 						zIdx := int(math.Round(z/rl.Config.Step))
// 						index := int(rl.PowerMap[zIdx][yIdx][xIdx])
// 						if (index >= rl.Config.WallMapNumber) {
// 							x += 2*newDx
// 							y += 2*newDy
// 							z += 2*dz
// 						}
// 						fmt.Printf("newDx: %v newDy: %v z: %v \n", x , y, z  )
// 						for (x >= 0 && x <= rl.Config.SizeX) && (y >= 0 && y < rl.Config.SizeY) && (z <= rl.Config.SizeZ) && currInteractions < rl.Config.NumOfInteractions && currPower >= rl.Config.MinimalRayPower {
// 							xIdx := int(math.Round(x/rl.Config.Step))
// 							yIdx := int(math.Round(y/rl.Config.Step))
// 							zIdx := int(math.Round(z/rl.Config.Step))
// 							index := int(rl.PowerMap[zIdx][yIdx][xIdx])
// 							fmt.Printf("xIdx: %v yIdx: %v zIdx: %v \n", xIdx, yIdx, zIdx  )
// 							if index >= rl.Config.WallMapNumber && index < rl.Config.RoofMapNumber && index != currWallIndex + rl.Config.WallMapNumber{ 	
// 								// 	check if there is wall and if its diffrent from previous one
// 								currWallIndex = index - rl.Config.WallMapNumber

// 								//get wall normal
// 								nx, ny, nz := rl.WallNormals[currWallIndex].Nx, rl.WallNormals[currWallIndex].Ny, rl.WallNormals[currWallIndex].Nz
// 								//!!! MAP IS MIRRORED BY Y SO ALL Y NORMALS SHOULD BE MIRRORED !!!
// 								ny = -ny
// 								dot := dx*nx + dy*ny + dz*nz
// 								if dot > 0 {
// 									nx, ny, nz = -nx, -ny, -nz
// 									dot = -dot
// 								}
// 								dot *= 2
// 								fmt.Printf("Promien %d, %d Nx=%.3f, Ny=%.3f, Nz=%.3f dot:%.3f \n", i, j,  nx, ny, nz, dot)

// 								dx = dx - dot*nx
// 								dy = dy - dot*ny
// 								dz = dz - dot*nz
// 								currInteractions++

// 								// sum distance and set new start position
// 								currSumRayLength += calculateDistance(currStartLengthPos, Point3D{X: x, Y: y, Z: z})
// 								currStartLengthPos = Point3D{X: x, Y: y, Z: z}
// 								cosTheta := -(dx*nx + dy*ny + dz*nz)
// 								if cosTheta > 1 {
// 									cosTheta = 1
// 								}
// 								if cosTheta < -1 {
// 									cosTheta = -1
// 								}
// 								theta := math.Acos(cosTheta)
// 								currReflectionFactor *= calculateReflectionFactor(theta, "concrete")
// 							} else {
// 								currRayLength = calculateDistance(currStartLengthPos, Point3D{X: x, Y: y, Z: z}) + currSumRayLength
// 								H := calculateTransmittance(currRayLength, rl.Config.WaveLength, currReflectionFactor)
// 								currPower = 10*math.Log10(rl.Config.TransmitterPower) + 20*math.Log10(cmplx.Abs(H))
// 								if rl.PowerMap[zIdx][yIdx][xIdx] == -150 || rl.PowerMap[zIdx][yIdx][xIdx] < currPower {
// 									rl.PowerMap[zIdx][yIdx][xIdx] = currPower
// 								} 
// 							}
// 							x += newDx
// 							y += newDy
// 							z += dz
// 						}
// 					}
	
// 				} 

// 				if index >= rl.Config.WallMapNumber && index < rl.Config.RoofMapNumber && index != currWallIndex + rl.Config.WallMapNumber{ 	// check if there is wall and if its diffrent from previous one
// 					currWallIndex = index - rl.Config.WallMapNumber

// 					//get wall normal
// 					nx, ny, nz := rl.WallNormals[currWallIndex].Nx, rl.WallNormals[currWallIndex].Ny, rl.WallNormals[currWallIndex].Nz
// 					//!!! MAP IS MIRRORED BY Y SO ALL Y NORMALS SHOULD BE MIRRORED !!!
// 					ny = -ny
// 					dot := dx*nx + dy*ny + dz*nz
// 					if dot > 0 {
// 						nx, ny, nz = -nx, -ny, -nz
// 						dot = -dot
// 					}
// 					dot *= 2
// 					fmt.Printf("Promien %d, %d Nx=%.3f, Ny=%.3f, Nz=%.3f dot:%.3f \n", i, j,  nx, ny, nz, dot)

// 					dx = dx - dot*nx
// 					dy = dy - dot*ny
// 					dz = dz - dot*nz
// 					currInteractions++

// 					// sum distance and set new start position
// 					currSumRayLength += calculateDistance(currStartLengthPos, Point3D{X: x, Y: y, Z: z})
// 					currStartLengthPos = Point3D{X: x, Y: y, Z: z}
// 					cosTheta := -(dx*nx + dy*ny + dz*nz)
// 					if cosTheta > 1 {
// 						cosTheta = 1
// 					}
// 					if cosTheta < -1 {
// 						cosTheta = -1
// 					}
// 					theta := math.Acos(cosTheta)
// 					currReflectionFactor *= calculateReflectionFactor(theta, "concrete")
// 				} else {
// 					// calculate distance and transmittance
// 					currRayLength = calculateDistance(currStartLengthPos, Point3D{X: x, Y: y, Z: z}) + currSumRayLength

// 					H := calculateTransmittance(currRayLength, rl.Config.WaveLength, currReflectionFactor)
// 					currPower = 10*math.Log10(rl.Config.TransmitterPower) + 20*math.Log10(cmplx.Abs(H)) - diffLossLdB
// 					if (diffLossLdB > 0.0) {
// 						// println("diffLoss: ",diffLossLdB)
// 						// println("currPorwe: ",currPower)
// 					}
// 					// update power map if power is higher than previous one
// 					if rl.PowerMap[zIdx][yIdx][xIdx] == -150 || rl.PowerMap[zIdx][yIdx][xIdx] < currPower {
// 						rl.PowerMap[zIdx][yIdx][xIdx] = currPower
// 					} 
// 					if targetRayIndex >= 0 {
// 						rl.RayPaths[targetRayIndex] = append(rl.RayPaths[targetRayIndex], RayPoint{
// 							X: float64(math.Round(x/rl.Config.Step)), 
// 							Y: float64(math.Round(y/rl.Config.Step)), 
// 							Z: float64(math.Round(z/rl.Config.Step)), 
// 							Power: currPower,
// 						})
// 					}
// 				}
// 				// println("currReflectionFactor: ",currReflectionFactor)
// 				// update position
// 				x += dx
// 				y += dy
// 				z += dz
// 			}
// 		}
// 	}
// }

func calculateDistance(p1, p2 Point3D) float64 {
	dist := math.Sqrt(math.Pow(p1.X-p2.X,2)+math.Pow(p1.Y-p2.Y,2)+math.Pow(p1.Z-p2.Z,2))
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
		if i - singleRay.Azimuth == 0 && j - singleRay.Elevation == 0{
			return idx
		}
	}
	return -1 
}

func calculateReflectionFactor(angle float64, material string) float64 {
	if angle > math.Pi/2 {
		angle = math.Pi - angle
	}
	var eta float64;
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
	root := math.Sqrt(eta - sinTheta*sinTheta)
	R_TE := (cosTheta - root)/(cosTheta + root)
	R_TM := (eta*cosTheta - root)/(eta*cosTheta + root)
	reflectionFactor := (math.Pow(R_TE, 2) + math.Pow(R_TM, 2)) / 2
	println("ANGLE: ", angle, "root: ", root,"R_TE: ", R_TE, " R_TM: ", R_TM, " reflectionFactor ", reflectionFactor)
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

	size :=3
	for dx := -size; dx <= size; dx++ {
		for dy := -size; dy <= size; dy++ {
			for dz := -size; dz <= size; dz++ {
				xprim := x + dx
				yprim := y + dy
				zprim := z + dz

				if xprim < 0 || yprim < 0 || zprim < 0 || xprim >= int(rl.Config.SizeX )|| yprim >= int(rl.Config.SizeY) || zprim >= int(rl.Config.SizeZ) {
					continue
				}

				index := int(rl.PowerMap[zprim][yprim][xprim])
				if index >= rl.Config.WallMapNumber && index < rl.Config.RoofMapNumber {
					currWallIndex := index - rl.Config.WallMapNumber
					if _, exists := neighborNormals[currWallIndex]; !exists {
						neighborNormals[currWallIndex] = rl.WallNormals[currWallIndex]
					}
				}
				if index == rl.Config.RoofMapNumber {
					neighborNormals[index] = Normal3D{Nx: 0, Ny: 0, Nz: 1}
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

func BergDiffractionLoss(d1, d2, lambda, alpha float64) float64 {
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

