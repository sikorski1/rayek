import numpy as np
import matplotlib.pyplot as plt
import cv2

z_dim, y_dim, x_dim = 30, 250, 250
walls = np.fromfile("wallsMatrix3D_raw.bin", dtype=np.float64).reshape((z_dim, y_dim, x_dim))
mapa = walls[0].copy()  
obraz_png = cv2.imread("aghlibrary_floor.png")
obraz_szary = cv2.cvtColor(obraz_png, cv2.COLOR_BGR2GRAY)
obraz_szary_odbity = cv2.flip(obraz_szary, 0)  
mapa_wypelniona = mapa.copy()

for y in range(y_dim):
    for x in range(x_dim):
        if mapa[y, x] == -150 and obraz_szary_odbity[y, x] == 0:
            mapa_wypelniona[y, x] = 1000
plt.figure(figsize=(15, 5))

plt.subplot(1, 3, 1)
plt.imshow(mapa, cmap='nipy_spectral', origin='lower')
plt.title("Oryginalna mapa (Z=0)")
plt.colorbar()
plt.grid(False)

plt.subplot(1, 3, 2)
plt.imshow(cv2.cvtColor(cv2.flip(obraz_png, 1), cv2.COLOR_BGR2RGB), origin='lower')
plt.title("Odbity lustrzanie obraz PNG")
plt.grid(False)

plt.subplot(1, 3, 3)
plt.imshow(mapa_wypelniona, cmap='nipy_spectral', origin='lower')
plt.title("Mapa wypełniona na podstawie odbitego PNG")
plt.colorbar()
plt.grid(False)

plt.tight_layout()
plt.show()

walls_wypelnione = walls.copy()
walls_wypelnione[0] = mapa_wypelniona

# Zapisz wypełnioną macierz 3D do pliku
walls_wypelnione.tofile("wallsMatrix3D_processed.bin")