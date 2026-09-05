//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package assets

import (
	"bytes"
	"encoding/json"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"slices"
	"strings"
	"testing"

	"github.com/ActiveMemory/ctx/internal/config/asset"
	cfgCodex "github.com/ActiveMemory/ctx/internal/config/codex"
	cfgHook "github.com/ActiveMemory/ctx/internal/config/hook"
	cfgMcpServer "github.com/ActiveMemory/ctx/internal/config/mcp/server"
)

// Guards for the embedded Codex plugin root
// (internal/assets/codex/) and the repo marketplace that points at
// it (.agents/plugins/marketplace.json). See
// specs/codex-integration.md "Validation Rules".

// claudeHookAnchor is the host-detecting prologue every command in the
// Claude Code manifest starts with. The Codex manifest uses
// [cfgCodex.HookAnchor] instead (Codex runs hooks with the session
// cwd, not a project-dir env var); parity is asserted on the
// `ctx …` tails that follow the anchors.
const claudeHookAnchor = `command -v ctx >/dev/null 2>&1 || exit 0; [ -n "${CLAUDE_PROJECT_DIR:-}" ] || exit 0; [ -d "$CLAUDE_PROJECT_DIR" ] || { echo "ctx: CLAUDE_PROJECT_DIR \"$CLAUDE_PROJECT_DIR\" is missing; restart the session at the project root" >&2; exit 1; }; cd "$CLAUDE_PROJECT_DIR" && `

// agentTailPrefix identifies the context-packet hook. Claude Code
// wires it under PreToolUse (plain stdout becomes context there);
// Codex ignores plain text on PreToolUse, so the same command is
// wired under SessionStart.
const agentTailPrefix = "ctx agent "

// claudeMatcherAlias maps a Claude Code matcher to the Codex
// matcher that stands in for it. Any other matcher must be carried
// over verbatim, or appear as one `|`-alternative of the Codex
// matcher (Codex's `apply_patch|Edit|Write` covers Claude's
// separate `Edit` and `Write` groups).
var claudeMatcherAlias = map[string]string{
	"EnterPlanMode": cfgCodex.ToolUpdatePlan,
}

// codexOnlySkillExclusions are the Claude Code skills that
// hack/sync-codex-skills.sh deliberately does not mirror into the
// Codex plugin (they operate on Claude Code-specific state). Keep
// this list identical to EXCLUDE in that script.
var codexOnlySkillExclusions = []string{
	"ctx-permission-sanitize",
	"ctx-plan-import",
	"ctx-dream",
	"ctx-skill-create",
}

// claudeHooksPath is the on-disk path of the Claude Code manifest
// relative to this package (it is not embedded; the plugin ships
// it from the repo tree).
var claudeHooksPath = filepath.Join(
	asset.DirClaude, cfgCodex.DirHooks, asset.FileHooksJSON,
)

// marketplacePath is the repo-root Codex marketplace catalog.
var marketplacePath = filepath.Join(
	repoRoot, cfgCodex.DirAgents, cfgCodex.DirMarketplacePlugins,
	cfgCodex.FileMarketplaceJSON,
)

// hookHandler is one command entry inside a matcher group.
type hookHandler struct {
	Type    string `json:"type"`
	Command string `json:"command"`
	Timeout int    `json:"timeout"`
}

// hookGroup is one matcher group under an event key.
type hookGroup struct {
	Matcher  string        `json:"matcher"`
	Handlers []hookHandler `json:"hooks"`
}

// hookManifest is the shared shape of the Claude Code and Codex
// hooks.json files.
type hookManifest struct {
	Hooks map[string][]hookGroup `json:"hooks"`
}

// hookEntry is a flattened (matcher, ctx tail) pair.
type hookEntry struct {
	Matcher string
	Tail    string
}

// decodeManifest parses a hooks.json body.
func decodeManifest(t *testing.T, label string, data []byte) hookManifest {
	t.Helper()
	var m hookManifest
	if err := json.Unmarshal(data, &m); err != nil {
		t.Fatalf("%s: parse: %v", label, err)
	}
	if m.Hooks == nil {
		t.Fatalf("%s: missing top-level %q key", label, cfgCodex.KeyHooks)
	}
	return m
}

