import mapboxgl from "mapbox-gl";

export type MapTypes = {
    title: string;
    coordinates: number[][][];
    center: mapboxgl.LngLatLike;
    bounds: mapboxgl.LngLatBoundsLike;
};