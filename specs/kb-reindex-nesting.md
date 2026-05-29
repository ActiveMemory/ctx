# Nesting-Aware `ctx kb reindex`

`ctx kb reindex` rebuilds the `CTX:KB:TOPICS` managed block in
`.context/kb/index.md` by enumerating topic pages under
`.context/kb/topics/`. The enumeration is one level deep, so a kb
that groups its topics into subfolders silently reindexes to zero
topics and blanks the managed block.

## Problem

`reindex.ListTopics` (`internal/cli/kb/core/reindex/topic.go`) does a
single `os.ReadDir(topicsDir)` and keeps each immediate child
directory that has a `<child>/index.md`. A grouped layout stores
topics at `topics/<group>/<slug>/index.md`; the `<group>` directory
has no topic `index.md` of its own (or only a group-landing page),
so:

- the real topics (`<group>/<slug>`) are never visited, and
- reindex reports "reindexed 0 topic(s)" and **blanks** the managed
  block (observed live in the `things-wtf-dr` kb after it
  reorganized 49 topics into grouped folders).

The package already documents the intended model — `ListTopics`'
doc comment promises "slashes preserved for vendor-namespaced
topology" and the entry template renders `topics/<slug>/` — but the
scan never produces slashed slugs because it stops at depth 1.

## Design

Make topic enumeration recursive and group-landing-aware:

- Walk `topicsDir` recursively (depth-unbounded), recording every
  directory that directly contains a `TopicIndex` (`index.md`),
  keyed by its slash-separated path relative to the topics root.
- A recorded directory is a **topic** iff no other recorded
  directory is a strict descendant of it. A directory whose
  `index.md` sits above nested topics is a **group-landing**
  (orientation) page and is excluded from enumeration.
- Slugs are slash-joined (`<group>/<slug>`), matching what
  `ctx kb topic new "<group>/<slug>"` creates (`slug.Path`) and what
  the `topics/<slug>/` link template already expects. Returned
  sorted.

Flat (`topics/<slug>/index.md`), grouped
(`topics/<group>/<slug>/index.md`), and mixed layouts all work;
arbitrary nesting depth is supported. A non-existent topics
directory still yields an empty list (the block shows the
"no topics yet" placeholder), never an error-blank.

`RenderBlock` is unchanged: sorted slashed slugs cluster by group
prefix in the flat list, and the existing `topics/<slug>/` link
template already targets nested paths correctly. Emitting explicit
per-group headings would change the managed-block format and is
deliberately out of scope here (possible later enhancement).

The recursive helpers live in a sibling `scan.go` (all-unexported)
so `topic.go` keeps a single exported `ListTopics`, per the
mixed-visibility convention.

## Tests

`internal/cli/kb/core/reindex/topic_test.go`:

- empty / non-existent topics dir → nil, no error.
- flat topics → `[a, b]`.
- grouped topics → `[g1/t1, g1/t2, g2/t3]`.
- mixed flat + grouped.
- group-landing excluded: `topics/g/index.md` + `topics/g/t/index.md`
  → `[g/t]` only (not `g`).
- a directory without an `index.md` is not a topic.
- slugs sorted; slashes preserved.

`internal/cli/kb/core/reindex/block_test.go`:

- a nested slug renders `- [`g/t`](topics/g/t/)`.

## Acceptance

- A grouped kb reindexes to its real topic count (not 0) and the
  managed block lists every `<group>/<slug>` with a working link.
- `make lint` clean; `make test` green; new regression tests pin the
  recursive + group-landing behavior.
