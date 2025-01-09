import SettingsDialog from "@/components/SettingsDialog/SettingsDialog";
import Map from "@/pages/SingleMap/Map/Map";
import { MapTypes, postComputeTypes } from "@/types/main";
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
		return response.data.mapConf;
	} catch (error) {
		console.log(error);
	}
};

const getBuildingsData = async ({
	mapTitle
}: {
	mapTitle:string
}): Promise<FeatureCollection | undefined> => {
	try {
		const response = await axios.get(url + `/raycheck/buildings/${mapTitle}`);
		return response.data.buildingsData;
	} catch (error) {
		console.log(error);
	}
};

const postCompute = async ({ freq, stationH }: postComputeTypes) => {
	let response;
	const data = {
		freq: freq,
		stationH: stationH,
	};
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

export default function SingleMap() {
	const [popSettings, setPopSettings] = useState<boolean>(false);
	const [frequency, setFrequency] = useState<string>("1000");
	const [stationHeight, setStationHeight] = useState<string>("0");
	const [stationPos, setStationPos] = useState<mapboxgl.LngLatLike | null>(null);
	const [mapData, setMapData] = useState<MapTypes | null>(null);
	const [buildingsData, setBuildingsData] = useState<FeatureCollection | null>(null);

	const { id } = useParams();
	const handleOnSettingsClose = () => {
		setPopSettings(false);
	};

	const handleDialogFormSubmit = (event: React.FormEvent<HTMLFormElement>) => {
		event.preventDefault();
		const target = event.target as typeof event.target & {
			frequency: { value: string };
			stationHeight: { value: string };
		};
		setFrequency(target.frequency.value);
		setStationHeight(target.stationHeight.value);
		setPopSettings(false);
	};

	const handleComputeBtn = async () => {
		const data = { freq: frequency, stationH: stationHeight };
		const response = await postCompute(data);
		console.log(response);
	};

	useEffect(() => {
		const fetchData = async () => {
			try {
				const mapResponse = (await getMapData({ mapTitle: id! })) || {
					title: "AGH fragment",
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
					setMapData(mapResponse);
					setStationPos(mapResponse.center);
				}
			} catch (error) {
				console.error("Error fetching map data:", error);
			}
		};
		fetchData();
	}, [id]);

	useEffect(() => {
		if (mapData) {
			const fetchBuildings = async () => {
				try {
					const buildingsResponse = await getBuildingsData({ mapTitle: id! });
					if (buildingsResponse) {
						setBuildingsData(buildingsResponse);
					}
				} catch (error) {
					console.error("Error fetching buildings data:", error);
				}
			};
			fetchBuildings();
		}
	}, [mapData]);

	return (
		<>
			{popSettings && (
				<SettingsDialog popState={popSettings} handleOnClose={handleOnSettingsClose}>
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
										defaultValue={frequency}
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
										defaultValue={stationHeight}
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
			{mapData && buildingsData && (
				<div className={styles.box}>
					<div className={styles.titleBox}>
						<h3>{id}</h3>
					</div>
					<div className={styles.mapBox}>
						<Map {...mapData!} stationPos={stationPos!} setStationPos={setStationPos} buildingsData={buildingsData} />
						<div className={styles.stationPosContener}>
							{stationPos && (
								<>
									<p>Station position</p>
									<p>Longitude: {parseFloat(stationPos.toString().split(",")[0]).toFixed(6)}</p>
									<p>Latitude: {parseFloat(stationPos.toString().split(",")[1]).toFixed(6)}</p>
								</>
							)}
						</div>
						<p className={styles.brandName}>ReyCheck</p>
					</div>
					<button className={styles.settingsBtn} onClick={() => setPopSettings(!popSettings)}>
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
