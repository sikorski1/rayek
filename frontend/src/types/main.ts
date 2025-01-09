
export type MapTypes = {
    title: string;
    coordinates: number[][][];
    center: mapboxgl.LngLatLike;
    bounds: mapboxgl.LngLatBoundsLike;
};

export type MapTypesExtended = MapTypes & {
    stationPos: mapboxgl.LngLatLike;
    setStationPos: (value: mapboxgl.LngLatLike | null | ((prevPos: mapboxgl.LngLatLike | null) => mapboxgl.LngLatLike | null)) => void;
    buildingsData: GeoJSON.FeatureCollection | null
}

export type postComputeTypes = {
	freq: string;
	stationH: string;
};