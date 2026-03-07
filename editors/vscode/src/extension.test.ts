import { describe, it, expect, vi, beforeEach } from "vitest";
import * as cp from "child_process";

// Mock vscode module (external, not bundled)
vi.mock("vscode", () => ({
  workspace: {
    getConfiguration: vi.fn(() => ({
      get: vi.fn(() => undefined),
    })),
    workspaceFolders: [{ uri: { fsPath: "/test/workspace" } }],
  },
  chat: {
    createChatParticipant: vi.fn(() => ({
      iconPath: null,
      followupProvider: null,
    })),
  },
  Uri: { joinPath: vi.fn() },
}));

vi.mock("child_process");

import {
  runCtx,
  getCtxPath,
  getWorkspaceRoot,
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
} from "./extension";

// Helper: create a fake CancellationToken
function fakeToken(cancelled = false) {
  const listeners: (() => void)[] = [];
  return {
    isCancellationRequested: cancelled,
    onCancellationRequested: vi.fn((cb: () => void) => {
      listeners.push(cb);
      return { dispose: vi.fn() };
    }),
    _fire: () => listeners.forEach((cb) => cb()),
  };
}

describe("getCtxPath", () => {
  it("returns 'ctx' when no config is set", () => {
    expect(getCtxPath()).toBe("ctx");
  });

  it("returns configured path when set", async () => {
    const vscode = await import("vscode");
    vi.mocked(vscode.workspace.getConfiguration).mockReturnValueOnce({
      get: vi.fn(() => "/custom/ctx"),
    } as never);
    expect(getCtxPath()).toBe("/custom/ctx");
  });
});

describe("getWorkspaceRoot", () => {
  it("returns first workspace folder path", () => {
    expect(getWorkspaceRoot()).toBe("/test/workspace");
  });

  it("returns undefined when no workspace is open", async () => {
    const vscode = await import("vscode");
    const original = vscode.workspace.workspaceFolders;
    (vscode.workspace as Record<string, unknown>).workspaceFolders = undefined;
    expect(getWorkspaceRoot()).toBeUndefined();
    (vscode.workspace as Record<string, unknown>).workspaceFolders = original;
  });
});

describe("runCtx", () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  it("resolves with stdout and stderr on success", async () => {
    vi.mocked(cp.execFile).mockImplementation(
      (_cmd: unknown, _args: unknown, _opts: unknown, cb: unknown) => {
        (cb as (e: null, out: string, err: string) => void)(
          null,
          "output",
          "errors"
        );
        return { kill: vi.fn() } as never;
      }
    );

    const result = await runCtx(["status"]);
    expect(result.stdout).toBe("output");
    expect(result.stderr).toBe("errors");
  });

  it("resolves on non-zero exit when output is present", async () => {
    vi.mocked(cp.execFile).mockImplementation(
      (_cmd: unknown, _args: unknown, _opts: unknown, cb: unknown) => {
        const err = new Error("exit 1");
        (cb as (e: Error, out: string, err: string) => void)(
          err,
          "",
          "drift detected"
        );
        return { kill: vi.fn() } as never;
      }
    );

    const result = await runCtx(["drift"]);
    expect(result.stderr).toBe("drift detected");
  });

  it("rejects on non-zero exit with no output", async () => {
    vi.mocked(cp.execFile).mockImplementation(
      (_cmd: unknown, _args: unknown, _opts: unknown, cb: unknown) => {
        const err = new Error("not found");
        (cb as (e: Error, out: string, err: string) => void)(err, "", "");
        return { kill: vi.fn() } as never;
      }
    );

    await expect(runCtx(["missing"])).rejects.toThrow("not found");
  });

  it("rejects immediately when token is already cancelled", async () => {
    const token = fakeToken(true);
    await expect(runCtx(["status"], "/test", token)).rejects.toThrow(
      "Cancelled"
    );
    expect(cp.execFile).not.toHaveBeenCalled();
  });

  it("kills child process when token fires cancellation", async () => {
    const killFn = vi.fn();
    let resolveCallback: (e: Error, out: string, err: string) => void;

    vi.mocked(cp.execFile).mockImplementation(
      (_cmd: unknown, _args: unknown, _opts: unknown, cb: unknown) => {
        resolveCallback = cb as typeof resolveCallback;
        return { kill: killFn } as never;
      }
    );

    const token = fakeToken();
    const promise = runCtx(["agent"], "/test", token);

    // Simulate cancellation
    token._fire();
    expect(killFn).toHaveBeenCalled();

    // Process exits after kill — no output so it rejects
    resolveCallback!(new Error("killed"), "", "");
    await expect(promise).rejects.toThrow("killed");
  });

  it("passes cwd to execFile", async () => {
    vi.mocked(cp.execFile).mockImplementation(
      (_cmd: unknown, _args: unknown, opts: unknown, cb: unknown) => {
        (cb as (e: null, out: string, err: string) => void)(null, "", "");
        return { kill: vi.fn() } as never;
      }
    );

    await runCtx(["status"], "/my/project");
    expect(cp.execFile).toHaveBeenCalledWith(
      "ctx",
      ["status"],
      expect.objectContaining({ cwd: "/my/project" }),
      expect.any(Function)
    );
  });

  it("disposes cancellation listener when process completes", async () => {
    const disposeFn = vi.fn();
    const token = {
      isCancellationRequested: false,
      onCancellationRequested: vi.fn(() => ({ dispose: disposeFn })),
    };

    vi.mocked(cp.execFile).mockImplementation(
      (_cmd: unknown, _args: unknown, _opts: unknown, cb: unknown) => {
        // Simulate async callback like real execFile
        process.nextTick(() =>
          (cb as (e: null, out: string, err: string) => void)(null, "done", "")
        );
        return { kill: vi.fn() } as never;
      }
    );

    await runCtx(["status"], "/test", token);
    expect(disposeFn).toHaveBeenCalled();
  });
});

describe("getPlatformInfo", () => {
  it("returns valid goos, goarch, and extension", () => {
    const info = getPlatformInfo();
    expect(["darwin", "linux", "windows"]).toContain(info.goos);
    expect(["amd64", "arm64"]).toContain(info.goarch);
    if (info.goos === "windows") {
      expect(info.ext).toBe(".exe");
    } else {
      expect(info.ext).toBe("");
    }
  });
});

// Helpers for handler tests
function fakeStream() {
  return {
    markdown: vi.fn(),
    progress: vi.fn(),
  };
}

