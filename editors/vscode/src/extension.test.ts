import { describe, it, expect, vi, beforeEach } from "vitest";
import * as cp from "child_process";
import * as vscode from "vscode";

vi.mock("vscode", async () => (await import("./vscodeMock")).createVscodeMock());
vi.mock("child_process");

import type { createVscodeMock } from "./vscodeMock";
// The mocked module's classes, with constructors the real typings hide.
const vs = vscode as unknown as ReturnType<typeof createVscodeMock>;

import {
  runCtx,
  getCtxPath,
  getWorkspaceRoot,
  getPlatformInfo,
  handler,
  tokenize,
  CLI_COMMANDS,
  SKILLS,
} from "./extension";

type ExecCallback = (e: unknown, out: string, err: string) => void;

// Helper: create a fake CancellationToken
function fakeToken(cancelled = false) {
  type Listener = (e: unknown) => unknown;
  const listeners: Listener[] = [];
  return {
    isCancellationRequested: cancelled,
    onCancellationRequested: vi.fn((cb: Listener) => {
      listeners.push(cb);
      return { dispose: vi.fn() };
    }),
    _fire: () => listeners.forEach((cb) => cb(undefined)),
  };
}

function fakeStream() {
  return {
    markdown: vi.fn(),
    progress: vi.fn(),
  };
}

/** execFile error for a process that exited with `code`. */
function exitError(code: number) {
  return Object.assign(new Error(`exit ${code}`), { code });
}

/** Every ctx call exits with `code` and prints `stdout`; git answers rev-parse. */
function mockExec(stdout: string, code = 0, stderr = "") {
  vi.mocked(cp.execFile).mockImplementation(((
    cmd: string,
    args: string[],
    _opts: unknown,
    cb: ExecCallback
  ) => {
    if (cmd === "git") {
      cb(null, args.includes("--abbrev-ref") ? "main\n" : "abc1234\n", "");
    } else {
      cb(code === 0 ? null : exitError(code), stdout, stderr);
    }
    return { kill: vi.fn() };
  }) as never);
}

/** argv of every ctx (non-git) call so far. */
function ctxCalls(): string[][] {
  return vi
    .mocked(cp.execFile)
    .mock.calls.filter((c) => c[0] !== "git")
    .map((c) => c[1] as unknown as string[]);
}

function markdownOf(stream: ReturnType<typeof fakeStream>): string {
  return stream.markdown.mock.calls.map((c) => c[0]).join("\n");
}

async function run(command: string, prompt: string) {
  const stream = fakeStream();
  const res = await CLI_COMMANDS[command](stream as never, prompt, "/test", fakeToken() as never);
  return { stream, res };
}

