export function getNormalizedValueFromLabel(label: string): number {
	switch (label) {
		case "< 0dbm":
			return 1.0;
		case "< -20dbm":
			return 0.85;
		case "< -40dbm":
			return 0.7;
		case "< -60dbm":
			return 0.55;
		case "< -80dbm":
			return 0.4;
		case "< -100dbm":
			return 0.25;
		case "< -120dbm":
			return 0.1;
		case "< -140dbm":
			return 0.02;
		default:
			return 0; 
	}
}