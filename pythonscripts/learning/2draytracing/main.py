import numpy as np
from math import sqrt, log10,pi,e
import matplotlib.pyplot as plt
from Vector import Vector
import time
class Raytracing:
    def __init__(self, matrixDimensions, tPos, tPower, tFreq, rFactor, oPos):
        self.step = 0.1
        self.transmitterPos = tPos
        self.transmitterPower = tPower #mW
        self.transmitterFreq = tFreq # GHz
        self.waveLength = 299792458 / tFreq / 10 ** 9;
        self.reflectionFactor = rFactor
        self.walls = oPos 
        self.powerMap = np.zeros((int(matrixDimensions[1]*(1/self.step)+1), int(matrixDimensions[0]*(1/self.step))+1))
        self.matrix = self.createMatrix(matrixDimensions)
        self.mirroredTransmittersPos = self.createMirroredTransmitters(self.walls)
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
                if self.checkLineOfSight(receiverPos, self.walls):
                    H += self.calculateTransmitation(receiverPos, self.transmitterPos) # add transmitantion from line of sight
                    
                H += self.calculateSingleWallReflection(receiverPos, self.walls)
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
            else:
                reflectionPoint = self.calculateCrossPoint(receiverPos, self.mirroredTransmittersPos[i], wall.A, wall.B)
            for j in range(len(walls) - 1): # all other walls and their reflection
                index = (i + j + 1) % len(walls) # other walls index, start at i + 1 and add j, if there is last index start at the beginning and loop over indexes before i
                if self.twoVectors(self.transmitterPos, reflectionPoint, walls[index].A, walls[index].B) >= 0:
                    break
                if self.twoVectors(reflectionPoint, receiverPos, walls[index].A, walls[index].B) >= 0:
                    break
            else:
                H += self.calculateTransmitation(receiverPos, self.mirroredTransmittersPos[i], self.reflectionFactor)
        return H

    def calculateCrossPoint(self, A, B, C, D): #calc cross point single reflection
        if A[0] == B[0]:
            # First line is vertical
            x = A[0]
            a2 = (D[1] - C[1]) / (D[0] - C[0])
            b2 = C[1] - a2 * C[0]
            y = a2 * x + b2
            return [x, y]
        elif C[0] == D[0]:
            # Second line is vertical
            x = C[0]
            a1 = (B[1] - A[1]) / (B[0] - A[0])
            b1 = A[1] - a1 * A[0]
            y = a1 * x + b1
            return [x, y]
        a1 = (B[1] - A[1]) / (B[0] - A[0])
        b1 = A[1] - a1 * A[0]
        a2 = (D[1] - C[1]) / (D[0] - C[0])
        b2 = C[1] - a2 * C[0]
        x = (b2 - b1) / (a1 - a2)
        y = a1 * x + b1
        return [x, y]

    def createMirroredTransmitters(self, walls):
        mirroredTransmittersPos = np.zeros((len(walls), 2))
        for i in range(len(walls)):
            wall = walls[i]
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
        for wall in self.walls:
            x_coords = [wall.A[0], wall.B[0]]
            y_coords = [wall.A[1], wall.B[1]]
            plt.plot(x_coords, y_coords, color='black', linewidth=1)
        # Plot transmitter
        plt.scatter(self.transmitterPos[0], self.transmitterPos[1], color='red', label='Transmitter', zorder=5)
        # Plot mirrored transmitters pos
      
        plt.legend()
        # Display the plot
        plt.show()


start = time.time()
wall1 = Vector([0, 3], [3, 6])
wall2 = Vector([1, 3], [6, 3])
wall3 = Vector([6, 10], [12, 12])

wall5 = Vector([25, 10], [25, 19])
wall6 = Vector([5, 30], [10, 35])
wall7 = Vector([23, 36], [25, 37])

raytracing = Raytracing([40, 40], [20, 20], 5, 10, 0.7, [wall1,wall2,wall3,wall5,wall6,wall7])
print(raytracing.mirroredTransmittersPos)
raytracing.calculateRayTracing()
end = time.time() - start
print(f"Computation time: {end}")
raytracing.displayPowerMap()

        