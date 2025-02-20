import numpy as np
import matplotlib.pyplot as plt
import pandas as pd
import json
import os
def display_power_map():
    # Wczytaj mapę mocy z pliku CSV
    go_dir = "../backendGo/utils/2draylaunching/"
    
    # Wczytaj mapę mocy z pliku CSV
    power_map = pd.read_csv(go_dir+'ray_map.csv', header=None).values
    
    # Wczytaj konfigurację
    with open(go_dir+'ray_config.json', 'r') as f:
        config = json.load(f)

    transmitter_pos = (config['tx'], config['ty'])
    walls = config['walls']
    
    # Określ rozmiar kroku (z Twojego oryginalnego kodu Go)
    step = 0.1
    
    # Utwórz wykres
    plt.figure(figsize=(10, 8))
    
    # Wyświetl mapę mocy
    extent = [0, power_map.shape[1]*step, 0, power_map.shape[0]*step]
    im = plt.imshow(power_map, origin='lower', cmap='jet', extent=extent, vmin=-150, vmax=0)
    
    # Dodaj pasek kolorów
    plt.colorbar(im, label='Power (dBm)')
    
    # Dodaj tytuł i etykiety osi
    plt.title('Power Map')
    plt.xlabel('X Coordinate (m)')
    plt.ylabel('Y Coordinate (m)')
    
    # Narysuj ściany
    for wall in walls:
        x_coords = [wall[0], wall[2]]
        y_coords = [wall[1], wall[3]]
        plt.plot(x_coords, y_coords, color='black', linewidth=1)
    
    # Narysuj pozycję nadajnika
    plt.scatter(transmitter_pos[0], transmitter_pos[1], color='red', label='Transmitter', zorder=5)
    
    # Dodaj legendę
    plt.legend()
    
    # Wyświetl wykres
    plt.savefig("power_map_python.png", dpi=300)
    plt.show()

if __name__ == "__main__":
    print(os.getcwd())
    display_power_map()
