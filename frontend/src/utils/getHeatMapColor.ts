import { Color } from "three";
export function getHeatMapColor(value: number): Color {
	let r = 0,
		g = 0,
		b = 0;
	if (value === -150) {
		return new Color(0, 0, 0);
	}
	if (value < 0.25) {
		r = 0;
		g = 255 * value * 4;
		b = 255;
	} else if (value < 0.5) {
		r = 0;
		g = 255;
		b = 255 * (2 - value * 6);
	} else if (value < 0.75) {
		r = 255 * (value * 4 - 2);
		g = 255;
		b = 0;
	} else {
		r = 255;
		g = 255 * (4 - value * 4);
		b = 0;
	}

	return new Color(r / 255, g / 255, b / 255);
}
