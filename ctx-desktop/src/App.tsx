import { useEffect, useState } from "react";
import Overview from "./screens/Overview";
import Tasks from "./screens/Tasks";
import Decisions from "./screens/Decisions";
import Learnings from "./screens/Learnings";
import ContextPacket from "./screens/ContextPacket";
import Journal from "./screens/Journal";
import { open } from "@tauri-apps/plugin-dialog";
import { ctxInfo, ctxDoctor, type CtxInfo, type DoctorReport } from "./adapter/ctx";

const RECENTS_KEY = "ctx.recents";

function loadRecents(): string[] {
  try {
    const v = JSON.parse(localStorage.getItem(RECENTS_KEY) ?? "[]");
    return Array.isArray(v) ? v : [];
  } catch {
    return [];
  }
}

// Default project = the ctx repo itself, so the app shows real
// data on first launch. Editable in the top bar.
const DEFAULT_DIR = "/Users/hamzaerbay/Code/ctx";

type View =
  | "overview"
  | "tasks"
  | "decisions"
  | "learnings"
  | "packet"
  | "journal";
const NAV: { id: View; label: string }[] = [
  { id: "overview", label: "Overview" },
  { id: "tasks", label: "Tasks" },
  { id: "decisions", label: "Decisions" },
  { id: "learnings", label: "Learnings" },
  { id: "packet", label: "Context Packet" },
  { id: "journal", label: "Journal" },
];

function HealthPill({ health }: { health: DoctorReport }) {
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
  return (
    <span className={`shrink-0 rounded-full px-3 py-1 text-xs ${cls}`}>
      doctor: {label}
    </span>
  );
}

function App() {
  const [dir, setDir] = useState(DEFAULT_DIR);
  const [draftDir, setDraftDir] = useState(DEFAULT_DIR);
  const [view, setView] = useState<View>("overview");
  const [info, setInfo] = useState<CtxInfo | null>(null);
  const [health, setHealth] = useState<DoctorReport | null>(null);
  const [recents, setRecents] = useState<string[]>(loadRecents);

  function applyDir(d: string) {
    if (!d) return;
    setDir(d);
    setDraftDir(d);
    setRecents((prev) => {
      const next = [d, ...prev.filter((p) => p !== d)].slice(0, 8);
      localStorage.setItem(RECENTS_KEY, JSON.stringify(next));
      return next;
    });
  }

  async function pickFolder() {
    const selected = await open({ directory: true, title: "Open a ctx project" });
    if (typeof selected === "string") applyDir(selected);
  }

  useEffect(() => {
    void ctxInfo().then(setInfo);
  }, []);

  useEffect(() => {
    ctxDoctor(dir)
      .then(setHealth)
      .catch(() => setHealth(null));
  }, [dir]);

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
      </aside>

      <div className="flex min-w-0 flex-1 flex-col">
        <header className="flex items-center gap-2 border-b border-border bg-panel px-4 py-2">
          <input
            value={draftDir}
            onChange={(e) => setDraftDir(e.target.value)}
            onKeyDown={(e) => {
              if (e.key === "Enter") applyDir(draftDir);
            }}
            spellCheck={false}
            className="flex-1 rounded-md border border-border bg-bg px-3 py-1.5 font-mono text-xs text-ink outline-none focus:border-accent"
            placeholder="/path/to/project (parent of .context)"
          />
          <button
            onClick={() => applyDir(draftDir)}
            className="rounded-md bg-accent px-3 py-1.5 text-xs font-medium text-bg"
          >
            Open
          </button>
          <button
            onClick={() => void pickFolder()}
            title="Open a folder…"
            className="rounded-md border border-border bg-bg px-3 py-1.5 text-xs text-ink hover:border-accent"
          >
            Folder…
          </button>
          {recents.length > 0 && (
            <select
              value=""
              onChange={(e) => e.target.value && applyDir(e.target.value)}
              title="Recent projects"
              className="max-w-40 rounded-md border border-border bg-bg px-2 py-1.5 text-xs text-ink outline-none focus:border-accent"
            >
              <option value="">Recent…</option>
              {recents.map((r) => (
                <option key={r} value={r}>
                  {r.split("/").pop() || r}
                </option>
              ))}
            </select>
          )}
          {health && <HealthPill health={health} />}
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
          {view === "packet" && <ContextPacket dir={dir} />}
          {view === "journal" && <Journal dir={dir} />}
        </main>
      </div>
    </div>
  );
}

export default App;
