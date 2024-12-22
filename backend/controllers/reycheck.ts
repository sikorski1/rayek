import { RequestHandler } from "express";

const computeReys: RequestHandler = (req, res, next) => {
	res.status(200).json({
		message: "Hello",
	});
    return
};

export { computeReys };
