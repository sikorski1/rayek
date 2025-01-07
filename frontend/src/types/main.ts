
export type MapTypes = {
    title: string;
    coordinates: number[][][];
    center: mapboxgl.LngLatLike;
    bounds: mapboxgl.LngLatBoundsLike;
    stationPos?: mapboxgl.LngLatLike;
    setStationPos?: (value: mapboxgl.LngLatLike | null | ((prevPos: mapboxgl.LngLatLike | null) => mapboxgl.LngLatLike | null)) => void;
};

export type postComputeTypes = {
	freq: string;
	stationH: string;
};