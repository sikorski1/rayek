import { ReactNode } from "react";
import styles from "./wrapper.module.scss";

export default function Wrapper({children}:{children:ReactNode}) {
	return (
		<div className={styles.wrapper}>
			{children}
		</div>
	);
}
