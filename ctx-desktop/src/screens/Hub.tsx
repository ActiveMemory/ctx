import { useCallback, useEffect, useState } from "react";
import { ctxConnectionStatus } from "../adapter/ctx";
import { useReloadOnCtxChange } from "../hooks/useReload";

export default function Hub({ dir }: { dir: string }) {
  const [status, setStatus] = useState("");
  const [error, setError] = useState<string | null>(null);
  const [loading, setLoading] = useState(false);
  const reload = useReloadOnCtxChange(dir);

  const load = useCallback(async (d: string) => {
    setLoading(true);
    setError(null);
    try {
      setStatus(await ctxConnectionStatus(d));
    } catch (e) {
      // No hub configured (.connect.enc missing) lands here.
      setError(String(e));
      setStatus("");
    } finally {
      setLoading(false);
    }
  }, []);

  useEffect(() => {
    void load(dir);
  }, [dir, load, reload]);

  const notConnected = !!error && /\.connect\.enc/.test(error);

  return (
    <div className="mx-auto max-w-3xl px-6 py-6">
      <div className="mb-4 flex items-center justify-between">
        <h1 className="text-lg font-semibold text-ink">Hub</h1>
        <button
          onClick={() => void load(dir)}
          disabled={loading}
          className="rounded-md border border-border bg-panel px-3 py-1.5 text-xs text-ink hover:border-accent disabled:opacity-50"
        >
          {loading ? "Checking…" : "Refresh"}
        </button>
      </div>

      {notConnected ? (
        <div className="rounded-lg border border-border bg-panel px-4 py-8 text-center text-sm text-muted">
          Not connected to a ctx Hub. Register with{" "}
          <code className="text-ink">ctx connection register</code> to share
          context across projects.
        </div>
      ) : error ? (
        <div className="rounded-md border border-border bg-panel p-3 font-mono text-xs text-err">
          {error}
        </div>
      ) : (
        <pre className="overflow-auto rounded-lg border border-border bg-panel p-4 font-mono text-xs leading-relaxed text-ink">
          {status.trim() || (loading ? "Checking…" : "No status.")}
        </pre>
      )}
    </div>
  );
}
