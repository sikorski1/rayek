import SettingsDialog from "@/components/SettingsDialog/SettingsDialog";
import Map from "@/pages/SingleMap/Map/Map";
import { PopupDataTypes, PostComputeTypes, SingleMapDataTypes } from "@/types/main";
import { url } from "@/utils/url";
import axios from "axios";
import { FeatureCollection } from "geojson";
import { useEffect, useState } from "react";
import { IoMdClose, IoMdSettings } from "react-icons/io";
import { useParams } from "react-router-dom";
import styles from "./singleMap.module.scss";
const getMapData = async ({ mapTitle }: { mapTitle: string }) => {
	try {
		const response = await axios.get(url + `/raycheck/${mapTitle}`);
		return response.data;
	} catch (error) {
		console.log(error);
	}
};

const getBuildingsData = async ({ mapTitle }: { mapTitle: string }): Promise<FeatureCollection | undefined> => {
	try {
		const response = await axios.get(url + `/raycheck/buildings/${mapTitle}`);
		return response.data;
	} catch (error) {
		console.log(error);
	}
};

const postCompute = async (data:PostComputeTypes) => {
	let response;
	try {
		response = await axios.post(url + "/raycheck/compute", data, {
			headers: {
				"Content-Type": "application/json",
			},
		});
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

	const { id } = useParams();

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
		const data: PostComputeTypes = { freq: frequency, stationH: stationHeight };
		const response = await postCompute(data);
		if (response) {
			setSingleMapData((prevSingleMapData) => {
				const updatedSingleMapData = {...prevSingleMapData, computationResult:response.data}
				return updatedSingleMapData
			})
		}
	};

	useEffect(() => {
		const fetchData = async () => {
			try {
				const mapResponse = (await getMapData({ mapTitle: id! })) || {
					title: "AGHFragment",
					coordinates: [
						[
							[19.914029, 50.065311], // Southwest
							[19.917527, 50.065311], // Southeast
							[19.917527, 50.067556], // Northeast
							[19.914029, 50.067556], // Northwest
							[19.914029, 50.065311], // Zamknięcie pętli
						],
					],
					center: [19.915778, 50.0664335],
					bounds: [
						[19.914029, 50.065311], // Southwest corner (dolny lewy róg)
						[19.917527, 50.067556], // Northeast corner (górny prawy róg)
					],
				};
				if (mapResponse) {
					setSingleMapData(prevSingleMapData => {
						const updatedSingleMapData = { ...prevSingleMapData, stationPos: mapResponse.center, mapData: mapResponse };
						return updatedSingleMapData;
					});
				}
			} catch (error) {
				console.error("Error fetching map data:", error);
			}
		};
		fetchData();
	}, [id]);

	useEffect(() => {
		if (singleMapData.mapData) {
			const fetchBuildings = async () => {
				try {
					const buildingsResponse = await getBuildingsData({ mapTitle: id! });
					if (buildingsResponse) {
						setSingleMapData(prevSingleMapData => {
							const updatedSingleMapData = { ...prevSingleMapData, buildingsData: buildingsResponse };
							return updatedSingleMapData;
						});
					}
				} catch (error) {
					console.error("Error fetching buildings data:", error);
				}
			};
			fetchBuildings();
		}
	}, [singleMapData.mapData])
	return (
		<>
			{popupData.isOpen && (
				<SettingsDialog popState={popupData.isOpen} handleOnClose={handleOnSettingsClose}>
					<div className={styles.dialogBox}>
						<div className={styles.closeBtnBox}>
							<button className={styles.closeBtn} onClick={handleOnSettingsClose}>
								<IoMdClose />
							</button>
						</div>
						<form onSubmit={handleDialogFormSubmit} className={styles.form}>
							<div className={styles.formBox}>
								<div className={styles.formInputBox}>
									<label htmlFor="frequency" className={styles.label}>
										Frequency (MHz):
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
								<div className={styles.formInputBox}>
									<label htmlFor="stationHeight" className={styles.label}>
										Station Height (m):
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
								<button className={styles.submitBtn} type="submit">
									Enter
								</button>
							</div>
						</form>
					</div>
				</SettingsDialog>
			)}
			{singleMapData.mapData && singleMapData.buildingsData && (
				<div className={styles.box}>
					<div className={styles.titleBox}>
						<h3>{singleMapData.mapData.title}</h3>
					</div>
					<div className={styles.mapBox}>
						<Map
							{...singleMapData.mapData!}
							stationPos={singleMapData.stationPos!}
							handleStationPosUpdate={handleStationPosUpdate}
							buildingsData={singleMapData.buildingsData}
							computationResult={singleMapData.computationResult}
						/>
						<div className={styles.stationPosContener}>
							{singleMapData.stationPos && (
								<>
									<p>Station position</p>
									<p>Longitude: {parseFloat(singleMapData.stationPos.toString().split(",")[0]).toFixed(6)}</p>
									<p>Latitude: {parseFloat(singleMapData.stationPos.toString().split(",")[1]).toFixed(6)}</p>
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
