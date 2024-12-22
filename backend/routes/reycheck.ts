import { Router } from "express";
import { computeReys } from "../controllers/reycheck";
const router = Router();
router.post("/compute", computeReys);

export default router;
