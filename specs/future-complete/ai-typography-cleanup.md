# AI Typography Cleanup

## Problem

LLM-generated text uses typographic punctuation from its training
corpus: em-dashes, en-dashes, smart quotes, and space-padded
double hyphens. These leak into docs and Go doc comments during
AI-assisted writing sessions. The detection script
(hack/detect-ai-typography.sh) was silently broken on macOS
because BSD grep does not support the -P flag.

## Approach

1. Fix detect-ai-typography.sh to work on both GNU grep and BSD
   grep (macOS) via runtime detection
2. Replace all AI typography in docs/ with contextually appropriate
   ASCII punctuation (semantic editing, not blind sed)
3. Replace all AI typography in internal/ Go and YAML files
4. Remove hack/agents/ (content migrated to docs/operations/runbooks/)

## Scope

- docs/: ~293 em-dashes across 44 files, plus en-dashes and
  double-hyphens
- internal/: ~1668 matches across 328 files (mostly doc.go from
  PR #69 enrichment)
- hack/detect-ai-typography.sh: BSD grep compatibility fix
- hack/agents/: deletion (migrated)

## Decision

Spec also covers the ctx backup deprecation planning (spec and
decision record only, no code removal yet):
- specs/deprecate-ctx-backup.md created
- Decision recorded in DECISIONS.md
- TASKS.md updated with deprecation task
