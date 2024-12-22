import express from "express";
import { Express } from "express";
import ReycheckRouter from "./routes/reycheck";
import { ErrorRequestHandler } from "express";
const app: Express = express();
const port = 3001;

app.use(express.json());

app.use("/reycheck", ReycheckRouter);

app.use(((error, req, res, next) => {
	const status = error.statusCode || 500;
	const message = error.message;
	res.status(status).json({
		status: status,
		message: message,
	});
}) as ErrorRequestHandler);

app.listen(port, () => {
	console.log(`Backend running at http://localhost:${port}`);
});
