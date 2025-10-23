import PowerInfoMarker from "@/components/PowerMapInfo/PowerInfoMarker";
import PowerInfoPanel from "@/components/PowerMapInfo/PowerInfoPanel";
import { MapTypesExtended } from "@/types/main";
import { geoToMatrixIndex } from "@/utils/geoToMatrixIndex";
import { getHeatMapColor } from "@/utils/getHeatMapColor";
import { getMatrixValue } from "@/utils/getMatrixValue";
import { normalizePower } from "@/utils/normalizePower";
import { FeatureCollection, Position } from "geojson";
import mapboxgl, { CustomLayerInterface, LngLatLike } from "mapbox-gl";
import "mapbox-gl/dist/mapbox-gl.css";
import { useEffect, useRef, useState } from "react";
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
	powerMapHeight,
	isPowerMapVisible,
}: MapTypesExtended) {
	const [clickedPowerInfo, setClickedPowerInfo] = useState<{
		lng: number;
		lat: number;
		power: number;
		x: number;
		y: number;
		buildingHeight: number | null;
	} | null>(null);
	const mapContainerRef = useRef<HTMLDivElement | null>(null);
	const mapRef = useRef<mapboxgl.Map | null>(null);
	const rendererRef = useRef<THREE.WebGLRenderer | null>(null);
	const powerMapUpdateTimeoutRef = useRef<NodeJS.Timeout>();
	const renderStateRef = useRef({
		isRendering: false,
		pendingRenders: new Set<Function>(),
	});

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
					coordinates: stationPos as unknown as Position,
				},
			},
		],
	};

	const getSharedRenderer = () => {
		if (!rendererRef.current && mapRef.current) {
			rendererRef.current = new THREE.WebGLRenderer({
				canvas: mapRef.current.getCanvas(),
				context: mapRef.current.painter.context.gl,
				antialias: true,
			});
			rendererRef.current.autoClear = false;
		}
		return rendererRef.current;
	};

	const createSafeRenderFunction = (originalRender: Function, _layerId: string) => {
		return (gl: WebGLRenderingContext, matrix: THREE.Matrix4) => {
			if (renderStateRef.current.isRendering) {
				renderStateRef.current.pendingRenders.add(() => originalRender(gl, matrix));
				return;
			}

			renderStateRef.current.isRendering = true;

			try {
				originalRender(gl, matrix);
			} finally {
				renderStateRef.current.isRendering = false;

				const pending = Array.from(renderStateRef.current.pendingRenders);
				renderStateRef.current.pendingRenders.clear();
				pending.forEach(render => render());
			}
		};
	};

	useEffect(() => {
		if (powerMapUpdateTimeoutRef.current) {
			clearTimeout(powerMapUpdateTimeoutRef.current);
		}

		powerMapUpdateTimeoutRef.current = setTimeout(() => {
			if (!powerMap || !mapRef.current?.isStyleLoaded()) return;

			for (let z = 0; z < 50; z++) {
				const layerId = `power-heatmap-${z}`;
				if (mapRef.current.getLayer(layerId)) {
					mapRef.current.removeLayer(layerId);
				}
				const sourceId = `power-layer-${z}`;
				if (mapRef.current.getSource(sourceId)) {
					mapRef.current.removeSource(sourceId);
				}
			}

			const depth = powerMap.length;
			const layerIds: string[] = [];

			const visibleLayers = [powerMapHeight];

			visibleLayers.forEach(z => {
				if (z >= 0 && z < depth) {
					const layerData = powerMap[z];
					const width = layerData[0].length;
					const height = layerData.length;

					const canvas = document.createElement("canvas");
					canvas.width = width;
					canvas.height = height;
					const ctx = canvas.getContext("2d");
					if (!ctx) return;

					const imageData = ctx.createImageData(width, height);
					const data = imageData.data;

					for (let y = 0; y < height; y++) {
						for (let x = 0; x < width; x++) {
							const flippedY = height - y - 1;
							const value = layerData[flippedY][x];

							const index = (y * width + x) * 4;
							if (value === -160) {
								data[index] = 255;
								data[index + 1] = 255;
								data[index + 2] = 255;
								data[index + 3] = 80; 
							} else if (value > 0) {
								data[index] = 255;
								data[index + 1] = 0;
								data[index + 2] = 0;
								data[index + 3] = 120; 
							} else {
								const normalized = normalizePower(value);
								const color = getHeatMapColor(normalized);
								data[index] = color.r * 255;
								data[index + 1] = color.g * 255;
								data[index + 2] = color.b * 255;
								const alpha =
									normalized < 0.3
										? 60 + normalized * 200 
										: 120 + normalized * 60; 
								data[index + 3] = alpha;
							}
						}
					}

					ctx.putImageData(imageData, 0, 0);
					const imageUrl = canvas.toDataURL();

					const sourceId = `power-layer-${z}`;
					const layerId = `power-heatmap-${z}`;

					mapRef.current?.addSource(sourceId, {
						type: "image",
						url: imageUrl,
						coordinates: [
							[coordinates[0][0][0], coordinates[0][0][1]],
							[coordinates[0][1][0], coordinates[0][1][1]],
							[coordinates[0][2][0], coordinates[0][2][1]],
							[coordinates[0][3][0], coordinates[0][3][1]],
						],
					});

					mapRef.current?.addLayer(
						{
							id: layerId,
							type: "raster",
							source: sourceId,
							paint: {
								"raster-opacity": 0.9,
							},
						},
						mapRef.current?.getStyle()?.layers?.[37]?.id
					);

					layerIds.push(layerId);
				}
			});
		}, 500);

		return () => {
			if (powerMapUpdateTimeoutRef.current) {
				clearTimeout(powerMapUpdateTimeoutRef.current);
			}
		};
	}, [powerMap, powerMapHeight]);
	useEffect(() => {
		if (!mapRef.current || !powerMap) return;

		for (let i = 0; i < powerMap.length; i++) {
			const layerId = `power-heatmap-${i}`;
			const opacity = isPowerMapVisible && i === powerMapHeight ? 0.8 : 0;

			if (mapRef.current.getLayer(layerId)) {
				mapRef.current.setPaintProperty(layerId, "raster-opacity", opacity);
			}
		}
	}, [isPowerMapVisible, powerMapHeight]);
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

			return {
				id: `sphere-ray${rayIndex}-point${pointIndex}`,
				type: "custom",
				renderingMode: "3d",
				onAdd: () => {},
				render: createSafeRenderFunction((_gl: WebGLRenderingContext, matrix: THREE.Matrix4) => {
					const renderer = getSharedRenderer();
					if (!renderer) return;

					const m = new THREE.Matrix4().fromArray(matrix as unknown as ArrayLike<number>);
					const l = new THREE.Matrix4()
						.makeTranslation(
							sphereModelTransform.translateX,
							sphereModelTransform.translateY,
							sphereModelTransform.translateZ
						)
						.scale(
							new THREE.Vector3(
								sphereModelTransform.scale,
								sphereModelTransform.scale,
								sphereModelTransform.scale
							)
						);
					camera.projectionMatrix = m.multiply(l);
					renderer.resetState();
					renderer.render(scene, camera);
					mapRef.current?.triggerRepaint();
				}, `sphere-ray${rayIndex}-point${pointIndex}`),
			};
		};

		if (!spherePositions || spherePositions.length === 0) return;

		const addSphereLayers = () => {
			spherePositions.forEach(({ rayIndex, positions }) => {
				positions.forEach(({ coord, power }, pointIndex) => {
					const customLayer = createSphereLayer(coord, power, rayIndex, pointIndex);
					mapRef.current?.addLayer(customLayer as unknown as CustomLayerInterface);
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

			if (checkBounds([coords.lng, coords.lat], bounds as number[][]) && value === -160) {
				towerModelOrigin = [parseFloat(coords.lng.toFixed(6)), parseFloat(coords.lat.toFixed(6))];
				lastValidCoords = towerModelOrigin;
			} else {
				towerModelOrigin = lastValidCoords ?? (center as [number, number]);
			}

			pointGeometry.coordinates = towerModelOrigin;

			const source = mapRef.current?.getSource("point") as mapboxgl.GeoJSONSource;
			source.setData(dragDropGeoJSON);

			const towerModelAsMercatorCoordinate = mapboxgl.MercatorCoordinate.fromLngLat(
				towerModelOrigin,
				stationHeight
			);
			towerModelTransform.translateX = towerModelAsMercatorCoordinate.x;
			towerModelTransform.translateY = towerModelAsMercatorCoordinate.y;
			towerModelTransform.translateZ = towerModelAsMercatorCoordinate.z;
			towerModelTransform.scale = towerModelAsMercatorCoordinate.meterInMercatorCoordinateUnits() * 5;
		}

		function onUp(_e: mapboxgl.MapMouseEvent | mapboxgl.MapTouchEvent) {
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

			return {
				id: "3d-5gtower",
				type: "custom",
				renderingMode: "3d",
				onAdd: () => {},
				render: createSafeRenderFunction((_gl: WebGLRenderingContext, matrix: THREE.Matrix4) => {
					const renderer = getSharedRenderer();
					if (!renderer) return;

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
						.scale(
							new THREE.Vector3(
								towerModelTransform.scale,
								-towerModelTransform.scale,
								towerModelTransform.scale
							)
						)
						.multiply(rotationX)
						.multiply(rotationY)
						.multiply(rotationZ);

					camera.projectionMatrix = m.multiply(l);
					renderer.resetState();
					renderer.render(scene, camera);
					mapRef.current?.triggerRepaint();
				}, "3d-5gtower"),
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
					layer =>
						layer.type === "symbol" && layer.layout && typeof layer.layout["text-field"] !== "undefined"
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
						"fill-extrusion-base": [
							"interpolate",
							["linear"],
							["zoom"],
							15,
							0,
							15.05,
							["get", "min_height"],
						],
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

		return () => {
			// Cleanup renderer
			if (rendererRef.current) {
				rendererRef.current.dispose();
				rendererRef.current = null;
			}
			// Cleanup timeouts
			if (powerMapUpdateTimeoutRef.current) {
				clearTimeout(powerMapUpdateTimeoutRef.current);
			}
			// Cleanup map
			mapRef.current?.remove();
		};
	}, [stationHeight]);

	useEffect(() => {
		const map = mapRef.current;
		if (!map) return;

		const handleMapClick = (e: mapboxgl.MapMouseEvent) => {
			const coords = e.lngLat;
			const lng = parseFloat(coords.lng.toFixed(6));
			const lat = parseFloat(coords.lat.toFixed(6));

			if (!checkBounds([lng, lat], bounds as number[][])) {
				setClickedPowerInfo(null);
				return;
			}

			const { i, j } = geoToMatrixIndex(
				lng,
				lat,
				coordinates[0][0][0],
				coordinates[0][2][0],
				coordinates[0][0][1],
				coordinates[0][2][1],
				size
			);
			let buildingHeight: number | null = null;
			if (wallMatrix && i >= 0 && j >= 0 && i < size && j < size) {
				const depth = 30;
				for (let z = depth - 1; z >= 0; z--) {
					const value = getMatrixValue(wallMatrix, i, j, z, size, size, depth);
					if (value && value >= 5000) {
						buildingHeight = z;
						break;
					}
				}
			}
			if (!powerMap || powerMap.length === 0) {
				setClickedPowerInfo({
					lng,
					lat,
					power: -160,
					x: i,
					y: j,
					buildingHeight,
				});
			} else {
				if (powerMapHeight >= 0 && powerMapHeight < powerMap.length) {
					const layerData = powerMap[powerMapHeight];

					if (j >= 0 && j < layerData.length && i >= 0 && i < layerData[0].length) {
						const power = layerData[j][i];

						setClickedPowerInfo({
							lng,
							lat,
							power,
							x: i,
							y: j,
							buildingHeight,
						});
					}
				}
			}
		};

		if (map.isStyleLoaded()) {
			map.on("click", handleMapClick);
		} else {
			map.once("load", () => {
				map.on("click", handleMapClick);
			});
		}

		return () => {
			if (map && map.loaded()) {
				map.off("click", handleMapClick);
			}
		};
	}, [powerMap, powerMapHeight, wallMatrix, stationHeight, coordinates]);
	return (
		<div style={{ position: "relative", height: "100%", width: "100%" }}>
			<div id={title} ref={mapContainerRef} style={{ height: "100%", width: "100%" }}></div>

			<PowerInfoMarker map={mapRef.current} powerInfo={clickedPowerInfo} />

			<PowerInfoPanel
				powerInfo={clickedPowerInfo}
				height={powerMapHeight}
				onClose={() => setClickedPowerInfo(null)}
			/>
		</div>
	);
}
