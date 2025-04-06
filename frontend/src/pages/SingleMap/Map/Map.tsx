import { MapTypesExtended } from "@/types/main";
import { FeatureCollection, Position } from "geojson";
import mapboxgl, { CustomLayerInterface } from "mapbox-gl";
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


export default function Map({
	title,
	coordinates,
	center,
	bounds,
	stationPos,
	handleStationPosUpdate,
	buildingsData,
	computationResult,
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
		function onMove(e: mapboxgl.MapMouseEvent | mapboxgl.MapTouchEvent) {
			const coords = e.lngLat;

			canvas.style.cursor = "grabbing";

			const pointGeometry = dragDropGeoJSON.features[0].geometry as GeoJSON.Point;
			let towerModelOrigin;
			if (checkBounds([parseFloat(coords.lng.toFixed(6)), parseFloat(coords.lat.toFixed(6))], bounds as number[][])) {
				pointGeometry.coordinates = [parseFloat(coords.lng.toFixed(6)), parseFloat(coords.lat.toFixed(6))];
				towerModelOrigin = [parseFloat(coords.lng.toFixed(6)), parseFloat(coords.lat.toFixed(6))];
			} else {
				pointGeometry.coordinates = center as Position;
				towerModelOrigin = center;
			}
			const source = mapRef.current?.getSource("point") as mapboxgl.GeoJSONSource;
			source.setData(dragDropGeoJSON);

			const towerModelAsMercatorCoordinate = mapboxgl.MercatorCoordinate.fromLngLat(
				towerModelOrigin as mapboxgl.LngLatLike,
				0
			);
			towerModelTransform.translateX = towerModelAsMercatorCoordinate.x;
			towerModelTransform.translateY = towerModelAsMercatorCoordinate.y;
			towerModelTransform.translateZ = towerModelAsMercatorCoordinate.z;
			towerModelTransform.scale = towerModelAsMercatorCoordinate.meterInMercatorCoordinateUnits() * 5;
		}

		function onUp(e: mapboxgl.MapMouseEvent | mapboxgl.MapTouchEvent) {
			const coords = e.lngLat;
			if (checkBounds([parseFloat(coords.lng.toFixed(6)), parseFloat(coords.lat.toFixed(6))], bounds as number[][])) {
				handleStationPosUpdate([parseFloat(coords.lng.toFixed(6)), parseFloat(coords.lat.toFixed(6))]);
			} else {
				handleStationPosUpdate(center);
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
		const towerModelAltitude = 0;
		const towerModelRotate = [Math.PI / 2, 0, 0];

		const towerModelAsMercatorCoordinate = mapboxgl.MercatorCoordinate.fromLngLat(
			towerModelOrigin!,
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
				id: "region-mask",
				type: "fill",
				source: "region-mask",
				paint: {
					"fill-color": "#333",
					"fill-opacity": 0.3,
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
	}, []);

	// useEffect(() => {
	// 	if (!computationResult) return;
	// 	const geojson = {
	// 		type: "FeatureCollection",
	// 		features: computationResult.map((position: mapboxgl.LngLatLike[])  => {
	// 			const feature  = {
	// 				type: "Feature",
	// 				properties: {},
	// 				geometry: {
	// 					type: "Point",
	// 					coordinates: position,
	// 				},
	// 			};
	// 			return feature
	// 		}),
	// 	};

	// 	mapRef.current?.on("style.load", () => {
	// 		mapRef.current?.addSource("spheres-2d", {
	// 			type: "geojson",
	// 			data: geojson as FeatureCollection,
	// 		});

	// 		mapRef.current?.addLayer({
	// 			id: "spheres-2d-layer",
	// 			type: "circle",
	// 			source: "spheres-2d",
	// 			paint: {
	// 				"circle-radius": 4,
	// 				"circle-color": "#d30000",
	// 			},
	// 		}, mapRef.current.getStyle()!.layers?.[30]?.id);
	// 	});
	// }, [computationResult]);
	return <div id={title} ref={mapContainerRef} style={{ height: "100%", width: "100%" }}></div>;
}
