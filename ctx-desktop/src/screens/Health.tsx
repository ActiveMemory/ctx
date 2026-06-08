import { useCallback, useEffect, useRef, useState } from "react";
import {
  ctxDoctor,
  ctxDrift,
  ctxCompact,
  type DoctorReport,
} from "../adapter/ctx";
import { useReloadOnCtxChange } from "../hooks/useReload";

type Result = DoctorReport["results"][number];

const DOT: Record<string, string> = {
  ok: "bg-ok",
  info: "bg-muted",
  warning: "bg-warn",
  error: "bg-err",
};

function Row({ r }: { r: Result }) {
  return (
    <li className="flex items-start gap-2 px-4 py-2">
      <span
        className={`mt-1.5 h-2 w-2 shrink-0 rounded-full ${DOT[r.status] ?? "bg-muted"}`}
      />
      <div className="min-w-0">
        <div className="text-sm text-ink">{r.message || r.name}</div>
        <div className="font-mono text-[11px] text-muted">
          {r.category} · {r.name}
        </div>
      </div>
    </li>
  );
}

export default function Health({ dir }: { dir: string }) {
  const [report, setReport] = useState<DoctorReport | null>(null);
  const [output, setOutput] = useState("");
  const [running, setRunning] = useState("");
  const [confirm, setConfirm] = useState<"" | "drift" | "compact">("");
  const [error, setError] = useState<string | null>(null);
  const reload = useReloadOnCtxChange();
  // Only the latest load applies — guards against stale data on a fast
  // project switch or overlapping reload after a fix.
  const reqId = useRef(0);

  const load = useCallback(async (d: string) => {
    const id = ++reqId.current;
    setError(null);
    try {
      const next = await ctxDoctor(d);
      if (id === reqId.current) setReport(next);
    } catch (e) {
      if (id === reqId.current) {
        setError(String(e));
        setReport(null);
      }
    }
  }, []);

  useEffect(() => {
    void load(dir);
  }, [dir, load, reload]);

  async function run(
    label: string,
    fn: () => Promise<string>,
    refresh = false,
  ) {
    setRunning(label);
    setError(null);
    setOutput("");
    setConfirm("");
    try {
      setOutput(await fn());
      if (refresh) await load(dir);
    } catch (e) {
      setError(String(e));
    } finally {
      setRunning("");
    }
  }

  const results = report?.results ?? [];
  const warnings = results.filter(
    (r) => r.status === "warning" || r.status === "error",
  );
  const others = results.filter(
    (r) => r.status === "ok" || r.status === "info",
  );
  const hasDrift = results.some((r) => r.name === "drift");
  const hasSize = results.some((r) => r.name === "context_size");

  function ConfirmButton({
    which,
    label,
    onConfirm,
  }: {
    which: "drift" | "compact";
    label: string;
    onConfirm: () => void;
  }) {
    if (confirm === which) {
      return (
        <span className="inline-flex items-center gap-1">
          <button
            onClick={onConfirm}
            className="rounded-md bg-warn px-3 py-1.5 text-xs font-medium text-bg"
          >
            Confirm
          </button>
          <button
            onClick={() => setConfirm("")}
            className="rounded-md border border-border px-3 py-1.5 text-xs text-muted hover:text-ink"
          >
            Cancel
          </button>
        </span>
      );
    }
    return (
      <button
        onClick={() => setConfirm(which)}
        disabled={!!running}
        className="rounded-md border border-border bg-bg px-3 py-1.5 text-xs text-ink hover:border-accent disabled:opacity-50"
      >
        {label}
      </button>
    );
  }

  return (
    <div className="mx-auto max-w-3xl px-6 py-6">
      <div className="mb-4 flex items-center justify-between">
        <h1 className="text-lg font-semibold text-ink">Health</h1>
        <button
          onClick={() => void load(dir)}
          className="rounded-md border border-border bg-panel px-3 py-1.5 text-xs text-ink hover:border-accent"
        >
          Re-run doctor
        </button>
      </div>

      {report && (
        <div className="mb-5 flex gap-3 font-mono text-xs">
          <span className="text-err">{report.errors} errors</span>
          <span className="text-warn">{report.warnings} warnings</span>
          <span className="text-muted">
            {results.filter((r) => r.status === "ok").length} ok ·{" "}
            {results.filter((r) => r.status === "info").length} info
          </span>
        </div>
      )}

      {error && (
        <div className="mb-3 rounded-md border border-border bg-panel p-3 font-mono text-xs text-err">
          {error}
        </div>
      )}

      {/* warnings + fixes */}
      {warnings.length > 0 && (
        <section className="mb-4">
          <h2 className="mb-2 text-sm font-medium text-warn">
            Needs attention
          </h2>
          <ul className="divide-y divide-border overflow-hidden rounded-lg border border-border bg-panel">
            {warnings.map((r, i) => (
              <Row key={i} r={r} />
            ))}
          </ul>
        </section>
      )}

      {/* fix actions */}
      <section className="mb-4 rounded-lg border border-border bg-panel p-4">
        <h2 className="mb-2 text-sm font-medium text-ink">Fixes</h2>
        <div className="flex flex-wrap items-center gap-2">
          {hasDrift && (
            <>
              <button
                onClick={() =>
                  void run("inspect", () => ctxDrift(dir, false))
                }
                disabled={!!running}
                className="rounded-md border border-border bg-bg px-3 py-1.5 text-xs text-ink hover:border-accent disabled:opacity-50"
              >
                {running === "inspect" ? "Running…" : "Inspect drift"}
              </button>
              <ConfirmButton
                which="drift"
                label="Auto-fix drift"
                onConfirm={() =>
                  void run("fix", () => ctxDrift(dir, true), true)
                }
              />
            </>
          )}
          {hasSize && (
            <ConfirmButton
              which="compact"
              label="Compact (archive completed)"
              onConfirm={() =>
                void run("compact", () => ctxCompact(dir, true), true)
              }
            />
          )}
          {!hasDrift && !hasSize && (
            <span className="text-xs text-muted">
              No guided fixes for the current findings.
            </span>
          )}
        </div>
        <p className="mt-2 text-xs text-muted">
          Auto-fix and compact write to <code className="text-ink">.context/</code>{" "}
          through <code className="text-ink">ctx</code> (git-backed, reversible).
        </p>
      </section>

      {output && (
        <section className="mb-4">
          <h2 className="mb-2 text-sm font-medium text-ink">Output</h2>
          <pre className="overflow-auto rounded-lg border border-border bg-panel p-4 font-mono text-xs leading-relaxed text-ink">
            {output}
          </pre>
        </section>
      )}

      {/* all checks */}
      {others.length > 0 && (
        <details className="rounded-lg border border-border bg-panel">
          <summary className="cursor-pointer px-4 py-2 text-sm text-muted">
            All checks ({results.length})
          </summary>
          <ul className="divide-y divide-border border-t border-border">
            {others.map((r, i) => (
              <Row key={i} r={r} />
            ))}
          </ul>
        </details>
      )}
    </div>
  );
}
