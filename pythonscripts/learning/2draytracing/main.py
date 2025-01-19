import numpy as np
class Raytracing:
    def __init__(self, matrixDimensions, tPos, tPower, tFreq, rFactor, oPos):
        self.matrix = self.createMatrix(matrixDimensions)
        self.transmitterPos = tPos
        self.transmitterPower = tPower #mW
        self.transmitterFreq = tFreq # GHz
        self.reflectionFactor = rFactor
        self.obstaclesPos = oPos 

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
        pass
        
    
        
raytracing = Raytracing([16, 28], [12.05, 7.05], 5, 3.6, 0.7, [[[0, 20.05],[10, 20.05]], [[13, 20.05],[16, 20.05]]])
raytracing.createMatrix()
    