function mockRunCtxSuccess(stdout: string, stderr = "") {
  vi.mocked(cp.execFile).mockImplementation(
    (_cmd: unknown, _args: unknown, _opts: unknown, cb: unknown) => {
      (cb as (e: null, out: string, err: string) => void)(null, stdout, stderr);
      return { kill: vi.fn() } as never;
    }
  );
}

function mockRunCtxError(message: string) {
  vi.mocked(cp.execFile).mockImplementation(
    (_cmd: unknown, _args: unknown, _opts: unknown, cb: unknown) => {
      const err = new Error(message);
      (cb as (e: Error, out: string, err: string) => void)(err, "", "");
      return { kill: vi.fn() } as never;
    }
  );
}

describe("handleComplete", () => {
  beforeEach(() => vi.clearAllMocks());

  it("shows usage when no task reference provided", async () => {
    const stream = fakeStream();
    const token = fakeToken();
    const result = await handleComplete(stream as never, "", "/test", token);
    expect(result.metadata.command).toBe("complete");
    expect(stream.markdown).toHaveBeenCalledWith(expect.stringContaining("Usage"));
  });

  it("runs complete command with task reference", async () => {
    mockRunCtxSuccess("Task 3 marked as done");
    const stream = fakeStream();
    const token = fakeToken();
    const result = await handleComplete(stream as never, "3", "/test", token);
    expect(result.metadata.command).toBe("complete");
    expect(stream.progress).toHaveBeenCalledWith("Marking task as completed...");
    expect(stream.markdown).toHaveBeenCalledWith(expect.stringContaining("Task 3 marked as done"));
  });

  it("runs complete with text reference", async () => {
    mockRunCtxSuccess("Completed: Fix login bug");
    const stream = fakeStream();
    const token = fakeToken();
    await handleComplete(stream as never, "Fix login bug", "/test", token);
    expect(cp.execFile).toHaveBeenCalledWith(
      "ctx",
      ["complete", "Fix login bug", "--no-color"],
      expect.anything(),
      expect.any(Function)
    );
  });

  it("handles errors gracefully", async () => {
    mockRunCtxError("task not found");
    const stream = fakeStream();
    const token = fakeToken();
    await handleComplete(stream as never, "99", "/test", token);
    expect(stream.markdown).toHaveBeenCalledWith(expect.stringContaining("Error"));
  });
});

describe("handleRemind", () => {
  beforeEach(() => vi.clearAllMocks());

  it("lists reminders when no subcommand given", async () => {
    mockRunCtxSuccess("1. Update docs\n2. Review PR");
    const stream = fakeStream();
    const token = fakeToken();
    await handleRemind(stream as never, "", "/test", token);
    expect(cp.execFile).toHaveBeenCalledWith(
      "ctx",
      ["remind", "list", "--no-color"],
      expect.anything(),
      expect.any(Function)
    );
  });

  it("adds reminder with 'add' subcommand", async () => {
    mockRunCtxSuccess("Reminder added");
    const stream = fakeStream();
    const token = fakeToken();
    await handleRemind(stream as never, "add Check CI status", "/test", token);
    expect(cp.execFile).toHaveBeenCalledWith(
      "ctx",
      ["remind", "add", "Check CI status", "--no-color"],
      expect.anything(),
      expect.any(Function)
    );
  });

  it("adds reminder when text provided without subcommand", async () => {
    mockRunCtxSuccess("Reminder added");
    const stream = fakeStream();
    const token = fakeToken();
    await handleRemind(stream as never, "Check CI status", "/test", token);
    expect(cp.execFile).toHaveBeenCalledWith(
      "ctx",
      ["remind", "add", "Check CI status", "--no-color"],
      expect.anything(),
      expect.any(Function)
    );
  });

  it("lists reminders with 'list' subcommand", async () => {
    mockRunCtxSuccess("No reminders");
    const stream = fakeStream();
    const token = fakeToken();
    await handleRemind(stream as never, "list", "/test", token);
    expect(cp.execFile).toHaveBeenCalledWith(
      "ctx",
      ["remind", "list", "--no-color"],
      expect.anything(),
      expect.any(Function)
    );
  });

  it("dismisses reminder by id", async () => {
    mockRunCtxSuccess("Dismissed reminder 2");
    const stream = fakeStream();
    const token = fakeToken();
    await handleRemind(stream as never, "dismiss 2", "/test", token);
    expect(cp.execFile).toHaveBeenCalledWith(
      "ctx",
      ["remind", "dismiss", "2", "--no-color"],
      expect.anything(),
      expect.any(Function)
    );
  });

  it("dismisses all when no id given", async () => {
    mockRunCtxSuccess("All dismissed");
    const stream = fakeStream();
    const token = fakeToken();
    await handleRemind(stream as never, "dismiss", "/test", token);
    expect(cp.execFile).toHaveBeenCalledWith(
      "ctx",
      ["remind", "dismiss", "--all", "--no-color"],
      expect.anything(),
      expect.any(Function)
    );
  });

  it("shows 'No reminders.' when output is empty", async () => {
    mockRunCtxSuccess("");
    const stream = fakeStream();
    const token = fakeToken();
    await handleRemind(stream as never, "list", "/test", token);
    expect(stream.markdown).toHaveBeenCalledWith("No reminders.");
  });

  it("handles errors gracefully", async () => {
    mockRunCtxError("failed");
    const stream = fakeStream();
    const token = fakeToken();
    await handleRemind(stream as never, "add test", "/test", token);
    expect(stream.markdown).toHaveBeenCalledWith(expect.stringContaining("Error"));
  });
});

