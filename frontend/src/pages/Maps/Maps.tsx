import sampleImg from "@/assets/imgs/sampleMapImg.jpg";
import PageTransition from "@/components/PageTransition/PageTransition";
import Title from "@/components/Title/Title";
import Wrapper from "@/components/Wrapper/Wrapper";
import { useGetMaps } from "@/hooks/useMap";
import { AnimatePresence, motion } from "framer-motion";
import { Link } from "react-router-dom";
import styles from "./maps.module.scss";

export default function Maps() {
	const { data } = useGetMaps();
	console.log(data);
	return (
		<PageTransition>
			<Wrapper>
				<div className={styles.box}>
					<Title>Available Maps</Title>
					<div className={styles.cardsBox}>
						<AnimatePresence>
							{data &&
								data.map((item, index: number) => (
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
											<Link key={item.id} to={item.id} className={styles.link}></Link>
											<div className={styles.imgBox}>
												<img
													className={styles.img}
													src={item.img || sampleImg}
													alt="map image"
												/>
											</div>
										</div>
										<div className={styles.sizeSign}>
											{item.size}m x {item.size}m
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
