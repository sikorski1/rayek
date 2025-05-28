export function matrixIndexToGeo(
	i: number,
	j: number,
	lonMin: number,
	lonMax: number,
	latMin: number,
	latMax: number,
	size: number
): { lon: number; lat: number } {
	const lon = lonMin + (i / (size - 1)) * (lonMax - lonMin);
	const lat = latMax - (j / (size - 1)) * (latMax - latMin);

	return { lon, lat };
}