// readCodexHooks parses the embedded Codex manifest.
func readCodexHooks(t *testing.T) hookManifest {
	t.Helper()
	data, err := FS.ReadFile(asset.PathCodexHooksJSON)
	if err != nil {
		t.Fatalf("read %s: %v", asset.PathCodexHooksJSON, err)
	}
	return decodeManifest(t, asset.PathCodexHooksJSON, data)
}

// readClaudeHooks parses the on-disk Claude Code manifest.
func readClaudeHooks(t *testing.T) hookManifest {
	t.Helper()
	data, err := os.ReadFile(claudeHooksPath)
	if err != nil {
		t.Fatalf("read %s: %v", claudeHooksPath, err)
	}
	return decodeManifest(t, claudeHooksPath, data)
}

// tails flattens a manifest into event → (matcher, tail) entries,
// stripping the given anchor from every command. Commands that do
// not start with the anchor are reported and skipped.
func tails(
	t *testing.T, label string, m hookManifest, anchor string,
) map[string][]hookEntry {
	t.Helper()
	out := make(map[string][]hookEntry, len(m.Hooks))
	for event, groups := range m.Hooks {
		for _, g := range groups {
			for _, h := range g.Handlers {
				if !strings.HasPrefix(h.Command, anchor) {
					t.Errorf(
						"%s: %s hook %q does not start with anchor %q",
						label, event, h.Command, anchor,
					)
					continue
				}
				out[event] = append(out[event], hookEntry{
					Matcher: g.Matcher,
					Tail:    strings.TrimPrefix(h.Command, anchor),
				})
			}
		}
	}
	return out
}

// codexEventFor returns the Codex event a Claude Code hook must be
// wired under.
func codexEventFor(claudeEvent, tail string) string {
	if claudeEvent == cfgHook.EventPreToolUse &&
		strings.HasPrefix(tail, agentTailPrefix) {
		return cfgCodex.EventSessionStart
	}
	return claudeEvent
}

// matcherCompatible reports whether a Codex hook under codexEvent
// with codexMatcher stands in for a Claude hook under claudeEvent
// with claudeMatcher: identical matchers, the aliased planning
// tool, one of the `|`-alternatives of a combined Codex matcher,
// or any matcher when the hook was relocated to a different event
// (SessionStart has no tool matcher to compare against).
func matcherCompatible(
	claudeEvent, claudeMatcher, codexEvent, codexMatcher string,
) bool {
	if claudeEvent != codexEvent || claudeMatcher == codexMatcher {
		return true
	}
	if alias, ok := claudeMatcherAlias[claudeMatcher]; ok {
		return alias == codexMatcher
	}
	return slices.Contains(strings.Split(codexMatcher, "|"), claudeMatcher)
}

// hasCounterpart reports whether entries contains an entry with the
// given tail whose matcher satisfies match.
func hasCounterpart(
	entries []hookEntry, tail string, match func(string) bool,
) bool {
	for _, e := range entries {
		if e.Tail == tail && match(e.Matcher) {
			return true
		}
	}
	return false
}

// TestCodexHooksManifestShape asserts the structural contract of
// the embedded Codex manifest: known events only, non-empty
// matcher groups, command-type handlers anchored to the git root,
// and a SessionEnd timeout within Codex's cap.
func TestCodexHooksManifestShape(t *testing.T) {
	m := readCodexHooks(t)
	if len(m.Hooks) == 0 {
		t.Fatal("codex hooks.json wires no events")
	}
	for event, groups := range m.Hooks {
		if !slices.Contains(cfgCodex.Events, event) {
			t.Errorf("event %q is not a Codex lifecycle event", event)
		}
		if len(groups) == 0 {
			t.Errorf("event %q has no matcher groups", event)
		}
		for i, g := range groups {
			if len(g.Handlers) == 0 {
				t.Errorf("%s group %d has no handlers", event, i)
			}
			for _, h := range g.Handlers {
				if h.Type != cfgCodex.HandlerTypeCommand {
					t.Errorf(
						"%s hook %q: type %q, want %q",
						event, h.Command, h.Type,
						cfgCodex.HandlerTypeCommand,
					)
				}
				if !strings.HasPrefix(h.Command, cfgCodex.HookPrologue) {
					t.Errorf(
						"%s hook %q does not start with the git-root anchor %q",
						event, h.Command, cfgCodex.HookPrologue,
					)
				}
				if event == cfgCodex.EventSessionEnd &&
					h.Timeout > cfgCodex.SessionEndTimeoutMax {
					t.Errorf(
						"SessionEnd hook %q: timeout %d exceeds Codex cap %d",
						h.Command, h.Timeout, cfgCodex.SessionEndTimeoutMax,
					)
				}
			}
		}
	}
}

