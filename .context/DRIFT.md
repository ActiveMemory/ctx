# Drift Detection

## Automatic Checks

### Documentation Staleness

- [ ] `docs/cli-reference.md` should be newer than or same age as `internal/cli/*.go`
- [ ] `docs/context-files.md` should be newer than or same age as `internal/templates/*.md`

### Path References

- [ ] All paths in docs/ reference existing files

## Manual Review Triggers

- [ ] After adding/removing CLI commands, update `docs/cli-reference.md`
- [ ] After changing context file formats, update `docs/context-files.md`
- [ ] After adding AI tool integration, update `docs/integrations.md`

## Staleness Indicators

| File                    | Stale If           | Action            |
|-------------------------|--------------------|-------------------|
| `docs/cli-reference.md` | CLI source newer   | Review and update |
| `docs/context-files.md` | Templates changed  | Review and update |
| `docs/integrations.md`  | Hook logic changed | Review and update |