describe("getCtxPath", () => {
  it("returns 'ctx' when no config is set", () => {
    expect(getCtxPath()).toBe("ctx");
  });

  it("returns configured path when set", () => {
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

  it("prefers the folder of the active editor in a multi-root window", () => {
    const win = vscode.window as { activeTextEditor: unknown };
    win.activeTextEditor = { document: { uri: { fsPath: "/other/file.ts" } } };
    vi.mocked(vscode.workspace.getWorkspaceFolder).mockReturnValueOnce({
      uri: { fsPath: "/other" },
    } as never);
    expect(getWorkspaceRoot()).toBe("/other");
    win.activeTextEditor = undefined;
  });

  it("returns undefined when no workspace is open", () => {
    const ws = vscode.workspace as { workspaceFolders: unknown };
    const original = ws.workspaceFolders;
    ws.workspaceFolders = undefined;
    expect(getWorkspaceRoot()).toBeUndefined();
    ws.workspaceFolders = original;
  });
});

describe("runCtx", () => {
  beforeEach(() => vi.clearAllMocks());

  it("resolves with output and exit code 0 on success", async () => {
    mockExec("output", 0, "errors");
    const result = await runCtx(["status"]);
    expect(result).toEqual({ stdout: "output", stderr: "errors", code: 0 });
  });

  it("resolves with the real exit code on a non-zero exit", async () => {
    mockExec("drift report", 1);
    const result = await runCtx(["drift"]);
    expect(result.code).toBe(1);
    expect(result.stdout).toBe("drift report");
  });

  it("resolves a non-zero exit even without output", async () => {
    mockExec("", 2);
    expect((await runCtx(["status"])).code).toBe(2);
  });

  it("rejects when the binary cannot start", async () => {
    vi.mocked(cp.execFile).mockImplementation(((_c: unknown, _a: unknown, _o: unknown, cb: ExecCallback) => {
      cb(Object.assign(new Error("spawn ctx ENOENT"), { code: "ENOENT" }), "", "");
      return { kill: vi.fn() };
    }) as never);
    await expect(runCtx(["status"])).rejects.toThrow("ENOENT");
  });

  it("rejects on timeout instead of showing partial output", async () => {
    vi.mocked(cp.execFile).mockImplementation(((_c: unknown, _a: unknown, _o: unknown, cb: ExecCallback) => {
      cb(Object.assign(new Error("timeout"), { killed: true, signal: "SIGTERM" }), "partial", "");
      return { kill: vi.fn() };
    }) as never);
    await expect(runCtx(["agent"])).rejects.toThrow("cancelled or timed out");
  });

  it("rejects immediately when token is already cancelled", async () => {
    await expect(runCtx(["status"], "/test", fakeToken(true) as never)).rejects.toThrow("Cancelled");
    expect(cp.execFile).not.toHaveBeenCalled();
  });

  it("kills the child and rejects when the token fires", async () => {
    const killFn = vi.fn();
    let finish: ExecCallback = () => {};
    vi.mocked(cp.execFile).mockImplementation(((_c: unknown, _a: unknown, _o: unknown, cb: ExecCallback) => {
      finish = cb;
      return { kill: killFn };
    }) as never);

    const token = fakeToken();
    const promise = runCtx(["agent"], "/test", token as never);
    token._fire();
    expect(killFn).toHaveBeenCalled();

    finish(Object.assign(new Error("killed"), { killed: true, signal: "SIGTERM" }), "half", "");
    await expect(promise).rejects.toThrow("cancelled or timed out");
  });

  it("passes cwd and never runs through a shell", async () => {
    mockExec("");
    await runCtx(["status"], "/my/project");
    const opts = vi.mocked(cp.execFile).mock.calls[0][2] as { cwd: string; shell?: unknown };
    expect(opts.cwd).toBe("/my/project");
    expect(opts.shell).toBeUndefined();
  });

  it("closes stdin so a prompting command cannot wait for input", async () => {
    const end = vi.fn();
    vi.mocked(cp.execFile).mockImplementation(((_c: unknown, _a: unknown, _o: unknown, cb: ExecCallback) => {
      process.nextTick(() => cb(null, "", ""));
      return { kill: vi.fn(), stdin: { end } };
    }) as never);
    await runCtx(["why"]);
    expect(end).toHaveBeenCalled();
  });

  it("disposes cancellation listener when process completes", async () => {
    const disposeFn = vi.fn();
    const token = {
      isCancellationRequested: false,
      onCancellationRequested: vi.fn(() => ({ dispose: disposeFn })),
    };
    vi.mocked(cp.execFile).mockImplementation(((_c: unknown, _a: unknown, _o: unknown, cb: ExecCallback) => {
      process.nextTick(() => cb(null, "done", ""));
      return { kill: vi.fn() };
    }) as never);

    await runCtx(["status"], "/test", token as never);
    expect(disposeFn).toHaveBeenCalled();
  });
});

describe("getPlatformInfo", () => {
  it("returns valid goos, goarch, and extension", () => {
    const info = getPlatformInfo();
    expect(["darwin", "linux", "windows"]).toContain(info.goos);
    expect(["amd64", "arm64"]).toContain(info.goarch);
    expect(info.ext).toBe(info.goos === "windows" ? ".exe" : "");
  });
});

describe("tokenize", () => {
  it("keeps double-quoted phrases together", () => {
    expect(tokenize('decision Use Postgres --context "need a db" --x y')).toEqual([
      "decision",
      "Use",
      "Postgres",
      "--context",
      "need a db",
      "--x",
      "y",
    ]);
  });
});

describe("result rendering", () => {
  beforeEach(() => vi.clearAllMocks());

  it("fences successful output", async () => {
    mockExec("3 tasks archived");
    const { stream } = await run("task", "archive");
    expect(markdownOf(stream)).toBe("```\n3 tasks archived\n```");
  });

  it("reports a non-zero exit as such, never as a result", async () => {
    mockExec("Error: unknown flag: --no-color", 1);
    const { stream } = await run("status", "");
    const md = markdownOf(stream);
    expect(md).toContain("`ctx status` exited with code 1.");
    expect(md).toContain("unknown flag");
  });

  it("points at /init when the folder has no .context/", async () => {
    mockExec("Error: no .context here", 1);
    const { stream } = await run("status", "");
    expect(markdownOf(stream)).toContain("@ctx /init");
  });

  it("renders spawn failures as errors", async () => {
    vi.mocked(cp.execFile).mockImplementation(((_c: unknown, _a: unknown, _o: unknown, cb: ExecCallback) => {
      cb(Object.assign(new Error("spawn ctx ENOENT"), { code: "ENOENT" }), "", "");
      return { kill: vi.fn() };
    }) as never);
    const { stream } = await run("drift", "");
    expect(markdownOf(stream)).toContain("**Error:**");
  });
});

describe("/task", () => {
  beforeEach(() => vi.clearAllMocks());

  it("shows usage when no subcommand given", async () => {
    mockExec("");
    const { stream, res } = await run("task", "");
    expect(res.metadata.command).toBe("task");
    expect(markdownOf(stream)).toContain("Usage");
    expect(cp.execFile).not.toHaveBeenCalled();
  });

  it("shows usage for complete without a reference", async () => {
    const { stream } = await run("task", "complete");
    expect(markdownOf(stream)).toContain("Usage");
  });

  it.each([
    ["complete Fix login bug", ["task", "complete", "Fix login bug"]],
    ["archive", ["task", "archive"]],
    ["snapshot pre-refactor", ["task", "snapshot", "pre-refactor"]],
    ["snapshot", ["task", "snapshot"]],
  ])("%s", async (prompt, argv) => {
    mockExec("ok");
    await run("task", prompt);
    expect(ctxCalls()).toEqual([argv]);
  });
});

describe("/remind", () => {
  beforeEach(() => vi.clearAllMocks());

  it.each([
    ["", ["remind", "list"]],
    ["list", ["remind", "list"]],
    ["add Check CI status", ["remind", "add", "Check CI status"]],
    ["Check CI status", ["remind", "add", "Check CI status"]],
    ["dismiss 2", ["remind", "dismiss", "2"]],
    ["dismiss", ["remind", "dismiss", "--all"]],
  ])("'%s'", async (prompt, argv) => {
    mockExec("ok");
    await run("remind", prompt);
    expect(ctxCalls()).toEqual([argv]);
  });

  it("shows 'No reminders.' when output is empty", async () => {
    mockExec("");
    const { stream } = await run("remind", "list");
    expect(markdownOf(stream)).toBe("No reminders.");
  });
});

describe("/pad", () => {
  beforeEach(() => vi.clearAllMocks());

  it.each([
    ["", ["pad"]],
    ["add my secret note", ["pad", "add", "my secret note"]],
    ["show 1", ["pad", "show", "1"]],
    ["rm 2 3", ["pad", "rm", "2", "3"]],
    // `ctx pad edit N [TEXT]`: the text must stay one argument
    ["edit 1 new text", ["pad", "edit", "1", "new text"]],
    ["mv 1 3", ["pad", "mv", "1", "3"]],
  ])("'%s'", async (prompt, argv) => {
    mockExec("ok");
    await run("pad", prompt);
    expect(ctxCalls()).toEqual([argv]);
  });

  it.each(["add", "rm", "edit", "import"])("shows usage for '%s' without arguments", async (sub) => {
    const { stream } = await run("pad", sub);
    expect(markdownOf(stream)).toContain("Usage");
    expect(cp.execFile).not.toHaveBeenCalled();
  });

  it("shows 'Scratchpad is empty.' when output is empty", async () => {
    mockExec("");
    const { stream } = await run("pad", "");
    expect(markdownOf(stream)).toBe("Scratchpad is empty.");
  });
});

describe("/notify", () => {
  beforeEach(() => vi.clearAllMocks());

  it("shows usage when no message given", async () => {
    const { stream } = await run("notify", "");
    expect(markdownOf(stream)).toContain("Usage");
  });

  it("sends setup to the terminal, so the webhook URL never enters the chat", async () => {
    const { stream } = await run("notify", "setup");
    expect(markdownOf(stream)).toContain("ctx hook notify setup");
    expect(cp.execFile).not.toHaveBeenCalled();
  });

  it.each([
    ["test", ["hook", "notify", "test"]],
    // `ctx hook notify [message]`: the message must stay one argument
    ["build done --event build", ["hook", "notify", "build done", "--event", "build"]],
  ])("'%s'", async (prompt, argv) => {
    mockExec("ok");
    await run("notify", prompt);
    expect(ctxCalls()).toEqual([argv]);
  });
});

describe("/system", () => {
  beforeEach(() => vi.clearAllMocks());

  it("shows usage when no subcommand given", async () => {
    const { stream } = await run("system", "");
    expect(markdownOf(stream)).toContain("Usage");
  });

  it.each([
    ["resources", ["sysinfo"]],
    ["bootstrap", ["system", "bootstrap"]],
    ["stats", ["usage"]],
    ["message", ["hook", "message", "list"]],
    ["message show check-freshness stale", ["hook", "message", "show", "check-freshness", "stale"]],
  ])("'%s'", async (prompt, argv) => {
    mockExec("ok");
    await run("system", prompt);
    expect(ctxCalls()).toEqual([argv]);
  });
});

describe("/why", () => {
  beforeEach(() => vi.clearAllMocks());

  it("defaults to the manifesto instead of the interactive menu", async () => {
    mockExec("# The ctx Manifesto");
    await run("why", "");
    expect(ctxCalls()).toEqual([["why", "manifesto"]]);
  });
});

describe("/add", () => {
  beforeEach(() => vi.clearAllMocks());

  it("shows usage for an unknown type", async () => {
    const { stream } = await run("add", "idea something");
    expect(markdownOf(stream)).toContain("Usage");
    expect(cp.execFile).not.toHaveBeenCalled();
  });

  it("adds a task with provenance", async () => {
    mockExec("✓ Added to TASKS.md");
    await run("add", 'task Fix login bug --section "Phase 1"');
    expect(ctxCalls()).toEqual([
      [
        "task", "add", "Fix login bug",
        "--section", "Phase 1",
        "--session-id", "01234567", "--branch", "main", "--commit", "abc1234",
      ],
    ]);
  });

  it("passes quoted flag values as single arguments and keeps user provenance", async () => {
    mockExec("✓ Added to DECISIONS.md");
    await run(
      "add",
      'decision Use PostgreSQL --context "Need a reliable DB" --rationale ACID --consequence "Ops training" --branch release'
    );
    expect(ctxCalls()).toEqual([
      [
        "decision", "add", "Use PostgreSQL",
        "--context", "Need a reliable DB", "--rationale", "ACID",
        "--consequence", "Ops training", "--branch", "release",
        "--session-id", "01234567", "--commit", "abc1234",
      ],
    ]);
  });

  it("adds a convention without provenance", async () => {
    mockExec("✓ Added to CONVENTIONS.md");
    await run("add", "convention Use camelCase --section Naming");
    expect(ctxCalls()).toEqual([["convention", "add", "Use camelCase", "--section", "Naming"]]);
  });

  it("surfaces the CLI's missing-field error", async () => {
    mockExec("Error: decision requires --context, --rationale, --consequence", 1);
    const { stream } = await run("add", "decision Use PostgreSQL");
    expect(markdownOf(stream)).toContain("exited with code 1");
  });
});

describe("skill-backed commands", () => {
  beforeEach(() => vi.clearAllMocks());

  function fakeModel(fragments: string[]) {
    return {
      sendRequest: vi.fn(async () => ({
        text: (async function* () {
          yield* fragments;
        })(),
      })),
    };
  }

  function request(command: string | undefined, prompt: string, model: unknown, references: unknown[] = []) {
    return { command, prompt, model, references } as never;
  }

  it("grounds the canonical skill in live ctx output and streams the answer", async () => {
    mockExec("# Context Packet");
    const model = fakeModel(["Recommended ", "next"]);
    const stream = fakeStream();
    const res = await handler(request("next", "", model), { history: [] } as never, stream as never, fakeToken() as never);

    expect(res).toEqual({ metadata: { command: "next" } });
    expect(ctxCalls()).toContainEqual(["agent"]);
    expect(ctxCalls()).toContainEqual(["journal", "source", "--limit", "3"]);
    const [messages] = model.sendRequest.mock.calls[0] as unknown as [Array<{ content: string }>];
    expect(messages[0].content).toContain(SKILLS.next.text);
    expect(messages[1].content).toContain("# Context Packet");
    expect(stream.markdown.mock.calls.map((c) => c[0])).toEqual(["Recommended ", "next"]);
  });

  it("inlines #file attachments", async () => {
    mockExec("packet");
    const model = fakeModel(["ok"]);
    const ref = { value: new vs.Uri("/test/workspace/specs/plans/m1.md") };
    await handler(request("implement", "", model, [ref]), { history: [] } as never, fakeStream() as never, fakeToken() as never);
    const [messages] = model.sendRequest.mock.calls[0] as unknown as [Array<{ content: string }>];
    expect(messages[1].content).toContain("attached text");
  });

  it("continues the previous skill when a reply has no slash command", async () => {
    mockExec("packet");
    const model = fakeModel(["next question"]);
    const history = [
      new vs.ChatRequestTurn("an auth idea", "brainstorm"),
      new vs.ChatResponseTurn(
        [new vs.ChatResponseMarkdownPart("What problem does it solve?")],
        { metadata: { command: "brainstorm" } }
      ),
    ];
    const res = await handler(
      request(undefined, "logins keep expiring, the status page is wrong", model),
      { history } as never,
      fakeStream() as never,
      fakeToken() as never
    );
    expect(res?.metadata?.command).toBe("brainstorm");
    const [messages] = model.sendRequest.mock.calls[0] as unknown as [Array<{ role: string; content: string }>];
    expect(messages[0].content).toContain(SKILLS.brainstorm.text);
    expect(messages).toContainEqual({ role: "user", content: "/brainstorm an auth idea" });
    expect(messages).toContainEqual({ role: "assistant", content: "What problem does it solve?" });
  });

  it("never sends CLI command turns (e.g. /pad) to the model", async () => {
    mockExec("packet");
    const model = fakeModel(["ok"]);
    const history = [
      new vs.ChatRequestTurn("add api-key=s3cret", "pad"),
      new vs.ChatResponseTurn([new vs.ChatResponseMarkdownPart("Added entry 1.")], {
        metadata: { command: "pad" },
      }),
    ];
    await handler(request("reflect", "", model), { history } as never, fakeStream() as never, fakeToken() as never);
    const [messages] = model.sendRequest.mock.calls[0] as unknown as [Array<{ content: string }>];
    expect(JSON.stringify(messages)).not.toContain("s3cret");
  });

  it("reports model failures instead of throwing", async () => {
    mockExec("packet");
    const model = { sendRequest: vi.fn(async () => Promise.reject(new Error("quota exceeded"))) };
    const stream = fakeStream();
    await handler(request("reflect", "", model), { history: [] } as never, stream as never, fakeToken() as never);
    expect(markdownOf(stream)).toContain("quota exceeded");
  });
});

describe("natural-language routing", () => {
  beforeEach(() => vi.clearAllMocks());

  it("routes to read-only commands and passes them no arguments", async () => {
    mockExec("ok");
    const res = await handler(
      { command: undefined, prompt: "show me the status of the login task", references: [] } as never,
      { history: [] } as never,
      fakeStream() as never,
      fakeToken() as never
    );
    expect(res?.metadata?.command).toBe("status");
    expect(ctxCalls()).toContainEqual(["status"]);
  });

  it("never mutates context on a keyword match", async () => {
    mockExec("ok");
    await handler(
      { command: undefined, prompt: "is the login task done? remind me to archive it", references: [] } as never,
      { history: [] } as never,
      fakeStream() as never,
      fakeToken() as never
    );
    const mutating = ctxCalls().filter((argv) => ["task", "remind", "pad", "add"].includes(argv[0]));
    expect(mutating).toEqual([]);
  });

  it("falls back to help", async () => {
    mockExec("ok");
    const res = await handler(
      { command: undefined, prompt: "hello", references: [] } as never,
      { history: [] } as never,
      fakeStream() as never,
      fakeToken() as never
    );
    expect(res?.metadata?.command).toBe("help");
  });
});
