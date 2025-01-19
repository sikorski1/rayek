import numpy as np
from math import sqrt, log10,pi
import matplotlib.pyplot as plt
class Raytracing:
    def __init__(self, matrixDimensions, tPos, tPower, tFreq, rFactor, oPos):
        self.matrix = self.createMatrix(matrixDimensions)
        self.transmitterPos = tPos
        self.transmitterPower = tPower #mW
        self.transmitterFreq = tFreq # GHz
        self.waveLength = 299792458 / tFreq / 10 ** 9;
        self.reflectionFactor = rFactor
        self.obstaclesPos = oPos 
        self.powerMap = np.zeros((matrixDimensions[1]*10+1, matrixDimensions[0]*10+1))
    def createMatrix(self, matrixDimensions):
        nx, ny = (matrixDimensions[0], matrixDimensions[1])
        step = 0.1
        intStepX = int(nx / step) + 1
        intStepY = int(ny / step) + 1
        x = np.linspace(0, nx, intStepX)
        y = np.linspace(0, ny, intStepY)

        xv, yv = np.meshgrid(x, y)
        matrix = np.stack((xv, yv), axis=-1)
        return matrix

    def twoVectors(self, A_x, A_y, B_x, B_y, C_x, C_y, D_x, D_y):
        result = (((C_x - A_x)*(B_y - A_y) - (B_x - A_x)*(C_y - A_y)) * ((D_x - A_x) * (B_y - A_y) - (B_x - A_x) * (D_y - A_y)))  #[C_x-A-x, B_x-A_x, vector CA, BA
                                                                                                                                   #C_y-A-y, B_y-A_y]
        if result > 0:
            return -1
        else:
            result2 = (((A_x - C_x)*(D_y - C_y) - (D_x - C_x)*(A_y - C_y)) * ((B_x - C_x) * (D_y - C_y) - (D_x - C_x) * (B_y - C_y)))
            if result2 > 0:
                return -1
            elif result < 0 and result2 < 0:
                return 1
            elif result == 0 and result2 < 0:
                return 0
            elif result < 0 and result2 == 0:
                return 0
            elif A_x < C_x and A_x < D_x and B_x < C_x and B_x < D_x:
                return -1
            elif A_y < C_y and A_y < D_y and B_y < C_y and B_y < D_y:
                return -1
            elif A_x > C_x and A_x > D_x and B_x > C_x and B_x > D_x:
                return -1
            elif A_y > C_y and A_y > D_y and B_y > C_y and B_y > D_y:
                return -1
            else:
                return 0
    def calculateRayTracing(self):
        for i in range(len(self.matrix)):
            for j in range(len(self.matrix[0])):
                receiverPos = self.matrix[i][j]
                if self.twoVectors(receiverPos[0], receiverPos[1], self.transmitterPos[0], self.transmitterPos[1], self.obstaclesPos[0][0][0], self.obstaclesPos[0][0][1], self.obstaclesPos[0][1][0], self.obstaclesPos[0][1][1]) == 1: #checking collision with wall
                    self.powerMap[i][j] = -150
                elif self.twoVectors(receiverPos[0], receiverPos[1], self.transmitterPos[0], self.transmitterPos[1], self.obstaclesPos[1][0][0], self.obstaclesPos[1][0][1], self.obstaclesPos[1][1][0], self.obstaclesPos[1][1][1]) == 1:
                    self.powerMap[i][j] = -150
                else:
                    self.powerMap[i][j] = self.calculatePower(receiverPos, self.transmitterPos)
    def calculatePower(self, p1, p2):
        power = 10*log10(self.transmitterPower*(self.waveLength/(4*pi*self.calculateDist(p1, p2)))**2) #10log(Pt*(lambda/(4*pi*r))^2)
        return power
    def calculateDist(self, p1, p2):
        dist = sqrt((p1[0] - p2[0])**2 + (p1[1] - p2[1])**2)
        return dist    
    def displayPowerMap(self):
        plt.figure(figsize=(10, 8))
        plt.imshow(self.powerMap, origin='lower', cmap='jet', extent=[0, self.matrix.shape[1]*0.1, 0, self.matrix.shape[0]*0.1])
        plt.colorbar(label='Power (dBm)')
        plt.title('Power Map')
        plt.xlabel('X Coordinate (m)')
        plt.ylabel('Y Coordinate (m)')
        plt.show()
        
raytracing = Raytracing([16, 28], [1.05, 7.05], 5, 3.6, 0.7, [[[0, 20.05],[10, 20.05]], [[11, 20.05],[16, 20.05]]])
raytracing.calculateRayTracing()
raytracing.displayPowerMap()
    