describe("handleTasks", () => {
  beforeEach(() => vi.clearAllMocks());

  it("shows usage when no subcommand given", async () => {
    const stream = fakeStream();
    const token = fakeToken();
    const result = await handleTasks(stream as never, "", "/test", token);
    expect(result.metadata.command).toBe("tasks");
    expect(stream.markdown).toHaveBeenCalledWith(expect.stringContaining("Usage"));
  });

  it("runs archive subcommand", async () => {
    mockRunCtxSuccess("Archived 3 tasks");
    const stream = fakeStream();
    const token = fakeToken();
    await handleTasks(stream as never, "archive", "/test", token);
    expect(cp.execFile).toHaveBeenCalledWith(
      "ctx",
      ["tasks", "archive", "--no-color"],
      expect.anything(),
      expect.any(Function)
    );
    expect(stream.progress).toHaveBeenCalledWith("Archiving completed tasks...");
  });

  it("runs snapshot subcommand with name", async () => {
    mockRunCtxSuccess("Snapshot created");
    const stream = fakeStream();
    const token = fakeToken();
    await handleTasks(stream as never, "snapshot pre-refactor", "/test", token);
    expect(cp.execFile).toHaveBeenCalledWith(
      "ctx",
      ["tasks", "snapshot", "pre-refactor", "--no-color"],
      expect.anything(),
      expect.any(Function)
    );
  });

  it("runs snapshot without name", async () => {
    mockRunCtxSuccess("Snapshot created");
    const stream = fakeStream();
    const token = fakeToken();
    await handleTasks(stream as never, "snapshot", "/test", token);
    expect(cp.execFile).toHaveBeenCalledWith(
      "ctx",
      ["tasks", "snapshot", "--no-color"],
      expect.anything(),
      expect.any(Function)
    );
  });

  it("shows fallback message when archive output is empty", async () => {
    mockRunCtxSuccess("");
    const stream = fakeStream();
    const token = fakeToken();
    await handleTasks(stream as never, "archive", "/test", token);
    expect(stream.markdown).toHaveBeenCalledWith("Completed tasks archived.");
  });

  it("handles errors gracefully", async () => {
    mockRunCtxError("no tasks file");
    const stream = fakeStream();
    const token = fakeToken();
    await handleTasks(stream as never, "archive", "/test", token);
    expect(stream.markdown).toHaveBeenCalledWith(expect.stringContaining("Error"));
  });
});

describe("handlePad", () => {
  beforeEach(() => vi.clearAllMocks());

  it("lists all entries when no subcommand given", async () => {
    mockRunCtxSuccess("1: secret key\n2: API token");
    const stream = fakeStream();
    const token = fakeToken();
    await handlePad(stream as never, "", "/test", token);
    expect(cp.execFile).toHaveBeenCalledWith(
      "ctx",
      ["pad", "--no-color"],
      expect.anything(),
      expect.any(Function)
    );
  });

  it("adds entry with 'add' subcommand", async () => {
    mockRunCtxSuccess("Entry added");
    const stream = fakeStream();
    const token = fakeToken();
    await handlePad(stream as never, "add my secret note", "/test", token);
    expect(cp.execFile).toHaveBeenCalledWith(
      "ctx",
      ["pad", "add", "my secret note", "--no-color"],
      expect.anything(),
      expect.any(Function)
    );
  });

  it("shows usage when 'add' has no content", async () => {
    const stream = fakeStream();
    const token = fakeToken();
    await handlePad(stream as never, "add", "/test", token);
    expect(stream.markdown).toHaveBeenCalledWith(expect.stringContaining("Usage"));
  });

  it("shows entry by number", async () => {
    mockRunCtxSuccess("secret value");
    const stream = fakeStream();
    const token = fakeToken();
    await handlePad(stream as never, "show 1", "/test", token);
    expect(cp.execFile).toHaveBeenCalledWith(
      "ctx",
      ["pad", "show", "1", "--no-color"],
      expect.anything(),
      expect.any(Function)
    );
  });

  it("removes entry by number", async () => {
    mockRunCtxSuccess("Entry removed");
    const stream = fakeStream();
    const token = fakeToken();
    await handlePad(stream as never, "rm 2", "/test", token);
    expect(cp.execFile).toHaveBeenCalledWith(
      "ctx",
      ["pad", "rm", "2", "--no-color"],
      expect.anything(),
      expect.any(Function)
    );
  });

  it("shows usage when 'rm' has no number", async () => {
    const stream = fakeStream();
    const token = fakeToken();
    await handlePad(stream as never, "rm", "/test", token);
    expect(stream.markdown).toHaveBeenCalledWith(expect.stringContaining("Usage"));
  });

  it("edits entry", async () => {
    mockRunCtxSuccess("Entry updated");
    const stream = fakeStream();
    const token = fakeToken();
    await handlePad(stream as never, "edit 1 new text", "/test", token);
    expect(cp.execFile).toHaveBeenCalledWith(
      "ctx",
      ["pad", "edit", "1", "new", "text", "--no-color"],
      expect.anything(),
      expect.any(Function)
    );
  });

  it("moves entry", async () => {
    mockRunCtxSuccess("Entry moved");
    const stream = fakeStream();
    const token = fakeToken();
    await handlePad(stream as never, "mv 1 3", "/test", token);
    expect(cp.execFile).toHaveBeenCalledWith(
      "ctx",
      ["pad", "mv", "1", "3", "--no-color"],
      expect.anything(),
      expect.any(Function)
    );
  });

  it("shows 'Scratchpad is empty.' when output is empty", async () => {
    mockRunCtxSuccess("");
    const stream = fakeStream();
    const token = fakeToken();
    await handlePad(stream as never, "", "/test", token);
    expect(stream.markdown).toHaveBeenCalledWith("Scratchpad is empty.");
  });

  it("handles errors gracefully", async () => {
    mockRunCtxError("no key");
    const stream = fakeStream();
    const token = fakeToken();
    await handlePad(stream as never, "add secret", "/test", token);
    expect(stream.markdown).toHaveBeenCalledWith(expect.stringContaining("Error"));
  });
});

describe("handleNotify", () => {
  beforeEach(() => vi.clearAllMocks());

  it("shows usage when no subcommand given", async () => {
    const stream = fakeStream();
    const token = fakeToken();
    const result = await handleNotify(stream as never, "", "/test", token);
    expect(result.metadata.command).toBe("notify");
    expect(stream.markdown).toHaveBeenCalledWith(expect.stringContaining("Usage"));
  });

  it("runs setup subcommand", async () => {
    mockRunCtxSuccess("Webhook configured");
    const stream = fakeStream();
    const token = fakeToken();
    await handleNotify(stream as never, "setup", "/test", token);
    expect(cp.execFile).toHaveBeenCalledWith(
      "ctx",
      ["notify", "setup", "--no-color"],
      expect.anything(),
      expect.any(Function)
    );
    expect(stream.progress).toHaveBeenCalledWith("Setting up webhook...");
  });

  it("runs test subcommand", async () => {
    mockRunCtxSuccess("Test OK");
    const stream = fakeStream();
    const token = fakeToken();
    await handleNotify(stream as never, "test", "/test", token);
    expect(cp.execFile).toHaveBeenCalledWith(
      "ctx",
      ["notify", "test", "--no-color"],
      expect.anything(),
      expect.any(Function)
    );
  });

  it("sends notification with message", async () => {
    mockRunCtxSuccess("Sent");
    const stream = fakeStream();
    const token = fakeToken();
    await handleNotify(stream as never, "build done --event build", "/test", token);
    expect(cp.execFile).toHaveBeenCalledWith(
      "ctx",
      ["notify", "build", "done", "--event", "build", "--no-color"],
      expect.anything(),
      expect.any(Function)
    );
  });

  it("shows fallback on empty setup output", async () => {
    mockRunCtxSuccess("");
    const stream = fakeStream();
    const token = fakeToken();
    await handleNotify(stream as never, "setup", "/test", token);
    expect(stream.markdown).toHaveBeenCalledWith("Webhook configured.");
  });

  it("shows fallback on empty test output", async () => {
    mockRunCtxSuccess("");
    const stream = fakeStream();
    const token = fakeToken();
    await handleNotify(stream as never, "test", "/test", token);
    expect(stream.markdown).toHaveBeenCalledWith("Test notification sent.");
  });

  it("handles errors gracefully", async () => {
    mockRunCtxError("webhook failed");
    const stream = fakeStream();
    const token = fakeToken();
    await handleNotify(stream as never, "test", "/test", token);
    expect(stream.markdown).toHaveBeenCalledWith(expect.stringContaining("Error"));
  });
});

