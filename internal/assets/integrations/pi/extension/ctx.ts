// ctx Pi extension — thin shim to ctx system subcommands.
// All real logic lives in the ctx Go binary; this extension just
// wires Pi lifecycle events to ctx system calls.
//
// Hook envelope requirement (load-bearing):
// The `ctx system` hook subcommands (post-commit, check-task-completion,
// check-persistence) parse a hook-JSON envelope from stdin via
// FullPreamble and bail silently without one. Every hook call below is
// spawned with piped stdin: `{"session_id": ..., "tool_input": {command}}`
// is written and stdin closed immediately (closing also avoids the
// 2s StdinReadTimeout stall on open pipes). `ctx agent` is a normal
// flag-driven command and needs no envelope.
//
// node:child_process is mandatory for the envelope calls: Pi's own exec
// helper hardcodes `stdio: ["ignore", "pipe", "pipe"]` (no stdin pipe),
// so the envelope cannot be delivered through it.
//
// Compaction interop is breadcrumb-mediated and stateless across reloads:
// an injected custom_message is a valid compaction cut point, so once it
// ages past keepRecentTokens Pi folds it into the lossy LLM summary. The
// branch-scan predicate below re-injects exactly when no ctx-injected
// custom_message exists after the most recent compaction entry, so
// extension-instance teardown on /new /resume /fork /reload (which resets
// all in-memory state) can only cause an extra injection, never a missed
// one. We deliberately do NOT take ownership of the Pi LLM summary via
// session_before_compact: a custom summary replaces Pi's, and cross-
// extension precedence is undocumented.
//
// The packet is produced once at session_start (warm-up, off the prompt
// path) and cached; the cache is dropped on session_compact so the next
// re-injection is fresh. All subprocess calls pass ctx.signal and a
// timeout so Esc cancels cleanly; non-zero exits are swallowed (nothrow)
// and an absent `ctx` binary makes the extension no-op silently.
//
// Tool name strings target Pi's built-ins: `bash` (shell), `edit`,
// `write`. A failed tool (isError: true) never triggers post-commit —
// that would re-score the previous commit. The agent's own shell tool is
// NOT anchored by the extension; users launch pi from the project root so
// the Go binary resolves $PWD/.context/ (per ctx's cwd-anchored model).
// If Pi renames an event or built-in tool, the corresponding branch
// silently no-ops; verify against the Pi docs when Pi bumps.
import type {
	ExtensionAPI,
	ExtensionContext,
} from "@earendil-works/pi-coding-agent";
import { execFile } from "node:child_process";

const CTX_CONTEXT_TYPE = "ctx-context";
const SHELL_TOOL = "bash";
const EDIT_TOOLS = new Set(["edit", "write"]);
// Match `git commit` but not `git commit-tree` / `git commit-graph`.
// The negative lookahead rejects `-` immediately after the boundary.
const GIT_COMMIT_RE = /\bgit\s+commit\b(?!-)/;
const AGENT_BUDGET = 4000;
const SUBPROCESS_TIMEOUT_MS = 15_000;
const MAX_BUFFER = 1024 * 1024;

// runCtx spawns the ctx binary and resolves with trimmed stdout, or
// undefined on any failure (nothrow equivalent). Envelope calls pass
// stdin; stdin is always ended immediately.
function runCtx(
	args: string[],
	cwd: string,
	signal: AbortSignal | undefined,
	stdin?: string,
): Promise<string | undefined> {
	return new Promise((resolve) => {
		const child = execFile(
			"ctx",
			args,
			{
				cwd,
				timeout: SUBPROCESS_TIMEOUT_MS,
				signal,
				maxBuffer: MAX_BUFFER,
			},
			(err, stdout) => {
				if (err) {
					resolve(undefined);
					return;
				}
				const text = String(stdout).trim();
				resolve(text.length > 0 ? text : undefined);
			},
		);
		const input = child.stdin;
		if (input) {
			if (stdin !== undefined) {
				input.write(stdin);
			}
			input.end();
		}
	});
}

// runCtxHook is the envelope variant: writes the ctx system hook-JSON
// envelope (session_id + tool_input.command) before closing stdin.
function runCtxHook(
	ctx: ExtensionContext,
	args: string[],
	command: string,
): Promise<string | undefined> {
	const envelope = JSON.stringify({
		session_id: ctx.sessionManager.getSessionId(),
		tool_input: { command },
	});
	return runCtx(args, ctx.cwd, ctx.signal, envelope);
}

// fetchPacket runs the normal (non-hook) ctx agent command.
function fetchPacket(
	cwd: string,
	signal: AbortSignal | undefined,
): Promise<string | undefined> {
	return runCtx(
		["agent", "--budget", String(AGENT_BUDGET)],
		cwd,
		signal,
	);
}

// needsContextInjection: true when no ctx-injected custom_message exists
// after the most recent compaction entry in the current branch. Fresh
// sessions (no compaction, no injection yet) return true, so the first
// turn always gets the packet.
function needsContextInjection(
	sm: ExtensionContext["sessionManager"],
): boolean {
	const branch = sm.getBranch();
	let lastCompaction = -1;
	for (let i = branch.length - 1; i >= 0; i--) {
		if (branch[i].type === "compaction") {
			lastCompaction = i;
			break;
		}
	}
	for (let i = branch.length - 1; i > lastCompaction; i--) {
		const entry = branch[i];
		if (entry.type === "custom_message" && entry.customType === CTX_CONTEXT_TYPE) {
			return false;
		}
	}
	return true;
}

export default function (pi: ExtensionAPI) {
	let packetCache: string | undefined;

	pi.on("session_start", async (_event, ctx) => {
		// Warm-up: keep ctx agent off the prompt path.
		packetCache = await fetchPacket(ctx.cwd, ctx.signal);
	});

	pi.on("before_agent_start", async (_event, ctx) => {
		if (!needsContextInjection(ctx.sessionManager)) {
			return;
		}
		let packet = packetCache;
		if (packet === undefined) {
			packet = await fetchPacket(ctx.cwd, ctx.signal);
		}
		if (packet === undefined) {
			return;
		}
		return {
			message: {
				customType: CTX_CONTEXT_TYPE,
				content: packet,
				display: true,
			},
		};
	});

	pi.on("session_compact", async (_event, _ctx) => {
		// The branch-scan predicate subsumes the flag; dropping the cache
		// makes the next re-injection fresh.
		packetCache = undefined;
	});

	pi.on("tool_result", async (event, ctx) => {
		if (event.isError) {
			return;
		}
		if (event.toolName === SHELL_TOOL) {
			const command =
				typeof event.input.command === "string"
					? event.input.command
					: "";
			if (GIT_COMMIT_RE.test(command)) {
				await runCtxHook(ctx, ["system", "post-commit"], command);
			}
			return;
		}
		if (EDIT_TOOLS.has(event.toolName)) {
			await runCtxHook(ctx, ["system", "check-task-completion"], "");
		}
	});

	pi.on("agent_settled", async (_event, ctx) => {
		await runCtxHook(ctx, ["system", "check-persistence"], "");
	});
}
