export function geoToMatrixIndex(
	lat: number,
	lon: number,
	latMin: number,
	latMax: number,
	lonMin: number,
	lonMax: number,
	size: number
): { i: number; j: number } {
	const y = ((lat - latMin) / (latMax - latMin)) * (size - 1);
	const x = ((lon - lonMin) / (lonMax - lonMin)) * (size - 1);

	const i = Math.round(x);
	const j = Math.round(y);

	return { i, j };
}
