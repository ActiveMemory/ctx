import {
  useCallback,
  useEffect,
  useRef,
  useState,
  type ReactNode,
} from "react";
import {
  projectSummary,
  ctxTasks,
  ctxJournal,
  ctxJournalShow,
  type Project,
  type ProjectSummary,
  type Task,
} from "../adapter/ctx";
import { useReloadOnCtxChange } from "../hooks/useReload";

// Above this many projects, summaries (2 ctx spawns each) aren't
// auto-loaded — the user opts in with "Load all" to avoid a process
// storm on a large workspace.
const AUTO_CAP = 60;

// Run `fn` over `items` at most `limit` at a time, so a big workspace
// doesn't spawn one `ctx` process per project all at once.
async function mapLimit<T, R>(
  items: T[],
  limit: number,
  fn: (item: T) => Promise<R>,
): Promise<R[]> {
  const out: R[] = new Array(items.length);
  let next = 0;
  async function worker() {
    while (next < items.length) {
      const i = next++;
      out[i] = await fn(items[i]);
    }
  }
  await Promise.all(
    Array.from({ length: Math.min(limit, items.length) }, worker),
  );
  return out;
}

// Per-project drill-down: either real task rows (when `task list --json`
// is available) or the journal activity feed (the 0.8.1 fallback).
interface Detail {
  state: "loading" | "tasks" | "journal" | "error";
  tasks?: Task[];
  journal?: string;
  message?: string;
}

function Pill({
  tone,
  children,
}: {
  tone: "ok" | "warn" | "err" | "muted" | "accent";
  children: ReactNode;
}) {
  const cls = {
    ok: "bg-ok/15 text-ok",
    warn: "bg-warn/15 text-warn",
    err: "bg-err/15 text-err",
    muted: "bg-border/40 text-muted",
    accent: "bg-accent/15 text-accent",
  }[tone];
  return (
    <span className={`rounded-full px-2 py-0.5 text-[11px] ${cls}`}>
      {children}
    </span>
  );
}

function num(n: number | null): string {
  return n === null ? "—" : n.toLocaleString();
}

