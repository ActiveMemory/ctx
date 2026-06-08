// ctx adapter (frontend half) — typed wrappers over the Rust
// commands in src-tauri/src/ctx_adapter.rs. All ctx access funnels
// through here so a CLI/output change is a one-file fix.
import { invoke } from "@tauri-apps/api/core";

export interface CtxInfo {
  found: boolean;
  version: string;
  error: string | null;
}

// Mirrors `ctx status --json`.
export interface CtxStatus {
  context_dir: string;
  total_files: number;
  total_tokens: number;
  total_size: number;
  files: {
    name: string;
    tokens: number;
    is_empty: boolean;
    summary?: string;
  }[];
}

// Mirrors `ctx task list --json`.
export interface Task {
  text: string;
  status: "pending" | "done";
  section: string;
  priority: string;
  is_sub: boolean;
  session: string;
  branch: string;
  commit: string;
  added: string;
}

// Mirrors `ctx decision list --json`.
export interface Decision {
  timestamp: string;
  date: string;
  title: string;
  status: string;
  context: string;
  rationale: string;
  consequence: string;
  superseded: boolean;
}

// Mirrors `ctx learning list --json`.
export interface Learning {
  timestamp: string;
  date: string;
  title: string;
  context: string;
  lesson: string;
  application: string;
  superseded: boolean;
}

// Mirrors `ctx agent --format json`. Section arrays render
// generically; an empty array means that section was dropped.
export interface AgentPacket {
  generated: string;
  budget: number;
  tokens_used: number;
  read_order: string[];
  constitution: string[];
  tasks: string[];
  conventions: string[];
  decisions: string[];
  learnings: string[];
  summaries: string[];
  instruction: string;
}

// A ctx project discovered under a workspace root.
export interface Project {
  path: string;
  name: string;
  has_git: boolean;
  branch: string; // current git branch, "" when absent/detached
}

/** Detect the ctx binary and read its version. */
export function ctxInfo(): Promise<CtxInfo> {
  return invoke<CtxInfo>("ctx_info");
}

/** Override the ctx binary path (empty string clears it → PATH lookup). */
export function setCtxPath(path: string): Promise<void> {
  return invoke<void>("set_ctx_path", { path });
}

/** True when `<dir>/.context` exists — validates a restored project. */
export function dirIsCtxProject(dir: string): Promise<boolean> {
  return invoke<boolean>("dir_is_ctx_project", { dir });
}

/** Watch every project's `.context/` so the dashboard sees external writes. */
export function watchProjects(dirs: string[]): Promise<void> {
  return invoke<void>("watch_projects", { dirs });
}

/** Scan a workspace `root` (up to `maxDepth` levels) for ctx projects. */
export function discoverProjects(
  root: string,
  maxDepth = 4,
): Promise<Project[]> {
  return invoke<Project[]>("discover_projects", { root, maxDepth });
}

/** Structured context packet from `ctx agent --format json --budget N`. */
export async function ctxAgentPacket(
  dir: string,
  budget: number,
): Promise<AgentPacket> {
  return JSON.parse(await invoke<string>("ctx_agent_json", { dir, budget }));
}

/** Paste-ready markdown packet from `ctx agent --budget N`. */
export function ctxAgentMarkdown(dir: string, budget: number): Promise<string> {
  return invoke<string>("ctx_agent_md", { dir, budget });
}

/** `ctx status --json` for the project at `dir`. */
export async function ctxStatus(dir: string): Promise<CtxStatus> {
  return JSON.parse(await invoke<string>("ctx_status", { dir }));
}

/**
 * Raw content of a canonical `.context/<name>` file (allowlisted in
 * the Rust adapter to CONSTITUTION.md / CONVENTIONS.md). Returns ""
 * when the file is absent.
 */
export function ctxReadDoc(dir: string, name: string): Promise<string> {
  return invoke<string>("ctx_read_doc", { dir, name });
}

