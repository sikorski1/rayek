import numpy as np
import matplotlib.pyplot as plt
import cv2

# === 1. Wczytaj macierz z pliku ===
z_dim, y_dim, x_dim = 30, 250, 250
walls = np.fromfile("wallsMatrix3D_raw.bin", dtype=np.float64).reshape((z_dim, y_dim, x_dim))

# === 2. Wyciągamy przekrój Z = 0 ===
mapa = walls[0].copy()  # tylko 2D

# === 3. Szukamy ścian (wartości >=1000) ===
sciany_bin = (mapa >= 1000).astype(np.uint8)

# === 4. Znajdź zewnętrzne kontury ===
kontury_zewnetrzne, _ = cv2.findContours(sciany_bin, cv2.RETR_EXTERNAL, cv2.CHAIN_APPROX_SIMPLE)

# === 5. Utworzenie maski do wypełnienia ===
maska_wypelnienia = np.zeros_like(sciany_bin)

# === 6. Wypełnij obszar zewnętrznymi konturami na biało (1) ===
cv2.drawContours(maska_wypelnienia, kontury_zewnetrzne, -1, 1, thickness=cv2.FILLED)

# === 7. Znajdź wszystkie kontury z hierarchią ===
kontury_wszystkie, hierarchia = cv2.findContours(sciany_bin, cv2.RETR_TREE, cv2.CHAIN_APPROX_SIMPLE)

# === 8. Iteruj po konturach i ich hierarchii ===
for i, kontur in enumerate(kontury_wszystkie):
    ma_dziadka = False
    ma_rodzica = hierarchia[0][i][3] >= 0
    
    if ma_rodzica:
        # Sprawdź czy kontur ma dziadka
        indeks_rodzica = hierarchia[0][i][3]
        ma_dziadka = hierarchia[0][indeks_rodzica][3] >= 0
    
    # Sprawdź czy kontur nie ma dziecka
    if hierarchia[0][i][2] < 0 and ma_dziadka:
        # Jeśli kontur nie ma dziecka i ma dziadka, wypełnij kontur czarnym kolorem (0)
        cv2.drawContours(maska_wypelnienia, [kontur], -1, 0, thickness=cv2.FILLED)

# === 9. Wypełnij wnętrze budynków wartością 1000 ===
mapa_wypelniona = mapa.copy()
mapa_wypelniona[(maska_wypelnienia == 1) & (mapa == -150)] = 1000

# === 10. Wizualizacja ===
plt.figure(figsize=(12, 6))

# Oryginalna mapa
plt.subplot(1, 2, 1)
plt.imshow(mapa, cmap='nipy_spectral', origin='lower')
plt.title("Oryginalna mapa (Z=0)")
plt.colorbar()
plt.grid(False)

# Wypełniona mapa
plt.subplot(1, 2, 2)
plt.imshow(mapa_wypelniona, cmap='nipy_spectral', origin='lower')
plt.title("Wypełnione wnętrza budynków (Z=0)")
plt.colorbar()
plt.grid(False)

plt.tight_layout()
plt.show()