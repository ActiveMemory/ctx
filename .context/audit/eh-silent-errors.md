# EH.1 — Silent Error Discard Catalogue

Audit of every error-discard site under `internal/` (non-test).
Generated for Phase EH (Error Handling Audit). Every `_ =`,
`_, _ =`, and `x, _ :=` site is surfaced and categorised — no
triage that quietly drops some. Recommended actions drive EH.2–EH.5.

Method: `grep -rnE '(_ :?= |, _ :?= |_, _ :?= )' internal/ --include='*.go'`
(excluding `_test.go`), then per-site classification by the discarded
value's type and the call's role. Tests are out of scope for this pass.

## Category legend

| Tag | Meaning | Action |
|-----|---------|--------|
| `B-data` | error dropped on a data-write path; corruption/loss | **return the error** |
| `NIL-DEREF` | error dropped then result dereferenced | **guard before deref** |
| `B-marshal` | Marshal/Unmarshal error dropped; bad data written/read | **surface** |
| `B-parse` | Parse error dropped; zero value used silently | **surface** |
| `A-close-WRITE` | defer-close on a write/append handle; flush-fail = data loss | **surface close error** |
| `A-close-read` | defer-close on a read handle | **stderr-log (EH.3)** |
| `A-close-bare` | non-defer close discard | **surface or stderr** |
| `SURFACE` | error dropped; stale/empty result flows on | **surface/stderr** |
| `C-state` | os.Remove/Rename/Shutdown dropped; stale state | **stderr warn (EH.4)** |
| `D-output` | fmt.Fprint to cmd out/err | **convert to cmd.Print* (drops the discard)** |
| `besteffort` | telemetry/display/callback; failure is tolerable | annotate intent |
| `OK*` | discarded value is `ok` bool / nil-safe / init-time programmer-error | annotate |
| `FALSE-POS` | discarded value is not an error (string, bool, compile assertion) | none |

## Summary

Counts by category (184 sites total):

| Category | Count |
|----------|------:|
| D-output | 47 |
| FALSE-POS | 21 |
| OK-flag | 20 |
| besteffort | 16 |
| A-close-read | 13 |
| A-close-WRITE | 11 |
| OK-typeassert | 10 |
| B-marshal | 10 |
| OK-markflag | 7 |
| C-state | 7 |
| OK-atoi | 4 |
| OK | 4 |
| SURFACE | 3 |
| OK-glob | 3 |
| B-data | 3 |
| A-close-bare | 3 |
| NIL-DEREF | 1 |
| B-parse | 1 |

### High-priority findings (data loss / crash)

These are not style nits — they silently lose or corrupt data, or crash:

- `internal/memory/publish.go:170` — `MergePublished` error dropped, then the
  (possibly empty) `merged` is written over the published memory file.
- `internal/cli/pad/core/store/store.go:257` — `ReadEntriesWithIDs` error
  dropped; existing pad IDs are lost when the store is rewritten (the same
  data-loss shape as `specs/fix-learning-add-index-data-loss.md`).
- `internal/hub/replicate.go:121` — `store.Append` error dropped; a replicated
  hub entry is silently lost.
- `internal/cli/memory/cmd/status/run.go:54` — `LoadState` error dropped, then
  `state.LastSync` dereferenced: potential nil-pointer panic.
- 11 × `A-close-WRITE` — defer-close on `SafeAppendFile`/`SafeCreateFile`
  handles, where a failed final flush silently loses the appended row.

## Full catalogue

