
export type MapTypes = {
    title: string;
    coordinates: number[][][];
    center: mapboxgl.LngLatLike;
    bounds: mapboxgl.LngLatBoundsLike;
    stationPos?: mapboxgl.LngLatLike;
};

export type postComputeTypes = {
	freq: string;
	stationH: string;
};