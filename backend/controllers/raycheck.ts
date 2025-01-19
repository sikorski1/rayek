import { RequestHandler } from "express";
import { buildingsData } from "../data/buildingsData";
import { mapData } from "../data/mapData";
import { MapTypes } from "../types/main";
const getMapConfiguration: RequestHandler = (req, res, next) => {
	try {
		const title = req.params.mapTitle;
		const data = mapData.find((map: MapTypes) => map.title === title);
		if (data) {
			res.status(200).json({ mapConf: data });
		} else {
			res.status(404).json({
				message: `Map configuration with title '${title}' not found.`,
			});
		}
	} catch (error) {
		next(error);
	}
};

const getBuildingsData: RequestHandler = (req, res, next) => {
	try {
		const title = req.params.mapTitle;
		const data = buildingsData;
		if (data) {
			res.status(200).json({ buildingsData: data });
			return;
		} else {
			res.status(404).json({
				message: `Buildings data with title '${title}' not found.`,
			});
		}
	} catch (error) {
		next(error);
	}
};

const computeRays: RequestHandler = (req, res, next) => {
	const southwest = [19.914029, 50.065311]
	const southeast = [19.917527, 50.065311]
	const northeast = [19.917527, 50.067556]
	const northwest = [19.914029, 50.067556]
	const numPoints = 500;

	const xStep = (southeast[0] - southwest[0]) / (numPoints - 1);
	const yStep = (northeast[1] - southeast[1]) / (numPoints - 1);
	const positions = [];
	for (let i = 0; i < numPoints; i++) {
		for (let j = 0; j < numPoints; j++) {
			positions.push([
				southwest[0] + xStep * i,
				southwest[1] + yStep * j
			]);
		}
	}
	res.status(200).json({
		message: "Hello",
		positions: positions
	});
	return;
};

export { computeRays, getBuildingsData, getMapConfiguration };
