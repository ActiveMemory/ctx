import { defineConfig } from "vitest/config";

export default defineConfig({
  // Mirror esbuild's `--loader:.md=text` so the bundled SKILL.md imports
  // resolve under test exactly as they do in the build.
  plugins: [
    {
      name: "md-as-text",
      transform(code, id) {
        return id.endsWith(".md") ? `export default ${JSON.stringify(code)};` : undefined;
      },
    },
  ],
  test: {
    include: ["src/**/*.test.ts"],
  },
});
