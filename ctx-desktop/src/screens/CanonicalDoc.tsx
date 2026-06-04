import { useEffect, useState, type ReactNode } from "react";
import { ctxReadDoc } from "../adapter/ctx";
import { useReloadOnCtxChange } from "../hooks/useReload";

// Inline formatting: **bold** and `code`. Everything else is text.
function inline(text: string): ReactNode[] {
  const nodes: ReactNode[] = [];
  const re = /\*\*([^*]+)\*\*|`([^`]+)`/g;
  let last = 0;
  let k = 0;
  let m: RegExpExecArray | null;
  while ((m = re.exec(text))) {
    if (m.index > last) nodes.push(text.slice(last, m.index));
    if (m[1] !== undefined) {
      nodes.push(
        <strong key={k++} className="font-semibold text-ink">
          {m[1]}
        </strong>,
      );
    } else if (m[2] !== undefined) {
      nodes.push(
        <code
          key={k++}
          className="rounded bg-border/40 px-1 font-mono text-[12px] text-ink"
        >
          {m[2]}
        </code>,
      );
    }
    last = m.index + m[0].length;
  }
  if (last < text.length) nodes.push(text.slice(last));
  return nodes;
}

// Minimal Markdown renderer: headings, bullets, checkboxes, rules,
// and paragraphs, with inline bold/code. HTML comments (editorial
// notes) are stripped. Good enough for the canonical context docs.
function render(md: string): ReactNode[] {
  const clean = md.replace(/<!--[\s\S]*?-->/g, "");
  const lines = clean.split("\n");
  const out: ReactNode[] = [];
  let key = 0;
  let bullets: ReactNode[] | null = null;

  const flush = () => {
    if (bullets) {
      out.push(
        <ul key={key++} className="mb-3 space-y-1">
          {bullets}
        </ul>,
      );
      bullets = null;
    }
  };

  for (const raw of lines) {
    const line = raw.trimEnd();
    if (!line.trim()) {
      flush();
      continue;
    }
    const h = /^(#{1,4})\s+(.*)$/.exec(line);
    if (h) {
      flush();
      const lvl = h[1].length;
      const cls =
        lvl === 1
          ? "mt-1 mb-3 text-lg font-semibold text-ink"
          : lvl === 2
            ? "mt-5 mb-2 text-base font-semibold text-ink"
            : "mt-4 mb-1.5 text-sm font-semibold text-ink";
      out.push(
        <div key={key++} className={cls}>
          {inline(h[2])}
        </div>,
      );
      continue;
    }
    if (/^(-{3,}|\*{3,}|_{3,})$/.test(line)) {
      flush();
      out.push(<hr key={key++} className="my-4 border-border" />);
      continue;
    }
    const cb = /^[-*]\s+\[([ xX])\]\s+(.*)$/.exec(line);
    if (cb) {
      (bullets ??= []).push(
        <li key={key++} className="flex gap-2 text-sm text-ink">
          <span className="select-none text-muted">
            {cb[1].trim() ? "☑" : "☐"}
          </span>
          <span>{inline(cb[2])}</span>
        </li>,
      );
      continue;
    }
    const b = /^[-*]\s+(.*)$/.exec(line);
    if (b) {
      (bullets ??= []).push(
        <li key={key++} className="flex gap-2 text-sm text-ink">
          <span className="select-none text-muted">•</span>
          <span>{inline(b[1])}</span>
        </li>,
      );
      continue;
    }
    flush();
    out.push(
      <p key={key++} className="mb-2 text-sm leading-relaxed text-ink">
        {inline(line)}
      </p>,
    );
  }
  flush();
  return out;
}

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
          {render(content)}
        </div>
      )}
    </div>
  );
}