describe("handleSystem", () => {
  beforeEach(() => vi.clearAllMocks());

  it("shows usage when no subcommand given", async () => {
    const stream = fakeStream();
    const token = fakeToken();
    const result = await handleSystem(stream as never, "", "/test", token);
    expect(result.metadata.command).toBe("system");
    expect(stream.markdown).toHaveBeenCalledWith(expect.stringContaining("Usage"));
  });

  it("runs resources subcommand", async () => {
    mockRunCtxSuccess("Memory: 4GB / 16GB\nDisk: 50%");
    const stream = fakeStream();
    const token = fakeToken();
    await handleSystem(stream as never, "resources", "/test", token);
    expect(cp.execFile).toHaveBeenCalledWith(
      "ctx",
      ["system", "resources", "--no-color"],
      expect.anything(),
      expect.any(Function)
    );
    expect(stream.progress).toHaveBeenCalledWith("Checking system resources...");
  });

  it("runs bootstrap subcommand", async () => {
    mockRunCtxSuccess("context_dir: .context");
    const stream = fakeStream();
    const token = fakeToken();
    await handleSystem(stream as never, "bootstrap", "/test", token);
    expect(cp.execFile).toHaveBeenCalledWith(
      "ctx",
      ["system", "bootstrap", "--no-color"],
      expect.anything(),
      expect.any(Function)
    );
    expect(stream.progress).toHaveBeenCalledWith("Running bootstrap...");
  });

  it("runs message subcommand with arguments", async () => {
    mockRunCtxSuccess("Hook messages listed");
    const stream = fakeStream();
    const token = fakeToken();
    await handleSystem(stream as never, "message list", "/test", token);
    expect(cp.execFile).toHaveBeenCalledWith(
      "ctx",
      ["system", "message", "list", "--no-color"],
      expect.anything(),
      expect.any(Function)
    );
  });

  it("shows 'No output.' when output is empty", async () => {
    mockRunCtxSuccess("");
    const stream = fakeStream();
    const token = fakeToken();
    await handleSystem(stream as never, "resources", "/test", token);
    expect(stream.markdown).toHaveBeenCalledWith("No output.");
  });

  it("handles errors gracefully", async () => {
    mockRunCtxError("system error");
    const stream = fakeStream();
    const token = fakeToken();
    await handleSystem(stream as never, "resources", "/test", token);
    expect(stream.markdown).toHaveBeenCalledWith(expect.stringContaining("Error"));
  });
});

describe("handleChanges", () => {
  beforeEach(() => vi.clearAllMocks());

  it("runs changes without --since when no prompt", async () => {
    mockRunCtxSuccess("3 files changed");
    const stream = fakeStream();
    const token = fakeToken();
    await handleChanges(stream as never, "", "/test", token);
    expect(cp.execFile).toHaveBeenCalledWith(
      "ctx",
      ["changes", "--no-color"],
      expect.anything(),
      expect.any(Function)
    );
  });

  it("passes --since when prompt provided", async () => {
    mockRunCtxSuccess("2 files changed");
    const stream = fakeStream();
    const token = fakeToken();
    await handleChanges(stream as never, "24h", "/test", token);
    expect(cp.execFile).toHaveBeenCalledWith(
      "ctx",
      ["changes", "--no-color", "--since", "24h"],
      expect.anything(),
      expect.any(Function)
    );
  });

  it("shows 'No changes detected.' when output is empty", async () => {
    mockRunCtxSuccess("");
    const stream = fakeStream();
    const token = fakeToken();
    await handleChanges(stream as never, "", "/test", token);
    expect(stream.markdown).toHaveBeenCalledWith("No changes detected.");
  });

  it("handles errors gracefully", async () => {
    mockRunCtxError("git error");
    const stream = fakeStream();
    const token = fakeToken();
    await handleChanges(stream as never, "", "/test", token);
    expect(stream.markdown).toHaveBeenCalledWith(expect.stringContaining("Error"));
  });
});

describe("handleConfig", () => {
  beforeEach(() => vi.clearAllMocks());

  it("shows usage when no subcommand given", async () => {
    const stream = fakeStream();
    const token = fakeToken();
    const result = await handleConfig(stream as never, "", "/test", token);
    expect(result.metadata.command).toBe("config");
    expect(stream.markdown).toHaveBeenCalledWith(expect.stringContaining("Usage"));
  });

  it("runs status subcommand", async () => {
    mockRunCtxSuccess("Profile: base");
    const stream = fakeStream();
    const token = fakeToken();
    await handleConfig(stream as never, "status", "/test", token);
    expect(cp.execFile).toHaveBeenCalledWith(
      "ctx",
      ["config", "status", "--no-color"],
      expect.anything(),
      expect.any(Function)
    );
  });

  it("runs switch subcommand with profile", async () => {
    mockRunCtxSuccess("Switched to dev");
    const stream = fakeStream();
    const token = fakeToken();
    await handleConfig(stream as never, "switch dev", "/test", token);
    expect(cp.execFile).toHaveBeenCalledWith(
      "ctx",
      ["config", "switch", "dev", "--no-color"],
      expect.anything(),
      expect.any(Function)
    );
  });

  it("handles errors gracefully", async () => {
    mockRunCtxError("config error");
    const stream = fakeStream();
    const token = fakeToken();
    await handleConfig(stream as never, "status", "/test", token);
    expect(stream.markdown).toHaveBeenCalledWith(expect.stringContaining("Error"));
  });
});

