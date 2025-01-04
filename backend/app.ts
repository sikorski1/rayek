import express, { ErrorRequestHandler, Express } from "express";
import RaycheckRouter from "./routes/raycheck";

import cors from "cors";
const app: Express = express();
const port = 3001;

app.use(
	cors({
		origin: "http://localhost:5173",
		credentials: true,
		allowedHeaders: ["Origin", "X-Requested-With", "Content-Type", "Authorization", "Accept", "input"],
		methods: ["GET", "POST", "PUT", "PATCH", "DELETE"],
	})
);

app.use(express.json());

app.use("/raycheck", RaycheckRouter);

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
