# assets/tpl/

Sprintf-based format string templates that cannot currently be expressed
in the YAML text system (`commands/text/*.yaml`).

## Why these exist

Each file here contains multi-line Markdown with `fmt.Sprintf`
placeholders (`%s`, `%d`). While the YAML text system technically
supports format strings (many `hooks.yaml` entries use `%s`), templates
with 10+ placeholders and complex multi-line structure are more
readable and maintainable as Go constants than as YAML block scalars.

**If a template has zero format placeholders, it does not belong here.**
Move it to a YAML text entry instead.

## Remaining files

| File              | Placeholders | Why it's here |
|-------------------|-------------|---------------|
| `tpl_entry.go`    | 15          | ctx add entry templates (decision, learning, convention, task) |
| `tpl_journal.go`  | 26          | Journal markdown rendering |
| `tpl_loop.go`     | 15          | Shell script generation for autonomous loops |
| `tpl_obsidian.go` | 1           | Obsidian vault README (borderline — could migrate) |
| `tpl_recall.go`   | 21          | Recall output rendering |
| `tpl_trigger.go`  | 2           | Trigger script scaffold (borderline — could migrate) |

## How they will be replaced

A Go `text/template` rendering pipeline will replace these constants
with proper `.tmpl` files in the embedded assets. Once that pipeline
exists:

1. Each `tpl_*.go` constant moves to an embedded `.tmpl` file
2. Callers use `template.Execute()` with typed data instead of `fmt.Sprintf`
3. Templates become localizable alongside the YAML text entries
4. This directory is deleted

## Adding a new template

Before adding here, check whether the template has format placeholders.
If it doesn't, put it in `commands/text/write.yaml` instead. If it does,
add it here with a comment listing the placeholder arguments.
