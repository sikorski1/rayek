import Wrapper from "@/components/Wrapper/Wrapper";
import { Link } from "react-router-dom";
import styles from "./home.module.scss"
export default function Home() {
	return (
		<Wrapper>
			<div className={styles.box}>
				<h1 className={styles.name}>RayCheck</h1>
				<Link className={styles.link} to="/maps">Maps</Link>
			</div>
		</Wrapper>
	);
}
