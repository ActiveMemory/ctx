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

async function handleInit(
  stream: vscode.ChatResponseStream,
  cwd: string,
  token: vscode.CancellationToken
): Promise<CtxResult> {
  stream.progress("Initializing .context/ directory...");
  try {
    const { stdout, stderr } = await runCtx(["init", "--no-color"], cwd, token);
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
  cwd: string,
  token: vscode.CancellationToken
): Promise<CtxResult> {
  stream.progress("Generating AI-ready context packet...");
  try {
    const { stdout, stderr } = await runCtx(["agent", "--no-color"], cwd, token);
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
  stream.progress("Searching session history...");
  try {
    const args = ["recall", "list", "--no-color"];
    if (prompt.trim()) {
      args.push("--query", prompt.trim());
    }
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
  const parts = prompt.trim().split(/\s+/);
  const type = parts[0];
  const content = parts.slice(1).join(" ");

  if (!type) {
    stream.markdown(
      "**Usage:** `@ctx /add <type> <content>`\n\n" +
        "Types: `task`, `decision`, `learning`\n\n" +
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
    const output = (stdout + stderr).trim();
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
  cwd: string,
  token: vscode.CancellationToken
): Promise<CtxResult> {
  stream.progress("Compacting context...");
  try {
    const { stdout, stderr } = await runCtx(["compact", "--no-color"], cwd, token);
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
  cwd: string,
  token: vscode.CancellationToken
): Promise<CtxResult> {
  stream.progress("Syncing context with codebase...");
  try {
    const { stdout, stderr } = await runCtx(["sync", "--no-color"], cwd, token);
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
  let progressMsg: string;

  switch (subcmd) {
    case "complete": {
      const taskRef = rest.trim();
      if (!taskRef) {
        stream.markdown(
          "**Usage:** `@ctx /task complete <task-id-or-text>`\n\n" +
            "Example: `@ctx /task complete 3` or " +
            "`@ctx /task complete Fix login bug`"
        );
        return { metadata: { command: "task" } };
      }
      args = ["task", "complete", taskRef];
      progressMsg = "Marking task as completed...";
      break;
    }
    case "archive":
      args = ["task", "archive"];
      progressMsg = "Archiving completed tasks...";
      break;
    case "snapshot":
      args = rest ? ["task", "snapshot", rest] : ["task", "snapshot"];
      progressMsg = "Creating task snapshot...";
      break;
    default:
      stream.markdown(
        "**Usage:** `@ctx /task <subcommand>`\n\n" +
          "| Subcommand | Description |\n" +
          "|------------|-------------|\n" +
          "| `complete <ref>` | Mark a task as completed |\n" +
          "| `archive` | Move completed tasks to archive |\n" +
          "| `snapshot [name]` | Create point-in-time snapshot |\n\n" +
          "Example: `@ctx /task complete 3` or " +
          "`@ctx /task archive`"
      );
      return { metadata: { command: "task" } };
  }
  args.push("--no-color");

  stream.progress(progressMsg);
  try {
    const { stdout, stderr } = await runCtx(args, cwd, token);
    const output = (stdout + stderr).trim();
    if (output) {
      stream.markdown("```\n" + output + "\n```");
    } else {
      switch (subcmd) {
        case "complete":
          stream.markdown(`Task **${rest.trim()}** marked as completed.`);
          break;
        case "archive":
          stream.markdown("Completed tasks archived.");
          break;
        default:
          stream.markdown("Task snapshot created.");
          break;
      }
    }
  } catch (err: unknown) {
    stream.markdown(
      `**Error:** Failed to ${subcmd} task.\n\n\`\`\`\n${err instanceof Error ? err.message : String(err)}\n\`\`\``
    );
  }
  return { metadata: { command: "task" } };
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
          "| `message list|show|edit|reset` | Manage hook messages |\n\n" +
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
      progressMsg = "Syncing memory bridge...";
      break;
    case "status":
      args = ["memory", "status"];
      progressMsg = "Checking memory status...";
      break;
    case "diff":
      args = ["memory", "diff"];
      progressMsg = "Diffing memory state...";
      break;
    case "import":
      args = ["memory", "import"];
      progressMsg = "Importing memory...";
      break;
    case "publish":
      args = ["memory", "publish"];
      progressMsg = "Publishing memory...";
      break;
    case "unpublish":
      args = ["memory", "unpublish"];
      progressMsg = "Unpublishing memory...";
      break;
    default:
      stream.markdown(
        "**Usage:** `@ctx /memory <subcommand>`\n\n" +
          "| Subcommand | Description |\n" +
          "|------------|-------------|\n" +
          "| `sync` | Synchronize memory bridge |\n" +
          "| `status` | Show memory bridge status |\n" +
          "| `diff` | Show memory diff |\n" +
          "| `import` | Import external memory |\n" +
          "| `publish` | Publish curated context |\n" +
          "| `unpublish` | Remove published context |\n\n" +
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
      stream.markdown(`Memory ${subcmd} completed.`);
    }
  } catch (err: unknown) {
    stream.markdown(
      `**Error:** Failed to run memory ${subcmd}.\n\n\`\`\`\n${err instanceof Error ? err.message : String(err)}\n\`\`\``
    );
  }
  return { metadata: { command: "memory" } };
}

async function handleJournal(
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
    case "site":
      args = ["journal", "site"];
      progressMsg = "Generating journal site...";
      break;
    case "obsidian":
      args = ["journal", "obsidian"];
      progressMsg = "Exporting journal to Obsidian...";
      break;
    default:
      stream.markdown(
        "**Usage:** `@ctx /journal <subcommand>`\n\n" +
          "| Subcommand | Description |\n" +
          "|------------|-------------|\n" +
          "| `site` | Generate journal site |\n" +
          "| `obsidian` | Export journal to Obsidian |\n\n" +
          "Example: `@ctx /journal site`"
      );
      return { metadata: { command: "journal" } };
  }
  args.push("--no-color");

  stream.progress(progressMsg);
  try {
    const { stdout, stderr } = await runCtx(args, cwd, token);
    const output = (stdout + stderr).trim();
    if (output) {
      stream.markdown("```\n" + output + "\n```");
    } else {
      stream.markdown(`Journal ${subcmd} completed.`);
    }
  } catch (err: unknown) {
    stream.markdown(
      `**Error:** Failed to run journal ${subcmd}.\n\n\`\`\`\n${err instanceof Error ? err.message : String(err)}\n\`\`\``
    );
  }
  return { metadata: { command: "journal" } };
}

async function handleDoctor(
  stream: vscode.ChatResponseStream,
  cwd: string,
  token: vscode.CancellationToken
): Promise<CtxResult> {
  stream.progress("Running context health diagnostics...");
  try {
    const { stdout, stderr } = await runCtx(["doctor", "--no-color"], cwd, token);
    const output = (stdout + stderr).trim();
    if (output) {
      stream.markdown("```\n" + output + "\n```");
    } else {
      stream.markdown("Context health check passed.");
    }
  } catch (err: unknown) {
    stream.markdown(
      `**Error:** Failed to run doctor.\n\n\`\`\`\n${err instanceof Error ? err.message : String(err)}\n\`\`\``
    );
  }
  return { metadata: { command: "doctor" } };
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
      args = rest ? ["config", "switch", rest] : ["config", "switch"];
      progressMsg = "Switching configuration...";
      break;
    case "status":
      args = ["config", "status"];
      progressMsg = "Checking configuration status...";
      break;
    case "schema":
      args = ["config", "schema"];
      progressMsg = "Showing configuration schema...";
      break;
    default:
      stream.markdown(
        "**Usage:** `@ctx /config <subcommand>`\n\n" +
          "| Subcommand | Description |\n" +
          "|------------|-------------|\n" +
          "| `switch` | Switch active configuration |\n" +
          "| `status` | Show current configuration |\n" +
          "| `schema` | Show configuration schema |\n\n" +
          "Example: `@ctx /config status` or `@ctx /config switch minimal`"
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
      stream.markdown(`Config ${subcmd} completed.`);
    }
  } catch (err: unknown) {
    stream.markdown(
      `**Error:** Failed to run config ${subcmd}.\n\n\`\`\`\n${err instanceof Error ? err.message : String(err)}\n\`\`\``
    );
  }
  return { metadata: { command: "config" } };
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
    case "add":
      args = rest ? ["prompt", "add", rest] : ["prompt", "add"];
      progressMsg = "Adding prompt template...";
      break;
    case "show":
      args = rest ? ["prompt", "show", rest] : ["prompt", "show"];
      progressMsg = "Showing prompt template...";
      break;
    case "rm":
      args = rest ? ["prompt", "rm", rest] : ["prompt", "rm"];
      progressMsg = "Removing prompt template...";
      break;
    default:
      stream.markdown(
        "**Usage:** `@ctx /prompt <subcommand>`\n\n" +
          "| Subcommand | Description |\n" +
          "|------------|-------------|\n" +
          "| `list` | List prompt templates |\n" +
          "| `add <name>` | Add a prompt template |\n" +
          "| `show <name>` | Show a prompt template |\n" +
          "| `rm <name>` | Remove a prompt template |\n\n" +
          "Example: `@ctx /prompt list` or `@ctx /prompt show review`"
      );
      return { metadata: { command: "prompt" } };
  }
  args.push("--no-color");

  stream.progress(progressMsg);
  try {
    const { stdout, stderr } = await runCtx(args, cwd, token);
    const output = (stdout + stderr).trim();
    if (output) {
      stream.markdown("```\n" + output + "\n```");
    } else {
      stream.markdown(`Prompt ${subcmd} completed.`);
    }
  } catch (err: unknown) {
    stream.markdown(
      `**Error:** Failed to run prompt ${subcmd}.\n\n\`\`\`\n${err instanceof Error ? err.message : String(err)}\n\`\`\``
    );
  }
  return { metadata: { command: "prompt" } };
}

async function handleWhy(
  stream: vscode.ChatResponseStream,
  prompt: string,
  cwd: string,
  token: vscode.CancellationToken
): Promise<CtxResult> {
  const filename = prompt.trim();
  const args = filename ? ["why", filename] : ["why"];
  args.push("--no-color");

  stream.progress("Looking up design rationale...");
  try {
    const { stdout, stderr } = await runCtx(args, cwd, token);
    const output = (stdout + stderr).trim();
    if (output) {
      stream.markdown("```\n" + output + "\n```");
    } else {
      stream.markdown("No rationale found.");
    }
  } catch (err: unknown) {
    stream.markdown(
      `**Error:** Failed to look up rationale.\n\n\`\`\`\n${err instanceof Error ? err.message : String(err)}\n\`\`\``
    );
  }
  return { metadata: { command: "why" } };
}

async function handleChange(
  stream: vscode.ChatResponseStream,
  cwd: string,
  token: vscode.CancellationToken
): Promise<CtxResult> {
  stream.progress("Checking recent codebase changes...");
  try {
    const { stdout, stderr } = await runCtx(["change", "--no-color"], cwd, token);
    const output = (stdout + stderr).trim();
    if (output) {
      stream.markdown("```\n" + output + "\n```");
    } else {
      stream.markdown("No recent changes found.");
    }
  } catch (err: unknown) {
    stream.markdown(
      `**Error:** Failed to check changes.\n\n\`\`\`\n${err instanceof Error ? err.message : String(err)}\n\`\`\``
    );
  }
  return { metadata: { command: "change" } };
}

async function handleDep(
  stream: vscode.ChatResponseStream,
  cwd: string,
  token: vscode.CancellationToken
): Promise<CtxResult> {
  stream.progress("Checking project dependencies...");
  try {
    const { stdout, stderr } = await runCtx(["dep", "--no-color"], cwd, token);
    const output = (stdout + stderr).trim();
    if (output) {
      stream.markdown("```\n" + output + "\n```");
    } else {
      stream.markdown("No dependencies found.");
    }
  } catch (err: unknown) {
    stream.markdown(
      `**Error:** Failed to check dependencies.\n\n\`\`\`\n${err instanceof Error ? err.message : String(err)}\n\`\`\``
    );
  }
  return { metadata: { command: "dep" } };
}

async function handleGuide(
  stream: vscode.ChatResponseStream,
  cwd: string,
  token: vscode.CancellationToken
): Promise<CtxResult> {
  stream.progress("Loading quick start guide...");
  try {
    const { stdout, stderr } = await runCtx(["guide", "--no-color"], cwd, token);
    const output = (stdout + stderr).trim();
    if (output) {
      stream.markdown(output);
    } else {
      stream.markdown("No guide content available.");
    }
  } catch (err: unknown) {
    stream.markdown(
      `**Error:** Failed to load guide.\n\n\`\`\`\n${err instanceof Error ? err.message : String(err)}\n\`\`\``
    );
  }
  return { metadata: { command: "guide" } };
}

async function handlePermission(
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
    case "snapshot":
      args = ["permission", "snapshot"];
      progressMsg = "Taking permission snapshot...";
      break;
    case "restore":
      args = ["permission", "restore"];
      progressMsg = "Restoring permissions...";
      break;
    default:
      stream.markdown(
        "**Usage:** `@ctx /permission <subcommand>`\n\n" +
          "| Subcommand | Description |\n" +
          "|------------|-------------|\n" +
          "| `snapshot` | Capture current file permissions |\n" +
          "| `restore` | Restore saved permissions |\n\n" +
          "Example: `@ctx /permission snapshot`"
      );
      return { metadata: { command: "permission" } };
  }
  args.push("--no-color");

  stream.progress(progressMsg);
  try {
    const { stdout, stderr } = await runCtx(args, cwd, token);
    const output = (stdout + stderr).trim();
    if (output) {
      stream.markdown("```\n" + output + "\n```");
    } else {
      stream.markdown(`Permission ${subcmd} completed.`);
    }
  } catch (err: unknown) {
    stream.markdown(
      `**Error:** Failed to ${subcmd} permissions.\n\n\`\`\`\n${err instanceof Error ? err.message : String(err)}\n\`\`\``
    );
  }
  return { metadata: { command: "permission" } };
}

async function handleSite(
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
    case "feed":
      args = ["site", "feed"];
      progressMsg = "Generating site feed...";
      break;
    default:
      stream.markdown(
        "**Usage:** `@ctx /site <subcommand>`\n\n" +
          "| Subcommand | Description |\n" +
          "|------------|-------------|\n" +
          "| `feed` | Generate documentation site feed |\n\n" +
          "Example: `@ctx /site feed`"
      );
      return { metadata: { command: "site" } };
  }
  args.push("--no-color");

  stream.progress(progressMsg);
  try {
    const { stdout, stderr } = await runCtx(args, cwd, token);
    const output = (stdout + stderr).trim();
    if (output) {
      stream.markdown("```\n" + output + "\n```");
    } else {
      stream.markdown(`Site ${subcmd} completed.`);
    }
  } catch (err: unknown) {
    stream.markdown(
      `**Error:** Failed to run site ${subcmd}.\n\n\`\`\`\n${err instanceof Error ? err.message : String(err)}\n\`\`\``
    );
  }
  return { metadata: { command: "site" } };
}

async function handleLoop(
  stream: vscode.ChatResponseStream,
  prompt: string,
  cwd: string,
  token: vscode.CancellationToken
): Promise<CtxResult> {
  const toolName = prompt.trim();
  const args = toolName ? ["loop", toolName] : ["loop"];
  args.push("--no-color");

  stream.progress("Generating iteration script...");
  try {
    const { stdout, stderr } = await runCtx(args, cwd, token);
    const output = (stdout + stderr).trim();
    if (output) {
      stream.markdown("```\n" + output + "\n```");
    } else {
      stream.markdown("No loop script generated.");
    }
  } catch (err: unknown) {
    stream.markdown(
      `**Error:** Failed to generate loop script.\n\n\`\`\`\n${err instanceof Error ? err.message : String(err)}\n\`\`\``
    );
  }
  return { metadata: { command: "loop" } };
}

async function handlePause(
  stream: vscode.ChatResponseStream,
  cwd: string,
  token: vscode.CancellationToken
): Promise<CtxResult> {
  stream.progress("Pausing context hooks...");
  try {
    const { stdout, stderr } = await runCtx(["pause", "--no-color"], cwd, token);
    const output = (stdout + stderr).trim();
    if (output) {
      stream.markdown("```\n" + output + "\n```");
    } else {
      stream.markdown("Context hooks paused for this session.");
    }
  } catch (err: unknown) {
    stream.markdown(
      `**Error:** Failed to pause hooks.\n\n\`\`\`\n${err instanceof Error ? err.message : String(err)}\n\`\`\``
    );
  }
  return { metadata: { command: "pause" } };
}

async function handleResume(
  stream: vscode.ChatResponseStream,
  cwd: string,
  token: vscode.CancellationToken
): Promise<CtxResult> {
  stream.progress("Resuming context hooks...");
  try {
    const { stdout, stderr } = await runCtx(["resume", "--no-color"], cwd, token);
    const output = (stdout + stderr).trim();
    if (output) {
      stream.markdown("```\n" + output + "\n```");
    } else {
      stream.markdown("Context hooks resumed.");
    }
  } catch (err: unknown) {
    stream.markdown(
      `**Error:** Failed to resume hooks.\n\n\`\`\`\n${err instanceof Error ? err.message : String(err)}\n\`\`\``
    );
  }
  return { metadata: { command: "resume" } };
}

async function handleReindex(
  stream: vscode.ChatResponseStream,
  cwd: string,
  token: vscode.CancellationToken
): Promise<CtxResult> {
  stream.progress("Rebuilding context file indices...");
  try {
    const { stdout, stderr } = await runCtx(["reindex", "--no-color"], cwd, token);
    const output = (stdout + stderr).trim();
    if (output) {
      stream.markdown("```\n" + output + "\n```");
    } else {
      stream.markdown("Context indices rebuilt.");
    }
  } catch (err: unknown) {
    stream.markdown(
      `**Error:** Failed to reindex.\n\n\`\`\`\n${err instanceof Error ? err.message : String(err)}\n\`\`\``
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
  if (prompt.includes("drift")) {
    return handleDrift(stream, cwd, token);
  }
  if (prompt.includes("recall") || prompt.includes("session") || prompt.includes("history")) {
    return handleRecall(stream, request.prompt, cwd, token);
  }
  if (prompt.includes("complete") || prompt.includes("done") || prompt.includes("finish")) {
    return handleTask(stream, "complete " + request.prompt, cwd, token);
  }
  if (prompt.includes("remind")) {
    return handleRemind(stream, request.prompt, cwd, token);
  }
  if (prompt.includes("task")) {
    return handleTask(stream, request.prompt, cwd, token);
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
  if (prompt.includes("memory")) {
    return handleMemory(stream, request.prompt, cwd, token);
  }
  if (prompt.includes("journal")) {
    return handleJournal(stream, request.prompt, cwd, token);
  }
  if (prompt.includes("doctor") || prompt.includes("health")) {
    return handleDoctor(stream, cwd, token);
  }
  if (prompt.includes("config") || prompt.includes("configuration")) {
    return handleConfig(stream, request.prompt, cwd, token);
  }
  if (prompt.includes("prompt") || prompt.includes("template")) {
    return handlePrompt(stream, request.prompt, cwd, token);
  }
  if (prompt.includes("why") || prompt.includes("rationale")) {
    return handleWhy(stream, request.prompt, cwd, token);
  }
  if (prompt.includes("change") || prompt.includes("recent")) {
    return handleChange(stream, cwd, token);
  }
  if (prompt.includes("dep") || prompt.includes("dependenc")) {
    return handleDep(stream, cwd, token);
  }
  if (prompt.includes("guide") || prompt.includes("quickstart") || prompt.includes("getting started")) {
    return handleGuide(stream, cwd, token);
  }
  if (prompt.includes("permission")) {
    return handlePermission(stream, request.prompt, cwd, token);
  }
  if (prompt.includes("site") || prompt.includes("feed")) {
    return handleSite(stream, request.prompt, cwd, token);
  }
  if (prompt.includes("loop") || prompt.includes("iterate")) {
    return handleLoop(stream, request.prompt, cwd, token);
  }
  if (prompt.includes("pause")) {
    return handlePause(stream, cwd, token);
  }
  if (prompt.includes("resume")) {
    return handleResume(stream, cwd, token);
  }
  if (prompt.includes("reindex") || prompt.includes("rebuild")) {
    return handleReindex(stream, cwd, token);
  }

  // Default: show help with available commands
  stream.markdown(
    "## ctx -- Persistent Context for AI\n\n" +
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
      "| `/task` | Task operations (complete, archive, snapshot) |\n" +
      "| `/remind` | Manage session reminders |\n" +
      "| `/pad` | Encrypted scratchpad |\n" +
      "| `/notify` | Webhook notifications |\n" +
      "| `/system` | System diagnostics |\n" +
      "| `/memory` | Memory bridge operations |\n" +
      "| `/journal` | Journal management |\n" +
      "| `/doctor` | Context health diagnostics |\n" +
      "| `/config` | Runtime configuration |\n" +
      "| `/prompt` | Prompt templates |\n" +
      "| `/why` | Design rationale for context files |\n" +
      "| `/change` | Recent codebase changes |\n" +
      "| `/dep` | Project dependencies |\n" +
      "| `/guide` | Quick start guide |\n" +
      "| `/permission` | Permission snapshot/restore |\n" +
      "| `/site` | Documentation site |\n" +
      "| `/loop` | Generate iteration scripts |\n" +
      "| `/pause` | Pause context hooks |\n" +
      "| `/resume` | Resume context hooks |\n" +
      "| `/reindex` | Rebuild context indices |\n\n" +
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
      return handleAgent(stream, cwd, token);
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
    case "task":
      return handleTask(stream, request.prompt, cwd, token);
    case "remind":
      return handleRemind(stream, request.prompt, cwd, token);
    case "pad":
      return handlePad(stream, request.prompt, cwd, token);
    case "notify":
      return handleNotify(stream, request.prompt, cwd, token);
    case "system":
      return handleSystem(stream, request.prompt, cwd, token);
    case "memory":
      return handleMemory(stream, request.prompt, cwd, token);
    case "journal":
      return handleJournal(stream, request.prompt, cwd, token);
    case "doctor":
      return handleDoctor(stream, cwd, token);
    case "config":
      return handleConfig(stream, request.prompt, cwd, token);
    case "prompt":
      return handlePrompt(stream, request.prompt, cwd, token);
    case "why":
      return handleWhy(stream, request.prompt, cwd, token);
    case "change":
      return handleChange(stream, cwd, token);
    case "dep":
      return handleDep(stream, cwd, token);
    case "guide":
      return handleGuide(stream, cwd, token);
    case "permission":
      return handlePermission(stream, request.prompt, cwd, token);
    case "site":
      return handleSite(stream, request.prompt, cwd, token);
    case "loop":
      return handleLoop(stream, request.prompt, cwd, token);
    case "pause":
      return handlePause(stream, cwd, token);
    case "resume":
      return handleResume(stream, cwd, token);
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
            { prompt: "Load full context", command: "load" },
            { prompt: "Run health check", command: "doctor" }
          );
          break;
        case "drift":
          followups.push(
            { prompt: "Sync context with codebase", command: "sync" },
            { prompt: "Show context status", command: "status" }
          );
          break;
        case "task":
          followups.push(
            { prompt: "Show context status", command: "status" },
            { prompt: "Compact context", command: "compact" }
          );
          break;
        case "remind":
          followups.push(
            { prompt: "Show context status", command: "status" }
          );
          break;
        case "pad":
          followups.push(
            { prompt: "List scratchpad", command: "pad" }
          );
          break;
        case "memory":
          followups.push(
            { prompt: "Check memory status", command: "memory" },
            { prompt: "Show context status", command: "status" }
          );
          break;
        case "doctor":
          followups.push(
            { prompt: "Show context status", command: "status" },
            { prompt: "Detect drift", command: "drift" }
          );
          break;
        case "config":
          followups.push(
            { prompt: "Show config status", command: "config" }
          );
          break;
        case "change":
          followups.push(
            { prompt: "Show context status", command: "status" }
          );
          break;
        case "pause":
          followups.push(
            { prompt: "Resume hooks", command: "resume" }
          );
          break;
        case "resume":
          followups.push(
            { prompt: "Show context status", command: "status" }
          );
          break;
        case "help":
          followups.push(
            { prompt: "Initialize project context", command: "init" },
            { prompt: "Show context status", command: "status" },
            { prompt: "Quick start guide", command: "guide" }
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
  handleTask,
  handleRemind,
  handlePad,
  handleNotify,
  handleSystem,
  handleMemory,
  handleJournal,
  handleDoctor,
  handleConfig,
  handlePrompt,
  handleWhy,
  handleChange,
  handleDep,
  handleGuide,
  handlePermission,
  handleSite,
  handleLoop,
  handlePause,
  handleResume,
  handleReindex,
};

export function deactivate() {}
