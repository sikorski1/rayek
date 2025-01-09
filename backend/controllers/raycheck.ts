import { RequestHandler } from "express";
import { MapTypes } from "../types/main";
import { mapData } from "../data/mapData";
import {buildingsData} from "../data/buildingsData"
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
		const title = req.params.mapTitle
		const data = buildingsData
		if (data) {
			res.status(200).json({buildingsData:data})
			return
		}
		else {
			res.status(404).json({
				message: `Buildings data with title '${title}' not found.`
			})
		}
	} catch (error) {
		next(error)
	}
}

const computeRays: RequestHandler = (req, res, next) => {
	res.status(200).json({
		message: "Hello",
	});
	return;
};

export { computeRays, getMapConfiguration, getBuildingsData };
