import numpy as np
from math import sqrt, log10, pi, e
import matplotlib.pyplot as plt
import time

# --- Definicja klasy Vector ---
class Vector:
    def __init__(self, A, B):
        self.A = A
        self.B = B

class Raytracing:
    def __init__(self, matrixDimensions, tPos, tPower, tFreq, rFactor, oPos):
        self.step = 0.1 # Krok siatki (im mniejszy, tym dokładniejsza mapa, ale wolniejsze obliczenia)
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
        result = (((C[0] - A[0])*(B[1] - A[1]) - (B[0] - A[0])*(C[1] - A[1])) * ((D[0] - A[0]) * (B[1] - A[1]) - (B[0] - A[0]) * (D[1] - A[1])))
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
        print(f"Rozpoczynam obliczenia dla {len(self.walls)} ścian...")
        for i in range(len(self.matrix)):
            for j in range(len(self.matrix[0])):
                H = 0
                receiverPos = self.matrix[i][j]
                if self.checkLineOfSight(receiverPos, self.walls):
                    H += self.calculateTransmitation(receiverPos, self.transmitterPos) 
                    
                H += self.calculateSingleWallReflection(receiverPos, self.walls)
                if H == 0:
                    self.powerMap[i][j] = -150
                else:
                    self.powerMap[i][j] = 10*log10(self.transmitterPower) + 20*log10(abs(H))
    
    def checkLineOfSight(self, receiverPos, walls):
        for wall in walls:
            if self.twoVectors(receiverPos, self.transmitterPos, wall.A, wall.B) >= 0:
                return False
        return True
    
    def calculateSingleWallReflection(self, receiverPos, walls):
        H = 0
        for i, wall in enumerate(walls):
            if self.twoVectors(receiverPos, self.mirroredTransmittersPos[i], wall.A, wall.B) <= 0:
                continue
            else:
                reflectionPoint = self.calculateCrossPoint(receiverPos, self.mirroredTransmittersPos[i], wall.A, wall.B)
            for j in range(len(walls) - 1):
                index = (i + j + 1) % len(walls)
                if self.twoVectors(self.transmitterPos, reflectionPoint, walls[index].A, walls[index].B) >= 0:
                    break
                if self.twoVectors(reflectionPoint, receiverPos, walls[index].A, walls[index].B) >= 0:
                    break
            else:
                H += self.calculateTransmitation(receiverPos, self.mirroredTransmittersPos[i], self.reflectionFactor)
        return H

    def calculateCrossPoint(self, A, B, C, D):
        if A[0] == B[0]:
            x = A[0]
            a2 = (D[1] - C[1]) / (D[0] - C[0])
            b2 = C[1] - a2 * C[0]
            y = a2 * x + b2
            return [x, y]
        elif C[0] == D[0]:
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
        
        # Rysowanie ścian
        for i, wall in enumerate(self.walls):
            x_coords = [wall.A[0], wall.B[0]]
            y_coords = [wall.A[1], wall.B[1]]
            # Dodajemy etykietę tylko raz do legendy
            label_text = "Wall" if i == 0 else None
            plt.plot(x_coords, y_coords, color='red', linewidth=2, label=label_text)

        # Plot transmitter
        plt.scatter(self.transmitterPos[0], self.transmitterPos[1], color='white', edgecolor='black', s=100, label='Transmitter', zorder=5)
      
        plt.legend(loc='upper right')
        plt.show()

# ---------------------------------------------------------
# --- KONFIGURACJA SCENARIUSZA I URUCHOMIENIE ---
# ---------------------------------------------------------

# WYBIERZ SCENARIUSZ: 1 (15 ścian) lub 2 (30 ścian)
SCENARIO_ID = 2 

walls_list = []

