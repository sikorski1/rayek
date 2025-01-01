export type MapTypes = {
    title: string;
    coordinates: number[][][];
    center: mapboxgl.LngLatLike;
    bounds: mapboxgl.LngLatBoundsLike;
};

export type postComputeTypes = {
	freq: string;
	stationH: string;
};