import FormField from "@/components/FormField/FormField";
import { SettingsDataTypes } from "@/types/main";
import { motion } from "framer-motion";
import { IoMdClose } from "react-icons/io";
import styles from "./singleraysettings.module.scss";
type Props = {
	formData: SettingsDataTypes;
	handleFormSubmit: (event: React.FormEvent<HTMLFormElement>) => void;
	handleAddRay: () => void;
	handleRemoveRay: (index: number) => void;
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
const removeBtnVariants = {
	initial: { opacity: 0, rotate: 0, scale: 0 },
	animate: (i: number) => ({
		opacity: 1,
		rotate: 360,
		scale: 1,
		transition: {
			delay: i * 0.1 + 0.2,
			duration: 0.5,
			ease: "easeOut",
		},
	}),
};
export default function SingleRaySettings({ formData, handleFormSubmit, handleAddRay, handleRemoveRay }: Props) {
	return (
		<form id="singleRay-form" onSubmit={handleFormSubmit} className={styles.formBox}>
			{formData.singleRays.length ? (
				formData.singleRays.map((ray, i) => (
					<div key={i} className={styles.singleRayBox}>
						<div style={{ display: "flex", flexDirection: "column", gap: "1rem" }}>
							<FormField
								key={ray.azimuth}
								index={i}
								label="Azimuth ray number"
								name="azimuth"
								defaultValue={ray.azimuth}
								placeholder={`Enter azimuth`}
								min={0}
								max={formData.numberOfRaysAzimuth - 1}
								step={1}
								required
							/>
						</div>
						<div style={{ display: "flex", flexDirection: "column", gap: "1rem" }}>
							<FormField
								key={ray.elevation}
								index={i}
								label="Elevation ray number"
								name="elevation"
								defaultValue={ray.elevation}
								placeholder={`Enter elevation`}
								min={0}
								max={formData.numberOfRaysElevation - 1}
								step={1}
								required
							/>
						</div>
						<motion.button
							type="button"
							onClick={() => handleRemoveRay(i)}
							className={styles.removeSingleRayBtn}
							variants={removeBtnVariants}
							custom={i}
							initial="initial"
							animate="animate">
							<IoMdClose />
						</motion.button>
					</div>
				))
			) : (
				<motion.p className={styles.noSingleRaysInfo} variants={itemVariants} initial="initial" animate="animate">
					Single rays are not configured...
				</motion.p>
			)}
			<motion.button
				type="button"
				onClick={handleAddRay}
				className={styles.addSingleRayBtn}
				disabled={formData.singleRays.length >= 4}
				variants={itemVariants}
				custom={1}
				initial="initial"
				animate="animate">
				Add ray
			</motion.button>
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
