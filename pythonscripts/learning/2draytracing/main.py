import numpy as np
from math import sqrt, log10,pi,e
import matplotlib.pyplot as plt
from Vector import Vector
class Raytracing:
    def __init__(self, matrixDimensions, tPos, tPower, tFreq, rFactor, oPos):
        self.step = 0.05
        self.transmitterPos = tPos
        self.transmitterPower = tPower #mW
        self.transmitterFreq = tFreq # GHz
        self.waveLength = 299792458 / tFreq / 10 ** 9;
        self.reflectionFactor = rFactor
        self.obstaclesPos = oPos 
        self.powerMap = np.zeros((matrixDimensions[1]*int((1/self.step))+1, matrixDimensions[0]*int((1/self.step))+1))
        self.matrix = self.createMatrix(matrixDimensions)
        self.mirroredTransmittersPos = self.createMirroredTransmitters()
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

                checkWall1Reflection = False
                checkWall2Reflection = False
                if self.twoVectors(receiverPos, self.mirroredTransmittersPos[0], self.obstaclesPos[0].A, self.obstaclesPos[0].B) == 1 and self.twoVectors(receiverPos, self.mirroredTransmittersPos[0], self.obstaclesPos[1].A, self.obstaclesPos[1].B) == -1:
                        distance = abs(self.obstaclesPos[1].B[1] - self.obstaclesPos[0].A[1])
                        bottomPoint = [self.obstaclesPos[1].A[0], self.obstaclesPos[0].A[1] + distance]
                        topPoint = [self.obstaclesPos[1].A[0], self.obstaclesPos[0].A[1] + distance + self.obstaclesPos[1].length]
                        mirroredWall = [bottomPoint, topPoint]
                        if self.twoVectors(receiverPos, self.mirroredTransmittersPos[0], mirroredWall[0], mirroredWall[1]) == -1:
                            checkWall1Reflection = True 
                if self.twoVectors(receiverPos, self.mirroredTransmittersPos[1], self.obstaclesPos[1].A, self.obstaclesPos[1].B) == 1 and self.twoVectors(receiverPos, self.mirroredTransmittersPos[1], self.obstaclesPos[0].A, self.obstaclesPos[0].B) == -1:
                        distance = abs(self.obstaclesPos[0].B[0] - self.obstaclesPos[1].B[0])
                        leftPoint = [self.obstaclesPos[1].A[0]+distance,self.obstaclesPos[0].B[1]]
                        rightPoint = [self.obstaclesPos[1].A[0]+distance+self.obstaclesPos[0].length,self.obstaclesPos[0].B[1]]
                        mirroredWall = [leftPoint, rightPoint]
                        if self.twoVectors(receiverPos, self.mirroredTransmittersPos[1], mirroredWall[0], mirroredWall[1]) == -1:
                            checkWall2Reflection = True

                if checkWall1Reflection: # add transmitation from 1 wall reflection
                    H += self.calculateTransmitation(receiverPos, self.mirroredTransmittersPos[0])
                if checkWall2Reflection: # add transmitation from 2 wall reflection
                    H += self.calculateTransmitation(receiverPos, self.mirroredTransmittersPos[1])
                    
                if H == 0:
                    self.powerMap[i][j] = -150
                else:
                    self.powerMap[i][j] = 10*log10(self.transmitterPower) + 20*log10(abs(H))
    
    def checkLineOfSight(self, receiverPos, walls):
        for wall in walls:
            if self.twoVectors(receiverPos, self.transmitterPos, wall.A, wall.B) >= 0:  #checking collision with wall line of sight
                return False
        return True
    def checkSingleWallReflection(self, receiverPos, walls):
        checkTable = [False for i in range(len(walls))]
        for i, wall in enumerate(walls):
            if self.twoVectors(receiverPos, self.mirroredTransmittersPos[i], wall.A, wall.B) < 0:
                continue
            for j in range(len(walls) - 1):
                index = (i + j + 1) % len(walls)
                mirroredWall = self.calculateMirroredWall(wall, walls[index])
                if self.twoVectors(receiverPos, self.mirroredTransmittersPos[i], walls[index].A, walls[index].B) > 0 or \
                self.twoVectors(receiverPos, self.mirroredTransmittersPos[i], mirroredWall.A, mirroredWall.B) > 0:
                    break
                
            else:
                checkTable[i] = True
        print(checkTable)
        return checkTable

    def calculateMirroredWall(self, mainWall, wallThatIsMirrored):
        if mainWall.B[0] == mainWall.A[0]:
            aPrimX = 2 * mainWall.A[0] - wallThatIsMirrored.A[0]
            bPrimX = 2 * mainWall.A[0] - wallThatIsMirrored.B[0]
            mirroredWall = Vector([aPrimX, wallThatIsMirrored.A[1]], [bPrimX, wallThatIsMirrored.A[1]])
            print(mirroredWall.A, mirroredWall.B)
            return mirroredWall
        if mainWall.B[1] == mainWall.A[1]:  
            y_prim = mainWall.A[1] 
            aPrimY = 2 * y_prim - wallThatIsMirrored.A[1]
            bPrimY = 2 * y_prim - wallThatIsMirrored.B[1]
            mirroredWall = Vector([wallThatIsMirrored.A[0], aPrimY], [wallThatIsMirrored.B[0], bPrimY])
            print(mirroredWall.A)
            print(mirroredWall.B)
            return mirroredWall
        m = (mainWall.B[1] - mainWall.A[1])/(mainWall.B[0] - mainWall.A[0])
        b =  mainWall.A[1] - m * mainWall.A[0]
        aPrimX = (wallThatIsMirrored.A[0] + m * (wallThatIsMirrored.A[1] - b)) / (1 + m**2)
        aPrimY = (m * wallThatIsMirrored.A[0] + m**2 * wallThatIsMirrored.A[1] - m * b) / (1 + m**2)
        bPrimX = (wallThatIsMirrored.B[0] + m * (wallThatIsMirrored.B[1] - b)) / (1 + m**2)
        bPrimY = (m * wallThatIsMirrored.B[0] + m**2 * wallThatIsMirrored.B[1] - m * b) / (1 + m**2)
        mirroredWall = Vector([aPrimX, aPrimY], [bPrimX, bPrimY])
        print(mirroredWall.A)
        print(mirroredWall.B)
        return mirroredWall

    def createMirroredTransmitters(self):
        mirroredTransmittersPos = np.zeros((len(self.obstaclesPos), 2))
        for i in range(len(self.obstaclesPos)):
            if self.obstaclesPos[i].A[0] == self.obstaclesPos[i].B[0]:
                mirroredTransmittersPos[i][1] = self.transmitterPos[1]
                distance = abs(self.obstaclesPos[i].A[0] - self.transmitterPos[0])
                if self.transmitterPos[0] < self.obstaclesPos[i].A[0]:
                    mirroredTransmittersPos[i][0] = self.obstaclesPos[i].A[0] + distance
                else:
                    mirroredTransmittersPos[i][0] = self.obstaclesPos[i].A[0] - distance
            elif self.obstaclesPos[i].A[1] == self.obstaclesPos[i].B[1]:
                mirroredTransmittersPos[i][0] = self.transmitterPos[0]
                distance = abs(self.obstaclesPos[i].A[1] - self.transmitterPos[1])
                if self.transmitterPos[1] < self.obstaclesPos[i].A[1]:
                    mirroredTransmittersPos[i][1] = self.obstaclesPos[i].A[1] + distance
                else:
                    mirroredTransmittersPos[i][1] = self.obstaclesPos[i].A[1] - distance
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

        for wall in self.obstaclesPos:
            x_coords = [wall.A, wall.A]
            y_coords = [wall.B, wall.B]
            plt.plot(x_coords, y_coords, color='black', linewidth=2, label='Wall')

        plt.legend()
        plt.show()

wall1 = Vector([0, 20],[10, 20])
wall2 = Vector([12, 6],[12, 12])
        
raytracing = Raytracing([16, 28], [13, 7.85], 5, 3.6, 0.7, [wall1, wall2])
raytracing.calculateMirroredWall(wall1, wall2)
raytracing.checkSingleWallReflection([5,7.5], [wall1, wall2])
raytracing.calculateRayTracing()
raytracing.displayPowerMap()