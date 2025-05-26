import PixelBackground from "../PixelBg/PixelBackground";

export default function PageTransition({ children }: { children: React.ReactNode }) {
	return (
		<>
			<PixelBackground />
			{children}
		</>
	);
}
