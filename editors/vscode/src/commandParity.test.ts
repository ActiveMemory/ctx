/**
 * Command-parity guard.
 *
 * The original chat-participant rewrite shipped ~12 commands that dispatched
 * to `ctx` subcommands which do not exist on the real binary (e.g. `ctx tasks`
 * instead of `ctx task`). Nothing caught it because the unit tests mock
 * `execFile` — any argv "passes". This test closes that gap by checking the
 * dispatcher against three sources of truth at once:
 *
 *   1. package.json `contributes.chatParticipants[].commands`  (the manifest)
 *   2. the `switch (request.command)` cases                    (the dispatcher)
 *   3. the real `ctx` command tree                             (the binary)
 *
 * (1) ↔ (2) is pure text and always runs. (2)/(3): every literal `ctx …`
 * invocation in the source must resolve to a real command path in a `ctx`
 * binary built from THIS commit. Provide the binary via the `CTX_BIN`
 * environment variable (CI builds it: `go build -o ctxbin ./cmd/ctx`); the
 * test also tries `ctx` on PATH. If neither resolves, the ground-truth block
 * is skipped LOUDLY (never silently passed) so CI can be wired to guarantee it.
 */
import { describe, it, expect } from "vitest";
import { execFileSync } from "node:child_process";
import * as fs from "node:fs";
import * as path from "node:path";

// vitest runs with the extension package as the working directory.
const SRC = path.resolve("src/extension.ts");
const PKG = path.resolve("package.json");

const source = fs.readFileSync(SRC, "utf8");

// ── (1) manifest ↔ (2) dispatcher ────────────────────────────────────────────

function manifestCommands(): string[] {
  const pkg = JSON.parse(fs.readFileSync(PKG, "utf8"));
  return pkg.contributes.chatParticipants[0].commands.map(
    (c: { name: string }) => c.name
  );
}

/** The `case "x":` labels inside `switch (request.command)`, up to `default:`. */
function dispatcherCommands(): string[] {
  const start = source.indexOf("switch (request.command)");
  const end = source.indexOf("default:", start);
  if (start < 0 || end < 0) throw new Error("dispatcher switch not found");
  const body = source.slice(start, end);
  return [...body.matchAll(/case "([^"]+)":/g)].map((m) => m[1]);
}

// ── (2)/(3) invocation command paths ↔ the real binary ───────────────────────

/**
 * Every `ctx` invocation's static command-path prefix. Invocations are the
 * literal `runCtx([...])` arrays and the `args = [...]` arrays that feed them.
 * We take the leading string-literal tokens (a command path is always literal
 * and precedes any flag/variable), stopping at the first flag (`--x`), spread,
 * variable, or template — so `["recall", "show", rest]` → `["recall","show"]`.
 */
function invocationPaths(): string[][] {
  const bodies: string[] = [];
  for (const m of source.matchAll(/runCtx\(\s*\[([\s\S]*?)\]/g)) bodies.push(m[1]);
  for (const m of source.matchAll(/\bargs\s*=\s*\[([\s\S]*?)\]/g)) bodies.push(m[1]);

  const paths: string[][] = [];
  for (const body of bodies) {
    const path: string[] = [];
    for (const raw of body.split(",")) {
      const t = raw.trim();
      const lit = t.match(/^["']([A-Za-z][A-Za-z0-9-]*)["']$/); // word literal, not a -flag
      if (!lit) break;
      path.push(lit[1]);
    }
    if (path.length) paths.push(path);
  }
  // Dedupe by joined path.
  const seen = new Set<string>();
  return paths.filter((p) => {
    const k = p.join(" ");
    if (seen.has(k)) return false;
    seen.add(k);
    return true;
  });
}

function resolveCtx(): string | null {
  const candidates = [process.env.CTX_BIN, "ctx"].filter(Boolean) as string[];
  for (const bin of candidates) {
    try {
      execFileSync(bin, ["--version"], { stdio: "pipe" });
      return bin;
    } catch {
      /* try next */
    }
  }
  return null;
}

/**
 * Whether `ctx <path…>` is a real command. Probes the binary directly (so
 * HIDDEN commands like `ctx system` are recognized) and requires the command's
 * own `Usage:` line to name the full path — a parent with a `RunE` that accepts
 * free args would otherwise echo its own help for a bogus subcommand and read
 * as a false positive. An unknown command exits non-zero → false. (Memoized.)
 */
const existsCache = new Map<string, boolean>();
function pathExists(ctxBin: string, path: string[]): boolean {
  const key = path.join(" ");
  const cached = existsCache.get(key);
  if (cached !== undefined) return cached;
  let ok = false;
  try {
    const out = execFileSync(ctxBin, [...path, "--help"], {
      encoding: "utf8",
      stdio: ["ignore", "pipe", "pipe"],
    });
    ok = out.includes(`ctx ${key}`);
  } catch {
    ok = false; // non-zero exit = unknown command
  }
  existsCache.set(key, ok);
  return ok;
}

/** Returns the shallowest sub-path that isn't a real command, else null. */
function firstDeadSegment(ctxBin: string, path: string[]): string | null {
  for (let i = 1; i <= path.length; i++) {
    const sub = path.slice(0, i);
    if (!pathExists(ctxBin, sub)) return sub.join(" ");
  }
  return null;
}

// ── tests ────────────────────────────────────────────────────────────────────

describe("manifest ↔ dispatcher parity", () => {
  it("declares exactly the commands it dispatches", () => {
    const manifest = manifestCommands().sort();
    const dispatched = dispatcherCommands().sort();
    expect(manifest).toEqual(dispatched);
  });
});

describe("dispatcher ↔ real ctx command tree", () => {
  const ctxBin = resolveCtx();

  if (!ctxBin) {
    it.skip("SKIPPED: no ctx binary (set CTX_BIN) — ground-truth check did NOT run", () => {
      /* Loud skip: CI must provide CTX_BIN so this cannot silently pass. */
    });
    console.warn(
      "[commandParity] No ctx binary found (CTX_BIN unset, `ctx` not on PATH). " +
        "CLI ground-truth parity was SKIPPED — wire CTX_BIN in CI to enforce it."
    );
    return;
  }

  it("every dispatched ctx invocation is a real command path", () => {
    const dead: string[] = [];
    for (const path of invocationPaths()) {
      const missing = firstDeadSegment(ctxBin, path);
      if (missing) dead.push(`ctx ${path.join(" ")}  →  no such command: "${missing}"`);
    }
    expect(dead).toEqual([]);
  });
});
