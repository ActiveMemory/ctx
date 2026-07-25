# workflow

## Code Change Heuristics

- **Present interpretations, don't pick silently**: If a request has multiple
  valid readings, lay them out rather than guessing
- **Push back when warranted**: If a simpler approach exists, say so
- **"Would a senior engineer call this overcomplicated?"**: If yes, simplify
- **Match existing style**: Even if you'd write it differently in a greenfield
- **Every changed line traces to the request**: If it doesn't, revert it

## Decision Heuristics

- **"Would I start this today?"**: If not, continuing is
  the sunk cost — evaluate only future value
- **"Reversible or one-way door?"**: Reversible decisions
  don't need deep analysis
- **"Does the analysis cost more than the decision?"**:
  Stop deliberating when the options are within an order
  of magnitude
- **"Order of magnitude, not precision"**: 10x better
  matters; 10% better usually doesn't

## Refactoring

- **Measure the end state, not the effort**: When refactoring, ask what the
  codebase looks like *after*, not how much work the change is
- **Three questions before restructuring**:
  1. What's the smallest codebase that solves this?
  2. Does the proposed change result in less total code?
  3. What can we delete now that this change makes obsolete?
- **Deletion is a feature**: Writing 50 lines that delete 200 is a net win

## Testing

- **Colocate tests**: Test files live next to source files
  - `foo.go` → `foo_test.go` in same package
  - Not a separate `tests/` folder
- **Test the unit, not the file**: One test file can test
  multiple related functions
- **Integration tests are separate**: `cli_test.go` for end-to-end binary tests

