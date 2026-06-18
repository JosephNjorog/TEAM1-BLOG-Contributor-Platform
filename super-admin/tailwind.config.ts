import type { Config } from "tailwindcss";
import sharedPreset from "@team1/shared/tailwind-preset";

export default {
  presets: [sharedPreset],
  content: [
    "./index.html",
    "./src/**/*.{ts,tsx}",
    "../packages/shared/src/**/*.{ts,tsx}",
  ],
} satisfies Config;
