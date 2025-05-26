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

const FormField: React.FC<FormFieldProps> = ({
	label,
	name,
	type = "number",
	defaultValue,
	placeholder,
	min,
	max,
	step,
	required = false,
}) => {
	return (
		<>
			<label htmlFor={name} className={styles.label}>
				{label}
			</label>
			<input
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
			/>
		</>
	);
};

export default FormField;
