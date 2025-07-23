import { MapTypesExtended } from "@/types/main";
import { geoToMatrixIndex } from "@/utils/geoToMatrixIndex";
import { getHeatMapColor } from "@/utils/getHeatMapColor";
import { getMatrixValue } from "@/utils/getMatrixValue";
import { FeatureCollection, Position } from "geojson";
import mapboxgl, { CustomLayerInterface, LngLatLike } from "mapbox-gl";
import "mapbox-gl/dist/mapbox-gl.css";
import { useEffect, useRef } from "react";
import * as THREE from "three";
import { GLTFLoader } from "three/examples/jsm/loaders/GLTFLoader.js";

const checkBounds = (coords: number[], bounds: number[][]) => {
	if (coords[0] > bounds[0][0] && coords[0] < bounds[1][0]) {
		if (coords[1] > bounds[0][1] && coords[1] < bounds[1][1]) {
			return true;
		} else return false;
	} else {
		return false;
	}
};

const normalizePower = (value: number) => {
	const minPowerDb = -160;
	const maxPowerDb = -0.1;

	const clamped = Math.max(minPowerDb, Math.min(maxPowerDb, value));

	return (clamped - minPowerDb) / (maxPowerDb - minPowerDb);
};

export default function Map({
	title,
	coordinates,
	center,
	bounds,
	size,
	stationPos,
	stationHeight,
	handleStationPosUpdate,
	buildingsData,
	spherePositions,
	wallMatrix,
	powerMap,
}: MapTypesExtended) {
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
					coordinates: coordinates,
				},
			},
		],
	};

	const dragDropGeoJSON: FeatureCollection = {
		type: "FeatureCollection",
		features: [
			{
				type: "Feature",
				properties: {},
				geometry: {
					type: "Point",
					coordinates: stationPos as Position,
				},
			},
		],
	};
	useEffect(() => {
		if (!powerMap || !mapRef.current?.isStyleLoaded()) return;

		const depth = powerMap.length;
		const layerIds: string[] = [];

		// Utwórz heatmapę dla każdej wysokości (co 5 poziomów dla wydajności)
		for (let z = 0; z < depth; z++) {
			const layerData = powerMap[z];
			const width = layerData[0].length;
			const height = layerData.length;

			// Stwórz canvas dla tej warstwy
			const canvas = document.createElement("canvas");
			canvas.width = width;
			canvas.height = height;
			const ctx = canvas.getContext("2d");

			if (!ctx) continue;

			const imageData = ctx.createImageData(width, height);
			const data = imageData.data;

			for (let y = 0; y < height; y++) {
				for (let x = 0; x < width; x++) {
					const flippedY = height - y - 1;
					const value = layerData[flippedY][x];
					const normalized = normalizePower(value);
					const color = getHeatMapColor(normalized);

					const index = (y * width + x) * 4;
					data[index] = color.r * 255;
					data[index + 1] = color.g * 255;
					data[index + 2] = color.b * 255;
					data[index + 3] = normalized * 180; // Trochę przezroczyste
				}
			}

			ctx.putImageData(imageData, 0, 0);
			const imageUrl = canvas.toDataURL();

			const sourceId = `power-layer-${z}`;
			const layerId = `power-heatmap-${z}`;

			mapRef.current.addSource(sourceId, {
				type: "image",
				url: imageUrl,
				coordinates: [
					[coordinates[0][0][0], coordinates[0][0][1]],
					[coordinates[0][1][0], coordinates[0][1][1]],
					[coordinates[0][2][0], coordinates[0][2][1]],
					[coordinates[0][3][0], coordinates[0][3][1]],
				],
			});

			mapRef.current.addLayer({
				id: layerId,
				type: "raster",
				source: sourceId,
				paint: {
					"raster-opacity": z === 0 ? 0.8 : 0, // Pierwsza warstwa najbardziej widoczna
				},
			});

			layerIds.push(layerId);
		}

		// Dodaj kontrolę dla przełączania warstw
		const addLayerControl = () => {
			const controlDiv = document.createElement("div");
			controlDiv.className = "mapboxgl-ctrl mapboxgl-ctrl-group";
			controlDiv.style.position = "absolute";
			controlDiv.style.top = "10px";
			controlDiv.style.right = "10px";
			controlDiv.style.background = "white";
			controlDiv.style.padding = "10px";
			controlDiv.style.borderRadius = "3px";

			const slider = document.createElement("input");
			slider.type = "range";
			slider.min = "0";
			slider.max = String(depth-1);
			slider.value = "0";
			slider.style.width = "150px";

			const label = document.createElement("div");
			label.textContent = `Wysokość: ${0 * 1}m`;
			label.style.marginBottom = "5px";
			label.style.fontSize = "12px";

			slider.addEventListener("input", e => {
				const selectedLayer = parseInt((e.target as HTMLInputElement).value);
				const heightInMeters = selectedLayer ; // każdy poziom to 2m, co 5 poziomów
				label.textContent = `Wysokość: ${heightInMeters}m`;

				// Ukryj wszystkie warstwy
				layerIds.forEach((layerId, index) => {
					const opacity = index === selectedLayer ? 0.8 : 0;
					mapRef.current?.setPaintProperty(layerId, "raster-opacity", opacity);
				});
			});

			controlDiv.appendChild(label);
			controlDiv.appendChild(slider);
			mapRef.current?.getContainer().appendChild(controlDiv);

			return controlDiv;
		};

		const control = addLayerControl();

		return () => {
			// Cleanup
			layerIds.forEach(layerId => {
				if (mapRef.current?.getLayer(layerId)) {
					mapRef.current.removeLayer(layerId);
				}
			});

			for (let z = 0; z < depth; z += 1) {
				const sourceId = `power-layer-${z}`;
				if (mapRef.current?.getSource(sourceId)) {
					mapRef.current.removeSource(sourceId);
				}
			}

			if (control && control.parentNode) {
				control.parentNode.removeChild(control);
			}
		};
	}, [powerMap]);
	useEffect(() => {
		const createSphereLayer = (
			position: mapboxgl.MercatorCoordinate,
			power: number,
			rayIndex: number,
			pointIndex: number
		) => {
			const camera = new THREE.Camera();
			const scene = new THREE.Scene();
			const sphereGeometry = new THREE.SphereGeometry(1, 8, 8);
			const minPower = -160;
			const maxPower = -0.1;

			const clampedPower = Math.min(Math.max(power, minPower), maxPower);
			const normalizedPower = (clampedPower - minPower) / (maxPower - minPower);
			const color = getHeatMapColor(normalizedPower);
			const sphereMaterial = new THREE.MeshBasicMaterial({
				color,
				transparent: true,
				opacity: 0.8,
			});

			const sphereMesh = new THREE.Mesh(sphereGeometry, sphereMaterial);
			scene.add(sphereMesh);

			const sphereModelTransform = {
				translateX: position.x,
				translateY: position.y,
				translateZ: position.z,
				rotateX: 0,
				rotateY: 0,
				rotateZ: 0,
				scale: position.meterInMercatorCoordinateUnits(),
			};

			const renderer = new THREE.WebGLRenderer({
				canvas: mapRef.current?.getCanvas(),
				context: mapRef.current?.painter.context.gl,
				antialias: true,
			});
			renderer.autoClear = false;

			return {
				id: `sphere-ray${rayIndex}-point${pointIndex}`,
				type: "custom",
				renderingMode: "3d",
				onAdd: () => {},
				render: (gl: WebGLRenderingContext, matrix: THREE.Matrix4) => {
					const m = new THREE.Matrix4().fromArray(matrix as unknown as ArrayLike<number>);
					const l = new THREE.Matrix4()
						.makeTranslation(
							sphereModelTransform.translateX,
							sphereModelTransform.translateY,
							sphereModelTransform.translateZ
						)
						.scale(
							new THREE.Vector3(sphereModelTransform.scale, sphereModelTransform.scale, sphereModelTransform.scale)
						);
					camera.projectionMatrix = m.multiply(l);
					renderer.resetState();
					renderer.render(scene, camera);
					mapRef.current?.triggerRepaint();
				},
			};
		};

		if (!spherePositions || spherePositions.length === 0) return;

		const addSphereLayers = () => {
			spherePositions.forEach(({ rayIndex, positions }) => {
				positions.forEach(({ coord, power }, pointIndex) => {
					const customLayer = createSphereLayer(coord, power, rayIndex, pointIndex);
					mapRef.current?.addLayer(customLayer as unknown as CustomLayerInterface, "waterway-label");
				});
			});
		};

		if (mapRef.current?.isStyleLoaded()) {
			addSphereLayers();
		} else {
			mapRef.current?.on("style.load", addSphereLayers);
		}

		return () => {
			if (mapRef.current) {
				spherePositions?.forEach(({ rayIndex, positions }) => {
					positions.forEach((_, pointIndex) => {
						const layerId = `sphere-ray${rayIndex}-point${pointIndex}`;
						if (mapRef.current?.getLayer(layerId)) {
							mapRef.current.removeLayer(layerId);
						}
					});
				});
			}
		};
	}, [spherePositions]);
	useEffect(() => {
		let lastValidCoords: [number, number] | null = null;
		function onMove(e: mapboxgl.MapMouseEvent | mapboxgl.MapTouchEvent) {
			const coords = e.lngLat;
			const { i, j } = geoToMatrixIndex(
				parseFloat(coords.lng.toFixed(6)),
				parseFloat(coords.lat.toFixed(6)),
				coordinates[0][0][0],
				coordinates[0][2][0],
				coordinates[0][0][1],
				coordinates[0][2][1],
				size
			);
			const value = getMatrixValue(wallMatrix, i, j, Number(stationHeight), size, size, 30);
			canvas.style.cursor = "grabbing";

			const pointGeometry = dragDropGeoJSON.features[0].geometry as GeoJSON.Point;
			let towerModelOrigin: [number, number];

			if (checkBounds([coords.lng, coords.lat], bounds as number[][]) && value === -150) {
				towerModelOrigin = [parseFloat(coords.lng.toFixed(6)), parseFloat(coords.lat.toFixed(6))];
				lastValidCoords = towerModelOrigin;
			} else {
				towerModelOrigin = lastValidCoords ?? (center as [number, number]);
			}

			pointGeometry.coordinates = towerModelOrigin;

			const source = mapRef.current?.getSource("point") as mapboxgl.GeoJSONSource;
			source.setData(dragDropGeoJSON);

			const towerModelAsMercatorCoordinate = mapboxgl.MercatorCoordinate.fromLngLat(towerModelOrigin, stationHeight);
			towerModelTransform.translateX = towerModelAsMercatorCoordinate.x;
			towerModelTransform.translateY = towerModelAsMercatorCoordinate.y;
			towerModelTransform.translateZ = towerModelAsMercatorCoordinate.z;
			towerModelTransform.scale = towerModelAsMercatorCoordinate.meterInMercatorCoordinateUnits() * 5;
		}

		function onUp(e: mapboxgl.MapMouseEvent | mapboxgl.MapTouchEvent) {
			if (lastValidCoords) {
				handleStationPosUpdate(lastValidCoords);
			}
			canvas.style.cursor = "";
			mapRef.current?.off("mousemove", onMove);
			mapRef.current?.off("touchmove", onMove);
		}
		const createCustomLayer = () => {
			const camera = new THREE.Camera();
			const scene = new THREE.Scene();

			const directionalLight1 = new THREE.DirectionalLight(0xffffff);
			directionalLight1.position.set(0, -70, 100).normalize();
			scene.add(directionalLight1);

			const directionalLight2 = new THREE.DirectionalLight(0xffffff);
			directionalLight2.position.set(0, 70, 100).normalize();
			scene.add(directionalLight2);

			const loader = new GLTFLoader();
			loader.load("/3dmodels/5_g_tower/scene.gltf", gltf => {
				scene.add(gltf.scene);
			});

			const renderer = new THREE.WebGLRenderer({
				canvas: mapRef.current?.getCanvas(),
				context: mapRef.current?.painter.context.gl,
				antialias: true,
			});
			renderer.autoClear = false;
			return {
				id: "3d-5gtower",
				type: "custom",
				renderingMode: "3d",
				onAdd: () => {},
				render: (gl: WebGLRenderingContext, matrix: THREE.Matrix4) => {
					const rotationX = new THREE.Matrix4().makeRotationAxis(
						new THREE.Vector3(1, 0, 0),
						towerModelTransform.rotateX
					);
					const rotationY = new THREE.Matrix4().makeRotationAxis(
						new THREE.Vector3(0, 1, 0),
						towerModelTransform.rotateY
					);
					const rotationZ = new THREE.Matrix4().makeRotationAxis(
						new THREE.Vector3(0, 0, 1),
						towerModelTransform.rotateZ
					);

					const m = new THREE.Matrix4().fromArray(matrix as unknown as ArrayLike<number>);
					const l = new THREE.Matrix4()
						.makeTranslation(
							towerModelTransform.translateX,
							towerModelTransform.translateY,
							towerModelTransform.translateZ
						)
						.scale(new THREE.Vector3(towerModelTransform.scale, -towerModelTransform.scale, towerModelTransform.scale))
						.multiply(rotationX)
						.multiply(rotationY)
						.multiply(rotationZ);

					camera.projectionMatrix = m.multiply(l);
					renderer.resetState();
					renderer.render(scene, camera);
					mapRef.current?.triggerRepaint();
				},
			};
		};

		mapboxgl.accessToken = import.meta.env.VITE_MAPBOX_ACCESS_TOKEN;
		const towerModelOrigin = stationPos;
		const towerModelAltitude = stationHeight;
		const towerModelRotate = [Math.PI / 2, 0, 0];

		const towerModelAsMercatorCoordinate = mapboxgl.MercatorCoordinate.fromLngLat(
			towerModelOrigin as unknown as LngLatLike,
			towerModelAltitude
		);
		const towerModelTransform = {
			translateX: towerModelAsMercatorCoordinate.x,
			translateY: towerModelAsMercatorCoordinate.y,
			translateZ: towerModelAsMercatorCoordinate.z,
			rotateX: towerModelRotate[0],
			rotateY: towerModelRotate[1],
			rotateZ: towerModelRotate[2],
			scale: towerModelAsMercatorCoordinate.meterInMercatorCoordinateUnits() * 5,
		};

		mapRef.current = new mapboxgl.Map({
			style: "mapbox://styles/mapbox/light-v11",
			center: center,
			zoom: 15.5,
			pitch: 45,
			bearing: -17.6,
			container: title,
			antialias: true,
		});
		const canvas = mapRef.current.getCanvasContainer();

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
			mapRef.current?.addSource("point", {
				type: "geojson",
				data: dragDropGeoJSON,
			});
			mapRef.current?.addSource("buildings-polygon", {
				type: "geojson",
				data: buildingsData!,
			});
			//add region-mask layer

			mapRef.current?.addLayer({
				id: "region-mask-border",
				type: "line",
				source: "region-mask",
				paint: {
					"line-color": "#000",
					"line-width": 1,
				},
			});
			//add dragdrop circle point layer
			mapRef.current?.addLayer({
				id: "point",
				type: "circle",
				source: "point",
				paint: {
					"circle-radius": 5,
					"circle-color": "#F84C4C",
				},
			});
			//add buildings walls polygon layer
			mapRef.current?.addLayer({
				id: "buildings-polygon",
				type: "line",
				source: "buildings-polygon",
				paint: {
					"line-color": "#000",
					"line-width": 1,
				},
			});

			//add 3d 5gtower model layer
			const customLayer = createCustomLayer();
			mapRef.current?.addLayer(customLayer as unknown as CustomLayerInterface, "waterway-label");
			//add map 3d buildings layer
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
			mapRef.current?.on("mouseenter", "point", () => {
				mapRef.current?.setPaintProperty("point", "circle-color", "#3bb2d0");
				canvas.style.cursor = "move";
			});

			mapRef.current?.on("mouseleave", "point", () => {
				mapRef.current?.setPaintProperty("point", "circle-color", "#3887be");
				canvas.style.cursor = "";
			});

			mapRef.current?.on("mousedown", "point", e => {
				e.preventDefault();
				canvas.style.cursor = "grab";
				mapRef.current?.on("mousemove", onMove);
				mapRef.current?.once("mouseup", onUp);
			});

			mapRef.current?.on("touchstart", "point", e => {
				if (e.points.length !== 1) return;
				e.preventDefault();
				mapRef.current?.on("touchmove", onMove);
				mapRef.current?.once("touchend", onUp);
			});
		});
		return () => mapRef.current?.remove();
	}, [stationHeight]);
	return <div id={title} ref={mapContainerRef} style={{ height: "100%", width: "100%" }}></div>;
}
