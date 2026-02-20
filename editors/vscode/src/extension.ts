import * as vscode from "vscode";
import { execFile } from "child_process";

const PARTICIPANT_ID = "ctx.participant";

interface CtxResult extends vscode.ChatResult {
  metadata: {
    command: string;
  };
}

function getCtxPath(): string {
  return (
    vscode.workspace.getConfiguration("ctx").get<string>("executablePath") ||
    "ctx"
  );
}

function getWorkspaceRoot(): string | undefined {
  return vscode.workspace.workspaceFolders?.[0]?.uri.fsPath;
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
    const child = execFile(
      ctxPath,
      args,
      { cwd, maxBuffer: 1024 * 1024, timeout: 30000 },
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
      "| `/sync` | Reconcile context with codebase |\n\n" +
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
    default:
      return handleFreeform(request, stream, cwd, token);
  }
};

export function activate(extensionContext: vscode.ExtensionContext) {
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

export { runCtx, getCtxPath, getWorkspaceRoot };

export function deactivate() {}
