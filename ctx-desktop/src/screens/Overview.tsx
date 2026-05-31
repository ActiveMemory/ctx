import { useCallback, useEffect, useState } from "react";
import {
  ctxInfo,
  ctxStatus,
  ctxTasks,
  ctxDecisions,
  ctxLearnings,
  type CtxInfo,
  type CtxStatus,
  type Task,
} from "../adapter/ctx";

// Default project = the ctx repo itself, so the shell shows real
// data on first launch. Editable in the field below.
const DEFAULT_DIR = "/Users/hamzaerbay/Code/ctx";

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

export default function Overview() {
  const [dir, setDir] = useState(DEFAULT_DIR);
  const [info, setInfo] = useState<CtxInfo | null>(null);
  const [status, setStatus] = useState<CtxStatus | null>(null);
  const [counts, setCounts] = useState<Counts | null>(null);
  const [loading, setLoading] = useState(false);
  const [errors, setErrors] = useState<string[]>([]);

  const load = useCallback(async (projectDir: string) => {
    setLoading(true);
    setErrors([]);
    const errs: string[] = [];

    setInfo(await ctxInfo());

    try {
      setStatus(await ctxStatus(projectDir));
    } catch (e) {
      setStatus(null);
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
    void load(DEFAULT_DIR);
  }, [load]);

  return (
    <div className="mx-auto max-w-3xl px-6 py-8">
      <header className="mb-6 flex items-baseline justify-between">
        <div>
          <h1 className="text-xl font-semibold text-ink">ctx Desktop</h1>
          <p className="text-sm text-muted">do you remember?</p>
        </div>
        {info && (
          <span
            className={`rounded-full px-3 py-1 font-mono text-xs ${
              info.found
                ? "bg-ok/15 text-ok"
                : "bg-err/15 text-err"
            }`}
            title={info.error ?? ""}
          >
            {info.found ? info.version : "ctx not found"}
          </span>
        )}
      </header>

      <div className="mb-6 flex gap-2">
        <input
          value={dir}
          onChange={(e) => setDir(e.target.value)}
          spellCheck={false}
          className="flex-1 rounded-md border border-border bg-panel px-3 py-2 font-mono text-sm text-ink outline-none focus:border-accent"
          placeholder="/path/to/project (parent of .context)"
        />
        <button
          onClick={() => void load(dir)}
          disabled={loading}
          className="rounded-md bg-accent px-4 py-2 text-sm font-medium text-bg disabled:opacity-50"
        >
          {loading ? "Loading…" : "Load"}
        </button>
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
            The list commands need the branch build of ctx. Run{" "}
            <code className="text-ink">make build &amp;&amp; sudo make install</code>{" "}
            from the ctx repo (feat/ctx-artifact-list-json), then Load again.
          </p>
        </section>
      )}
    </div>
  );
}