if SCENARIO_ID == 1:
    # OPCJA 1: 15 Ścian (Mieszanka długich i ukośnych)
    walls_list = [
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
    print("Wybrano scenariusz 1: 15 ścian.")

elif SCENARIO_ID == 2:
    # OPCJA 2: 30 Ścian (Bardziej skomplikowany labirynt/miejski kanion)
    walls_list = [
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
    print("Wybrano scenariusz 2: 30 ścian.")

   
walls_scenario_3 = [
    # --- 1. ZEWNĘTRZNY PIERŚCIEŃ (GRANICE) ---
    # Lewa krawędź
    Vector([2, 2], [2, 10]),
    Vector([2, 12], [2, 20]),
    Vector([2, 22], [2, 30]),  # Zagęszczone
    Vector([2, 32], [2, 38]),

    # Prawa krawędź
    Vector([38, 2], [38, 12]),
    Vector([38, 15], [38, 25]),
    Vector([38, 28], [38, 38]),

    # Dolna krawędź
    Vector([5, 38], [12, 38]),
    Vector([15, 38], [25, 38]),
    Vector([28, 38], [35, 38]),

    # Górna krawędź
    Vector([5, 2], [15, 2]),
    Vector([18, 2], [28, 2]),
    Vector([32, 2], [36, 2]),

    # --- 2. WEWNĘTRZNE STRUKTURY (ZYGZAKI I KORYTARZE) ---
    # Lewy górny róg - labirynt
    Vector([8, 5], [8, 10]),
    Vector([8, 10], [12, 12]),
    Vector([12, 12], [15, 8]),
    Vector([5, 8], [8, 12]),  # Dodatkowa ukośna

    # Lewy dolny róg - "pokoje"
    Vector([8, 32], [12, 28]),
    Vector([12, 28], [16, 32]),
    Vector([5, 25], [10, 25]),
    Vector([10, 25], [10, 30]),

    # --- 3. CENTRUM (WOKÓŁ NADAJNIKA 20,20) ---
    # Skomplikowane otoczenie środka
    Vector([18, 18], [22, 18]),  # Daszek nad środkiem
    Vector([18, 22], [18, 25]),  # Lewa ścianka dolna
    Vector([22, 22], [25, 25]),  # Prawa ukośna
    Vector([15, 20], [17, 18]),  # Osłona lewa
    Vector([23, 18], [25, 20]),  # Osłona prawa

    # --- 4. PRAWA STRONA - REFLEKTORY I ROZPRASZACZE ---
    # Górne "X"
    Vector([30, 10], [35, 15]),
    Vector([30, 15], [35, 10]),

    # Dolne "X"
    Vector([30, 25], [35, 30]),
    Vector([30, 30], [35, 25]),

    # Pionowe przegrody "grzebień"
    Vector([28, 5], [28, 10]),
    Vector([28, 20], [28, 15]),  # Środkowa przegroda
    Vector([28, 30], [28, 35]),

    # --- 5. DROBNE PRZESZKODY I FILARY (HARD MODE DLA RAYTRACINGU) ---
    # Rozrzucone małe elementy
    Vector([15, 5], [15, 7]),    # Filar góra
    Vector([25, 5], [25, 7]),    # Filar góra 2
    Vector([32, 20], [35, 20]),  # Pozioma belka
    Vector([12, 35], [14, 32]),  # Mały ukos dół
    Vector([20, 35], [20, 32]),  # Pionowy słupek dół
    Vector([22, 5], [24, 3]),    # Zamykający ukos góra
    Vector([3, 18], [6, 18]),    # Pozioma belka lewa
    Vector([16, 15], [17, 14]),  # Mikro przeszkoda 1
    Vector([24, 15], [23, 14]),  # Mikro przeszkoda 2
    Vector([20, 30], [22, 28]),  # "V" shape part 1
    Vector([22, 28], [24, 30]),  # "V" shape part 2
    Vector([32, 35], [35, 36]),  # Płaski ukos róg
    Vector([5, 35], [7, 33]),    # Narożnik lewy dół
]
print("Wybrano scenariusz 3: 50+ ścian.")

# Uruchomienie symulacji
start = time.time()
raytracing = Raytracing([40, 40], [20, 20], 5, 10, 0.7, walls_list)
raytracing.calculateRayTracing()
end = time.time() - start
print(f"Computation time: {end:.2f}s")
raytracing.displayPowerMap()