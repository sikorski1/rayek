import FormField from "@/components/FormField/FormField";
import SpinWifi from "@/components/Loaders/SpinWifi";
import Modal from "@/components/Modal/Modal";
import { useGetMapById, useRayLaunching } from "@/hooks/useMap";
import Map from "@/pages/SingleMap/Map/Map";
import { PopupDataTypes, RayLaunchType } from "@/types/main";
import { geoToMatrixIndex } from "@/utils/geoToMatrixIndex";
import { getMatrixValue } from "@/utils/getMatrixValue";
import { loadWallMatrix } from "@/utils/loadWallMatrix";
import { matrixIndexToGeo } from "@/utils/matrixIndexToGeo";
import { useEffect, useMemo, useState } from "react";
import { IoMdSettings } from "react-icons/io";
import mapboxgl from "mapbox-gl";
import { useParams } from "react-router-dom";
import styles from "./singleMap.module.scss";
const initialPopupData: PopupDataTypes = {
	isOpen: false,
	frequency: 5,
	stationHeight: 5,
	numberOfRaysAzimuth: 360,
	numberOfRaysElevation: 360,
	numberOfInteractions: 5,
	reflectionFactor: 0.5,
	stationPower: 5,
	minimalRayPower: -120,
};

export default function SingleMap() {
	const [popupData, setPopupData] = useState<PopupDataTypes>(initialPopupData);
	const [wallMatrix, setWallMatrix] = useState<Float64Array | null>(null);
	const [rayLaunchData, setRayLaunchData] = useState<RayLaunchType | null>(null);
	const { id } = useParams();
	const { data, isLoading, error } = useGetMapById(id!);

	const handleStationPosUpdate = (stationPos: mapboxgl.LngLatLike) => {
		setPopupData(prev => {
			const updatedSingleMapData = { ...prev, stationPos: stationPos };
			return updatedSingleMapData;
		});
	};

	const handleOnSettingsClose = () => {
		setPopupData(prevPopupData => {
			const updatedPopupData = { ...prevPopupData, isOpen: false };
			return updatedPopupData;
		});
	};

	const handleDialogFormSubmit = (event: React.FormEvent<HTMLFormElement>) => {
		event.preventDefault();

		const form = event.currentTarget;
		const formData = new FormData(form);

		const updatedPopupData: PopupDataTypes = {
			isOpen: false,
			numberOfRaysAzimuth: Number(formData.get("raysAzimuth")),
			numberOfRaysElevation: Number(formData.get("raysElevation")),
			frequency: Number(formData.get("frequency")),
			stationHeight: Number(formData.get("stationHeight")),
			reflectionFactor: Number(formData.get("relfectionFactor")),
			numberOfInteractions: Number(formData.get("interactions")),
			stationPower: Number(formData.get("stationPower")),
			minimalRayPower: Number(formData.get("minimalRayPower")),
		};
		setPopupData(prev => ({ ...prev, ...updatedPopupData }));
	};
	const handleOnSuccess = (data: any) => {
		setRayLaunchData(data.rayPath);
	};
	const { mutate, isPending: isPendingRayLaunch } = useRayLaunching(handleOnSuccess);
	const handleComputeBtn = async () => {
		const { stationPos, stationHeight, ...restData } = popupData;
		mutate({
			mapTitle: id!,
			configData: { stationPos: { x: i, y: j, z: Number(stationHeight) }, ...restData },
		});
	};
	useEffect(() => {
		if (!data) return;
		setPopupData(prev => ({ ...prev, stationPos: data.mapData.center }));
	}, [data]);

	useEffect(() => {
		if (!id) return;
		loadWallMatrix(id).then(setWallMatrix);
	}, [id]);

	const { matrixIndexValue, i, j } = useMemo(() => {
		if (!wallMatrix || !popupData?.stationHeight || !popupData?.stationPos || popupData.stationPos.length < 2) {
			return { matrixIndexValue: undefined, i: undefined, j: undefined };
		}
		const { i, j } = geoToMatrixIndex(
			popupData.stationPos[0] as unknown as number,
			popupData.stationPos[1] as unknown as number,
			data.mapData.coordinates[0][0][0],
			data.mapData.coordinates[0][2][0],
			data.mapData.coordinates[0][0][1],
			data.mapData.coordinates[0][2][1],
			250
		);
		const matrixIndexValue = getMatrixValue(wallMatrix, i, j, Number(popupData.stationHeight));
		return { matrixIndexValue, i, j };
	}, [wallMatrix, popupData?.stationHeight, popupData?.stationPos, data?.mapData?.coordinates]);
	const spherePositions = useMemo(() => {
		if (!data?.mapData?.coordinates || !rayLaunchData) return;
		const coordinates = data.mapData.coordinates
		return rayLaunchData.map(({ x, y, z }) => {
			const { lon, lat } = matrixIndexToGeo(
				x,
				y,
				coordinates[0][0][0],
				coordinates[0][2][0],
				coordinates[0][0][1],
				coordinates[0][2][1],
				250
			);
			console.log(lon,lat);
			return mapboxgl.MercatorCoordinate.fromLngLat([lon, lat], z);
		});
	}, [rayLaunchData, data?.mapData?.coordinates]);
	return (
		<>
			{popupData.isOpen && (
				<Modal onClose={handleOnSettingsClose}>
					<div className={styles.dialogBox}>
						<form onSubmit={handleDialogFormSubmit} className={styles.formBox}>
							<div className={styles.formInputBox}>
								<FormField
									label="RAYS AZIMUTH"
									name="raysAzimuth"
									defaultValue={popupData.numberOfRaysAzimuth}
									placeholder="Enter number of rays azimuth"
									min={1}
									max={1440}
									step={1}
									required
								/>
								<FormField
									label="STATION POWER (watt)"
									name="stationPower"
									defaultValue={popupData.stationPower}
									placeholder="Enter station power in watt"
									min={0.1}
									max={100}
									step={0.1}
									required
								/>
								<FormField
									label="FREQUENCY (GHz)"
									name="frequency"
									defaultValue={popupData.frequency}
									placeholder="Enter frequency in MHz"
									min={0.1}
									max={100}
									step={0.1}
									required
								/>
								<FormField
									label="REFLECTION FACTOR"
									name="relfectionFactor"
									defaultValue={popupData.reflectionFactor}
									placeholder="Enter reflection factor"
									min={0}
									max={1}
									step={0.01}
									required
								/>
							</div>

							<div className={styles.formInputBox}>
								<FormField
									label="RAYS ELEVATION"
									name="raysElevation"
									defaultValue={popupData.numberOfRaysElevation}
									placeholder="Enter number of rays elevation"
									min={1}
									max={1440}
									step={1}
									required
								/>
								<FormField
									label="STATION HEIGHT (m)"
									name="stationHeight"
									defaultValue={popupData.stationHeight}
									placeholder="Enter station height in meters"
									min={0}
									max={29}
									step={1}
									required
								/>
								<FormField
									label="INTERACTIONS"
									name="interactions"
									defaultValue={popupData.numberOfInteractions}
									placeholder="Enter number of interactions"
									min={1}
									max={10}
									step={1}
									required
								/>
								<FormField
									label="MINIMAL RAY POWER (dBm)"
									name="minimalRayPower"
									defaultValue={popupData.minimalRayPower}
									placeholder="Enter minimal ray power in dBm"
									min={-160}
									max={-60}
									step={0.01}
									required
								/>
							</div>
							<button className={styles.submitBtn} type="submit">
								Enter
							</button>
						</form>
					</div>
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
						{popupData.stationPos && matrixIndexValue && (
							<Map
								{...data.mapData}
								stationPos={popupData.stationPos}
								stationHeight={popupData.stationHeight}
								handleStationPosUpdate={handleStationPosUpdate}
								buildingsData={data.buildingsData}
								spherePositions={spherePositions}
								wallMatrix={wallMatrix}
							/>
						)}
						<div className={styles.stationPosContener}>
							{popupData.stationPos && (
								<>
									<p>Station position</p>
									<p>
										Longitude: {parseFloat(popupData.stationPos.toString().split(",")[0]).toFixed(6)} |{" "}
										{
											geoToMatrixIndex(
												popupData.stationPos[0] as unknown as number,
												popupData.stationPos[1] as unknown as number,
												data.mapData.coordinates[0][0][0],
												data.mapData.coordinates[0][2][0],
												data.mapData.coordinates[0][0][1],
												data.mapData.coordinates[0][2][1],
												250
											).i
										}{" "}
									</p>
									<p>
										Latitude: {parseFloat(popupData.stationPos.toString().split(",")[1]).toFixed(6)} |{" "}
										{
											geoToMatrixIndex(
												popupData.stationPos[0] as unknown as number,
												popupData.stationPos[1] as unknown as number,
												data.mapData.coordinates[0][0][0],
												data.mapData.coordinates[0][2][0],
												data.mapData.coordinates[0][0][1],
												data.mapData.coordinates[0][2][1],
												250
											).j
										}
									</p>
									{matrixIndexValue && <p>index: {matrixIndexValue}</p>}
								</>
							)}
						</div>
						<p className={styles.brandName}>Rayek</p>
					</div>
					<button
						className={styles.settingsBtn}
						onClick={() =>
							setPopupData(prevPopupData => {
								const updatedPopupData = { ...prevPopupData, isOpen: !prevPopupData.isOpen };
								return updatedPopupData;
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
