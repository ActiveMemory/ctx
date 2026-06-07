import { useCallback, useEffect, useMemo, useState } from "react";
import { ctxDecisions, ctxDecisionAdd, type Decision } from "../adapter/ctx";
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

export default function Decisions({ dir }: { dir: string }) {
  const [decisions, setDecisions] = useState<Decision[]>([]);
  const [query, setQuery] = useState("");
  const [error, setError] = useState<string | null>(null);
  // Keyed by decision timestamp (stable) — NOT array index, which
  // would point at the wrong row once the list is filtered/reordered.
  const [expanded, setExpanded] = useState<Set<string>>(new Set());

  // authoring form
  const [formOpen, setFormOpen] = useState(false);
  const [title, setTitle] = useState("");
  const [context, setContext] = useState("");
  const [rationale, setRationale] = useState("");
  const [consequence, setConsequence] = useState("");
  const [busy, setBusy] = useState(false);
  const reload = useReloadOnCtxChange();

  const load = useCallback(async (d: string) => {
    setError(null);
    try {
      setDecisions(await ctxDecisions(d));
    } catch (e) {
      setError(String(e));
      setDecisions([]);
    }
  }, []);

  useEffect(() => {
    void load(dir);
  }, [dir, load, reload]);

  const filtered = useMemo(() => {
    const q = query.trim().toLowerCase();
    if (!q) return decisions;
    return decisions.filter((d) =>
      [d.title, d.context, d.rationale, d.consequence]
        .join(" ")
        .toLowerCase()
        .includes(q),
    );
  }, [decisions, query]);

  const canSave =
    title.trim() && context.trim() && rationale.trim() && consequence.trim();

  async function save() {
    if (!canSave) return;
    setBusy(true);
    setError(null);
    try {
      await ctxDecisionAdd(dir, title.trim(), context.trim(), rationale.trim(), consequence.trim());
      setTitle("");
      setContext("");
      setRationale("");
      setConsequence("");
      setFormOpen(false);
      await load(dir);
    } catch (e) {
      setError(String(e));
    } finally {
      setBusy(false);
    }
  }

  function toggle(id: string) {
    setExpanded((prev) => {
      const next = new Set(prev);
      next.has(id) ? next.delete(id) : next.add(id);
      return next;
    });
  }

  return (
    <div className="mx-auto max-w-3xl px-6 py-6">
      <div className="mb-4 flex items-center justify-between">
        <h1 className="text-lg font-semibold text-ink">Decisions</h1>
        <button
          onClick={() => setFormOpen((v) => !v)}
          className="rounded-md bg-accent px-3 py-1.5 text-sm font-medium text-bg"
        >
          {formOpen ? "Cancel" : "New decision"}
        </button>
      </div>

      {formOpen && (
        <div className="mb-5 rounded-lg border border-border bg-panel p-4">
          <input
            value={title}
            onChange={(e) => setTitle(e.target.value)}
            placeholder="Decision title"
            className="mb-2 w-full rounded-md border border-border bg-bg px-3 py-2 text-sm text-ink outline-none focus:border-accent"
          />
          {(
            [
              ["Context — what prompted this?", context, setContext],
              ["Rationale — why this over alternatives?", rationale, setRationale],
              ["Consequence — what changes as a result?", consequence, setConsequence],
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
              All four fields are required (this is an ADR).
            </span>
            <button
              onClick={() => void save()}
              disabled={busy || !canSave}
              className="rounded-md bg-accent px-4 py-2 text-sm font-medium text-bg disabled:opacity-50"
            >
              {busy ? "Saving…" : "Save decision"}
            </button>
          </div>
        </div>
      )}

      <input
        value={query}
        onChange={(e) => setQuery(e.target.value)}
        placeholder="Search decisions…"
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
            No decisions to show.
          </div>
        )}
        {filtered.map((d) => (
          <div
            key={d.timestamp}
            className="rounded-lg border border-border bg-panel"
          >
            <button
              onClick={() => toggle(d.timestamp)}
              className="flex w-full items-center gap-3 px-4 py-3 text-left"
            >
              <div className="min-w-0 flex-1">
                <div
                  className={`truncate text-sm font-medium ${
                    d.superseded ? "text-muted line-through" : "text-ink"
                  }`}
                >
                  {d.title}
                </div>
                <div className="mt-0.5 font-mono text-[11px] text-muted">
                  {d.date}
                </div>
              </div>
              {d.status && (
                <span className="shrink-0 rounded-full bg-ok/15 px-2 py-0.5 text-[11px] text-ok">
                  {d.status}
                </span>
              )}
              <span className="shrink-0 text-muted">
                {expanded.has(d.timestamp) ? "−" : "+"}
              </span>
            </button>
            {expanded.has(d.timestamp) && (
              <div className="border-t border-border px-4 py-3">
                <Field label="Context" value={d.context} />
                <Field label="Rationale" value={d.rationale} />
                <Field label="Consequence" value={d.consequence} />
              </div>
            )}
          </div>
        ))}
      </div>
    </div>
  );
}
