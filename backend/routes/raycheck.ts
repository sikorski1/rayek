import { Router } from "express";
import { getMapConfiguration, computeRays, getBuildingsData } from "../controllers/raycheck";
const router = Router();
router.get("/:mapTitle", getMapConfiguration);
router.get("/buildings/:mapTitle", getBuildingsData);
router.post("/compute", computeRays);

export default router;
