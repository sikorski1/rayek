import { ReactNode, useEffect, useRef } from "react";
import styles from "./settingsDialog.module.scss";

type typeDialog = {
	children: ReactNode;
    popState: boolean
	handleOnClose?: VoidFunction;
};
export default function SettingsDialog({ children, popState, handleOnClose }: typeDialog) {
	const dialogRef = useRef<HTMLDialogElement | null>(null);

	const closeDialog = () => {
		if (!dialogRef.current) {
			return;
		}
		if (handleOnClose) handleOnClose();
		dialogRef.current?.close();
	};

	useEffect(() => {
		if (!dialogRef.current) {
			return;
		}
		dialogRef.current.showModal();

		dialogRef.current.addEventListener("close", closeDialog);
		return () => {
			dialogRef.current?.removeEventListener("close", closeDialog);
		};
	}, [popState]);

	return (
		<dialog className={styles.dialog} ref={dialogRef}>
			{children}
		</dialog>
	);
}
