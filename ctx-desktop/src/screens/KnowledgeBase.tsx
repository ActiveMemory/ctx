import { useEffect, useMemo, useState, type ReactNode } from "react";
import { kbInfo, kbRead, type KbInfo } from "../adapter/ctx";
import { useReloadOnCtxChange } from "../hooks/useReload";
import { renderMarkdown } from "../lib/markdown";

const DOC_LABEL: Record<string, string> = {
  "index.md": "Overview",
  "evidence-index.md": "Evidence Index",
  "source-map.md": "Source Map",
  "source-coverage.md": "Source Coverage",
  "grounding-sources.md": "Grounding Sources",
  "outstanding-questions.md": "Outstanding Questions",
};

function RailButton({
  active,
  onClick,
  children,
  indent,
}: {
  active: boolean;
  onClick: () => void;
  children: ReactNode;
  indent?: boolean;
}) {
  return (
    <button
      onClick={onClick}
      className={`mb-0.5 block w-full truncate rounded-md px-2 py-1.5 text-left text-xs ${
        indent ? "pl-4" : ""
      } ${
        active
          ? "bg-accent/15 text-accent"
          : "text-muted hover:bg-border/40 hover:text-ink"
      }`}
    >
      {children}
    </button>
  );
}

export default function KnowledgeBase({ dir }: { dir: string }) {
  const [info, setInfo] = useState<KbInfo | null>(null);
  const [sel, setSel] = useState("");
  const [content, setContent] = useState("");
  const [error, setError] = useState<string | null>(null);
  const reload = useReloadOnCtxChange();

  useEffect(() => {
    setError(null);
    setInfo(null);
    setSel("");
    kbInfo(dir)
      .then((i) => {
        setInfo(i);
        const first = i.docs[0]
          ? i.docs[0]
          : i.topics[0]
            ? `topics/${i.topics[0]}/index.md`
            : "";
        setSel(first);
      })
      .catch((e) => setError(String(e)));
  }, [dir, reload]);

  useEffect(() => {
    if (!sel) {
      setContent("");
      return;
    }
    kbRead(dir, sel)
      .then(setContent)
      .catch((e) => {
        setError(String(e));
        setContent("");
      });
  }, [dir, sel, reload]);

  // Topics grouped by parent path (the grouped layout from reindex).
  const groups = useMemo(() => {
    const map: Record<string, { slug: string; leaf: string }[]> = {};
    const order: string[] = [];
    for (const slug of info?.topics ?? []) {
      const cut = slug.lastIndexOf("/");
      const group = cut < 0 ? "" : slug.slice(0, cut);
      const leaf = cut < 0 ? slug : slug.slice(cut + 1);
      if (!(group in map)) {
        map[group] = [];
        order.push(group);
      }
      map[group].push({ slug, leaf });
    }
    return { map, order };
  }, [info]);

  if (info && !info.exists) {
    return (
      <div className="mx-auto max-w-3xl px-6 py-6">
        <h1 className="mb-4 text-lg font-semibold text-ink">Knowledge Base</h1>
        <div className="rounded-lg border border-border bg-panel px-4 py-8 text-center text-sm text-muted">
          This project has no <code>.context/kb/</code>. Create a topic with{" "}
          <code className="text-ink">ctx kb topic new "&lt;name&gt;"</code>.
        </div>
      </div>
    );
  }

  return (
    <div className="flex h-full min-h-0">
      <aside className="w-56 shrink-0 overflow-auto border-r border-border bg-panel px-2 py-4">
        {info?.docs.length ? (
          <>
            <div className="px-2 pb-1 text-[10px] uppercase tracking-wide text-muted">
              Documents
            </div>
            {info.docs.map((d) => (
              <RailButton key={d} active={sel === d} onClick={() => setSel(d)}>
                {DOC_LABEL[d] ?? d}
              </RailButton>
            ))}
          </>
        ) : null}

        {info?.topics.length ? (
          <div className="mt-3">
            <div className="px-2 pb-1 text-[10px] uppercase tracking-wide text-muted">
              Topics
            </div>
            {groups.order.map((g) => (
              <div key={g || "_"} className="mb-1">
                {g && (
                  <div className="px-2 pt-1 text-[11px] font-medium text-ink">
                    {g}
                  </div>
                )}
                {groups.map[g].map(({ slug, leaf }) => {
                  const rel = `topics/${slug}/index.md`;
                  return (
                    <RailButton
                      key={slug}
                      active={sel === rel}
                      onClick={() => setSel(rel)}
                      indent={!!g}
                    >
                      {leaf}
                    </RailButton>
                  );
                })}
              </div>
            ))}
          </div>
        ) : null}
      </aside>

      <div className="min-w-0 flex-1 overflow-auto">
        <div className="mx-auto max-w-3xl px-6 py-6">
          {error && (
            <div className="mb-3 rounded-md border border-border bg-panel p-3 font-mono text-xs text-err">
              {error}
            </div>
          )}
          {!error && !content.trim() && (
            <div className="rounded-lg border border-border bg-panel px-4 py-6 text-center text-sm text-muted">
              {info ? "Select an item from the left." : "Loading…"}
            </div>
          )}
          {content.trim() && (
            <div className="rounded-lg border border-border bg-panel px-5 py-4">
              {renderMarkdown(content)}
            </div>
          )}
        </div>
      </div>
    </div>
  );
}