// Inventory of a project's `.context/kb/` for the KB browser.
export interface KbInfo {
  exists: boolean;
  docs: string[]; // present top-level kb files, in display order
  topics: string[]; // slash-joined topic slugs, sorted
}

/** `.context/kb/` inventory (existence, top-level docs, topics). */
export function kbInfo(dir: string): Promise<KbInfo> {
  return invoke<KbInfo>("kb_info", { dir });
}

/** Raw content of a kb-relative file under `.context/kb/` ("" if absent). */
export function kbRead(dir: string, rel: string): Promise<string> {
  return invoke<string>("kb_read", { dir, rel });
}

// Mirrors `ctx doctor --json`.
export interface DoctorReport {
  results: {
    name: string;
    category: string;
    status: string;
    message?: string;
  }[];
  warnings: number;
  errors: number;
}

/** `ctx doctor --json` health report for the project at `dir`. */
export async function ctxDoctor(dir: string): Promise<DoctorReport> {
  return JSON.parse(await invoke<string>("ctx_doctor", { dir }));
}

/** Raw `ctx journal source --limit N` table text for `dir`. */
export function ctxJournal(dir: string, limit: number): Promise<string> {
  return invoke<string>("ctx_journal", { dir, limit });
}

/** Raw `ctx journal source --show <session>` text for one session. */
export function ctxJournalShow(dir: string, session: string): Promise<string> {
  return invoke<string>("ctx_journal_show", { dir, session });
}

// One project's at-a-glance health for the multi-project dashboard.
// Derived entirely from `status --json` + `doctor --json`, so it works
// on stock ctx 0.8.1 (no `task list --json` needed). Any failed source
// is recorded in `errors` and leaves its fields null rather than
// throwing — one broken project must not blank the whole grid.
export interface ProjectSummary {
  tasksOpen: number | null;
  tasksDone: number | null;
  decisions: number | null;
  totalFiles: number | null;
  totalTokens: number | null;
  warnings: number | null;
  errors: number | null;
  hasDrift: boolean;
  problems: string[];
}

// Parses a TASKS.md status summary like "237 active, 11 completed".
function parseTaskSummary(s: string): { open: number; done: number } | null {
  const open = /(\d+)\s+active/.exec(s);
  const done = /(\d+)\s+completed/.exec(s);
  if (!open && !done) return null;
  return { open: open ? +open[1] : 0, done: done ? +done[1] : 0 };
}

// Parses a DECISIONS.md summary like "112 decisions".
function parseCount(s: string): number | null {
  const m = /(\d+)/.exec(s);
  return m ? +m[1] : null;
}

/**
 * At-a-glance summary for one project, aggregating `status --json` and
 * `doctor --json`. Never rejects: per-source failures land in
 * `problems` so the dashboard can render a degraded card.
 */
export async function projectSummary(dir: string): Promise<ProjectSummary> {
  const out: ProjectSummary = {
    tasksOpen: null,
    tasksDone: null,
    decisions: null,
    totalFiles: null,
    totalTokens: null,
    warnings: null,
    errors: null,
    hasDrift: false,
    problems: [],
  };

  try {
    const st = await ctxStatus(dir);
    out.totalFiles = st.total_files;
    out.totalTokens = st.total_tokens;
    for (const f of st.files) {
      const summary = f.summary ?? "";
      if (f.name === "TASKS.md") {
        const t = parseTaskSummary(summary);
        if (t) {
          out.tasksOpen = t.open;
          out.tasksDone = t.done;
        }
      } else if (f.name === "DECISIONS.md") {
        out.decisions = parseCount(summary);
      }
    }
  } catch (e) {
    out.problems.push(`status: ${String(e)}`);
  }

  try {
    const doc = await ctxDoctor(dir);
    out.warnings = doc.warnings;
    out.errors = doc.errors;
    out.hasDrift = doc.results.some(
      (r) => r.name === "drift" && (r.status === "warning" || r.status === "error"),
    );
  } catch (e) {
    out.problems.push(`doctor: ${String(e)}`);
  }

  return out;
}

