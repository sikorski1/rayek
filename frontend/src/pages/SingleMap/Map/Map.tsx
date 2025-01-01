import { FeatureCollection } from "geojson";
import mapboxgl from "mapbox-gl";
import "mapbox-gl/dist/mapbox-gl.css";
import { useEffect, useRef } from "react";
export default function Map() {
	const mapContainerRef = useRef<HTMLDivElement | null>(null);
	const mapRef = useRef<mapboxgl.Map | null>(null);
	const regionGeoJSON: FeatureCollection = {
		type: "FeatureCollection",
		features: [
			{
				type: "Feature",
				properties: {},
				geometry: {
					type: "Polygon",
					coordinates: [
						[
							[19.914029, 50.065311], // Southwest
							[19.917527, 50.065311], // Southeast
							[19.917527, 50.067556], // Northeast
							[19.914029, 50.067556], // Northwest
							[19.914029, 50.065311], // Zamknięcie pętli
						],
					],
				},
			},
		],
	};
	useEffect(() => {
		mapboxgl.accessToken =
			"pk.eyJ1Ijoic2lrb3Jza2kxIiwiYSI6ImNtNWNjNW9vaDJycXYyanNnNjA2MG5rdGwifQ.w5aqAQkqzWFtkJF1_JelGg";

		mapRef.current = new mapboxgl.Map({
			style: "mapbox://styles/mapbox/light-v11",
			center: [19.915778, 50.0664335],
			zoom: 15.5,
			pitch: 45,
			bearing: -17.6,
			container: "map",
			antialias: true,
		});

		const bounds: mapboxgl.LngLatBoundsLike = [
			[19.914029, 50.065311], // Southwest corner (dolny lewy róg)
			[19.917527, 50.067556], // Northeast corner (górny prawy róg)
		];
		mapRef.current.fitBounds(bounds, { padding: 20 });

		mapRef.current.on("style.load", () => {
			const layers = mapRef.current?.getStyle()?.layers;
			let labelLayerId;

			if (layers) {
				labelLayerId = layers.find(
					layer => layer.type === "symbol" && layer.layout && typeof layer.layout["text-field"] !== "undefined"
				)?.id;
			}
			mapRef.current?.addSource("region-mask", {
				type: "geojson",
				data: regionGeoJSON,
			});

			mapRef.current?.addLayer({
				id: "region-mask",
				type: "fill",
				source: "region-mask",
				paint: {
					"fill-color": "#333",
					"fill-opacity": 0.3,
				},
			});

			mapRef.current?.addLayer(
				{
					id: "add-3d-buildings",
					source: "composite",
					"source-layer": "building",
					filter: ["==", "extrude", "true"],
					type: "fill-extrusion",
					minzoom: 15,
					paint: {
						"fill-extrusion-color": "#aaa",
						"fill-extrusion-height": ["interpolate", ["linear"], ["zoom"], 15, 0, 15.05, ["get", "height"]],
						"fill-extrusion-base": ["interpolate", ["linear"], ["zoom"], 15, 0, 15.05, ["get", "min_height"]],
						"fill-extrusion-opacity": 0.6,
					},
				},
				labelLayerId
			);
		});
		return () => mapRef.current?.remove();
	}, []);
	return <div id="map" ref={mapContainerRef} style={{ height: "100%", width:"100%"}}></div>;
}
