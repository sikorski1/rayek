export function geoToMatrixIndex(
	lon: number,
	lat: number,
	lonMin: number,
	lonMax: number,
	latMin: number,
	latMax: number,
	size: number
): { i: number; j: number } {
	const y = ((latMax - lat) / (latMax - latMin)) * (size - 1);
	const x = ((lon - lonMin) / (lonMax - lonMin)) * (size - 1);

	const i = Math.round(x);
	const j = Math.round(y);

	return { i, j };
}