// TestCodexHooksParityWithClaude asserts the Codex manifest wires
// the same `ctx …` commands as the Claude Code manifest, under the
// documented event/matcher mapping (specs/codex-integration.md
// "Event mapping"). A hook added to the Claude manifest fails here
// until it is mirrored; a Codex-only hook fails until it is
// justified by a Claude counterpart.
func TestCodexHooksParityWithClaude(t *testing.T) {
	claude := tails(t, claudeHooksPath, readClaudeHooks(t), claudeHookAnchor)
	codex := tails(t, asset.PathCodexHooksJSON, readCodexHooks(t), cfgCodex.HookPrologue)

	// Claude → Codex: every Claude hook has a Codex counterpart
	// under the mapped event with a compatible matcher.
	for claudeEvent, entries := range claude {
		for _, e := range entries {
			want := codexEventFor(claudeEvent, e.Tail)
			found := hasCounterpart(
				codex[want], e.Tail,
				func(codexMatcher string) bool {
					return matcherCompatible(
						claudeEvent, e.Matcher, want, codexMatcher,
					)
				},
			)
			if !found {
				t.Errorf(
					"Claude %s [%q] %q has no Codex counterpart under %s",
					claudeEvent, e.Matcher, e.Tail, want,
				)
			}
		}
	}

	// Codex → Claude: every Codex hook traces back to a Claude hook.
	for codexEvent, entries := range codex {
		if codexEvent == cfgCodex.EventStop {
			// Codex-only: the async Stop-hook journal import has
			// no Claude counterpart (Claude imports on SessionEnd
			// without Codex's 3-second cap).
			continue
		}
		for _, e := range entries {
			found := false
			for claudeEvent, claudeEntries := range claude {
				found = hasCounterpart(
					claudeEntries, e.Tail,
					func(claudeMatcher string) bool {
						return codexEventFor(claudeEvent, e.Tail) == codexEvent &&
							matcherCompatible(
								claudeEvent, claudeMatcher, codexEvent, e.Matcher,
							)
					},
				)
				if found {
					break
				}
			}
			if !found {
				t.Errorf(
					"Codex %s [%q] %q has no Claude counterpart",
					codexEvent, e.Matcher, e.Tail,
				)
			}
		}
	}

	// The context packet moves from Claude PreToolUse to Codex
	// SessionStart and must not also fire on PreToolUse.
	for _, e := range codex[cfgCodex.EventPreToolUse] {
		if strings.HasPrefix(e.Tail, agentTailPrefix) {
			t.Errorf(
				"Codex PreToolUse wires %q; the context packet belongs under SessionStart",
				e.Tail,
			)
		}
	}
	if len(codex[cfgCodex.EventSessionStart]) == 0 {
		t.Error("Codex SessionStart wires nothing; expected the ctx agent packet")
	}

	// The planning nudge must target Codex's planning tool, and
	// file-edit hooks must cover Codex's apply_patch.
	for _, e := range codex[cfgCodex.EventPreToolUse] {
		if e.Matcher == "EnterPlanMode" {
			t.Errorf(
				"Codex PreToolUse uses Claude matcher %q; want %q",
				e.Matcher, cfgCodex.ToolUpdatePlan,
			)
		}
	}
	for _, e := range codex[cfgCodex.EventPostToolUse] {
		alts := strings.Split(e.Matcher, "|")
		if slices.Contains(alts, "Edit") &&
			!slices.Contains(alts, cfgCodex.ToolApplyPatch) {
			t.Errorf(
				"Codex PostToolUse matcher %q covers Edit but not %q",
				e.Matcher, cfgCodex.ToolApplyPatch,
			)
		}
	}
}

