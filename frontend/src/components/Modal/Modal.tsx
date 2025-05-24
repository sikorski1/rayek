import { AnimatePresence, motion } from "framer-motion";
import { useEffect, useState } from "react";
import { createPortal } from "react-dom";
import styles from "./modal.module.scss";
import { IoMdClose } from "react-icons/io"; 

type ModalProps = {
	children: React.ReactNode;
	onClose: () => void;
};

export default function Modal({ children, onClose }: ModalProps) {
	const [isVisible, setIsVisible] = useState(true);

	useEffect(() => {
		const onEsc = (e: KeyboardEvent) => e.key === "Escape" && handleClose();
		window.addEventListener("keydown", onEsc);
		return () => window.removeEventListener("keydown", onEsc);
	}, []);

	const handleClose = () => {
		setIsVisible(false);
	};

	const modalRoot = document.getElementById("modal");
	if (!modalRoot) return null;

	return createPortal(
		<AnimatePresence onExitComplete={onClose}>
			{isVisible && (
				<motion.div
					initial={{ opacity: 0 }}
					animate={{ opacity: 1 }}
					exit={{ opacity: 0 }}
					className={styles.modalOverlay}
					onClick={handleClose}
				>
					<motion.div
						initial={{ scale: 0.8, opacity: 0 }}
						animate={{ scale: 1, opacity: 1 }}
						exit={{ scale: 0.8, opacity: 0 }}
						transition={{ duration: 0.2 }}
						className={styles.modalContent}
						onClick={(e) => e.stopPropagation()}
					>
						<div className={styles.closeBtnBox}>
							<button className={styles.closeBtn} onClick={handleClose}>
								<IoMdClose />
							</button>
						</div>
						{children}
					</motion.div>
				</motion.div>
			)}
		</AnimatePresence>,
		modalRoot
	);
}