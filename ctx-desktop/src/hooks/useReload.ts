import { useEffect, useState } from "react";
import { listen } from "@tauri-apps/api/event";

// Returns a counter that increments (debounced) whenever `event` fires.
// Add it to a load effect's deps to refetch on external writes. Defaults
// to "ctx-changed" (the active project); the dashboard passes
// "ctx-projects-changed" to react to writes in ANY project.
export function useReloadOnCtxChange(event = "ctx-changed"): number {
  const [n, setN] = useState(0);
  useEffect(() => {
    let unlisten: (() => void) | undefined;
    let timer: ReturnType<typeof setTimeout> | undefined;
    void listen(event, () => {
      clearTimeout(timer);
      timer = setTimeout(() => setN((x) => x + 1), 300);
    }).then((u) => {
      unlisten = u;
    });
    return () => {
      unlisten?.();
      clearTimeout(timer);
    };
  }, [event]);
  return n;
}
