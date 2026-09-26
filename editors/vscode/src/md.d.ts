// SKILL.md files are imported as text: esbuild bundles them with
// `--loader:.md=text` and vitest.config.ts mirrors that loader.
declare module "*.md" {
  const text: string;
  export default text;
}
