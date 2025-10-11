export const loadWallMatrix = async (mapTitle: string): Promise<Int16Array> => {
	const response = await fetch(`/maps/${mapTitle}/wallsMatrix3D_processed.bin`);
	const buffer = await response.arrayBuffer();
	return new Int16Array(buffer);
};
