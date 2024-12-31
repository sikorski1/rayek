import sampleImg from "@/assets/imgs/sampleMapImg.jpg";
import SettingsDialog from "@/components/SettingsDialog/SettingsDialog";
import Title from "@/components/Title/Title";
import Wrapper from "@/components/Wrapper/Wrapper";
import { url } from "@/utils/url";
import axios from "axios";
import { useState } from "react";
import { IoMdClose, IoMdSettings } from "react-icons/io";
import { useParams } from "react-router-dom";
import styles from "./singleMap.module.scss";
import Map from "@/pages/SingleMap/Map/Map"
type postComputeTypes = {
	freq: string;
	stationH: string;
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
	return response
};

export default function SingleMap() {
	const [popSettings, setPopSettings] = useState<boolean>(false);
	const [frequency, setFrequency] = useState<string>("1000");
	const [stationHeight, setStationHeight] = useState<string>("0");
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
		const response = await postCompute(data)
		console.log(response);
	};
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
			<Wrapper>
				<div className={styles.box}>
					<div className={styles.titleBox}>
						<Title>{id}</Title>
					</div>
					<div className={styles.mapBox}>
						<div className={styles.btnsBox}>
							<button className={styles.settingsBtn} onClick={() => setPopSettings(!popSettings)}>
								<IoMdSettings />
							</button>
						</div>
						<div className={styles.imgBox}>
							<Map/>
						</div>
						<button onClick={handleComputeBtn} className={styles.computeBtn}>
							Compute
						</button>
					</div>
				</div>
			</Wrapper>
		</>
	);
}
