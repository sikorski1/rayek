import { MapTypes } from "@/types/main";
import { FeatureCollection } from "geojson";
import mapboxgl, { CustomLayerInterface } from "mapbox-gl";
import "mapbox-gl/dist/mapbox-gl.css";
import { useEffect, useRef } from "react";
import * as THREE from "three";
import { GLTFLoader } from "three/examples/jsm/loaders/GLTFLoader.js";
export default function Map({ title, coordinates, center, bounds, stationPos }: MapTypes) {
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
	useEffect(() => {
		mapboxgl.accessToken = import.meta.env.VITE_MAPBOX_ACCESS_TOKEN;

		mapRef.current = new mapboxgl.Map({
			style: "mapbox://styles/mapbox/light-v11",
			center: center,
			zoom: 15.5,
			pitch: 45,
			bearing: -17.6,
			container: title,
			antialias: true,
		});

		const towerModelOrigin = stationPos;
		const towerModelAltitude = 0;
		const towerModelRotate = [Math.PI / 2, 0, 0];

		const towerModelAsMercatorCoordinate = mapboxgl.MercatorCoordinate.fromLngLat(towerModelOrigin!, towerModelAltitude);

		const towerModelTransform = {
			translateX: towerModelAsMercatorCoordinate.x,
			translateY: towerModelAsMercatorCoordinate.y,
			translateZ: towerModelAsMercatorCoordinate.z,
			rotateX: towerModelRotate[0],
			rotateY: towerModelRotate[1],
			rotateZ: towerModelRotate[2],
			scale: towerModelAsMercatorCoordinate.meterInMercatorCoordinateUnits() * 5,
		};

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
				onAdd: () => {
					// Add logic that runs on layer addition if necessary.
				},
				render: (gl:WebGLRenderingContext, matrix:THREE.Matrix4) => {
					const rotationX = new THREE.Matrix4().makeRotationAxis(new THREE.Vector3(1, 0, 0), towerModelTransform.rotateX);
					const rotationY = new THREE.Matrix4().makeRotationAxis(new THREE.Vector3(0, 1, 0), towerModelTransform.rotateY);
					const rotationZ = new THREE.Matrix4().makeRotationAxis(new THREE.Vector3(0, 0, 1), towerModelTransform.rotateZ);

					const m = new THREE.Matrix4().fromArray(matrix as unknown as ArrayLike<number>);
					const l = new THREE.Matrix4()
						.makeTranslation(towerModelTransform.translateX, towerModelTransform.translateY, towerModelTransform.translateZ)
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

			const customLayer = createCustomLayer();
			mapRef.current?.addLayer(customLayer as unknown as CustomLayerInterface, 'waterway-label');
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
	return <div id={title} ref={mapContainerRef} style={{ height: "100%", width: "100%" }}></div>;
}
