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
  handleComplete,
  handleRemind,
  handleTasks,
  handlePad,
  handleNotify,
  handleSystem,
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
