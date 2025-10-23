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
	value?: number;
	step?: number;
	required?: boolean;
	onChange?: (value: any) => void;
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
	value,
	placeholder,
	min,
	max,
	step,
	required = false,
	index = 0,
	onChange,
}) => {
	const inputProps = onChange ? { value: value } : { defaultValue: value };

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
				{...inputProps}
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
				onChange={e => {
					if (onChange) {
						onChange(parseFloat(e.target.value) || 0);
					}
				}}
			/>
		</>
	);
};

export default FormField;
