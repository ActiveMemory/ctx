import { useEffect, useState } from "react";
import Overview from "./screens/Overview";
import Search from "./screens/Search";
import Tasks from "./screens/Tasks";
import Reminders from "./screens/Reminders";
import Decisions from "./screens/Decisions";
import Learnings from "./screens/Learnings";
import CanonicalDoc from "./screens/CanonicalDoc";
import ContextPacket from "./screens/ContextPacket";
import KnowledgeBase from "./screens/KnowledgeBase";
import Pad from "./screens/Pad";
import Journal from "./screens/Journal";
import Drift from "./screens/Drift";
import Health from "./screens/Health";
import Hub from "./screens/Hub";
import Projects from "./screens/Projects";
import { open } from "@tauri-apps/plugin-dialog";
import { invoke } from "@tauri-apps/api/core";
import {
  ctxInfo,
  ctxDoctor,
  discoverProjects,
  dirIsCtxProject,
  setCtxPath,
  watchProjects,
  type CtxInfo,
  type DoctorReport,
  type Project,
} from "./adapter/ctx";

const DIR_KEY = "ctx.dir";
const WORKSPACES_KEY = "ctx.workspaces";
const LEGACY_WORKSPACE_KEY = "ctx.workspace";
const CTX_BIN_KEY = "ctx.bin";

// No baked-in project: until the user adds a workspace (persisted in
// localStorage), `dir` is empty and the main area shows a chooser.
// A hardcoded path only ever pointed at one machine's checkout.
const DEFAULT_DIR = "";

// Reads the persisted workspace roots, migrating the old single-value
// `ctx.workspace` key into the array form on first run.
function loadWorkspaces(): string[] {
  try {
    const raw = localStorage.getItem(WORKSPACES_KEY);
    if (raw) {
      const arr = JSON.parse(raw);
      if (Array.isArray(arr))
        return arr.filter((x): x is string => typeof x === "string");
    }
  } catch {
    // malformed value — fall through to legacy / empty
  }
  const legacy = localStorage.getItem(LEGACY_WORKSPACE_KEY);
  return legacy ? [legacy] : [];
}

type View =
  | "projects"
  | "overview"
  | "search"
  | "tasks"
  | "reminders"
  | "decisions"
  | "learnings"
  | "conventions"
  | "constitution"
  | "packet"
  | "kb"
  | "pad"
  | "journal"
  | "drift"
  | "health"
  | "hub";
const NAV: { id: View; label: string }[] = [
  { id: "projects", label: "Projects" },
  { id: "overview", label: "Overview" },
  { id: "search", label: "Search" },
  { id: "tasks", label: "Tasks" },
  { id: "reminders", label: "Reminders" },
  { id: "decisions", label: "Decisions" },
  { id: "learnings", label: "Learnings" },
  { id: "conventions", label: "Conventions" },
  { id: "constitution", label: "Constitution" },
  { id: "packet", label: "Context Packet" },
  { id: "kb", label: "Knowledge Base" },
  { id: "pad", label: "Scratchpad" },
  { id: "journal", label: "Journal" },
  { id: "drift", label: "Drift" },
  { id: "health", label: "Health" },
  { id: "hub", label: "Hub" },
];

function HealthPill({
  health,
  onClick,
}: {
  health: DoctorReport;
  onClick: () => void;
}) {
  const level =
    health.errors > 0 ? "err" : health.warnings > 0 ? "warn" : "ok";
  const label =
    level === "err"
      ? `${health.errors} error${health.errors > 1 ? "s" : ""}`
      : level === "warn"
        ? `${health.warnings} warning${health.warnings > 1 ? "s" : ""}`
        : "healthy";
  const cls =
    level === "err"
      ? "bg-err/15 text-err"
      : level === "warn"
        ? "bg-warn/15 text-warn"
        : "bg-ok/15 text-ok";
  const detail =
    health.results
      .filter((r) => r.status === "warning" || r.status === "error")
      .map((r) => `• [${r.category}] ${r.message || r.name}`)
      .join("\n") || "All structural checks passed.";
  return (
    <button
      onClick={onClick}
      className={`shrink-0 cursor-pointer rounded-full px-3 py-1 text-xs ${cls}`}
      title={`${detail}\n\n(click for the Health screen)`}
    >
      doctor: {label}
    </button>
  );
}