describe("handleDoctor", () => {
  beforeEach(() => vi.clearAllMocks());

  it("runs doctor command", async () => {
    mockRunCtxSuccess("All checks passed");
    const stream = fakeStream();
    const token = fakeToken();
    await handleDoctor(stream as never, "/test", token);
    expect(cp.execFile).toHaveBeenCalledWith(
      "ctx",
      ["doctor", "--no-color"],
      expect.anything(),
      expect.any(Function)
    );
    expect(stream.progress).toHaveBeenCalledWith("Running health checks...");
  });

  it("handles errors gracefully", async () => {
    mockRunCtxError("doctor error");
    const stream = fakeStream();
    const token = fakeToken();
    await handleDoctor(stream as never, "/test", token);
    expect(stream.markdown).toHaveBeenCalledWith(expect.stringContaining("Error"));
  });
});

describe("handleGuide", () => {
  beforeEach(() => vi.clearAllMocks());

  it("runs guide without flags when no prompt", async () => {
    mockRunCtxSuccess("ctx cheat sheet");
    const stream = fakeStream();
    const token = fakeToken();
    await handleGuide(stream as never, "", "/test", token);
    expect(cp.execFile).toHaveBeenCalledWith(
      "ctx",
      ["guide", "--no-color"],
      expect.anything(),
      expect.any(Function)
    );
  });

  it("passes --skills flag", async () => {
    mockRunCtxSuccess("skills list");
    const stream = fakeStream();
    const token = fakeToken();
    await handleGuide(stream as never, "skills", "/test", token);
    expect(cp.execFile).toHaveBeenCalledWith(
      "ctx",
      ["guide", "--no-color", "--skills"],
      expect.anything(),
      expect.any(Function)
    );
  });

  it("passes --commands flag", async () => {
    mockRunCtxSuccess("commands list");
    const stream = fakeStream();
    const token = fakeToken();
    await handleGuide(stream as never, "commands", "/test", token);
    expect(cp.execFile).toHaveBeenCalledWith(
      "ctx",
      ["guide", "--no-color", "--commands"],
      expect.anything(),
      expect.any(Function)
    );
  });

  it("handles errors gracefully", async () => {
    mockRunCtxError("guide error");
    const stream = fakeStream();
    const token = fakeToken();
    await handleGuide(stream as never, "", "/test", token);
    expect(stream.markdown).toHaveBeenCalledWith(expect.stringContaining("Error"));
  });
});

describe("handleWhy", () => {
  beforeEach(() => vi.clearAllMocks());

  it("runs why with document name", async () => {
    mockRunCtxSuccess("The ctx Manifesto");
    const stream = fakeStream();
    const token = fakeToken();
    await handleWhy(stream as never, "manifesto", "/test", token);
    expect(cp.execFile).toHaveBeenCalledWith(
      "ctx",
      ["why", "manifesto", "--no-color"],
      expect.anything(),
      expect.any(Function)
    );
  });

  it("runs why without document for interactive menu", async () => {
    mockRunCtxSuccess("1. manifesto\n2. about");
    const stream = fakeStream();
    const token = fakeToken();
    await handleWhy(stream as never, "", "/test", token);
    expect(cp.execFile).toHaveBeenCalledWith(
      "ctx",
      ["why", "--no-color"],
      expect.anything(),
      expect.any(Function)
    );
  });

  it("handles errors gracefully", async () => {
    mockRunCtxError("doc not found");
    const stream = fakeStream();
    const token = fakeToken();
    await handleWhy(stream as never, "manifesto", "/test", token);
    expect(stream.markdown).toHaveBeenCalledWith(expect.stringContaining("Error"));
  });
});

describe("handleMemory", () => {
  beforeEach(() => vi.clearAllMocks());

  it("shows usage when no subcommand given", async () => {
    const stream = fakeStream();
    const token = fakeToken();
    const result = await handleMemory(stream as never, "", "/test", token);
    expect(result.metadata.command).toBe("memory");
    expect(stream.markdown).toHaveBeenCalledWith(expect.stringContaining("Usage"));
  });

  it("runs sync subcommand", async () => {
    mockRunCtxSuccess("Memory synced");
    const stream = fakeStream();
    const token = fakeToken();
    await handleMemory(stream as never, "sync", "/test", token);
    expect(cp.execFile).toHaveBeenCalledWith(
      "ctx",
      ["memory", "sync", "--no-color"],
      expect.anything(),
      expect.any(Function)
    );
  });

  it("runs status subcommand", async () => {
    mockRunCtxSuccess("Status: in sync");
    const stream = fakeStream();
    const token = fakeToken();
    await handleMemory(stream as never, "status", "/test", token);
    expect(cp.execFile).toHaveBeenCalledWith(
      "ctx",
      ["memory", "status", "--no-color"],
      expect.anything(),
      expect.any(Function)
    );
  });

  it("handles errors gracefully", async () => {
    mockRunCtxError("memory error");
    const stream = fakeStream();
    const token = fakeToken();
    await handleMemory(stream as never, "sync", "/test", token);
    expect(stream.markdown).toHaveBeenCalledWith(expect.stringContaining("Error"));
  });
});

describe("handlePrompt", () => {
  beforeEach(() => vi.clearAllMocks());

  it("lists prompts when no subcommand given", async () => {
    mockRunCtxSuccess("review\nrefactor");
    const stream = fakeStream();
    const token = fakeToken();
    await handlePrompt(stream as never, "", "/test", token);
    expect(cp.execFile).toHaveBeenCalledWith(
      "ctx",
      ["prompt", "list", "--no-color"],
      expect.anything(),
      expect.any(Function)
    );
  });

  it("shows a prompt by name", async () => {
    mockRunCtxSuccess("Review the code...");
    const stream = fakeStream();
    const token = fakeToken();
    await handlePrompt(stream as never, "show review", "/test", token);
    expect(cp.execFile).toHaveBeenCalledWith(
      "ctx",
      ["prompt", "show", "review", "--no-color"],
      expect.anything(),
      expect.any(Function)
    );
  });

  it("shows usage when 'rm' has no name", async () => {
    const stream = fakeStream();
    const token = fakeToken();
    await handlePrompt(stream as never, "rm", "/test", token);
    expect(stream.markdown).toHaveBeenCalledWith(expect.stringContaining("Usage"));
  });

  it("shows 'No prompt templates found.' when empty", async () => {
    mockRunCtxSuccess("");
    const stream = fakeStream();
    const token = fakeToken();
    await handlePrompt(stream as never, "", "/test", token);
    expect(stream.markdown).toHaveBeenCalledWith("No prompt templates found.");
  });

  it("handles errors gracefully", async () => {
    mockRunCtxError("prompt error");
    const stream = fakeStream();
    const token = fakeToken();
    await handlePrompt(stream as never, "show missing", "/test", token);
    expect(stream.markdown).toHaveBeenCalledWith(expect.stringContaining("Error"));
  });
});

