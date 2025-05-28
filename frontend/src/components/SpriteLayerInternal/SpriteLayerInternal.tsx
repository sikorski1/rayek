import { LngLatLike, default as mapboxgl, MercatorCoordinate } from "mapbox-gl";
import React, { useEffect, useState } from "react";

import { FeatureCollection } from "geojson";
import * as THREE from "three";
import { GLTFLoader } from "three/examples/jsm/loaders/GLTFLoader.js";

export interface SpritePaint {
	gltfPath: string;
	/** Apply a scaling factor to the model's coordinates. After this, they should be in meters. */
	scale?: number;
	/** Rotate the model by the given amount along each axis. */
	rotateDeg?: {
		x?: number;
		y?: number;
		z?: number;
	};
}

export interface Props {
	id: string;
	spritePaint: SpritePaint;
	data: FeatureCollection;
	map: mapboxgl.Map;
}

interface Model {
	path: string;
	scale: number;
	rotate: number[];
}

interface Sprite {
	model: Model;
	position: LngLatLike;
	altitude: number;
}

// The approach in this file is based on this Mapbox GL demo:
// https://docs.mapbox.com/mapbox-gl-js/example/add-3d-model/

function getSpriteMatrix(sprite: Sprite, center: mapboxgl.MercatorCoordinate): THREE.Matrix4 {
	const { model, position, altitude } = sprite;
	const { scale, rotate } = model;
	const rotationX = new THREE.Matrix4().makeRotationAxis(new THREE.Vector3(1, 0, 0), rotate[0]);
	const rotationY = new THREE.Matrix4().makeRotationAxis(new THREE.Vector3(0, 1, 0), rotate[1]);
	const rotationZ = new THREE.Matrix4().makeRotationAxis(new THREE.Vector3(0, 0, 1), rotate[2]);

	const coord = MercatorCoordinate.fromLngLat(position, altitude);
	return new THREE.Matrix4()
		.makeTranslation(coord.x - center.x, coord.y - center.y, coord.z! - center.z!)
		.scale(new THREE.Vector3(scale, -scale, scale))
		.multiply(rotationX)
		.multiply(rotationY)
		.multiply(rotationZ);
}

/**
 * Load a 3D model and render it at specific Lat/Lngs.
 * This renders a THREE.js scene in the same WebGL canvas as Mapbox GL.
 */
class SpriteCustomLayer implements mapboxgl.CustomLayerInterface {
	type = "custom" as const;
	renderingMode = "3d" as const;

	id: string;
	options: SpritePaint;
	camera!: THREE.Camera;
	scene!: THREE.Scene;
	map!: mapboxgl.Map;
	renderer!: THREE.WebGLRenderer;
	center!: mapboxgl.MercatorCoordinate;
	cameraTransform!: THREE.Matrix4;
	model: Promise<THREE.Group>;
	modelConfig: Model;
	features: FeatureCollection | null;

	constructor(id: string, options: SpritePaint) {
		this.id = id;
		this.options = options;
		this.modelConfig = {
			path: options.gltfPath,
			scale: options.scale || 1,
			rotate: [
				options.rotateDeg ? options.rotateDeg.x || 0 : 0,
				options.rotateDeg ? options.rotateDeg.y || 0 : 0,
				options.rotateDeg ? options.rotateDeg.z || 0 : 0,
			].map(deg => (Math.PI / 180) * deg),
		};
		this.model = new Promise<THREE.Group>((resolve, reject) => {
			const loader = new GLTFLoader();
			loader.load(
				options.gltfPath,
				gltf => {
					resolve(gltf.scene);
				},
				() => {
					// progress is being made; bytes loaded = xhr.loaded / xhr.total
				},
				e => {
					console.log(e);
				}
			);
		});
		this.features = null;
	}

	onAdd(map: mapboxgl.Map, gl: WebGLRenderingContext) {
		this.camera = new THREE.Camera();

		this.center = MercatorCoordinate.fromLngLat(map.getCenter(), 0);
		const { x, y, z } = this.center;
		this.cameraTransform = new THREE.Matrix4().makeTranslation(x, y, z!);

		this.map = map;
		this.scene = this.makeScene();

		// use the Mapbox GL JS map canvas for three.js
		this.renderer = new THREE.WebGLRenderer({
			canvas: map.getCanvas(),
			context: gl,
			antialias: true,
            alpha: true
		});

		// From https://threejs.org/docs/#examples/en/loaders/GLTFLoader
		this.renderer.outputColorSpace = THREE.SRGBColorSpace;

		this.renderer.autoClear = false;
	}

	makeScene() {
		const scene = new THREE.Scene();

		// TODO(danvk): fiddle with lighting
		const ambientLight = new THREE.AmbientLight(0x916262, 0.5);
		scene.add(ambientLight);

		const light = new THREE.HemisphereLight(0xffffbb, 0x080820, 1);
		scene.add(light);
		return scene;
	}

	async setData(geojson: FeatureCollection) {
		this.features = geojson;
		const model = await this.model;
		if (this.features !== geojson) {
			return; // there was another call
		}

		this.scene = this.makeScene(); // clear the old scene
		const spriteScenes = geojson.features.map(f => {
			const { geometry } = f;
			if (geometry.type !== "Point") {
				throw new Error(`Sprite layers must have Point geometries; got ${f.geometry.type}`);
			}
			const { coordinates } = geometry;
			const scene = model.clone();
			scene.applyMatrix4(
				getSpriteMatrix(
					{
						model: this.modelConfig,
						position: {
							lng: coordinates[0],
							lat: coordinates[1],
						},
						altitude: f.properties?.altitude ?? 0,
					},
					this.center
				)
			);
			return scene;
		});

		for (const scene of spriteScenes) {
			this.scene.add(scene);
		}
	}

	render(gl: WebGLRenderingContext, matrix: number[]) {
		this.camera.projectionMatrix = new THREE.Matrix4().fromArray(matrix).multiply(this.cameraTransform);
		this.renderer.state.reset();
		this.renderer.render(this.scene, this.camera);
		this.map.triggerRepaint();
	}
}

const SpriteLayerInternal: React.FunctionComponent<Props> = props => {
	const { map, id, spritePaint, data } = props;
	const [spriteLayer, setSpriteLayer] = useState<SpriteCustomLayer | null>(null);

	useEffect(() => {
		const layer = new SpriteCustomLayer(id, spritePaint);
		map.addLayer(layer);
		setSpriteLayer(layer);

		return () => {
			map.removeLayer(id);
		};
	}, []);

	useEffect(() => {
		if (spriteLayer) {
			spriteLayer.setData(data);
		}
	}, [spriteLayer, data]);

	return null;
};

export const SpriteLayer = SpriteLayerInternal;
