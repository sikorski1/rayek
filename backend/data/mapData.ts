import { MapTypes } from "../types/main";

export const mapData: MapTypes[] = [
	{
		title: "AGH fragment",
		coordinates: [
			[
				[19.914029, 50.065311], // Southwest
				[19.917527, 50.065311], // Southeast
				[19.917527, 50.067556], // Northeast
				[19.914029, 50.067556], // Northwest
				[19.914029, 50.065311], // Zamknięcie pętli
			],
		],
		center: [19.915778, 50.0664335],
		bounds: [
			[19.914029, 50.065311], // Southwest corner (dolny lewy róg)
			[19.917527, 50.067556], // Northeast corner (górny prawy róg)
		],
	},
];