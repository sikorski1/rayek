from math import sqrt
class Vector:
    def __init__(self, firstPoint, secondPoint):
        self.A = firstPoint
        self.B = secondPoint
        self.AB = [secondPoint[0]-firstPoint[0], secondPoint[1]-firstPoint[1]]
        self.length = self.calculateLength()
        pass
    def calculateLength(self):
        length = sqrt(self.AB[0] ** 2 + self.AB[1] ** 2)
        return length