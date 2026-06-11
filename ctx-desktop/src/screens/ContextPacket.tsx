import { useCallback, useEffect, useRef, useState } from "react";
import {
  ctxAgentPacket,
  ctxAgentMarkdown,
  type AgentPacket,
} from "../adapter/ctx";
import { useReloadOnCtxChange } from "../hooks/useReload";

const MIN = 1000;
const MAX = 16000;
const STEP = 500;

// Section arrays in display order, with labels.
const SECTIONS: [keyof AgentPacket, string][] = [
  ["constitution", "Constitution"],
  ["tasks", "Tasks"],
  ["conventions", "Conventions"],
  ["decisions", "Decisions"],
  ["learnings", "Learnings"],
  ["summaries", "Summaries"],
];

export default function ContextPacket({ dir }: { dir: string }) {
  const [budget, setBudget] = useState(8000);
  const [packet, setPacket] = useState<AgentPacket | null>(null);
  const [error, setError] = useState<string | null>(null);
  const [loading, setLoading] = useState(false);
  const [copied, setCopied] = useState<"" | "packet" | "command">("");
  const reload = useReloadOnCtxChange(dir);
  // Only the latest load applies — guards against an earlier budget/dir
  // load landing after a newer one (slider drags fire rapidly).
  const reqId = useRef(0);

  const load = useCallback(async (d: string, b: number) => {
    const id = ++reqId.current;
    setLoading(true);
    setError(null);
    try {
      const next = await ctxAgentPacket(d, b);
      if (id === reqId.current) setPacket(next);
    } catch (e) {
      if (id === reqId.current) {
        setError(String(e));
        setPacket(null);
      }
    } finally {
      if (id === reqId.current) setLoading(false);
    }
  }, []);

  // Debounce so dragging the slider doesn't spawn a call per pixel.
  useEffect(() => {
    const t = setTimeout(() => void load(dir, budget), 250);
    return () => clearTimeout(t);
  }, [dir, budget, load, reload]);

  async function copyPacket() {
    try {
      const md = await ctxAgentMarkdown(dir, budget);
      await navigator.clipboard.writeText(md);
      flash("packet");
    } catch (e) {
      setError(String(e));
    }
  }

  async function copyCommand() {
    await navigator.clipboard.writeText(`ctx agent --budget ${budget}`);
    flash("command");
  }

  function flash(which: "packet" | "command") {
    setCopied(which);
    setTimeout(() => setCopied(""), 1200);
  }

  const used = packet?.tokens_used ?? 0;
  const pct = packet ? Math.min(100, Math.round((used / budget) * 100)) : 0;
  const over = used > budget;

  return (
    <div className="mx-auto max-w-3xl px-6 py-6">
      <div className="mb-4 flex items-center justify-between">
        <h1 className="text-lg font-semibold text-ink">Context Packet</h1>
        <div className="flex gap-2">
          <button
            onClick={() => void copyPacket()}
            className="rounded-md border border-border bg-panel px-3 py-1.5 text-xs text-ink hover:border-accent"
          >
            {copied === "packet" ? "Copied ✓" : "Copy packet"}
          </button>
          <button
            onClick={() => void copyCommand()}
            className="rounded-md border border-border bg-panel px-3 py-1.5 text-xs text-ink hover:border-accent"
          >
            {copied === "command" ? "Copied ✓" : "Copy command"}
          </button>
        </div>
      </div>

      {/* budget control */}
      <div className="mb-5 rounded-lg border border-border bg-panel p-4">
        <div className="mb-2 flex items-center gap-3">
          <label className="text-sm text-muted">Budget</label>
          <input
            type="number"
            min={MIN}
            max={MAX}
            step={STEP}
            value={budget}
            onChange={(e) => setBudget(Number(e.target.value) || MIN)}
            className="w-28 rounded-md border border-border bg-bg px-2 py-1 font-mono text-sm text-ink outline-none focus:border-accent"
          />
          <span className="font-mono text-xs text-muted">tokens</span>
          {loading && (
            <span className="ml-auto text-xs text-muted">loading…</span>
          )}
        </div>
        <input
          type="range"
          min={MIN}
          max={MAX}
          step={STEP}
          value={budget}
          onChange={(e) => setBudget(Number(e.target.value))}
          className="w-full accent-accent"
        />
        {packet && (
          <div className="mt-3">
            <div className="h-2 overflow-hidden rounded-full bg-bg">
              <div
                className={`h-full ${over ? "bg-warn" : "bg-accent"}`}
                style={{ width: `${pct}%` }}
              />
            </div>
            <div className="mt-1 font-mono text-[11px] text-muted">
              ~{used.toLocaleString()} tokens used / {budget.toLocaleString()}{" "}
              budget
              {over && (
                <span className="text-warn">
                  {" "}
                  · over budget — lowest-scored entries become title-only
                </span>
              )}
            </div>
          </div>
        )}
      </div>

      {error && (
        <div className="mb-3 rounded-md border border-border bg-panel p-3 font-mono text-xs text-err">
          {error}
        </div>
      )}

      {packet && (
        <>
          {/* read order */}
          <section className="mb-4 rounded-lg border border-border bg-panel p-4">
            <h2 className="mb-2 text-sm font-medium text-ink">Read order</h2>
            <ol className="space-y-0.5 font-mono text-xs text-muted">
              {packet.read_order.map((f, i) => (
                <li key={i}>
                  {i + 1}. {f.split("/").pop()}
                </li>
              ))}
            </ol>
          </section>

          {/* sections with included/dropped affordance */}
          {SECTIONS.map(([key, label]) => {
            const items = (packet[key] as string[]) ?? [];
            return (
              <section
                key={key}
                className="mb-3 rounded-lg border border-border bg-panel p-4"
              >
                <div className="mb-2 flex items-center justify-between">
                  <h2 className="text-sm font-medium text-ink">{label}</h2>
                  <span
                    className={`font-mono text-[11px] ${
                      items.length ? "text-ok" : "text-muted"
                    }`}
                  >
                    {items.length
                      ? `${items.length} included`
                      : "dropped at this budget"}
                  </span>
                </div>
                {items.length > 0 && (
                  <ul className="space-y-0.5 font-mono text-xs text-muted">
                    {items.slice(0, 12).map((it, i) => (
                      <li key={i} className="truncate">
                        {it}
                      </li>
                    ))}
                    {items.length > 12 && (
                      <li className="text-muted">
                        …and {items.length - 12} more
                      </li>
                    )}
                  </ul>
                )}
              </section>
            );
          })}
        </>
      )}
    </div>
  );
}