/** `ctx drift` report, or `ctx drift --fix` to auto-correct, for `dir`. */
export function ctxDrift(dir: string, fix = false): Promise<string> {
  return invoke<string>("ctx_drift", { dir, fix });
}

/** `ctx compact`, or `ctx compact --archive` (mutating), for `dir`. */
export function ctxCompact(dir: string, archive = false): Promise<string> {
  return invoke<string>("ctx_compact", { dir, archive });
}

/** `ctx remind list` — raw text of pending reminders. */
export function ctxRemindList(dir: string): Promise<string> {
  return invoke<string>("ctx_remind_list", { dir });
}

/** `ctx remind add <text>`. */
export function ctxRemindAdd(dir: string, text: string): Promise<string> {
  return invoke<string>("ctx_remind_add", { dir, text });
}

/** `ctx remind dismiss <target>` — number or "all". */
export function ctxRemindDismiss(dir: string, target: string): Promise<string> {
  return invoke<string>("ctx_remind_dismiss", { dir, target });
}

/** `ctx pad` — raw text list of scratchpad entries. */
export function ctxPadList(dir: string): Promise<string> {
  return invoke<string>("ctx_pad_list", { dir });
}

/** `ctx pad add <text>`. */
export function ctxPadAdd(dir: string, text: string): Promise<string> {
  return invoke<string>("ctx_pad_add", { dir, text });
}

/** `ctx pad rm <n>`. */
export function ctxPadRm(dir: string, n: string): Promise<string> {
  return invoke<string>("ctx_pad_rm", { dir, n });
}

/** `ctx pad show <n>` — raw text of a single entry. */
export function ctxPadShow(dir: string, n: string): Promise<string> {
  return invoke<string>("ctx_pad_show", { dir, n });
}

/** `ctx connection status` — hub status (rejects when no hub configured). */
export function ctxConnectionStatus(dir: string): Promise<string> {
  return invoke<string>("ctx_connection_status", { dir });
}

/** `ctx task list --json` for the project at `dir`. */
export async function ctxTasks(dir: string): Promise<Task[]> {
  const out = JSON.parse(await invoke<string>("ctx_task_list", { dir }));
  return out.tasks ?? [];
}

/** `ctx decision list --json` for the project at `dir`. */
export async function ctxDecisions(dir: string): Promise<Decision[]> {
  const out = JSON.parse(await invoke<string>("ctx_decision_list", { dir }));
  return out.decisions ?? [];
}

/** `ctx learning list --json` for the project at `dir`. */
export async function ctxLearnings(dir: string): Promise<Learning[]> {
  const out = JSON.parse(await invoke<string>("ctx_learning_list", { dir }));
  return out.learnings ?? [];
}

/**
 * `ctx task add` — provenance (session id, branch, commit) is
 * synthesized in the Rust adapter. Empty priority/section omitted.
 */
export function ctxTaskAdd(
  dir: string,
  text: string,
  priority = "",
  section = "",
): Promise<string> {
  return invoke<string>("ctx_task_add", { dir, text, priority, section });
}

/** `ctx task complete <id-or-text>`. */
export function ctxTaskComplete(dir: string, target: string): Promise<string> {
  return invoke<string>("ctx_task_complete", { dir, target });
}

/**
 * `ctx decision add` — all three ADR fields required; provenance
 * (session id, branch, commit) synthesized in the Rust adapter.
 */
export function ctxDecisionAdd(
  dir: string,
  title: string,
  context: string,
  rationale: string,
  consequence: string,
): Promise<string> {
  return invoke<string>("ctx_decision_add", {
    dir,
    title,
    context,
    rationale,
    consequence,
  });
}

/**
 * `ctx learning add` — all three fields required; provenance
 * synthesized in the Rust adapter.
 */
export function ctxLearningAdd(
  dir: string,
  title: string,
  context: string,
  lesson: string,
  application: string,
): Promise<string> {
  return invoke<string>("ctx_learning_add", {
    dir,
    title,
    context,
    lesson,
    application,
  });
}
