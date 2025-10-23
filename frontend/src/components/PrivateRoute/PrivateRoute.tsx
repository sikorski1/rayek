import { hashPassword } from "@/utils/hashPassword";
import { ReactNode, useEffect, useState } from "react";
import { Navigate } from "react-router-dom";

export default function PrivateRoute({ children }: { children: ReactNode }) {
	const CORRECT_PASSWORD = import.meta.env.VITE_PASSWORD;
	const storedHash = localStorage.getItem("authHash");

	const [isValid, setIsValid] = useState<boolean | null>(null);

	useEffect(() => {
		(async () => {
			if (!storedHash) return setIsValid(false);
			const correctHash = await hashPassword(CORRECT_PASSWORD);
			setIsValid(storedHash === correctHash);
		})();
	}, [storedHash, CORRECT_PASSWORD]);

	if (isValid === null) return null;
	return isValid ? children : <Navigate to="/login" replace />;
}
