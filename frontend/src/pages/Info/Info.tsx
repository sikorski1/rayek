import PageTransition from "@/components/PageTransition/PageTransition";
import Title from "@/components/Title/Title";
import { motion } from "framer-motion";
import styles from "./info.module.scss";

export default function Info() {
	return (
		<PageTransition>
			<main className={styles.box}>
				<div className={styles.container}>
					<Title>About Rayek</Title>
					<motion.section
						className={styles.section}
						initial={{ opacity: 0, y: 40 }}
						whileInView={{ opacity: 1, y: 0 }}
						transition={{ duration: 0.6 }}
						viewport={{ once: true }}>
						<h2 className={styles.subtitle}>Overview</h2>
						<p>
							<strong>Rayek</strong> is an advanced <strong>radio wave propagation simulator</strong>,
							developed as the practical part of an engineering thesis. It combines physically based
							computation with interactive 3D visualization to model how electromagnetic waves propagate,
							reflect, and diffract in realistic environments.
						</p>
					</motion.section>

					<motion.section
						className={styles.section}
						initial={{ opacity: 0, y: 40 }}
						whileInView={{ opacity: 1, y: 0 }}
						transition={{ duration: 0.6, delay: 0.2 }}
						viewport={{ once: true }}>
						<h2 className={styles.subtitle}>Core Algorithm</h2>
						<p>
							The simulation is powered by a <strong>3D ray launching engine</strong>, capable of modeling{" "}
							<strong>multiple reflections</strong> and
							<strong> diffraction</strong> effects. Each ray represents a discrete propagation path that
							interacts with the environment, enabling the analysis of both line-of-sight and
							non-line-of-sight conditions.
						</p>
						<p>
							Calculations are performed in the backend written in <strong>Go (Golang)</strong>, which
							ensures high efficiency and concurrency when processing thousands of rays in parallel.
						</p>
					</motion.section>
					<motion.section
						className={styles.section}
						initial={{ opacity: 0, y: 40 }}
						whileInView={{ opacity: 1, y: 0 }}
						transition={{ duration: 0.6, delay: 0.4 }}
						viewport={{ once: true }}>
						<h2 className={styles.subtitle}>Configuration</h2>
						<p style={{ marginBottom: "3rem" }}>
							The simulator allows users to precisely configure parameters that influence the behavior and
							accuracy of the propagation model. Each setting affects how rays are emitted, interact, and
							decay.
						</p>
						<ul className={styles.list}>
							<li>
								<strong>Rays Azimuth</strong> — defines the number of rays distributed horizontally in
								the azimuth plane.
							</li>
							<li>
								<strong>Rays Elevation</strong> — defines the number of rays distributed vertically in
								the elevation plane.
							</li>
							<li>
								<strong>Station Height</strong> — height of the transmitting antenna above ground level.
							</li>
							<li>
								<strong>Station Power (W)</strong> — transmission power of the station in watts.
							</li>
							<li>
								<strong>Frequency (GHz)</strong> — operating frequency; higher frequencies reduce
								coverage but increase resolution.
							</li>
							<li>
								<strong>Reflection Factor</strong> — determines the proportion of the signal reflected
								from surfaces (0 = none, 1 = full).
							</li>
							<li>
								<strong>Diffraction Ray Number</strong> — controls how many rays are generated to
								simulate diffraction effects.
							</li>
							<li>
								<strong>Interactions</strong> — sets the maximum number of reflections, diffractions, or
								transmissions a ray can experience.
							</li>
							<li>
								<strong>Minimal Ray Power (dBm)</strong> — defines the threshold below which rays are
								discarded from computation.
							</li>
						</ul>
						<p>
							These parameters enable full control over simulation precision, computation time, and
							realism — making the model suitable for both fast approximations and in-depth analysis.
						</p>
					</motion.section>

					<motion.section
						className={styles.section}
						initial={{ opacity: 0, y: 40 }}
						whileInView={{ opacity: 1, y: 0 }}
						transition={{ duration: 0.6, delay: 0.6 }}
						viewport={{ once: true }}>
						<h2 className={styles.subtitle}>Simulation Results</h2>
						<p style={{ marginBottom: "3rem" }}>
							After computation, the simulator can display results in several visualization modes:
						</p>
						<ul className={styles.list}>
							<li>
								<strong>Power Heatmap (0–29 m)</strong> — shows signal strength distribution per
								height.
							</li>
							<li>
								<strong>Ray Visualization (1–4 selected rays)</strong> — displays the actual propagation
								paths in the 3D scene.
							</li>
							<li>
								<strong>Coverage Map</strong> — a side panel view presenting overall radio coverage and
								signal reach.
							</li>
						</ul>
						<p>
							This allows both macro-scale overview and detailed inspection of individual propagation
							paths.
						</p>
					</motion.section>

					<motion.section
						className={styles.section}
						initial={{ opacity: 0, y: 40 }}
						whileInView={{ opacity: 1, y: 0 }}
						transition={{ duration: 0.6, delay: 0.8 }}
						viewport={{ once: true }}>
						<h2 className={styles.subtitle}>Data & Visualization Tools</h2>
						<p>
							The simulation environment is built upon geospatial data from
							<strong> OpenStreetMap (OSM)</strong>. This data provides the base layout of buildings,
							streets, and terrain for the simulation.
						</p>
						<p>
							Due to limited access to accurate 3D elevation data, <strong>building heights</strong>
							{" "}were assigned manually. As a result, some objects may differ from their real-world
							dimensions.
						</p>
						<p>
							The visualization of the map is powered by <strong>Mapbox</strong>, enabling smooth
							navigation, dynamic rendering, and high-quality map layers during simulation and result
							analysis.
						</p>
					</motion.section>

					{/* Purpose */}
					<motion.section
						className={styles.section}
						initial={{ opacity: 0, y: 40 }}
						whileInView={{ opacity: 1, y: 0 }}
						transition={{ duration: 0.6, delay: 1.0 }}
						viewport={{ once: true }}>
						<h2 className={styles.subtitle}>Purpose</h2>
						<p>
							Rayek bridges the gap between theoretical radio propagation models and real-world
							visualization. Its goal is to make electromagnetic simulations more intuitive, interactive,
							and visually comprehensible.
						</p>
						<p>
							The project demonstrates the use of modern web technologies with physics-based modeling —
							combining <strong>Go</strong> for computation and <strong>React</strong> with{" "}
							<strong>Framer Motion</strong> for a dynamic front-end experience.
						</p>
					</motion.section>
				</div>
			</main>
		</PageTransition>
	);
}
