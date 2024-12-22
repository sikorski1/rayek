import { Router } from "express";
import { computeRays } from "../controllers/raycheck";
const router = Router();
router.post("/compute", computeRays);

export default router;
