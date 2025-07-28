import SpinWifi from "@/components/Loaders/SpinWifi";
import GlobalSettings from "@/components/Modal/GlobalSettings/GlobalSettings";
import Modal from "@/components/Modal/Modal";
import SingleRaySettings from "@/components/Modal/SingleRaySettings/SingleRaySettings";
import { useGetMapById, useRayLaunching } from "@/hooks/useMap";
import Map from "@/pages/SingleMap/Map/Map";
import { RayLaunchType, SettingsDataTypes } from "@/types/main";
import { geoToMatrixIndex } from "@/utils/geoToMatrixIndex";
import { getMatrixValue } from "@/utils/getMatrixValue";
import { loadWallMatrix } from "@/utils/loadWallMatrix";
import { matrixIndexToGeo } from "@/utils/matrixIndexToGeo";
import { AnimatePresence, motion } from "framer-motion";
import mapboxgl from "mapbox-gl";
import { useEffect, useMemo, useState } from "react";
import { IoMdSettings } from "react-icons/io";
import { useParams } from "react-router-dom";
import styles from "./singleMap.module.scss";
const initialSettingsData: SettingsDataTypes = {
	settingsType: "global",
	isOpen: false,
	frequency: 5,
	stationHeight: 5,
	numberOfRaysAzimuth: 360,
	numberOfRaysElevation: 360,
	numberOfInteractions: 5,
	reflectionFactor: 0.5,
	stationPower: 5,
	minimalRayPower: -120,
	singleRays: [],
	powerMapHeight: 5,
};

