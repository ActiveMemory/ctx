import * as vscode from "vscode";
import { execFile } from "child_process";
import * as fs from "fs";
import * as os from "os";
import * as path from "path";
import * as https from "https";

// Canonical ctx skills, bundled as text at build time (esbuild's
// `--loader:.md=text`; vitest.config.ts mirrors it). Bundling pins each
// skill body to the ctx commit the extension is built from, and a renamed
// or deleted skill fails the build instead of shipping a dead command.
import blogSkill from "../../../internal/assets/claude/skills/ctx-blog/SKILL.md";
import brainstormSkill from "../../../internal/assets/claude/skills/ctx-brainstorm/SKILL.md";
import consolidateSkill from "../../../internal/assets/claude/skills/ctx-consolidate/SKILL.md";
import implementSkill from "../../../internal/assets/claude/skills/ctx-implement/SKILL.md";
import nextSkill from "../../../internal/assets/claude/skills/ctx-next/SKILL.md";
import reflectSkill from "../../../internal/assets/claude/skills/ctx-reflect/SKILL.md";
import rememberSkill from "../../../internal/assets/claude/skills/ctx-remember/SKILL.md";
import specSkill from "../../../internal/assets/claude/skills/ctx-spec/SKILL.md";
import wrapUpSkill from "../../../internal/assets/claude/skills/ctx-wrap-up/SKILL.md";

const PARTICIPANT_ID = "ctx.participant";
const GITHUB_REPO = "ActiveMemory/ctx";

interface CtxResult extends vscode.ChatResult {
  metadata: {
    command: string;
  };
}

/** A CLI-backed slash command. `prompt` is the text after the command. */
type Handler = (
  stream: vscode.ChatResponseStream,
  prompt: string,
  cwd: string,
  token: vscode.CancellationToken
) => Promise<CtxResult>;

// Resolved path to ctx binary — set during bootstrap
let resolvedCtxPath: string | undefined;

// Extension context — set during activation
let extensionCtx: vscode.ExtensionContext | undefined;

// Status bar item for context reminders
let reminderStatusBar: vscode.StatusBarItem | undefined;

function getCtxPath(): string {
  if (resolvedCtxPath) {
    return resolvedCtxPath;
  }
  return (
    vscode.workspace.getConfiguration("ctx").get<string>("executablePath") ||
    "ctx"
  );
}

/**
 * The project root ctx runs in: the workspace folder of the active editor,
 * so a multi-root window targets the project being worked on, falling back
 * to the first folder.
 */
function getWorkspaceRoot(): string | undefined {
  const active = vscode.window.activeTextEditor?.document.uri;
  const folder =
    (active && vscode.workspace.getWorkspaceFolder(active)) ||
    vscode.workspace.workspaceFolders?.[0];
  return folder?.uri.fsPath;
}

/**
 * Map Node.js os values to Go GOOS/GOARCH used in release binary names.
 */
function getPlatformInfo(): { goos: string; goarch: string; ext: string } {
  const platform = os.platform();
  const arch = os.arch();

  let goos: string;
  switch (platform) {
    case "darwin":
      goos = "darwin";
      break;
    case "win32":
      goos = "windows";
      break;
    default:
      goos = "linux";
      break;
  }

  let goarch: string;
  switch (arch) {
    case "arm64":
    case "aarch64":
      goarch = "arm64";
      break;
    default:
      goarch = "amd64";
      break;
  }

  const ext = goos === "windows" ? ".exe" : "";
  return { goos, goarch, ext };
}

/**
 * Fetch JSON from a URL (follows redirects).
 */
function fetchJSON(url: string): Promise<unknown> {
  return new Promise((resolve, reject) => {
    const get = (reqUrl: string, redirectCount: number) => {
      if (redirectCount > 5) {
        reject(new Error("Too many redirects"));
        return;
      }
      https
        .get(reqUrl, { headers: { "User-Agent": "ctx-vscode" } }, (res) => {
          if (
            res.statusCode &&
            res.statusCode >= 300 &&
            res.statusCode < 400 &&
            res.headers.location
          ) {
            get(res.headers.location, redirectCount + 1);
            return;
          }
          if (res.statusCode !== 200) {
            reject(new Error(`HTTP ${res.statusCode} fetching ${reqUrl}`));
            return;
          }
          const chunks: Buffer[] = [];
          res.on("data", (chunk: Buffer) => chunks.push(chunk));
          res.on("end", () => {
            try {
              resolve(JSON.parse(Buffer.concat(chunks).toString()));
            } catch (e) {
              reject(e);
            }
          });
          res.on("error", reject);
        })
        .on("error", reject);
    };
    get(url, 0);
  });
}

/**
 * Download a file from a URL to a local path (follows redirects).
 */
function downloadFile(url: string, destPath: string): Promise<void> {
  return new Promise((resolve, reject) => {
    const get = (reqUrl: string, redirectCount: number) => {
      if (redirectCount > 5) {
        reject(new Error("Too many redirects"));
        return;
      }
      https
        .get(reqUrl, { headers: { "User-Agent": "ctx-vscode" } }, (res) => {
          if (
            res.statusCode &&
            res.statusCode >= 300 &&
            res.statusCode < 400 &&
            res.headers.location
          ) {
            get(res.headers.location, redirectCount + 1);
            return;
          }
          if (res.statusCode !== 200) {
            reject(new Error(`HTTP ${res.statusCode} downloading ${reqUrl}`));
            return;
          }
          const file = fs.createWriteStream(destPath);
          res.pipe(file);
          file.on("finish", () => {
            file.close();
            resolve();
          });
          file.on("error", (err) => {
            fs.unlink(destPath, () => {});
            reject(err);
          });
        })
        .on("error", (err) => {
          fs.unlink(destPath, () => {});
          reject(err);
        });
    };
    get(url, 0);
  });
}

/**
 * Check if a binary is executable by attempting to run it.
 */
function isCtxExecutable(binPath: string): Promise<boolean> {
  return new Promise((resolve) => {
    execFile(binPath, ["--version"], { timeout: 5000 }, (error) => {
      resolve(!error);
    });
  });
}

/**
 * Ensure the ctx CLI binary is available. If not found on PATH or at the
 * configured path, automatically downloads the correct platform binary
 * from GitHub releases into the extension's global storage directory.
 */
