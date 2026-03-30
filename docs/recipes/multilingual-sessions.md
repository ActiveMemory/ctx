---
#   /    ctx:                         https://ctx.ist
# ,'`./    do you remember?
# `.,'\
#   \    Copyright 2026-present Context contributors.
#                 SPDX-License-Identifier: Apache-2.0

title: "Multilingual Session Parsing"
icon: lucide/languages
---

![ctx](../images/ctx-banner.png)

## The Problem

Your team works across languages. Session files written by AI tools
might use headers like `# Oturum: 2026-01-15 - API Düzeltme` (Turkish)
or `# セッション: 2026-01-15 - テスト` (Japanese) instead of
`# Session: 2026-01-15 - Fix API`.

By default, ctx only recognizes `Session:` as a session header prefix.
Files with other prefixes are silently skipped during journal import and
journal generation: They look like regular Markdown, not sessions.

## TL;DR

Add recognized prefixes to `.ctxrc`:

```yaml
session_prefixes:
  - "Session:"      # English (include to keep default)
  - "Oturum:"       # Turkish
  - "セッション:"     # Japanese
```

Restart your session. All configured prefixes are now recognized.

## How It Works

The Markdown session parser detects session files by looking for an H1
header that starts with a known prefix followed by a date:

```markdown
# Session: 2026-01-15 - Fix API Rate Limiting
# Oturum: 2026-01-15 - API Düzeltme
# セッション: 2026-01-15 - テスト
```

The list of recognized prefixes comes from `session_prefixes` in
`.ctxrc`. When the key is absent or empty, ctx falls back to the
built-in default: `["Session:"]`.

Date-only headers (`# 2026-01-15 - Morning Work`) are always recognized
regardless of prefix configuration.

## Configuration

### Adding a language

Add the prefix with a trailing colon to your `.ctxrc`:

```yaml
session_prefixes:
  - "Session:"
  - "Sesión:"       # Spanish
```

!!! warning "Include Session: explicitly"
    When you override `session_prefixes`, **the default is replaced**,
    not extended. If you still want English headers recognized, include
    `"Session:"` in your list.

### Team setup

Commit `.ctxrc` to the repo so all team members share the same prefix
list. This ensures `ctx journal import` and journal generation pick up
sessions from all team members regardless of language.

### Common prefixes

| Language   | Prefix     |
|------------|------------|
| English    | `Session:` |
| Turkish    | `Oturum:`  |
| Spanish    | `Sesión:`  |
| French     | `Session:` |
| German     | `Sitzung:` |
| Japanese   | `セッション:`   |
| Korean     | `세션:`      |
| Portuguese | `Sessão:`  |
| Chinese    | `会话:`      |

### Verifying

After configuring, test with `ctx journal source`. Sessions with the new
prefixes should appear in the output.

## What This Does NOT Do

- **Change the interface language**: `ctx` output is always English.
  This setting only controls which session files ctx can *parse*.
- **Generate headers**: `ctx` never writes session headers. The prefix
  list is recognition-only (input, not output).
- **Affect JSONL sessions**: Claude Code JSONL transcripts don't use
  header prefixes. This only applies to Markdown session files in
  `.context/sessions/`.

## See Also

*See also: [Setup Across AI Tools](multi-tool-setup.md) - complete
multi-tool setup including Markdown session configuration.*

*See also: [CLI Reference](../cli/index.md) - full `.ctxrc` field
reference including `session_prefixes`.*
