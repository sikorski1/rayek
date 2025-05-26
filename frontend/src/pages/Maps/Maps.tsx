import sampleImg from "@/assets/imgs/sampleMapImg.jpg";
import PageTransition from "@/components/PageTransition/PageTransition";
import Title from "@/components/Title/Title";
import Wrapper from "@/components/Wrapper/Wrapper";
import { useGetMaps } from "@/hooks/useMap";
import { AnimatePresence, motion } from "framer-motion";
import { Link } from "react-router-dom";
import styles from "./maps.module.scss";
type MapData = {
	id: string;
	name: string;
	description: string;
	img: string;
};
export default function Maps() {
	const { data, isLoading, error } = useGetMaps();
	return (
		<PageTransition>
			<Wrapper>
				<div className={styles.box}>
					<Title>Available Maps</Title>
					<div className={styles.cardsBox}>
						<AnimatePresence>
							{data &&
								data.map((item: MapData, index: number) => (
									<motion.div
										className={styles.card}
										key={item.id}
										initial={{ opacity: 0, y: -20 }}
										animate={{ opacity: 1, y: 0 }}
										transition={{ duration: 0.4, delay: 0.3 + index * 0.1 }}>
										<motion.div
											className={styles.titleBox}
											initial={{ opacity: 0, x: -20 }}
											animate={{ opacity: 1, x: 0 }}
											transition={{ duration: 0.2, delay: 0.4 + index * 0.1 }}>
											<h3 className={styles.cardTitle}>{item.name}</h3>
										</motion.div>
										<div className={styles.cardBox}>
											<Link
												key={item.id}
												to={item.name.replace(/\s+/g, "").toLocaleLowerCase()}
												className={styles.link}></Link>
											<div className={styles.bgGradient}></div>
											<div className={styles.imgBox}>
												<img className={styles.img} src={item.img || sampleImg} alt="map image" />
											</div>
										</div>
									</motion.div>
								))}
						</AnimatePresence>
					</div>
				</div>
			</Wrapper>
		</PageTransition>
	);
}
