import numpy as np
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

    _, height_str = match.groups()
    height = int(height_str)
    if height < 0 or height >= z_dim:
        continue

    png_path = os.path.join(png_folder, file_name)
    img = cv2.imread(png_path)
    if img is None:
        print(f"Could not read image: {png_path}")
        continue

    img_rgb = cv2.cvtColor(img, cv2.COLOR_BGR2RGB)

    gray_mask = np.all(img_rgb == [192, 192, 192], axis=-1)
    magenta_mask = np.all(img_rgb == [255, 0, 255], axis=-1)
    magenta_dark_mask = np.all(img_rgb == [128, 0, 128], axis=-1)
    yellow_mask = np.all(img_rgb == [255, 255, 0], axis=-1)

    gray_mask_flipped = np.flipud(gray_mask)
    magenta_mask_flipped = np.flipud(magenta_mask)
    magenta_dark_mask_flipped = np.flipud(magenta_dark_mask)
    yellow_mask_flipped = np.flipud(yellow_mask)

    layer = walls[height].copy()

    condition_mask = np.logical_or(layer == -150, layer >= 1000)
    layer[np.logical_and(layer == -150, gray_mask_flipped)] = 5000
    layer[np.logical_and(condition_mask, magenta_mask_flipped)] = 10000
    layer[np.logical_and(layer == -150, yellow_mask_flipped)] = 20000
    layer[np.logical_and(condition_mask, magenta_dark_mask_flipped)] = 10001

    walls[height] = np.flipud(layer)

    print(f"Processed height {height}")

walls.tofile("wallsMatrix3D_processed.bin")
print("Saved wallsMatrix3D_processed.bin")