async function ensureCtxAvailable(): Promise<void> {
  // 1. Check if user-configured or PATH-resolved ctx works
  const configuredPath = getCtxPath();
  if (await isCtxExecutable(configuredPath)) {
    resolvedCtxPath = configuredPath;
    return;
  }

  // 2. Check if we already downloaded it to global storage
  if (extensionCtx) {
    const { ext } = getPlatformInfo();
    const storagePath = extensionCtx.globalStorageUri.fsPath;
    const localBin = path.join(storagePath, `ctx${ext}`);
    if (fs.existsSync(localBin) && (await isCtxExecutable(localBin))) {
      resolvedCtxPath = localBin;
      return;
    }
  }

  // 3. Download from GitHub releases
  if (!extensionCtx) {
    throw new Error(
      "ctx binary not found and extension context unavailable for auto-install."
    );
  }

  const { goos, goarch, ext } = getPlatformInfo();
  const storagePath = extensionCtx.globalStorageUri.fsPath;
  fs.mkdirSync(storagePath, { recursive: true });

  // Fetch latest release info from GitHub API
  const apiUrl = `https://api.github.com/repos/${GITHUB_REPO}/releases/latest`;
  const release = (await fetchJSON(apiUrl)) as {
    tag_name: string;
    assets: Array<{ name: string; browser_download_url: string }>;
  };

  const version = release.tag_name.replace(/^v/, "");
  const expectedName = `ctx-${version}-${goos}-${goarch}${ext}`;
  const asset = release.assets.find((a) => a.name === expectedName);

  if (!asset) {
    throw new Error(
      `No release binary found for ${goos}/${goarch} (looked for ${expectedName}). ` +
        `Install ctx manually: https://github.com/${GITHUB_REPO}/releases`
    );
  }

  const localBin = path.join(storagePath, `ctx${ext}`);
  await downloadFile(asset.browser_download_url, localBin);

  // Make executable on Unix
  if (goos !== "windows") {
    fs.chmodSync(localBin, 0o755);
  }

  // Verify the downloaded binary works
  if (!(await isCtxExecutable(localBin))) {
    fs.unlinkSync(localBin);
    throw new Error(
      "Downloaded ctx binary failed verification. " +
        `Install ctx manually: https://github.com/${GITHUB_REPO}/releases`
    );
  }

  resolvedCtxPath = localBin;
}

// Bootstrap state — ensures we only download once per session
let bootstrapPromise: Promise<void> | undefined;
let bootstrapDone = false;

async function bootstrap(): Promise<void> {
  if (bootstrapDone) {
    return;
  }
  if (!bootstrapPromise) {
    bootstrapPromise = ensureCtxAvailable().then(
      () => {
        bootstrapDone = true;
      },
      (err) => {
        // Reset so next attempt can retry
        bootstrapPromise = undefined;
        throw err;
      }
    );
  }
  return bootstrapPromise;
}

/**
 * Merge stdout and stderr without duplicating lines that appear in both.
 * Cobra prints errors to both streams — naive concatenation doubles them.
 */
function mergeOutput(stdout: string, stderr: string): string {
  const out = stdout.trim();
  const err = stderr.trim();
  if (!out) return err;
  if (!err) return out;
  // If stderr content already appears in stdout, skip it
  if (out.includes(err)) return out;
  if (err.includes(out)) return err;
  return out + "\n" + err;
}

/**
 * Result of a completed `ctx` process. `code` is the exit code: callers
 * must check it, because a failing command still prints output (an unknown
 * flag prints usage and exits 1).
 */
interface CtxRun {
  stdout: string;
  stderr: string;
  code: number;
}

/**
 * Run ctx without a shell: arguments reach the binary verbatim, so free
 * text from the chat prompt can neither split into extra arguments nor
 * inject shell syntax. (Node resolves `ctx` to `ctx.exe` on PATH on
 * Windows without a shell.) stdin is closed at once, so a command that
 * would prompt (`ctx why` with no document, an add with no content) fails
 * fast with its own message instead of waiting for the timeout.
 *
 * Resolves for any exit code; rejects only when there is no complete
 * result to show: the process could not start, its output overflowed the
 * buffer, or it was cancelled, timed out, or killed.
 */
function runCtx(
  args: string[],
  cwd?: string,
  token?: vscode.CancellationToken
): Promise<CtxRun> {
  const ctxPath = getCtxPath();
  return new Promise((resolve, reject) => {
    if (token?.isCancellationRequested) {
      reject(new Error("Cancelled"));
      return;
    }
    let disposed = false;
    // `disposable` is declared, not const-initialized: the cancellation
    // listener can only register after `child` exists, and mocked
    // execFile (vitest) fires the callback synchronously, so a
    // const-declared-later pattern would TDZ-trap.
    // eslint-disable-next-line prefer-const
    let disposable: { dispose(): void } | undefined;
    const child = execFile(
      ctxPath,
      args,
      { cwd, maxBuffer: 1024 * 1024, timeout: 30000 },
      (error, stdout, stderr) => {
        if (!disposed) {
          disposed = true;
          disposable?.dispose();
        }
        if (!error) {
          resolve({ stdout, stderr, code: 0 });
          return;
        }
        const err = error as NodeJS.ErrnoException & {
          killed?: boolean;
          signal?: NodeJS.Signals | null;
          code?: number | string;
        };
        // Killed by cancellation, the 30s timeout, or a signal: the
        // output is partial and must not be presented as a result.
        if (err.killed || err.signal || token?.isCancellationRequested) {
          reject(
            new Error(
              `\`ctx ${args.join(" ")}\` was cancelled or timed out` +
                (err.signal ? ` (${err.signal})` : "")
            )
          );
          return;
        }
        // A string code is a spawn or buffer failure (ENOENT,
        // ERR_CHILD_PROCESS_STDIO_MAXBUFFER), not an exit status.
        if (typeof err.code !== "number") {
          reject(error);
          return;
        }
        resolve({ stdout, stderr, code: err.code });
      }
    );
    child.stdin?.end();
    disposable = token?.onCancellationRequested(() => {
      child.kill();
    });
  });
}

