import { useEffect, useMemo, useState } from "react";
import {
  ctxTasks,
  ctxDecisions,
  ctxLearnings,
  ctxReadDoc,
} from "../adapter/ctx";
import { useReloadOnCtxChange } from "../hooks/useReload";

// A target screen id the result links to (matches App's View ids).
type View =
  | "tasks"
  | "decisions"
  | "learnings"
  | "conventions"
  | "constitution";

interface Hit {
  kind: string;
  view: View;
  title: string;
  snippet: string;
}

// Cap matches per category so a broad query stays readable.
const PER_KIND = 25;

// Loaded once per project; searched in-memory as the user types.
interface Corpus {
  tasks: { text: string; section: string; priority: string }[];
  decisions: {
    title: string;
    context: string;
    rationale: string;
    consequence: string;
  }[];
  learnings: {
    title: string;
    context: string;
    lesson: string;
    application: string;
  }[];
  conventions: string[];
  constitution: string[];
}

const EMPTY: Corpus = {
  tasks: [],
  decisions: [],
  learnings: [],
  conventions: [],
  constitution: [],
};

// Non-empty, non-comment markdown lines, for line-level doc search.
function docLines(md: string): string[] {
  return md
    .replace(/<!--[\s\S]*?-->/g, "")
    .split("\n")
    .map((l) => l.trim())
    .filter((l) => l && !/^#{1,4}\s/.test(l) && !/^-{3,}$/.test(l));
}

function clip(s: string, q: string): string {
  const i = s.toLowerCase().indexOf(q);
  if (i < 0) return s.slice(0, 140);
  const start = Math.max(0, i - 40);
  return (start > 0 ? "…" : "") + s.slice(start, start + 140);
}

export default function Search({
  dir,
  onOpen,
}: {
  dir: string;
  onOpen: (view: View) => void;
}) {
  const [corpus, setCorpus] = useState<Corpus>(EMPTY);
  const [query, setQuery] = useState("");
  const [error, setError] = useState<string | null>(null);
  const reload = useReloadOnCtxChange(dir);

  useEffect(() => {
    let alive = true;
    setError(null);
    (async () => {
      const [tasks, decisions, learnings, conv, cons] = await Promise.all([
        ctxTasks(dir).catch(() => []),
        ctxDecisions(dir).catch(() => []),
        ctxLearnings(dir).catch(() => []),
        ctxReadDoc(dir, "CONVENTIONS.md").catch(() => ""),
        ctxReadDoc(dir, "CONSTITUTION.md").catch(() => ""),
      ]);
      if (!alive) return;
      setCorpus({
        tasks: tasks.map((t) => ({
          text: t.text,
          section: t.section,
          priority: t.priority,
        })),
        decisions: decisions.map((d) => ({
          title: d.title,
          context: d.context,
          rationale: d.rationale,
          consequence: d.consequence,
        })),
        learnings: learnings.map((l) => ({
          title: l.title,
          context: l.context,
          lesson: l.lesson,
          application: l.application,
        })),
        conventions: docLines(conv),
        constitution: docLines(cons),
      });
    })().catch((e) => alive && setError(String(e)));
    return () => {
      alive = false;
    };
  }, [dir, reload]);

  const hits = useMemo<Hit[]>(() => {
    const q = query.trim().toLowerCase();
    if (q.length < 2) return [];
    const out: Hit[] = [];
    const push = (
      kind: string,
      view: View,
      title: string,
      hay: string,
      snippetSrc: string,
    ) => {
      if (hay.toLowerCase().includes(q))
        out.push({ kind, view, title, snippet: clip(snippetSrc, q) });
    };

    let n = 0;
    for (const t of corpus.tasks) {
      if (n >= PER_KIND) break;
      const before = out.length;
      push("Task", "tasks", t.text, `${t.text} ${t.section}`, t.text);
      if (out.length > before) n++;
    }
    n = 0;
    for (const d of corpus.decisions) {
      if (n >= PER_KIND) break;
      const hay = `${d.title} ${d.context} ${d.rationale} ${d.consequence}`;
      const before = out.length;
      push("Decision", "decisions", d.title, hay, hay);
      if (out.length > before) n++;
    }
    n = 0;
    for (const l of corpus.learnings) {
      if (n >= PER_KIND) break;
      const hay = `${l.title} ${l.context} ${l.lesson} ${l.application}`;
      const before = out.length;
      push("Learning", "learnings", l.title, hay, hay);
      if (out.length > before) n++;
    }
    n = 0;
    for (const line of corpus.conventions) {
      if (n >= PER_KIND) break;
      const before = out.length;
      push("Convention", "conventions", "Conventions", line, line);
      if (out.length > before) n++;
    }
    n = 0;
    for (const line of corpus.constitution) {
      if (n >= PER_KIND) break;
      const before = out.length;
      push("Constitution", "constitution", "Constitution", line, line);
      if (out.length > before) n++;
    }
    return out;
  }, [corpus, query]);

  const BADGE: Record<string, string> = {
    Task: "bg-accent/15 text-accent",
    Decision: "bg-ok/15 text-ok",
    Learning: "bg-warn/15 text-warn",
    Convention: "bg-muted/15 text-muted",
    Constitution: "bg-err/15 text-err",
  };

  return (
    <div className="mx-auto max-w-3xl px-6 py-6">
      <h1 className="mb-4 text-lg font-semibold text-ink">Search</h1>

      <input
        autoFocus
        value={query}
        onChange={(e) => setQuery(e.target.value)}
        placeholder="Search tasks, decisions, learnings, conventions, constitution…"
        className="mb-3 w-full rounded-md border border-border bg-panel px-3 py-2 text-sm text-ink outline-none focus:border-accent"
      />

      {error && (
        <div className="mb-3 rounded-md border border-border bg-panel p-3 font-mono text-xs text-err">
          {error}
        </div>
      )}

      {query.trim().length >= 2 && (
        <div className="mb-2 text-xs text-muted">
          {hits.length} match{hits.length === 1 ? "" : "es"}
        </div>
      )}

      <div className="space-y-2">
        {query.trim().length >= 2 && hits.length === 0 && (
          <div className="rounded-lg border border-border bg-panel px-4 py-6 text-center text-sm text-muted">
            No matches.
          </div>
        )}
        {hits.map((h, i) => (
          <button
            key={i}
            onClick={() => onOpen(h.view)}
            className="block w-full rounded-lg border border-border bg-panel px-4 py-3 text-left hover:border-accent"
          >
            <div className="flex items-center gap-2">
              <span
                className={`shrink-0 rounded-full px-2 py-0.5 text-[11px] ${BADGE[h.kind] ?? "bg-muted/15 text-muted"}`}
              >
                {h.kind}
              </span>
              <span className="truncate text-sm font-medium text-ink">
                {h.title}
              </span>
            </div>
            <div className="mt-1 text-xs text-muted">{h.snippet}</div>
          </button>
        ))}
      </div>
    </div>
  );
}