// "what's happening" for one task: provenance + lazy journal entry.
function TaskRow({ dir, task }: { dir: string; task: Task }) {
  const [open, setOpen] = useState(false);
  const [entry, setEntry] = useState<string | null>(null);
  const [loading, setLoading] = useState(false);

  async function toggle() {
    const next = !open;
    setOpen(next);
    if (next && entry === null && task.session) {
      setLoading(true);
      try {
        setEntry((await ctxJournalShow(dir, task.session)) || "(no journal entry for this session)");
      } catch (e) {
        setEntry(`(could not load session: ${String(e)})`);
      } finally {
        setLoading(false);
      }
    }
  }

  return (
    <li className="px-3 py-2">
      <button
        onClick={() => void toggle()}
        className="flex w-full items-start gap-2 text-left"
      >
        <span
          className={`mt-0.5 grid h-4 w-4 shrink-0 place-items-center rounded border text-[10px] ${
            task.status === "done"
              ? "border-ok bg-ok/20 text-ok"
              : "border-border text-muted"
          }`}
        >
          {task.status === "done" ? "✓" : ""}
        </span>
        <span className="min-w-0 flex-1">
          <span
            className={`block text-sm ${task.status === "done" ? "text-muted line-through" : "text-ink"}`}
          >
            {task.text}
          </span>
          <span className="mt-0.5 flex flex-wrap gap-2 font-mono text-[11px] text-muted">
            {task.section && <span>{task.section}</span>}
            {task.priority && <span className="text-warn">#{task.priority}</span>}
            {task.branch && <span>⎇ {task.branch}</span>}
            {task.commit && <span>@{task.commit}</span>}
            {task.added && <span>{task.added}</span>}
          </span>
        </span>
        {task.session && (
          <span className="shrink-0 text-muted">{open ? "−" : "+"}</span>
        )}
      </button>
      {open && (
        <pre className="mt-2 max-h-72 overflow-auto whitespace-pre-wrap rounded-md border border-border bg-bg p-3 font-mono text-[11px] leading-relaxed text-muted">
          {loading ? "loading session…" : entry}
        </pre>
      )}
    </li>
  );
}

export default function Projects({
  workspaces,
  deadRoots,
  projects,
  scanning,
  onOpen,
  onAddWorkspace,
  onRemoveWorkspace,
}: {
  workspaces: string[];
  deadRoots: string[];
  projects: Project[];
  scanning: boolean;
  onOpen: (dir: string) => void;
  onAddWorkspace: () => void;
  onRemoveWorkspace: (path: string) => void;
}) {
  const [summaries, setSummaries] = useState<Record<string, ProjectSummary>>(
    {},
  );
  const [loading, setLoading] = useState(false);
  const [loadAll, setLoadAll] = useState(false);
  const [expanded, setExpanded] = useState<string | null>(null);
  const [details, setDetails] = useState<Record<string, Detail>>({});
  const reload = useReloadOnCtxChange();

  // The set of projects we actually survey: capped unless the user opts
  // into loading every one (see AUTO_CAP).
  const surveyed =
    loadAll || projects.length <= AUTO_CAP
      ? projects
      : projects.slice(0, AUTO_CAP);

  // Only the latest scan applies — a fast add→remove of workspaces can
  // otherwise let an older summary scan land last and show stale cards.
  const reqId = useRef(0);

  const loadSummaries = useCallback(async (ps: Project[]) => {
    const id = ++reqId.current;
    setLoading(true);
    setSummaries({});
    const results = await mapLimit(ps, 4, async (p) => ({
      path: p.path,
      summary: await projectSummary(p.path),
    }));
    if (id !== reqId.current) return; // superseded by a newer scan
    setSummaries(Object.fromEntries(results.map((r) => [r.path, r.summary])));
    setLoading(false);
  }, []);

  useEffect(() => {
    const list =
      loadAll || projects.length <= AUTO_CAP
        ? projects
        : projects.slice(0, AUTO_CAP);
    if (list.length) void loadSummaries(list);
    else setSummaries({});
    // `reload` re-scans when any watched project's .context changes.
  }, [projects, loadAll, loadSummaries, reload]);

  async function expand(dir: string) {
    if (expanded === dir) {
      setExpanded(null);
      return;
    }
    setExpanded(dir);
    if (details[dir]) return;
    setDetails((d) => ({ ...d, [dir]: { state: "loading" } }));
    try {
      // Prefer real task rows; fall back to the journal feed when the
      // installed ctx has no `task list --json`.
      const tasks = await ctxTasks(dir);
      setDetails((d) => ({ ...d, [dir]: { state: "tasks", tasks } }));
    } catch {
      try {
        const journal = await ctxJournal(dir, 8);
        setDetails((d) => ({ ...d, [dir]: { state: "journal", journal } }));
      } catch (e) {
        setDetails((d) => ({
          ...d,
          [dir]: { state: "error", message: String(e) },
        }));
      }
    }
  }

  // Attention-first: errors, then warnings/drift, float to the top.
  const ordered = [...surveyed].sort((a, b) => {
    const sa = summaries[a.path];
    const sb = summaries[b.path];
    const score = (s?: ProjectSummary) =>
      s ? (s.errors ?? 0) * 100 + (s.warnings ?? 0) * 10 + (s.hasDrift ? 1 : 0) : 0;
    return score(sb) - score(sa);
  });

  // Workspace-wide rollup across every loaded project.
  const loaded = Object.values(summaries);
  const sum = (pick: (s: ProjectSummary) => number | null) =>
    loaded.reduce((acc, s) => acc + (pick(s) ?? 0), 0);
  const totals = {
    projects: projects.length,
    open: sum((s) => s.tasksOpen),
    warnings: sum((s) => s.warnings),
    errors: sum((s) => s.errors),
    drifted: loaded.filter((s) => s.hasDrift).length,
  };

  if (!workspaces.length) {
    return (
      <div className="mx-auto max-w-md px-6 py-20 text-center">
        <div className="text-lg font-semibold text-ink">No workspaces</div>
        <p className="mt-2 text-sm text-muted">
          Add one or more workspace folders to survey every ctx project across
          them.
        </p>
        <button
          onClick={onAddWorkspace}
          className="mt-5 inline-flex h-9 items-center rounded-md bg-accent px-4 text-sm font-medium text-bg"
        >
          Add workspace…
        </button>
      </div>
    );
  }

  return (
    <div className="mx-auto max-w-5xl px-6 py-6">
      <div className="mb-3 flex items-center justify-between">
        <div>
          <h1 className="text-lg font-semibold text-ink">Projects</h1>
          <p className="text-xs text-muted">
            {projects.length} project{projects.length === 1 ? "" : "s"} across{" "}
            {workspaces.length} workspace{workspaces.length === 1 ? "" : "s"}
            {scanning && " · scanning…"}
          </p>
        </div>
        <button
          onClick={() => void loadSummaries(surveyed)}
          disabled={loading}
          className="rounded-md border border-border bg-panel px-3 py-1.5 text-xs text-ink hover:border-accent disabled:opacity-50"
        >
          {loading ? "Refreshing…" : "Refresh"}
        </button>
      </div>

      {/* Workspace roots: add / remove the folders that get scanned. */}
      <div className="mb-4 flex flex-wrap items-center gap-1.5">
        {workspaces.map((w) => {
          const dead = deadRoots.includes(w);
          return (
            <span
              key={w}
              title={dead ? `${w} (not found)` : w}
              className={`group inline-flex items-center gap-1 rounded-full border py-1 pl-3 pr-1.5 text-xs ${
                dead
                  ? "border-err/40 bg-err/10 text-err"
                  : "border-border bg-panel text-ink"
              }`}
            >
              <span className="max-w-44 truncate font-mono">
                {w.split("/").pop() || w}
              </span>
              {dead && <span className="text-[10px]">missing</span>}
              <button
                onClick={() => onRemoveWorkspace(w)}
                title="Remove this workspace"
                className="grid h-4 w-4 place-items-center rounded-full text-muted hover:bg-err/15 hover:text-err"
              >
                ✕
              </button>
            </span>
          );
        })}
        <button
          onClick={onAddWorkspace}
          className="rounded-full border border-dashed border-border px-3 py-1 text-xs text-muted hover:border-accent hover:text-accent"
        >
          + Add workspace
        </button>
      </div>

      {!loadAll && projects.length > AUTO_CAP && (
        <div className="mb-4 flex flex-wrap items-center gap-2 rounded-lg border border-warn/40 bg-warn/10 px-4 py-3 text-xs">
          <span className="text-ink">
            Showing the first {AUTO_CAP} of {projects.length} projects to avoid
            a slow scan.
          </span>
          <button
            onClick={() => setLoadAll(true)}
            className="ml-auto rounded-md border border-border bg-panel px-3 py-1 text-ink hover:border-accent"
          >
            Load all {projects.length}
          </button>
        </div>
      )}

      {projects.length > 0 && (
        <div className="mb-4 flex flex-wrap items-center gap-2 rounded-lg border border-border bg-panel px-4 py-3 text-xs">
          <span className="font-mono text-ink">{totals.projects}</span>
          <span className="text-muted">projects</span>
          <span className="mx-1 text-border">·</span>
          <span className="font-mono text-accent">{totals.open.toLocaleString()}</span>
          <span className="text-muted">open tasks</span>
          {totals.errors > 0 && (
            <>
              <span className="mx-1 text-border">·</span>
              <span className="font-mono text-err">{totals.errors}</span>
              <span className="text-muted">errors</span>
            </>
          )}
          {totals.warnings > 0 && (
            <>
              <span className="mx-1 text-border">·</span>
              <span className="font-mono text-warn">{totals.warnings}</span>
              <span className="text-muted">warnings</span>
            </>
          )}
          {totals.drifted > 0 && (
            <>
              <span className="mx-1 text-border">·</span>
              <span className="font-mono text-warn">{totals.drifted}</span>
              <span className="text-muted">drifted</span>
            </>
          )}
          {loading && <span className="ml-auto text-muted">scanning…</span>}
        </div>
      )}

      {projects.length === 0 && (
        <div className="rounded-lg border border-border bg-panel px-4 py-10 text-center text-sm text-muted">
          No ctx projects found in these workspaces.
        </div>
      )}

      <div className="grid gap-3 sm:grid-cols-2">
        {ordered.map((p) => {
          const s = summaries[p.path];
          const isOpen = expanded === p.path;
          const detail = details[p.path];
          return (
            <div
              key={p.path}
              className={`rounded-lg border bg-panel ${isOpen ? "border-accent sm:col-span-2" : "border-border"}`}
            >
              <div className="flex items-start justify-between gap-2 px-4 pt-3">
                <button
                  onClick={() => void expand(p.path)}
                  className="min-w-0 flex-1 text-left"
                >
                  <div className="flex items-center gap-2">
                    <span className="truncate text-sm font-medium text-ink">
                      {p.name}
                    </span>
                    {!p.has_git && <Pill tone="muted">no git</Pill>}
                  </div>
                  <div className="flex items-center gap-2 truncate font-mono text-[11px] text-muted">
                    {p.branch && <span className="text-accent">⎇ {p.branch}</span>}
                    <span className="truncate">{p.path}</span>
                  </div>
                </button>
                <button
                  onClick={() => onOpen(p.path)}
                  title="Open this project"
                  className="shrink-0 rounded-md bg-accent px-2.5 py-1 text-[11px] font-medium text-bg"
                >
                  Open
                </button>
              </div>

              <div className="flex flex-wrap items-center gap-1.5 px-4 py-3">
                {!s ? (
                  <span className="text-[11px] text-muted">loading…</span>
                ) : (
                  <>
                    <Pill tone="accent">{num(s.tasksOpen)} open</Pill>
                    <Pill tone="muted">{num(s.tasksDone)} done</Pill>
                    {s.errors !== null && s.errors > 0 && (
                      <Pill tone="err">{s.errors} err</Pill>
                    )}
                    {s.warnings !== null && s.warnings > 0 && (
                      <Pill tone="warn">{s.warnings} warn</Pill>
                    )}
                    {s.errors === 0 && s.warnings === 0 && (
                      <Pill tone="ok">healthy</Pill>
                    )}
                    {s.hasDrift && <Pill tone="warn">drift</Pill>}
                    <span className="ml-auto font-mono text-[11px] text-muted">
                      {num(s.totalFiles)} files · {num(s.totalTokens)} tok
                    </span>
                  </>
                )}
              </div>

              {s?.problems.length ? (
                <div className="border-t border-border px-4 py-2 font-mono text-[11px] text-muted">
                  {s.problems.join(" · ")}
                </div>
              ) : null}

              {isOpen && (
                <div className="border-t border-border px-4 py-3">
                  {!detail || detail.state === "loading" ? (
                    <div className="text-xs text-muted">loading detail…</div>
                  ) : detail.state === "tasks" ? (
                    <ul className="divide-y divide-border overflow-hidden rounded-md border border-border bg-bg">
                      {(detail.tasks ?? []).length === 0 && (
                        <li className="px-3 py-3 text-center text-xs text-muted">
                          No tasks.
                        </li>
                      )}
                      {(detail.tasks ?? []).map((t, i) => (
                        <TaskRow key={i} dir={p.path} task={t} />
                      ))}
                    </ul>
                  ) : detail.state === "journal" ? (
                    <>
                      <div className="mb-1 text-[11px] text-muted">
                        Recent activity (task list unavailable on this ctx —
                        showing journal)
                      </div>
                      <pre className="max-h-80 overflow-auto whitespace-pre-wrap rounded-md border border-border bg-bg p-3 font-mono text-[11px] leading-relaxed text-muted">
                        {detail.journal || "(no recent sessions)"}
                      </pre>
                    </>
                  ) : (
                    <div className="font-mono text-[11px] text-err">
                      {detail.message}
                    </div>
                  )}
                </div>
              )}
            </div>
          );
        })}
      </div>
    </div>
  );
}
