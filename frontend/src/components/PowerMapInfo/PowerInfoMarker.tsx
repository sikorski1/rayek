import { getHeatMapColor } from "@/utils/getHeatMapColor";
import { normalizePower } from "@/utils/normalizePower";
import { useEffect } from "react";
import { Color } from "three";
interface PowerInfoMarkerProps {
	map: mapboxgl.Map | null;
	powerInfo: {
		lng: number;
		lat: number;
		power: number;
		x: number;
		y: number;
		buildingHeight: number | null;
	} | null;
}
export default function PowerInfoMarker({ map, powerInfo }: PowerInfoMarkerProps) {
	useEffect(() => {
		if (!map || !powerInfo) return;
		if ((map as any)._removed) return;

		let color: Color;
		if (powerInfo.power > 0 || powerInfo.power <= -160) {
			color = new Color(1, 1, 1);
		} else {
			const normalized = normalizePower(powerInfo.power);
			color = getHeatMapColor(normalized);
		}

		const layerId = "power-info-marker";
		const sourceId = "power-info-marker";

		try {
			if (map.getLayer(layerId)) map.removeLayer(layerId);
			if (map.getSource(sourceId)) map.removeSource(sourceId);

			if (!map.isStyleLoaded()) return;

			map.addSource(sourceId, {
				type: "geojson",
				data: {
					type: "Feature",
					properties: {},
					geometry: {
						type: "Point",
						coordinates: [powerInfo.lng, powerInfo.lat],
					},
				},
			});

			map.addLayer({
				id: layerId,
				type: "circle",
				source: sourceId,
				paint: {
					"circle-radius": 6,
					"circle-color": `rgb(${color.r * 255}, ${color.g * 255}, ${color.b * 255})`,
					"circle-stroke-width": 2,
					"circle-stroke-color": "#000000",
				},
			});
		} catch (e) {
			console.warn("PowerInfoMarker render aborted (map unmounted or removed):", e);
		}

		return () => {
			if (!map || (map as any)._removed) return;
			try {
				if (map.getLayer(layerId)) map.removeLayer(layerId);
				if (map.getSource(sourceId)) map.removeSource(sourceId);
			} catch {
				console.warn("Map has been destroyed");
			}
		};
	}, [map, powerInfo]);

	return null;
}