/**
 * Run `git` with a timeout and cancellation, mirroring runCtx. Resolves
 * with stdout; rejects on non-zero exit, timeout, or cancel.
 */
function execGit(
  args: string[],
  cwd: string,
  token?: vscode.CancellationToken
): Promise<string> {
  return new Promise((resolve, reject) => {
    if (token?.isCancellationRequested) {
      reject(new Error("Cancelled"));
      return;
    }
    let disposed = false;
    // eslint-disable-next-line prefer-const
    let disposable: { dispose(): void } | undefined;
    const child = execFile(
      "git",
      args,
      { cwd, timeout: 30000, maxBuffer: 1024 * 1024 },
      (error, stdout) => {
        if (!disposed) {
          disposed = true;
          disposable?.dispose();
        }
        if (error) {
          reject(error);
          return;
        }
        resolve(stdout);
      }
    );
    disposable = token?.onCancellationRequested(() => child.kill());
  });
}

/**
 * Check if .context/ directory exists in the workspace root.
 */
function hasContextDir(cwd: string): boolean {
  return fs.existsSync(path.join(cwd, ".context"));
}

function fence(text: string): string {
  return "```\n" + text + "\n```";
}

function errorMarkdown(title: string, err: unknown): string {
  return `**Error:** ${title}.\n\n` + fence(err instanceof Error ? err.message : String(err));
}

function result(command: string): CtxResult {
  return { metadata: { command } };
}

/**
 * Run ctx and render the outcome. A non-zero exit is shown as such, with
 * the CLI's own output, never as a normal result; when the folder has no
 * .context/ yet, it also points at /init. Returns whether ctx exited 0.
 */
async function runAndRender(
  stream: vscode.ChatResponseStream,
  cwd: string,
  token: vscode.CancellationToken,
  args: string[],
  progress: string,
  emptyMessage: string,
  fenced = true
): Promise<boolean> {
  stream.progress(progress);
  let run: CtxRun;
  try {
    run = await runCtx(args, cwd, token);
  } catch (err: unknown) {
    stream.markdown(errorMarkdown(`\`ctx ${args.join(" ")}\` did not complete`, err));
    return false;
  }
  const output = mergeOutput(run.stdout, run.stderr);
  if (run.code !== 0) {
    stream.markdown(
      `**\`ctx ${args.join(" ")}\` exited with code ${run.code}.**\n\n` +
        fence(output || "(no output)") +
        (hasContextDir(cwd)
          ? ""
          : "\n\nThis folder has no `.context/` yet. Run `@ctx /init` to set it up.")
    );
    return false;
  }
  stream.markdown(!output ? emptyMessage : fenced ? fence(output) : output);
  return true;
}

/** A command that runs one fixed ctx invocation and ignores its prompt. */
function simple(
  command: string,
  args: string[],
  progress: string,
  emptyMessage: string,
  fenced = true
): Handler {
  return async (stream, _prompt, cwd, token) => {
    await runAndRender(stream, cwd, token, args, progress, emptyMessage, fenced);
    return result(command);
  };
}

/** A command whose first word selects one of `allowed` ctx subcommands. */
function subcommands(
  command: string,
  allowed: string[],
  usage: string
): Handler {
  return async (stream, prompt, cwd, token) => {
    const sub = prompt.trim().split(/\s+/)[0]?.toLowerCase();
    if (!sub || !allowed.includes(sub)) {
      stream.markdown(usage);
      return result(command);
    }
    const args = [command, sub];
    await runAndRender(stream, cwd, token, args, `Running ctx ${args.join(" ")}...`, `\`ctx ${args.join(" ")}\` completed.`);
    return result(command);
  };
}

/**
 * Split a prompt into words, keeping "double-quoted phrases" together so
 * flag values like `--context "why we chose it"` survive as one argument.
 */
function tokenize(prompt: string): string[] {
  return [...prompt.matchAll(/"([^"]*)"|(\S+)/g)].map((m) => m[1] ?? m[2]);
}

/** Leading words joined as free text, then everything from the first `--flag` on. */
function splitFlags(words: string[]): { text: string; flags: string[] } {
  const at = words.findIndex((w) => w.startsWith("--"));
  return at < 0
    ? { text: words.join(" "), flags: [] }
    : { text: words.slice(0, at).join(" "), flags: words.slice(at) };
}

const SESSION_START = ["system", "session-event", "--type", "start", "--caller", "vscode"];
const SESSION_END = ["system", "session-event", "--type", "end", "--caller", "vscode"];
const REMINDER_CHECK = ["remind", "list"];

/** ctx invocations the extension makes outside chat requests. */
const BACKGROUND_INVOCATIONS = [SESSION_START, SESSION_END, REMINDER_CHECK];

async function handleInit(
  stream: vscode.ChatResponseStream,
  _prompt: string,
  cwd: string,
  token: vscode.CancellationToken
): Promise<CtxResult> {
  const ok = await runAndRender(
    stream,
    cwd,
    token,
    ["init", "--caller", "vscode"],
    "Initializing .context/ directory...",
    "`.context/` initialized. Run `@ctx /status` to see your project context."
  );
  if (ok) {
    // Copilot reads .github/copilot-instructions.md natively, so plain
    // Copilot Chat picks up the project context too.
    await runAndRender(
      stream,
      cwd,
      token,
      ["setup", "copilot", "--write"],
      "Generating Copilot instructions...",
      "`.github/copilot-instructions.md` generated for Copilot context loading."
    );
    // activate() skipped session-start: .context/ did not exist yet.
    runCtx(SESSION_START, cwd).catch(() => {});
  }
  return result("init");
}

async function handleAgent(
  stream: vscode.ChatResponseStream,
  prompt: string,
  cwd: string,
  token: vscode.CancellationToken
): Promise<CtxResult> {
  const args = ["agent"];
  const budget = prompt.match(/(?:--budget\s+|budget\s+)(\d+)/);
  if (budget) {
    args.push("--budget", budget[1]);
  }
  await runAndRender(stream, cwd, token, args, "Generating AI-ready context packet...", "Empty context packet.", false);
  return result("agent");
}

