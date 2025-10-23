export function normalizePower(value: number): number {
	if (value === -160) return -1;

	const strongSignalRange = 0.8;
	const weakSignalRange = 0.2;

	if (value >= -120) {
		const ratio = -value / 120;
		return 1.0 - ratio * strongSignalRange;
	} else {
		const ratio = (-value - 120) / 40;
		return weakSignalRange * (1.0 - ratio);
	}
}