export default function SingleMap() {
	const [settingsData, setSettingsData] = useState<SettingsDataTypes>(initialSettingsData);
	const [wallMatrix, setWallMatrix] = useState<Float64Array | null>(null);
	const [rayLaunchData, setRayLaunchData] = useState<RayLaunchType | null>(null);
	const [powerMapData, setPowerMapData] = useState<any>(null);
	const { id } = useParams();
	const { data, isLoading, error } = useGetMapById(id!);
	const handleStationPosUpdate = (stationPos: mapboxgl.LngLatLike) => {
		setSettingsData(prev => {
			const updatedSingleMapData = { ...prev, stationPos: stationPos };
			return updatedSingleMapData;
		});
	};

	const handleOnSettingsClose = () => {
		setSettingsData(prevSettingsData => {
			const updatedsettingsData = { ...prevSettingsData, isOpen: false };
			return updatedsettingsData;
		});
	};
	const handleAddRay = () => {
		if (settingsData.singleRays.length >= 4) return;

		setSettingsData(prev => ({
			...prev,
			singleRays: [...prev.singleRays, { azimuth: 0, elevation: 0 }],
		}));
	};
	const handleRemoveRay = (index: number) => {
		setSettingsData(prev => ({
			...prev,
			singleRays: prev.singleRays.filter((_, i) => i !== index),
		}));
	};
	const handleGlobalSettingsSubmit = (event: React.FormEvent<HTMLFormElement>) => {
		event.preventDefault();

		const form = event.currentTarget;
		const formData = new FormData(form);

		const updatedSettingsData: Omit<SettingsDataTypes, "settingsType" | "singleRays" | "isOpen"> = {
			numberOfRaysAzimuth: Number(formData.get("raysAzimuth")),
			numberOfRaysElevation: Number(formData.get("raysElevation")),
			frequency: Number(formData.get("frequency")),
			stationHeight: Number(formData.get("stationHeight")),
			reflectionFactor: Number(formData.get("relfectionFactor")),
			numberOfInteractions: Number(formData.get("interactions")),
			stationPower: Number(formData.get("stationPower")),
			minimalRayPower: Number(formData.get("minimalRayPower")),
		};
		if (
			updatedSettingsData.numberOfRaysAzimuth !== settingsData.numberOfRaysAzimuth ||
			updatedSettingsData.numberOfRaysElevation !== settingsData.numberOfRaysElevation
		) {
			setSettingsData(prev => ({ ...prev, ...updatedSettingsData, singleRays: [] }));
		} else {
			setSettingsData(prev => ({ ...prev, ...updatedSettingsData }));
		}
	};
	const handleSingleRaySettingsSubmit = (event: React.FormEvent<HTMLFormElement>) => {
		event.preventDefault();

		const form = event.currentTarget;
		const formData = new FormData(form);

		const newSingleRays = settingsData.singleRays.map((_, index) => ({
			azimuth: Number(formData.getAll("azimuth")[index]),
			elevation: Number(formData.getAll("elevation")[index]),
		}));

		setSettingsData(prev => ({
			...prev,
			singleRays: newSingleRays,
		}));
	};

	const handleOnSuccess = (data: any) => {
		setRayLaunchData(data.rayPaths);
		setPowerMapData(data.powerMap);
	};
	const { mutate, isPending: isPendingRayLaunch } = useRayLaunching(handleOnSuccess);
	const handleComputeBtn = async () => {
		const { stationPos, stationHeight, isOpen, ...restData } = settingsData;
		mutate({
			mapTitle: id!,
			configData: { stationPos: { x: i, y: j, z: Number(stationHeight) }, size: data.mapData.size, ...restData },
		});
	};
	useEffect(() => {
		if (!data) return;
		setSettingsData(prev => ({ ...prev, stationPos: data.mapData.center }));
	}, [data]);

	useEffect(() => {
		if (!id) return;
		loadWallMatrix(id).then(setWallMatrix);
	}, [id]);

	const { matrixIndexValue, i, j } = useMemo(() => {
		if (
			!wallMatrix ||
			!settingsData?.stationHeight ||
			!settingsData?.stationPos ||
			settingsData.stationPos.length < 2
		) {
			return { matrixIndexValue: undefined, i: undefined, j: undefined };
		}
		const { i, j } = geoToMatrixIndex(
			settingsData.stationPos[0] as unknown as number,
			settingsData.stationPos[1] as unknown as number,
			data.mapData.coordinates[0][0][0],
			data.mapData.coordinates[0][2][0],
			data.mapData.coordinates[0][0][1],
			data.mapData.coordinates[0][2][1],
			data.mapData.size
		);
		const matrixIndexValue = getMatrixValue(
			wallMatrix,
			i,
			j,
			Number(settingsData.stationHeight),
			data.mapData.size,
			data.mapData.size,
			30
		);
		return { matrixIndexValue, i, j };
	}, [wallMatrix, settingsData?.stationHeight, settingsData?.stationPos, data?.mapData?.coordinates]);
	const spherePositions = useMemo(() => {
		if (!data?.mapData?.coordinates || !rayLaunchData || !Array.isArray(rayLaunchData)) return [];
		const coordinates = data.mapData.coordinates;

		return rayLaunchData.map((rayPath, rayIndex) => {
			const filteredRayData = rayPath.filter((_, index) => index % 3 === 0);
			return {
				rayIndex,
				positions: filteredRayData.map(({ x, y, z, power }) => {
					const { lon, lat } = matrixIndexToGeo(
						x,
						y,
						coordinates[0][0][0],
						coordinates[0][2][0],
						coordinates[0][0][1],
						coordinates[0][2][1],
						data.mapData.size
					);
					const coord = mapboxgl.MercatorCoordinate.fromLngLat([lon, lat], z ?? 0);
					return {
						coord,
						power,
					};
				}),
			};
		});
	}, [rayLaunchData, data?.mapData?.coordinates]);
	return (
		<>
			{settingsData.isOpen && (
				<Modal onClose={handleOnSettingsClose}>
					<div className={styles.settingsBtnsBox}>
						<button
							onClick={() => setSettingsData(prev => ({ ...prev, settingsType: "global" }))}
							className={`${styles.settingsTypeBtn} ${
								settingsData.settingsType === "global" ? styles.settingsTypeBtnActive : ""
							}`}>
							Global
						</button>
						<button
							onClick={() => setSettingsData(prev => ({ ...prev, settingsType: "singleRay" }))}
							className={`${styles.settingsTypeBtn} ${
								settingsData.settingsType === "singleRay" ? styles.settingsTypeBtnActive : ""
							}`}>
							Single Ray
						</button>
					</div>
					<motion.div className={styles.dialogBox}>
						<AnimatePresence mode="wait">
							{settingsData.settingsType === "global" && (
								<GlobalSettings handleFormSubmit={handleGlobalSettingsSubmit} formData={settingsData} />
							)}
							{settingsData.settingsType === "singleRay" && (
								<SingleRaySettings
									handleFormSubmit={handleSingleRaySettingsSubmit}
									formData={settingsData}
									handleAddRay={handleAddRay}
									handleRemoveRay={handleRemoveRay}
								/>
							)}
						</AnimatePresence>
					</motion.div>
				</Modal>
			)}
			{isPendingRayLaunch && (
				<div className={styles.loadingScreen}>
					<SpinWifi />
				</div>
			)}
			{!isLoading && data && wallMatrix && (
				<div className={styles.box}>
					<div className={styles.titleBox}>
						<h3>{data.mapData.title}</h3>
					</div>
					<div className={styles.mapBox}>
						{settingsData.stationPos && matrixIndexValue && (
							<Map
								{...data.mapData}
								{...settingsData}
								handleStationPosUpdate={handleStationPosUpdate}
								buildingsData={data.buildingsData}
								spherePositions={spherePositions}
								wallMatrix={wallMatrix}
								powerMap={powerMapData}
							/>
						)}
						<div className={styles.stationPosContener}>
							{settingsData.stationPos && (
								<>
									<p>Station position</p>
									<p>
										Longitude: {parseFloat(settingsData.stationPos.toString().split(",")[0]).toFixed(6)} |{" "}
										{
											geoToMatrixIndex(
												settingsData.stationPos[0] as unknown as number,
												settingsData.stationPos[1] as unknown as number,
												data.mapData.coordinates[0][0][0],
												data.mapData.coordinates[0][2][0],
												data.mapData.coordinates[0][0][1],
												data.mapData.coordinates[0][2][1],
												data.mapData.size
											).i
										}{" "}
									</p>
									<p>
										Latitude: {parseFloat(settingsData.stationPos.toString().split(",")[1]).toFixed(6)} |{" "}
										{
											geoToMatrixIndex(
												settingsData.stationPos[0] as unknown as number,
												settingsData.stationPos[1] as unknown as number,
												data.mapData.coordinates[0][0][0],
												data.mapData.coordinates[0][2][0],
												data.mapData.coordinates[0][0][1],
												data.mapData.coordinates[0][2][1],
												data.mapData.size
											).j
										}
									</p>
									{matrixIndexValue && <p>index: {matrixIndexValue}</p>}
								</>
							)}
						</div>
						<p className={styles.brandName}>Rayek</p>
					</div>
					<motion.div
						className={styles.sliderBox}
						initial={{ opacity: 0, y: -20 }}
						animate={{ opacity: 1, y: 0 }}
						transition={{ duration: 0.5 }}>
						<label>Height: {settingsData.powerMapHeight}m</label>
						<input
							className={styles.slider}
							type="range"
							value={settingsData.powerMapHeight}
							onChange={e => setSettingsData(prev => ({ ...prev, powerMapHeight: Number(e.target.value) }))}
							min={0}
							max={29}
						/>
					</motion.div>
					<button
						className={styles.settingsBtn}
						onClick={() =>
							setSettingsData(prevSettingsData => {
								const updatedsettingsData = { ...prevSettingsData, isOpen: !prevSettingsData.isOpen };
								return updatedsettingsData;
							})
						}>
						<IoMdSettings />
					</button>
					<button onClick={handleComputeBtn} className={styles.computeBtn}>
						raylaunch
					</button>
				</div>
			)}
		</>
	);
}