describe("handleDecisions", () => {
  beforeEach(() => vi.clearAllMocks());

  it("runs decisions command", async () => {
    mockRunCtxSuccess("1. Use Redis");
    const stream = fakeStream();
    const token = fakeToken();
    await handleDecisions(stream as never, "", "/test", token);
    expect(cp.execFile).toHaveBeenCalledWith(
      "ctx",
      ["decisions", "--no-color"],
      expect.anything(),
      expect.any(Function)
    );
  });

  it("passes subcommand arguments", async () => {
    mockRunCtxSuccess("Decision listed");
    const stream = fakeStream();
    const token = fakeToken();
    await handleDecisions(stream as never, "list", "/test", token);
    expect(cp.execFile).toHaveBeenCalledWith(
      "ctx",
      ["decisions", "list", "--no-color"],
      expect.anything(),
      expect.any(Function)
    );
  });

  it("handles errors gracefully", async () => {
    mockRunCtxError("decisions error");
    const stream = fakeStream();
    const token = fakeToken();
    await handleDecisions(stream as never, "", "/test", token);
    expect(stream.markdown).toHaveBeenCalledWith(expect.stringContaining("Error"));
  });
});

describe("handleLearnings", () => {
  beforeEach(() => vi.clearAllMocks());

  it("runs learnings command", async () => {
    mockRunCtxSuccess("1. Go embed trick");
    const stream = fakeStream();
    const token = fakeToken();
    await handleLearnings(stream as never, "", "/test", token);
    expect(cp.execFile).toHaveBeenCalledWith(
      "ctx",
      ["learnings", "--no-color"],
      expect.anything(),
      expect.any(Function)
    );
  });

  it("handles errors gracefully", async () => {
    mockRunCtxError("learnings error");
    const stream = fakeStream();
    const token = fakeToken();
    await handleLearnings(stream as never, "", "/test", token);
    expect(stream.markdown).toHaveBeenCalledWith(expect.stringContaining("Error"));
  });
});

describe("handleDeps", () => {
  beforeEach(() => vi.clearAllMocks());

  it("runs deps command", async () => {
    mockRunCtxSuccess("internal/mcp -> internal/config");
    const stream = fakeStream();
    const token = fakeToken();
    await handleDeps(stream as never, "", "/test", token);
    expect(cp.execFile).toHaveBeenCalledWith(
      "ctx",
      ["deps", "--no-color"],
      expect.anything(),
      expect.any(Function)
    );
    expect(stream.progress).toHaveBeenCalledWith("Analyzing dependencies...");
  });

  it("shows fallback when no output", async () => {
    mockRunCtxSuccess("");
    const stream = fakeStream();
    const token = fakeToken();
    await handleDeps(stream as never, "", "/test", token);
    expect(stream.markdown).toHaveBeenCalledWith("No dependency information available.");
  });

  it("handles errors gracefully", async () => {
    mockRunCtxError("deps error");
    const stream = fakeStream();
    const token = fakeToken();
    await handleDeps(stream as never, "", "/test", token);
    expect(stream.markdown).toHaveBeenCalledWith(expect.stringContaining("Error"));
  });
});

describe("handleJournal", () => {
  beforeEach(() => vi.clearAllMocks());

  it("runs journal command", async () => {
    mockRunCtxSuccess("Session analysis");
    const stream = fakeStream();
    const token = fakeToken();
    await handleJournal(stream as never, "", "/test", token);
    expect(cp.execFile).toHaveBeenCalledWith(
      "ctx",
      ["journal", "--no-color"],
      expect.anything(),
      expect.any(Function)
    );
  });

  it("passes subcommand arguments", async () => {
    mockRunCtxSuccess("Analysis complete");
    const stream = fakeStream();
    const token = fakeToken();
    await handleJournal(stream as never, "synthesize", "/test", token);
    expect(cp.execFile).toHaveBeenCalledWith(
      "ctx",
      ["journal", "synthesize", "--no-color"],
      expect.anything(),
      expect.any(Function)
    );
  });

  it("handles errors gracefully", async () => {
    mockRunCtxError("journal error");
    const stream = fakeStream();
    const token = fakeToken();
    await handleJournal(stream as never, "", "/test", token);
    expect(stream.markdown).toHaveBeenCalledWith(expect.stringContaining("Error"));
  });
});

describe("handleReindex", () => {
  beforeEach(() => vi.clearAllMocks());

  it("runs reindex command", async () => {
    mockRunCtxSuccess("Indices updated");
    const stream = fakeStream();
    const token = fakeToken();
    await handleReindex(stream as never, "/test", token);
    expect(cp.execFile).toHaveBeenCalledWith(
      "ctx",
      ["reindex", "--no-color"],
      expect.anything(),
      expect.any(Function)
    );
    expect(stream.progress).toHaveBeenCalledWith("Regenerating indices...");
  });

  it("shows fallback when no output", async () => {
    mockRunCtxSuccess("");
    const stream = fakeStream();
    const token = fakeToken();
    await handleReindex(stream as never, "/test", token);
    expect(stream.markdown).toHaveBeenCalledWith("Indices regenerated.");
  });

  it("handles errors gracefully", async () => {
    mockRunCtxError("reindex error");
    const stream = fakeStream();
    const token = fakeToken();
    await handleReindex(stream as never, "/test", token);
    expect(stream.markdown).toHaveBeenCalledWith(expect.stringContaining("Error"));
  });
});

describe("splitArgs", () => {
  it("splits simple space-separated args", () => {
    expect(splitArgs("task Fix login bug")).toEqual(["task", "Fix", "login", "bug"]);
  });

  it("handles double-quoted strings", () => {
    expect(splitArgs('decision "Use PostgreSQL" --context "Need DB"')).toEqual([
      "decision", "Use PostgreSQL", "--context", "Need DB",
    ]);
  });

  it("returns empty array for empty input", () => {
    expect(splitArgs("")).toEqual([]);
  });

  it("handles single arg", () => {
    expect(splitArgs("task")).toEqual(["task"]);
  });

  it("handles mixed quoted and unquoted args", () => {
    expect(splitArgs('decision "Use Redis" --rationale ACID')).toEqual([
      "decision", "Use Redis", "--rationale", "ACID",
    ]);
  });
});

