import mapboxgl from "mapbox-gl";
import { useEffect, useRef } from "react";
import { FeatureCollection } from "geojson";
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
							[-74.25909, 40.477399], // Southwest
							[-73.700272, 40.477399], // Southeast
							[-73.700272, 40.917576], // Northeast
							[-74.25909, 40.917576], // Northwest
							[-74.25909, 40.477399], // Zamknięcie pętli
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
			center: [-74.0066, 40.7135],
			zoom: 15.5,
			pitch: 45,
			bearing: -17.6,
			container: "map",
			antialias: true,
		});

        const bounds: mapboxgl.LngLatBoundsLike = [
			[-74.25909, 40.477399], // Southwest corner (dolny lewy róg)
			[-73.700272, 40.917576], // Northeast corner (górny prawy róg)
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
	return <div id="map" ref={mapContainerRef} style={{ height: "100%" }}></div>;
}
