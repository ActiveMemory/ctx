import { useEffect, useState } from "react";
import { ctxReadDoc } from "../adapter/ctx";
import { useReloadOnCtxChange } from "../hooks/useReload";
import { renderMarkdown } from "../lib/markdown";

// CanonicalDoc renders a read-only canonical context file
// (Constitution, Conventions) as formatted Markdown.
export default function CanonicalDoc({
  dir,
  title,
  file,
}: {
  dir: string;
  title: string;
  file: string;
}) {
  const [content, setContent] = useState("");
  const [error, setError] = useState<string | null>(null);
  const reload = useReloadOnCtxChange();

  useEffect(() => {
    setError(null);
    ctxReadDoc(dir, file)
      .then(setContent)
      .catch((e) => {
        setError(String(e));
        setContent("");
      });
  }, [dir, file, reload]);

  return (
    <div className="mx-auto max-w-3xl px-6 py-6">
      <h1 className="mb-4 text-lg font-semibold text-ink">{title}</h1>

      {error && (
        <div className="mb-3 rounded-md border border-border bg-panel p-3 font-mono text-xs text-err">
          {error}
        </div>
      )}

      {!error && !content.trim() && (
        <div className="rounded-lg border border-border bg-panel px-4 py-6 text-center text-sm text-muted">
          No {title} file in this project's <code>.context/</code>.
        </div>
      )}

      {content.trim() && (
        <div className="rounded-lg border border-border bg-panel px-5 py-4">
          {renderMarkdown(content)}
        </div>
      )}
    </div>
  );
}
