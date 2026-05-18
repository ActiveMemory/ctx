//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package kb

// Subdirectory names under .context/.
const (
	// KBSubdir is the knowledge-base subdirectory name.
	KBSubdir = "kb"
	// IngestSubdir is the editorial-pipeline working subdirectory.
	IngestSubdir = "ingest"
	// SiteSubdir is the gitignored rendered-site output root.
	SiteSubdir = "site"
)

// Sub-paths under the KB / ingest / site / archive subdirectories.
const (
	// TopicsSubdir holds folder-shaped topic pages under .context/kb/.
	TopicsSubdir = "topics"
	// CloseoutsSubdir holds per-pass closeout files under
	// .context/ingest/.
	CloseoutsSubdir = "closeouts"
	// SchemasSubdir holds schema templates under .context/ingest/.
	SchemasSubdir = "schemas"
	// ArchiveCloseoutsSubdir holds folded closeouts under
	// .context/archive/.
	ArchiveCloseoutsSubdir = "closeouts"
	// SiteKBSubdir holds rendered KB output under .context/site/.
	SiteKBSubdir = "kb"
)

// KB-side filenames under .context/kb/.
const (
	// KBIndex is the kb landing page (carries scope + topic
	// managed block).
	KBIndex = "index.md"
	// Glossary is the kb-scoped domain glossary.
	Glossary = "glossary.md"
	// Contradictions records EV rows that disagree.
	Contradictions = "contradictions.md"
	// OutstandingQuestions records Q-### entries the kb cannot
	// yet answer.
	OutstandingQuestions = "outstanding-questions.md"
	// DomainDecisions holds kb-scoped decisions (distinct from
	// project-level DECISIONS.md; different schema, different
	// authority).
	DomainDecisions = "domain-decisions.md"
	// Timeline records dated events worth recording.
	Timeline = "timeline.md"
	// SourceMap is the discovery + admission map (kind, url, etc.).
	SourceMap = "source-map.md"
	// RelationshipMap records cross-topic + cross-source ties.
	RelationshipMap = "relationship-map.md"
	// TopicIndex is the per-topic landing filename under
	// topics/<slug>/.
	TopicIndex = "index.md"
)

// Ingest-side filenames under .context/ingest/.
const (
	// Rules is the editorial constitution (the kb's
	// counterpart to .context/CONSTITUTION.md, named to avoid
	// the collision per DECISIONS.md 2026-05-10).
	Rules = "KB-RULES.md"
	// GroundingSources lists external grounding inputs for
	// ctx kb ground.
	GroundingSources = "grounding-sources.md"
	// Prompt is the hand-fallback auto-router used when no
	// skill is installed.
	Prompt = "PROMPT.md"
	// Findings is the lazy-init destination for ctx kb note.
	Findings = "findings.md"
)

// Closeout-file shape: <TIMESTAMP>-<mode>-closeout.md per the
// generalized naming in specs/kb-editorial-pipeline.md (the
// upstream uses single-mode "ingest-closeout.md"; ctx is multi-mode).
const (
	// CloseoutSuffix is appended after the mode token.
	CloseoutSuffix = "-closeout.md"
	// CloseoutModeIngest names ingest-mode closeouts.
	CloseoutModeIngest = "ingest"
	// CloseoutModeAsk names ask-mode closeouts.
	CloseoutModeAsk = "ask"
	// CloseoutModeGround names ground-mode closeouts.
	CloseoutModeGround = "ground"
)

// Source-coverage ledger states. These string values are
// written verbatim into `.context/kb/source-coverage.md` rows
// and parsed back by the source-coverage writer. The allowed
// transitions between them live in
// [github.com/ActiveMemory/ctx/internal/write/kb/sourcecoverage].
const (
	StateDiscovered          = "discovered"
	StateAdmitted            = "admitted"
	StateHighlightsExtracted = "highlights-extracted"
	StatePartiallyIngested   = "partially-ingested"
	StateTopicPageDrafted    = "topic-page-drafted"
	StateComprehensive       = "comprehensive"
	StateSkipped             = "skipped"
)

// Confidence bands per the spec's confidence-laddering rules.
const (
	ConfidenceHigh        = "high"
	ConfidenceMedium      = "medium"
	ConfidenceLow         = "low"
	ConfidenceSpeculative = "speculative"
)

// EVIDPrefix is the prefix for evidence rows (`EV-###`).
const EVIDPrefix = "EV"

// EVIDDigits is the zero-padded width for evidence row numbers.
const EVIDDigits = 3

// EvidenceOnlyTag is the additive tag applied to rows minted in
// evidence-only mode passes. A future topic-page pass must
// re-read the source before citing such rows.
const EvidenceOnlyTag = "evidence-only"

// Topic-page template asset paths (under internal/assets/).
const (
	// AssetTemplatesIngest is the embedded-asset root for ingest
	// templates.
	AssetTemplatesIngest = "kb/templates/ingest"
	// AssetTemplatesSchemas is the embedded-asset root for ingest
	// schema templates.
	AssetTemplatesSchemas = "kb/templates/ingest/schemas"
	// AssetTemplateKBIndex is the embedded-asset path for the kb
	// landing template.
	AssetTemplateKBIndex = "kb/templates/kb/index.md"
	// AssetTemplateTopicIndex is the embedded-asset path for the
	// topic-page index.md template.
	AssetTemplateTopicIndex = "kb/templates/kb/topics/_template/index.md"
)
