import Modal from "@/components/Modal/Modal";
import { useGetMapById, useRayLaunching } from "@/hooks/useMap";
import Map from "@/pages/SingleMap/Map/Map";
import { PopupDataTypes, PostComputeTypes, SingleMapDataTypes } from "@/types/main";
import { geoToMatrixIndex } from "@/utils/geoToMatrixIndex";
import { getMatrixValue } from "@/utils/getMatrixValue";
import { loadWallMatrix } from "@/utils/loadWallMatrix";
import { url } from "@/utils/url";
import axios from "axios";
import { useEffect, useMemo, useState } from "react";
import { IoMdSettings } from "react-icons/io";
import { useParams } from "react-router-dom";
import styles from "./singleMap.module.scss";

const postCompute = async (data: PostComputeTypes, mapTitle: string) => {
	let response;
	try {
		response = await axios.post(
			url + `/raycheck/rayLaunch/${mapTitle}`,
			{},
			{
				headers: {
					"Content-Type": "application/json",
				},
			}
		);
	} catch (error) {
		console.log(error);
	}
	return response;
};

const initialPopupData: PopupDataTypes = {
	isOpen: false,
	frequency: "1000",
	stationHeight: "0",
};

export default function SingleMap() {
	const [popupData, setPopupData] = useState<PopupDataTypes>(initialPopupData);
	const [singleMapData, setSingleMapData] = useState<SingleMapDataTypes>({} as SingleMapDataTypes);
	const [wallMatrix, setWallMatrix] = useState<Float64Array | null>(null);
	const { id } = useParams();
	const { data, isLoading, error } = useGetMapById(id!);
	const { mutate } = useRayLaunching();
	const handleStationPosUpdate = (stationPos: mapboxgl.LngLatLike) => {
		setSingleMapData(prevSingleMapData => {
			const updatedSingleMapData = { ...prevSingleMapData, stationPos: stationPos };
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
		const target = event.target as typeof event.target & {
			frequency: { value: string };
			stationHeight: { value: string };
		};
		setPopupData({ isOpen: false, frequency: target.frequency.value, stationHeight: target.stationHeight.value });
	};

	const handleComputeBtn = async () => {
		const { frequency, stationHeight } = popupData;
		mutate({
			mapTitle: id!,
			configData: { stationPos: { x: i, y: j, z: Number(stationHeight) } },
		});
	};
	useEffect(() => {
		if (!data) return;
		setSingleMapData(prev => ({ ...prev, stationPos: data.mapData.center }));
	}, [data]);

	useEffect(() => {
		if (!id) return;
		loadWallMatrix(id).then(setWallMatrix);
	}, [id]);

	const { matrixIndexValue, i, j } = useMemo(() => {
		if (!wallMatrix || !popupData?.stationHeight || !singleMapData?.stationPos || singleMapData.stationPos.length < 2) {
			return { matrixIndexValue: undefined, i: undefined, j: undefined };
		}
		const { i, j } = geoToMatrixIndex(
			singleMapData.stationPos[0] as unknown as number,
			singleMapData.stationPos[1] as unknown as number,
			data.mapData.coordinates[0][0][0],
			data.mapData.coordinates[0][2][0],
			data.mapData.coordinates[0][0][1],
			data.mapData.coordinates[0][2][1],
			250
		);
		const matrixIndexValue = getMatrixValue(wallMatrix, i, j, Number(popupData.stationHeight));
		return { matrixIndexValue, i, j };
	}, [wallMatrix, popupData?.stationHeight, singleMapData?.stationPos, data?.mapData?.coordinates]);
	return (
		<>
			{popupData.isOpen && (
				<Modal onClose={handleOnSettingsClose}>
					<div className={styles.dialogBox}>
						<form onSubmit={handleDialogFormSubmit} className={styles.formBox}>
							<div className={styles.formInputBox}>
								<div>
									<label htmlFor="frequency" className={styles.label}>
										FREQUENCY (MHz)
									</label>
									<input
										type="number"
										name="frequency"
										defaultValue={popupData.frequency}
										className={styles.input}
										placeholder="Enter frequency in MHz"
										min="100"
										max="100000"
										step="1"
										required
									/>
								</div>
								<div>
									<label htmlFor="stationHeight" className={styles.label}>
										STATION HEIGHT (m)
									</label>
									<input
										type="number"
										name="stationHeight"
										defaultValue={popupData.stationHeight}
										className={styles.input}
										placeholder="Enter station height in meters"
										min="0"
										max="10"
										step="0.1"
										required
									/>
								</div>
							</div>
							<div className={styles.formInputBox}>
								<div>
									<label htmlFor="frequency" className={styles.label}>
										FREQUENCY (MHz)
									</label>
									<input
										type="number"
										name="frequency"
										defaultValue={popupData.frequency}
										className={styles.input}
										placeholder="Enter frequency in MHz"
										min="100"
										max="100000"
										step="1"
										required
									/>
								</div>
								<div>
									<label htmlFor="stationHeight" className={styles.label}>
										STATION HEIGHT (m)
									</label>
									<input
										type="number"
										name="stationHeight"
										defaultValue={popupData.stationHeight}
										className={styles.input}
										placeholder="Enter station height in meters"
										min="0"
										max="10"
										step="0.1"
										required
									/>
								</div>
							</div>
							<div className={styles.formInputBox}>
								<div>
									<label htmlFor="frequency" className={styles.label}>
										FREQUENCY (MHz)
									</label>
									<input
										type="number"
										name="frequency"
										defaultValue={popupData.frequency}
										className={styles.input}
										placeholder="Enter frequency in MHz"
										min="100"
										max="100000"
										step="1"
										required
									/>
								</div>
								<div>
									<label htmlFor="stationHeight" className={styles.label}>
										STATION HEIGHT (m)
									</label>
									<input
										type="number"
										name="stationHeight"
										defaultValue={popupData.stationHeight}
										className={styles.input}
										placeholder="Enter station height in meters"
										min="0"
										max="10"
										step="0.1"
										required
									/>
								</div>
							</div>
							<div className={styles.formInputBox}></div>
							<button className={styles.submitBtn} type="submit">
								Enter
							</button>
						</form>
					</div>
				</Modal>
			)}
			{!isLoading && data && wallMatrix && (
				<div className={styles.box}>
					<div className={styles.titleBox}>
						<h3>{data.mapData.title}</h3>
					</div>
					<div className={styles.mapBox}>
						{singleMapData.stationPos && matrixIndexValue && (
							<Map
								{...data.mapData}
								stationPos={singleMapData.stationPos}
								stationHeight={popupData.stationHeight}
								handleStationPosUpdate={handleStationPosUpdate}
								buildingsData={data.buildingsData}
								computationResult={data.computationResult}
								wallMatrix={wallMatrix}
							/>
						)}
						<div className={styles.stationPosContener}>
							{singleMapData.stationPos && (
								<>
									<p>Station position</p>
									<p>
										Longitude: {parseFloat(singleMapData.stationPos.toString().split(",")[0]).toFixed(6)} |{" "}
										{
											geoToMatrixIndex(
												singleMapData.stationPos[0] as unknown as number,
												singleMapData.stationPos[1] as unknown as number,
												data.mapData.coordinates[0][0][0],
												data.mapData.coordinates[0][2][0],
												data.mapData.coordinates[0][0][1],
												data.mapData.coordinates[0][2][1],
												250
											).i
										}{" "}
									</p>
									<p>
										Latitude: {parseFloat(singleMapData.stationPos.toString().split(",")[1]).toFixed(6)} |{" "}
										{
											geoToMatrixIndex(
												singleMapData.stationPos[0] as unknown as number,
												singleMapData.stationPos[1] as unknown as number,
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
						<p className={styles.brandName}>ReyCheck</p>
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
						Compute
					</button>
				</div>
			)}
		</>
	);
}
