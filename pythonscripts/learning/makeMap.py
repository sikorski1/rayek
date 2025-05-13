import numpy as np
import matplotlib.pyplot as plt
import cv2
import os
import re

z_dim, y_dim, x_dim = 30, 250, 250
walls = np.fromfile("wallsMatrix3D_raw.bin", dtype=np.float64).reshape((z_dim, y_dim, x_dim))

png_folder = "final"

file_pattern = re.compile(r"^(.*)-(\d+)m\.png$")

for file_name in os.listdir(png_folder):
    match = file_pattern.match(file_name)
    if not match:
        continue

    base_name, height_str = match.groups()
    height = int(height_str)
    if height < 0 or height >= z_dim:
        continue  

    png_path = os.path.join(png_folder, file_name)
    img = cv2.imread(png_path)

    if img is None:
        print(f"Could not read image: {png_path}")
        continue

    img_rgb = cv2.cvtColor(img, cv2.COLOR_BGR2RGB)
    img_flipped = cv2.flip(img_rgb, 0)

    gray_mask = np.all(img_flipped == [192, 192, 192], axis=-1)

    layer = walls[height].copy()
    layer[np.logical_and(layer == -150, gray_mask)] = 5000
    walls[height] = layer
    print(f"Processed height {height} from file: {file_name}")

walls.tofile("wallsMatrix3D_processed.bin")
print("3D matrix saved to wallsMatrix3D_processed.bin")