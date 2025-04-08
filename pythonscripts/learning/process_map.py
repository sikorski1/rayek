import numpy as np
import cv2
import argparse
import sys
import os

def process_z_slice(walls_3d_slice):
    """
    Przetwarza podany plaster 2D macierzy używając Flood Fill.
    Zwraca przetworzony plaster 2D.
    """
    # Pracujemy na kopii plastra, aby nie modyfikować oryginału nieoczekiwanie
    mapa = walls_3d_slice.copy()
    h, w = mapa.shape
    print(f"Python (process_z_slice): Wymiary plastra: {h}x{w}")

    # === Kod Flood Fill (Edge Barrier Method) ===
    sciany_bin = (mapa >= 1000).astype(np.uint8)
    flood_map = 1 - sciany_bin
    fill_value = 1
    padding = 1
    flood_map_padded = cv2.copyMakeBorder(flood_map, padding, padding, padding, padding,
                                          cv2.BORDER_CONSTANT, value=fill_value)
    h_pad, w_pad = flood_map_padded.shape
    flood_map_padded[0, :] = 0
    flood_map_padded[h_pad - 1, :] = 0
    flood_map_padded[:, 0] = 0
    flood_map_padded[:, w_pad - 1] = 0
    mask_ff = np.zeros((h_pad + 2, w_pad + 2), dtype=np.uint8)
    seed_point = (padding, padding)
    if flood_map_padded[seed_point] != fill_value:
        print(f"Python WARNING: Punkt startowy {seed_point} ma wartość {flood_map_padded[seed_point]}, oczekiwano {fill_value}.", file=sys.stderr)

    cv2.floodFill(flood_map_padded, mask_ff, seed_point, 0,
                  flags=4 | cv2.FLOODFILL_MASK_ONLY | (1 << 8))
    mask_cropped = mask_ff[1:-1, 1:-1]
    mask_inverted = 1 - mask_cropped
    mask_interiors_and_walls = mask_inverted[padding:-padding, padding:-padding]
    final_mask = (mask_interiors_and_walls == 1) & (sciany_bin == 0)
    final_mask = final_mask.astype(np.uint8)
    fill_count = np.sum(final_mask)
    print(f"Python (process_z_slice): Znaleziono {fill_count} pikseli wnętrz do wypełnienia.")
    mapa[final_mask == 1] = 1000.0 # Używamy float

    print(f"Python (process_z_slice): Zakończono przetwarzanie plastra.")
    # Zwracamy przetworzony plaster, upewniając się, że jest float64
    return mapa.astype(np.float64)

def main():
    # Zaktualizowany opis
    parser = argparse.ArgumentParser(description="Przetwarza plaster Z=0 macierzy 3D (raw binary) i zapisuje całą zmodyfikowaną macierz 3D.")
    parser.add_argument("--input", required=True, help="Ścieżka do wejściowego pliku raw binary 3D.")
    # Zaktualizowany opis wyjścia
    parser.add_argument("--output", required=True, help="Ścieżka do wyjściowego pliku raw binary 3D (z przetworzonym plastrem Z=0).")
    parser.add_argument("--dims", required=True, help="Wymiary macierzy 3D w formacie 'z,y,x' (np. '30,250,250').")

    args = parser.parse_args()

    try:
        dims_str = args.dims.split(',')
        if len(dims_str) != 3:
            raise ValueError("Wymiary muszą być w formacie 'z,y,x'")
        z_dim = int(dims_str[0])
        y_dim = int(dims_str[1])
        x_dim = int(dims_str[2])
        print(f"Python: Oczekiwane wymiary: z={z_dim}, y={y_dim}, x={x_dim}")
    except Exception as e:
        print(f"Python ERROR: Błąd parsowania wymiarów '{args.dims}': {e}", file=sys.stderr)
        sys.exit(1)

    if not os.path.exists(args.input):
        print(f"Python ERROR: Plik wejściowy nie istnieje: {args.input}", file=sys.stderr)
        sys.exit(1)

    try:
        print(f"Python: Wczytywanie pliku: {args.input}")
        walls_3d = np.fromfile(args.input, dtype=np.float64)
        expected_size = z_dim * y_dim * x_dim
        # Sprawdzanie rozmiaru pliku - tak jak poprzednio
        if walls_3d.size != expected_size:
             print(f"Python WARNING: Rozmiar wczytanych danych ({walls_3d.size}) różni się od oczekiwanego ({expected_size}).", file=sys.stderr)
             if walls_3d.size < expected_size:
                 print(f"Python ERROR: Niewystarczająca ilość danych w pliku wejściowym.", file=sys.stderr)
                 sys.exit(1)
             walls_3d = walls_3d[:expected_size]

        # Reshape do 3D
        walls_3d = walls_3d.reshape((z_dim, y_dim, x_dim))
        print(f"Python: Pomyślnie wczytano i zreshapowano dane do {walls_3d.shape}")

        # === KLUCZOWA ZMIANA: PRZETWARZANIE I PODMIANA PLASTRA ===
        # Wyodrębnij plaster Z=0
        slice_z0 = walls_3d[0]
        # Przetwórz ten plaster
        print(f"\nPython: Rozpoczynanie przetwarzania plastra Z=0...")
        processed_slice_z0 = process_z_slice(slice_z0)
        # Podmień oryginalny plaster Z=0 w macierzy 3D na przetworzony
        print(f"Python: Podmiana plastra Z=0 w macierzy 3D...")
        walls_3d[0] = processed_slice_z0
        print(f"Python: Plaster Z=0 podmieniony.")
        # ==========================================================

        # === ZMIANA: ZAPIS CAŁEJ MACIERZY 3D ===
        print(f"\nPython: Zapisywanie całej zmodyfikowanej macierzy 3D ({walls_3d.shape}) do: {args.output}")
        # Upewnij się, że typ jest poprawny przed zapisem
        walls_3d.astype(np.float64).tofile(args.output)
        print(f"Python: Zapis zakończony pomyślnie.")
        # ========================================

    except Exception as e:
        print(f"Python ERROR: Wystąpił błąd podczas przetwarzania: {e}", file=sys.stderr)
        sys.exit(1)

    sys.exit(0)

if __name__ == "__main__":
    main()