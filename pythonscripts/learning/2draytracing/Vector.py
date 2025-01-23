from math import sqrt
class Vector:
    def __init__(self, firstPoint, secondPoint):
        self.A = firstPoint
        self.B = secondPoint
        self.length = self.calculateLength(self.A, self.B)
        pass
    def calculateLength(self):
        length = sqrt((self.A.x - self.B.x) ** 2 + (self.A.y - self.B.y) ** 2)
        return length