// TestCodexHooksUserPromptSubmitOrder asserts the UserPromptSubmit
// commands match the Claude Code manifest exactly and in order
// (the nudge order is part of the contract: context size first,
// heartbeat last).
func TestCodexHooksUserPromptSubmitOrder(t *testing.T) {
	claude := tails(t, claudeHooksPath, readClaudeHooks(t), claudeHookAnchor)
	codex := tails(t, asset.PathCodexHooksJSON, readCodexHooks(t), cfgCodex.HookPrologue)

	var want, got []string
	for _, e := range claude[cfgCodex.EventUserPromptSubmit] {
		want = append(want, e.Tail)
	}
	for _, e := range codex[cfgCodex.EventUserPromptSubmit] {
		got = append(got, e.Tail)
	}
	if len(want) == 0 {
		t.Fatal("Claude manifest wires no UserPromptSubmit hooks")
	}
	if !slices.Equal(got, want) {
		t.Errorf(
			"UserPromptSubmit commands differ\n codex: %q\nclaude: %q",
			got, want,
		)
	}
}

// codexPluginManifest is the subset of .codex-plugin/plugin.json
// the guard inspects.
type codexPluginManifest struct {
	Name       string `json:"name"`
	Version    string `json:"version"`
	Skills     string `json:"skills"`
	Hooks      string `json:"hooks"`
	MCPServers string `json:"mcpServers"`
	Interface  struct {
		DisplayName string `json:"displayName"`
	} `json:"interface"`
}

// TestCodexPluginManifest asserts the plugin manifest names the
// ctx plugin, carries the repo version, points its component
// fields at files that exist in the embedded plugin root, and has
// a display name.
func TestCodexPluginManifest(t *testing.T) {
	data, err := FS.ReadFile(asset.PathCodexPluginJSON)
	if err != nil {
		t.Fatalf("read %s: %v", asset.PathCodexPluginJSON, err)
	}
	var m codexPluginManifest
	if parseErr := json.Unmarshal(data, &m); parseErr != nil {
		t.Fatalf("parse %s: %v", asset.PathCodexPluginJSON, parseErr)
	}

	if m.Name != cfgCodex.PluginName {
		t.Errorf("name = %q, want %q", m.Name, cfgCodex.PluginName)
	}
	if want := repoVersion(t); m.Version != want {
		t.Errorf("version = %q, want %q (VERSION file)", m.Version, want)
	}
	if strings.TrimSpace(m.Interface.DisplayName) == "" {
		t.Error("interface.displayName is empty")
	}

	for field, rel := range map[string]string{
		"hooks":      m.Hooks,
		"skills":     m.Skills,
		"mcpServers": m.MCPServers,
	} {
		if rel == "" {
			t.Errorf("%s field is empty", field)
			continue
		}
		embedded := path.Join(asset.DirCodex, path.Clean(rel))
		if _, statErr := fs.Stat(FS, embedded); statErr != nil {
			t.Errorf(
				"%s = %q does not resolve in the embedded plugin root (%s): %v",
				field, rel, embedded, statErr,
			)
		}
	}
}

// mcpServerEntry is one server in the plugin's .mcp.json map.
type mcpServerEntry struct {
	Command string   `json:"command"`
	Args    []string `json:"args"`
}

// TestCodexMCPServerMap asserts .mcp.json registers the ctx MCP
// server with the canonical launch command and arguments.
func TestCodexMCPServerMap(t *testing.T) {
	data, err := FS.ReadFile(asset.PathCodexMCPJSON)
	if err != nil {
		t.Fatalf("read %s: %v", asset.PathCodexMCPJSON, err)
	}
	var servers map[string]mcpServerEntry
	if parseErr := json.Unmarshal(data, &servers); parseErr != nil {
		t.Fatalf("parse %s: %v", asset.PathCodexMCPJSON, parseErr)
	}
	entry, ok := servers[cfgMcpServer.Name]
	if !ok {
		t.Fatalf("missing server %q; have %v", cfgMcpServer.Name, servers)
	}
	if entry.Command != cfgMcpServer.Command {
		t.Errorf("command = %q, want %q", entry.Command, cfgMcpServer.Command)
	}
	if want := cfgMcpServer.Args(); !slices.Equal(entry.Args, want) {
		t.Errorf("args = %q, want %q", entry.Args, want)
	}
}

// marketplaceCatalog is the subset of .agents/plugins/marketplace.json
// the guard inspects.
type marketplaceCatalog struct {
	Name     string `json:"name"`
	Metadata struct {
		Version string `json:"version"`
	} `json:"metadata"`
	Plugins []struct {
		Name   string `json:"name"`
		Source struct {
			Source string `json:"source"`
			Path   string `json:"path"`
		} `json:"source"`
		Policy struct {
			Installation   string `json:"installation"`
			Authentication string `json:"authentication"`
		} `json:"policy"`
		Category string `json:"category"`
	} `json:"plugins"`
}