async function handleRecall(
  stream: vscode.ChatResponseStream,
  prompt: string,
  cwd: string,
  token: vscode.CancellationToken
): Promise<CtxResult> {
  const words = prompt.trim().split(/\s+/).filter(Boolean);
  let args: string[];
  if (words[0]?.toLowerCase() === "show") {
    const id = words.slice(1).join(" ");
    if (!id) {
      stream.markdown("**Usage:** `@ctx /recall show <session-id>`");
      return result("recall");
    }
    args = ["journal", "source", "--show", id];
  } else {
    args = ["journal", "source"];
    const limit = prompt.match(/(?:--limit\s+|limit\s+)(\d+)/);
    if (limit) {
      args.push("--limit", limit[1]);
    }
  }
  await runAndRender(stream, cwd, token, args, "Loading session history...", "No session history found.");
  return result("recall");
}

async function handleSetup(
  stream: vscode.ChatResponseStream,
  prompt: string,
  cwd: string,
  token: vscode.CancellationToken
): Promise<CtxResult> {
  const words = prompt.trim().split(/\s+/).filter(Boolean);
  const preview = words.includes("preview") || words.includes("--preview");
  const tool = words.find((w) => w !== "preview" && w !== "--preview") || "copilot";
  const args = ["setup", tool];
  if (!preview) {
    args.push("--write");
  }
  await runAndRender(
    stream,
    cwd,
    token,
    args,
    preview ? `Previewing ${tool} integration config...` : `Generating ${tool} integration config...`,
    preview ? `No output for **${tool}** preview.` : `Integration config for **${tool}** generated.`
  );
  return result("setup");
}

const ENTRY_TYPES = ["task", "decision", "learning", "convention"];

/**
 * Provenance flags for `ctx task|decision|learning add`, which the CLI
 * requires. VS Code exposes no AI session ID, so the window session ID
 * stands in for it.
 */
async function provenanceFlags(
  cwd: string,
  token: vscode.CancellationToken
): Promise<string[]> {
  const git = (...args: string[]) =>
    execGit(args, cwd, token).then(
      (out) => out.trim() || "unknown",
      () => "unknown"
    );
  const [branch, commit] = await Promise.all([
    git("rev-parse", "--abbrev-ref", "HEAD"),
    git("rev-parse", "--short", "HEAD"),
  ]);
  return [
    "--session-id",
    vscode.env.sessionId.slice(0, 8),
    "--branch",
    branch,
    "--commit",
    commit,
  ];
}

async function handleAdd(
  stream: vscode.ChatResponseStream,
  prompt: string,
  cwd: string,
  token: vscode.CancellationToken
): Promise<CtxResult> {
  const [first, ...words] = tokenize(prompt);
  const type = first?.toLowerCase();
  if (!type || !ENTRY_TYPES.includes(type)) {
    stream.markdown(
      "**Usage:** `@ctx /add <type> <text> [flags]`\n\n" +
        "| Type | Required flags |\n" +
        "|------|----------------|\n" +
        '| `task` | `--section "<phase>"` |\n' +
        '| `decision` | `--context "..." --rationale "..." --consequence "..."` |\n' +
        '| `learning` | `--context "..." --lesson "..." --application "..."` |\n' +
        '| `convention` | `--section "<section>"` |\n\n' +
        "Quote multi-word values. Provenance (`--session-id`, `--branch`, " +
        "`--commit`) is filled in automatically.\n\n" +
        'Example: `@ctx /add task Fix the login redirect loop --section "Phase 1"`'
    );
    return result("add");
  }
  const { text, flags } = splitFlags(words);
  const args = [type, "add"];
  if (text) {
    args.push(text);
  }
  // No --section default: the CLI makes the caller choose the section,
  // so a catch-all never quietly collects every entry.
  args.push(...flags);
  if (type !== "convention") {
    const provenance = await provenanceFlags(cwd, token);
    for (let i = 0; i < provenance.length; i += 2) {
      if (!flags.includes(provenance[i])) {
        args.push(provenance[i], provenance[i + 1]);
      }
    }
  }
  await runAndRender(stream, cwd, token, args, `Adding ${type}...`, `Added **${type}**.`);
  return result("add");
}

async function handleTask(
  stream: vscode.ChatResponseStream,
  prompt: string,
  cwd: string,
  token: vscode.CancellationToken
): Promise<CtxResult> {
  const parts = prompt.trim().split(/\s+/);
  const subcmd = parts[0]?.toLowerCase();
  const rest = parts.slice(1).join(" ");

  let args: string[];
  switch (subcmd) {
    case "complete":
      if (!rest) {
        stream.markdown(
          "**Usage:** `@ctx /task complete <task-id-or-text>`\n\n" +
            "Example: `@ctx /task complete 3` or `@ctx /task complete Fix login bug`"
        );
        return result("task");
      }
      args = ["task", "complete", rest];
      break;
    case "archive":
      args = ["task", "archive"];
      break;
    case "snapshot":
      args = rest ? ["task", "snapshot", rest] : ["task", "snapshot"];
      break;
    default:
      stream.markdown(
        "**Usage:** `@ctx /task <subcommand>`\n\n" +
          "| Subcommand | Description |\n" +
          "|------------|-------------|\n" +
          "| `complete <ref>` | Mark a task as completed |\n" +
          "| `archive` | Move completed tasks to archive |\n" +
          "| `snapshot [name]` | Create point-in-time snapshot |\n\n" +
          "Add tasks with `@ctx /add task ...`."
      );
      return result("task");
  }
  await runAndRender(stream, cwd, token, args, `Running ctx ${args.join(" ")}...`, `\`ctx ${args.join(" ")}\` completed.`);
  return result("task");
}

async function handleRemind(
  stream: vscode.ChatResponseStream,
  prompt: string,
  cwd: string,
  token: vscode.CancellationToken
): Promise<CtxResult> {
  const parts = prompt.trim().split(/\s+/);
  const subcmd = parts[0]?.toLowerCase();
  const rest = parts.slice(1).join(" ");

  let args: string[];
  switch (subcmd) {
    case "dismiss":
    case "rm":
      args = rest ? ["remind", "dismiss", ...rest.split(/\s+/)] : ["remind", "dismiss", "--all"];
      break;
    case "list":
    case "ls":
      args = ["remind", "list"];
      break;
    case "add":
      args = rest ? ["remind", "add", rest] : ["remind", "list"];
      break;
    default:
      // Text without a subcommand is a new reminder.
      args = subcmd ? ["remind", "add", prompt.trim()] : ["remind", "list"];
      break;
  }
  await runAndRender(stream, cwd, token, args, "Managing reminders...", "No reminders.");
  return result("remind");
}

