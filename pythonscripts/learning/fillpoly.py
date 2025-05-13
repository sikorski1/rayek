import numpy as np
import matplotlib.pyplot as plt
import cv2

# === 1. Load the matrix from the file ===
z_dim, y_dim, x_dim = 30, 250, 250
walls = np.fromfile("wallsMatrix3D_raw.bin", dtype=np.float64).reshape((z_dim, y_dim, x_dim))

# Create an empty array to store the processed data
processed_walls = np.zeros_like(walls)

# === 2. Iterate through all Z slices ===
for z in range(z_dim):
    print(f"Processing level Z = {z}")

    # --- Read the 2D slice ---
    map_slice = walls[z].copy()

    # --- Binary wall map (walls = 1, others = 0) ---
    walls_bin = (map_slice >= 1000).astype(np.uint8) * 255  # Important: for OpenCV, the image must be 0 or 255

    # --- Find external and internal contours ---
    contours, hierarchy = cv2.findContours(walls_bin, cv2.RETR_TREE, cv2.CHAIN_APPROX_SIMPLE)

    # --- Create a mask for filling ---
    fill_mask = np.zeros_like(walls_bin)

    # --- Iterate through contours ---
    for i, contour in enumerate(contours):
        # Check if the contour has a parent (internal contour) or is external (building contour)
        if hierarchy[0][i][1] == -1:  # If the contour has no parent, it is an external contour
            # Fill the external contour (building)
            cv2.drawContours(fill_mask, [contour], -1, 255, cv2.FILLED)
        else:
            # If the contour has a parent, it may be an internal contour (courtyard)
            # Check if it has a child, i.e., if it is not surrounded by other contours
            if hierarchy[0][i][2] == -1:  # Contour has no child, so it's a courtyard
                # Fill the internal contour (courtyard) with black color
                cv2.drawContours(fill_mask, [contour], -1, 255, cv2.FILLED)
            else:
                if hierarchy[0][i][3] == -1:
                    cv2.drawContours(fill_mask, [contour], -1, 0, cv2.FILLED)

    # --- Convert the mask 255/0 -> 1/0 ---
    logical_mask = (fill_mask == 255)

    # --- Fill buildings in empty areas (-150) with value 1000 ---
    filled_map = map_slice.copy()
    filled_map[(logical_mask) & (map_slice == -150)] = 1000

    # Save the processed slice into the processed_walls array
    processed_walls[z] = filled_map

    # === Visualization ===
plt.figure(figsize=(12, 6))

plt.subplot(1, 2, 1)
plt.imshow(walls[0], cmap='nipy_spectral', origin='lower')
plt.title(f"Original map (Z={0})")
plt.colorbar()
plt.grid(False)

plt.subplot(1, 2, 2)
plt.imshow(processed_walls[0], cmap='nipy_spectral', origin='lower')
plt.title(f"Filled map (Z={0})")
plt.colorbar()
plt.grid(False)

plt.tight_layout()
plt.show()

# === 3. Save the processed matrix to a binary file ===
processed_walls.tofile("wallsMatrix3D_processed.bin")
print("Processed matrix saved to 'wallsMatrix3D_processed.bin'")