// TestCodexMarketplace asserts the repo marketplace catalog names
// the marketplace ctx documents, lists exactly one plugin (ctx)
// sourced from the embedded plugin root, and carries the repo
// version.
func TestCodexMarketplace(t *testing.T) {
	data, err := os.ReadFile(marketplacePath)
	if err != nil {
		t.Fatalf("read %s: %v", marketplacePath, err)
	}
	var cat marketplaceCatalog
	if parseErr := json.Unmarshal(data, &cat); parseErr != nil {
		t.Fatalf("parse %s: %v", marketplacePath, parseErr)
	}

	if cat.Name != cfgCodex.MarketplaceID {
		t.Errorf("name = %q, want %q", cat.Name, cfgCodex.MarketplaceID)
	}
	if want := repoVersion(t); cat.Metadata.Version != want {
		t.Errorf(
			"metadata.version = %q, want %q (VERSION file)",
			cat.Metadata.Version, want,
		)
	}
	if len(cat.Plugins) != 1 {
		t.Fatalf("plugins: got %d entries, want exactly 1", len(cat.Plugins))
	}

	p := cat.Plugins[0]
	if p.Name != cfgCodex.PluginName {
		t.Errorf("plugins[0].name = %q, want %q", p.Name, cfgCodex.PluginName)
	}
	if p.Source.Source != cfgCodex.PluginVersionLocal {
		t.Errorf(
			"plugins[0].source.source = %q, want %q",
			p.Source.Source, cfgCodex.PluginVersionLocal,
		)
	}
	wantPath := "./" + path.Join("internal", "assets", asset.DirCodex)
	if p.Source.Path != wantPath {
		t.Errorf("plugins[0].source.path = %q, want %q", p.Source.Path, wantPath)
	}
	if p.Policy.Installation == "" {
		t.Error("plugins[0].policy.installation is empty")
	}
	if p.Policy.Authentication == "" {
		t.Error("plugins[0].policy.authentication is empty")
	}
	if p.Category == "" {
		t.Error("plugins[0].category is empty")
	}
}

// skillDirs returns the sorted skill directory names under an
// embedded skill tree.
func skillDirs(t *testing.T, tree string) []string {
	t.Helper()
	entries, err := FS.ReadDir(tree)
	if err != nil {
		t.Fatalf("read %s: %v", tree, err)
	}
	var names []string
	for _, e := range entries {
		if e.IsDir() {
			names = append(names, e.Name())
		}
	}
	slices.Sort(names)
	return names
}

// stripAllowedTools drops `allowed-tools:` lines, the transform
// hack/sync-codex-skills.sh applies when mirroring a skill.
func stripAllowedTools(body string) string {
	lines := strings.Split(body, "\n")
	kept := lines[:0]
	for _, line := range lines {
		if strings.HasPrefix(line, "allowed-tools:") {
			continue
		}
		kept = append(kept, line)
	}
	return strings.Join(kept, "\n")
}

