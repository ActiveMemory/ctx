import * as vscode from "vscode";
import { execFile } from "child_process";
import * as fs from "fs";
import * as os from "os";
import * as path from "path";
import * as https from "https";

const PARTICIPANT_ID = "ctx.participant";
const GITHUB_REPO = "ActiveMemory/ctx";

/**
 * Split a string into args, respecting double-quoted segments.
 */
function splitArgs(input: string): string[] {
  const args: string[] = [];
  const regex = /(?:[^\s"]+|"[^"]*")+/g;
  let match;
  while ((match = regex.exec(input)) !== null) {
    let arg = match[0];
    if (arg.startsWith('"') && arg.endsWith('"') && arg.length >= 2) {
      arg = arg.slice(1, -1);
    }
    args.push(arg);
  }
  return args;
}

interface CtxResult extends vscode.ChatResult {
  metadata: {
    command: string;
  };
}

// Resolved path to ctx binary — set during bootstrap
let resolvedCtxPath: string | undefined;

// Extension context — set during activation
let extensionCtx: vscode.ExtensionContext | undefined;

function getCtxPath(): string {
  if (resolvedCtxPath) {
    return resolvedCtxPath;
  }
  return (
    vscode.workspace.getConfiguration("ctx").get<string>("executablePath") ||
    "ctx"
  );
}

function getWorkspaceRoot(): string | undefined {
  return vscode.workspace.workspaceFolders?.[0]?.uri.fsPath;
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
    const useShell = os.platform() === "win32";
    execFile(binPath, ["--version"], { timeout: 5000, shell: useShell }, (error, stdout) => {
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
  } else {
  }
  return bootstrapPromise;
}

function runCtx(
  args: string[],
  cwd?: string,
  token?: vscode.CancellationToken
): Promise<{ stdout: string; stderr: string }> {
  const ctxPath = getCtxPath();
  return new Promise((resolve, reject) => {
    if (token?.isCancellationRequested) {
      reject(new Error("Cancelled"));
      return;
    }
    let disposed = false;
    let disposable: { dispose(): void } | undefined;
    // Use shell on Windows so execFile can resolve PATH executables
    // without requiring the .exe extension.
    const useShell = os.platform() === "win32";
    const child = execFile(
      ctxPath,
      args,
      { cwd, maxBuffer: 1024 * 1024, timeout: 30000, shell: useShell },
      (error, stdout, stderr) => {
        if (!disposed) {
          disposed = true;
          disposable?.dispose();
        }
        if (error) {
          // Still return output even on non-zero exit — ctx drift uses exit 1
          // for "drift detected" which is a valid result
          if (stdout || stderr) {
            resolve({ stdout, stderr });
            return;
          }
          reject(error);
          return;
        }
        resolve({ stdout, stderr });
      }
    );
    // Close stdin immediately so the child process never blocks waiting
    // for interactive input (e.g. y/n prompts that --force/--merge skip).
    child.stdin?.end();
    disposable = token?.onCancellationRequested(() => {
      child.kill();
    });
  });
}

async function handleInit(
  stream: vscode.ChatResponseStream,
  cwd: string,
  token: vscode.CancellationToken
): Promise<CtxResult> {
  stream.progress("Initializing .context/ directory...");
  try {
    const { stdout, stderr } = await runCtx(["init", "--force", "--merge", "--no-color", "--caller", "vscode"], cwd, token);
    const output = (stdout + stderr).trim();
    if (output) {
      stream.markdown("```\n" + output + "\n```");
    }

    // Auto-generate .github/copilot-instructions.md so Copilot gets
    // project context automatically.
    stream.progress("Generating Copilot instructions...");
    try {
      const hookResult = await runCtx(
        ["hook", "copilot", "--write", "--no-color"],
        cwd,
        token
      );
      const hookOutput = (hookResult.stdout + hookResult.stderr).trim();
      if (hookOutput) {
        stream.markdown(
          "\n**Copilot integration:**\n```\n" + hookOutput + "\n```"
        );
      } else {
        stream.markdown(
          "\n`.github/copilot-instructions.md` generated for Copilot context loading."
        );
      }
    } catch {
      // Non-fatal — init succeeded, hook is a bonus
      stream.markdown(
        "\n> **Note:** Could not generate `.github/copilot-instructions.md`. " +
          "Run `@ctx /hook copilot` manually."
      );
    }

    if (!output) {
      stream.markdown(
        "`.context/` directory initialized. Run `@ctx /status` to see your project context."
      );
    }
  } catch (err: unknown) {
    stream.markdown(
      `**Error:** Failed to initialize context.\n\n\`\`\`\n${err instanceof Error ? err.message : String(err)}\n\`\`\``
    );
  }
  return { metadata: { command: "init" } };
}

async function handleStatus(
  stream: vscode.ChatResponseStream,
  cwd: string,
  token: vscode.CancellationToken
): Promise<CtxResult> {
  stream.progress("Checking context status...");
  try {
    const { stdout, stderr } = await runCtx(["status", "--no-color"], cwd, token);
    const output = (stdout + stderr).trim();
    stream.markdown("```\n" + output + "\n```");
  } catch (err: unknown) {
    stream.markdown(
      `**Error:** Failed to get status.\n\n\`\`\`\n${err instanceof Error ? err.message : String(err)}\n\`\`\``
    );
  }
  return { metadata: { command: "status" } };
}

async function handleAgent(
  stream: vscode.ChatResponseStream,
  prompt: string,
  cwd: string,
  token: vscode.CancellationToken
): Promise<CtxResult> {
  const args = ["agent", "--no-color"];
  const parts = prompt.trim().split(/\s+/);
  for (let i = 0; i < parts.length; i++) {
    const p = parts[i]?.toLowerCase();
    if ((p === "--budget" || p === "budget") && parts[i + 1]) {
      args.push("--budget", parts[++i]);
    } else if ((p === "--format" || p === "format") && parts[i + 1]) {
      args.push("--format", parts[++i]);
    }
  }
  stream.progress("Generating AI-ready context packet...");
  try {
    const { stdout, stderr } = await runCtx(args, cwd, token);
    const output = (stdout + stderr).trim();
    stream.markdown(output);
  } catch (err: unknown) {
    stream.markdown(
      `**Error:** Failed to generate agent context.\n\n\`\`\`\n${err instanceof Error ? err.message : String(err)}\n\`\`\``
    );
  }
  return { metadata: { command: "agent" } };
}

async function handleDrift(
  stream: vscode.ChatResponseStream,
  cwd: string,
  token: vscode.CancellationToken
): Promise<CtxResult> {
  stream.progress("Detecting context drift...");
  try {
    const { stdout, stderr } = await runCtx(["drift", "--no-color"], cwd, token);
    const output = (stdout + stderr).trim();
    stream.markdown("```\n" + output + "\n```");
  } catch (err: unknown) {
    stream.markdown(
      `**Error:** Failed to detect drift.\n\n\`\`\`\n${err instanceof Error ? err.message : String(err)}\n\`\`\``
    );
  }
  return { metadata: { command: "drift" } };
}

async function handleRecall(
  stream: vscode.ChatResponseStream,
  prompt: string,
  cwd: string,
  token: vscode.CancellationToken
): Promise<CtxResult> {
  const parts = prompt.trim().split(/\s+/);
  const subcmd = parts[0]?.toLowerCase();
  const rest = parts.slice(1);

  let args: string[];
  let progressMsg: string;

  switch (subcmd) {
    case "show":
      args = rest.length ? ["recall", "show", ...rest] : ["recall", "show", "--latest"];
      progressMsg = "Showing session details...";
      break;
    case "export":
      args = ["recall", "export", "--all", ...rest];
      progressMsg = "Exporting sessions...";
      break;
    case "lock":
      args = rest.length ? ["recall", "lock", ...rest] : ["recall", "lock"];
      progressMsg = "Locking journal entries...";
      break;
    case "unlock":
      args = rest.length ? ["recall", "unlock", ...rest] : ["recall", "unlock", "--all"];
      progressMsg = "Unlocking journal entries...";
      break;
    case "sync":
      args = ["recall", "sync"];
      progressMsg = "Syncing lock state...";
      break;
    case "list":
      args = ["recall", "list", ...rest];
      progressMsg = "Listing sessions...";
      break;
    default:
      // No subcommand or unknown — default to list with optional query
      args = ["recall", "list"];
      if (subcmd) {
        args.push("--query", prompt.trim());
      }
      progressMsg = "Searching session history...";
      break;
  }
  args.push("--no-color");

  stream.progress(progressMsg);
  try {
    const { stdout, stderr } = await runCtx(args, cwd, token);
    const output = (stdout + stderr).trim();
    if (output) {
      stream.markdown("```\n" + output + "\n```");
    } else {
      stream.markdown("No session history found.");
    }
  } catch (err: unknown) {
    stream.markdown(
      `**Error:** Failed to recall sessions.\n\n\`\`\`\n${err instanceof Error ? err.message : String(err)}\n\`\`\``
    );
  }
  return { metadata: { command: "recall" } };
}

async function handleHook(
  stream: vscode.ChatResponseStream,
  prompt: string,
  cwd: string,
  token: vscode.CancellationToken
): Promise<CtxResult> {
  const parts = prompt.trim().split(/\s+/);
  const tool = parts[0] || "copilot";
  const preview = parts.includes("preview") || parts.includes("--preview");

  const args = ["hook", tool];
  if (!preview) {
    args.push("--write");
  }
  args.push("--no-color");

  stream.progress(
    preview
      ? `Previewing ${tool} integration config...`
      : `Generating ${tool} integration config...`
  );
  try {
    const { stdout, stderr } = await runCtx(args, cwd, token);
    const output = (stdout + stderr).trim();
    if (output) {
      stream.markdown("```\n" + output + "\n```");
    } else {
      stream.markdown(
        preview
          ? `No output for **${tool}** preview.`
          : `Integration config for **${tool}** generated.`
      );
    }
  } catch (err: unknown) {
    stream.markdown(
      `**Error:** Failed to generate hook.\n\n\`\`\`\n${err instanceof Error ? err.message : String(err)}\n\`\`\``
    );
  }
  return { metadata: { command: "hook" } };
}

async function handleAdd(
  stream: vscode.ChatResponseStream,
  prompt: string,
  cwd: string,
  token: vscode.CancellationToken
): Promise<CtxResult> {
  const parts = splitArgs(prompt.trim());
  const type = parts[0]?.toLowerCase();

  if (!type) {
    stream.markdown(
      "**Usage:** `@ctx /add <type> <content> [flags]`\n\n" +
        "Types: `task`, `decision`, `learning`, `convention`\n\n" +
        "**Task:** `@ctx /add task Implement auth --priority high`\n\n" +
        "**Decision:**\n```\n@ctx /add decision \"Use PostgreSQL\" --context \"Need reliable DB\" --rationale \"ACID + JSON\" --consequences \"Team training\"\n```\n\n" +
        "**Learning:**\n```\n@ctx /add learning \"Go embed trick\" --context \"Tried parent dir\" --lesson \"Same package only\" --application \"Put in internal/\"\n```"
    );
    return { metadata: { command: "add" } };
  }

  stream.progress(`Adding ${type}...`);
  try {
    const args = ["add", type, ...parts.slice(1), "--no-color"];
    const { stdout, stderr } = await runCtx(args, cwd, token);
    const output = (stdout + stderr).trim();
    if (output) {
      stream.markdown("```\n" + output + "\n```");
    } else {
      stream.markdown(`Added **${type}** entry.`);
    }
  } catch (err: unknown) {
    stream.markdown(
      `**Error:** Failed to add ${type}.\n\n\`\`\`\n${err instanceof Error ? err.message : String(err)}\n\`\`\``
    );
  }
  return { metadata: { command: "add" } };
}

async function handleLoad(
  stream: vscode.ChatResponseStream,
  prompt: string,
  cwd: string,
  token: vscode.CancellationToken
): Promise<CtxResult> {
  const args = ["load", "--no-color"];
  const parts = prompt.trim().split(/\s+/);
  for (let i = 0; i < parts.length; i++) {
    const p = parts[i]?.toLowerCase();
    if ((p === "--budget" || p === "budget") && parts[i + 1]) {
      args.push("--budget", parts[++i]);
    } else if (p === "--raw" || p === "raw") {
      args.push("--raw");
    }
  }
  stream.progress("Loading assembled context...");
  try {
    const { stdout, stderr } = await runCtx(args, cwd, token);
    const output = (stdout + stderr).trim();
    stream.markdown(output);
  } catch (err: unknown) {
    stream.markdown(
      `**Error:** Failed to load context.\n\n\`\`\`\n${err instanceof Error ? err.message : String(err)}\n\`\`\``
    );
  }
  return { metadata: { command: "load" } };
}

async function handleCompact(
  stream: vscode.ChatResponseStream,
  prompt: string,
  cwd: string,
  token: vscode.CancellationToken
): Promise<CtxResult> {
  const args = ["compact", "--no-color"];
  const lower = prompt.trim().toLowerCase();
  if (lower.includes("archive") || lower.includes("--archive")) {
    args.push("--archive");
  }
  stream.progress("Compacting context...");
  try {
    const { stdout, stderr } = await runCtx(args, cwd, token);
    const output = (stdout + stderr).trim();
    if (output) {
      stream.markdown("```\n" + output + "\n```");
    } else {
      stream.markdown("Context compacted successfully.");
    }
  } catch (err: unknown) {
    stream.markdown(
      `**Error:** Failed to compact context.\n\n\`\`\`\n${err instanceof Error ? err.message : String(err)}\n\`\`\``
    );
  }
  return { metadata: { command: "compact" } };
}

async function handleSync(
  stream: vscode.ChatResponseStream,
  prompt: string,
  cwd: string,
  token: vscode.CancellationToken
): Promise<CtxResult> {
  const args = ["sync", "--no-color"];
  const lower = prompt.trim().toLowerCase();
  if (lower.includes("dry-run") || lower.includes("--dry-run") || lower.includes("dry run")) {
    args.push("--dry-run");
  }
  stream.progress("Syncing context with codebase...");
  try {
    const { stdout, stderr } = await runCtx(args, cwd, token);
    const output = (stdout + stderr).trim();
    if (output) {
      stream.markdown("```\n" + output + "\n```");
    } else {
      stream.markdown("Context synced with codebase.");
    }
  } catch (err: unknown) {
    stream.markdown(
      `**Error:** Failed to sync context.\n\n\`\`\`\n${err instanceof Error ? err.message : String(err)}\n\`\`\``
    );
  }
  return { metadata: { command: "sync" } };
}

async function handleComplete(
  stream: vscode.ChatResponseStream,
  prompt: string,
  cwd: string,
  token: vscode.CancellationToken
): Promise<CtxResult> {
  const taskRef = prompt.trim();
  if (!taskRef) {
    stream.markdown(
      "**Usage:** `@ctx /complete <task-id-or-text>`\n\n" +
        "Example: `@ctx /complete 3` or `@ctx /complete Fix login bug`"
    );
    return { metadata: { command: "complete" } };
  }

  stream.progress("Marking task as completed...");
  try {
    const { stdout, stderr } = await runCtx(
      ["complete", taskRef, "--no-color"],
      cwd,
      token
    );
    const output = (stdout + stderr).trim();
    if (output) {
      stream.markdown("```\n" + output + "\n```");
    } else {
      stream.markdown(`Task **${taskRef}** marked as completed.`);
    }
  } catch (err: unknown) {
    stream.markdown(
      `**Error:** Failed to complete task.\n\n\`\`\`\n${err instanceof Error ? err.message : String(err)}\n\`\`\``
    );
  }
  return { metadata: { command: "complete" } };
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
  let progressMsg: string;

  switch (subcmd) {
    case "dismiss":
    case "rm":
      args = rest ? ["remind", "dismiss", rest] : ["remind", "dismiss", "--all"];
      progressMsg = "Dismissing reminder(s)...";
      break;
    case "list":
    case "ls":
      args = ["remind", "list"];
      progressMsg = "Listing reminders...";
      break;
    case "add":
      args = rest ? ["remind", "add", rest] : ["remind", "list"];
      progressMsg = rest ? "Adding reminder..." : "Listing reminders...";
      break;
    default:
      // If text provided without subcommand, treat as "add"
      if (subcmd) {
        args = ["remind", "add", prompt.trim()];
        progressMsg = "Adding reminder...";
      } else {
        args = ["remind", "list"];
        progressMsg = "Listing reminders...";
      }
      break;
  }
  args.push("--no-color");

  stream.progress(progressMsg);
  try {
    const { stdout, stderr } = await runCtx(args, cwd, token);
    const output = (stdout + stderr).trim();
    if (output) {
      stream.markdown("```\n" + output + "\n```");
    } else {
      stream.markdown("No reminders.");
    }
  } catch (err: unknown) {
    stream.markdown(
      `**Error:** Failed to manage reminders.\n\n\`\`\`\n${err instanceof Error ? err.message : String(err)}\n\`\`\``
    );
  }
  return { metadata: { command: "remind" } };
}

async function handleTasks(
  stream: vscode.ChatResponseStream,
  prompt: string,
  cwd: string,
  token: vscode.CancellationToken
): Promise<CtxResult> {
  const parts = prompt.trim().split(/\s+/);
  const subcmd = parts[0]?.toLowerCase();
  const rest = parts.slice(1).join(" ");

  let args: string[];
  let progressMsg: string;

  switch (subcmd) {
    case "archive":
      args = ["tasks", "archive"];
      progressMsg = "Archiving completed tasks...";
      break;
    case "snapshot":
      args = rest ? ["tasks", "snapshot", rest] : ["tasks", "snapshot"];
      progressMsg = "Creating task snapshot...";
      break;
    default:
      stream.markdown(
        "**Usage:** `@ctx /tasks <subcommand>`\n\n" +
          "| Subcommand | Description |\n" +
          "|------------|-------------|\n" +
          "| `archive` | Move completed tasks to archive |\n" +
          "| `snapshot [name]` | Create point-in-time snapshot |\n\n" +
          "Example: `@ctx /tasks archive` or `@ctx /tasks snapshot pre-refactor`"
      );
      return { metadata: { command: "tasks" } };
  }
  args.push("--no-color");

  stream.progress(progressMsg);
  try {
    const { stdout, stderr } = await runCtx(args, cwd, token);
    const output = (stdout + stderr).trim();
    if (output) {
      stream.markdown("```\n" + output + "\n```");
    } else {
      stream.markdown(
        subcmd === "archive"
          ? "Completed tasks archived."
          : "Task snapshot created."
      );
    }
  } catch (err: unknown) {
    stream.markdown(
      `**Error:** Failed to ${subcmd} tasks.\n\n\`\`\`\n${err instanceof Error ? err.message : String(err)}\n\`\`\``
    );
  }
  return { metadata: { command: "tasks" } };
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

  let args: string[];
  let progressMsg: string;

  switch (subcmd) {
    case "add":
      if (!rest) {
        stream.markdown("**Usage:** `@ctx /pad add <text>`");
        return { metadata: { command: "pad" } };
      }
      args = ["pad", "add", rest];
      progressMsg = "Adding scratchpad entry...";
      break;
    case "show":
      args = rest ? ["pad", "show", rest] : ["pad"];
      progressMsg = "Showing scratchpad entry...";
      break;
    case "rm":
      if (!rest) {
        stream.markdown("**Usage:** `@ctx /pad rm <number>`");
        return { metadata: { command: "pad" } };
      }
      args = ["pad", "rm", rest];
      progressMsg = "Removing scratchpad entry...";
      break;
    case "edit":
      if (!rest) {
        stream.markdown("**Usage:** `@ctx /pad edit <number> [text]`");
        return { metadata: { command: "pad" } };
      }
      args = ["pad", "edit", ...parts.slice(1)];
      progressMsg = "Editing scratchpad entry...";
      break;
    case "mv":
      args = ["pad", "mv", ...parts.slice(1)];
      progressMsg = "Moving scratchpad entry...";
      break;
    default:
      // No subcommand or unknown — list all entries
      args = ["pad"];
      progressMsg = "Listing scratchpad...";
      break;
  }
  args.push("--no-color");

  stream.progress(progressMsg);
  try {
    const { stdout, stderr } = await runCtx(args, cwd, token);
    const output = (stdout + stderr).trim();
    if (output) {
      stream.markdown("```\n" + output + "\n```");
    } else {
      stream.markdown("Scratchpad is empty.");
    }
  } catch (err: unknown) {
    stream.markdown(
      `**Error:** Failed to access scratchpad.\n\n\`\`\`\n${err instanceof Error ? err.message : String(err)}\n\`\`\``
    );
  }
  return { metadata: { command: "pad" } };
}

async function handleNotify(
  stream: vscode.ChatResponseStream,
  prompt: string,
  cwd: string,
  token: vscode.CancellationToken
): Promise<CtxResult> {
  const parts = prompt.trim().split(/\s+/);
  const subcmd = parts[0]?.toLowerCase();
  const rest = parts.slice(1).join(" ");

  let args: string[];
  let progressMsg: string;

  switch (subcmd) {
    case "setup":
      args = ["notify", "setup"];
      progressMsg = "Setting up webhook...";
      break;
    case "test":
      args = ["notify", "test"];
      progressMsg = "Sending test notification...";
      break;
    default: {
      // Send a notification — require --event flag
      if (!subcmd) {
        stream.markdown(
          "**Usage:** `@ctx /notify <subcommand>`\n\n" +
            "| Subcommand | Description |\n" +
            "|------------|-------------|\n" +
            "| `setup` | Configure webhook URL |\n" +
            "| `test` | Send test notification |\n" +
            "| `<message> --event <name>` | Send notification |\n\n" +
            "Example: `@ctx /notify test` or `@ctx /notify setup`"
        );
        return { metadata: { command: "notify" } };
      }
      args = ["notify", ...parts];
      progressMsg = "Sending notification...";
      break;
    }
  }
  args.push("--no-color");

  stream.progress(progressMsg);
  try {
    const { stdout, stderr } = await runCtx(args, cwd, token);
    const output = (stdout + stderr).trim();
    if (output) {
      stream.markdown("```\n" + output + "\n```");
    } else {
      stream.markdown(
        subcmd === "setup"
          ? "Webhook configured."
          : subcmd === "test"
            ? "Test notification sent."
            : "Notification sent."
      );
    }
  } catch (err: unknown) {
    stream.markdown(
      `**Error:** Failed to send notification.\n\n\`\`\`\n${err instanceof Error ? err.message : String(err)}\n\`\`\``
    );
  }
  return { metadata: { command: "notify" } };
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
  let progressMsg: string;

  switch (subcmd) {
    case "resources":
      args = ["system", "resources"];
      progressMsg = "Checking system resources...";
      break;
    case "bootstrap":
      args = ["system", "bootstrap"];
      progressMsg = "Running bootstrap...";
      break;
    case "message":
      args = ["system", "message", ...parts.slice(1)];
      progressMsg = "Managing hook messages...";
      break;
    default:
      stream.markdown(
        "**Usage:** `@ctx /system <subcommand>`\n\n" +
          "| Subcommand | Description |\n" +
          "|------------|-------------|\n" +
          "| `resources` | Show system resource usage |\n" +
          "| `bootstrap` | Print context location for AI agents |\n" +
          "| `message list\|show\|edit\|reset` | Manage hook messages |\n\n" +
          "Example: `@ctx /system resources` or `@ctx /system bootstrap`"
      );
      return { metadata: { command: "system" } };
  }
  args.push("--no-color");

  stream.progress(progressMsg);
  try {
    const { stdout, stderr } = await runCtx(args, cwd, token);
    const output = (stdout + stderr).trim();
    if (output) {
      stream.markdown("```\n" + output + "\n```");
    } else {
      stream.markdown("No output.");
    }
  } catch (err: unknown) {
    stream.markdown(
      `**Error:** System command failed.\n\n\`\`\`\n${err instanceof Error ? err.message : String(err)}\n\`\`\``
    );
  }
  return { metadata: { command: "system" } };
}

async function handleChanges(
  stream: vscode.ChatResponseStream,
  prompt: string,
  cwd: string,
  token: vscode.CancellationToken
): Promise<CtxResult> {
  const args = ["changes", "--no-color"];
  const since = prompt.trim();
  if (since) {
    args.push("--since", since);
  }
  stream.progress("Checking changes since last session...");
  try {
    const { stdout, stderr } = await runCtx(args, cwd, token);
    const output = (stdout + stderr).trim();
    if (output) {
      stream.markdown("```\n" + output + "\n```");
    } else {
      stream.markdown("No changes detected.");
    }
  } catch (err: unknown) {
    stream.markdown(
      `**Error:** Failed to check changes.\n\n\`\`\`\n${err instanceof Error ? err.message : String(err)}\n\`\`\``
    );
  }
  return { metadata: { command: "changes" } };
}

async function handleConfig(
  stream: vscode.ChatResponseStream,
  prompt: string,
  cwd: string,
  token: vscode.CancellationToken
): Promise<CtxResult> {
  const parts = prompt.trim().split(/\s+/);
  const subcmd = parts[0]?.toLowerCase();

  let args: string[];
  let progressMsg: string;

  switch (subcmd) {
    case "switch":
      args = parts[1] ? ["config", "switch", parts[1]] : ["config", "switch"];
      progressMsg = "Switching config profile...";
      break;
    case "status":
      args = ["config", "status"];
      progressMsg = "Checking config status...";
      break;
    case "schema":
      args = ["config", "schema"];
      progressMsg = "Printing config schema...";
      break;
    default:
      stream.markdown(
        "**Usage:** `@ctx /config <subcommand>`\n\n" +
          "| Subcommand | Description |\n" +
          "|------------|-------------|\n" +
          "| `switch [dev\\|base]` | Switch .ctxrc profile |\n" +
          "| `status` | Show active profile |\n" +
          "| `schema` | Print JSON Schema for .ctxrc |\n\n" +
          "Example: `@ctx /config status` or `@ctx /config switch dev`"
      );
      return { metadata: { command: "config" } };
  }
  args.push("--no-color");

  stream.progress(progressMsg);
  try {
    const { stdout, stderr } = await runCtx(args, cwd, token);
    const output = (stdout + stderr).trim();
    if (output) {
      stream.markdown("```\n" + output + "\n```");
    } else {
      stream.markdown("No output.");
    }
  } catch (err: unknown) {
    stream.markdown(
      `**Error:** Config command failed.\n\n\`\`\`\n${err instanceof Error ? err.message : String(err)}\n\`\`\``
    );
  }
  return { metadata: { command: "config" } };
}

async function handleDoctor(
  stream: vscode.ChatResponseStream,
  cwd: string,
  token: vscode.CancellationToken
): Promise<CtxResult> {
  stream.progress("Running health checks...");
  try {
    const { stdout, stderr } = await runCtx(["doctor", "--no-color"], cwd, token);
    const output = (stdout + stderr).trim();
    stream.markdown("```\n" + output + "\n```");
  } catch (err: unknown) {
    stream.markdown(
      `**Error:** Doctor check failed.\n\n\`\`\`\n${err instanceof Error ? err.message : String(err)}\n\`\`\``
    );
  }
  return { metadata: { command: "doctor" } };
}

async function handleGuide(
  stream: vscode.ChatResponseStream,
  prompt: string,
  cwd: string,
  token: vscode.CancellationToken
): Promise<CtxResult> {
  const args = ["guide", "--no-color"];
  const flag = prompt.trim().toLowerCase();
  if (flag === "skills" || flag === "--skills") {
    args.push("--skills");
  } else if (flag === "commands" || flag === "--commands") {
    args.push("--commands");
  }
  stream.progress("Loading guide...");
  try {
    const { stdout, stderr } = await runCtx(args, cwd, token);
    const output = (stdout + stderr).trim();
    stream.markdown(output);
  } catch (err: unknown) {
    stream.markdown(
      `**Error:** Failed to load guide.\n\n\`\`\`\n${err instanceof Error ? err.message : String(err)}\n\`\`\``
    );
  }
  return { metadata: { command: "guide" } };
}

async function handleWhy(
  stream: vscode.ChatResponseStream,
  prompt: string,
  cwd: string,
  token: vscode.CancellationToken
): Promise<CtxResult> {
  const doc = prompt.trim().toLowerCase();
  const args = ["why"];
  if (doc) {
    args.push(doc);
  }
  args.push("--no-color");
  stream.progress("Loading philosophy document...");
  try {
    const { stdout, stderr } = await runCtx(args, cwd, token);
    const output = (stdout + stderr).trim();
    stream.markdown(output);
  } catch (err: unknown) {
    stream.markdown(
      `**Error:** Failed to load document.\n\n\`\`\`\n${err instanceof Error ? err.message : String(err)}\n\`\`\``
    );
  }
  return { metadata: { command: "why" } };
}

async function handleMemory(
  stream: vscode.ChatResponseStream,
  prompt: string,
  cwd: string,
  token: vscode.CancellationToken
): Promise<CtxResult> {
  const parts = prompt.trim().split(/\s+/);
  const subcmd = parts[0]?.toLowerCase();

  let args: string[];
  let progressMsg: string;

  switch (subcmd) {
    case "sync":
      args = ["memory", "sync"];
      progressMsg = "Syncing memory...";
      break;
    case "status":
      args = ["memory", "status"];
      progressMsg = "Checking memory status...";
      break;
    case "diff":
      args = ["memory", "diff"];
      progressMsg = "Showing memory diff...";
      break;
    case "import":
      args = ["memory", "import"];
      progressMsg = "Importing memory entries...";
      break;
    case "publish":
      args = ["memory", "publish"];
      progressMsg = "Publishing to MEMORY.md...";
      break;
    case "unpublish":
      args = ["memory", "unpublish"];
      progressMsg = "Removing published block...";
      break;
    default:
      stream.markdown(
        "**Usage:** `@ctx /memory <subcommand>`\n\n" +
          "| Subcommand | Description |\n" +
          "|------------|-------------|\n" +
          "| `sync` | Copy MEMORY.md to mirror |\n" +
          "| `status` | Show drift and timestamps |\n" +
          "| `diff` | Show changes since last sync |\n" +
          "| `import` | Classify entries into .context/ files |\n" +
          "| `publish` | Push context to MEMORY.md |\n" +
          "| `unpublish` | Remove published block |\n\n" +
          "Example: `@ctx /memory status` or `@ctx /memory sync`"
      );
      return { metadata: { command: "memory" } };
  }
  args.push("--no-color");

  stream.progress(progressMsg);
  try {
    const { stdout, stderr } = await runCtx(args, cwd, token);
    const output = (stdout + stderr).trim();
    if (output) {
      stream.markdown("```\n" + output + "\n```");
    } else {
      stream.markdown("No output.");
    }
  } catch (err: unknown) {
    stream.markdown(
      `**Error:** Memory command failed.\n\n\`\`\`\n${err instanceof Error ? err.message : String(err)}\n\`\`\``
    );
  }
  return { metadata: { command: "memory" } };
}

async function handlePrompt(
  stream: vscode.ChatResponseStream,
  prompt: string,
  cwd: string,
  token: vscode.CancellationToken
): Promise<CtxResult> {
  const parts = prompt.trim().split(/\s+/);
  const subcmd = parts[0]?.toLowerCase();
  const rest = parts.slice(1).join(" ");

  let args: string[];
  let progressMsg: string;

  switch (subcmd) {
    case "list":
    case "ls":
      args = ["prompt", "list"];
      progressMsg = "Listing prompt templates...";
      break;
    case "show":
      args = rest ? ["prompt", "show", rest] : ["prompt", "list"];
      progressMsg = rest ? "Showing prompt template..." : "Listing prompt templates...";
      break;
    case "add":
      args = rest ? ["prompt", "add", rest] : ["prompt", "list"];
      progressMsg = rest ? "Creating prompt template..." : "Listing prompt templates...";
      break;
    case "rm":
      if (!rest) {
        stream.markdown("**Usage:** `@ctx /prompt rm <name>`");
        return { metadata: { command: "prompt" } };
      }
      args = ["prompt", "rm", rest];
      progressMsg = "Removing prompt template...";
      break;
    default:
      args = ["prompt", "list"];
      progressMsg = "Listing prompt templates...";
      break;
  }
  args.push("--no-color");

  stream.progress(progressMsg);
  try {
    const { stdout, stderr } = await runCtx(args, cwd, token);
    const output = (stdout + stderr).trim();
    if (output) {
      stream.markdown("```\n" + output + "\n```");
    } else {
      stream.markdown("No prompt templates found.");
    }
  } catch (err: unknown) {
    stream.markdown(
      `**Error:** Prompt command failed.\n\n\`\`\`\n${err instanceof Error ? err.message : String(err)}\n\`\`\``
    );
  }
  return { metadata: { command: "prompt" } };
}

async function handleDecisions(
  stream: vscode.ChatResponseStream,
  prompt: string,
  cwd: string,
  token: vscode.CancellationToken
): Promise<CtxResult> {
  const args = ["decisions"];
  const subcmd = prompt.trim();
  if (subcmd) {
    args.push(...subcmd.split(/\s+/));
  }
  args.push("--no-color");
  stream.progress("Managing decisions...");
  try {
    const { stdout, stderr } = await runCtx(args, cwd, token);
    const output = (stdout + stderr).trim();
    if (output) {
      stream.markdown("```\n" + output + "\n```");
    } else {
      stream.markdown("No decisions found.");
    }
  } catch (err: unknown) {
    stream.markdown(
      `**Error:** Decisions command failed.\n\n\`\`\`\n${err instanceof Error ? err.message : String(err)}\n\`\`\``
    );
  }
  return { metadata: { command: "decisions" } };
}

async function handleLearnings(
  stream: vscode.ChatResponseStream,
  prompt: string,
  cwd: string,
  token: vscode.CancellationToken
): Promise<CtxResult> {
  const args = ["learnings"];
  const subcmd = prompt.trim();
  if (subcmd) {
    args.push(...subcmd.split(/\s+/));
  }
  args.push("--no-color");
  stream.progress("Managing learnings...");
  try {
    const { stdout, stderr } = await runCtx(args, cwd, token);
    const output = (stdout + stderr).trim();
    if (output) {
      stream.markdown("```\n" + output + "\n```");
    } else {
      stream.markdown("No learnings found.");
    }
  } catch (err: unknown) {
    stream.markdown(
      `**Error:** Learnings command failed.\n\n\`\`\`\n${err instanceof Error ? err.message : String(err)}\n\`\`\``
    );
  }
  return { metadata: { command: "learnings" } };
}

async function handleDeps(
  stream: vscode.ChatResponseStream,
  prompt: string,
  cwd: string,
  token: vscode.CancellationToken
): Promise<CtxResult> {
  const args = ["deps", "--no-color"];
  const parts = prompt.trim().split(/\s+/);
  for (let i = 0; i < parts.length; i++) {
    const p = parts[i]?.toLowerCase();
    if ((p === "--format" || p === "format") && parts[i + 1]) {
      args.push("--format", parts[++i]);
    } else if ((p === "--type" || p === "type") && parts[i + 1]) {
      args.push("--type", parts[++i]);
    } else if (p === "--external" || p === "external") {
      args.push("--external");
    }
  }
  stream.progress("Analyzing dependencies...");
  try {
    const { stdout, stderr } = await runCtx(args, cwd, token);
    const output = (stdout + stderr).trim();
    if (output) {
      stream.markdown("```\n" + output + "\n```");
    } else {
      stream.markdown("No dependency information available.");
    }
  } catch (err: unknown) {
    stream.markdown(
      `**Error:** Failed to analyze dependencies.\n\n\`\`\`\n${err instanceof Error ? err.message : String(err)}\n\`\`\``
    );
  }
  return { metadata: { command: "deps" } };
}

async function handleJournal(
  stream: vscode.ChatResponseStream,
  prompt: string,
  cwd: string,
  token: vscode.CancellationToken
): Promise<CtxResult> {
  const args = ["journal"];
  const subcmd = prompt.trim();
  if (subcmd) {
    args.push(...subcmd.split(/\s+/));
  }
  args.push("--no-color");
  stream.progress("Analyzing sessions...");
  try {
    const { stdout, stderr } = await runCtx(args, cwd, token);
    const output = (stdout + stderr).trim();
    if (output) {
      stream.markdown("```\n" + output + "\n```");
    } else {
      stream.markdown("No journal data available.");
    }
  } catch (err: unknown) {
    stream.markdown(
      `**Error:** Journal command failed.\n\n\`\`\`\n${err instanceof Error ? err.message : String(err)}\n\`\`\``
    );
  }
  return { metadata: { command: "journal" } };
}

async function handleReindex(
  stream: vscode.ChatResponseStream,
  cwd: string,
  token: vscode.CancellationToken
): Promise<CtxResult> {
  stream.progress("Regenerating indices...");
  try {
    const { stdout, stderr } = await runCtx(["reindex", "--no-color"], cwd, token);
    const output = (stdout + stderr).trim();
    if (output) {
      stream.markdown("```\n" + output + "\n```");
    } else {
      stream.markdown("Indices regenerated.");
    }
  } catch (err: unknown) {
    stream.markdown(
      `**Error:** Failed to regenerate indices.\n\n\`\`\`\n${err instanceof Error ? err.message : String(err)}\n\`\`\``
    );
  }
  return { metadata: { command: "reindex" } };
}

async function handleFreeform(
  request: vscode.ChatRequest,
  stream: vscode.ChatResponseStream,
  cwd: string,
  token: vscode.CancellationToken
): Promise<CtxResult> {
  const prompt = request.prompt.trim().toLowerCase();

  // Try to infer intent from natural language
  if (prompt.includes("init")) {
    return handleInit(stream, cwd, token);
  }
  if (prompt.includes("status")) {
    return handleStatus(stream, cwd, token);
  }
  if (prompt.includes("agent") || prompt.includes("context packet")) {
    return handleAgent(stream, request.prompt, cwd, token);
  }
  if (prompt.includes("load") || prompt.includes("assembled")) {
    return handleLoad(stream, request.prompt, cwd, token);
  }
  if (prompt.includes("compact") || prompt.includes("archive")) {
    return handleCompact(stream, request.prompt, cwd, token);
  }
  if (prompt.includes("sync") || prompt.includes("reconcile")) {
    return handleSync(stream, request.prompt, cwd, token);
  }
  if (prompt.includes("drift")) {
    return handleDrift(stream, cwd, token);
  }
  if (prompt.includes("recall") || prompt.includes("session") || prompt.includes("history")) {
    return handleRecall(stream, request.prompt, cwd, token);
  }
  if (prompt.includes("complete") || prompt.includes("done") || prompt.includes("finish")) {
    return handleComplete(stream, request.prompt, cwd, token);
  }
  if (prompt.includes("remind")) {
    return handleRemind(stream, request.prompt, cwd, token);
  }
  if (prompt.includes("task")) {
    return handleTasks(stream, request.prompt, cwd, token);
  }
  if (prompt.includes("pad") || prompt.includes("scratchpad") || prompt.includes("scratch")) {
    return handlePad(stream, request.prompt, cwd, token);
  }
  if (prompt.includes("notify") || prompt.includes("webhook")) {
    return handleNotify(stream, request.prompt, cwd, token);
  }
  if (prompt.includes("system") || prompt.includes("resource") || prompt.includes("bootstrap")) {
    return handleSystem(stream, request.prompt, cwd, token);
  }
  if (prompt.includes("change") || prompt.includes("diff") || prompt.includes("since")) {
    return handleChanges(stream, request.prompt, cwd, token);
  }
  if (prompt.includes("config") || prompt.includes("profile")) {
    return handleConfig(stream, request.prompt, cwd, token);
  }
  if (prompt.includes("doctor") || prompt.includes("health")) {
    return handleDoctor(stream, cwd, token);
  }
  if (prompt.includes("guide") || prompt.includes("cheat")) {
    return handleGuide(stream, request.prompt, cwd, token);
  }
  if (prompt.includes("why") || prompt.includes("philosophy") || prompt.includes("manifesto")) {
    return handleWhy(stream, request.prompt, cwd, token);
  }
  if (prompt.includes("memory") || prompt.includes("mirror")) {
    return handleMemory(stream, request.prompt, cwd, token);
  }
  if (prompt.includes("prompt") || prompt.includes("template")) {
    return handlePrompt(stream, request.prompt, cwd, token);
  }
  if (prompt.includes("decision")) {
    return handleDecisions(stream, request.prompt, cwd, token);
  }
  if (prompt.includes("learning")) {
    return handleLearnings(stream, request.prompt, cwd, token);
  }
  if (prompt.includes("dep") || prompt.includes("dependency")) {
    return handleDeps(stream, request.prompt, cwd, token);
  }
  if (prompt.includes("journal") || prompt.includes("session")) {
    return handleJournal(stream, request.prompt, cwd, token);
  }
  if (prompt.includes("reindex") || prompt.includes("index")) {
    return handleReindex(stream, cwd, token);
  }

  // Default: show help with available commands
  stream.markdown(
    "## ctx — Persistent Context for AI\n\n" +
      "Available commands:\n\n" +
      "| Command | Description |\n" +
      "|---------|-------------|\n" +
      "| `/init` | Initialize `.context/` directory |\n" +
      "| `/status` | Show context summary |\n" +
      "| `/agent` | Print AI-ready context packet |\n" +
      "| `/drift` | Detect stale or invalid context |\n" +
      "| `/recall` | Browse session history |\n" +
      "| `/hook` | Generate tool integration configs |\n" +
      "| `/add` | Add task, decision, or learning |\n" +
      "| `/load` | Output assembled context |\n" +
      "| `/compact` | Archive completed tasks |\n" +
      "| `/sync` | Reconcile context with codebase |\n" +
      "| `/complete` | Mark a task as completed |\n" +
      "| `/remind` | Manage session reminders |\n" +
      "| `/tasks` | Archive or snapshot tasks |\n" +
      "| `/pad` | Encrypted scratchpad |\n" +
      "| `/notify` | Webhook notifications |\n" +
      "| `/system` | System diagnostics |\n" +
      "| `/changes` | What changed since last session |\n" +
      "| `/config` | Manage runtime configuration |\n" +
      "| `/doctor` | Structural health check |\n" +
      "| `/guide` | Quick-reference cheat sheet |\n" +
      "| `/why` | Philosophy behind ctx |\n" +
      "| `/memory` | Bridge Claude Code auto memory |\n" +
      "| `/prompt` | Manage prompt templates |\n" +
      "| `/decisions` | Manage DECISIONS.md |\n" +
      "| `/learnings` | Manage LEARNINGS.md |\n" +
      "| `/deps` | Package dependency graph |\n" +
      "| `/journal` | Analyze AI sessions |\n" +
      "| `/reindex` | Regenerate indices |\n\n" +
      "Example: `@ctx /status` or `@ctx /add task Fix login bug`"
  );
  return { metadata: { command: "help" } };
}

const handler: vscode.ChatRequestHandler = async (
  request: vscode.ChatRequest,
  _context: vscode.ChatContext,
  stream: vscode.ChatResponseStream,
  token: vscode.CancellationToken
): Promise<CtxResult> => {
  const cwd = getWorkspaceRoot();
  if (!cwd) {
    stream.markdown(
      "**Error:** No workspace folder is open. Open a project folder first."
    );
    return { metadata: { command: request.command || "none" } };
  }
  // Auto-bootstrap: ensure ctx binary is available before any command
  try {
    stream.progress("Checking ctx installation...");
    await bootstrap();
  } catch (err: unknown) {
    stream.markdown(
      `**Error:** ctx CLI not found and auto-install failed.\n\n` +
        `\`\`\`\n${err instanceof Error ? err.message : String(err)}\n\`\`\`\n\n` +
        `Install manually: \`go install github.com/ActiveMemory/ctx/cmd/ctx@latest\` ` +
        `or download from [GitHub Releases](https://github.com/${GITHUB_REPO}/releases).`
    );
    return { metadata: { command: request.command || "none" } };
  }

  switch (request.command) {
    case "init":
      return handleInit(stream, cwd, token);
    case "status":
      return handleStatus(stream, cwd, token);
    case "agent":
      return handleAgent(stream, request.prompt, cwd, token);
    case "drift":
      return handleDrift(stream, cwd, token);
    case "recall":
      return handleRecall(stream, request.prompt, cwd, token);
    case "hook":
      return handleHook(stream, request.prompt, cwd, token);
    case "add":
      return handleAdd(stream, request.prompt, cwd, token);
    case "load":
      return handleLoad(stream, request.prompt, cwd, token);
    case "compact":
      return handleCompact(stream, request.prompt, cwd, token);
    case "sync":
      return handleSync(stream, request.prompt, cwd, token);
    case "complete":
      return handleComplete(stream, request.prompt, cwd, token);
    case "remind":
      return handleRemind(stream, request.prompt, cwd, token);
    case "tasks":
      return handleTasks(stream, request.prompt, cwd, token);
    case "pad":
      return handlePad(stream, request.prompt, cwd, token);
    case "notify":
      return handleNotify(stream, request.prompt, cwd, token);
    case "system":
      return handleSystem(stream, request.prompt, cwd, token);
    case "changes":
      return handleChanges(stream, request.prompt, cwd, token);
    case "config":
      return handleConfig(stream, request.prompt, cwd, token);
    case "doctor":
      return handleDoctor(stream, cwd, token);
    case "guide":
      return handleGuide(stream, request.prompt, cwd, token);
    case "why":
      return handleWhy(stream, request.prompt, cwd, token);
    case "memory":
      return handleMemory(stream, request.prompt, cwd, token);
    case "prompt":
      return handlePrompt(stream, request.prompt, cwd, token);
    case "decisions":
      return handleDecisions(stream, request.prompt, cwd, token);
    case "learnings":
      return handleLearnings(stream, request.prompt, cwd, token);
    case "deps":
      return handleDeps(stream, request.prompt, cwd, token);
    case "journal":
      return handleJournal(stream, request.prompt, cwd, token);
    case "reindex":
      return handleReindex(stream, cwd, token);
    default:
      return handleFreeform(request, stream, cwd, token);
  }
};

export function activate(extensionContext: vscode.ExtensionContext) {
  // Store extension context for auto-bootstrap binary downloads
  extensionCtx = extensionContext;

  // Kick off background bootstrap — don't block activation
  bootstrap().catch((err) => {
    // Errors will surface when user invokes a command
  });
  const participant = vscode.chat.createChatParticipant(
    PARTICIPANT_ID,
    handler
  );
  participant.iconPath = vscode.Uri.joinPath(
    extensionContext.extensionUri,
    "icon.png"
  );
  participant.followupProvider = {
    provideFollowups(
      result: CtxResult,
      _context: vscode.ChatContext,
      _token: vscode.CancellationToken
    ) {
      const followups: vscode.ChatFollowup[] = [];

      switch (result.metadata.command) {
        case "init":
          followups.push(
            { prompt: "Show my context status", command: "status" },
            {
              prompt: "Generate copilot integration",
              command: "hook",
            }
          );
          break;
        case "status":
          followups.push(
            { prompt: "Detect context drift", command: "drift" },
            { prompt: "Load full context", command: "load" }
          );
          break;
        case "drift":
          followups.push(
            { prompt: "Sync context with codebase", command: "sync" },
            { prompt: "Show context status", command: "status" }
          );
          break;
        case "complete":
          followups.push(
            { prompt: "Show context status", command: "status" },
            { prompt: "Archive completed tasks", command: "tasks" }
          );
          break;
        case "remind":
          followups.push(
            { prompt: "Show context status", command: "status" }
          );
          break;
        case "tasks":
          followups.push(
            { prompt: "Show context status", command: "status" },
            { prompt: "Compact context", command: "compact" }
          );
          break;
        case "pad":
          followups.push(
            { prompt: "List scratchpad", command: "pad" }
          );
          break;
        case "help":
          followups.push(
            { prompt: "Initialize project context", command: "init" },
            { prompt: "Show context status", command: "status" }
          );
          break;
      }

      return followups;
    },
  };

  extensionContext.subscriptions.push(participant);
}

export {
  runCtx,
  getCtxPath,
  getWorkspaceRoot,
  ensureCtxAvailable,
  bootstrap,
  getPlatformInfo,
  splitArgs,
  handleAdd,
  handleAgent,
  handleLoad,
  handleCompact,
  handleSync,
  handleRecall,
  handleComplete,
  handleRemind,
  handleTasks,
  handlePad,
  handleNotify,
  handleSystem,
  handleChanges,
  handleConfig,
  handleDoctor,
  handleGuide,
  handleWhy,
  handleMemory,
  handlePrompt,
  handleDecisions,
  handleLearnings,
  handleDeps,
  handleJournal,
  handleReindex,
};

export function deactivate() {}
