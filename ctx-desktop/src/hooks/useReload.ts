import { useEffect, useState } from "react";
import { listen } from "@tauri-apps/api/event";

// Returns a counter that increments (debounced) whenever the active
// project's .context/ changes. Add it to a load effect's deps to
// refetch when the CLI or an AI agent writes externally.
export function useReloadOnCtxChange(): number {
  const [n, setN] = useState(0);
  useEffect(() => {
    let unlisten: (() => void) | undefined;
    let timer: ReturnType<typeof setTimeout> | undefined;
    void listen("ctx-changed", () => {
      clearTimeout(timer);
      timer = setTimeout(() => setN((x) => x + 1), 300);
    }).then((u) => {
      unlisten = u;
    });
    return () => {
      unlisten?.();
      clearTimeout(timer);
    };
  }, []);
  return n;
}