| Category | Site | Expression | Recommended action |
|----------|------|------------|--------------------|
| B-data | `internal/cli/pad/core/store/store.go:257` | `existing, _ := ReadEntriesWithIDs()` | ReadEntriesWithIDs error dropped; existing IDs lost on rewrite — surface/return |
| B-data | `internal/hub/replicate.go:121` | `_, _ = store.Append([]Entry{entry})` | store.Append error dropped; replicated entry silently lost — stderr/return |
| B-data | `internal/memory/publish.go:170` | `merged, _ := MergePublished(string(existing), formatted)` | MergePublished error dropped; (possibly empty) merged then written — return the error |
| NIL-DEREF | `internal/cli/memory/cmd/status/run.go:54` | `state, _ := mem.LoadState(contextDir)` | LoadState error dropped then state.LastSync deref — verify state non-nil; potential panic |
| B-marshal | `internal/cli/initialize/core/vscode/extension.go:59` | `data, _ := json.MarshalIndent(content, "", token.Indent2)` | surface: empty/partial data written on failure |
| B-marshal | `internal/cli/initialize/core/vscode/mcp.go:50` | `data, _ := json.MarshalIndent(file, "", token.Indent2)` | surface: empty/partial data written on failure |
| B-marshal | `internal/cli/initialize/core/vscode/tasks.go:59` | `data, _ := json.MarshalIndent(file, "", token.Indent2)` | surface: empty/partial data written on failure |
| B-marshal | `internal/cli/setup/core/copilot/vscode.go:57` | `data, _ := json.MarshalIndent(mcpCfg, "", token.Indent2)` | surface: empty/partial data written on failure |
| B-marshal | `internal/cli/system/cmd/blocknonpathctx/run.go:82` | `data, _ := json.Marshal(resp)` | surface: empty/partial data written on failure |
| B-marshal | `internal/cli/system/core/session/session.go:91` | `_ = json.Unmarshal(res.data, &input)` | json.Unmarshal of hook stdin; best-effort w/ timeout — annotate, or stderr |
| B-marshal | `internal/cli/trigger/cmd/test/cmd.go:109` | `inputJSON, _ := json.MarshalIndent(input, "", token.Indent2)` | surface: empty/partial data written on failure |
| B-marshal | `internal/steering/format.go:147` | `raw, _ := yaml.Marshal(fm)` | surface: empty/partial data written on failure |
| B-marshal | `internal/steering/format.go:195` | `raw, _ := yaml.Marshal(fm)` | surface: empty/partial data written on failure |
| B-marshal | `internal/steering/parse.go:79` | `raw, _ := yaml.Marshal(sf)` | surface: empty/partial data written on failure |
| B-parse | `internal/write/kb/sourcecoverage/parse.go:42` | `updated, _ := time.Parse(` | surface: zero value silently on failure |
| A-close-WRITE | `internal/cli/kb/cmd/note/run.go:52` | `defer func() { _ = f.Close() }()` | write handle: surface close error (data loss on flush fail) |
| A-close-WRITE | `internal/skill/copy.go:75` | `defer func() { _ = in.Close() }()` | write handle: surface close error (data loss on flush fail) |
| A-close-WRITE | `internal/skill/copy.go:81` | `defer func() { _ = out.Close() }()` | write handle: surface close error (data loss on flush fail) |
| A-close-WRITE | `internal/trace/jsonl.go:40` | `defer func() { _ = f.Close() }()` | write handle: surface close error (data loss on flush fail) |
| A-close-WRITE | `internal/trace/jsonl.go:89` | `defer func() { _ = f.Close() }()` | write handle: surface close error (data loss on flush fail) |
| A-close-WRITE | `internal/write/kb/evidence/append.go:70` | `defer func() { _ = f.Close() }()` | write handle: surface close error (data loss on flush fail) |
| A-close-WRITE | `internal/write/kb/glossary/append.go:49` | `defer func() { _ = f.Close() }()` | write handle: surface close error (data loss on flush fail) |
| A-close-WRITE | `internal/write/kb/relationship/append.go:50` | `defer func() { _ = f.Close() }()` | write handle: surface close error (data loss on flush fail) |
| A-close-WRITE | `internal/write/kb/row/append.go:53` | `defer func() { _ = f.Close() }()` | write handle: surface close error (data loss on flush fail) |
| A-close-WRITE | `internal/write/kb/sourcemap/append.go:49` | `defer func() { _ = f.Close() }()` | write handle: surface close error (data loss on flush fail) |
| A-close-WRITE | `internal/write/kb/timeline/append.go:49` | `defer func() { _ = f.Close() }()` | write handle: surface close error (data loss on flush fail) |
| SURFACE | `internal/cli/drift/cmd/root/run.go:75` | `ctx, _ = load.Do("")` | load.Do error dropped; stale/empty ctx re-detected — surface |
| SURFACE | `internal/cli/setup/core/opencode/opencode.go:44` | `mcpPath, _ := globalConfigPath()` | globalConfigPath resolve error dropped — surface |
| SURFACE | `internal/journal/parser/query.go:49` | `sessions, _ := ScanDirectory(resolved)` | ScanDirectory error dropped; a dir's sessions silently skipped — stderr |
| C-state | `internal/cli/connection/core/sync/state.go:52` | `release := func() { _ = os.Remove(lockPath) }` | stderr warn (stale state on failure) |
| C-state | `internal/cli/hub/core/server/daemon.go:120` | `_ = os.Remove(pidPath)` | stderr warn (stale state on failure) |
| C-state | `internal/cli/trace/core/hook/hook.go:128` | `_ = os.Remove(path)` | stderr warn (stale state on failure) |
| C-state | `internal/hub/server.go:60` | `_ = s.cluster.Shutdown()` | cluster.Shutdown error dropped on teardown — stderr |
| C-state | `internal/io/security.go:200` | `cleanup := func() { _ = os.Remove(tmpPath) }` | stderr warn (stale state on failure) |
| C-state | `internal/mcp/handler/violations.go:39` | `_ = os.Remove(filepath.Join(stateDir, file.Violations))` | stderr warn (stale state on failure) |
| C-state | `internal/skill/install.go:56` | `_ = os.RemoveAll(destDir)` | stderr warn (stale state on failure) |
| A-close-bare | `internal/hub/failover.go:57` | `_ = conn.Close()` | surface or stderr-log |
| A-close-bare | `internal/io/security.go:202` | `_ = tmp.Close()` | surface or stderr-log |
| A-close-bare | `internal/io/security.go:207` | `_ = tmp.Close()` | surface or stderr-log |
| A-close-read | `internal/cli/connection/core/listen/listen.go:44` | `defer func() { _ = client.Close() }()` | log close failure to stderr (EH.3) |
| A-close-read | `internal/cli/connection/core/publish/publish.go:44` | `defer func() { _ = client.Close() }()` | log close failure to stderr (EH.3) |
| A-close-read | `internal/cli/connection/core/register/register.go:43` | `defer func() { _ = client.Close() }()` | log close failure to stderr (EH.3) |
| A-close-read | `internal/cli/connection/core/status/status.go:39` | `defer func() { _ = client.Close() }()` | log close failure to stderr (EH.3) |
| A-close-read | `internal/cli/connection/core/sync/sync.go:48` | `defer func() { _ = client.Close() }()` | log close failure to stderr (EH.3) |
| A-close-read | `internal/cli/hub/core/status/status.go:40` | `defer func() { _ = client.Close() }()` | log close failure to stderr (EH.3) |
| A-close-read | `internal/cli/system/core/hubsync/sync.go:76` | `defer func() { _ = client.Close() }()` | log close failure to stderr (EH.3) |
| A-close-read | `internal/hub/replicate.go:84` | `defer func() { _ = conn.Close() }()` | log close failure to stderr (EH.3) |
| A-close-read | `internal/journal/parser/copilot.go:115` | `defer func() { _ = f.Close() }()` | log close failure to stderr (EH.3) |
| A-close-read | `internal/journal/parser/copilot.go:70` | `defer func() { _ = f.Close() }()` | log close failure to stderr (EH.3) |
| A-close-read | `internal/journal/parser/copilot_cli.go:77` | `defer func() { _ = f.Close() }()` | log close failure to stderr (EH.3) |
| A-close-read | `internal/trace/resolve_entry.go:85` | `defer func() { _ = f.Close() }()` | log close failure to stderr (EH.3) |
| A-close-read | `internal/trace/working_tasks.go:38` | `defer func() { _ = f.Close() }()` | log close failure to stderr (EH.3) |
| D-output | `internal/cli/journal/core/section/section.go:104` | `_, _ = fmt.Fprintf(sb,` | best-effort CLI output; convert to cmd.Print* to drop the discard |
| D-output | `internal/cli/journal/core/section/section.go:108` | `_, _ = fmt.Fprintf(sb, tpl.JournalIndexSummary+nl, e.Summary)` | best-effort CLI output; convert to cmd.Print* to drop the discard |
| D-output | `internal/cli/journal/core/section/section.go:177` | `_, _ = fmt.Fprintf(sb, ltTpl+nl, label, e.Title, link)` | best-effort CLI output; convert to cmd.Print* to drop the discard |
| D-output | `internal/cli/journal/core/section/section.go:97` | `_, _ = fmt.Fprintf(sb, tpl.JournalMonthHeading+nl+nl, month)` | best-effort CLI output; convert to cmd.Print* to drop the discard |
| D-output | `internal/log/warn/warn.go:35` | `_, _ = fmt.Fprintf(` | best-effort CLI output; convert to cmd.Print* to drop the discard |
| D-output | `internal/mcp/handler/tool.go:168` | `_, _ = fmt.Fprintf(` | best-effort CLI output; convert to cmd.Print* to drop the discard |
| D-output | `internal/mcp/handler/tool.go:179` | `_, _ = fmt.Fprintf(` | best-effort CLI output; convert to cmd.Print* to drop the discard |
| D-output | `internal/mcp/handler/tool.go:190` | `_, _ = fmt.Fprintf(` | best-effort CLI output; convert to cmd.Print* to drop the discard |
| D-output | `internal/mcp/handler/tool.go:246` | `_, _ = fmt.Fprintf(` | best-effort CLI output; convert to cmd.Print* to drop the discard |
| D-output | `internal/mcp/handler/tool.go:252` | `_, _ = fmt.Fprintf(` | best-effort CLI output; convert to cmd.Print* to drop the discard |
| D-output | `internal/mcp/handler/tool.go:257` | `_, _ = fmt.Fprintf(` | best-effort CLI output; convert to cmd.Print* to drop the discard |
| D-output | `internal/mcp/handler/tool.go:264` | `_, _ = fmt.Fprintf(` | best-effort CLI output; convert to cmd.Print* to drop the discard |
| D-output | `internal/mcp/handler/tool.go:394` | `_, _ = fmt.Fprintf(` | best-effort CLI output; convert to cmd.Print* to drop the discard |
| D-output | `internal/mcp/handler/tool.go:412` | `_, _ = fmt.Fprintf(` | best-effort CLI output; convert to cmd.Print* to drop the discard |
| D-output | `internal/mcp/handler/tool.go:554` | `_, _ = fmt.Fprintf(` | best-effort CLI output; convert to cmd.Print* to drop the discard |
| D-output | `internal/mcp/handler/tool.go:560` | `_, _ = fmt.Fprintf(` | best-effort CLI output; convert to cmd.Print* to drop the discard |
| D-output | `internal/mcp/handler/tool.go:71` | `_, _ = fmt.Fprintf(` | best-effort CLI output; convert to cmd.Print* to drop the discard |
| D-output | `internal/mcp/server/resource/resource.go:108` | `_, _ = fmt.Fprintf(` | best-effort CLI output; convert to cmd.Print* to drop the discard |
| D-output | `internal/mcp/server/route/prompt/entry.go:45` | `_, _ = fmt.Fprintf(` | best-effort CLI output; convert to cmd.Print* to drop the discard |
| D-output | `internal/mcp/server/route/prompt/prompt.go:111` | `_, _ = fmt.Fprintf(` | best-effort CLI output; convert to cmd.Print* to drop the discard |
| D-output | `internal/mcp/server/route/prompt/prompt.go:60` | `_, _ = fmt.Fprintf(` | best-effort CLI output; convert to cmd.Print* to drop the discard |
| D-output | `internal/memory/diff.go:46` | `_, _ = fmt.Fprintf(&buf, desc.Text(text.DescKeyMemoryDiffOldFormat), oldPath)` | best-effort CLI output; convert to cmd.Print* to drop the discard |
| D-output | `internal/memory/diff.go:47` | `_, _ = fmt.Fprintf(&buf, desc.Text(text.DescKeyMemoryDiffNewFormat), newPath)` | best-effort CLI output; convert to cmd.Print* to drop the discard |
| D-output | `internal/write/journal/source.go:183` | `_, _ = fmt.Fprintf(cmd.ErrOrStderr(),` | best-effort CLI output; convert to cmd.Print* to drop the discard |
| D-output | `internal/write/journal/source.go:187` | `_, _ = fmt.Fprintf(cmd.ErrOrStderr(),` | best-effort CLI output; convert to cmd.Print* to drop the discard |
| D-output | `internal/write/journal/source.go:191` | `_, _ = fmt.Fprintf(cmd.ErrOrStderr(),` | best-effort CLI output; convert to cmd.Print* to drop the discard |
| D-output | `internal/write/journal/source.go:262` | `_, _ = fmt.Fprintf(cmd.OutOrStdout(),` | best-effort CLI output; convert to cmd.Print* to drop the discard |
| D-output | `internal/write/journal/source.go:266` | `_, _ = fmt.Fprintf(cmd.OutOrStdout(),` | best-effort CLI output; convert to cmd.Print* to drop the discard |
| D-output | `internal/write/journal/source.go:282` | `_, _ = fmt.Fprintf(cmd.OutOrStdout(), format, values...)` | best-effort CLI output; convert to cmd.Print* to drop the discard |
| D-output | `internal/write/journal/source.go:294` | `_, _ = fmt.Fprintln(cmd.OutOrStdout())` | best-effort CLI output; convert to cmd.Print* to drop the discard |
| D-output | `internal/write/journal/source.go:351` | `_, _ = fmt.Fprintf(cmd.OutOrStdout(),` | best-effort CLI output; convert to cmd.Print* to drop the discard |
| D-output | `internal/write/journal/source.go:367` | `_, _ = fmt.Fprintf(cmd.OutOrStdout(),` | best-effort CLI output; convert to cmd.Print* to drop the discard |
| D-output | `internal/write/journal/source.go:384` | `_, _ = fmt.Fprintf(cmd.OutOrStdout(),` | best-effort CLI output; convert to cmd.Print* to drop the discard |
| D-output | `internal/write/journal/source.go:388` | `_, _ = fmt.Fprintln(cmd.OutOrStdout())` | best-effort CLI output; convert to cmd.Print* to drop the discard |
| D-output | `internal/write/journal/source.go:399` | `_, _ = fmt.Fprintln(cmd.OutOrStdout())` | best-effort CLI output; convert to cmd.Print* to drop the discard |
| D-output | `internal/write/journal/source.go:413` | `_, _ = fmt.Fprintf(cmd.OutOrStdout(),` | best-effort CLI output; convert to cmd.Print* to drop the discard |
| D-output | `internal/write/journal/source.go:416` | `_, _ = fmt.Fprintln(cmd.OutOrStdout())` | best-effort CLI output; convert to cmd.Print* to drop the discard |
| D-output | `internal/write/journal/source.go:428` | `_, _ = fmt.Fprintln(cmd.OutOrStdout(), body)` | best-effort CLI output; convert to cmd.Print* to drop the discard |
| D-output | `internal/write/journal/source.go:429` | `_, _ = fmt.Fprintln(cmd.OutOrStdout())` | best-effort CLI output; convert to cmd.Print* to drop the discard |
| D-output | `internal/write/journal/source.go:441` | `_, _ = fmt.Fprintf(cmd.OutOrStdout(),` | best-effort CLI output; convert to cmd.Print* to drop the discard |
| D-output | `internal/write/journal/source.go:456` | `_, _ = fmt.Fprintf(cmd.OutOrStdout(),` | best-effort CLI output; convert to cmd.Print* to drop the discard |
| D-output | `internal/write/journal/source.go:471` | `_, _ = fmt.Fprintf(cmd.OutOrStdout(),` | best-effort CLI output; convert to cmd.Print* to drop the discard |
| D-output | `internal/write/journal/source.go:485` | `_, _ = fmt.Fprintf(cmd.OutOrStdout(),` | best-effort CLI output; convert to cmd.Print* to drop the discard |
| D-output | `internal/write/load/load.go:70` | `_, _ = fmt.Fprintf(&sb, tpl.LoadBudget+nl+nl, budget, totalTokens)` | best-effort CLI output; convert to cmd.Print* to drop the discard |
| D-output | `internal/write/load/load.go:82` | `_, _ = fmt.Fprintf(&sb, nl+sep+nl+nl+tpl.LoadTruncated+nl, f.Name)` | best-effort CLI output; convert to cmd.Print* to drop the discard |
| D-output | `internal/write/load/load.go:86` | `_, _ = fmt.Fprintf(&sb, tpl.LoadSectionHeading+nl+nl, titleFn(f.Name))` | best-effort CLI output; convert to cmd.Print* to drop the discard |
| D-output | `internal/write/stat/stream.go:20` | `_, _ = fmt.Fprintln(w, line)` | best-effort CLI output; convert to cmd.Print* to drop the discard |
| besteffort | `internal/cli/add/core/run/run.go:115` | `_ = trace.Record(fType+cfgTrace.RefFirstEntry, stateDir)` | trace.Record telemetry; best-effort — annotate |
| besteffort | `internal/cli/change/cmd/root/run.go:42` | `ctxChanges, _ := scan.FindContextChanges(refTime)` | FindContextChanges for display; best-effort — annotate |
| besteffort | `internal/cli/change/cmd/root/run.go:43` | `codeChanges, _ := scan.SummarizeCodeChanges(refTime)` | SummarizeCodeChanges for display; best-effort — annotate |
| besteffort | `internal/cli/initialize/core/claudecheck/detail.go:51` | `marketplaces, _ := readMarketplaces(marketplacesPath)` | readMarketplaces detection; best-effort — annotate |
| besteffort | `internal/cli/system/cmd/checkcontextsize/run.go:93` | `info, _ := coreSession.ReadTokenInfo(sessionID)` | ReadTokenInfo for display hint; best-effort — annotate |
| besteffort | `internal/cli/system/cmd/contextloadgate/run.go:135` | `ctxChanges, _ := changeCore.FindContextChanges(refTime)` | FindContextChanges in gate; best-effort — annotate |
| besteffort | `internal/cli/system/cmd/contextloadgate/run.go:136` | `codeChanges, _ := changeCore.SummarizeCodeChanges(refTime)` | SummarizeCodeChanges in gate; best-effort — annotate |
| besteffort | `internal/cli/system/cmd/heartbeat/run.go:95` | `info, _ := session.ReadTokenInfo(sessionID)` | ReadTokenInfo for display hint; best-effort — annotate |
| besteffort | `internal/cli/system/core/message/message.go:125` | `line, _ := ctxContext.DirLine()` | DirLine for display; best-effort — annotate |
| besteffort | `internal/cli/system/core/provenance/provenance.go:66` | `info, _ := coreSession.ReadTokenInfo(sessionID)` | ReadTokenInfo for display hint; best-effort — annotate |
| besteffort | `internal/cli/task/cmd/complete/run.go:45` | `_ = trace.Record(ref, stateDir)` | trace.Record telemetry; best-effort — annotate |
| besteffort | `internal/cli/trace/core/collect/collect.go:49` | `_ = trace.TruncatePending(stateDir)` | TruncatePending cleanup; best-effort — stderr at most |
| besteffort | `internal/cli/trace/core/collect/collect.go:69` | `_ = trace.TruncatePending(stateDir)` | TruncatePending cleanup; best-effort — stderr at most |
| besteffort | `internal/cli/trace/core/show/show.go:48` | `msg, _ := trace.CommitMessage(fullHash)` | CommitMessage for display; best-effort — annotate |
| besteffort | `internal/cli/trace/core/show/show.go:59` | `msg, _ := trace.CommitMessage(fullHash)` | CommitMessage for display; best-effort — annotate |
| besteffort | `internal/mcp/server/server.go:50` | `_ = srv.out.WriteJSON(n)` | WriteJSON in poller callback; no return path — stderr |
| OK | `internal/cli/add/core/extract/content.go:64` | `stat, _ := os.Stdin.Stat()` | os.Stdin.Stat for pipe detection; nil-safe fallback — annotate |
| OK | `internal/cli/setup/core/copilotcli/mcp.go:59` | `servers, _ := existing[cfgHook.KeyMCPServers].(map[string]interface{})` | type-assert ok-discard; zero map handled |
| OK | `internal/cli/setup/core/opencode/mcp.go:108` | `servers, _ := existing[cfgHook.KeyMCP].(map[string]interface{})` | type-assert ok-discard; zero map handled |
| OK | `internal/cli/system/core/session/session_token.go:81` | `if initialized, _ := state.Initialized(); !initialized {` | state.Initialized bool used; ok-discard |
| OK-flag | `internal/cli/doctor/cmd/root/cmd.go:36` | `jsonOut, _ := cmd.Flags().GetBool(cFlag.JSON)` | error only if flag unregistered (programmer err); annotate |
| OK-flag | `internal/cli/event/run.go:29` | `hook, _ := cmd.Flags().GetString(cFlag.Hook)` | error only if flag unregistered (programmer err); annotate |
| OK-flag | `internal/cli/event/run.go:30` | `session, _ := cmd.Flags().GetString(cFlag.Session)` | error only if flag unregistered (programmer err); annotate |
| OK-flag | `internal/cli/event/run.go:31` | `event, _ := cmd.Flags().GetString(cFlag.Event)` | error only if flag unregistered (programmer err); annotate |
| OK-flag | `internal/cli/event/run.go:32` | `last, _ := cmd.Flags().GetInt(cFlag.Last)` | error only if flag unregistered (programmer err); annotate |
| OK-flag | `internal/cli/event/run.go:33` | `jsonOut, _ := cmd.Flags().GetBool(cFlag.JSON)` | error only if flag unregistered (programmer err); annotate |
| OK-flag | `internal/cli/event/run.go:34` | `includeAll, _ := cmd.Flags().GetBool(cFlag.All)` | error only if flag unregistered (programmer err); annotate |
| OK-flag | `internal/cli/message/cmd/list/run.go:57` | `jsonFlag, _ := cmd.Flags().GetBool(cFlag.JSON)` | error only if flag unregistered (programmer err); annotate |
| OK-flag | `internal/cli/pause/cmd/root/cmd.go:30` | `sessionID, _ := cmd.Flags().GetString(cFlag.SessionID)` | error only if flag unregistered (programmer err); annotate |
| OK-flag | `internal/cli/resolve/tool.go:34` | `v, _ := cmd.Flags().GetString(flag.Tool)` | error only if flag unregistered (programmer err); annotate |
| OK-flag | `internal/cli/resume/cmd/root/cmd.go:31` | `sessionID, _ := cmd.Flags().GetString(cFlag.SessionID)` | error only if flag unregistered (programmer err); annotate |
| OK-flag | `internal/cli/sysinfo/run.go:31` | `jsonFlag, _ := cmd.Flags().GetBool(cFlag.JSON)` | error only if flag unregistered (programmer err); annotate |
| OK-flag | `internal/cli/system/cmd/bootstrap/run.go:51` | `quiet, _ := cmd.Flags().GetBool(cFlag.Quiet)` | error only if flag unregistered (programmer err); annotate |
| OK-flag | `internal/cli/system/cmd/bootstrap/run.go:66` | `jsonFlag, _ := cmd.Flags().GetBool(cFlag.JSON)` | error only if flag unregistered (programmer err); annotate |
| OK-flag | `internal/cli/system/cmd/markjournal/run.go:42` | `check, _ := cmd.Flags().GetBool(cFlag.Check)` | error only if flag unregistered (programmer err); annotate |
| OK-flag | `internal/cli/system/core/check/pause_preamble.go:53` | `sessionID, _ := cmd.Flags().GetString(cFlag.SessionID)` | error only if flag unregistered (programmer err); annotate |
| OK-flag | `internal/cli/usage/run.go:31` | `follow, _ := cmd.Flags().GetBool(cFlag.Follow)` | error only if flag unregistered (programmer err); annotate |
| OK-flag | `internal/cli/usage/run.go:32` | `session, _ := cmd.Flags().GetString(cFlag.Session)` | error only if flag unregistered (programmer err); annotate |
| OK-flag | `internal/cli/usage/run.go:33` | `last, _ := cmd.Flags().GetInt(cFlag.Last)` | error only if flag unregistered (programmer err); annotate |
| OK-flag | `internal/cli/usage/run.go:34` | `jsonOut, _ := cmd.Flags().GetBool(cFlag.JSON)` | error only if flag unregistered (programmer err); annotate |
| OK-markflag | `internal/cli/add/core/build/build.go:119` | `_ = c.RegisterFlagCompletionFunc(` | RegisterFlagCompletionFunc; init-time, bad-name only — annotate |
| OK-markflag | `internal/cli/connection/cmd/register/cmd.go:48` | `_ = c.MarkFlagRequired(cFlag.Token)` | init-time; bad flag name only; annotate or handle |
| OK-markflag | `internal/cli/handover/cmd/write/cmd.go:67` | `_ = c.MarkFlagRequired(cFlag.Summary)` | init-time; bad flag name only; annotate or handle |
| OK-markflag | `internal/cli/handover/cmd/write/cmd.go:68` | `_ = c.MarkFlagRequired(cFlag.Next)` | init-time; bad flag name only; annotate or handle |
| OK-markflag | `internal/cli/system/cmd/sessionevent/cmd.go:46` | `_ = c.MarkFlagRequired(cFlag.Type)` | init-time; bad flag name only; annotate or handle |
| OK-markflag | `internal/cli/system/cmd/sessionevent/cmd.go:47` | `_ = c.MarkFlagRequired(cFlag.Caller)` | init-time; bad flag name only; annotate or handle |
| OK-markflag | `internal/cli/trace/cmd/tag/cmd.go:37` | `_ = c.MarkFlagRequired(cFlag.Note)` | init-time; bad flag name only; annotate or handle |
| OK-glob | `internal/cli/sync/core/validate/validate.go:93` | `matches, _ := filepath.Glob(cfg.Pattern)` | filepath.Glob; static pattern, nil-safe range — annotate |
| OK-glob | `internal/cli/system/core/stats/stats.go:288` | `matches, _ := filepath.Glob(globPat)` | filepath.Glob; static pattern, nil-safe range — annotate |
| OK-glob | `internal/cli/system/core/stats/stats.go:300` | `matches, _ = filepath.Glob(globPat)` | filepath.Glob; static pattern, nil-safe range — annotate |
| OK-atoi | `internal/cli/journal/core/normalize/boundary.go:33` | `num, _ := strconv.Atoi(m[1])` | regex-guaranteed digits; annotate why safe |
| OK-atoi | `internal/cli/journal/core/normalize/normalize.go:407` | `num, _ := strconv.Atoi(m[1])` | regex-guaranteed digits; annotate why safe |
| OK-atoi | `internal/cli/pad/core/parse/entry.go:52` | `id, _ := strconv.Atoi(match[1])` | regex-guaranteed digits; annotate why safe |
| OK-atoi | `internal/cli/system/core/nudge/pause.go:63` | `count, _ := strconv.Atoi(strings.TrimSpace(string(data)))` | regex-guaranteed digits; annotate why safe |
| OK-typeassert | `internal/mcp/server/extract/extract.go:33` | `entryType, _ := args[cli.AttrType].(string)` | discards ok bool; zero value handled |
| OK-typeassert | `internal/mcp/server/extract/extract.go:34` | `content, _ := args[field.Content].(string)` | discards ok bool; zero value handled |
| OK-typeassert | `internal/mcp/server/route/tool/steering.go:35` | `prompt, _ := args[field.Prompt].(string)` | discards ok bool; zero value handled |
| OK-typeassert | `internal/mcp/server/route/tool/steering.go:54` | `query, _ := args[field.Query].(string)` | discards ok bool; zero value handled |
| OK-typeassert | `internal/mcp/server/route/tool/steering.go:78` | `summary, _ := args[field.Summary].(string)` | discards ok bool; zero value handled |
| OK-typeassert | `internal/mcp/server/route/tool/tool.go:113` | `if sinceStr, _ := args[field.Since].(string); sinceStr != "" {` | discards ok bool; zero value handled |
| OK-typeassert | `internal/mcp/server/route/tool/tool.go:209` | `recentAction, _ := args[field.RecentAction].(string)` | discards ok bool; zero value handled |
| OK-typeassert | `internal/mcp/server/route/tool/tool.go:228` | `eventType, _ := args[cli.AttrType].(string)` | discards ok bool; zero value handled |
| OK-typeassert | `internal/mcp/server/route/tool/tool.go:234` | `caller, _ := args[field.Caller].(string)` | discards ok bool; zero value handled |
| OK-typeassert | `internal/mcp/server/route/tool/tool.go:78` | `query, _ := args[field.Query].(string)` | discards ok bool; zero value handled |
| FALSE-POS | `internal/cli/config/cmd/status/cmd.go:23` | `short, _ := desc.Command(cmd.DescKeyConfigStatus)` | not an error discard |
| FALSE-POS | `internal/cli/initialize/core/project/getting_started.go:34` | `_ = contextDir` | not an error discard |
| FALSE-POS | `internal/cli/message/cmd/edit/cmd.go:21` | `short, _ := desc.Command(cmd.DescKeyMessageEdit)` | not an error discard |
| FALSE-POS | `internal/cli/message/cmd/list/cmd.go:23` | `short, _ := desc.Command(cmd.DescKeyMessageList)` | not an error discard |
| FALSE-POS | `internal/cli/message/cmd/reset/cmd.go:21` | `short, _ := desc.Command(cmd.DescKeyMessageReset)` | not an error discard |
| FALSE-POS | `internal/cli/message/cmd/show/cmd.go:21` | `short, _ := desc.Command(cmd.DescKeyMessageShow)` | not an error discard |
| FALSE-POS | `internal/cli/pad/cmd/add/cmd.go:26` | `short, _ := desc.Command(cmd.DescKeyPadAdd)` | not an error discard |
| FALSE-POS | `internal/cli/pad/cmd/mv/cmd.go:23` | `short, _ := desc.Command(cmd.DescKeyPadMv)` | not an error discard |
| FALSE-POS | `internal/cli/pad/cmd/normalize/cmd.go:21` | `short, _ := desc.Command(cmd.DescKeyPadNormalize)` | not an error discard |
| FALSE-POS | `internal/cli/pad/cmd/rm/cmd.go:25` | `short, _ := desc.Command(cmd.DescKeyPadRm)` | not an error discard |
| FALSE-POS | `internal/cli/remind/cmd/add/cmd.go:26` | `short, _ := desc.Command(cmd.DescKeyRemindAdd)` | not an error discard |
| FALSE-POS | `internal/cli/remind/cmd/dismiss/cmd.go:30` | `short, _ := desc.Command(cmd.DescKeyRemindDismiss)` | not an error discard |
| FALSE-POS | `internal/cli/remind/cmd/list/cmd.go:21` | `short, _ := desc.Command(cmd.DescKeyRemindList)` | not an error discard |
| FALSE-POS | `internal/cli/skill/cmd/list/cmd.go:27` | `short, _ := desc.Command(cmd.DescKeySkillList)` | not an error discard |
| FALSE-POS | `internal/cli/steering/cmd/list/cmd.go:28` | `short, _ := desc.Command(cmd.DescKeySteeringList)` | not an error discard |
| FALSE-POS | `internal/cli/system/cmd/checkmemorydrift/cmd.go:23` | `short, _ := desc.Command(cmd.DescKeySystemCheckMemoryDrift)` | not an error discard |
| FALSE-POS | `internal/cli/system/cmd/checkreminder/cmd.go:23` | `short, _ := desc.Command(cmd.DescKeySystemCheckReminder)` | not an error discard |
| FALSE-POS | `internal/cli/trigger/cmd/list/cmd.go:25` | `short, _ := desc.Command(cmd.DescKeyTriggerList)` | not an error discard |
| FALSE-POS | `internal/format/format.go:125` | `line, _, _ := strings.Cut(s, token.NewlineLF)` | strings.Cut discards after+found bool, not an error |
| FALSE-POS | `internal/hub/replicate.go:22` | `var _ = startReplication` | not an error discard |
| FALSE-POS | `internal/hub/store_sequence.go:11` | `var _ = (*Store).lastSequence` | var _ = method value: compile-time assertion, not an error |
