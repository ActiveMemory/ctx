import { useCallback, useEffect, useState } from "react";
import { ctxJournal } from "../adapter/ctx";
import { useReloadOnCtxChange } from "../hooks/useReload";

const LIMITS = [5, 10, 20, 50];

export default function Journal({ dir }: { dir: string }) {
  const [limit, setLimit] = useState(10);
  const [text, setText] = useState("");
  const [error, setError] = useState<string | null>(null);
  const [loading, setLoading] = useState(false);
  const reload = useReloadOnCtxChange(dir);

  const load = useCallback(async (d: string, n: number) => {
    setLoading(true);
    setError(null);
    try {
      setText(await ctxJournal(d, n));
    } catch (e) {
      setError(String(e));
      setText("");
    } finally {
      setLoading(false);
    }
  }, []);

  useEffect(() => {
    void load(dir, limit);
  }, [dir, limit, load, reload]);

  return (
    <div className="mx-auto max-w-3xl px-6 py-6">
      <div className="mb-4 flex items-center justify-between">
        <h1 className="text-lg font-semibold text-ink">Journal</h1>
        <div className="flex items-center gap-2">
          {loading && <span className="text-xs text-muted">loading…</span>}
          <label className="text-xs text-muted">Last</label>
          <select
            value={limit}
            onChange={(e) => setLimit(Number(e.target.value))}
            className="rounded-md border border-border bg-panel px-2 py-1 text-sm text-ink outline-none focus:border-accent"
          >
            {LIMITS.map((n) => (
              <option key={n} value={n}>
                {n}
              </option>
            ))}
          </select>
        </div>
      </div>

      {error && (
        <div className="mb-3 rounded-md border border-border bg-panel p-3 font-mono text-xs text-err">
          {error}
        </div>
      )}

      <pre className="overflow-auto rounded-lg border border-border bg-panel p-4 font-mono text-xs leading-relaxed text-ink">
        {text || "No session history."}
      </pre>

      <p className="mt-3 text-xs text-muted">
        Rendered verbatim from <code className="text-ink">ctx journal source</code>.
        A structured timeline awaits a <code className="text-ink">journal
        source --json</code> mode upstream (same pattern as the list commands).
      </p>
    </div>
  );
}