describe("handleAdd", () => {
  beforeEach(() => vi.clearAllMocks());

  it("shows usage when no type provided", async () => {
    const stream = fakeStream();
    const token = fakeToken();
    const result = await handleAdd(stream as never, "", "/test", token);
    expect(result.metadata.command).toBe("add");
    expect(stream.markdown).toHaveBeenCalledWith(expect.stringContaining("Usage"));
  });

  it("adds a task", async () => {
    mockRunCtxSuccess("Task added");
    const stream = fakeStream();
    const token = fakeToken();
    await handleAdd(stream as never, "task Fix login bug", "/test", token);
    expect(cp.execFile).toHaveBeenCalledWith(
      "ctx",
      ["add", "task", "Fix", "login", "bug", "--no-color"],
      expect.anything(),
      expect.any(Function)
    );
  });

  it("adds a decision with flags", async () => {
    mockRunCtxSuccess("Decision added");
    const stream = fakeStream();
    const token = fakeToken();
    await handleAdd(
      stream as never,
      'decision "Use PostgreSQL" --context "Need DB" --rationale "ACID" --consequences "Training"',
      "/test",
      token
    );
    expect(cp.execFile).toHaveBeenCalledWith(
      "ctx",
      ["add", "decision", "Use PostgreSQL", "--context", "Need DB", "--rationale", "ACID", "--consequences", "Training", "--no-color"],
      expect.anything(),
      expect.any(Function)
    );
  });

  it("shows fallback message when output is empty", async () => {
    mockRunCtxSuccess("");
    const stream = fakeStream();
    const token = fakeToken();
    await handleAdd(stream as never, "task Test", "/test", token);
    expect(stream.markdown).toHaveBeenCalledWith("Added **task** entry.");
  });

  it("handles errors gracefully", async () => {
    mockRunCtxError("add failed");
    const stream = fakeStream();
    const token = fakeToken();
    await handleAdd(stream as never, "task Test", "/test", token);
    expect(stream.markdown).toHaveBeenCalledWith(expect.stringContaining("Error"));
  });
});

describe("handleAgent", () => {
  beforeEach(() => vi.clearAllMocks());

  it("runs agent command", async () => {
    mockRunCtxSuccess("Context packet...");
    const stream = fakeStream();
    const token = fakeToken();
    await handleAgent(stream as never, "", "/test", token);
    expect(cp.execFile).toHaveBeenCalledWith(
      "ctx",
      ["agent", "--no-color"],
      expect.anything(),
      expect.any(Function)
    );
  });

  it("passes --budget flag", async () => {
    mockRunCtxSuccess("Context packet...");
    const stream = fakeStream();
    const token = fakeToken();
    await handleAgent(stream as never, "--budget 4000", "/test", token);
    expect(cp.execFile).toHaveBeenCalledWith(
      "ctx",
      ["agent", "--no-color", "--budget", "4000"],
      expect.anything(),
      expect.any(Function)
    );
  });

  it("passes --format flag", async () => {
    mockRunCtxSuccess("{}");
    const stream = fakeStream();
    const token = fakeToken();
    await handleAgent(stream as never, "--format json", "/test", token);
    expect(cp.execFile).toHaveBeenCalledWith(
      "ctx",
      ["agent", "--no-color", "--format", "json"],
      expect.anything(),
      expect.any(Function)
    );
  });

  it("handles errors gracefully", async () => {
    mockRunCtxError("agent error");
    const stream = fakeStream();
    const token = fakeToken();
    await handleAgent(stream as never, "", "/test", token);
    expect(stream.markdown).toHaveBeenCalledWith(expect.stringContaining("Error"));
  });
});

describe("handleLoad", () => {
  beforeEach(() => vi.clearAllMocks());

  it("runs load command", async () => {
    mockRunCtxSuccess("assembled context");
    const stream = fakeStream();
    const token = fakeToken();
    await handleLoad(stream as never, "", "/test", token);
    expect(cp.execFile).toHaveBeenCalledWith(
      "ctx",
      ["load", "--no-color"],
      expect.anything(),
      expect.any(Function)
    );
  });

  it("passes --budget flag", async () => {
    mockRunCtxSuccess("trimmed context");
    const stream = fakeStream();
    const token = fakeToken();
    await handleLoad(stream as never, "--budget 2000", "/test", token);
    expect(cp.execFile).toHaveBeenCalledWith(
      "ctx",
      ["load", "--no-color", "--budget", "2000"],
      expect.anything(),
      expect.any(Function)
    );
  });

  it("passes --raw flag", async () => {
    mockRunCtxSuccess("raw context");
    const stream = fakeStream();
    const token = fakeToken();
    await handleLoad(stream as never, "--raw", "/test", token);
    expect(cp.execFile).toHaveBeenCalledWith(
      "ctx",
      ["load", "--no-color", "--raw"],
      expect.anything(),
      expect.any(Function)
    );
  });

  it("handles errors gracefully", async () => {
    mockRunCtxError("load error");
    const stream = fakeStream();
    const token = fakeToken();
    await handleLoad(stream as never, "", "/test", token);
    expect(stream.markdown).toHaveBeenCalledWith(expect.stringContaining("Error"));
  });
});

describe("handleCompact", () => {
  beforeEach(() => vi.clearAllMocks());

  it("runs compact command", async () => {
    mockRunCtxSuccess("Compacted");
    const stream = fakeStream();
    const token = fakeToken();
    await handleCompact(stream as never, "", "/test", token);
    expect(cp.execFile).toHaveBeenCalledWith(
      "ctx",
      ["compact", "--no-color"],
      expect.anything(),
      expect.any(Function)
    );
  });

  it("passes --archive flag when keyword used", async () => {
    mockRunCtxSuccess("Archived");
    const stream = fakeStream();
    const token = fakeToken();
    await handleCompact(stream as never, "archive", "/test", token);
    expect(cp.execFile).toHaveBeenCalledWith(
      "ctx",
      ["compact", "--no-color", "--archive"],
      expect.anything(),
      expect.any(Function)
    );
  });

  it("shows fallback when output is empty", async () => {
    mockRunCtxSuccess("");
    const stream = fakeStream();
    const token = fakeToken();
    await handleCompact(stream as never, "", "/test", token);
    expect(stream.markdown).toHaveBeenCalledWith("Context compacted successfully.");
  });

  it("handles errors gracefully", async () => {
    mockRunCtxError("compact error");
    const stream = fakeStream();
    const token = fakeToken();
    await handleCompact(stream as never, "", "/test", token);
    expect(stream.markdown).toHaveBeenCalledWith(expect.stringContaining("Error"));
  });
});

