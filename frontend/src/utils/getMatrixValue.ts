export const getMatrixValue = (
	matrix: Float64Array,
	x: number,
	y: number,
	z: number,
	sizeX = 400,
	sizeY = 400,
	sizeZ = 30
): number | null => {
	if (x < 0 || x >= sizeX || y < 0 || y >= sizeY || z < 0 || z >= sizeZ) {
		return null;
	}
	const index = z * sizeY * sizeX + y * sizeX + x;
	return matrix[index];
};

