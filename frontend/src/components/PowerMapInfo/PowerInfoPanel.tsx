import { getHeatMapColor } from "@/utils/getHeatMapColor";
import { normalizePower } from "@/utils/normalizePower";
import { useEffect } from "react";
import styles from "./powermapinfo.module.scss";

interface PowerInfoPanelProps {
	powerInfo: {
		lng: number;
		lat: number;
		power: number;
		x: number;
		y: number;
		buildingHeight: number | null;
	} | null;
	height: number;
	onClose: () => void;
}

export default function PowerInfoPanel({ powerInfo, height, onClose }: PowerInfoPanelProps) {
	useEffect(() => {
		const handleEscape = (e: KeyboardEvent) => {
			if (e.key === "Escape") {
				onClose();
			}
		};

		window.addEventListener("keydown", handleEscape);
		return () => window.removeEventListener("keydown", handleEscape);
	}, [onClose]);

	if (!powerInfo) return null;

	const normalized = normalizePower(powerInfo.power);
	const color = getHeatMapColor(normalized);
	return (
		<div className={styles.panel}>
			<div className={styles.header}>
				<strong className={styles.title}>Signal Power</strong>
				<button onClick={onClose} className={styles.closeButton} aria-label="Close">
					×
				</button>
			</div>
			<div className={styles.content}>
				<div className={styles.row}>
					<span className={styles.label}>Power:</span>{" "}
					{powerInfo.power > 0 || powerInfo.power <= -160 ? (
						<strong>0 dBm</strong>
					) : (
						<strong
							style={{
								color: `rgb(${color.r * 255}, ${color.g * 255}, ${color.b * 255})`,
							}}>
							{powerInfo.power === -160 ? "Wall" : `${powerInfo.power.toFixed(1)} dBm`}
						</strong>
					)}
				</div>
				<div className={styles.row}>
					<span className={styles.label}>Height:</span> <strong>{height}m</strong>
				</div>
				{powerInfo.buildingHeight !== null && (
					<div className={styles.row}>
						<span className={styles.label}>Building:</span>
						<strong className={styles.buildingHeight}>{powerInfo.buildingHeight}m</strong>
					</div>
				)}
				<div className={styles.coordinates}>
					Coordinates: {powerInfo.lat.toFixed(6)}, {powerInfo.lng.toFixed(6)}
				</div>
				<div className={styles.index}>
					Index: ({powerInfo.x}, {powerInfo.y})
				</div>
			</div>
		</div>
	);
}
