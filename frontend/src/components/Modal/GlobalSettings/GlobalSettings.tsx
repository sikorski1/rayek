import FormField from "@/components/FormField/FormField";

import ToolTip from "@/components/ToolTip/ToolTip";
import { SettingsDataTypes } from "@/types/main";
import { motion } from "framer-motion";
import { useState } from "react";
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
const MAX_TOTAL_RAYS = 2880;
export default function GlobalSettings({ formData, handleFormSubmit }: Props) {
	const [raysAzimuth, setRaysAzimuth] = useState(formData.numberOfRaysAzimuth);
	const [raysElevation, setRaysElevation] = useState(formData.numberOfRaysElevation);

	const maxAzimuth = MAX_TOTAL_RAYS - raysElevation;
	const maxElevation = MAX_TOTAL_RAYS - raysAzimuth;
	return (
		<form id="global-form" onSubmit={handleFormSubmit} className={styles.formBox}>
			<div className={styles.formInputBox}>
				{[
					{
						label: "RAYS AZIMUTH",
						name: "raysAzimuth",
						value: raysAzimuth,
						min: 1,
						max: maxAzimuth,
						step: 1,
						onChange: (value: number) => setRaysAzimuth(value),
						toolTipText:
							"Defines the number of rays distributed horizontally (azimuth plane). Higher values increase precision but also computation time.",
					},
					{
						label: "STATION POWER (watt)",
						name: "stationPower",
						value: formData.stationPower,
						min: 0.01,
						max: 100,
						step: 0.01,
						toolTipText:
							"Specifies the transmission power of the station in watts. Affects the signal strength and propagation range.",
					},
					{
						label: "FREQUENCY (GHz)",
						name: "frequency",
						value: formData.frequency,
						min: 0.1,
						max: 100,
						step: 0.1,
						toolTipText:
							"Sets the transmission frequency in gigahertz. Higher frequencies offer more bandwidth but weaker signal range.",
					},
					{
						label: "REFLECTION FACTOR",
						name: "relfectionFactor",
						value: formData.reflectionFactor,
						min: 0,
						max: 1,
						step: 0.01,
						toolTipText:
							"Determines how much of the signal is reflected by surfaces. A value of almost 0 means no reflection, 1 means full reflection.",
					},
					{
						label: "DIFFRACTION RAY NUMBER",
						name: "diffractionRayNumber",
						value: formData.diffractionRayNumber,
						min: 0,
						max: 120,
						step: 1,
						toolTipText:
							"Controls how many rays are used to model diffraction effects around obstacles. More rays increase accuracy but slow down computation.",
					},
				].map((field, i) => (
					<div key={field.name} className={styles.singleInputBox}>
						<ToolTip name={field.name} index={i} toolTipText={field.toolTipText} place="top-end" />
						<FormField
							index={i}
							label={field.label}
							name={field.name}
							defaultValue={field.value}
							value={field.value}
							placeholder={`Enter ${field.label.toLowerCase()}`}
							min={field.min}
							max={field.max}
							step={field.step}
							onChange={field.onChange}
							required
						/>
					</div>
				))}
			</div>

			<div className={styles.formInputBox}>
				{[
					{
						label: "RAYS ELEVATION",
						name: "raysElevation",
						value: raysElevation,
						min: 1,
						max: maxElevation,
						step: 1,
						onChange: (value: number) => setRaysElevation(value),
						toolTipText:
							"Defines the number of rays distributed vertically (elevation plane). More rays improve 3D accuracy but require more processing.",
					},
					{
						label: "STATION HEIGHT (m)",
						name: "stationHeight",
						value: formData.stationHeight,
						min: 1,
						max: 29,
						step: 1,
						toolTipText:
							"Specifies the station’s height above ground level in meters. It affects line-of-sight and coverage area.",
					},
					{
						label: "INTERACTIONS",
						name: "interactions",
						value: formData.numberOfInteractions,
						min: 1,
						max: 10,
						step: 1,
						toolTipText:
							"Sets the maximum number of reflections, diffractions, or transmissions a ray can undergo. Higher values simulate more complex paths.",
					},
					{
						label: "MINIMAL RAY POWER (dBm)",
						name: "minimalRayPower",
						value: formData.minimalRayPower,
						min: -160,
						max: -60,
						step: 0.01,
						toolTipText:
							"Defines the minimum power threshold in dBm for rays to be considered. Rays weaker than this value are ignored.",
					},
				].map((field, i) => (
					<div key={field.name} className={styles.singleInputBox}>
						<ToolTip name={field.name} index={i} toolTipText={field.toolTipText} place="top-end" />
						<FormField
							index={i}
							label={field.label}
							name={field.name}
							defaultValue={field.value}
							value={field.value}
							placeholder={`Enter ${field.label.toLowerCase()}`}
							min={field.min}
							max={field.max}
							step={field.step}
							onChange={field.onChange}
							required
						/>
					</div>
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
