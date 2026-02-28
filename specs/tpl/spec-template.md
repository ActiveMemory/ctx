# Spec Template

Use this template when writing specifications for ctx features.
Copy to a new file: `specs/{feature-name}.md`

---

# {Feature Name}

## Problem

What user-visible problem does this solve? Why does it matter now?

## Approach

High-level design. How does this work? Where does it fit in the
existing architecture?

## Behavior

### Happy Path

Step-by-step: what happens when everything goes right.

### Edge Cases

Enumerate explicitly. For each:

| Case | Expected behavior |
|------|-------------------|
| Empty input | ... |
| Partial failure (e.g., network drop mid-operation) | ... |
| Concurrent access | ... |
| Already exists / duplicate | ... |
| Missing dependencies | ... |

### Validation Rules

What input constraints are enforced? Where is validation performed?

### Error Handling

| Error condition | User-facing message | Recovery |
|-----------------|---------------------|----------|
| ... | ... | ... |

## Interface

### CLI (if applicable)

```
ctx {command} [flags]
```

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| ... | ... | ... | ... |

### Skill (if applicable)

```
/ctx-{name}
```

Trigger phrases: ...

## Implementation

### Files to Create/Modify

| File | Change |
|------|--------|
| ... | ... |

### Key Functions

```go
func ... { ... }
```

### Helpers to Reuse

List existing code to leverage. Avoid reinventing.

## Configuration

Any `.ctxrc` keys, environment variables, or settings.

## Testing

- Unit: ...
- Integration: ...
- Edge cases: ...

## Non-Goals

What this intentionally does NOT do. Be explicit.

## Open Questions

Unresolved items. Remove this section when all are answered.
