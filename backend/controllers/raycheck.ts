import { RequestHandler } from "express";

const computeRays: RequestHandler = (req, res, next) => {
	res.status(200).json({
		message: "Hello",
	});
    return
};

export { computeRays };
