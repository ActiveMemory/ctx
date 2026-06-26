import { useCallback, useEffect, useState } from "react";
import { ctxDrift } from "../adapter/ctx";
import { useReloadOnCtxChange } from "../hooks/useReload";

// Drift surfaces `ctx drift` on demand (always available, unlike the
// Health screen which only offers it when the doctor flags
// staleness) and exposes the mutating `--fix` behind a confirm.
export default function Drift({ dir }: { dir: string }) {
  const [report, setReport] = useState("");
  const [running, setRunning] = useState<"" | "inspect" | "fix">("");
  const [confirm, setConfirm] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const reload = useReloadOnCtxChange(dir);

  const inspect = useCallback(async (d: string) => {
    setRunning("inspect");
    setError(null);
    try {
      setReport(await ctxDrift(d, false));
    } catch (e) {
      setError(String(e));
      setReport("");
    } finally {
      setRunning("");
    }
  }, []);

  useEffect(() => {
    void inspect(dir);
  }, [dir, inspect, reload]);

  async function fix() {
    setRunning("fix");
    setError(null);
    setConfirm(false);
    try {
      setReport(await ctxDrift(dir, true));
    } catch (e) {
      setError(String(e));
    } finally {
      setRunning("");
    }
  }

  const busy = !!running;

  return (
    <div className="mx-auto max-w-3xl px-6 py-6">
      <div className="mb-1 flex items-center justify-between">
        <h1 className="text-lg font-semibold text-ink">Drift</h1>
        <div className="flex items-center gap-2">
          <button
            onClick={() => void inspect(dir)}
            disabled={busy}
            className="rounded-md border border-border bg-panel px-3 py-1.5 text-xs text-ink hover:border-accent disabled:opacity-50"
          >
            {running === "inspect" ? "Inspecting…" : "Re-inspect"}
          </button>
          {confirm ? (
            <span className="inline-flex items-center gap-1">
              <button
                onClick={() => void fix()}
                className="rounded-md bg-warn px-3 py-1.5 text-xs font-medium text-bg"
              >
                {running === "fix" ? "Fixing…" : "Confirm fix"}
              </button>
              <button
                onClick={() => setConfirm(false)}
                className="rounded-md border border-border px-3 py-1.5 text-xs text-muted hover:text-ink"
              >
                Cancel
              </button>
            </span>
          ) : (
            <button
              onClick={() => setConfirm(true)}
              disabled={busy}
              className="rounded-md border border-border bg-bg px-3 py-1.5 text-xs text-ink hover:border-accent disabled:opacity-50"
            >
              Auto-fix
            </button>
          )}
        </div>
      </div>

      <p className="mb-4 text-xs text-muted">
        Detects stale paths, broken references, and constitution
        violations in the context files. Auto-fix writes to{" "}
        <code className="text-ink">.context/</code> through{" "}
        <code className="text-ink">ctx</code> (git-backed, reversible).
      </p>

      {error && (
        <div className="mb-3 rounded-md border border-border bg-panel p-3 font-mono text-xs text-err">
          {error}
        </div>
      )}

      {!error && (
        <pre className="overflow-auto rounded-lg border border-border bg-panel p-4 font-mono text-xs leading-relaxed text-ink">
          {report.trim() || (busy ? "Running…" : "No output.")}
        </pre>
      )}
    </div>
  );
}
