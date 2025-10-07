import PageTransition from "@/components/PageTransition/PageTransition";
import { Link } from "react-router-dom";
import styles from "./home.module.scss";
export default function Home() {
	return (
		<PageTransition>
			<main className={styles.box}>
				<div className={styles.container}>
					<div className={styles.bgImage}></div>
					<h1 className={styles.name}>Rayek</h1>
					<div className={styles.buttonsBox}>
						<Link className={styles.link} to="/maps">
							Maps
						</Link>
					</div>
				</div>
			</main>
		</PageTransition>
	);
}
