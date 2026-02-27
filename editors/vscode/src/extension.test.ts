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
