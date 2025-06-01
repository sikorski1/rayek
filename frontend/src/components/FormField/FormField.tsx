import React from "react";
import styles from "./formfield.module.scss";

interface FormFieldProps {
	label: string;
	name: string;
	type?: string;
	defaultValue?: number;
	placeholder?: string;
	min?: number;
	max?: number;
	step?: number;
	required?: boolean;
}

import { motion } from "framer-motion";

const labelVariants = {
	initial: { opacity: 0, y: 20 },
	animate: (custom: number) => ({
		opacity: 1,
		y: 0,
		transition: { delay: custom * 0.1, duration: 0.3 },
	}),
};

const inputVariants = labelVariants;

const FormField: React.FC<FormFieldProps & { index?: number }> = ({
	label,
	name,
	type = "number",
	defaultValue,
	placeholder,
	min,
	max,
	step,
	required = false,
	index = 0,
}) => {
	return (
		<>
			<motion.label
				className={styles.label}
				variants={labelVariants}
				initial="initial"
				animate="animate"
				custom={index}>
				{label}
			</motion.label>
			<motion.input
				id={name}
				name={name}
				type={type}
				defaultValue={defaultValue}
				placeholder={placeholder}
				className={styles.input}
				min={min}
				max={max}
				step={step}
				required={required}
				variants={inputVariants}
				initial="initial"
				animate="animate"
				custom={index}
			/>
		</>
	);
};

export default FormField;
