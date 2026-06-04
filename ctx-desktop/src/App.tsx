import { useEffect, useState } from "react";
import Overview from "./screens/Overview";
import Tasks from "./screens/Tasks";
import Decisions from "./screens/Decisions";
import Learnings from "./screens/Learnings";
import CanonicalDoc from "./screens/CanonicalDoc";
import ContextPacket from "./screens/ContextPacket";
import Journal from "./screens/Journal";
import Drift from "./screens/Drift";
import Health from "./screens/Health";
import { open } from "@tauri-apps/plugin-dialog";
import { invoke } from "@tauri-apps/api/core";
import {
  ctxInfo,
  ctxDoctor,
  discoverProjects,
  type CtxInfo,
  type DoctorReport,
  type Project,
} from "./adapter/ctx";

const DIR_KEY = "ctx.dir";
const WORKSPACE_KEY = "ctx.workspace";

// Fallback project = the ctx repo itself, so the app shows real
// data on first launch before a workspace is chosen.
const DEFAULT_DIR = "/Users/hamzaerbay/Code/ctx";

type View =
  | "overview"
  | "tasks"
  | "decisions"
  | "learnings"
  | "conventions"
  | "constitution"
  | "packet"
  | "journal"
  | "drift"
  | "health";
const NAV: { id: View; label: string }[] = [
  { id: "overview", label: "Overview" },
  { id: "tasks", label: "Tasks" },
  { id: "decisions", label: "Decisions" },
  { id: "learnings", label: "Learnings" },
  { id: "conventions", label: "Conventions" },
  { id: "constitution", label: "Constitution" },
  { id: "packet", label: "Context Packet" },
  { id: "journal", label: "Journal" },
  { id: "drift", label: "Drift" },
  { id: "health", label: "Health" },
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
  const [workspace, setWorkspace] = useState(
    () => localStorage.getItem(WORKSPACE_KEY) || "",
  );
  const [projects, setProjects] = useState<Project[]>([]);
  const [scanning, setScanning] = useState(false);
  const [view, setView] = useState<View>("overview");
  const [info, setInfo] = useState<CtxInfo | null>(null);
  const [health, setHealth] = useState<DoctorReport | null>(null);

  function applyDir(d: string) {
    if (!d) return;
    setDir(d);
    localStorage.setItem(DIR_KEY, d);
  }

  async function scan(root: string) {
    if (!root) return;
    setScanning(true);
    try {
      setProjects(await discoverProjects(root, 4));
    } catch {
      setProjects([]);
    } finally {
      setScanning(false);
    }
  }

  async function pickWorkspace() {
    const selected = await open({
      directory: true,
      title: "Choose a workspace folder",
    });
    if (typeof selected !== "string") return;
    setWorkspace(selected);
    localStorage.setItem(WORKSPACE_KEY, selected);
    await scan(selected);
  }

  useEffect(() => {
    void ctxInfo().then(setInfo);
    // Re-scan a previously chosen workspace so the dropdown is ready.
    if (workspace) void scan(workspace);
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, []);

  useEffect(() => {
    ctxDoctor(dir)
      .then(setHealth)
      .catch(() => setHealth(null));
  }, [dir]);

  // Watch the active project's .context/ for external writes.
  useEffect(() => {
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
        {workspace && (
          <div className="mt-auto border-t border-border px-4 py-3">
            <div className="text-[11px] uppercase tracking-wide text-muted">
              Workspace
            </div>
            <div
              className="truncate font-mono text-[11px] text-ink"
              title={workspace}
            >
              {workspace.split("/").pop() || workspace}
            </div>
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
                  : "No workspace"}
            </option>
            {projects.map((p) => (
              <option key={p.path} value={p.path}>
                {p.name}
                {p.has_git ? "" : " (no git)"}
              </option>
            ))}
          </select>
          <button
            onClick={() => void pickWorkspace()}
            title="Choose a workspace folder to scan for ctx projects"
            className="flex h-8 shrink-0 items-center rounded-md bg-accent px-3 text-xs font-medium text-bg"
          >
            Workspace…
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
        </header>

        <main className="min-h-0 flex-1 overflow-auto">
          {view === "overview" && <Overview dir={dir} />}
          {view === "tasks" && <Tasks dir={dir} />}
          {view === "decisions" && <Decisions dir={dir} />}
          {view === "learnings" && <Learnings dir={dir} />}
          {view === "conventions" && (
            <CanonicalDoc dir={dir} title="Conventions" file="CONVENTIONS.md" />
          )}
          {view === "constitution" && (
            <CanonicalDoc
              dir={dir}
              title="Constitution"
              file="CONSTITUTION.md"
            />
          )}
          {view === "packet" && <ContextPacket dir={dir} />}
          {view === "journal" && <Journal dir={dir} />}
          {view === "drift" && <Drift dir={dir} />}
          {view === "health" && <Health dir={dir} />}
        </main>
      </div>
    </div>
  );
}

export default App;
