import { AnimatePresence, motion } from "framer-motion";
import { useEffect, useState } from "react";
import styles from "./pixelbg.module.scss";

const anim = {
	initial: { opacity: 1 },
	open: (index: number) => ({
		opacity: 1,
		transition: { duration: 0.3, delay: 0.02 * index },
	}),
	closed: (index: number) => ({
		opacity: 0,
		transition: { duration: 0.3, delay: 0.02 * index },
	}),
};

export default function PixelBackground() {
	const [visible, setVisible] = useState(true);

	useEffect(() => {
		const timer = setTimeout(() => setVisible(false), 50);
		return () => clearTimeout(timer);
	}, []);

	const shuffle = (a: number[]) => {
		for (let i = a.length - 1; i > 0; i--) {
			const j = Math.floor(Math.random() * (i + 1));
			[a[i], a[j]] = [a[j], a[i]];
		}
		return a;
	};

	const getBlocks = (index: number) => {
		if (!visible) return null; // <- kluczowa linia
		const { innerWidth, innerHeight } = window;
		const blockSize = innerWidth * 0.05;
		const amountOfBlocks = Math.ceil(innerHeight / blockSize);
		const delays = shuffle([...Array(amountOfBlocks)].map((_, i) => i));
		return delays.map((randomDelay, i) => (
			<motion.div
				key={`block-${index}-${i}`}
				initial="initial"
				animate="closed"
				exit="closed"
				variants={anim}
				custom={index + randomDelay}
				className={styles.block}
			/>
		));
	};

	return (
		<div className={styles.pixelBackground}>
			<AnimatePresence>
				{visible && [...Array(20)].map((_, index) => (
					<div key={index} className={styles.column}>
						{getBlocks(index)}
					</div>
				))}
			</AnimatePresence>
		</div>
	);
}