// TestCodexSkillsMirrorClaude asserts the Codex skill tree is the
// Claude Code skill tree minus the documented exclusions, that
// every Codex SKILL.md exists and carries no `allowed-tools:`
// frontmatter, and that each body equals its Claude source after
// the sync transform (the in-process twin of `make
// check-codex-skills`).
func TestCodexSkillsMirrorClaude(t *testing.T) {
	claudeSkills := skillDirs(t, asset.DirClaudeSkills)
	codexSkills := skillDirs(t, asset.DirCodexSkills)

	var want []string
	for _, name := range claudeSkills {
		if !slices.Contains(codexOnlySkillExclusions, name) {
			want = append(want, name)
		}
	}
	if !slices.Equal(codexSkills, want) {
		t.Errorf(
			"Codex skill set differs from Claude minus exclusions; "+
				"run hack/sync-codex-skills.sh\n  got: %v\n want: %v",
			codexSkills, want,
		)
	}
	for _, excluded := range codexOnlySkillExclusions {
		if !slices.Contains(claudeSkills, excluded) {
			t.Errorf(
				"exclusion %q names no Claude skill; prune it here and in hack/sync-codex-skills.sh",
				excluded,
			)
		}
	}

	for _, name := range codexSkills {
		codexPath := path.Join(asset.DirCodexSkills, name, asset.FileSKILLMd)
		codexBody, readErr := FS.ReadFile(codexPath)
		if readErr != nil {
			t.Errorf("%s: %v", codexPath, readErr)
			continue
		}
		if string(codexBody) != stripAllowedTools(string(codexBody)) {
			t.Errorf("%s: carries Claude-only `allowed-tools:` frontmatter", codexPath)
		}

		claudePath := path.Join(asset.DirClaudeSkills, name, asset.FileSKILLMd)
		claudeBody, claudeErr := FS.ReadFile(claudePath)
		if claudeErr != nil {
			continue // already reported by the set comparison
		}
		if string(codexBody) != stripAllowedTools(string(claudeBody)) {
			t.Errorf(
				"%s is stale relative to %s; run hack/sync-codex-skills.sh",
				codexPath, claudePath,
			)
		}
	}
}

// TestCodexSkillReferencesShipped asserts that every references/
// path cited by a shipped Codex SKILL.md exists in the embedded
// asset tree. Guards against skills instructing agents to read
// files the sync script or embed globs failed to ship.
func TestCodexSkillReferencesShipped(t *testing.T) {
	entries, dirErr := fs.ReadDir(FS, asset.DirCodexSkills)
	if dirErr != nil {
		t.Fatalf("read codex skills dir: %v", dirErr)
	}
	re := regexp.MustCompile("references/[A-Za-z0-9._-]+")
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		name := entry.Name()
		content, readErr := FS.ReadFile(
			path.Join(asset.DirCodexSkills, name, asset.FileSKILLMd))
		if readErr != nil {
			t.Fatalf("%s: %v", name, readErr)
		}
		for _, ref := range re.FindAllString(string(content), -1) {
			refPath := path.Join(asset.DirCodexSkills, name, ref)
			if _, statErr := fs.Stat(FS, refPath); statErr != nil {
				t.Errorf(
					"%s cites %s but it is not embedded: %v",
					name, ref, statErr,
				)
			}
		}
	}
}

// TestClaudeRootDualManifest guards the dual-manifest defense: the
// Claude plugin root must carry a .codex-plugin manifest whose
// hooks entry points at a byte-copy of the canonical Codex hooks
// manifest, so a Codex that resolves the legacy
// .claude-plugin/marketplace.json still installs working hooks.
func TestClaudeRootDualManifest(t *testing.T) {
	manifestPath := filepath.Join(
		"claude", cfgCodex.DirPluginManifest, "plugin.json",
	)
	data, readErr := os.ReadFile(filepath.Clean(manifestPath))
	if readErr != nil {
		t.Fatalf("dual manifest missing: %v", readErr)
	}
	var manifest struct {
		Name    string `json:"name"`
		Version string `json:"version"`
		Hooks   string `json:"hooks"`
	}
	if jsonErr := json.Unmarshal(data, &manifest); jsonErr != nil {
		t.Fatalf("dual manifest parse: %v", jsonErr)
	}
	if manifest.Name != cfgCodex.PluginName {
		t.Errorf("name = %q, want %q", manifest.Name, cfgCodex.PluginName)
	}
	if manifest.Version != repoVersion(t) {
		t.Errorf("version = %q, want VERSION %q",
			manifest.Version, repoVersion(t))
	}
	if manifest.Hooks != "./hooks/codex.json" {
		t.Errorf("hooks = %q, want ./hooks/codex.json", manifest.Hooks)
	}

	claudeCopy, copyErr := os.ReadFile(filepath.Clean(
		filepath.Join("claude", "hooks", "codex.json"),
	))
	if copyErr != nil {
		t.Fatalf("hooks/codex.json missing: %v", copyErr)
	}
	canonical, canonErr := FS.ReadFile(asset.PathCodexHooksJSON)
	if canonErr != nil {
		t.Fatalf("embedded codex manifest: %v", canonErr)
	}
	if !bytes.Equal(claudeCopy, canonical) {
		t.Error("claude/hooks/codex.json diverges from " +
			"codex/hooks/hooks.json — run hack/sync-codex-skills.sh")
	}
}
