import { type ReactNode } from "react";

// Inline formatting: **bold** and `code`. Everything else is text.
export function inline(text: string): ReactNode[] {
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
// notes) are stripped. Good enough for the context/kb docs.
export function renderMarkdown(md: string): ReactNode[] {
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
