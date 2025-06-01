import { MercatorCoordinate } from "mapbox-gl";
export type Maps = {
	id: string;
	name: string;
	description: string;
	img: string;
};

export type SingleRay = {
	azimuth: number;
	elevation: number;
};
export type SettingsDataTypes = {
	settingsType: string;
	isOpen: boolean;
	frequency: number;
	stationHeight: number;
	numberOfRaysAzimuth: number;
	numberOfRaysElevation: number;
	numberOfInteractions: number;
	reflectionFactor: number;
	stationPower: number;
	minimalRayPower: number;
	stationPos?: [number, number];
	singleRays: SingleRay[];
};

export type RayLaunchType = {
	x: number;
	y: number;
	z: number;
	power: number;
}[];

export type SingleMapDataTypes = {
	stationPos: mapboxgl.LngLatLike[];
	computationResult: mapboxgl.LngLatLike[][];
};

export type MapTypes = {
	title: string;
	coordinates: number[][][];
	center: mapboxgl.LngLatLike;
	bounds: mapboxgl.LngLatBoundsLike;
};

export type MapTypesExtended = MapTypes & {
	stationPos: mapboxgl.LngLatLike[];
	stationHeight: number;
	minimalRayPower: number;
	handleStationPosUpdate: (value: mapboxgl.LngLatLike) => void;
	buildingsData: GeoJSON.FeatureCollection | null;
	computationResult: mapboxgl.LngLatLike[][];
	wallMatrix: Float64Array;
	spherePositions: { coord: MercatorCoordinate; power: number }[] | undefined;
};

export type PostComputeTypes = {
	freq: string;
	stationH: string;
};