async function handlePad(
  stream: vscode.ChatResponseStream,
  prompt: string,
  cwd: string,
  token: vscode.CancellationToken
): Promise<CtxResult> {
  const parts = prompt.trim().split(/\s+/);
  const subcmd = parts[0]?.toLowerCase();
  const rest = parts.slice(1).join(" ");
  const usage: Record<string, string> = {
    add: "`@ctx /pad add <text>`",
    rm: "`@ctx /pad rm <number> [number...]`",
    edit: "`@ctx /pad edit <number> [text]`",
    mv: "`@ctx /pad mv <from> <to>`",
    import: "`@ctx /pad import <file>`",
    merge: "`@ctx /pad merge <file> [file...]`",
  };

  let args: string[];
  switch (subcmd) {
    case "add":
      args = ["pad", "add", rest];
      break;
    case "show":
      args = rest ? ["pad", "show", rest] : ["pad"];
      break;
    case "rm":
    case "mv":
    case "merge":
      args = ["pad", subcmd, ...parts.slice(1)];
      break;
    case "edit": {
      // `ctx pad edit N [TEXT]`: the replacement text is one argument.
      const text = parts.slice(2).join(" ");
      args = text ? ["pad", "edit", parts[1], text] : ["pad", "edit", parts[1]];
      break;
    }
    case "import":
      args = ["pad", "import", rest];
      break;
    case "export":
      args = rest ? ["pad", "export", rest] : ["pad", "export"];
      break;
    case "resolve":
      args = ["pad", "resolve"];
      break;
    default:
      // No subcommand or unknown — list all entries
      args = ["pad"];
      break;
  }
  if (subcmd && Object.hasOwn(usage, subcmd) && !rest) {
    stream.markdown(`**Usage:** ${usage[subcmd]}`);
    return result("pad");
  }
  await runAndRender(stream, cwd, token, args, "Accessing scratchpad...", "Scratchpad is empty.");
  return result("pad");
}

async function handleNotify(
  stream: vscode.ChatResponseStream,
  prompt: string,
  cwd: string,
  token: vscode.CancellationToken
): Promise<CtxResult> {
  const words = tokenize(prompt);
  const subcmd = words[0]?.toLowerCase();

  if (subcmd === "setup") {
    // `ctx hook notify setup` prompts for the webhook URL. The URL is a
    // secret, so it is entered in a terminal, never in the chat history.
    stream.markdown(
      "Run `ctx hook notify setup` in a terminal. It prompts for the webhook " +
        "URL and stores it encrypted; then check it with `@ctx /notify test`."
    );
    return result("notify");
  }
  let args: string[];
  if (subcmd === "test") {
    args = ["hook", "notify", "test"];
  } else {
    const { text, flags } = splitFlags(words);
    if (!text) {
      stream.markdown(
        "**Usage:** `@ctx /notify <subcommand>`\n\n" +
          "| Subcommand | Description |\n" +
          "|------------|-------------|\n" +
          "| `setup` | How to configure the webhook URL |\n" +
          "| `test` | Send test notification |\n" +
          "| `<message> --event <name>` | Send notification |\n\n" +
          "Example: `@ctx /notify test` or `@ctx /notify build done --event build`"
      );
      return result("notify");
    }
    args = ["hook", "notify", text, ...flags];
  }
  await runAndRender(stream, cwd, token, args, "Sending notification...", "Notification sent.");
  return result("notify");
}

async function handleSystem(
  stream: vscode.ChatResponseStream,
  prompt: string,
  cwd: string,
  token: vscode.CancellationToken
): Promise<CtxResult> {
  const parts = prompt.trim().split(/\s+/);
  const subcmd = parts[0]?.toLowerCase();

  let args: string[];
  switch (subcmd) {
    case "resources":
      args = ["sysinfo"];
      break;
    case "bootstrap":
      args = ["system", "bootstrap"];
      break;
    case "stats":
      args = ["usage"];
      break;
    case "message": {
      const action = parts[1]?.toLowerCase();
      args = ["show", "edit", "reset"].includes(action ?? "")
        ? ["hook", "message", action as string, ...parts.slice(2)]
        : ["hook", "message", "list"];
      break;
    }
    default:
      stream.markdown(
        "**Usage:** `@ctx /system <subcommand>`\n\n" +
          "| Subcommand | Runs |\n" +
          "|------------|------|\n" +
          "| `resources` | `ctx sysinfo`: memory, swap, disk, load |\n" +
          "| `bootstrap` | `ctx system bootstrap`: context location for AI agents |\n" +
          "| `stats` | `ctx usage`: session token usage |\n" +
          "| `message [list]` | `ctx hook message list` |\n" +
          "| `message show\\|edit\\|reset <hook> <variant>` | `ctx hook message ...` |\n\n" +
          "Example: `@ctx /system resources`"
      );
      return result("system");
  }
  await runAndRender(stream, cwd, token, args, `Running ctx ${args.join(" ")}...`, "No output.");
  return result("system");
}

async function handleConfig(
  stream: vscode.ChatResponseStream,
  prompt: string,
  cwd: string,
  token: vscode.CancellationToken
): Promise<CtxResult> {
  const parts = prompt.trim().split(/\s+/);
  const subcmd = parts[0]?.toLowerCase();
  const profile = parts[1];

  let args: string[];
  if (subcmd === "switch" && profile) {
    args = ["config", "switch", profile];
  } else if (subcmd === "status" || subcmd === "schema") {
    args = ["config", subcmd];
  } else {
    stream.markdown(
      "**Usage:** `@ctx /config <subcommand>`\n\n" +
        "| Subcommand | Description |\n" +
        "|------------|-------------|\n" +
        "| `switch <dev\\|base>` | Switch the runtime config profile |\n" +
        "| `status` | Show the active profile |\n" +
        "| `schema` | Show the .ctxrc schema |\n\n" +
        "Example: `@ctx /config switch dev`"
    );
    return result("config");
  }
  await runAndRender(stream, cwd, token, args, `Running ctx ${args.join(" ")}...`, "No output.");
  return result("config");
}

