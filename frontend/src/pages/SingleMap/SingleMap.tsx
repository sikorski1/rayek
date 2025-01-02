import SettingsDialog from "@/components/SettingsDialog/SettingsDialog";
import Title from "@/components/Title/Title";
import Map from "@/pages/SingleMap/Map/Map";
import { MapTypes, postComputeTypes } from "@/types/main";
import { url } from "@/utils/url";
import axios from "axios";
import { useEffect, useState } from "react";
import { IoMdClose, IoMdSettings } from "react-icons/io";
import { useParams } from "react-router-dom";
import styles from "./singleMap.module.scss";

const getMapData = async ({ mapTitle }: { mapTitle: string }) => {
	try {
		const response = await axios.get(url + `/raycheck/:${mapTitle}`);
		return response.data;
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

const data: MapTypes[] = [
	{
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
	},
];

export default function SingleMap() {
	const [popSettings, setPopSettings] = useState<boolean>(false);
	const [frequency, setFrequency] = useState<string>("1000");
	const [stationHeight, setStationHeight] = useState<string>("0");
	const [mapData, setMapData] = useState<MapTypes | null>(null);
	const [isLoading, setIsLoading] = useState<boolean>(true);
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
			let response;
			try {
				response = await getMapData({ mapTitle: id! });
				setMapData(response);
				setIsLoading(false);
			} catch (error) {
				console.log(error);
			}
			return response;
		};
		fetchData();
	}, [id]);
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
			{!isLoading && (
				<div className={styles.box}>
					<div className={styles.titleBox}>
						<Title>{id}</Title>
					</div>
					<div className={styles.mapBox}>
						<Map {...mapData!} />
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