describe("handleSync", () => {
  beforeEach(() => vi.clearAllMocks());

  it("runs sync command", async () => {
    mockRunCtxSuccess("Synced");
    const stream = fakeStream();
    const token = fakeToken();
    await handleSync(stream as never, "", "/test", token);
    expect(cp.execFile).toHaveBeenCalledWith(
      "ctx",
      ["sync", "--no-color"],
      expect.anything(),
      expect.any(Function)
    );
  });

  it("passes --dry-run flag when keyword used", async () => {
    mockRunCtxSuccess("Would sync...");
    const stream = fakeStream();
    const token = fakeToken();
    await handleSync(stream as never, "dry-run", "/test", token);
    expect(cp.execFile).toHaveBeenCalledWith(
      "ctx",
      ["sync", "--no-color", "--dry-run"],
      expect.anything(),
      expect.any(Function)
    );
  });

  it("shows fallback when output is empty", async () => {
    mockRunCtxSuccess("");
    const stream = fakeStream();
    const token = fakeToken();
    await handleSync(stream as never, "", "/test", token);
    expect(stream.markdown).toHaveBeenCalledWith("Context synced with codebase.");
  });

  it("handles errors gracefully", async () => {
    mockRunCtxError("sync error");
    const stream = fakeStream();
    const token = fakeToken();
    await handleSync(stream as never, "", "/test", token);
    expect(stream.markdown).toHaveBeenCalledWith(expect.stringContaining("Error"));
  });
});

describe("handleRecall", () => {
  beforeEach(() => vi.clearAllMocks());

  it("defaults to list when no subcommand", async () => {
    mockRunCtxSuccess("session 1\nsession 2");
    const stream = fakeStream();
    const token = fakeToken();
    await handleRecall(stream as never, "", "/test", token);
    expect(cp.execFile).toHaveBeenCalledWith(
      "ctx",
      ["recall", "list", "--no-color"],
      expect.anything(),
      expect.any(Function)
    );
  });

  it("routes show subcommand", async () => {
    mockRunCtxSuccess("session details");
    const stream = fakeStream();
    const token = fakeToken();
    await handleRecall(stream as never, "show abc123", "/test", token);
    expect(cp.execFile).toHaveBeenCalledWith(
      "ctx",
      ["recall", "show", "abc123", "--no-color"],
      expect.anything(),
      expect.any(Function)
    );
  });

  it("routes export subcommand", async () => {
    mockRunCtxSuccess("exported");
    const stream = fakeStream();
    const token = fakeToken();
    await handleRecall(stream as never, "export abc123", "/test", token);
    expect(cp.execFile).toHaveBeenCalledWith(
      "ctx",
      ["recall", "export", "--all", "abc123", "--no-color"],
      expect.anything(),
      expect.any(Function)
    );
  });

  it("routes lock subcommand", async () => {
    mockRunCtxSuccess("locked");
    const stream = fakeStream();
    const token = fakeToken();
    await handleRecall(stream as never, "lock abc123", "/test", token);
    expect(cp.execFile).toHaveBeenCalledWith(
      "ctx",
      ["recall", "lock", "abc123", "--no-color"],
      expect.anything(),
      expect.any(Function)
    );
  });

  it("routes unlock subcommand with --all when no id", async () => {
    mockRunCtxSuccess("unlocked");
    const stream = fakeStream();
    const token = fakeToken();
    await handleRecall(stream as never, "unlock", "/test", token);
    expect(cp.execFile).toHaveBeenCalledWith(
      "ctx",
      ["recall", "unlock", "--all", "--no-color"],
      expect.anything(),
      expect.any(Function)
    );
  });

  it("routes sync subcommand", async () => {
    mockRunCtxSuccess("synced");
    const stream = fakeStream();
    const token = fakeToken();
    await handleRecall(stream as never, "sync", "/test", token);
    expect(cp.execFile).toHaveBeenCalledWith(
      "ctx",
      ["recall", "sync", "--no-color"],
      expect.anything(),
      expect.any(Function)
    );
  });

  it("treats unknown text as query for list", async () => {
    mockRunCtxSuccess("matching sessions");
    const stream = fakeStream();
    const token = fakeToken();
    await handleRecall(stream as never, "refactoring work", "/test", token);
    expect(cp.execFile).toHaveBeenCalledWith(
      "ctx",
      ["recall", "list", "--query", "refactoring work", "--no-color"],
      expect.anything(),
      expect.any(Function)
    );
  });

  it("shows 'No session history found.' when empty", async () => {
    mockRunCtxSuccess("");
    const stream = fakeStream();
    const token = fakeToken();
    await handleRecall(stream as never, "", "/test", token);
    expect(stream.markdown).toHaveBeenCalledWith("No session history found.");
  });

  it("handles errors gracefully", async () => {
    mockRunCtxError("recall error");
    const stream = fakeStream();
    const token = fakeToken();
    await handleRecall(stream as never, "", "/test", token);
    expect(stream.markdown).toHaveBeenCalledWith(expect.stringContaining("Error"));
  });
});

describe("handleDeps with flags", () => {
  beforeEach(() => vi.clearAllMocks());

  it("passes --format flag", async () => {
    mockRunCtxSuccess("json output");
    const stream = fakeStream();
    const token = fakeToken();
    await handleDeps(stream as never, "--format json", "/test", token);
    expect(cp.execFile).toHaveBeenCalledWith(
      "ctx",
      ["deps", "--no-color", "--format", "json"],
      expect.anything(),
      expect.any(Function)
    );
  });

  it("passes --external flag", async () => {
    mockRunCtxSuccess("external deps");
    const stream = fakeStream();
    const token = fakeToken();
    await handleDeps(stream as never, "--external", "/test", token);
    expect(cp.execFile).toHaveBeenCalledWith(
      "ctx",
      ["deps", "--no-color", "--external"],
      expect.anything(),
      expect.any(Function)
    );
  });

  it("passes --type flag", async () => {
    mockRunCtxSuccess("go deps");
    const stream = fakeStream();
    const token = fakeToken();
    await handleDeps(stream as never, "--type go", "/test", token);
    expect(cp.execFile).toHaveBeenCalledWith(
      "ctx",
      ["deps", "--no-color", "--type", "go"],
      expect.anything(),
      expect.any(Function)
    );
  });
});
