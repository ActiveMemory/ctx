# Tasks — Context CLI

## In Progress

## Phase 1

- [ ] Verify all Markdown files by "actually reading them"; take notes for
  follow-up actions.
- [ ] All go code should have godoc and testing.
- [ ] GitHub CI linter is giving errors that need fixing.
- [ ] Manual code review. take notes.
- [ ] Add tests per file.
- [ ] validate everything in the docs with a skeptical eye.
- [ ] consider the case where `ctx` is not called from within AI prompt:
  - does the command still make sense?
  - does it create the expected output?
- [ ] Cut the first release.


## Next Up

- [ ] Enforce test coverage targets in CI/Makefile #priority:medium #area:quality
  - internal/cli: 60% (currently 62.8%)
  - internal/context: 80% (currently 86.8%)
  - internal/drift: 80% (currently 88.0%)
  - internal/claude: 80% (currently 87.5%)
  - internal/templates: 80% (currently 88.9%)

## Completed (Recent)

## Blocked

## Reference

**Specs** (in `specs/` directory):
- `core-architecture.md` — Overall design philosophy
- `go-cli-implementation.md` — Go project structure and patterns
- `cli.md` — All CLI commands and their behavior
- `context-file-formats.md` — File format specifications
- `context-loader.md` — Loading and parsing logic
- `context-updater.md` — Update command handling

**Task Status Labels**:
- `[ ]` — pending
- `[x]` — completed
- `[-]` — skipped (with reason)
- `#in-progress` — currently being worked on (add inline, don't move task)

**Archives**: See `.context/archive/` for completed tasks from previous phases.
