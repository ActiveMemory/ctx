import { useCallback, useEffect, useRef, useState } from "react";
import {
  ctxStatus,
  ctxTasks,
  ctxDecisions,
  ctxLearnings,
  type CtxStatus,
  type Task,
} from "../adapter/ctx";
import { useReloadOnCtxChange } from "../hooks/useReload";

interface Counts {
  tasksOpen: number;
  tasksDone: number;
  decisions: number;
  learnings: number;
}

function Stat({ label, value }: { label: string; value: string | number }) {
  return (
    <div className="rounded-lg border border-border bg-panel px-4 py-3">
      <div className="font-mono text-2xl text-ink">{value}</div>
      <div className="mt-1 text-xs text-muted">{label}</div>
    </div>
  );
}

export default function Overview({ dir }: { dir: string }) {
  const [status, setStatus] = useState<CtxStatus | null>(null);
  const [counts, setCounts] = useState<Counts | null>(null);
  const [loading, setLoading] = useState(false);
  const [errors, setErrors] = useState<string[]>([]);
  const reload = useReloadOnCtxChange();
  // Monotonic request id: only the latest load applies its results, so a
  // fast project switch can't let a slow earlier load overwrite state.
  const reqId = useRef(0);

  const load = useCallback(async (projectDir: string) => {
    const id = ++reqId.current;
    setLoading(true);
    const errs: string[] = [];

    let status: CtxStatus | null = null;
    try {
      status = await ctxStatus(projectDir);
    } catch (e) {
      errs.push(`status: ${String(e)}`);
    }

    const [tasks, decisions, learnings] = await Promise.all([
      ctxTasks(projectDir).catch((e) => {
        errs.push(`task list: ${String(e)}`);
        return [] as Task[];
      }),
      ctxDecisions(projectDir).catch((e) => {
        errs.push(`decision list: ${String(e)}`);
        return [];
      }),
      ctxLearnings(projectDir).catch((e) => {
        errs.push(`learning list: ${String(e)}`);
        return [];
      }),
    ]);

    if (id !== reqId.current) return; // a newer load superseded this one
    setStatus(status);
    setCounts({
      tasksOpen: tasks.filter((t) => t.status === "pending").length,
      tasksDone: tasks.filter((t) => t.status === "done").length,
      decisions: decisions.length,
      learnings: learnings.length,
    });
    setErrors(errs);
    setLoading(false);
  }, []);

  useEffect(() => {
    void load(dir);
  }, [dir, load, reload]);

  return (
    <div className="mx-auto max-w-3xl px-6 py-6">
      <div className="mb-4 flex items-baseline justify-between">
        <h1 className="text-lg font-semibold text-ink">Overview</h1>
        {loading && <span className="text-xs text-muted">loading…</span>}
      </div>

      <section className="grid grid-cols-2 gap-3 sm:grid-cols-4">
        <Stat label="Open tasks" value={counts?.tasksOpen ?? "—"} />
        <Stat label="Done tasks" value={counts?.tasksDone ?? "—"} />
        <Stat label="Decisions" value={counts?.decisions ?? "—"} />
        <Stat label="Learnings" value={counts?.learnings ?? "—"} />
      </section>

      {status && (
        <section className="mt-3 grid grid-cols-2 gap-3">
          <Stat label="Context files" value={status.total_files} />
          <Stat
            label="Total tokens"
            value={status.total_tokens.toLocaleString()}
          />
        </section>
      )}

      {errors.length > 0 && (
        <section className="mt-6 rounded-lg border border-border bg-panel p-4">
          <h2 className="mb-2 text-sm font-medium text-warn">
            Some commands didn't return
          </h2>
          <ul className="space-y-1 font-mono text-xs text-muted">
            {errors.map((e, i) => (
              <li key={i}>{e}</li>
            ))}
          </ul>
          <p className="mt-3 text-xs text-muted">
            The task/decision/learning counts need the branch build of ctx. Run{" "}
            <code className="text-ink">make build &amp;&amp; sudo make install</code>{" "}
            from the ctx repo (feat/ctx-artifact-list-json), then reopen.
          </p>
        </section>
      )}
    </div>
  );
}