function App() {
  const [dir, setDir] = useState(
    () => localStorage.getItem(DIR_KEY) || DEFAULT_DIR,
  );
  const [workspaces, setWorkspaces] = useState<string[]>(loadWorkspaces);
  const [projects, setProjects] = useState<Project[]>([]);
  const [deadRoots, setDeadRoots] = useState<string[]>([]);
  const [scanning, setScanning] = useState(false);
  const [view, setView] = useState<View>("overview");
  const [info, setInfo] = useState<CtxInfo | null>(null);
  const [health, setHealth] = useState<DoctorReport | null>(null);

  function applyDir(d: string) {
    if (!d) return;
    setDir(d);
    localStorage.setItem(DIR_KEY, d);
  }

  // Scans every workspace root, merging results and de-duping by path
  // (overlapping roots can surface the same project). Roots that fail to
  // scan (deleted/unmounted) are tracked so the UI can flag them.
  async function scanAll(roots: string[]) {
    if (!roots.length) {
      setProjects([]);
      setDeadRoots([]);
      void watchProjects([]);
      return;
    }
    setScanning(true);
    try {
      const results = await Promise.all(
        roots.map((r) =>
          discoverProjects(r, 4).then(
            (ps) => ({ root: r, ps, ok: true }),
            () => ({ root: r, ps: [] as Project[], ok: false }),
          ),
        ),
      );
      const byPath = new Map<string, Project>();
      for (const res of results) for (const p of res.ps) byPath.set(p.path, p);
      const merged = [...byPath.values()].sort((a, b) =>
        a.name.toLowerCase().localeCompare(b.name.toLowerCase()),
      );
      setProjects(merged);
      setDeadRoots(results.filter((r) => !r.ok).map((r) => r.root));
      // Watch every project's .context/ so the dashboard reflects
      // external writes to any project, not just the active one.
      void watchProjects(merged.map((p) => p.path)).catch(() => {});
    } finally {
      setScanning(false);
    }
  }

  function saveWorkspaces(list: string[]) {
    setWorkspaces(list);
    localStorage.setItem(WORKSPACES_KEY, JSON.stringify(list));
  }

  async function addWorkspace() {
    const selected = await open({
      directory: true,
      title: "Add a workspace folder",
    });
    if (typeof selected !== "string" || workspaces.includes(selected)) return;
    const next = [...workspaces, selected];
    saveWorkspaces(next);
    await scanAll(next);
  }

  function removeWorkspace(path: string) {
    const next = workspaces.filter((w) => w !== path);
    saveWorkspaces(next);
    void scanAll(next);
  }

  // Let the user point at a ctx binary the app can't find on PATH.
  async function pickCtxBinary() {
    const selected = await open({
      title: "Locate the ctx binary",
      directory: false,
    });
    if (typeof selected !== "string") return;
    localStorage.setItem(CTX_BIN_KEY, selected);
    await setCtxPath(selected).catch(() => {});
    setInfo(await ctxInfo());
  }

  useEffect(() => {
    void (async () => {
      // Apply a saved ctx path BEFORE detecting the binary.
      const savedBin = localStorage.getItem(CTX_BIN_KEY);
      if (savedBin) await setCtxPath(savedBin).catch(() => {});
      void ctxInfo().then(setInfo);
      // Drop a restored active project that no longer exists, so the app
      // shows the chooser instead of erroring on every screen.
      if (dir) {
        const ok = await dirIsCtxProject(dir).catch(() => false);
        if (!ok) {
          setDir("");
          localStorage.removeItem(DIR_KEY);
        }
      }
      // Re-scan previously added workspaces so the dropdown/grid are ready.
      if (workspaces.length) void scanAll(workspaces);
    })();
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, []);

  useEffect(() => {
    if (!dir) {
      setHealth(null);
      return;
    }
    ctxDoctor(dir)
      .then(setHealth)
      .catch(() => setHealth(null));
  }, [dir]);

  // Watch the active project's .context/ for external writes.
  useEffect(() => {
    if (!dir) return;
    void invoke("watch_context", { dir }).catch(() => {});
  }, [dir]);

  const dirInProjects = projects.some((p) => p.path === dir);

  return (
    <div className="flex min-h-screen bg-bg text-ink">
      <aside className="flex w-48 shrink-0 flex-col border-r border-border bg-panel">
        <div className="px-4 py-4">
          <div className="text-sm font-semibold">ctx Desktop</div>
          <div className="text-xs text-muted">do you remember?</div>
        </div>
        <nav className="px-2">
          {NAV.map((n) => (
            <button
              key={n.id}
              onClick={() => setView(n.id)}
              className={`mb-1 w-full rounded-md px-3 py-2 text-left text-sm ${
                view === n.id
                  ? "bg-accent/15 text-accent"
                  : "text-muted hover:bg-border/40 hover:text-ink"
              }`}
            >
              {n.label}
            </button>
          ))}
        </nav>
        {workspaces.length > 0 && (
          <div className="mt-auto border-t border-border px-4 py-3">
            <div className="mb-1 text-[11px] uppercase tracking-wide text-muted">
              Workspaces ({workspaces.length})
            </div>
            <ul className="space-y-0.5">
              {workspaces.map((w) => {
                const dead = deadRoots.includes(w);
                return (
                  <li
                    key={w}
                    className="group flex items-center gap-1"
                    title={dead ? `${w} (not found)` : w}
                  >
                    <span
                      className={`truncate font-mono text-[11px] ${
                        dead ? "text-err line-through" : "text-ink"
                      }`}
                    >
                      {w.split("/").pop() || w}
                    </span>
                    {dead && (
                      <span className="shrink-0 text-[10px] text-err">
                        missing
                      </span>
                    )}
                    <button
                      onClick={() => removeWorkspace(w)}
                      title="Remove this workspace"
                      className="ml-auto shrink-0 text-muted opacity-0 hover:text-err group-hover:opacity-100"
                    >
                      ✕
                    </button>
                  </li>
                );
              })}
            </ul>
          </div>
        )}
      </aside>

      <div className="flex min-w-0 flex-1 flex-col">
        <header className="flex items-center gap-2 border-b border-border bg-panel px-4 py-2">
          {/* Workspace project switcher */}
          <select
            value={dirInProjects ? dir : ""}
            onChange={(e) => e.target.value && applyDir(e.target.value)}
            title="Switch project"
            className="h-8 min-w-48 max-w-72 rounded-md border border-border bg-bg px-2 text-xs text-ink outline-none focus:border-accent"
          >
            <option value="">
              {scanning
                ? "Scanning…"
                : projects.length
                  ? `${projects.length} project${projects.length > 1 ? "s" : ""}…`
                  : "No workspaces"}
            </option>
            {projects.map((p) => (
              <option key={p.path} value={p.path}>
                {p.name}
                {p.has_git ? "" : " (no git)"}
              </option>
            ))}
          </select>
          <button
            onClick={() => void addWorkspace()}
            title="Add a workspace folder to scan for ctx projects"
            className="flex h-8 shrink-0 items-center rounded-md bg-accent px-3 text-xs font-medium text-bg"
          >
            Add workspace…
          </button>

          <div className="flex-1" />

          {health && (
            <HealthPill health={health} onClick={() => setView("health")} />
          )}
          {info && (
            <span
              className={`shrink-0 rounded-full px-3 py-1 font-mono text-xs ${
                info.found ? "bg-ok/15 text-ok" : "bg-err/15 text-err"
              }`}
              title={info.error ?? ""}
            >
              {info.found ? info.version : "ctx not found"}
            </span>
          )}
          {info && !info.found && (
            <button
              onClick={() => void pickCtxBinary()}
              title="Locate the ctx binary if it isn't on PATH"
              className="shrink-0 rounded-md border border-border bg-panel px-3 py-1 text-xs text-ink hover:border-accent"
            >
              Set ctx path…
            </button>
          )}
        </header>

        <main className="min-h-0 flex-1 overflow-auto">
          {/* Projects is workspace-level: it surveys every discovered
              project, so it renders regardless of the active `dir`. */}
          {view === "projects" && (
            <Projects
              workspaces={workspaces}
              deadRoots={deadRoots}
              projects={projects}
              scanning={scanning}
              onOpen={(d) => {
                applyDir(d);
                setView("overview");
              }}
              onAddWorkspace={() => void addWorkspace()}
              onRemoveWorkspace={removeWorkspace}
            />
          )}
          {view !== "projects" && !dir && (
            <div className="mx-auto max-w-md px-6 py-20 text-center">
              <div className="text-lg font-semibold text-ink">
                No project selected
              </div>
              <p className="mt-2 text-sm text-muted">
                Add one or more workspace folders and ctx Desktop will scan them
                for projects (any directory with a <code>.context/</code>), then
                pick one from the <strong>Projects</strong> screen or the top-bar
                dropdown.
              </p>
              <button
                onClick={() => void addWorkspace()}
                className="mt-5 inline-flex h-9 items-center rounded-md bg-accent px-4 text-sm font-medium text-bg"
              >
                Add workspace…
              </button>
            </div>
          )}
          {view !== "projects" && dir && (
            <>
              {view === "overview" && <Overview dir={dir} />}
              {view === "search" && <Search dir={dir} onOpen={setView} />}
              {view === "tasks" && <Tasks dir={dir} />}
              {view === "reminders" && <Reminders dir={dir} />}
              {view === "decisions" && <Decisions dir={dir} />}
              {view === "learnings" && <Learnings dir={dir} />}
              {view === "conventions" && (
                <CanonicalDoc
                  dir={dir}
                  title="Conventions"
                  file="CONVENTIONS.md"
                />
              )}
              {view === "constitution" && (
                <CanonicalDoc
                  dir={dir}
                  title="Constitution"
                  file="CONSTITUTION.md"
                />
              )}
              {view === "packet" && <ContextPacket dir={dir} />}
              {view === "kb" && <KnowledgeBase dir={dir} />}
              {view === "pad" && <Pad dir={dir} />}
              {view === "journal" && <Journal dir={dir} />}
              {view === "drift" && <Drift dir={dir} />}
              {view === "health" && <Health dir={dir} />}
              {view === "hub" && <Hub dir={dir} />}
            </>
          )}
        </main>
      </div>
    </div>
  );
}

export default App;
