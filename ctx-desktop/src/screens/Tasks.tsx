import { useCallback, useEffect, useMemo, useRef, useState } from "react";
import {
  ctxTasks,
  ctxTaskAdd,
  ctxTaskComplete,
  type Task,
} from "../adapter/ctx";
import { useReloadOnCtxChange } from "../hooks/useReload";

type Filter = "all" | "open" | "done";

const FILTERS: Filter[] = ["open", "done", "all"];
const PRIORITIES = ["", "high", "medium", "low"];

export default function Tasks({ dir }: { dir: string }) {
  const [tasks, setTasks] = useState<Task[]>([]);
  const [filter, setFilter] = useState<Filter>("open");
  const [text, setText] = useState("");
  const [priority, setPriority] = useState("");
  const [section, setSection] = useState("");
  const [busy, setBusy] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const reload = useReloadOnCtxChange(dir);
  // Only the latest load applies — guards against stale data on a fast
  // project switch or overlapping reload after an add/complete.
  const reqId = useRef(0);

  const load = useCallback(async (d: string) => {
    const id = ++reqId.current;
    setError(null);
    try {
      const next = await ctxTasks(d);
      if (id === reqId.current) setTasks(next);
    } catch (e) {
      if (id === reqId.current) {
        setError(String(e));
        setTasks([]);
      }
    }
  }, []);

  useEffect(() => {
    void load(dir);
  }, [dir, load, reload]);

  async function add() {
    if (!text.trim()) return;
    setBusy(true);
    setError(null);
    try {
      // ctx requires a target section/phase; default to Misc.
      await ctxTaskAdd(dir, text.trim(), priority, section.trim() || "Misc");
      setText("");
      setPriority("");
      // keep section so several tasks can be added to the same phase
      await load(dir);
    } catch (e) {
      setError(String(e));
    } finally {
      setBusy(false);
    }
  }

  // Complete by the task's exact text, not a locally computed pending
  // number: the CLI's numbering reflects ITS current file order, which
  // can drift from this list (external writes between load and click),
  // silently completing the wrong task. With the exact text, the CLI
  // matches the right task or fails loudly; an ambiguity error
  // (duplicate texts) surfaces through the normal error panel.
  async function complete(t: Task) {
    setBusy(true);
    setError(null);
    try {
      await ctxTaskComplete(dir, t.text);
      await load(dir);
    } catch (e) {
      setError(String(e));
    } finally {
      setBusy(false);
    }
  }

  const shown = tasks.filter((t) =>
    filter === "all"
      ? true
      : filter === "open"
        ? t.status === "pending"
        : t.status === "done",
  );

  const sections = useMemo(
    () =>
      Array.from(new Set(tasks.map((t) => t.section).filter(Boolean))).sort(),
    [tasks],
  );

  return (
    <div className="mx-auto max-w-3xl px-6 py-6">
      <h1 className="mb-4 text-lg font-semibold text-ink">Tasks</h1>

      {/* inline add */}
      <div className="mb-4 flex gap-2">
        <input
          value={text}
          onChange={(e) => setText(e.target.value)}
          onKeyDown={(e) => {
            if (e.key === "Enter") void add();
          }}
          placeholder="Add a task…"
          className="flex-1 rounded-md border border-border bg-panel px-3 py-2 text-sm text-ink outline-none focus:border-accent"
        />
        <input
          list="task-sections"
          value={section}
          onChange={(e) => setSection(e.target.value)}
          placeholder="Section (Misc)"
          title="Target phase/section in TASKS.md"
          className="w-32 rounded-md border border-border bg-panel px-3 py-2 text-sm text-ink outline-none focus:border-accent"
        />
        <datalist id="task-sections">
          {sections.map((s) => (
            <option key={s} value={s} />
          ))}
        </datalist>
        <select
          value={priority}
          onChange={(e) => setPriority(e.target.value)}
          className="rounded-md border border-border bg-panel px-2 py-2 text-sm text-ink outline-none focus:border-accent"
        >
          {PRIORITIES.map((p) => (
            <option key={p || "none"} value={p}>
              {p || "priority"}
            </option>
          ))}
        </select>
        <button
          onClick={() => void add()}
          disabled={busy || !text.trim()}
          className="rounded-md bg-accent px-4 py-2 text-sm font-medium text-bg disabled:opacity-50"
        >
          Add
        </button>
      </div>

      {/* filter tabs */}
      <div className="mb-3 flex gap-1">
        {FILTERS.map((f) => (
          <button
            key={f}
            onClick={() => setFilter(f)}
            className={`rounded-md px-3 py-1 text-xs capitalize ${
              filter === f
                ? "bg-accent/15 text-accent"
                : "text-muted hover:text-ink"
            }`}
          >
            {f}
          </button>
        ))}
        <span className="ml-auto self-center font-mono text-xs text-muted">
          {shown.length} shown
        </span>
      </div>

      {error && (
        <div className="mb-3 rounded-md border border-border bg-panel p-3 font-mono text-xs text-err">
          {error}
        </div>
      )}

      {/* list */}
      <ul className="divide-y divide-border overflow-hidden rounded-lg border border-border">
        {shown.length === 0 && (
          <li className="px-4 py-6 text-center text-sm text-muted">
            No tasks to show.
          </li>
        )}
        {shown.map((t) => (
          <li
            key={`${t.section}|${t.added}|${t.text}`}
            className="flex items-start gap-3 bg-panel px-4 py-2.5"
          >
            <button
              onClick={() => void complete(t)}
              disabled={busy || t.status === "done"}
              title={t.status === "done" ? "Done" : "Mark complete"}
              className={`mt-0.5 grid h-4 w-4 shrink-0 place-items-center rounded border ${
                t.status === "done"
                  ? "border-ok bg-ok/20 text-ok"
                  : "border-border hover:border-accent"
              }`}
            >
              {t.status === "done" ? "✓" : ""}
            </button>
            <div className="min-w-0 flex-1">
              <div
                className={`text-sm ${
                  t.status === "done" ? "text-muted line-through" : "text-ink"
                } ${t.is_sub ? "pl-3" : ""}`}
              >
                {t.text}
              </div>
              <div className="mt-0.5 flex flex-wrap gap-2 font-mono text-[11px] text-muted">
                {t.section && <span>{t.section}</span>}
                {t.priority && (
                  <span className="text-warn">#{t.priority}</span>
                )}
              </div>
            </div>
          </li>
        ))}
      </ul>
    </div>
  );
}