async function handleWhy(
  stream: vscode.ChatResponseStream,
  prompt: string,
  cwd: string,
  token: vscode.CancellationToken
): Promise<CtxResult> {
  // Bare `ctx why` opens an interactive menu, so default to a document.
  const args = ["why", prompt.trim() || "manifesto"];
  await runAndRender(stream, cwd, token, args, "Loading philosophy...", "No philosophy content available.", false);
  return result("why");
}

async function handleChange(
  stream: vscode.ChatResponseStream,
  prompt: string,
  cwd: string,
  token: vscode.CancellationToken
): Promise<CtxResult> {
  const args = ["change"];
  const since = prompt.match(/--since\s+(\S+)/);
  if (since) {
    args.push("--since", since[1]);
  }
  await runAndRender(stream, cwd, token, args, "Checking what changed...", "No changes detected since last session.", false);
  return result("change");
}

async function handleGuide(
  stream: vscode.ChatResponseStream,
  prompt: string,
  cwd: string,
  token: vscode.CancellationToken
): Promise<CtxResult> {
  const args = ["guide"];
  if (prompt.includes("--skills")) {
    args.push("--skills");
  } else if (prompt.includes("--commands")) {
    args.push("--commands");
  }
  await runAndRender(stream, cwd, token, args, "Loading guide...", "No guide output.");
  return result("guide");
}

/** CLI-backed slash commands: each dispatches to real `ctx` subcommands. */
const CLI_COMMANDS: Record<string, Handler> = {
  init: handleInit,
  status: simple("status", ["status"], "Checking context status...", "No status output."),
  agent: handleAgent,
  drift: simple("drift", ["drift"], "Detecting context drift...", "No drift output."),
  recall: handleRecall,
  setup: handleSetup,
  add: handleAdd,
  // posix join: ctx accepts forward slashes on every OS, and the argv
  // stays identical across platforms (see ctx-cli-surface.json).
  decision: simple(
    "decision",
    ["index", path.posix.join(".context", "DECISIONS.md")],
    "Loading decisions...",
    "No decisions recorded yet. Add one with `@ctx /add decision ...`."
  ),
  learning: simple(
    "learning",
    ["index", path.posix.join(".context", "LEARNINGS.md")],
    "Loading learnings...",
    "No learnings recorded yet. Add one with `@ctx /add learning ...`."
  ),
  load: simple("load", ["load"], "Loading assembled context...", "No context loaded.", false),
  compact: simple("compact", ["compact"], "Compacting context...", "Context compacted."),
  sync: simple("sync", ["sync"], "Syncing context with codebase...", "Context synced with codebase."),
  task: handleTask,
  remind: handleRemind,
  pad: handlePad,
  notify: handleNotify,
  system: handleSystem,
  memory: subcommands(
    "memory",
    ["sync", "status", "diff", "import", "publish", "unpublish"],
    "**Usage:** `@ctx /memory <sync|status|diff|import|publish|unpublish>`\n\n" +
      "Bridges Claude Code auto memory into `.context/`. Example: `@ctx /memory status`"
  ),
  journal: subcommands(
    "journal",
    ["site", "obsidian"],
    "**Usage:** `@ctx /journal <site|obsidian>`\n\n" +
      "Exports the session journal. Browse sessions with `@ctx /recall`."
  ),
  doctor: simple("doctor", ["doctor"], "Running context health diagnostics...", "Context health check passed."),
  config: handleConfig,
  why: handleWhy,
  change: handleChange,
  guide: handleGuide,
  permission: subcommands(
    "permission",
    ["snapshot", "restore"],
    "**Usage:** `@ctx /permission <snapshot|restore>`\n\n" +
      "Saves or restores the Claude Code permission golden image."
  ),
  pause: simple("pause", ["hook", "pause"], "Pausing context hooks...", "Context hooks paused."),
  resume: simple("resume", ["hook", "resume"], "Resuming context hooks...", "Context hooks resumed."),
};

interface SkillCommand {
  /** Skill directory under internal/assets/claude/skills. */
  skill: string;
  /** The skill's SKILL.md, bundled at build time. */
  text: string;
  /** Read-only ctx invocations the skill's process relies on; `ctx agent` always runs. */
  reads: string[][];
}

/**
 * Skill-backed slash commands: `/<name>` runs the canonical `ctx-<name>`
 * skill through the chat model, grounded in live ctx output. Only skills
 * that can work from that output, the @ctx conversation, and #file
 * attachments are exposed: the model here cannot run commands or read the
 * workspace on its own.
 */
const SKILLS: Record<string, SkillCommand> = {
  blog: { skill: "ctx-blog", text: blogSkill, reads: [["journal", "source", "--limit", "10"]] },
  brainstorm: { skill: "ctx-brainstorm", text: brainstormSkill, reads: [] },
  consolidate: { skill: "ctx-consolidate", text: consolidateSkill, reads: [["drift", "--json"]] },
  implement: { skill: "ctx-implement", text: implementSkill, reads: [] },
  next: { skill: "ctx-next", text: nextSkill, reads: [["journal", "source", "--limit", "3"]] },
  reflect: { skill: "ctx-reflect", text: reflectSkill, reads: [] },
  remember: { skill: "ctx-remember", text: rememberSkill, reads: [["journal", "source", "--limit", "3"]] },
  spec: { skill: "ctx-spec", text: specSkill, reads: [] },
  "wrap-up": { skill: "ctx-wrap-up", text: wrapUpSkill, reads: [] },
};

function skillPreamble(skill: string): string {
  return (
    `You are running the ctx skill \`${skill}\` for the user inside the VS Code \`@ctx\` chat participant. ` +
    "Follow the skill below within the limits of this environment:\n" +
    "- You cannot run commands or read or edit files. The ctx CLI output and any files the user attached " +
    "are included in the next message. If the skill needs something that is missing, ask the user to " +
    "attach the file with #file or to paste the command output.\n" +
    "- Where the skill says to run a command or change a file, give the exact command or change for the " +
    "user to apply (`ctx` commands also exist as `@ctx` slash commands, e.g. `@ctx /add`). Never claim you " +
    "ran or changed anything.\n\n"
  );
}

