import { hashPassword } from "@/utils/hashPassword";
import { FormEvent, useState } from "react";
import { useNavigate } from "react-router-dom";
import styles from "./login.module.scss";

export default function Login() {
	const [password, setPassword] = useState("");
	const [error, setError] = useState("");
	const navigate = useNavigate();
	const CORRECT_PASSWORD = import.meta.env.VITE_PASSWORD;

	const handleSubmit = async (e:FormEvent<HTMLFormElement>) => {
		e.preventDefault();

		const enteredHash = await hashPassword(password);
		const correctHash = await hashPassword(CORRECT_PASSWORD);

		if (enteredHash === correctHash) {
			localStorage.setItem("authHash", enteredHash);
			navigate("/");
		} else {
			setError("Invalid password");
		}
	};

	return (
		<div className={styles.loginBox}>
			<form onSubmit={handleSubmit} className={styles.form}>
				<h1 className={styles.title}>Rayek Access</h1>
				<input
					type="password"
					placeholder="Enter password"
					value={password}
					onChange={e => setPassword(e.target.value)}
					className={styles.input}
				/>
				{error ? <p className={styles.error}>{error}</p> : <div className={styles.error}></div>}
				<button type="submit" className={styles.btn}>
					Login
				</button>
			</form>
		</div>
	);
}
