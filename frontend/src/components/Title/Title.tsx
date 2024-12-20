import { ReactNode } from "react";
import styles from "./title.module.scss"
export default function Title({ children }: { children: ReactNode }) {
	return <h2 className={styles.title}>{children}</h2>;
}
