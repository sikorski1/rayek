import numpy as np
from math import sqrt, log10,pi,e
import matplotlib.pyplot as plt
from Vector import Vector
class Raytracing:
    def __init__(self, matrixDimensions, tPos, tPower, tFreq, rFactor, oPos):
        self.step = 0.1
        self.transmitterPos = tPos
        self.transmitterPower = tPower #mW
        self.transmitterFreq = tFreq # GHz
        self.waveLength = 299792458 / tFreq / 10 ** 9;
        self.reflectionFactor = rFactor
        self.obstaclesPos = oPos 
        self.powerMap = np.zeros((matrixDimensions[1]*int((1/self.step))+1, matrixDimensions[0]*int((1/self.step))+1))
        self.matrix = self.createMatrix(matrixDimensions)
        self.mirroredTransmittersPos = self.createMirroredTransmitters()
        self.mirroredWalls = self.calculateMirroredWalls(self.obstaclesPos, self.obstaclesPos)
    def createMatrix(self, matrixDimensions):
        x = np.linspace(0, matrixDimensions[0], int(matrixDimensions[0]/self.step)+1)
        y = np.linspace(0, matrixDimensions[1], int(matrixDimensions[1]/self.step)+1)

        xv, yv = np.meshgrid(x, y)
        matrix = np.stack((xv, yv), axis=-1)
        return matrix

    def twoVectors(self, A, B, C, D):
        result = (((C[0] - A[0])*(B[1] - A[1]) - (B[0] - A[0])*(C[1] - A[1])) * ((D[0] - A[0]) * (B[1] - A[1]) - (B[0] - A[0]) * (D[1] - A[1])))  #[C[0]-A-x, B[0]-A[0], vector CA, BA
                                                                                                                                   #C[1]-A-y, B[1]-A[1]]
        if result > 0:
            return -1
        else:
            result2 = (((A[0] - C[0])*(D[1] - C[1]) - (D[0] - C[0])*(A[1] - C[1])) * ((B[0] - C[0]) * (D[1] - C[1]) - (D[0] - C[0]) * (B[1] - C[1])))
            if result2 > 0:
                return -1
            elif result < 0 and result2 < 0:
                return 1
            elif result == 0 and result2 < 0:
                return 0
            elif result < 0 and result2 == 0:
                return 0
            elif A[0] < C[0] and A[0] < D[0] and B[0] < C[0] and B[0] < D[0]:
                return -1
            elif A[1] < C[1] and A[1] < D[1] and B[1] < C[1] and B[1] < D[1]:
                return -1
            elif A[0] > C[0] and A[0] > D[0] and B[0] > C[0] and B[0] > D[0]:
                return -1
            elif A[1] > C[1] and A[1] > D[1] and B[1] > C[1] and B[1] > D[1]:
                return -1
            else:
                return 0
    def calculateRayTracing(self):
        for i in range(len(self.matrix)):
            for j in range(len(self.matrix[0])):
                H = 0
                receiverPos = self.matrix[i][j]
                if self.checkLineOfSight(receiverPos, self.obstaclesPos):
                    H += self.calculateTransmitation(receiverPos, self.transmitterPos) # add transmitantion from line of sight
                    
                H += self.calculateSingleWallReflection(receiverPos, self.obstaclesPos)
                if H == 0:
                    self.powerMap[i][j] = -150
                else:
                    self.powerMap[i][j] = 10*log10(self.transmitterPower) + 20*log10(abs(H))
    
    def checkLineOfSight(self, receiverPos, walls):
        for wall in walls:
            if self.twoVectors(receiverPos, self.transmitterPos, wall.A, wall.B) >= 0:  #checking collision with wall line of sight
                return False
        return True
    def calculateSingleWallReflection(self, receiverPos, walls):
        H = 0
        for i, wall in enumerate(walls): # main wall, vectors should colliding
            if self.twoVectors(receiverPos, self.mirroredTransmittersPos[i], wall.A, wall.B) <= 0:
                continue
            for j in range(len(walls) - 1): # all other walls and their reflection
                index = (i + j + 1) % len(walls) # other walls index, start at i + 1 and add j, if there is last index start at the beginning and loop over indexes before i
                if self.twoVectors(self.transmitterPos, walls[index].A, wall.A, wall.B) <= 0 and self.twoVectors(self.transmitterPos, walls[index].B, wall.A, wall.B) <= 0:
                    if self.twoVectors(receiverPos, self.mirroredTransmittersPos[i], walls[index].A, walls[index].B) >= 0:
                        break
                    if self.mirroredWalls[i][index] != None and self.twoVectors(receiverPos, self.mirroredTransmittersPos[i], self.mirroredWalls[i][index].A, self.mirroredWalls[i][index].B) >= 0:
                        break
            else:
                H += self.calculateTransmitation(receiverPos, self.mirroredTransmittersPos[i], self.reflectionFactor)
        return H

    def calculateMirroredWalls(self, mainWalls, wallsThatAreMirrored):
        mirroredWalls = np.zeros((len(mainWalls), len(wallsThatAreMirrored)), dtype=Vector)
        for i, mainWall in enumerate(mainWalls):
            for j, wallThatIsMirrored in enumerate(wallsThatAreMirrored):
                if mainWall.B[0] == mainWall.A[0]: # if main wall is vertical
                    isTransmitterLeft = self.transmitterPos[0] <= mainWall.A[0] 
                    isWallMirroredLeft = wallThatIsMirrored.A[0] <= mainWall.A[0] and wallThatIsMirrored.B[0] <= mainWall.A[0]
                    isTrue = isTransmitterLeft == isWallMirroredLeft
                    if isTrue:         
                        aPrimX = 2 * mainWall.A[0] - wallThatIsMirrored.A[0]
                        bPrimX = 2 * mainWall.A[0] - wallThatIsMirrored.B[0]
                        mirroredWalls[i][j] = Vector([aPrimX, wallThatIsMirrored.A[1]], [bPrimX, wallThatIsMirrored.B[1]])
                    else:
                        mirroredWalls[i][j] = None
                    continue
                if mainWall.B[1] == mainWall.A[1]: # if main wall is horizontal
                    isTransmitterBelow = self.transmitterPos[1] <= mainWall.A[1] 
                    isWallMirroredBelow = wallThatIsMirrored.A[1] <= mainWall.A[1] and wallThatIsMirrored.B[1] <= mainWall.A[1]
                    isTrue = isTransmitterBelow == isWallMirroredBelow
                    if isTrue:
                        aPrimY = 2 * mainWall.A[1]  - wallThatIsMirrored.A[1]
                        bPrimY = 2 * mainWall.A[1]  - wallThatIsMirrored.B[1]
                        mirroredWalls[i][j] = Vector([wallThatIsMirrored.A[0], aPrimY], [wallThatIsMirrored.B[0], bPrimY])
                    else:
                        mirroredWalls[i][j] = None
                    continue
                # if main wall is oblique
                if self.twoVectors(self.transmitterPos, wallThatIsMirrored.A, mainWall.A, mainWall.B) <= 0 and self.twoVectors(self.transmitterPos, wallThatIsMirrored.B, mainWall.A, mainWall.B) <= 0:
                    m = (mainWall.B[1] - mainWall.A[1])/(mainWall.B[0] - mainWall.A[0])
                    b =  mainWall.A[1] - m * mainWall.A[0]
                    m2 = -1/m
                    b2 = wallThatIsMirrored.A[1] - m2 * wallThatIsMirrored.A[0]
                    x1 = (b2-b)/(m-m2)
                    y1 = m*x1 + b
                    m = (mainWall.B[1] - mainWall.A[1])/(mainWall.B[0] - mainWall.A[0])
                    b =  mainWall.A[1] - m * mainWall.A[0]
                    m2 = -1/m
                    b2 = wallThatIsMirrored.B[1] - m2 * wallThatIsMirrored.B[0]
                    x2 = (b2-b)/(m-m2)
                    y2 = m*x2 + b
                    aPrimX = 2*x1 - wallThatIsMirrored.A[0]
                    aPrimY = 2*y1 - wallThatIsMirrored.A[1]
                    bPrimX = 2*x2 - wallThatIsMirrored.B[0]
                    bPrimY = 2*y2 - wallThatIsMirrored.B[1]
                    mirroredWalls[i][j] = Vector([aPrimX, aPrimY], [bPrimX, bPrimY])
                else:
                    mirroredWalls[i][j] = None
        return mirroredWalls

    def createMirroredTransmitters(self):
        mirroredTransmittersPos = np.zeros((len(self.obstaclesPos), 2))
        for i in range(len(self.obstaclesPos)):
            wall = self.obstaclesPos[i]
            if wall.A[0] == wall.B[0]: 
                mirroredTransmittersPos[i][1] = self.transmitterPos[1]
                distance = abs(wall.A[0] - self.transmitterPos[0])
                if self.transmitterPos[0] < wall.A[0]:
                    mirroredTransmittersPos[i][0] = wall.A[0] + distance
                else:
                    mirroredTransmittersPos[i][0] = wall.A[0] - distance
                continue

            if wall.A[1] == wall.B[1]: 
                mirroredTransmittersPos[i][0] = self.transmitterPos[0]
                distance = abs(wall.A[1] - self.transmitterPos[1])
                if self.transmitterPos[1] < wall.A[1]:
                    mirroredTransmittersPos[i][1] = wall.A[1] + distance
                else:
                    mirroredTransmittersPos[i][1] = wall.A[1] - distance
                continue
            m = (wall.B[1] - wall.A[1]) / (wall.B[0] - wall.A[0])
            b = wall.A[1] - m * wall.A[0]
            m2 = -1/m
            b2 = self.transmitterPos[1] - m2 * self.transmitterPos[0]
            x = (b2-b)/(m-m2)
            y = m*x + b
    
            mirroredTransmittersPos[i][0] = 2*x - self.transmitterPos[0]
            mirroredTransmittersPos[i][1] = 2*y - self.transmitterPos[1] 
            print(mirroredTransmittersPos[i][0], mirroredTransmittersPos[i][1])
        return mirroredTransmittersPos
            
    def calculateTransmitation(self, p1, p2, reflectionRef=1):
        r = self.calculateDist(p1, p2)
        if r > 0:
            H = reflectionRef*self.waveLength/(4*pi*r)*e**(-2j*pi*r/self.waveLength)
        else:
            H = 0
        return H
    def calculateDist(self, p1, p2):
        dist = sqrt((p1[0] - p2[0])**2 + (p1[1] - p2[1])**2)
        return dist    
    def displayPowerMap(self):
        plt.figure(figsize=(10, 8))
        plt.imshow(self.powerMap, origin='lower', cmap='jet', 
                extent=[0, self.matrix.shape[1]*self.step, 0, self.matrix.shape[0]*self.step])
        plt.colorbar(label='Power (dBm)')
        plt.title('Power Map')
        plt.xlabel('X Coordinate (m)')
        plt.ylabel('Y Coordinate (m)')

        # Plot walls
        for wall in self.obstaclesPos:
            x_coords = [wall.A[0], wall.B[0]]
            y_coords = [wall.A[1], wall.B[1]]
            plt.plot(x_coords, y_coords, color='black', linewidth=2, label='Wall')

   
        colors = ["green", "red", "blue", "yellow", "black", "purple", "orange", "grey"]
        for i, mirroredWallsForWall in enumerate(self.mirroredWalls):
            for mirroredWall in mirroredWallsForWall:
                if mirroredWall:
                    x_coords = [mirroredWall.A[0], mirroredWall.B[0]]
                    y_coords = [mirroredWall.A[1], mirroredWall.B[1]]
                    plt.plot(x_coords, y_coords, color=colors[i], linewidth=1, linestyle='--', label=f'Mirrored Wall {i+1}')

        # Plot transmitter
        plt.scatter(self.transmitterPos[0], self.transmitterPos[1], color='red', label='Transmitter', zorder=5)

        # Plot mirrored transmitters
        for idx, mirroredPos in enumerate(self.mirroredTransmittersPos):
            plt.scatter(mirroredPos[0], mirroredPos[1], color=colors[idx], label=f'Mirrored Transmitter {idx+1}', zorder=5)

        # Adjust legend to avoid duplicate entries
        handles, labels = plt.gca().get_legend_handles_labels()
        by_label = dict(zip(labels, handles))
        plt.legend(by_label.values(), by_label.keys())

        plt.show()

wall1 = Vector([5, 16],[8, 16])
wall2 = Vector([12, 16],[16, 16])
wall3 = Vector([1, 14],[1, 17])
wall4 = Vector([16, 1],[16, 5])
wall5 = Vector([16, 8],[16, 12])
wall6 = Vector([8, 8],[10, 12])
wall7 = Vector([4, 16],[8, 22])
wall8 = Vector([12, 20],[16, 24])

        
raytracing = Raytracing([30, 30], [7, 11], 5, 3.6, 0.7, [wall1, wall2, wall3, wall4,wall5, wall6, wall7, wall8])
# print(raytracing.calculateSingleWallReflection([9.9, 19.9], [wall1, wall2]))
# print(raytracing.checkLineOfSight([9.9, 19.9], [wall1, wall2]))
# print(raytracing.calculateSingleWallReflection([8.4, 11], [wall1, wall2, wall3]))

raytracing.calculateRayTracing()
raytracing.displayPowerMap()