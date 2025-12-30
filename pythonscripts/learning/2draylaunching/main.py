import numpy as np
import json
import matplotlib.pyplot as plt
from matplotlib.colors import LinearSegmentedColormap
import time

class Point:
    def __init__(self, x, y):
        self.X = x
        self.Y = y
    
    def __repr__(self):
        return f"Point(X={self.X}, Y={self.Y})"

class Vector:
    def __init__(self, A, B):
        if isinstance(A, list):
            self.A = Point(A[0], A[1])
        else:
            self.A = A
        if isinstance(B, list):
            self.B = Point(B[0], B[1])
        else:
            self.B = B
    
    def __repr__(self):
        return f"Vector(A={self.A}, B={self.B})"

class Normal:
    def __init__(self, nx, ny):
        self.Nx = nx
        self.Ny = ny
    
    def __repr__(self):
        return f"Normal(Nx={self.Nx}, Ny={self.Ny})"

class RayLaunching:
    def __init__(self, matrixDimensions, tPos, tPower, tFreq, rFactor, wallPos):
        self.step = 0.1
        self.numberOfRays = 2880
        self.numberOfInteractions = 2
        self.minimalPower = -150.0
        self.wallMapNumber = 1000
        self.transmitterPos = tPos
        self.transmitterPower = tPower
        self.transmitterFreq = tFreq
        self.waveLength = 299792458 / (tFreq * 10**9)
        self.reflectionFactor = rFactor
        
        rows = int(matrixDimensions.Y * (1/self.step)) + 1
        cols = int(matrixDimensions.X * (1/self.step)) + 1
        
        self.Map = np.full((rows, cols), -150.0)
        self.wallNormals = [Normal(0, 0) for _ in range(len(wallPos))]
        
        self.setWallsIn2DMap(wallPos)
    
    def calculateDist(self, p1, p2):
        return np.sqrt((p1.X - p2.X)**2 + (p1.Y - p2.Y)**2)
    
    def calculateTransmittanceWithLength(self, r, waveLength, reflectionFactor):
        if r > 0:
            H = reflectionFactor * waveLength / (4 * np.pi * r) * np.exp(-2j * np.pi * r / waveLength)
        else:
            H = 0
        return H
    
    def calculateRayLaunching(self):
        maxSizeX = (len(self.Map[0]) - 1) * self.step
        maxSizeY = (len(self.Map) - 1) * self.step
        
        for i in range(self.numberOfRays):
            currInteractions = 0
            currPower = 0.0
            dRadians = (2 * np.pi) / self.numberOfRays * i
            dx = np.cos(dRadians) * self.step
            dy = np.sin(dRadians) * self.step
            dx = np.round(dx * 1e15) / 1e15
            dy = np.round(dy * 1e15) / 1e15
            
            x = self.transmitterPos.X + dx
            y = self.transmitterPos.Y + dy
            currWallIndex = 0
            currRayLength = 0.0
            sumRayLength = 0.0
            currStartLengthPos = Point(self.transmitterPos.X, self.transmitterPos.Y)
            
            while ((x >= 0 and x <= maxSizeX) and (y >= 0 and y <= maxSizeY) and 
                   currInteractions < self.numberOfInteractions and currPower >= self.minimalPower):
                
                xIdx = int(np.round(x / self.step))
                yIdx = int(np.round(y / self.step))
                
                # Sprawdzenie granic
                if yIdx >= len(self.Map) or xIdx >= len(self.Map[0]):
                    break
                
                index = int(self.Map[yIdx][xIdx])
                
                # Sprawdzenie czy jest ściana i czy różni się od poprzedniej
                if index >= self.wallMapNumber and index != currWallIndex + self.wallMapNumber:
                    currWallIndex = index - self.wallMapNumber
                    nx = self.wallNormals[currWallIndex].Nx
                    ny = self.wallNormals[currWallIndex].Ny
                    dot = 2 * (dx * nx + dy * ny)
                    dx = dx - dot * nx
                    dy = dy - dot * ny
                    currInteractions += 1
                    sumRayLength += self.calculateDist(currStartLengthPos, Point(x, y))
                    currStartLengthPos = Point(x, y)
                else:
                    currRayLength = self.calculateDist(currStartLengthPos, Point(x, y)) + sumRayLength
                    H = self.calculateTransmittanceWithLength(
                        currRayLength, 
                        self.waveLength, 
                        self.reflectionFactor ** currInteractions
                    )
                    currPower = 10 * np.log10(self.transmitterPower) + 20 * np.log10(abs(H))
                    
                    if self.Map[yIdx][xIdx] != -150:
                        existingPowerLin = 10 ** (self.Map[yIdx][xIdx] / 10)
                        currPowerLin = 10 ** (currPower / 10)
                        newPowerDb = 10 * np.log10(existingPowerLin + currPowerLin)
                        self.Map[yIdx][xIdx] = newPowerDb
                    else:
                        self.Map[yIdx][xIdx] = currPower
                
                x += dx
                y += dy
    
    def setWallsIn2DMap(self, walls):
        for i, wall in enumerate(walls):
            x1, y1 = wall.A.X, wall.A.Y
            x2, y2 = wall.B.X, wall.B.Y
            x1Idx = int(np.round(x1 / self.step))
            y1Idx = int(np.round(y1 / self.step))
            x2Idx = int(np.round(x2 / self.step))
            y2Idx = int(np.round(y2 / self.step))
            
            dx = x2 - x1
            dy = y2 - y1
            length = np.hypot(dx, dy)
            
            if length != 0:
                nx = -dy / length
                ny = dx / length
                self.wallNormals[i] = Normal(nx, ny)
            
            if x1 == x2:
                if y1 > y2:
                    y1Idx, y2Idx = y2Idx, y1Idx
                for y in range(y1Idx, y2Idx + 1):
                    if y < len(self.Map) and x1Idx < len(self.Map[0]):
                        self.Map[y][x1Idx] = self.wallMapNumber + i
            elif y1 == y2:
                if x1 > x2:
                    x1Idx, x2Idx = x2Idx, x1Idx
                for x in range(x1Idx, x2Idx + 1):
                    if y1Idx < len(self.Map) and x < len(self.Map[0]):
                        self.Map[y1Idx][x] = self.wallMapNumber + i
            else:
                steps = int(max(abs(dx / self.step), abs(dy / self.step)))
                prevXIdx = int(np.round(x1 / self.step))
                prevYIdx = int(np.round(y1 / self.step))
                
                for j in range(steps + 1):
                    x = x1 + (dx * j) / steps
                    y = y1 + (dy * j) / steps
                    xIdx = int(np.round(x / self.step))
                    yIdx = int(np.round(y / self.step))
                    
                    if yIdx >= len(self.Map) or xIdx >= len(self.Map[0]):
                        continue
                    
                    # Zapewnienie ciągłości ścian
                    if ((prevXIdx < xIdx and prevYIdx < yIdx) or 
                        (prevXIdx < xIdx and prevYIdx > yIdx)):
                        if yIdx < len(self.Map) and prevXIdx < len(self.Map[0]):
                            self.Map[yIdx][prevXIdx] = self.wallMapNumber + i
                    
                    if ((prevXIdx > xIdx and prevYIdx < yIdx) or 
                        (prevXIdx > xIdx and prevYIdx > yIdx)):
                        if prevYIdx < len(self.Map) and xIdx < len(self.Map[0]):
                            self.Map[prevYIdx][xIdx] = self.wallMapNumber + i
                    
                    self.Map[yIdx][xIdx] = self.wallMapNumber + i
                    prevXIdx = xIdx
                    prevYIdx = yIdx
    
    def displayPowerMap(self, walls):
        plt.figure(figsize=(10, 8))
        plt.imshow(self.Map, origin='lower', cmap='jet',
                   extent=[0, self.Map.shape[1] * self.step, 0, self.Map.shape[0] * self.step],
                   vmin=-150, vmax=10)
        plt.colorbar(label='Power (dBm)')
        plt.title('Ray Launching Power Map')
        plt.xlabel('X Coordinate (m)')
        plt.ylabel('Y Coordinate (m)')
        
        # Rysowanie ścian
        for i, wall in enumerate(walls):
            x_coords = [wall.A.X, wall.B.X]
            y_coords = [wall.A.Y, wall.B.Y]
            label_text = "Wall" if i == 0 else None
            plt.plot(x_coords, y_coords, color='red', linewidth=2, label=label_text)
        
        # Transmiter
        plt.scatter(self.transmitterPos.X, self.transmitterPos.Y, 
                    color='white', edgecolor='black', s=100, label='Transmitter', zorder=5)
        
        plt.legend(loc='upper right')
        plt.show()

