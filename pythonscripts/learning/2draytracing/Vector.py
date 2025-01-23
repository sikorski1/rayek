from math import sqrt
class Vector:
    def __init__(self, firstPoint, secondPoint):
        self.A = firstPoint
        self.B = secondPoint
        self.length = self.calculateLength(self.A, self.B)
        pass
    def calculateLength(self):
        length = sqrt((self.A[0] - self.B[0]) ** 2 + (self.A[1] - self.B[1]) ** 2)
        return length