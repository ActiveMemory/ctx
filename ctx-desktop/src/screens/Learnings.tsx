import { useCallback, useEffect, useMemo, useState } from "react";
import { ctxLearnings, ctxLearningAdd, type Learning } from "../adapter/ctx";
import { useReloadOnCtxChange } from "../hooks/useReload";

function Field({ label, value }: { label: string; value: string }) {
  if (!value) return null;
  return (
    <div className="mt-2">
      <div className="text-[11px] font-medium uppercase tracking-wide text-muted">
        {label}
      </div>
      <div className="text-sm text-ink">{value}</div>
    </div>
  );
}

export default function Learnings({ dir }: { dir: string }) {
  const [learnings, setLearnings] = useState<Learning[]>([]);
  const [query, setQuery] = useState("");
  const [error, setError] = useState<string | null>(null);
  const [expanded, setExpanded] = useState<Set<number>>(new Set());

  const [formOpen, setFormOpen] = useState(false);
  const [title, setTitle] = useState("");
  const [context, setContext] = useState("");
  const [lesson, setLesson] = useState("");
  const [application, setApplication] = useState("");
  const [busy, setBusy] = useState(false);
  const reload = useReloadOnCtxChange();

  const load = useCallback(async (d: string) => {
    setError(null);
    try {
      setLearnings(await ctxLearnings(d));
    } catch (e) {
      setError(String(e));
      setLearnings([]);
    }
  }, []);

  useEffect(() => {
    void load(dir);
  }, [dir, load, reload]);

  const filtered = useMemo(() => {
    const q = query.trim().toLowerCase();
    if (!q) return learnings;
    return learnings.filter((l) =>
      [l.title, l.context, l.lesson, l.application]
        .join(" ")
        .toLowerCase()
        .includes(q),
    );
  }, [learnings, query]);

  const canSave =
    title.trim() && context.trim() && lesson.trim() && application.trim();

  async function save() {
    if (!canSave) return;
    setBusy(true);
    setError(null);
    try {
      await ctxLearningAdd(dir, title.trim(), context.trim(), lesson.trim(), application.trim());
      setTitle("");
      setContext("");
      setLesson("");
      setApplication("");
      setFormOpen(false);
      await load(dir);
    } catch (e) {
      setError(String(e));
    } finally {
      setBusy(false);
    }
  }

  function toggle(i: number) {
    setExpanded((prev) => {
      const next = new Set(prev);
      if (next.has(i)) next.delete(i);
      else next.add(i);
      return next;
    });
  }

  return (
    <div className="mx-auto max-w-3xl px-6 py-6">
      <div className="mb-4 flex items-center justify-between">
        <h1 className="text-lg font-semibold text-ink">Learnings</h1>
        <button
          onClick={() => setFormOpen((v) => !v)}
          className="rounded-md bg-accent px-3 py-1.5 text-sm font-medium text-bg"
        >
          {formOpen ? "Cancel" : "New learning"}
        </button>
      </div>

      {formOpen && (
        <div className="mb-5 rounded-lg border border-border bg-panel p-4">
          <input
            value={title}
            onChange={(e) => setTitle(e.target.value)}
            placeholder="Learning title"
            className="mb-2 w-full rounded-md border border-border bg-bg px-3 py-2 text-sm text-ink outline-none focus:border-accent"
          />
          {(
            [
              ["Context — what prompted this?", context, setContext],
              ["Lesson — the key insight", lesson, setLesson],
              ["Application — how to apply it going forward", application, setApplication],
            ] as const
          ).map(([ph, val, set]) => (
            <textarea
              key={ph}
              value={val}
              onChange={(e) => set(e.target.value)}
              placeholder={ph}
              rows={2}
              className="mb-2 w-full resize-y rounded-md border border-border bg-bg px-3 py-2 text-sm text-ink outline-none focus:border-accent"
            />
          ))}
          <div className="flex items-center justify-end gap-2">
            <span className="mr-auto text-xs text-muted">
              All four fields are required.
            </span>
            <button
              onClick={() => void save()}
              disabled={busy || !canSave}
              className="rounded-md bg-accent px-4 py-2 text-sm font-medium text-bg disabled:opacity-50"
            >
              {busy ? "Saving…" : "Save learning"}
            </button>
          </div>
        </div>
      )}

      <input
        value={query}
        onChange={(e) => setQuery(e.target.value)}
        placeholder="Search learnings…"
        className="mb-3 w-full rounded-md border border-border bg-panel px-3 py-2 text-sm text-ink outline-none focus:border-accent"
      />

      {error && (
        <div className="mb-3 rounded-md border border-border bg-panel p-3 font-mono text-xs text-err">
          {error}
        </div>
      )}

      <div className="space-y-2">
        {filtered.length === 0 && (
          <div className="rounded-lg border border-border bg-panel px-4 py-6 text-center text-sm text-muted">
            No learnings to show.
          </div>
        )}
        {filtered.map((l, i) => (
          <div key={i} className="rounded-lg border border-border bg-panel">
            <button
              onClick={() => toggle(i)}
              className="flex w-full items-center gap-3 px-4 py-3 text-left"
            >
              <div className="min-w-0 flex-1">
                <div
                  className={`truncate text-sm font-medium ${
                    l.superseded ? "text-muted line-through" : "text-ink"
                  }`}
                >
                  {l.title}
                </div>
                <div className="mt-0.5 font-mono text-[11px] text-muted">
                  {l.date}
                </div>
              </div>
              <span className="shrink-0 text-muted">
                {expanded.has(i) ? "−" : "+"}
              </span>
            </button>
            {expanded.has(i) && (
              <div className="border-t border-border px-4 py-3">
                <Field label="Context" value={l.context} />
                <Field label="Lesson" value={l.lesson} />
                <Field label="Application" value={l.application} />
              </div>
            )}
          </div>
        ))}
      </div>
    </div>
  );
}