def saveMapToCSV(Map, filename):
    np.savetxt(filename, Map, delimiter=',', fmt='%f')

def saveConfigToJSON(transmitter, walls, filename):
    wallsData = []
    for wall in walls:
        wallsData.append([wall.A.X, wall.A.Y, wall.B.X, wall.B.Y])
    
    config = {
        'tx': transmitter.X,
        'ty': transmitter.Y,
        'walls': wallsData
    }
    
    with open(filename, 'w') as f:
        json.dump(config, f, indent=2)

# ---------------------------------------------------------
# --- MAIN ---
# ---------------------------------------------------------

if __name__ == "__main__":
    start = time.time()
    
    matrixDimensions = Point(40, 40)
    transmitterPos = Point(20, 20)
    transmitterPower = 5.0
    transmitterFreq = 2.4
    reflectionFactor = 0.8
    
    walls1 = [
        Vector(Point(0, 3), Point(3, 6)), Vector(Point(1, 3), Point(6, 3)),
        Vector(Point(12, 12), Point(6, 10)), Vector(Point(25, 10), Point(25, 15)),
        Vector(Point(10, 36), Point(5, 30)), Vector(Point(23, 36), Point(25, 39)),
        Vector(Point(1, 24), Point(1, 26)), Vector(Point(1, 28), Point(1, 30)),
        Vector(Point(1, 37), Point(1, 40)), Vector(Point(35, 36), Point(30, 28)),
        Vector(Point(40, 1), Point(36, 2)), Vector(Point(24, 3), Point(25, 6)),
        Vector(Point(16, 21), Point(18, 22)), Vector(Point(12, 18), Point(12, 20)),
        Vector(Point(18, 36), Point(12, 36))
    ]
    
    walls_scenario_1 =walls_list = [
        Vector([5, 5], [5, 15]),   Vector([5, 15], [15, 15]),  # Lewy dolny róg
        Vector([10, 20], [10, 30]),Vector([10, 30], [20, 30]), # Środkowa sekcja
        Vector([25, 5], [35, 15]), Vector([35, 15], [35, 5]),  # Trójkąt
        Vector([25, 25], [35, 35]),Vector([35, 25], [25, 35]), # X kształt
        Vector([2, 38], [15, 38]), Vector([38, 2], [38, 15]),  # Graniczne bariery
        Vector([18, 14], [22, 18]), # Mała przeszkoda przy środku
        Vector([0, 20], [5, 25]),   # Ukośna wejściowa
        Vector([30, 0], [35, 5]),   # Ukośna dolna
        Vector([15, 5], [20, 10]),  # Rozpraszacz
        Vector([20, 35], [25, 40])  # Górna
    ]
    
    walls_scenario_2 = [
        # Zewnętrzny pierścień (przerywany)
        Vector([2, 2], [2, 10]), Vector([2, 12], [2, 20]), Vector([2, 22], [2, 38]),
        Vector([38, 2], [38, 15]), Vector([38, 18], [38, 38]),
        Vector([5, 38], [15, 38]), Vector([20, 38], [35, 38]),
        Vector([5, 2], [20, 2]), Vector([25, 2], [35, 2]),
        
        # Wewnętrzne struktury
        Vector([8, 8], [12, 12]), Vector([12, 12], [16, 8]),   # Zygzak
        Vector([8, 32], [12, 28]), Vector([12, 28], [16, 32]), # Zygzak góra
        
        # Centralne przeszkody (wokół transmitera 20,20)
        Vector([18, 22], [22, 22]), Vector([22, 22], [22, 18]), # Kąt przy środku
        Vector([18, 18], [15, 15]), Vector([22, 25], [25, 28]), # Rozchodzące się
        
        # Rozproszone ukośne (reflektory)
        Vector([30, 10], [35, 15]), Vector([30, 15], [35, 10]),
        Vector([30, 25], [35, 30]), Vector([30, 30], [35, 25]),
        
        # Pionowe/Poziome przegrody
        Vector([10, 15], [10, 25]),
        Vector([28, 5], [28, 15]),
        Vector([28, 25], [28, 35]),
        Vector([15, 5], [25, 5]),
        
        # Dodatkowe drobne elementy
        Vector([5, 25], [8, 25]),
        Vector([32, 20], [35, 20]),
        Vector([12, 35], [15, 32]),
        Vector([25, 8], [22, 5])
    ]
    
    wallsSet = [walls1, walls_scenario_1, walls_scenario_2]
    
    # WYBIERZ SCENARIUSZ
    walls = wallsSet[1]  # 0, 1, lub 2
    
    raylaunching = RayLaunching(matrixDimensions, transmitterPos, 
                                transmitterPower, transmitterFreq, 
                                reflectionFactor, walls)
    
    print(f"Rozpoczynam obliczenia Ray Launching dla {len(walls)} ścian...")
    raylaunching.calculateRayLaunching()
    
    stop = time.time() - start
    print(f"Computation time: {stop:.2f}s")
    
    saveMapToCSV(raylaunching.Map, "ray_map.csv")
    saveConfigToJSON(transmitterPos, walls, "ray_config.json")
    
    raylaunching.displayPowerMap(walls)
    
    print("Done! Files saved: ray_map.csv, ray_config.json, and plot displayed.")