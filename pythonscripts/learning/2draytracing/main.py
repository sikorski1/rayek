import numpy as np
class Raytracing:
    def __init__(self, matrixDimensions):
        self.matrixDimensions = matrixDimensions

    def createMatrix(self):
        nx, ny = (self.matrixDimensions[0], self.matrixDimensions[1])
        step = 0.1
        intStepX = (int(nx/step)+1)
        intStepY = (int(ny/step)+1)
        x = np.arange(0, nx+step, step)
        y = np.arange(0, ny+step, step)
        matrix = np.zeros((intStepY, intStepX, 2)) # create ny/step x nx/step matrix 
        for i in range(intStepX):
            for j in range(intStepY):
                matrix[j][i] = [x[i], y[j]] #fill matrix with [x, y] positions
        print(matrix[0])

        
    
        
raytracing = Raytracing([16,28])
raytracing.createMatrix()
    