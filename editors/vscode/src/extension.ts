import * as vscode from "vscode";
import { execFile } from "child_process";
import * as fs from "fs";
import * as os from "os";
import * as path from "path";
import * as https from "https";

const PARTICIPANT_ID = "ctx.participant";
const GITHUB_REPO = "ActiveMemory/ctx";

interface CtxResult extends vscode.ChatResult {
  metadata: {
    command: string;
  };
}

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
    disposable = token?.onCancellationRequested(() => {
      child.kill();
    });
  });
}

/**
 * Check if .context/ directory exists in the workspace root.
 */
function hasContextDir(cwd: string): boolean {
  return fs.existsSync(path.join(cwd, ".context"));
}

async function handleInit(
  stream: vscode.ChatResponseStream,
  cwd: string,
  token: vscode.CancellationToken
): Promise<CtxResult> {
  stream.progress("Initializing .context/ directory...");
  try {
    const { stdout, stderr } = await runCtx(["init", "--caller", "vscode", "--no-color"], cwd, token);
    const output = mergeOutput(stdout, stderr);
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
      const hookOutput = mergeOutput(hookResult.stdout, hookResult.stderr);
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

    // Fire session-start since activate() missed it (no .context/ at activation time)
    runCtx(["system", "session-event", "--type", "start", "--caller", "vscode", "--no-color"], cwd).catch(() => {});
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
    const output = mergeOutput(stdout, stderr);
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
  stream.progress("Generating AI-ready context packet...");
  try {
    const args = ["agent", "--no-color"];
    const budgetMatch = prompt.match(/(?:--budget\s+|budget\s+)(\d+)/);
    if (budgetMatch) {
      args.splice(1, 0, "--budget", budgetMatch[1]);
    }
    const { stdout, stderr } = await runCtx(args, cwd, token);
    const output = mergeOutput(stdout, stderr);
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
    const output = mergeOutput(stdout, stderr);
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
  const rest = parts.slice(1).join(" ");

  let args: string[];
  let progressMsg: string;

  switch (subcmd) {
    case "show":
      if (!rest) {
        stream.markdown("**Usage:** `@ctx /recall show <id>`");
        return { metadata: { command: "recall" } };
      }
      args = ["recall", "show", rest];
      progressMsg = "Loading session details...";
      break;
    case "export":
      if (!rest) {
        stream.markdown("**Usage:** `@ctx /recall export <id>`");
        return { metadata: { command: "recall" } };
      }
      args = ["recall", "export", rest];
      progressMsg = "Exporting session...";
      break;
    case "lock":
      if (!rest) {
        stream.markdown("**Usage:** `@ctx /recall lock <id>`");
        return { metadata: { command: "recall" } };
      }
      args = ["recall", "lock", rest];
      progressMsg = "Locking session...";
      break;
    case "unlock":
      if (!rest) {
        stream.markdown("**Usage:** `@ctx /recall unlock <id>`");
        return { metadata: { command: "recall" } };
      }
      args = ["recall", "unlock", rest];
      progressMsg = "Unlocking session...";
      break;
    case "sync":
      args = ["recall", "sync"];
      progressMsg = "Syncing recall database...";
      break;
    case "list":
    default: {
      args = ["recall", "list"];
      progressMsg = "Searching session history...";
      const limitMatch = prompt.match(/(?:--limit\s+|limit\s+)(\d+)/);
      if (limitMatch) {
        args.push("--limit", limitMatch[1]);
      }
      const query = (subcmd === "list" ? rest : prompt.trim()).replace(/--limit\s+\d+/, "").replace(/limit\s+\d+/, "").trim();
      if (query) {
        args.push("--query", query);
      }
      break;
    }
  }
  args.push("--no-color");

  stream.progress(progressMsg);
  try {
    const { stdout, stderr } = await runCtx(args, cwd, token);
    const output = mergeOutput(stdout, stderr);
    if (output) {
      stream.markdown("```\n" + output + "\n```");
    } else {
      stream.markdown(subcmd === "list" || !subcmd ? "No session history found." : "No output.");
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
    const output = mergeOutput(stdout, stderr);
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
  const parts = prompt.trim().split(/\s+/);
  const type = parts[0];
  const content = parts.slice(1).join(" ");

  if (!type) {
    stream.markdown(
      "**Usage:** `@ctx /add <type> <content>`\n\n" +
        "Types: `task`, `decision`, `learning`, `convention`\n\n" +
        "Example: `@ctx /add task Implement user authentication`"
    );
    return { metadata: { command: "add" } };
  }

  stream.progress(`Adding ${type}...`);
  try {
    const args = ["add", type];
    if (content) {
      args.push(content);
    }
    const { stdout, stderr } = await runCtx(args, cwd, token);
    const output = mergeOutput(stdout, stderr);
    if (output) {
      stream.markdown("```\n" + output + "\n```");
    } else {
      stream.markdown(`Added **${type}**: ${content}`);
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
  cwd: string,
  token: vscode.CancellationToken
): Promise<CtxResult> {
  stream.progress("Loading assembled context...");
  try {
    const { stdout, stderr } = await runCtx(["load", "--no-color"], cwd, token);
    const output = mergeOutput(stdout, stderr);
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
  cwd: string,
  token: vscode.CancellationToken
): Promise<CtxResult> {
  stream.progress("Compacting context...");
  try {
    const { stdout, stderr } = await runCtx(["compact", "--no-color"], cwd, token);
    const output = mergeOutput(stdout, stderr);
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
  cwd: string,
  token: vscode.CancellationToken
): Promise<CtxResult> {
  stream.progress("Syncing context with codebase...");
  try {
    const { stdout, stderr } = await runCtx(["sync", "--no-color"], cwd, token);
    const output = mergeOutput(stdout, stderr);
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
    const output = mergeOutput(stdout, stderr);
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
    const output = mergeOutput(stdout, stderr);
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
    const output = mergeOutput(stdout, stderr);
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
    case "resolve":
      args = ["pad", "resolve"];
      progressMsg = "Resolving scratchpad conflicts...";
      break;
    case "import":
      if (!rest) {
        stream.markdown("**Usage:** `@ctx /pad import <path>`");
        return { metadata: { command: "pad" } };
      }
      args = ["pad", "import", rest];
      progressMsg = "Importing scratchpad archive...";
      break;
    case "export":
      args = rest ? ["pad", "export", rest] : ["pad", "export"];
      progressMsg = "Exporting scratchpad...";
      break;
    case "merge":
      args = rest ? ["pad", "merge", rest] : ["pad", "merge"];
      progressMsg = "Merging scratchpads...";
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
    const output = mergeOutput(stdout, stderr);
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
    const output = mergeOutput(stdout, stderr);
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
    case "doctor":
      args = ["doctor"];
      progressMsg = "Running diagnostics...";
      break;
    case "message":
      args = ["system", "message", ...parts.slice(1)];
      if (parts.length < 2 || !["list", "show", "edit", "reset"].includes(parts[1]?.toLowerCase())) {
        args = ["system", "message", "list"];
      }
      progressMsg = "Managing hook messages...";
      break;
    case "stats":
      args = ["system", "stats"];
      progressMsg = "Loading system stats...";
      break;
    case "backup":
      args = ["system", "backup"];
      progressMsg = "Running backup...";
      break;
    default:
      stream.markdown(
        "**Usage:** `@ctx /system <subcommand>`\n\n" +
          "| Subcommand | Description |\n" +
          "|------------|-------------|\n" +
          "| `resources` | Show system resource usage |\n" +
          "| `doctor` | Diagnose context health |\n" +
          "| `bootstrap` | Print context location for AI agents |\n" +
          "| `stats` | Show session and context stats |\n" +
          "| `backup` | Backup context data |\n" +
          "| `message list` | List hook message templates |\n" +
          "| `message show <hook>` | Show a hook message |\n" +
          "| `message edit <hook>` | Edit a hook message |\n" +
          "| `message reset <hook>` | Reset a hook message |\n\n" +
          "Example: `@ctx /system resources` or `@ctx /system message list`"
      );
      return { metadata: { command: "system" } };
  }
  args.push("--no-color");

  stream.progress(progressMsg);
  try {
    const { stdout, stderr } = await runCtx(args, cwd, token);
    const output = mergeOutput(stdout, stderr);
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

async function handleWrapup(
  stream: vscode.ChatResponseStream,
  cwd: string,
  token: vscode.CancellationToken
): Promise<CtxResult> {
  stream.progress("Generating session wrap-up...");
  try {
    // Gather status + drift in parallel for a comprehensive wrap-up
    const [statusResult, driftResult] = await Promise.all([
      runCtx(["status", "--no-color"], cwd, token),
      runCtx(["drift", "--no-color"], cwd, token),
    ]);
    const statusOutput = mergeOutput(statusResult.stdout, statusResult.stderr);
    const driftOutput = mergeOutput(driftResult.stdout, driftResult.stderr);

    stream.markdown("## Session Wrap-up\n\n");
    stream.markdown("### Context Status\n```\n" + statusOutput + "\n```\n\n");
    stream.markdown("### Drift Check\n```\n" + driftOutput + "\n```\n\n");
    stream.markdown(
      "**Before closing:** Review any open tasks in `.context/TASKS.md`. " +
        "Record decisions or learnings with `@ctx /add decision ...` or `@ctx /add learning ...`.\n"
    );

    // 2.15: Journal audit
    try {
      const stateDir = path.join(cwd, ".context", "state");
      if (fs.existsSync(stateDir)) {
        const journalFiles = fs.readdirSync(stateDir).filter((f) => f.includes("journal") || f.includes("event"));
        if (journalFiles.length > 0) {
          const latest = journalFiles.sort().slice(-1)[0];
          const content = fs.readFileSync(path.join(stateDir, latest), "utf-8");
          const lines = content.split("\n").filter((l) => l.trim());
          stream.markdown(`\n### Journal\n${lines.length} entries in \`${latest}\`. `);
          const today = new Date().toISOString().split("T")[0];
          if (!content.includes(today)) {
            stream.markdown("**No entries today.** Consider logging your work.\n");
          } else {
            stream.markdown("Today's entries present.\n");
          }
        }
      }
    } catch { /* non-fatal */ }

    // 2.18: Memory drift
    try {
      const memDir = path.join(cwd, ".context", "memory");
      if (fs.existsSync(memDir)) {
        const memFiles = fs.readdirSync(memDir).filter((f) => f.endsWith(".md"));
        if (memFiles.length > 0) {
          const contextFiles = ["DECISIONS.md", "LEARNINGS.md", "CONVENTIONS.md", "TASKS.md"];
          const drifts: string[] = [];
          for (const memFile of memFiles) {
            const memStat = fs.statSync(path.join(memDir, memFile));
            for (const ctxFile of contextFiles) {
              const ctxPath = path.join(cwd, ".context", ctxFile);
              if (fs.existsSync(ctxPath)) {
                const ctxStat = fs.statSync(ctxPath);
                if (memStat.mtimeMs < ctxStat.mtimeMs - 86400000) {
                  drifts.push(`\`memory/${memFile}\` older than \`${ctxFile}\``);
                }
              }
            }
          }
          if (drifts.length > 0) {
            stream.markdown("\n### Memory Drift\n" + drifts.map((d) => `- ${d}`).join("\n") + "\n");
          }
        }
      }
    } catch { /* non-fatal */ }

    // Record session end
    runCtx(["system", "session-event", "--type", "end", "--caller", "vscode", "--no-color"], cwd).catch(() => {});
  } catch (err: unknown) {
    stream.markdown(
      `**Error:** Wrap-up failed.\n\n\`\`\`\n${err instanceof Error ? err.message : String(err)}\n\`\`\``
    );
  }
  return { metadata: { command: "wrapup" } };
}

async function handleRemember(
  stream: vscode.ChatResponseStream,
  prompt: string,
  cwd: string,
  token: vscode.CancellationToken
): Promise<CtxResult> {
  stream.progress("Loading recent sessions...");
  try {
    const args = ["recall", "list", "--no-color"];
    const limitMatch = prompt.match(/(?:--limit\s+|limit\s+)(\d+)/);
    args.push("--limit", limitMatch ? limitMatch[1] : "3");
    const { stdout, stderr } = await runCtx(args, cwd, token);
    const output = mergeOutput(stdout, stderr);
    if (output) {
      stream.markdown("## Recent Sessions\n\n```\n" + output + "\n```");
    } else {
      stream.markdown("No recent sessions found. Start working and sessions will be recorded.");
    }
  } catch (err: unknown) {
    stream.markdown(
      `**Error:** Failed to load sessions.\n\n\`\`\`\n${err instanceof Error ? err.message : String(err)}\n\`\`\``
    );
  }
  return { metadata: { command: "remember" } };
}

async function handleNext(
  stream: vscode.ChatResponseStream,
  cwd: string,
  token: vscode.CancellationToken
): Promise<CtxResult> {
  stream.progress("Finding next task...");
  try {
    const tasksPath = path.join(cwd, ".context", "TASKS.md");
    if (!fs.existsSync(tasksPath)) {
      stream.markdown("No `.context/TASKS.md` found. Add tasks with `@ctx /add task ...`.");
      return { metadata: { command: "next" } };
    }
    const content = fs.readFileSync(tasksPath, "utf-8");
    const lines = content.split("\n");
    const openTasks = lines.filter((l) => /^\s*-\s*\[ \]/.test(l));
    if (openTasks.length === 0) {
      stream.markdown("All tasks completed! Add new tasks with `@ctx /add task ...`.");
    } else {
      stream.markdown("## Next Task\n\n" + openTasks[0].trim() + "\n");
      if (openTasks.length > 1) {
        stream.markdown(
          `\n*${openTasks.length - 1} more open task(s) remaining.*`
        );
      }
    }
  } catch (err: unknown) {
    stream.markdown(
      `**Error:** Failed to read tasks.\n\n\`\`\`\n${err instanceof Error ? err.message : String(err)}\n\`\`\``
    );
  }
  return { metadata: { command: "next" } };
}

async function handleBrainstorm(
  stream: vscode.ChatResponseStream,
  prompt: string,
  cwd: string,
  _token: vscode.CancellationToken
): Promise<CtxResult> {
  stream.progress("Loading ideas...");
  try {
    const ideasDir = path.join(cwd, "ideas");
    if (!fs.existsSync(ideasDir)) {
      stream.markdown(
        "No `ideas/` directory found. Run `@ctx /init` first, then add ideas to `ideas/README.md`."
      );
      return { metadata: { command: "brainstorm" } };
    }
    const readmePath = path.join(ideasDir, "README.md");
    if (fs.existsSync(readmePath)) {
      const content = fs.readFileSync(readmePath, "utf-8").trim();
      stream.markdown("## Current Ideas\n\n" + content + "\n");
    } else {
      stream.markdown("Ideas directory exists but `ideas/README.md` is empty.\n");
    }
    // List any other files in ideas/
    const files = fs.readdirSync(ideasDir).filter((f) => f !== "README.md" && f.endsWith(".md"));
    if (files.length > 0) {
      stream.markdown(
        "\n### Idea Files\n" + files.map((f) => `- \`ideas/${f}\``).join("\n") + "\n"
      );
    }
    if (prompt.trim()) {
      stream.markdown(
        "\n---\n\nTo develop **" + prompt.trim() + "** into a spec, create `specs/" +
          prompt.trim().toLowerCase().replace(/\s+/g, "-") + ".md` with your design."
      );
    } else {
      stream.markdown(
        "\n---\nTo develop an idea into a spec, run `@ctx /brainstorm <idea name>`."
      );
    }
  } catch (err: unknown) {
    stream.markdown(
      `**Error:** Failed to load ideas.\n\n\`\`\`\n${err instanceof Error ? err.message : String(err)}\n\`\`\``
    );
  }
  return { metadata: { command: "brainstorm" } };
}

async function handleReflect(
  stream: vscode.ChatResponseStream,
  cwd: string,
  token: vscode.CancellationToken
): Promise<CtxResult> {
  stream.progress("Reflecting on session...");
  try {
    const [statusResult, driftResult] = await Promise.all([
      runCtx(["status", "--no-color"], cwd, token),
      runCtx(["drift", "--no-color"], cwd, token),
    ]);
    const statusOutput = mergeOutput(statusResult.stdout, statusResult.stderr);
    const driftOutput = mergeOutput(driftResult.stdout, driftResult.stderr);

    stream.markdown("## Session Reflection\n\n");
    stream.markdown("### Current State\n```\n" + statusOutput + "\n```\n\n");
    if (driftOutput) {
      stream.markdown("### Drift Detected\n```\n" + driftOutput + "\n```\n\n");
    }
    stream.markdown(
      "### Worth Persisting?\n\n" +
        "Consider what happened this session:\n" +
        "- **Decision?** Did you make a design choice? → `@ctx /add decision ...`\n" +
        "- **Learning?** Hit a gotcha or discovered something? → `@ctx /add learning ...`\n" +
        "- **Convention?** Established a pattern? → `@ctx /add convention ...`\n" +
        "- **Task?** Identified work for later? → `@ctx /add task ...`\n"
    );

    // 2.15: Journal audit — check journal completeness
    try {
      const stateDir = path.join(cwd, ".context", "state");
      if (fs.existsSync(stateDir)) {
        const journalFiles = fs.readdirSync(stateDir).filter((f) => f.includes("journal") || f.includes("event"));
        if (journalFiles.length > 0) {
          const latest = journalFiles.sort().slice(-1)[0];
          const content = fs.readFileSync(path.join(stateDir, latest), "utf-8");
          const lines = content.split("\n").filter((l) => l.trim());
          stream.markdown(`\n### Journal\n${lines.length} entries in \`${latest}\`. `);
          const today = new Date().toISOString().split("T")[0];
          if (!content.includes(today)) {
            stream.markdown("**No entries today.** Consider logging your work.\n");
          } else {
            stream.markdown("Today's entries present.\n");
          }
        }
      }
    } catch { /* non-fatal */ }

    // 2.18: Memory drift — compare memory with context files
    try {
      const memDir = path.join(cwd, ".context", "memory");
      if (fs.existsSync(memDir)) {
        const memFiles = fs.readdirSync(memDir).filter((f) => f.endsWith(".md"));
        if (memFiles.length > 0) {
          const contextFiles = ["DECISIONS.md", "LEARNINGS.md", "CONVENTIONS.md", "TASKS.md"];
          const drifts: string[] = [];
          for (const memFile of memFiles) {
            const memStat = fs.statSync(path.join(memDir, memFile));
            for (const ctxFile of contextFiles) {
              const ctxPath = path.join(cwd, ".context", ctxFile);
              if (fs.existsSync(ctxPath)) {
                const ctxStat = fs.statSync(ctxPath);
                // Memory older than context by 24+ hours = potentially stale
                if (memStat.mtimeMs < ctxStat.mtimeMs - 86400000) {
                  drifts.push(`\`memory/${memFile}\` older than \`${ctxFile}\``);
                }
              }
            }
          }
          if (drifts.length > 0) {
            stream.markdown("\n### Memory Drift\n" + drifts.map((d) => `- ${d}`).join("\n") + "\n");
          }
        }
      }
    } catch { /* non-fatal */ }
  } catch (err: unknown) {
    stream.markdown(
      `**Error:** Reflection failed.\n\n\`\`\`\n${err instanceof Error ? err.message : String(err)}\n\`\`\``
    );
  }
  return { metadata: { command: "reflect" } };
}

async function handleSpec(
  stream: vscode.ChatResponseStream,
  prompt: string,
  cwd: string,
  _token: vscode.CancellationToken
): Promise<CtxResult> {
  stream.progress("Loading specs...");
  try {
    const specsDir = path.join(cwd, "specs");
    const tplDir = path.join(specsDir, "tpl");
    if (!prompt.trim()) {
      const specs = fs.existsSync(specsDir) ? fs.readdirSync(specsDir).filter((f) => f.endsWith(".md")) : [];
      const templates = fs.existsSync(tplDir) ? fs.readdirSync(tplDir).filter((f) => f.endsWith(".md")) : [];
      stream.markdown("## Specs\n\n");
      if (specs.length) {
        stream.markdown("### Existing\n" + specs.map((f) => `- \`specs/${f}\``).join("\n") + "\n\n");
      }
      if (templates.length) {
        stream.markdown("### Templates\n" + templates.map((f) => `- \`specs/tpl/${f}\``).join("\n") + "\n\nScaffold: `@ctx /spec <name>`\n");
      } else {
        stream.markdown("No templates in `specs/tpl/`. Create one to enable scaffolding.\n");
      }
    } else {
      const name = prompt.trim().toLowerCase().replace(/\s+/g, "-").replace(/[^a-z0-9\-]/g, "");
      const target = path.join(specsDir, `${name}.md`);
      if (fs.existsSync(target)) {
        const content = fs.readFileSync(target, "utf-8");
        stream.markdown(`\`specs/${name}.md\` exists:\n\n\`\`\`markdown\n${content}\n\`\`\``);
      } else {
        const templates = fs.existsSync(tplDir) ? fs.readdirSync(tplDir).filter((f) => f.endsWith(".md")) : [];
        let content: string;
        if (templates.length > 0) {
          content = fs.readFileSync(path.join(tplDir, templates[0]), "utf-8")
            .replace(/\{\{name\}\}/gi, name)
            .replace(/\{\{title\}\}/gi, prompt.trim());
        } else {
          content = `# ${prompt.trim()}\n\n## Problem\n\n## Proposal\n\n## Implementation\n\n## Verification\n`;
        }
        if (!fs.existsSync(specsDir)) { fs.mkdirSync(specsDir, { recursive: true }); }
        fs.writeFileSync(target, content, "utf-8");
        stream.markdown(`Created \`specs/${name}.md\`.\n\n\`\`\`markdown\n${content}\n\`\`\``);
      }
    }
  } catch (err: unknown) {
    stream.markdown(`**Error:** ${err instanceof Error ? err.message : String(err)}`);
  }
  return { metadata: { command: "spec" } };
}

async function handleImplement(
  stream: vscode.ChatResponseStream,
  cwd: string,
  _token: vscode.CancellationToken
): Promise<CtxResult> {
  stream.progress("Loading implementation plan...");
  try {
    const planPath = path.join(cwd, "IMPLEMENTATION_PLAN.md");
    if (!fs.existsSync(planPath)) {
      stream.markdown("No `IMPLEMENTATION_PLAN.md` found in project root.");
      return { metadata: { command: "implement" } };
    }
    const content = fs.readFileSync(planPath, "utf-8");
    const lines = content.split("\n");
    const done = lines.filter((l) => /^\s*-\s*\[x\]/i.test(l)).length;
    const open = lines.filter((l) => /^\s*-\s*\[ \]/.test(l)).length;
    const total = done + open;
    if (total > 0) {
      stream.markdown(`## Implementation Plan  (${done}/${total} steps done)\n\n`);
      const nextStep = lines.find((l) => /^\s*-\s*\[ \]/.test(l));
      if (nextStep) {
        stream.markdown("**Next step:** " + nextStep.replace(/^\s*-\s*\[ \]\s*/, "").trim() + "\n\n---\n\n");
      }
    } else {
      stream.markdown("## Implementation Plan\n\n");
    }
    stream.markdown(content);
  } catch (err: unknown) {
    stream.markdown(`**Error:** ${err instanceof Error ? err.message : String(err)}`);
  }
  return { metadata: { command: "implement" } };
}

async function handleVerify(
  stream: vscode.ChatResponseStream,
  cwd: string,
  token: vscode.CancellationToken
): Promise<CtxResult> {
  stream.progress("Running verification checks...");
  try {
    const results: string[] = [];
    try {
      const { stdout, stderr } = await runCtx(["doctor", "--no-color"], cwd, token);
      results.push("### Context Health\n```\n" + mergeOutput(stdout, stderr) + "\n```");
    } catch (err: unknown) {
      results.push("### Context Health\n```\nFailed: " + (err instanceof Error ? err.message : String(err)) + "\n```");
    }
    try {
      const { stdout, stderr } = await runCtx(["drift", "--no-color"], cwd, token);
      const output = mergeOutput(stdout, stderr);
      if (output) { results.push("### Drift\n```\n" + output + "\n```"); }
    } catch { /* non-fatal */ }
    stream.markdown("## Verification Report\n\n" + results.join("\n\n"));
  } catch (err: unknown) {
    stream.markdown(`**Error:** ${err instanceof Error ? err.message : String(err)}`);
  }
  return { metadata: { command: "verify" } };
}

async function handleMap(
  stream: vscode.ChatResponseStream,
  cwd: string,
  _token: vscode.CancellationToken
): Promise<CtxResult> {
  stream.progress("Mapping dependencies...");
  try {
    const results: string[] = [];
    const gomodPath = path.join(cwd, "go.mod");
    if (fs.existsSync(gomodPath)) {
      const content = fs.readFileSync(gomodPath, "utf-8");
      const moduleLine = content.match(/^module\s+(.+)$/m);
      const requires = content.match(/require\s*\(([\s\S]*?)\)/);
      results.push("### Go Module: " + (moduleLine ? moduleLine[1] : "unknown"));
      if (requires) {
        const deps = requires[1].trim().split("\n").filter((l) => l.trim() && !l.trim().startsWith("//"));
        results.push(`${deps.length} dependencies:\n\`\`\`\n${deps.map((d) => d.trim()).join("\n")}\n\`\`\``);
      }
    }
    const pkgPath = path.join(cwd, "package.json");
    if (fs.existsSync(pkgPath)) {
      const pkg = JSON.parse(fs.readFileSync(pkgPath, "utf-8"));
      const deps = Object.keys(pkg.dependencies || {});
      const devDeps = Object.keys(pkg.devDependencies || {});
      results.push("### Node Package: " + (pkg.name || "unknown"));
      if (deps.length) { results.push(`Dependencies: ${deps.join(", ")}`); }
      if (devDeps.length) { results.push(`Dev: ${devDeps.join(", ")}`); }
    }
    if (results.length === 0) {
      stream.markdown("No `go.mod` or `package.json` found.");
    } else {
      stream.markdown("## Dependency Map\n\n" + results.join("\n\n"));
    }
  } catch (err: unknown) {
    stream.markdown(`**Error:** ${err instanceof Error ? err.message : String(err)}`);
  }
  return { metadata: { command: "map" } };
}

async function handlePromptTpl(
  stream: vscode.ChatResponseStream,
  prompt: string,
  cwd: string,
  token: vscode.CancellationToken
): Promise<CtxResult> {
  stream.progress("Loading prompt templates...");
  try {
    const promptsDir = path.join(cwd, ".context", "prompts");
    const parts = prompt.trim().split(/\s+/);
    const subcmd = parts[0]?.toLowerCase();
    const rest = parts.slice(1).join(" ");

    if (subcmd === "add") {
      if (!rest) {
        stream.markdown("**Usage:** `@ctx /prompt add <file-path>`\n\nAdds a file as a prompt template.");
        return { metadata: { command: "prompt" } };
      }
      const args = ["prompt", "add", rest, "--no-color"];
      const { stdout, stderr } = await runCtx(args, cwd, token);
      const output = mergeOutput(stdout, stderr);
      stream.markdown(output ? "```\n" + output + "\n```" : `Template **${rest}** added.`);
      return { metadata: { command: "prompt" } };
    }

    if (subcmd === "rm" || subcmd === "remove" || subcmd === "delete") {
      if (!rest) {
        stream.markdown("**Usage:** `@ctx /prompt rm <name>`");
        return { metadata: { command: "prompt" } };
      }
      const args = ["prompt", "rm", rest, "--no-color"];
      const { stdout, stderr } = await runCtx(args, cwd, token);
      const output = mergeOutput(stdout, stderr);
      stream.markdown(output ? "```\n" + output + "\n```" : `Template **${rest}** removed.`);
      return { metadata: { command: "prompt" } };
    }

    if (!fs.existsSync(promptsDir)) {
      stream.markdown("No `.context/prompts/` directory found.");
      return { metadata: { command: "prompt" } };
    }
    const files = fs.readdirSync(promptsDir).filter((f) => f.endsWith(".md"));
    if (!prompt.trim()) {
      if (files.length === 0) {
        stream.markdown("`.context/prompts/` is empty. Add prompt templates with `@ctx /prompt add <file>`.");
      } else {
        stream.markdown("## Prompt Templates\n\n" + files.map((f) => `- \`${f}\``).join("\n") +
          "\n\nView: `@ctx /prompt <name>`\nAdd: `@ctx /prompt add <file>`\nRemove: `@ctx /prompt rm <name>`");
      }
    } else {
      const name = prompt.trim();
      const match = files.find((f) => f.toLowerCase().includes(name.toLowerCase()));
      if (match) {
        const content = fs.readFileSync(path.join(promptsDir, match), "utf-8");
        stream.markdown(`## ${match}\n\n${content}`);
      } else {
        stream.markdown(`No template matching "${name}". Available: ${files.join(", ")}`);
      }
    }
  } catch (err: unknown) {
    stream.markdown(`**Error:** ${err instanceof Error ? err.message : String(err)}`);
  }
  return { metadata: { command: "prompt" } };
}

async function handleBlog(
  stream: vscode.ChatResponseStream,
  prompt: string,
  cwd: string,
  _token: vscode.CancellationToken
): Promise<CtxResult> {
  stream.progress("Drafting blog post...");
  try {
    const sections: string[] = [];
    const decisionsPath = path.join(cwd, ".context", "DECISIONS.md");
    const learningsPath = path.join(cwd, ".context", "LEARNINGS.md");
    if (fs.existsSync(decisionsPath)) {
      const entries = fs.readFileSync(decisionsPath, "utf-8").split("\n").filter((l) => l.startsWith("- "));
      if (entries.length) { sections.push("## Key Decisions\n\n" + entries.slice(-5).join("\n")); }
    }
    if (fs.existsSync(learningsPath)) {
      const entries = fs.readFileSync(learningsPath, "utf-8").split("\n").filter((l) => l.startsWith("- "));
      if (entries.length) { sections.push("## Lessons Learned\n\n" + entries.slice(-5).join("\n")); }
    }
    const title = prompt.trim() || "Untitled Post";
    const date = new Date().toISOString().split("T")[0];
    stream.markdown(
      `# Blog Draft: ${title}\n\n*Date: ${date}*\n\n` +
        (sections.length ? sections.join("\n\n") : "No decisions or learnings to draw from.") +
        "\n\n---\n*Edit and refine this draft, then save to `docs/blog/`.*"
    );
  } catch (err: unknown) {
    stream.markdown(`**Error:** ${err instanceof Error ? err.message : String(err)}`);
  }
  return { metadata: { command: "blog" } };
}

async function handleChangelog(
  stream: vscode.ChatResponseStream,
  cwd: string,
  _token: vscode.CancellationToken
): Promise<CtxResult> {
  stream.progress("Generating changelog...");
  try {
    const result = await new Promise<string>((resolve, reject) => {
      execFile("git", ["log", "--oneline", "--no-decorate", "-20"], { cwd }, (err, stdout) => {
        if (err) { reject(err); } else { resolve(stdout); }
      });
    });
    if (result.trim()) {
      stream.markdown("## Recent Commits\n\n```\n" + result.trim() + "\n```\n\n" +
        "Use these to draft release notes or a changelog blog post.");
    } else {
      stream.markdown("No commits found.");
    }
  } catch (err: unknown) {
    stream.markdown(`**Error:** ${err instanceof Error ? err.message : String(err)}`);
  }
  return { metadata: { command: "changelog" } };
}

async function handleCheckLinks(
  stream: vscode.ChatResponseStream,
  cwd: string,
  _token: vscode.CancellationToken
): Promise<CtxResult> {
  stream.progress("Checking links in context files...");
  try {
    const contextDir = path.join(cwd, ".context");
    if (!fs.existsSync(contextDir)) {
      stream.markdown("No `.context/` directory found.");
      return { metadata: { command: "check-links" } };
    }
    const files = fs.readdirSync(contextDir).filter((f) => f.endsWith(".md"));
    const broken: string[] = [];
    let total = 0;
    for (const file of files) {
      const content = fs.readFileSync(path.join(contextDir, file), "utf-8");
      const linkRegex = /\[([^\]]*)\]\(([^)]+)\)/g;
      let m;
      while ((m = linkRegex.exec(content)) !== null) {
        const target = m[2];
        if (target.startsWith("http://") || target.startsWith("https://")) { continue; }
        total++;
        const resolved = path.resolve(contextDir, target);
        if (!fs.existsSync(resolved)) {
          broken.push(`- \`${file}\`: [${m[1]}](${target}) \u2192 not found`);
        }
      }
    }
    stream.markdown(`## Link Check\n\nChecked ${total} local links in ${files.length} context files.\n\n`);
    if (broken.length) {
      stream.markdown("### Broken Links\n" + broken.join("\n") + "\n");
    } else {
      stream.markdown("All local links are valid.\n");
    }
  } catch (err: unknown) {
    stream.markdown(`**Error:** ${err instanceof Error ? err.message : String(err)}`);
  }
  return { metadata: { command: "check-links" } };
}

async function handleJournal(
  stream: vscode.ChatResponseStream,
  prompt: string,
  cwd: string,
  token: vscode.CancellationToken
): Promise<CtxResult> {
  stream.progress("Checking journal...");
  try {
    const parts = prompt.trim().split(/\s+/);
    const subcmd = parts[0]?.toLowerCase();
    const args = ["journal"];
    let progressOverride: string | undefined;
    if (subcmd === "site") {
      args.push("site", ...parts.slice(1));
      progressOverride = "Exporting journal to static site...";
    } else if (subcmd === "obsidian") {
      args.push("obsidian", ...parts.slice(1));
      progressOverride = "Exporting journal to Obsidian...";
    } else if (prompt.trim()) {
      args.push(...parts);
    }
    args.push("--no-color");
    if (progressOverride) { stream.progress(progressOverride); }
    const { stdout, stderr } = await runCtx(args, cwd, token);
    const output = mergeOutput(stdout, stderr);
    if (output) { stream.markdown("```\n" + output + "\n```"); }
    else { stream.markdown("No journal output."); }
  } catch {
    try {
      const stateDir = path.join(cwd, ".context", "state");
      if (fs.existsSync(stateDir)) {
        const files = fs.readdirSync(stateDir).filter((f) => f.includes("journal") || f.includes("event"));
        if (files.length) {
          stream.markdown("## Journal Entries\n\n");
          for (const f of files.slice(-3)) {
            try {
              const content = fs.readFileSync(path.join(stateDir, f), "utf-8").trim();
              const preview = content.split("\n").slice(0, 10).join("\n");
              stream.markdown(`### ${f}\n\`\`\`\n${preview}\n\`\`\`\n\n`);
            } catch { /* skip unreadable */ }
          }
        } else {
          stream.markdown("No journal or event log files found in `.context/state/`.");
        }
      } else {
        stream.markdown("No `.context/state/` directory found.");
      }
    } catch (err: unknown) {
      stream.markdown(`**Error:** ${err instanceof Error ? err.message : String(err)}`);
    }
  }
  return { metadata: { command: "journal" } };
}

async function handleConsolidate(
  stream: vscode.ChatResponseStream,
  cwd: string,
  _token: vscode.CancellationToken
): Promise<CtxResult> {
  stream.progress("Scanning for overlapping entries...");
  try {
    const contextDir = path.join(cwd, ".context");
    const targetFiles = ["DECISIONS.md", "LEARNINGS.md", "CONVENTIONS.md", "TASKS.md"];
    const findings: string[] = [];
    for (const file of targetFiles) {
      const filePath = path.join(contextDir, file);
      if (!fs.existsSync(filePath)) { continue; }
      const entries = fs.readFileSync(filePath, "utf-8").split("\n")
        .filter((l) => /^\s*-\s/.test(l))
        .map((l) => l.replace(/^\s*-\s*(\[.\]\s*)?/, "").trim().toLowerCase());
      const seen = new Map<string, number>();
      for (const entry of entries) {
        const key = entry.replace(/[^a-z0-9\s]/g, "").replace(/\s+/g, " ");
        seen.set(key, (seen.get(key) || 0) + 1);
      }
      const dupes = [...seen.entries()].filter(([, count]) => count > 1);
      if (dupes.length) {
        findings.push(`### ${file}\n` + dupes.map(([text, count]) => `- "${text}" (\u00d7${count})`).join("\n"));
      }
    }
    stream.markdown("## Consolidation Report\n\n");
    if (findings.length) {
      stream.markdown(findings.join("\n\n") + "\n\nReview and merge manually.");
    } else {
      stream.markdown("No duplicate entries found across context files.");
    }
  } catch (err: unknown) {
    stream.markdown(`**Error:** ${err instanceof Error ? err.message : String(err)}`);
  }
  return { metadata: { command: "consolidate" } };
}

async function handleAudit(
  stream: vscode.ChatResponseStream,
  cwd: string,
  token: vscode.CancellationToken
): Promise<CtxResult> {
  stream.progress("Running alignment audit...");
  try {
    const { stdout: driftOut, stderr: driftErr } = await runCtx(["drift", "--no-color"], cwd, token);
    stream.markdown("## Alignment Audit\n\n### Drift\n```\n" + mergeOutput(driftOut, driftErr) + "\n```\n\n");
    const convPath = path.join(cwd, ".context", "CONVENTIONS.md");
    if (fs.existsSync(convPath)) {
      const entries = fs.readFileSync(convPath, "utf-8").split("\n").filter((l) => /^\s*-\s/.test(l));
      stream.markdown(`### Conventions: ${entries.length} documented\n\n`);
      if (entries.length === 0) {
        stream.markdown("**Warning:** No conventions documented. Run `@ctx /add convention ...`\n");
      }
    }
    const archPath = path.join(cwd, ".context", "ARCHITECTURE.md");
    if (fs.existsSync(archPath)) {
      const lines = fs.readFileSync(archPath, "utf-8").split("\n").filter((l) => l.trim() && !l.startsWith("#"));
      if (lines.length < 3) {
        stream.markdown("**Warning:** `ARCHITECTURE.md` appears sparse. Consider documenting system structure.\n");
      }
    }
  } catch (err: unknown) {
    stream.markdown(`**Error:** ${err instanceof Error ? err.message : String(err)}`);
  }
  return { metadata: { command: "audit" } };
}

async function handleWorktree(
  stream: vscode.ChatResponseStream,
  prompt: string,
  cwd: string,
  _token: vscode.CancellationToken
): Promise<CtxResult> {
  stream.progress("Managing worktrees...");
  try {
    const subcmd = prompt.trim().split(/\s+/)[0]?.toLowerCase() || "list";
    if (subcmd === "list" || !prompt.trim()) {
      const result = await new Promise<string>((resolve, reject) => {
        execFile("git", ["worktree", "list"], { cwd }, (err, stdout) => {
          if (err) { reject(err); } else { resolve(stdout); }
        });
      });
      stream.markdown("## Git Worktrees\n\n```\n" + result.trim() + "\n```\n\nCreate: `@ctx /worktree add <branch>`");
    } else if (subcmd === "add") {
      const branch = prompt.trim().split(/\s+/).slice(1).join("-").replace(/[^a-zA-Z0-9_\-\/]/g, "");
      if (!branch) {
        stream.markdown("Usage: `@ctx /worktree add <branch-name>`");
      } else {
        const worktreePath = path.join(path.dirname(cwd), path.basename(cwd) + "-" + branch);
        const result = await new Promise<string>((resolve, reject) => {
          execFile("git", ["worktree", "add", worktreePath, "-b", branch], { cwd }, (err, stdout, stderr) => {
            if (err) { reject(err); } else { resolve(stdout + stderr); }
          });
        });
        stream.markdown(`Worktree created at \`${worktreePath}\`.\n\n\`\`\`\n${result.trim()}\n\`\`\``);
      }
    } else {
      stream.markdown("**Usage:** `@ctx /worktree [list|add <branch>]`");
    }
  } catch (err: unknown) {
    stream.markdown(`**Error:** ${err instanceof Error ? err.message : String(err)}`);
  }
  return { metadata: { command: "worktree" } };
}

async function handlePause(
  stream: vscode.ChatResponseStream,
  cwd: string,
  token: vscode.CancellationToken
): Promise<CtxResult> {
  stream.progress("Pausing session...");
  try {
    const stateDir = path.join(cwd, ".context", "state");
    if (!fs.existsSync(stateDir)) { fs.mkdirSync(stateDir, { recursive: true }); }
    const { stdout } = await runCtx(["status", "--no-color"], cwd, token);
    const state = { paused_at: new Date().toISOString(), status: stdout.trim(), cwd };
    fs.writeFileSync(path.join(stateDir, "paused-session.json"), JSON.stringify(state, null, 2), "utf-8");
    stream.markdown("Session paused. State saved to `.context/state/paused-session.json`.\n\nResume with `@ctx /resume`.");
  } catch (err: unknown) {
    stream.markdown(`**Error:** ${err instanceof Error ? err.message : String(err)}`);
  }
  return { metadata: { command: "pause" } };
}

async function handleResume(
  stream: vscode.ChatResponseStream,
  cwd: string,
  token: vscode.CancellationToken
): Promise<CtxResult> {
  stream.progress("Resuming session...");
  try {
    const statePath = path.join(cwd, ".context", "state", "paused-session.json");
    if (!fs.existsSync(statePath)) {
      stream.markdown("No paused session found. Start fresh with `@ctx /status`.");
      return { metadata: { command: "resume" } };
    }
    const state = JSON.parse(fs.readFileSync(statePath, "utf-8"));
    stream.markdown("## Resuming Session\n\n" + `Paused at: ${state.paused_at}\n\n` +
      "### Status at pause\n```\n" + state.status + "\n```\n\n");
    try {
      const { stdout } = await runCtx(["status", "--no-color"], cwd, token);
      stream.markdown("### Current Status\n```\n" + stdout.trim() + "\n```\n");
    } catch { /* non-fatal */ }
    fs.unlinkSync(statePath);
    stream.markdown("\nSession resumed. Pause file removed.");
  } catch (err: unknown) {
    stream.markdown(`**Error:** ${err instanceof Error ? err.message : String(err)}`);
  }
  return { metadata: { command: "resume" } };
}

async function handleMemory(
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
    case "sync":
      args = ["memory", "sync"];
      progressMsg = "Syncing memory bridge...";
      break;
    case "status":
      args = ["memory", "status"];
      progressMsg = "Checking memory bridge status...";
      break;
    case "diff":
      args = ["memory", "diff"];
      progressMsg = "Comparing memory with context...";
      break;
    case "import":
      args = rest ? ["memory", "import", rest] : ["memory", "import"];
      progressMsg = "Importing memory...";
      break;
    case "publish":
      args = rest ? ["memory", "publish", rest] : ["memory", "publish"];
      progressMsg = "Publishing memory...";
      break;
    case "unpublish":
      args = rest ? ["memory", "unpublish", rest] : ["memory", "unpublish"];
      progressMsg = "Unpublishing memory...";
      break;
    default:
      stream.markdown(
        "**Usage:** `@ctx /memory <subcommand>`\n\n" +
          "| Subcommand | Description |\n" +
          "|------------|-------------|\n" +
          "| `sync` | Sync Claude Code memory bridge |\n" +
          "| `status` | Show memory bridge status |\n" +
          "| `diff` | Compare memory with context files |\n" +
          "| `import` | Import from Claude Code memory |\n" +
          "| `publish` | Publish context to Claude Code memory |\n" +
          "| `unpublish` | Remove published memory |\n\n" +
          "Example: `@ctx /memory sync` or `@ctx /memory diff`"
      );
      return { metadata: { command: "memory" } };
  }
  args.push("--no-color");

  stream.progress(progressMsg);
  try {
    const { stdout, stderr } = await runCtx(args, cwd, token);
    const output = mergeOutput(stdout, stderr);
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

async function handleDecisions(
  stream: vscode.ChatResponseStream,
  prompt: string,
  cwd: string,
  token: vscode.CancellationToken
): Promise<CtxResult> {
  const subcmd = prompt.trim().split(/\s+/)[0]?.toLowerCase();

  if (subcmd === "reindex") {
    stream.progress("Reindexing decisions...");
    try {
      const { stdout, stderr } = await runCtx(["decisions", "reindex", "--no-color"], cwd, token);
      const output = mergeOutput(stdout, stderr);
      stream.markdown(output ? "```\n" + output + "\n```" : "Decision index rebuilt.");
    } catch (err: unknown) {
      stream.markdown(`**Error:** ${err instanceof Error ? err.message : String(err)}`);
    }
  } else {
    stream.progress("Loading decisions...");
    try {
      const { stdout, stderr } = await runCtx(["decisions", "--no-color"], cwd, token);
      const output = mergeOutput(stdout, stderr);
      stream.markdown(output ? "```\n" + output + "\n```" : "No decisions found.");
    } catch (err: unknown) {
      stream.markdown(`**Error:** ${err instanceof Error ? err.message : String(err)}`);
    }
  }
  return { metadata: { command: "decisions" } };
}

async function handleLearnings(
  stream: vscode.ChatResponseStream,
  prompt: string,
  cwd: string,
  token: vscode.CancellationToken
): Promise<CtxResult> {
  const subcmd = prompt.trim().split(/\s+/)[0]?.toLowerCase();

  if (subcmd === "reindex") {
    stream.progress("Reindexing learnings...");
    try {
      const { stdout, stderr } = await runCtx(["learnings", "reindex", "--no-color"], cwd, token);
      const output = mergeOutput(stdout, stderr);
      stream.markdown(output ? "```\n" + output + "\n```" : "Learning index rebuilt.");
    } catch (err: unknown) {
      stream.markdown(`**Error:** ${err instanceof Error ? err.message : String(err)}`);
    }
  } else {
    stream.progress("Loading learnings...");
    try {
      const { stdout, stderr } = await runCtx(["learnings", "--no-color"], cwd, token);
      const output = mergeOutput(stdout, stderr);
      stream.markdown(output ? "```\n" + output + "\n```" : "No learnings found.");
    } catch (err: unknown) {
      stream.markdown(`**Error:** ${err instanceof Error ? err.message : String(err)}`);
    }
  }
  return { metadata: { command: "learnings" } };
}

async function handleConfig(
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
    case "switch":
      if (!rest) {
        stream.markdown("**Usage:** `@ctx /config switch <profile>`\n\nExample: `@ctx /config switch dev`");
        return { metadata: { command: "config" } };
      }
      args = ["config", "switch", rest];
      progressMsg = `Switching to profile "${rest}"...`;
      break;
    case "status":
      args = ["config", "status"];
      progressMsg = "Checking config status...";
      break;
    case "schema":
      args = ["config", "schema"];
      progressMsg = "Loading config schema...";
      break;
    default:
      stream.markdown(
        "**Usage:** `@ctx /config <subcommand>`\n\n" +
          "| Subcommand | Description |\n" +
          "|------------|-------------|\n" +
          "| `switch <profile>` | Switch to a config profile |\n" +
          "| `status` | Show current config profile |\n" +
          "| `schema` | Show config schema |\n\n" +
          "Example: `@ctx /config switch dev`"
      );
      return { metadata: { command: "config" } };
  }
  args.push("--no-color");

  stream.progress(progressMsg);
  try {
    const { stdout, stderr } = await runCtx(args, cwd, token);
    const output = mergeOutput(stdout, stderr);
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

async function handlePermissions(
  stream: vscode.ChatResponseStream,
  prompt: string,
  cwd: string,
  token: vscode.CancellationToken
): Promise<CtxResult> {
  const subcmd = prompt.trim().split(/\s+/)[0]?.toLowerCase();

  let args: string[];
  let progressMsg: string;

  switch (subcmd) {
    case "snapshot":
      args = ["permissions", "snapshot"];
      progressMsg = "Saving permissions snapshot...";
      break;
    case "restore":
      args = ["permissions", "restore"];
      progressMsg = "Restoring permissions...";
      break;
    default:
      stream.markdown(
        "**Usage:** `@ctx /permissions <subcommand>`\n\n" +
          "| Subcommand | Description |\n" +
          "|------------|-------------|\n" +
          "| `snapshot` | Backup current Claude settings |\n" +
          "| `restore` | Restore settings from backup |\n\n" +
          "Example: `@ctx /permissions snapshot`"
      );
      return { metadata: { command: "permissions" } };
  }
  args.push("--no-color");

  stream.progress(progressMsg);
  try {
    const { stdout, stderr } = await runCtx(args, cwd, token);
    const output = mergeOutput(stdout, stderr);
    if (output) {
      stream.markdown("```\n" + output + "\n```");
    } else {
      stream.markdown(subcmd === "snapshot" ? "Permissions snapshot saved." : "Permissions restored.");
    }
  } catch (err: unknown) {
    stream.markdown(
      `**Error:** Permissions command failed.\n\n\`\`\`\n${err instanceof Error ? err.message : String(err)}\n\`\`\``
    );
  }
  return { metadata: { command: "permissions" } };
}

async function handleChanges(
  stream: vscode.ChatResponseStream,
  prompt: string,
  cwd: string,
  token: vscode.CancellationToken
): Promise<CtxResult> {
  stream.progress("Checking what changed...");
  try {
    const args = ["changes"];
    const sinceMatch = prompt.match(/--since\s+(\S+)/);
    if (sinceMatch) {
      args.push("--since", sinceMatch[1]);
    }
    args.push("--no-color");
    const { stdout, stderr } = await runCtx(args, cwd, token);
    const output = mergeOutput(stdout, stderr);
    if (output) {
      stream.markdown("```\n" + output + "\n```");
    } else {
      stream.markdown("No changes detected since last session.");
    }
  } catch (err: unknown) {
    stream.markdown(
      `**Error:** Failed to check changes.\n\n\`\`\`\n${err instanceof Error ? err.message : String(err)}\n\`\`\``
    );
  }
  return { metadata: { command: "changes" } };
}

async function handleDeps(
  stream: vscode.ChatResponseStream,
  prompt: string,
  cwd: string,
  token: vscode.CancellationToken
): Promise<CtxResult> {
  stream.progress("Analyzing dependencies...");
  try {
    const args = ["deps"];
    const formatMatch = prompt.match(/--format\s+(\S+)/);
    if (formatMatch) {
      args.push("--format", formatMatch[1]);
    }
    if (prompt.includes("--external")) {
      args.push("--external");
    }
    args.push("--no-color");
    const { stdout, stderr } = await runCtx(args, cwd, token);
    const output = mergeOutput(stdout, stderr);
    if (output) {
      stream.markdown(output);
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

async function handleGuide(
  stream: vscode.ChatResponseStream,
  prompt: string,
  cwd: string,
  token: vscode.CancellationToken
): Promise<CtxResult> {
  stream.progress("Loading guide...");
  try {
    const args = ["guide"];
    if (prompt.includes("--skills")) {
      args.push("--skills");
    } else if (prompt.includes("--commands")) {
      args.push("--commands");
    }
    args.push("--no-color");
    const { stdout, stderr } = await runCtx(args, cwd, token);
    const output = mergeOutput(stdout, stderr);
    if (output) {
      stream.markdown(output);
    } else {
      stream.markdown("No guide output.");
    }
  } catch (err: unknown) {
    stream.markdown(
      `**Error:** Failed to load guide.\n\n\`\`\`\n${err instanceof Error ? err.message : String(err)}\n\`\`\``
    );
  }
  return { metadata: { command: "guide" } };
}

async function handleReindex(
  stream: vscode.ChatResponseStream,
  cwd: string,
  token: vscode.CancellationToken
): Promise<CtxResult> {
  stream.progress("Regenerating indices...");
  try {
    const { stdout, stderr } = await runCtx(["reindex", "--no-color"], cwd, token);
    const output = mergeOutput(stdout, stderr);
    if (output) {
      stream.markdown("```\n" + output + "\n```");
    } else {
      stream.markdown("Indices regenerated for DECISIONS.md and LEARNINGS.md.");
    }
  } catch (err: unknown) {
    stream.markdown(
      `**Error:** Failed to reindex.\n\n\`\`\`\n${err instanceof Error ? err.message : String(err)}\n\`\`\``
    );
  }
  return { metadata: { command: "reindex" } };
}

async function handleWhy(
  stream: vscode.ChatResponseStream,
  prompt: string,
  cwd: string,
  token: vscode.CancellationToken
): Promise<CtxResult> {
  stream.progress("Loading philosophy...");
  try {
    const args = ["why"];
    if (prompt.trim()) {
      args.push(prompt.trim());
    }
    args.push("--no-color");
    const { stdout, stderr } = await runCtx(args, cwd, token);
    const output = mergeOutput(stdout, stderr);
    if (output) {
      stream.markdown(output);
    } else {
      stream.markdown("No philosophy content available.");
    }
  } catch (err: unknown) {
    stream.markdown(
      `**Error:** Failed to load philosophy.\n\n\`\`\`\n${err instanceof Error ? err.message : String(err)}\n\`\`\``
    );
  }
  return { metadata: { command: "why" } };
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
  if (prompt.includes("wrap") || prompt.includes("end session") || prompt.includes("closing")) {
    return handleWrapup(stream, cwd, token);
  }
  if (prompt.includes("remember") || prompt.includes("last session") || prompt.includes("what were we")) {
    return handleRemember(stream, request.prompt, cwd, token);
  }
  if (prompt.includes("next") || prompt.includes("what should") || prompt.includes("pick task")) {
    return handleNext(stream, cwd, token);
  }
  if (prompt.includes("brainstorm") || prompt.includes("idea")) {
    return handleBrainstorm(stream, request.prompt, cwd, token);
  }
  if (prompt.includes("reflect") || prompt.includes("persist") || prompt.includes("worth saving")) {
    return handleReflect(stream, cwd, token);
  }
  if (prompt.includes("spec") || prompt.includes("scaffold")) {
    return handleSpec(stream, request.prompt, cwd, token);
  }
  if (prompt.includes("implement") || prompt.includes("execution plan")) {
    return handleImplement(stream, cwd, token);
  }
  if (prompt.includes("verify") || prompt.includes("qa") || prompt.includes("lint")) {
    return handleVerify(stream, cwd, token);
  }
  if (prompt.includes("map") || prompt.includes("dependencies") || prompt.includes("deps")) {
    return handleMap(stream, cwd, token);
  }
  if (prompt.includes("prompt template") || prompt.includes("prompts")) {
    return handlePromptTpl(stream, request.prompt, cwd, token);
  }
  if (prompt.includes("blog") && prompt.includes("changelog")) {
    return handleChangelog(stream, cwd, token);
  }
  if (prompt.includes("blog") || prompt.includes("post")) {
    return handleBlog(stream, request.prompt, cwd, token);
  }
  if (prompt.includes("changelog") || prompt.includes("release notes")) {
    return handleChangelog(stream, cwd, token);
  }
  if (prompt.includes("link") || prompt.includes("broken") || prompt.includes("dead link")) {
    return handleCheckLinks(stream, cwd, token);
  }
  if (prompt.includes("journal") || prompt.includes("log entries")) {
    return handleJournal(stream, request.prompt, cwd, token);
  }
  if (prompt.includes("consolidate") || prompt.includes("merge entries") || prompt.includes("duplicate")) {
    return handleConsolidate(stream, cwd, token);
  }
  if (prompt.includes("memory bridge") || prompt.includes("memory sync") || prompt.includes("memory status") || prompt.includes("memory diff") || prompt.includes("memory import") || prompt.includes("memory publish")) {
    return handleMemory(stream, request.prompt, cwd, token);
  }
  if (prompt.includes("decisions") || prompt.includes("decision list") || prompt.includes("decision reindex")) {
    return handleDecisions(stream, request.prompt, cwd, token);
  }
  if (prompt.includes("learnings") || prompt.includes("learning list") || prompt.includes("learning reindex")) {
    return handleLearnings(stream, request.prompt, cwd, token);
  }
  if (prompt.includes("config") || prompt.includes("profile") || prompt.includes("switch profile")) {
    return handleConfig(stream, request.prompt, cwd, token);
  }
  if (prompt.includes("permissions") || prompt.includes("permission snapshot") || prompt.includes("permission restore")) {
    return handlePermissions(stream, request.prompt, cwd, token);
  }
  if (prompt.includes("changes") || prompt.includes("what changed") || prompt.includes("since last session")) {
    return handleChanges(stream, request.prompt, cwd, token);
  }
  if (prompt.includes("deps") || prompt.includes("dependency graph") || prompt.includes("package graph")) {
    return handleDeps(stream, request.prompt, cwd, token);
  }
  if (prompt.includes("guide") || prompt.includes("cheat sheet") || prompt.includes("quick reference")) {
    return handleGuide(stream, request.prompt, cwd, token);
  }
  if (prompt.includes("reindex") || prompt.includes("rebuild index") || prompt.includes("regenerate index")) {
    return handleReindex(stream, cwd, token);
  }
  if (prompt.includes("why") || prompt.includes("philosophy") || prompt.includes("manifesto")) {
    return handleWhy(stream, request.prompt, cwd, token);
  }
  if (prompt.includes("audit") || prompt.includes("alignment")) {
    return handleAudit(stream, cwd, token);
  }
  if (prompt.includes("worktree")) {
    return handleWorktree(stream, request.prompt, cwd, token);
  }
  if (prompt.includes("pause") || prompt.includes("save state")) {
    return handlePause(stream, cwd, token);
  }
  if (prompt.includes("resume") || prompt.includes("restore state") || prompt.includes("continue session")) {
    return handleResume(stream, cwd, token);
  }

  // 2.5: Specs nudge — remind about specs when planning
  if (prompt.includes("plan") || prompt.includes("design") || prompt.includes("architect")) {
    const specsDir = path.join(cwd, "specs");
    if (fs.existsSync(specsDir)) {
      const specs = fs.readdirSync(specsDir).filter((f) => f.endsWith(".md"));
      if (specs.length > 0) {
        stream.markdown("> **specs/** has " + specs.length + " spec(s). Review with `@ctx /spec` before designing.\n\n");
      }
    }
  }

  // Default: show help with available commands
  stream.markdown(
    "## ctx — Persistent Context for AI\n\n" +
      "Available commands:\n\n" +
      "| Command | Description |\n" +
      "|---------|-------------|\n" +
      "| `/init` | Initialize `.context/` directory |\n" +
      "| `/status` | Show context summary |\n" +
      "| `/agent [--budget N]` | Print AI-ready context packet |\n" +
      "| `/drift` | Detect stale or invalid context |\n" +
      "| `/recall [show\\|export\\|lock\\|unlock\\|sync]` | Browse session history |\n" +
      "| `/hook` | Generate tool integration configs |\n" +
      "| `/add` | Add task, decision, learning, or convention |\n" +
      "| `/load` | Output assembled context |\n" +
      "| `/compact` | Archive completed tasks |\n" +
      "| `/sync` | Reconcile context with codebase |\n" +
      "| `/complete` | Mark a task as completed |\n" +
      "| `/remind` | Manage session reminders |\n" +
      "| `/tasks` | Archive or snapshot tasks |\n" +
      "| `/decisions [reindex]` | List or reindex decisions |\n" +
      "| `/learnings [reindex]` | List or reindex learnings |\n" +
      "| `/pad [resolve\\|import\\|export\\|merge]` | Encrypted scratchpad |\n" +
      "| `/notify` | Webhook notifications |\n" +
      "| `/memory [sync\\|status\\|diff\\|import\\|publish]` | Claude Code memory bridge |\n" +
      "| `/system [stats\\|backup\\|message]` | System diagnostics |\n" +
      "| `/config [switch\\|status\\|schema]` | Config profile management |\n" +
      "| `/permissions [snapshot\\|restore]` | Claude settings backup |\n" +
      "| `/wrapup` | End-of-session wrap-up |\n" +
      "| `/remember [--limit N]` | Recall recent sessions |\n" +
      "| `/next` | Show next open task |\n" +
      "| `/brainstorm` | Browse and develop ideas |\n" +
      "| `/reflect` | Surface items worth persisting |\n" +
      "| `/spec` | List or scaffold feature specs |\n" +
      "| `/implement` | Show implementation plan |\n" +
      "| `/verify` | Run verification checks |\n" +
      "| `/map` | Show dependency map |\n" +
      "| `/prompt [add\\|rm]` | Manage prompt templates |\n" +
      "| `/blog` | Draft blog post from context |\n" +
      "| `/changelog` | Recent commits for changelog |\n" +
      "| `/check-links` | Audit local links in context |\n" +
      "| `/journal [site\\|obsidian]` | View or export journal |\n" +
      "| `/consolidate` | Find duplicate entries |\n" +
      "| `/audit` | Alignment audit (drift + conventions) |\n" +
      "| `/worktree` | Git worktree management |\n" +
      "| `/pause` | Save session state |\n" +
      "| `/resume` | Restore paused session |\n" +
      "| `/changes [--since duration]` | Show what changed since last session |\n" +
      "| `/deps [--format mermaid\\|table\\|json]` | Package dependency graph |\n" +
      "| `/guide [--skills\\|--commands]` | Quick-reference cheat sheet |\n" +
      "| `/reindex` | Regenerate decision/learning indices |\n" +
      "| `/why [topic]` | Read ctx philosophy |\n\n" +
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

  // For non-init commands, verify .context/ exists
  if (request.command !== "init" && !hasContextDir(cwd)) {
    stream.markdown(
      "**Not initialized.** This project doesn't have a `.context/` directory yet.\n\n" +
        "Run `@ctx /init` to set up persistent context for this project."
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
      return handleLoad(stream, cwd, token);
    case "compact":
      return handleCompact(stream, cwd, token);
    case "sync":
      return handleSync(stream, cwd, token);
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
    case "wrapup":
      return handleWrapup(stream, cwd, token);
    case "remember":
      return handleRemember(stream, request.prompt, cwd, token);
    case "next":
      return handleNext(stream, cwd, token);
    case "brainstorm":
      return handleBrainstorm(stream, request.prompt, cwd, token);
    case "reflect":
      return handleReflect(stream, cwd, token);
    case "spec":
      return handleSpec(stream, request.prompt, cwd, token);
    case "implement":
      return handleImplement(stream, cwd, token);
    case "verify":
      return handleVerify(stream, cwd, token);
    case "map":
      return handleMap(stream, cwd, token);
    case "prompt":
      return handlePromptTpl(stream, request.prompt, cwd, token);
    case "blog":
      return handleBlog(stream, request.prompt, cwd, token);
    case "changelog":
      return handleChangelog(stream, cwd, token);
    case "check-links":
      return handleCheckLinks(stream, cwd, token);
    case "journal":
      return handleJournal(stream, request.prompt, cwd, token);
    case "consolidate":
      return handleConsolidate(stream, cwd, token);
    case "audit":
      return handleAudit(stream, cwd, token);
    case "worktree":
      return handleWorktree(stream, request.prompt, cwd, token);
    case "pause":
      return handlePause(stream, cwd, token);
    case "resume":
      return handleResume(stream, cwd, token);
    case "memory":
      return handleMemory(stream, request.prompt, cwd, token);
    case "decisions":
      return handleDecisions(stream, request.prompt, cwd, token);
    case "learnings":
      return handleLearnings(stream, request.prompt, cwd, token);
    case "config":
      return handleConfig(stream, request.prompt, cwd, token);
    case "permissions":
      return handlePermissions(stream, request.prompt, cwd, token);
    case "changes":
      return handleChanges(stream, request.prompt, cwd, token);
    case "deps":
      return handleDeps(stream, request.prompt, cwd, token);
    case "guide":
      return handleGuide(stream, request.prompt, cwd, token);
    case "reindex":
      return handleReindex(stream, cwd, token);
    case "why":
      return handleWhy(stream, request.prompt, cwd, token);
    default:
      return handleFreeform(request, stream, cwd, token);
  }
};

export function activate(extensionContext: vscode.ExtensionContext) {
  // Store extension context for auto-bootstrap binary downloads
  extensionCtx = extensionContext;

  // Kick off background bootstrap — don't block activation
  bootstrap().catch(() => {
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
        case "wrapup":
          followups.push(
            { prompt: "Add a decision", command: "add" },
            { prompt: "Add a learning", command: "add" }
          );
          break;
        case "remember":
          followups.push(
            { prompt: "Show context status", command: "status" },
            { prompt: "Load full context", command: "load" }
          );
          break;
        case "next":
          followups.push(
            { prompt: "Mark task completed", command: "complete" },
            { prompt: "Show all tasks", command: "status" }
          );
          break;
        case "brainstorm":
          followups.push(
            { prompt: "Show context status", command: "status" },
            { prompt: "Add a task", command: "add" }
          );
          break;
        case "reflect":
          followups.push(
            { prompt: "Add a decision", command: "add" },
            { prompt: "Add a learning", command: "add" },
            { prompt: "Wrap up session", command: "wrapup" }
          );
          break;
        case "spec":
          followups.push(
            { prompt: "Show implementation plan", command: "implement" },
            { prompt: "Run verification", command: "verify" }
          );
          break;
        case "implement":
          followups.push(
            { prompt: "Show next task", command: "next" },
            { prompt: "Run verification", command: "verify" }
          );
          break;
        case "verify":
          followups.push(
            { prompt: "Show context status", command: "status" },
            { prompt: "Run alignment audit", command: "audit" }
          );
          break;
        case "map":
          followups.push(
            { prompt: "Show context status", command: "status" }
          );
          break;
        case "blog":
        case "changelog":
          followups.push(
            { prompt: "Show context status", command: "status" }
          );
          break;
        case "consolidate":
          followups.push(
            { prompt: "Run alignment audit", command: "audit" },
            { prompt: "Show context status", command: "status" }
          );
          break;
        case "audit":
          followups.push(
            { prompt: "Fix drift", command: "sync" },
            { prompt: "Add a convention", command: "add" }
          );
          break;
        case "pause":
          followups.push(
            { prompt: "Resume session", command: "resume" }
          );
          break;
        case "resume":
          followups.push(
            { prompt: "Show next task", command: "next" },
            { prompt: "Show context status", command: "status" }
          );
          break;
        case "memory":
          followups.push(
            { prompt: "Show memory diff", command: "memory" },
            { prompt: "Show context status", command: "status" }
          );
          break;
        case "decisions":
          followups.push(
            { prompt: "Add a decision", command: "add" },
            { prompt: "Reindex decisions", command: "decisions" }
          );
          break;
        case "learnings":
          followups.push(
            { prompt: "Add a learning", command: "add" },
            { prompt: "Reindex learnings", command: "learnings" }
          );
          break;
        case "config":
          followups.push(
            { prompt: "Show config status", command: "config" },
            { prompt: "Show context status", command: "status" }
          );
          break;
        case "permissions":
          followups.push(
            { prompt: "Show context status", command: "status" }
          );
          break;
        case "changes":
          followups.push(
            { prompt: "Show context status", command: "status" },
            { prompt: "Load full context", command: "load" }
          );
          break;
        case "deps":
          followups.push(
            { prompt: "Show dependency map", command: "map" },
            { prompt: "Show context status", command: "status" }
          );
          break;
        case "guide":
          followups.push(
            { prompt: "Show context status", command: "status" },
            { prompt: "Read philosophy", command: "why" }
          );
          break;
        case "reindex":
          followups.push(
            { prompt: "List decisions", command: "decisions" },
            { prompt: "List learnings", command: "learnings" }
          );
          break;
        case "why":
          followups.push(
            { prompt: "Show guide", command: "guide" },
            { prompt: "Show context status", command: "status" }
          );
          break;
      }

      return followups;
    },
  };

  extensionContext.subscriptions.push(participant);

  // --- Command palette entries — open chat with the right slash command ---
  const paletteCommands: Array<[string, string]> = [
    ["ctx.init", "/init"],
    ["ctx.status", "/status"],
    ["ctx.agent", "/agent"],
    ["ctx.drift", "/drift"],
    ["ctx.recall", "/recall"],
    ["ctx.hook", "/hook"],
    ["ctx.add", "/add"],
    ["ctx.load", "/load"],
    ["ctx.compact", "/compact"],
    ["ctx.sync", "/sync"],
    ["ctx.complete", "/complete"],
    ["ctx.remind", "/remind"],
    ["ctx.tasks", "/tasks"],
    ["ctx.pad", "/pad"],
    ["ctx.notify", "/notify"],
    ["ctx.system", "/system"],
    ["ctx.wrapup", "/wrapup"],
    ["ctx.remember", "/remember"],
    ["ctx.next", "/next"],
    ["ctx.brainstorm", "/brainstorm"],
    ["ctx.reflect", "/reflect"],
    ["ctx.spec", "/spec"],
    ["ctx.implement", "/implement"],
    ["ctx.verify", "/verify"],
    ["ctx.map", "/map"],
    ["ctx.prompt", "/prompt"],
    ["ctx.blog", "/blog"],
    ["ctx.changelog", "/changelog"],
    ["ctx.checkLinks", "/check-links"],
    ["ctx.journal", "/journal"],
    ["ctx.consolidate", "/consolidate"],
    ["ctx.audit", "/audit"],
    ["ctx.worktree", "/worktree"],
    ["ctx.pause", "/pause"],
    ["ctx.resume", "/resume"],
    ["ctx.memory", "/memory"],
    ["ctx.decisions", "/decisions"],
    ["ctx.learnings", "/learnings"],
    ["ctx.config", "/config"],
    ["ctx.permissions", "/permissions"],
    ["ctx.changes", "/changes"],
    ["ctx.deps", "/deps"],
    ["ctx.guide", "/guide"],
    ["ctx.reindex", "/reindex"],
    ["ctx.why", "/why"],
  ];
  for (const [cmdId, slash] of paletteCommands) {
    extensionContext.subscriptions.push(
      vscode.commands.registerCommand(cmdId, () => {
        vscode.commands.executeCommand("workbench.action.chat.open", {
          query: `@ctx ${slash}`,
        });
      })
    );
  }

  // --- VS Code native hooks (equivalent to Claude Code hooks.json) ---
  const cwd = getWorkspaceRoot();

  // 2.6: onDidSave → task completion check (PostToolUse Edit/Write equivalent)
  const saveWatcher = vscode.workspace.onDidSaveTextDocument((doc) => {
    if (!cwd || !bootstrapDone || !hasContextDir(cwd)) {
      return;
    }
    // Only trigger for files inside the workspace, not for .context/ files themselves
    const rel = path.relative(cwd, doc.uri.fsPath);
    if (rel.startsWith(".context")) {
      return;
    }
    // Fire and forget — non-blocking background check
    runCtx(["system", "check-task-completion", "--no-color"], cwd).catch(
      () => {}
    );
  });
  extensionContext.subscriptions.push(saveWatcher);

  // 2.7: Git post-commit — detect commits and nudge for context capture
  try {
    const gitExtension = vscode.extensions.getExtension("vscode.git");
    if (gitExtension) {
      const setupGitHook = (git: any) => {
        try {
          const api = git.getAPI(1);
          if (api && api.repositories.length > 0) {
            const repo = api.repositories[0];
            let lastCommit = repo.state.HEAD?.commit;
            const commitListener = repo.state.onDidChange(() => {
              const currentCommit = repo.state.HEAD?.commit;
              if (currentCommit && currentCommit !== lastCommit) {
                lastCommit = currentCommit;
                if (!cwd || !bootstrapDone || !hasContextDir(cwd)) {
                  return;
                }
                vscode.window
                  .showInformationMessage(
                    "Commit succeeded. Record context or run QA?",
                    "Add Decision",
                    "Add Learning",
                    "Verify",
                    "Skip"
                  )
                  .then((choice) => {
                    if (choice === "Add Decision") {
                      vscode.commands.executeCommand(
                        "workbench.action.chat.open",
                        { query: "@ctx /add decision " }
                      );
                    } else if (choice === "Add Learning") {
                      vscode.commands.executeCommand(
                        "workbench.action.chat.open",
                        { query: "@ctx /add learning " }
                      );
                    } else if (choice === "Verify") {
                      vscode.commands.executeCommand(
                        "workbench.action.chat.open",
                        { query: "@ctx /verify" }
                      );
                    }
                  });
              }
            });
            extensionContext.subscriptions.push(commitListener);
          }
        } catch {
          // Git API not available
        }
      };

      if (gitExtension.isActive) {
        setupGitHook(gitExtension.exports);
      } else {
        gitExtension.activate().then(setupGitHook, () => {});
      }
    }
  } catch {
    // Git extension not available
  }

  // 2.9: Watch .context/ for external changes — refresh reminders
  if (cwd) {
    const contextWatcher = vscode.workspace.createFileSystemWatcher(
      new vscode.RelativePattern(cwd, ".context/**")
    );
    const onContextChange = () => {
      if (hasContextDir(cwd)) {
        updateReminderStatus(cwd);
        // Re-generate copilot-instructions.md when context files change
        runCtx(["hook", "copilot", "--write", "--no-color"], cwd).catch(() => {});
      }
    };
    contextWatcher.onDidChange(onContextChange);
    contextWatcher.onDidCreate(onContextChange);
    contextWatcher.onDidDelete(onContextChange);
    extensionContext.subscriptions.push(contextWatcher);
  }

  // 2.17: Watch dependency files for staleness
  if (cwd) {
    const depWatcher = vscode.workspace.createFileSystemWatcher(
      new vscode.RelativePattern(cwd, "{go.mod,go.sum,package.json,package-lock.json}")
    );
    depWatcher.onDidChange(() => {
      vscode.window.showInformationMessage(
        "Dependencies changed. Review with @ctx /map?",
        "View Map"
      ).then((choice) => {
        if (choice === "View Map") {
          vscode.commands.executeCommand("workbench.action.chat.open", { query: "@ctx /map" });
        }
      });
    });
    extensionContext.subscriptions.push(depWatcher);
  }

  // 2.10: Status bar reminder indicator
  reminderStatusBar = vscode.window.createStatusBarItem(
    vscode.StatusBarAlignment.Right,
    50
  );
  reminderStatusBar.name = "ctx Reminders";
  reminderStatusBar.command = undefined; // informational only
  extensionContext.subscriptions.push(reminderStatusBar);

  // Check reminders periodically (every 5 minutes)
  if (cwd && hasContextDir(cwd)) {
    updateReminderStatus(cwd);
    const reminderInterval = setInterval(() => {
      updateReminderStatus(cwd);
      // 2.14: Heartbeat — record session-alive timestamp
      try {
        const stateDir = path.join(cwd, ".context", "state");
        if (!fs.existsSync(stateDir)) { fs.mkdirSync(stateDir, { recursive: true }); }
        fs.writeFileSync(path.join(stateDir, "heartbeat"), new Date().toISOString(), "utf-8");
      } catch { /* non-fatal */ }
    }, 5 * 60 * 1000);
    extensionContext.subscriptions.push({
      dispose: () => clearInterval(reminderInterval),
    });

    // 2.12: Session start ceremony
    runCtx(
      ["system", "session-event", "--type", "start", "--caller", "vscode", "--no-color"],
      cwd
    ).catch(() => {});
  }
}

/**
 * Update the status bar reminder indicator by checking due reminders.
 */
function updateReminderStatus(cwd: string): void {
  if (!bootstrapDone || !reminderStatusBar) {
    return;
  }
  runCtx(["system", "check-reminders", "--no-color"], cwd)
    .then(({ stdout }) => {
      const trimmed = stdout.trim();
      if (trimmed && !trimmed.includes("no reminders")) {
        reminderStatusBar!.text = "$(bell) ctx";
        reminderStatusBar!.tooltip = trimmed;
        reminderStatusBar!.show();
      } else {
        reminderStatusBar!.hide();
      }
    })
    .catch(() => {
      reminderStatusBar!.hide();
    });
}

export {
  runCtx,
  getCtxPath,
  getWorkspaceRoot,
  ensureCtxAvailable,
  bootstrap,
  getPlatformInfo,
  handleComplete,
  handleRemind,
  handleTasks,
  handlePad,
  handleNotify,
  handleSystem,
  handleMemory,
  handleDecisions,
  handleLearnings,
  handleConfig,
  handlePermissions,
  handleChanges,
  handleDeps,
  handleGuide,
  handleReindex,
  handleWhy,
};

export function deactivate() {
  // 2.12: Session end ceremony
  const cwd = getWorkspaceRoot();
  if (cwd && hasContextDir(cwd)) {
    runCtx(
      ["system", "session-event", "--type", "end", "--caller", "vscode", "--no-color"],
      cwd
    ).catch(() => {});
  }
}
