import Title from "@/components/Title/Title";
import Wrapper from "@/components/Wrapper/Wrapper";
import { Link } from "react-router-dom";
import styles from "./maps.module.scss";
export default function Maps() {
	return (
		<Wrapper>
			<div className={styles.box}>
				<Title>Available Maps</Title>
				<div>
					<Link to="1" className={styles.link}>
						First map
					</Link>
					<Link to="2" className={styles.link}>
						Second map
					</Link>
					<Link to="3" className={styles.link}>
						Third map
					</Link>
				</div>
			</div>
		</Wrapper>
	);
}
