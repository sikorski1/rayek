import { DoorClosed } from "lucide-react";
import { useNavigate } from "react-router-dom";
import styles from "./logoutbutton.module.scss";
const LogoutButton = () => {
	const navigate = useNavigate();

	const handleLogout = () => {
		localStorage.removeItem("authHash");
		navigate("/");
	};

	return (
		<button onClick={handleLogout} className={styles.logOutBtn}>
			<DoorClosed size={40} />
		</button>
	);
};

export default LogoutButton;
