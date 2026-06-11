import { useCallback, useEffect, useState } from "react";
import { ctxPadList, ctxPadAdd, ctxPadRm } from "../adapter/ctx";
import { useReloadOnCtxChange } from "../hooks/useReload";

export default function Pad({ dir }: { dir: string }) {
  const [list, setList] = useState("");
  const [text, setText] = useState("");
  const [target, setTarget] = useState("");
  const [busy, setBusy] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const reload = useReloadOnCtxChange(dir);

  const load = useCallback(async (d: string) => {
    setError(null);
    try {
      setList(await ctxPadList(d));
    } catch (e) {
      setError(String(e));
      setList("");
    }
  }, []);

  useEffect(() => {
    void load(dir);
  }, [dir, load, reload]);

  async function run(fn: () => Promise<string>) {
    setBusy(true);
    setError(null);
    try {
      await fn();
      await load(dir);
    } catch (e) {
      setError(String(e));
    } finally {
      setBusy(false);
    }
  }

  return (
    <div className="mx-auto max-w-3xl px-6 py-6">
      <h1 className="mb-1 text-lg font-semibold text-ink">Scratchpad</h1>
      <p className="mb-4 text-xs text-muted">
        Short, AES-256-GCM-encrypted one-liners that travel with the project
        (<code className="text-ink">.context/scratchpad.enc</code>). Shown
        decrypted here.
      </p>

      <div className="mb-3 flex gap-2">
        <input
          value={text}
          onChange={(e) => setText(e.target.value)}
          onKeyDown={(e) => {
            if (e.key === "Enter" && text.trim()) {
              void run(() => ctxPadAdd(dir, text.trim())).then(() =>
                setText(""),
              );
            }
          }}
          placeholder="New entry…"
          className="flex-1 rounded-md border border-border bg-panel px-3 py-2 text-sm text-ink outline-none focus:border-accent"
        />
        <button
          onClick={() =>
            void run(() => ctxPadAdd(dir, text.trim())).then(() => setText(""))
          }
          disabled={busy || !text.trim()}
          className="rounded-md bg-accent px-4 py-2 text-sm font-medium text-bg disabled:opacity-50"
        >
          Add
        </button>
      </div>

      <div className="mb-4 flex items-center gap-2">
        <input
          value={target}
          onChange={(e) => setTarget(e.target.value)}
          placeholder="Remove #"
          className="w-28 rounded-md border border-border bg-panel px-3 py-1.5 text-xs text-ink outline-none focus:border-accent"
        />
        <button
          onClick={() =>
            void run(() => ctxPadRm(dir, target.trim())).then(() =>
              setTarget(""),
            )
          }
          disabled={busy || !target.trim()}
          className="rounded-md border border-border bg-bg px-3 py-1.5 text-xs text-ink hover:border-accent disabled:opacity-50"
        >
          Remove
        </button>
      </div>

      {error && (
        <div className="mb-3 rounded-md border border-border bg-panel p-3 font-mono text-xs text-err">
          {error}
        </div>
      )}

      <pre className="overflow-auto rounded-lg border border-border bg-panel p-4 font-mono text-xs leading-relaxed text-ink">
        {list.trim() || "Scratchpad is empty."}
      </pre>
    </div>
  );
}
