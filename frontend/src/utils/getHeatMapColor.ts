import { Color } from "three";
export function getHeatMapColor(value: number): Color {
  if (value === -160 || value === -1) {
    return new Color(1, 1, 1);
  }
  
  let r = 0, g = 0, b = 0;
  
  const brightness = 0.9; 
  
  if (value < 0.2) {
    const t = value / 0.2;
    r = 0;
    g = 0;
    b = (80 + 175 * t) * brightness;
  } else if (value < 0.4) {
    const t = (value - 0.2) / 0.2;
    r = 0;
    g = 200 * t * brightness;
    b = 255 * brightness;
  } else if (value < 0.6) {
    const t = (value - 0.4) / 0.2;
    r = 0;
    g = (200 + 55 * t) * brightness;
    b = 255 * (1 - t) * brightness;
  } else if (value < 0.8) {
    const t = (value - 0.6) / 0.2;
    r = 220 * t * brightness;
    g = 255 * brightness;
    b = 0;
  } else {
    const t = (value - 0.8) / 0.2;
    r = (220 + 35 * t) * brightness;
    g = 255 * (1 - t * 0.7) * brightness;
    b = 0;
  }
  
  return new Color(r / 255, g / 255, b / 255);
}