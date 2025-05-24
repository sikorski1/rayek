export const loadWallMatrix = async (mapTitle: string): Promise<Float64Array> => {
	const response = await fetch(`/maps/${mapTitle}/wallsMatrix3D_processed.bin`);
	const buffer = await response.arrayBuffer();
	return new Float64Array(buffer);
};