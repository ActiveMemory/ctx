//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package warn

// Format strings for file I/O warnings. Each takes (path, error).
const (
	// Close is the format for file close failures.
	Close = "close %s: %v"

	// Write is the format for file write failures.
	Write = "write %s: %v"

	// Remove is the format for file remove failures.
	Remove = "remove %s: %v"

	// Mkdir is the format for directory creation failures.
	Mkdir = "mkdir %s: %v"

	// Rename is the format for file rename failures.
	Rename = "rename %s: %v"

	// Walk is the format for directory walk failures.
	Walk = "walk %s: %v"

	// Getwd is the format for working directory resolution failures.
	Getwd = "getwd: %v"

	// Marshal is the format for JSON marshal failures. Takes (error).
	Marshal = "marshal: %v"

	// Readdir is the format for directory read failures.
	Readdir = "readdir %s: %v"

	// CloseResponse is the format for HTTP response body close failures.
	CloseResponse = "close response body: %v"

	// ParseConfig is the format for config file parse failures.
	ParseConfig = "warning: failed to parse %s: %v (using defaults)"

	// CopilotClose is the format for Copilot CLI file close failures.
	CopilotClose = "copilot-cli: close %s: %v"

	// RelayUnknownSubcommand is the format for a best-effort relay
	// failure when `ctx system` reports an unknown subcommand. The
	// stdout box already reached the agent; this only logs that the
	// event-log/webhook leg could not be recorded.
	RelayUnknownSubcommand = "relay unknown-subcommand event: %v"

	// JSONEncode is the JSON-safe error for encoding failures.
	JSONEncode = `{"error": "json encode: %v"}`

	// ContextDirResolve is the stderr format for unexpected
	// rc.ContextDir failures in hook paths that must not propagate.
	// The declared-vs-undeclared split is matched with errors.Is at
	// each call site; this constant is used only when that match
	// fails, which should never happen with the current single-error
	// return but catches future regressions loudly.
	ContextDirResolve = "resolve context dir: %v"

	// RCNoContextDir is the stderr message emitted by rc.load when
	// it observes ErrNoCtxHere. Exempt commands (init,
	// activate, doctor, hub *, etc.) legitimately reach this state;
	// they call accessors and want defaults. Operating commands
	// should never reach it because [bootstrap/cmd.go]'s
	// PersistentPreRunE gate calls RequireContextDir first. The
	// warning is the breadcrumb that catches a missed-gate
	// regression: an operating command added without the gate
	// would silently get default config (token_budget = 8000,
	// auto_archive = true, etc.) regardless of what the user's
	// .ctxrc says, with no diagnostic. This message makes the
	// silence visible so the call site can be evaluated.
	RCNoContextDir = "rc.RC: no .context/ at $PWD; " +
		"defaults applied " +
		"(investigate calling command if unexpected)"

	// ReadMapTracking is the stderr format for map-tracking.json
	// read / parse failures in the check-map-staleness hook. The
	// hook can't fail the user's tool call, so it logs and returns
	// nil; the log line keeps the failure visible instead of having
	// the staleness check silently stop firing.
	ReadMapTracking = "read map tracking: %v"

	// CheckKnowledge is the stderr format for check-knowledge hook
	// failures downstream of rc.ContextDir resolution. Same shape
	// as ReadMapTracking: hook surfaces the error rather than
	// silently going dark.
	CheckKnowledge = "check knowledge: %v"

	// HubConnectedProbe is the stderr format for failures inside
	// [hubsync.Connected] beyond "no context dir declared" and
	// "connect file missing." Surfacing the error keeps operators
	// from wondering why the hub silently stopped syncing after a
	// broken .ctxrc or permissions regression.
	HubConnectedProbe = "probe hub connection: %v"

	// CloseHubClient is the stderr format for a failed hub gRPC
	// client/connection close. The close runs in a defer after the
	// command's real work, so the error is not actionable but should
	// not vanish.
	CloseHubClient = "close hub client: %v"

	// HubReplicateAppend is the stderr format for a failed
	// [Store.Append] inside the follower replication stream. The
	// loop is best-effort and has no return path, so a dropped
	// append would silently lose a replicated entry; warning keeps
	// the loss visible.
	HubReplicateAppend = "hub replicate append: %v"

	// StateInitializedProbe is the stderr format for failures
	// inside [state.Initialized] beyond "no context dir declared."
	// Hooks bail on false either way, but a visible warning shows
	// operators why the hook stopped firing instead of letting the
	// failure vanish into the gap between "initialized" and "not."
	StateInitializedProbe = "probe state initialized: %v"

	// StateDirProbe is the stderr format for failures inside
	// [state.Dir] beyond "no context dir declared." Callers use
	// the returned path as a filepath.Join base; a warning here
	// explains why the state directory resolution went sideways
	// before the caller surfaces an empty-path error.
	StateDirProbe = "probe state dir: %v"

	// SteeringUnfilled is the stderr format for steering files
	// that still carry the cfgSteering.Tombstone placeholder
	// marker. The file is skipped on every load path (agent
	// context packet, MCP ctx_steering_get, sync to Cursor /
	// Cline / Kiro). The warning is the breadcrumb that tells
	// the user a scaffolded steering file is silently inert
	// until the tombstone line is removed.
	SteeringUnfilled = "skipping unfilled steering file %s " +
		"(remove the tombstone line to activate)"

	// TemplateRender is the stderr format for an embedded-template
	// render failure. Parse is gated by TestTemplatesParse and the
	// data is typed, so Execute cannot fail in a correct build; the
	// warning catches a future regression loudly instead of letting
	// [tpl.RenderOr]'s fallback silently blank a section.
	TemplateRender = "render template: %v"
)

// Pad history warning formats.
const (
	// PadHistoryPruneFile is the format for per-file prune
	// failures in the scratchpad history directory.
	// Takes (snapshot filename, error).
	PadHistoryPruneFile = "pad history: prune %s: %v"

	// PadHistoryPrune is the format for an overall prune
	// failure when iterating the scratchpad history.
	// Takes (error).
	PadHistoryPrune = "pad history: prune: %v"
)

// Warn context identifiers for index generation.
const (
	// IndexHeader is the context label for index header write errors.
	IndexHeader = "index-header"
	// IndexSeparator is the context label for index separator write
	// errors.
	IndexSeparator = "index-separator"
	// IndexRow is the context label for index row write errors.
	IndexRow = "index-row"
	// ResponseBody is the context label for HTTP response body
	// close errors.
	ResponseBody = "response body"
)