/**
 * Earlier skill exchanges in this @ctx conversation, so multi-turn skills
 * keep their thread. CLI command turns are left out: /pad, /notify, and
 * /add can carry secrets, and their output must not reach the model.
 */
function historyMessages(context: vscode.ChatContext): vscode.LanguageModelChatMessage[] {
  const messages: vscode.LanguageModelChatMessage[] = [];
  context.history.forEach((turn, i) => {
    if (!(turn instanceof vscode.ChatResponseTurn)) {
      return;
    }
    const command = (turn.result as CtxResult).metadata?.command;
    if (!command || !Object.hasOwn(SKILLS, command)) {
      return;
    }
    const asked = context.history[i - 1];
    if (asked instanceof vscode.ChatRequestTurn) {
      messages.push(
        vscode.LanguageModelChatMessage.User((asked.command ? `/${asked.command} ` : "") + asked.prompt)
      );
    }
    const answer = turn.response
      .map((part) => (part instanceof vscode.ChatResponseMarkdownPart ? part.value.value : ""))
      .join("");
    if (answer) {
      messages.push(vscode.LanguageModelChatMessage.Assistant(answer));
    }
  });
  return messages;
}

async function handleSkill(
  command: string,
  request: vscode.ChatRequest,
  context: vscode.ChatContext,
  stream: vscode.ChatResponseStream,
  cwd: string,
  token: vscode.CancellationToken
): Promise<CtxResult> {
  const { skill, text, reads } = SKILLS[command];
  stream.progress(`Loading project context for ${skill}...`);
  try {
    const sections: string[] = [];
    for (const args of [["agent"], ...reads]) {
      const run = await runCtx(args, cwd, token);
      sections.push(
        run.code === 0
          ? `### \`ctx ${args.join(" ")}\`\n\n${run.stdout.trim()}`
          : `### \`ctx ${args.join(" ")}\` (exited with code ${run.code})\n\n${mergeOutput(run.stdout, run.stderr)}`
      );
    }
    // ponytail: attachments are inlined whole; cap them if large files start blowing the model's context.
    for (const ref of request.references) {
      if (ref.value instanceof vscode.Uri) {
        const content = new TextDecoder().decode(await vscode.workspace.fs.readFile(ref.value));
        sections.push(`### Attached: ${vscode.workspace.asRelativePath(ref.value)}\n\n${content}`);
      }
    }

    const model =
      request.model ?? (await vscode.lm.selectChatModels({ vendor: "copilot" }))[0];
    if (!model) {
      stream.markdown("**Error:** No chat model is available. Sign in to GitHub Copilot Chat and retry.");
      return result(command);
    }
    const messages = [
      vscode.LanguageModelChatMessage.User(skillPreamble(skill) + text),
      vscode.LanguageModelChatMessage.User("## Project context\n\n" + sections.join("\n\n")),
      ...historyMessages(context),
      vscode.LanguageModelChatMessage.User(request.prompt.trim() || `Run the ${skill} skill.`),
    ];
    stream.progress(`Running ${skill}...`);
    const response = await model.sendRequest(messages, {}, token);
    for await (const fragment of response.text) {
      stream.markdown(fragment);
    }
  } catch (err: unknown) {
    stream.markdown(errorMarkdown(`Failed to run ${skill}`, err));
  }
  return result(command);
}

/**
 * Natural-language routing for prompts without a slash command. Only
 * read-only targets: a keyword match must never mutate context. CLI
 * targets run with no arguments; skill targets get the prompt.
 */
const FREEFORM: Array<[string[], string]> = [
  [["remember", "last session", "what were we"], "remember"],
  [["what should i", "work on next", "next task"], "next"],
  [["wrap up", "wrap-up", "end of session"], "wrap-up"],
  [["reflect", "worth saving", "worth persisting"], "reflect"],
  [["brainstorm", "idea"], "brainstorm"],
  [["status"], "status"],
  [["drift"], "drift"],
  [["doctor", "health"], "doctor"],
  [["what changed", "since last session"], "change"],
  [["recall", "session history"], "recall"],
  [["guide", "cheat sheet", "getting started"], "guide"],
  [["philosophy", "manifesto"], "why"],
];

function helpMarkdown(): string {
  const commands: Array<{ name: string; description: string }> =
    extensionCtx?.extension.packageJSON?.contributes?.chatParticipants?.[0]?.commands ?? [];
  return (
    "## ctx: Persistent Context for AI\n\n" +
    "| Command | Description |\n" +
    "|---------|-------------|\n" +
    commands.map((c) => `| \`/${c.name}\` | ${c.description} |`).join("\n") +
    "\n\nExample: `@ctx /status` or `@ctx /add task Fix login bug`"
  );
}

async function dispatch(
  command: string,
  prompt: string,
  request: vscode.ChatRequest,
  context: vscode.ChatContext,
  stream: vscode.ChatResponseStream,
  cwd: string,
  token: vscode.CancellationToken
): Promise<CtxResult> {
  if (Object.hasOwn(SKILLS, command)) {
    return handleSkill(command, request, context, stream, cwd, token);
  }
  return CLI_COMMANDS[command](stream, prompt, cwd, token);
}

const handler: vscode.ChatRequestHandler = async (
  request: vscode.ChatRequest,
  context: vscode.ChatContext,
  stream: vscode.ChatResponseStream,
  token: vscode.CancellationToken
): Promise<CtxResult> => {
  const cwd = getWorkspaceRoot();
  if (!cwd) {
    stream.markdown(
      "**Error:** No workspace folder is open. Open a project folder first."
    );
    return result(request.command || "none");
  }

  // Auto-bootstrap: ensure ctx binary is available before any command
  try {
    stream.progress("Checking ctx installation...");
    await bootstrap();
  } catch (err: unknown) {
    stream.markdown(
      `**Error:** ctx CLI not found and auto-install failed.\n\n` +
        fence(err instanceof Error ? err.message : String(err)) +
        `\n\nInstall manually: \`go install github.com/ActiveMemory/ctx/cmd/ctx@latest\` ` +
        `or download from [GitHub Releases](https://github.com/${GITHUB_REPO}/releases).`
    );
    return result(request.command || "none");
  }

  // No init gate here: the CLI decides which commands need .context/
  // (AnnotationSkipInit), and runAndRender points at /init when one fails.
  const command = request.command;
  if (command && (Object.hasOwn(SKILLS, command) || Object.hasOwn(CLI_COMMANDS, command))) {
    return dispatch(command, request.prompt, request, context, stream, cwd, token);
  }

  // A plain reply right after a skill turn continues that skill, so
  // multi-turn workflows (e.g. /brainstorm) keep their thread.
  const last = context.history[context.history.length - 1];
  const previous =
    last instanceof vscode.ChatResponseTurn ? (last.result as CtxResult).metadata?.command : undefined;
  if (previous && Object.hasOwn(SKILLS, previous)) {
    return handleSkill(previous, request, context, stream, cwd, token);
  }

  const text = request.prompt.trim().toLowerCase();
  const hit = FREEFORM.find(([keywords]) => keywords.some((k) => text.includes(k)));
  if (hit) {
    return dispatch(hit[1], "", request, context, stream, cwd, token);
  }
  stream.markdown(helpMarkdown());
  return result("help");
};

