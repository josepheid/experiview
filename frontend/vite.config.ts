/// vite.config.ts
import { defineConfig } from "vite";
import react from "@vitejs/plugin-react";
import { configDefaults } from "vitest/config"; // ✅ Ensure correct import

export default defineConfig({
    plugins: [react()],
    test: {
        globals: true,
        environment: "jsdom",
        setupFiles: "./setupTests.ts",
        exclude: [...configDefaults.exclude, "e2e/**"], // Optional
    },
});
