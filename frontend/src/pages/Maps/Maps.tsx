import aghWiet from "@/assets/imgs/agh_fragment.png";
import sampleImg from "@/assets/imgs/sampleMapImg.jpg";
import Title from "@/components/Title/Title";
import Wrapper from "@/components/Wrapper/Wrapper";
import { Link } from "react-router-dom";
import styles from "./maps.module.scss";
type MapData = {
	id: string;
	name: string;
	description: string;
	img: string;
};
export default function Maps() {
	const sampleData: MapData[] = [
		{
			id: "1",
			name: "AGH WIET",
			description: "",
			img: aghWiet,
		},
		{
			id: "2",
			name: "AGH Library",
			description: "",
			img: "",
		},
		{
			id: "3",
			name: "Cracow Nowa Huta",
			description: "",
			img: "",
		},
		{
			id: "4",
			name: "Cracow Wawel",
			description: "",
			img: "",
		},
	];
	return (
		<Wrapper>
			<div className={styles.box}>
				<Title>Available Maps</Title>
				<div className={styles.cardsBox}>
					{sampleData.map((item: MapData) => (
						<div className={styles.card} key={item.id}>
							<div className={styles.titleBox}>
								<h3 className={styles.cardTitle}>{item.name}</h3>
							</div>
							<div className={styles.cardBox}>
								<Link key={item.id} to={item.name.replace(/\s+/g, "").toLocaleLowerCase()} className={styles.link}></Link>
								<div className={styles.bgGradient}></div>
								<div className={styles.imgBox}>
									<img className={styles.img} src={item.img || sampleImg} alt="map image" />
								</div>
							</div>
						</div>
					))}
				</div>
			</div>
		</Wrapper>
	);
}
