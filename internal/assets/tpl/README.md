# assets/tpl/

Sprintf-based format string templates that cannot currently be expressed
in the YAML text system (`commands/text/*.yaml`).

## Why these exist

The YAML `short:` / `long:` format supports single values looked up by
key. These templates contain multi-line Markdown with `fmt.Sprintf`
placeholders (`%s`, `%d`), which don't fit that model.

## How they will be replaced

A Go `text/template` rendering pipeline (tracked in TASKS.md:
"Migrate Sprintf-based templates (tpl_*.go) to Go text/template")
will replace these constants with proper `.tmpl` files in the
embedded assets. Once that pipeline exists:

1. Each `tpl_*.go` constant moves to an embedded `.tmpl` file
2. Callers use `template.Execute()` with typed data instead of `fmt.Sprintf`
3. Templates become localizable alongside the YAML text entries
4. This directory is deleted

## Adding a new template

If you need a multi-line format string with placeholders, add it here
with a comment explaining why it can't go to YAML (see `tpl_obsidian.go`
for the pattern). Prefer this over scattering `const` blocks in CLI
packages — at least they're collected in one place.
