import { Router } from "express";
import { getMapConfiguration, computeRays } from "../controllers/raycheck";
const router = Router();
router.get("/:mapTitle", getMapConfiguration);
router.post("/compute", computeRays);

export default router;
