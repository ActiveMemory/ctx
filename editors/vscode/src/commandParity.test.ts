/**
 * Command-parity guard.
 *
 * The chat participant once shipped commands that dispatched to `ctx`
 * subcommands the binary does not have. Unit tests could not see it: they
 * mock `execFile`, so any argv "passes". This suite closes the gap in two
 * halves:
 *
 * 1. Here: package.json commands == the dispatcher's commands; every
 *    skill-backed command bundles the skill it names; follow-ups and
 *    natural-language routes target real commands. Then every scenario
 *    below is driven through the real chat handler and each `ctx` argv it
 *    produces is recorded in `ctx-cli-surface.json` (a file snapshot, so
 *    an unreviewed change to what the extension runs fails this test).
 *
 * 2. In Go: internal/bootstrap/vscode_surface_test.go parses every argv in
 *    that file against the real cobra command tree (unknown commands,
 *    unknown flags, bad positional arguments, missing required fields) and
 *    checks every listed skill exists. It runs in the main Go test job, so
 *    a CLI rename that strands a VS Code command fails CI there.
 *
 * Refresh the snapshot after changing what a command runs:
 *   npx vitest run -u      (from editors/vscode)
 */
import { describe, it, expect, vi, beforeAll } from "vitest";
import * as cp from "child_process";
import * as fs from "fs";
import * as path from "path";

vi.mock("vscode", async () => (await import("./vscodeMock")).createVscodeMock());
vi.mock("child_process");

import {
  handler,
  CLI_COMMANDS,
  SKILLS,
  FREEFORM,
  FOLLOWUPS,
  BACKGROUND_INVOCATIONS,
} from "./extension";

// vitest runs with the extension package as the working directory.
const SKILLS_DIR = path.resolve("..", "..", "internal", "assets", "claude", "skills");

const manifest: string[] = JSON.parse(fs.readFileSync("package.json", "utf8"))
  .contributes.chatParticipants[0].commands.map((c: { name: string }) => c.name);

/**
 * Prompts that exercise every branch that runs ctx. A new command or a new
 * branch needs a row here, or its argv is never validated.
 */
const SCENARIOS: Record<string, string[]> = {
  init: [""],
  status: [""],
  agent: ["", "--budget 4000"],
  drift: [""],
  recall: ["", "--limit 5", "show 1a2b3c4d"],
  setup: ["", "claude-code", "cursor preview"],
  add: [
    'task Fix login bug --section "Phase 1"',
    "task Fix login bug --section Maintenance --priority high",
    'decision Use PostgreSQL --context "Need a reliable DB" --rationale "ACID and JSON" --consequence "Ops training"',
    'learning Go embed is package-local --context "Embedding a parent dir failed" --lesson "No parent paths" --application "Keep assets beside the package"',
    "convention Use camelCase for functions --section Naming",
  ],
  decision: [""],
  learning: [""],
  load: [""],
  compact: [""],
  sync: [""],
  task: ["complete 3", "archive", "snapshot", "snapshot pre-refactor"],
  remind: ["", "add Check CI", "Check CI", "dismiss 2", "dismiss"],
  pad: [
    "",
    "add a secret note",
    "show 1",
    "rm 1 2",
    "edit 1 new text",
    "mv 1 3",
    "resolve",
    "import notes.txt",
    "export",
    "export out",
    "merge a.enc b.enc",
  ],
  notify: ["setup", "test", "build done --event build"],
  system: ["resources", "bootstrap", "stats", "message", "message show check-freshness stale"],
  memory: ["sync", "status", "diff", "import", "publish", "unpublish"],
  journal: ["site", "obsidian"],
  doctor: [""],
  config: ["switch dev", "status", "schema"],
  why: ["", "manifesto"],
  change: ["", "--since 2h"],
  guide: ["", "--skills", "--commands"],
  permission: ["snapshot", "restore"],
  pause: [""],
  resume: [""],
  ...Object.fromEntries(Object.keys(SKILLS).map((c) => [c, [""]])),
};

describe("manifest ↔ dispatcher", () => {
  it("declares exactly the commands it dispatches", () => {
    const dispatched = [...Object.keys(CLI_COMMANDS), ...Object.keys(SKILLS)];
    expect([...manifest].sort()).toEqual([...dispatched].sort());
    expect(new Set(dispatched).size).toBe(dispatched.length);
  });

  it("routes follow-ups and natural language only to real commands", () => {
    const targets = [
      ...Object.keys(FOLLOWUPS).filter((c) => c !== "help"),
      ...Object.values(FOLLOWUPS).flat().map((f) => f.command),
      ...FREEFORM.map(([, command]) => command),
    ];
    expect(targets.filter((c) => !c || !manifest.includes(c))).toEqual([]);
  });
});

describe("skill-backed commands", () => {
  it.each(Object.entries(SKILLS))("/%s bundles the skill it names", (command, def) => {
    expect(def.skill).toBe(`ctx-${command}`);
    const file = path.join(SKILLS_DIR, def.skill, "SKILL.md");
    expect(fs.existsSync(file), `${file} missing`).toBe(true);
    expect(def.text).toBe(fs.readFileSync(file, "utf8"));
  });
});

describe("ctx CLI surface", () => {
  const argvs: string[][] = [];

  beforeAll(async () => {
    vi.mocked(cp.execFile).mockImplementation(((
      cmd: string,
      args: string[],
      _opts: unknown,
      cb: (e: unknown, out: string, err: string) => void
    ) => {
      if (cmd === "git") {
        cb(null, args.includes("--abbrev-ref") ? "main\n" : "abc1234\n", "");
      } else {
        // `--version` is the bootstrap probe, not a command.
        if (args[0] !== "--version") {
          argvs.push(args);
        }
        cb(null, "ok", "");
      }
      return { kill: () => {} };
    }) as never);
    const model = {
      sendRequest: async () => ({ text: (async function* () {})() }),
    };
    const token = { isCancellationRequested: false, onCancellationRequested: () => ({ dispose() {} }) };
    const stream = { markdown: () => {}, progress: () => {} };
    for (const [command, prompts] of Object.entries(SCENARIOS)) {
      for (const prompt of prompts) {
        await handler(
          { command, prompt, model, references: [] } as never,
          { history: [] } as never,
          stream as never,
          token as never
        );
      }
    }
  });

  it("has a scenario for every command", () => {
    expect(manifest.filter((c) => !SCENARIOS[c])).toEqual([]);
  });

  it("matches ctx-cli-surface.json (refresh: npx vitest run -u)", async () => {
    const unique = new Map<string, string[]>();
    for (const argv of [...argvs, ...BACKGROUND_INVOCATIONS]) {
      unique.set(JSON.stringify(argv), argv);
    }
    const invocations = [...unique.keys()].sort();
    const skills = Object.values(SKILLS).map((s) => s.skill).sort();
    const surface =
      "{\n" +
      '  "_generated": "By editors/vscode/src/commandParity.test.ts (refresh: npx vitest run -u). ' +
      'Validated against the real ctx command tree by internal/bootstrap/vscode_surface_test.go.",\n' +
      `  "skills": ${JSON.stringify(skills)},\n` +
      '  "invocations": [\n' +
      invocations.map((argv) => `    ${argv}`).join(",\n") +
      "\n  ]\n}\n";
    await expect(surface).toMatchFileSnapshot("./ctx-cli-surface.json");
  });
});