/**
 * Follow-up suggestions per command. `prompt` is what the follow-up sends
 * as the command's arguments; `label` is what the user sees.
 */
const FOLLOWUPS: Record<string, vscode.ChatFollowup[]> = {
  init: [
    { label: "Show context status", prompt: "", command: "status" },
    { label: "What should I work on next?", prompt: "", command: "next" },
  ],
  status: [
    { label: "Detect context drift", prompt: "", command: "drift" },
    { label: "What should I work on next?", prompt: "", command: "next" },
  ],
  drift: [
    { label: "Sync context with codebase", prompt: "", command: "sync" },
    { label: "Run health check", prompt: "", command: "doctor" },
  ],
  doctor: [
    { label: "Show context status", prompt: "", command: "status" },
    { label: "Detect context drift", prompt: "", command: "drift" },
  ],
  task: [
    { label: "Show context status", prompt: "", command: "status" },
    { label: "Compact context", prompt: "", command: "compact" },
  ],
  remind: [{ label: "List reminders", prompt: "list", command: "remind" }],
  pad: [{ label: "List scratchpad", prompt: "", command: "pad" }],
  pause: [{ label: "Resume hooks", prompt: "", command: "resume" }],
  resume: [{ label: "Show context status", prompt: "", command: "status" }],
  remember: [{ label: "What should I work on next?", prompt: "", command: "next" }],
  next: [{ label: "Show context status", prompt: "", command: "status" }],
  reflect: [{ label: "Wrap up the session", prompt: "", command: "wrap-up" }],
  "wrap-up": [{ label: "Record an entry", prompt: "", command: "add" }],
  brainstorm: [{ label: "Turn this into a spec", prompt: "Turn the design above into a spec.", command: "spec" }],
  spec: [{ label: "Plan the implementation", prompt: "Plan the implementation of the spec above.", command: "implement" }],
  help: [
    { label: "Initialize project context", prompt: "", command: "init" },
    { label: "Show context status", prompt: "", command: "status" },
    { label: "Quick-reference guide", prompt: "", command: "guide" },
  ],
};

/**
 * Show `$(bell) ctx` while `ctx remind list` has entries. `remind list`
 * is read-only, so the .context/** watcher that calls this cannot loop.
 */
function updateReminderStatus(cwd: string): void {
  if (!bootstrapDone || !reminderStatusBar) {
    return;
  }
  const bar = reminderStatusBar;
  runCtx(REMINDER_CHECK, cwd)
    .then(({ stdout, code }) => {
      const trimmed = stdout.trim();
      if (code === 0 && trimmed && !/no reminders/i.test(trimmed)) {
        bar.text = "$(bell) ctx";
        bar.tooltip = trimmed;
        bar.show();
      } else {
        bar.hide();
      }
    })
    .catch(() => bar.hide());
}

export function activate(extensionContext: vscode.ExtensionContext) {
  // Store extension context for auto-bootstrap binary downloads
  extensionCtx = extensionContext;

  const participant = vscode.chat.createChatParticipant(PARTICIPANT_ID, handler);
  participant.iconPath = vscode.Uri.joinPath(extensionContext.extensionUri, "icon.png");
  participant.followupProvider = {
    provideFollowups(result: CtxResult) {
      return FOLLOWUPS[result.metadata.command] ?? [];
    },
  };
  extensionContext.subscriptions.push(participant);

  reminderStatusBar = vscode.window.createStatusBarItem(vscode.StatusBarAlignment.Right, 50);
  reminderStatusBar.name = "ctx Reminders";
  extensionContext.subscriptions.push(reminderStatusBar);

  const cwd = getWorkspaceRoot();
  if (!cwd) {
    bootstrap().catch(() => {});
    return;
  }
  const refresh = () => updateReminderStatus(cwd);
  const watcher = vscode.workspace.createFileSystemWatcher(
    new vscode.RelativePattern(cwd, ".context/**")
  );
  watcher.onDidChange(refresh);
  watcher.onDidCreate(refresh);
  watcher.onDidDelete(refresh);
  const interval = setInterval(refresh, 5 * 60 * 1000);
  extensionContext.subscriptions.push(watcher, { dispose: () => clearInterval(interval) });

  // Background bootstrap; errors surface when the user invokes a command.
  // Reminder and session-start calls wait for it, so they use the resolved
  // (possibly auto-downloaded) binary.
  bootstrap().then(
    () => {
      if (hasContextDir(cwd)) {
        refresh();
        runCtx(SESSION_START, cwd).catch(() => {});
      }
    },
    () => {}
  );
}

export function deactivate(): Thenable<void> | undefined {
  const cwd = getWorkspaceRoot();
  if (!cwd || !bootstrapDone || !hasContextDir(cwd)) {
    return undefined;
  }
  return runCtx(SESSION_END, cwd).then(
    () => undefined,
    () => undefined
  );
}

export {
  runCtx,
  getCtxPath,
  getWorkspaceRoot,
  ensureCtxAvailable,
  bootstrap,
  getPlatformInfo,
  handler,
  tokenize,
  CLI_COMMANDS,
  SKILLS,
  FREEFORM,
  FOLLOWUPS,
  BACKGROUND_INVOCATIONS,
};
