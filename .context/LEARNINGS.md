# Learnings

<!--
UPDATE WHEN:
- Discover a gotcha, bug, or unexpected behavior
- Debugging reveals non-obvious root cause
- External dependency has quirks worth documenting
- "I wish I knew this earlier" moments
- Production incidents reveal gaps

DO NOT UPDATE FOR:
- Well-documented behavior (link to docs instead)
- Temporary workarounds (use TASKS.md for follow-up)
- Opinions without evidence
-->

<!-- INDEX:START -->
| Date | Learning |
|------|--------|
| 2026-04-26 | OpenCode auto-loads only flat .ts files under .opencode/plugins/; subdirectories are ignored |
| 2026-04-26 | OpenCode opencode.json MCP shape: command is Array<string>, no separate args field |
| 2026-04-26 | make test exit code unreliable due to -cover covdata tooling issue |
| 2026-04-26 | Trailing word boundary in regex matches commit-tree as git commit |
| 2026-04-26 | ctx system help can list project-local hooks not in the Go binary |
| 2026-04-25 | Confident code comments can pull an LLM away from first-principles knowledge |
| 2026-04-25 | filepath.Join('', rel) returns rel as CWD-relative, not error |
| 2026-04-25 | Parallel go test ./... packages can race on ~/.claude/settings.json |
<!-- INDEX:END -->

<!-- Add gotchas, tips, and lessons learned here -->
## [2026-04-26-180000] OpenCode auto-loads only flat .ts files under .opencode/plugins/; subdirectories are ignored

**Context**: Initial OpenCode integration deployed the plugin as `.opencode/plugins/ctx/index.ts` (a directory with index.ts inside, mirroring npm package conventions). End-to-end smoke testing showed the plugin file was present and the binary was current, yet OpenCode never invoked any of the plugin's hooks (no `module-load` trace fired even with `--print-logs --log-level DEBUG`). Copying the same content to a flat `.opencode/plugins/ctx.ts` file made the plugin load and fire correctly.

**Lesson**: OpenCode's plugin auto-discovery only scans top-level files under `.opencode/plugins/` and `~/.config/opencode/plugins/`. Subdirectories are silently skipped — there is no log line indicating a subdirectory was found and ignored. The official docs at opencode.ai/docs/plugins/ say only "files in these directories are automatically loaded at startup" without specifying the rule, so this is easy to miss. The `opencode plugin <module>` CLI registers npm modules (a different code path) and accepts only npm names, not local paths.

**Application**: Deploy single-file plugins as `.opencode/plugins/<name>.ts`, not `.opencode/plugins/<name>/index.ts`. No `package.json` is required when the plugin uses type-only imports (`import type` is erased at compile time) and the host runtime injects the plugin context. To verify a plugin is actually loaded, add a top-of-module side effect (e.g. `appendFileSync` to a known path) and confirm it fires before debugging hook contracts.

---

## [2026-04-26-165500] OpenCode opencode.json MCP shape: command is Array<string>, no separate args field

**Context**: `ctx setup opencode --write` was generating `opencode.json` with the Copilot CLI MCP shape (`{type: "local", command: "ctx", args: ["mcp", "serve"]}`). OpenCode rejected the file at startup with `Configuration is invalid… Expected array, got "ctx" mcp.ctx.command` and `Missing key mcp.ctx.enabled`.

**Lesson**: OpenCode's `McpLocalConfig` (in `@opencode-ai/sdk`) defines `command: Array<string>` as a single field that holds the binary AND its arguments — there is no separate `args` field. It also requires `enabled: boolean` at runtime even though the TS type marks it optional. The Copilot CLI MCP shape is similar in spirit but structurally different; do not copy-paste between them.

**Application**: For OpenCode MCP entries always use `command: ["ctx", "mcp", "serve"]` and include `enabled: true`. If you add a new editor integration with its own MCP file format, read the upstream type definitions from `node_modules/@<vendor>/sdk/dist/gen/types.gen.d.ts` (or equivalent) before reusing an existing generator.

---

## [2026-04-26-152850] make test exit code unreliable due to -cover covdata tooling issue

**Context**: make test exited 1 even with all 123 packages passing on this Go install; root cause is missing covdata tool when -cover is enabled

**Lesson**: Don't trust make test exit code alone when verifying changes. The -cover flag in the test target can fail with 'no such tool covdata' even when every package passes.

**Application**: When make test fails, fall back to 'go test ./...' (no -cover) and tally ^ok / ^FAIL counts to distinguish real failures from tooling issues.

---

## [2026-04-26-152842] Trailing word boundary in regex matches commit-tree as git commit

**Context**: First post-commit filter regex \bgit\s+commit\b in the OpenCode plugin would have triggered on git commit-tree because \b matches between t and -

**Lesson**: A trailing word boundary doesn't exclude hyphenated continuations — \b matches every word/non-word transition. Use (?!-) negative lookahead to specifically reject hyphen-suffixed siblings.

**Application**: For any porcelain with hyphenated cousins (commit-tree, commit-graph, for-each-ref), append (?!-) to the boundary.

---

## [2026-04-26-152836] ctx system help can list project-local hooks not in the Go binary

**Context**: PR #72 plugin called 'ctx system block-dangerous-commands'; user's installed ctx 0.7.2 listed it in help, but no directory exists under internal/cli/system/cmd/ — it's a Claude Code plugin-local hook surfaced via wrapper

**Lesson**: ctx system help output is a union of compiled Go subcommands and project-local Claude wrappers; non-Claude integrations only see the Go subset

**Application**: When porting plugin behavior to a new editor, only call subcommands that have a directory under internal/cli/system/cmd/. Don't trust ctx system help output as the canonical surface.

---

## [2026-04-25-014704] Confident code comments can pull an LLM away from first-principles knowledge

**Context**: cli_test.go had a comment claiming 'parent's t.Setenv doesn't propagate to exec'd children unless we build it into cmd.Env' which is wrong. I patched the helper's CTX_DIR dedup instead of questioning the helper itself, despite knowing t.Setenv semantics.

**Lesson**: A comment that explains why a stdlib mechanism 'doesn't work' is doing extra rhetorical work to talk a reader out of the obvious approach. That's exactly when to verify from first principles instead of trusting the surrounding-code frame.

**Application**: When an existing comment justifies a non-canonical approach contradicting stdlib knowledge: pause, verify against memory of the actual API before patching within the existing frame.

---

## [2026-04-25-014704] filepath.Join('', rel) returns rel as CWD-relative, not error

**Context**: Recurring orphan jsonl-path-<sessionID> appeared at project root. Older state.Dir() returned ('', nil) when CTX_DIR was undeclared, so filepath.Join('', 'jsonl-path-XXX') = 'jsonl-path-XXX', writing relative to CWD.

**Lesson**: Functions returning a path-string must never return ('', nil). Sentinel errors force callers to gate, closing the silent CWD-relative write.

**Application**: Audit any (string, error) path-returner that historically had a ('', nil) shortcut. Closed for state.Dir and rc.ContextDir; check remaining resolvers.

---

## [2026-04-25-014704] Parallel go test ./... packages can race on ~/.claude/settings.json

**Context**: make test runs packages in parallel processes. Fourteen test files invoked initialize.Cmd().Execute(), which read-modify-writes ~/.claude/settings.json without HOME isolation.

**Lesson**: Under load the races materialized as flaky 'FAIL coverage: [no statements]' in cli/watch/core. Run alone the package passed; under parallel make test it failed intermittently.

**Application**: testctx.Declare now sets HOME alongside CTX_DIR. Centralized fix; future tests automatically isolate user-home writes.
