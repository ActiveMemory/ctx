import { type ReactNode } from "react";
import { openUrl } from "@tauri-apps/plugin-opener";

// Inline formatting: **bold**, `code`, and [text](url) links. Links open
// in the system browser (not the webview) via the opener plugin, so an
// external URL can't navigate the app or trip the CSP. Everything else
// is plain text.
export function inline(text: string): ReactNode[] {
  const nodes: ReactNode[] = [];
  const re = /\*\*([^*]+)\*\*|`([^`]+)`|\[([^\]]+)\]\(([^)]+)\)/g;
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
    } else if (m[3] !== undefined && m[4] !== undefined) {
      const label = m[3];
      const url = m[4];
      nodes.push(
        <a
          key={k++}
          href={url}
          onClick={(e) => {
            e.preventDefault();
            void openUrl(url).catch(() => {});
          }}
          className="cursor-pointer text-accent underline hover:opacity-80"
        >
          {label}
        </a>,
      );
    }
    last = m.index + m[0].length;
  }
  if (last < text.length) nodes.push(text.slice(last));
  return nodes;
}

// Splits a table row into trimmed cells, dropping the outer pipes.
function tableCells(line: string): string[] {
  let s = line.trim();
  if (s.startsWith("|")) s = s.slice(1);
  if (s.endsWith("|")) s = s.slice(0, -1);
  return s.split("|").map((c) => c.trim());
}

// A row of a GitHub-style table: contains a pipe and isn't a fence.
function isTableRow(line: string): boolean {
  return line.includes("|") && !line.trimStart().startsWith("```");
}

// The |---|:--:| separator line under a table header.
function isTableSep(line: string): boolean {
  const cells = tableCells(line);
  return (
    cells.length > 0 && cells.every((c) => /^:?-+:?$/.test(c.replace(/\s/g, "")))
  );
}

// Column alignment classes parsed from the separator row.
function alignClass(sep: string): string {
  const c = sep.replace(/\s/g, "");
  const l = c.startsWith(":");
  const r = c.endsWith(":");
  return l && r ? "text-center" : r ? "text-right" : "text-left";
}

// Minimal Markdown renderer: headings, tables, code fences, bullets,
// checkboxes, rules, links, and paragraphs, with inline bold/code/link.
// HTML comments (editorial notes) are stripped. Good enough for the
// context/kb docs.
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

  let i = 0;
  while (i < lines.length) {
    const raw = lines[i];
    const line = raw.trimEnd();

    // Fenced code block: collect verbatim until the closing fence.
    if (line.trimStart().startsWith("```")) {
      flush();
      const body: string[] = [];
      i++;
      while (i < lines.length && !lines[i].trimStart().startsWith("```")) {
        body.push(lines[i]);
        i++;
      }
      i++; // consume closing fence
      out.push(
        <pre
          key={key++}
          className="mb-3 overflow-auto rounded-md border border-border bg-bg p-3 font-mono text-xs leading-relaxed text-ink"
        >
          {body.join("\n")}
        </pre>,
      );
      continue;
    }

    if (!line.trim()) {
      flush();
      i++;
      continue;
    }

    // Table: a row followed by a |---| separator.
    if (
      isTableRow(line) &&
      i + 1 < lines.length &&
      isTableSep(lines[i + 1].trimEnd())
    ) {
      flush();
      const header = tableCells(line);
      const aligns = tableCells(lines[i + 1].trimEnd()).map(alignClass);
      i += 2;
      const rows: string[][] = [];
      while (i < lines.length && isTableRow(lines[i].trimEnd()) && lines[i].trim()) {
        rows.push(tableCells(lines[i].trimEnd()));
        i++;
      }
      out.push(
        // Scroll a too-wide table inside its own box instead of letting
        // it blow out the surrounding doc width.
        <div key={key++} className="mb-3 overflow-x-auto">
          <table className="w-full border-collapse text-sm">
            <thead>
              <tr>
                {header.map((h, ci) => (
                  <th
                    key={ci}
                    className={`break-words border border-border px-2 py-1 align-top font-semibold text-ink ${aligns[ci] ?? "text-left"}`}
                  >
                    {inline(h)}
                  </th>
                ))}
              </tr>
            </thead>
            <tbody>
              {rows.map((row, ri) => (
                <tr key={ri}>
                  {row.map((cell, ci) => (
                    <td
                      key={ci}
                      className={`break-words border border-border px-2 py-1 align-top text-ink ${aligns[ci] ?? "text-left"}`}
                    >
                      {inline(cell)}
                    </td>
                  ))}
                </tr>
              ))}
            </tbody>
          </table>
        </div>,
      );
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
      i++;
      continue;
    }
    if (/^(-{3,}|\*{3,}|_{3,})$/.test(line)) {
      flush();
      out.push(<hr key={key++} className="my-4 border-border" />);
      i++;
      continue;
    }
    // Bullets/checkboxes, with light indentation for nested items.
    const indent = raw.length - raw.trimStart().length;
    const pad = Math.min(3, Math.floor(indent / 2)) * 14;
    const cb = /^[-*]\s+\[([ xX])\]\s+(.*)$/.exec(line.trimStart());
    if (cb) {
      (bullets ??= []).push(
        <li
          key={key++}
          className="flex gap-2 text-sm text-ink"
          style={pad ? { marginLeft: pad } : undefined}
        >
          <span className="select-none text-muted">
            {cb[1].trim() ? "☑" : "☐"}
          </span>
          <span>{inline(cb[2])}</span>
        </li>,
      );
      i++;
      continue;
    }
    const b = /^[-*]\s+(.*)$/.exec(line.trimStart());
    if (b) {
      (bullets ??= []).push(
        <li
          key={key++}
          className="flex gap-2 text-sm text-ink"
          style={pad ? { marginLeft: pad } : undefined}
        >
          <span className="select-none text-muted">•</span>
          <span>{inline(b[1])}</span>
        </li>,
      );
      i++;
      continue;
    }
    flush();
    out.push(
      <p key={key++} className="mb-2 text-sm leading-relaxed text-ink">
        {inline(line)}
      </p>,
    );
    i++;
  }
  flush();
  return out;
}
