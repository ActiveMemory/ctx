import { useEffect, useRef, useState } from "react";
import { listen } from "@tauri-apps/api/event";

// Returns a counter that increments (debounced) whenever the active
// project's "ctx-changed" event fires. Add it to a load effect's deps
// to refetch on external writes.
//
// The backend emits the project root as the event payload; pass
// `forDir` to ignore events for other projects (a payload-less event
// — older backend — always counts, so filtering is fail-open).
export function useReloadOnCtxChange(
  forDir?: string,
  event = "ctx-changed",
): number {
  const [n, setN] = useState(0);
  useEffect(() => {
    let unlisten: (() => void) | undefined;
    let timer: ReturnType<typeof setTimeout> | undefined;
    void listen<string>(event, (e) => {
      if (forDir && e.payload && e.payload !== forDir) return;
      clearTimeout(timer);
      timer = setTimeout(() => setN((x) => x + 1), 300);
    }).then((u) => {
      unlisten = u;
    });
    return () => {
      unlisten?.();
      clearTimeout(timer);
    };
  }, [event, forDir]);
  return n;
}

// Calls `onChange(dir)` whenever any watched project's `.context/`
// mutates ("ctx-projects-changed"), debounced PER project dir so a
// burst of writes in one project collapses to one callback without
// delaying another project's. A payload-less event (older backend)
// fires with "" — callers should treat that as "refresh everything".
// The latest `onChange` is always used; it need not be memoized.
export function useProjectsChanged(onChange: (dir: string) => void): void {
  const cb = useRef(onChange);
  cb.current = onChange;
  useEffect(() => {
    let unlisten: (() => void) | undefined;
    const timers = new Map<string, ReturnType<typeof setTimeout>>();
    void listen<string>("ctx-projects-changed", (e) => {
      const dir = e.payload || "";
      clearTimeout(timers.get(dir));
      timers.set(
        dir,
        setTimeout(() => {
          timers.delete(dir);
          cb.current(dir);
        }, 300),
      );
    }).then((u) => {
      unlisten = u;
    });
    return () => {
      unlisten?.();
      for (const t of timers.values()) clearTimeout(t);
    };
  }, []);
}
