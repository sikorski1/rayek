import sampleImg from "@/assets/imgs/sampleMapImg.jpg";
import SettingsDialog from "@/components/SettingsDialog/SettingsDialog";
import Title from "@/components/Title/Title";
import Wrapper from "@/components/Wrapper/Wrapper";
import { useState } from "react";
import { IoMdClose, IoMdSettings } from "react-icons/io";
import { useParams } from "react-router-dom";
import styles from "./singleMap.module.scss";
export default function SingleMap() {
	const [popSettings, setPopSettings] = useState<boolean>();
	const { id } = useParams();
	const handleOnSettingsClose = () => {
		setPopSettings(false);
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
						<form className={styles.form}>
							<div className={styles.formBox}>
								<div className={styles.formInputBox}>
									<label htmlFor="frequency" className={styles.label}>
										Frequency (MHz):
									</label>
									<input
										type="number"
										name="frequency"
										className={styles.input}
										placeholder="Enter frequency in MHz"
										min="0"
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
										className={styles.input}
										placeholder="Enter station height in meters"
										min="0"
										max="10"
										step="0.1"
										required
									/>
								</div>
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
							<img src={sampleImg} className={styles.img} alt="mapa" />
						</div>
					</div>
				</div>
			</Wrapper>
		</>
	);
}
