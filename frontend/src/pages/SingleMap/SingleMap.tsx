import SpinWifi from "@/components/Loaders/SpinWifi";
import LogoutButton from "@/components/LogoutButton/LogoutButton";
import GlobalSettings from "@/components/Modal/GlobalSettings/GlobalSettings";
import Modal from "@/components/Modal/Modal";
import SingleRaySettings from "@/components/Modal/SingleRaySettings/SingleRaySettings";
import { useGetMapById, useRayLaunching } from "@/hooks/useMap";
import Map from "@/pages/SingleMap/Map/Map";
import { PowerMapLegendEntry, PowerMapLegendType, RayLaunchType, SettingsDataTypes } from "@/types/main";
import { geoToMatrixIndex } from "@/utils/geoToMatrixIndex";
import { getHeatMapColor } from "@/utils/getHeatMapColor";
import { getMatrixValue } from "@/utils/getMatrixValue";
import { getNormalizedValueFromLabel } from "@/utils/getNormalizedValueFromLabel";
import { loadWallMatrix } from "@/utils/loadWallMatrix";
import { matrixIndexToGeo } from "@/utils/matrixIndexToGeo";
import { AnimatePresence, motion } from "framer-motion";
import mapboxgl from "mapbox-gl";
import { useEffect, useMemo, useState } from "react";
import { IoMdSettings } from "react-icons/io";
import { Link, useParams } from "react-router-dom";
import styles from "./singleMap.module.scss";
const initialSettingsData: SettingsDataTypes = {
	settingsType: "global",
	isOpen: false,
	frequency: 3.8,
	stationHeight: 5,
	numberOfRaysAzimuth: 360,
	numberOfRaysElevation: 360,
	numberOfInteractions: 5,
	reflectionFactor: 1,
	stationPower: 5,
	minimalRayPower: -160,
	singleRays: [],
	powerMapHeight: 5,
	isPowerMapVisible: true,
	diffractionRayNumber: 20,
};

const createDefaultLegendEntry = (): PowerMapLegendEntry => ({
	total: 0,
	"< 0dbm": 0,
	"< -20dbm": 0,
	"< -40dbm": 0,
	"< -60dbm": 0,
	"< -80dbm": 0,
	"< -100dbm": 0,
	"< -120dbm": 0,
	"< -140dbm": 0,
});

const defaultPowerMapLegend: PowerMapLegendType = Object.fromEntries(
	Array.from({ length: 10 }, (_, z) => [z, createDefaultLegendEntry()])
);

export default function SingleMap() {
	const [settingsData, setSettingsData] = useState<SettingsDataTypes>(initialSettingsData);
	const [wallMatrix, setWallMatrix] = useState<Int16Array | null>(null);
	const [rayLaunchData, setRayLaunchData] = useState<RayLaunchType | null>(null);
	const [powerMapData, setPowerMapData] = useState<any>(null);
	const [powerMapLegend, setPowerMapLegend] = useState<PowerMapLegendType>(defaultPowerMapLegend);
	const { id } = useParams();
	const { data, isLoading } = useGetMapById(id!);

	const handleStationPosUpdate = (stationPos: [number, number]) => {
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

		const updatedSettingsData: Omit<
			SettingsDataTypes,
			"settingsType" | "singleRays" | "isOpen" | "isPowerMapVisible" | "powerMapHeight"
		> = {
			numberOfRaysAzimuth: Number(formData.get("raysAzimuth")),
			numberOfRaysElevation: Number(formData.get("raysElevation")),
			frequency: Number(formData.get("frequency")),
			stationHeight: Number(formData.get("stationHeight")),
			reflectionFactor: Number(formData.get("relfectionFactor")),
			numberOfInteractions: Number(formData.get("interactions")),
			stationPower: Number(formData.get("stationPower")),
			minimalRayPower: Number(formData.get("minimalRayPower")),
			diffractionRayNumber: Number(formData.get("diffractionRayNumber")),
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
		setPowerMapLegend(data.powerMapLegend);
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
			if (rayPath.length === 0) return [];
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
						<div className={styles.legendBox}>
							<div className={styles.legendHeaderBox}>
								<p className={styles.legendHeaderText}>Total coverage</p>
								<p className={styles.legendHeaderPercentage}>
									{powerMapLegend[settingsData.powerMapHeight].total.toFixed(2)}%
								</p>
							</div>

							<div className={styles.powerRangesBox}>
								{Object.entries(powerMapLegend[settingsData.powerMapHeight])
									.filter(([power]) => power !== "total")
									.map(([power, percentage]) => {
										const normalized = getNormalizedValueFromLabel(power);
										const color = getHeatMapColor(normalized);
										const [main, unit] = power.split(/(dbm)/gi);

										return (
											<div className={styles.powerRangeEntry} key={power}>
												<div
													className={styles.powerCircle}
													style={{
														backgroundColor: `rgb(${color.r * 255}, ${color.g * 255}, ${
															color.b * 255
														})`,
													}}
												/>
												<p className={styles.powerRangeText}>
													{main}
													{unit && (
														<span style={{ fontSize: "0.70em", marginLeft: 1 }}>
															[{unit}]
														</span>
													)}
												</p>
												<p className={styles.percentageText}>{percentage.toFixed(2)}%</p>
											</div>
										);
									})}
							</div>
						</div>
						<div className={styles.stationPosContener}>
							{settingsData.stationPos && (
								<>
									<p>Station position</p>
									<p>
										Longitude:{" "}
										{parseFloat(settingsData.stationPos.toString().split(",")[0]).toFixed(6)} |{" "}
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
										Latitude:{" "}
										{parseFloat(settingsData.stationPos.toString().split(",")[1]).toFixed(6)} |{" "}
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
								</>
							)}
						</div>
						<p className={styles.brandName}>Rayek</p>
					</div>
					{powerMapData && (
						<motion.div
							className={styles.sliderBox}
							initial={{ opacity: 0, y: -20 }}
							animate={{ opacity: 1, y: 0 }}
							transition={{ duration: 0.5 }}>
							<label>Height {settingsData.powerMapHeight}m</label>
							<input
								className={styles.slider}
								type="range"
								value={settingsData.powerMapHeight}
								onChange={e =>
									setSettingsData(prev => ({ ...prev, powerMapHeight: Number(e.target.value) }))
								}
								min={0}
								max={29}
							/>
						</motion.div>
					)}
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
					<div className={styles.infoBox}>
						<LogoutButton />
						<Link to="/info" className={styles.infoLink}>
							i
						</Link>
					</div>
				</div>
			)}
		</>
	);
}
