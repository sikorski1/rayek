import FormField from "@/components/FormField/FormField";

import { SettingsDataTypes } from "@/types/main";
import { motion } from "framer-motion";
import styles from "./globalsettings.module.scss";
type Props = {
	formData: SettingsDataTypes;
	handleFormSubmit: (event: React.FormEvent<HTMLFormElement>) => void;
};
const itemVariants = {
	initial: { opacity: 0, x: 20 },
	animate: (i: number) => ({
		opacity: 1,
		x: 0,
		transition: {
			delay: i * 0.1,
			duration: 0.3,
			ease: "easeOut",
		},
	}),
};
export default function GlobalSettings({ formData, handleFormSubmit }: Props) {
	return (
		<form id="global-form" onSubmit={handleFormSubmit} className={styles.formBox}>
			<div className={styles.formInputBox}>
				{[
					{
						label: "RAYS AZIMUTH",
						name: "raysAzimuth",
						value: formData.numberOfRaysAzimuth,
						min: 1,
						max: 1440,
						step: 1,
					},
					{
						label: "STATION POWER (watt)",
						name: "stationPower",
						value: formData.stationPower,
						min: 0.1,
						max: 100,
						step: 0.1,
					},
					{ label: "FREQUENCY (GHz)", name: "frequency", value: formData.frequency, min: 0.1, max: 100, step: 0.1 },
					{
						label: "REFLECTION FACTOR",
						name: "relfectionFactor",
						value: formData.reflectionFactor,
						min: 0,
						max: 1,
						step: 0.01,
					},
					{
						label: "DIFFRACTION RAY NUMBER",
						name: "diffractionRayNumber",
						value: formData.diffractionRayNumber,
						min: 4,
						max: 120,
						step:1,
					},
				].map((field, i) => (
					<FormField
						key={field.name}
						index={i}
						label={field.label}
						name={field.name}
						defaultValue={field.value}
						placeholder={`Enter ${field.label.toLowerCase()}`}
						min={field.min}
						max={field.max}
						step={field.step}
						required
					/>
				))}
			</div>

			<div className={styles.formInputBox}>
				{[
					{
						label: "RAYS ELEVATION",
						name: "raysElevation",
						value: formData.numberOfRaysElevation,
						min: 1,
						max: 1440,
						step: 1,
					},
					{
						label: "STATION HEIGHT (m)",
						name: "stationHeight",
						value: formData.stationHeight,
						min: 0,
						max: 29,
						step: 1,
					},
					{
						label: "INTERACTIONS",
						name: "interactions",
						value: formData.numberOfInteractions,
						min: 1,
						max: 10,
						step: 1,
					},
					{
						label: "MINIMAL RAY POWER (dBm)",
						name: "minimalRayPower",
						value: formData.minimalRayPower,
						min: -160,
						max: -60,
						step: 0.01,
					},
				].map((field, i) => (
					<FormField
						key={field.name}
						index={i}
						label={field.label}
						name={field.name}
						defaultValue={field.value}
						placeholder={`Enter ${field.label.toLowerCase()}`}
						min={field.min}
						max={field.max}
						step={field.step}
						required
					/>
				))}
			</div>

			<motion.button
				className={styles.submitBtn}
				type="submit"
				variants={itemVariants}
				custom={5}
				initial="initial"
				animate="animate">
				Save
			</motion.button>
		</form>
	);
}
