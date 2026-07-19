# context-model-and-state

## [2026-06-07-180006] CWD-anchored context model (consolidated)

**Consolidated from**: 5 entries (2026-04-13 to 2026-05-21)

- Walk boundary uses git as a hint, not a requirement: walkForContextDir
  consults findGitRoot to anchor ancestor .context candidates and falls back to
  CWD when no git is found — fixes nested-repo binding without making git
  mandatory or relying on unreliable project markers.
- ctx activate is strict-CWD (drop upward walk): state-setting commands follow
  git's read-vs-state pattern (read walks, state refuses to cross repo
  boundaries); workspace-shared layouts are preserved by user action (cd first),
  not inferred walk.
- Anchor ctx to CWD entirely: drop activate/deactivate, the env-var (CTX_DIR)
  resolver, and all walks. With .context/ mandated as .git/'s sibling, every
  resolver collapses to os.Stat; keeping any walk would force maintaining two
  implementations. Mental model matches helm/terraform/Claude Code; ~600-1000
  LOC net deletion (specs/cwd-anchored-context.md).
- Spec steps 1+2 (resolver swap + init-guard removal) merged into one commit
  because step 1 cannot compile without step 2; cleanest commit boundaries beat
  strict spec adherence — remaining steps stay discrete (4-commit
  decomposition, not the spec's 5).
- Substrate vs. artifact placement: cognitive substrate (read AND written via
  ctx-mediated paths) lives under .context/; project artifacts (read/edited
  directly by humans, e.g. specs/, CLAUDE.md, docs/) live at root. kb passes all
  three coupling tests (mediated queries, pipeline coupling, skill discipline)
  so it stays under .context/.

---

## [2026-06-07-180007] Encryption key resolution and migration (consolidated)

**Consolidated from**: 3 entries (2026-03-01 to 2026-06-02)

- Single global key at ~/.ctx/.ctx.key (matches ~/.claude/ convention); one key
  per machine covers ~99% of users. Replaced the over-engineered slug-based
  per-project key system; project-local key-next-to-ciphertext was a security
  antipattern that broke in worktrees. [Original 2026-03-01 entry was marked
  Superseded by the 2026-03-02 simplification.]
- Legacy-key auto-migration replaced with a stderr warning only: warn-only is
  simpler, avoids silent file operations, and keeps the (small, alpha) userbase
  in control; docs carry migration instructions.
- Removed the implicit project-local .context/.ctx.key auto-detection tier from
  ResolveKeyPath: resolution is now (1) explicit .ctxrc key_path, (2) global
  ~/.ctx/.ctx.key, (3) project-local only as a degenerate fallback when home is
  unavailable. The local tier was the only thing making worktrees differ from
  side-by-side terminals; its removal is net deletion, and the previously-silent
  fire-path decrypt failure is now surfaced.

---

## [2026-05-24-092912] Pad snapshot-on-mutate at the store.WriteEntries choke point

**Status**: Accepted

**Context**: Adding a safety net for accidental `ctx pad rm` (and any other
destructive pad mutation) required choosing where to insert the snapshot logic:
per-subcommand (in each cmd/<op>/run.go), or at the persistence choke point
(store.WriteEntriesWithIDs).

**Decision**: Pad snapshot-on-mutate at the store.WriteEntries choke point

**Rationale**: store.WriteEntriesWithIDs is invoked by every mutating pad
subcommand (add/edit/mv/rm/merge/normalize/resolve/tag and undo itself);
instrumenting it once gives universal coverage with one site of truth.
Per-subcommand instrumentation would need maintenance every time a new pad
mutation lands and is easy to forget. The snapshot itself is a byte-for-byte
copy of the existing pad blob (no re-encryption), so plaintext and encrypted
modes use identical logic; the existing ciphertext IS the snapshot.

**Consequence**: All future pad mutations get the safety net automatically
without per-command wiring. The op label for the snapshot filename is derived
from cmd.Name() at the call site, so the cmd parameter that already flowed in
for diagnostic output now carries semantic weight too. New constraint: any
future code path that bypasses WriteEntriesWithIDs to mutate the pad will
silently bypass the safety net — a guardrail test could enforce this if/when
that risk materializes.

---

## [2026-05-20-214753] Gitignore .context/handovers/; track only .gitkeep

**Status**: Accepted

**Context**: Per-session, operator-specific artifacts that grow without bound
and can leak host/internal identifiers into public
mirrors when the project's .context/ is committed.

**Decision**: Gitignore .context/handovers/; track only .gitkeep

**Rationale**: Aligns with the existing per-personal-state gitignore family
(journal, memory, state, logs, reminders.json, scratchpad.enc); the directory's
.gitkeep keeps the read-side missing-dir gate passing on fresh clones; the rest
of the closeout-fold pipeline already lives in .context/archive/closeouts/ which
IS tracked.

**Consequence**: ctx init template (internal/config/file/ignore.go) added
.context/handovers/* and !.context/handovers/.gitkeep; existing tracked
handovers untracked via git rm --cached but kept on disk; the 'handover is the
sole authoritative recall artifact' phrasing in KB-RULES.md still holds — it's
local-machine authoritative.

---

## [2026-04-11-180000] `Entry.Author` is server-authoritative, not client-authoritative

**Status**: Accepted

**Context**: The `Entry.Author` field on hub entries is copied verbatim from
the client's publish request (`handler.go:82`). It's optional, freeform, and
unauthenticated — a client with a valid token for project `alpha` can publish
entries claiming `Author: "bob@acme.com"` regardless of who actually
authenticated. This is the same spoofing pattern as `Origin` (audit finding
H-04) and was flagged as audit finding H-22 with three options: keep, drop,
override, or promote. The decision was never formally closed.

The premise that resolved it: **identity is eventually part of the token**.
Under the sysadmin-registry MVP, the server already knows `{user_id, project}`
from the authenticated token. Under the PKI stretch, the signed claim carries
identity cryptographically. In both models, the client has nothing to say about
authorship that the server doesn't already know with higher confidence.

**Decision**: `Entry.Author` is **server-authoritative**. The server stamps it
from the authenticated identity source on every publish. The client's
`pe.Author` input is ignored (or rejected — implementation choice, not
semantic difference). The field stays in the wire format but its semantics
change from "whatever the client said" to "whatever the server's auth layer
resolved."

Stamping source by phase:

- **Today (pre-registry)**: `Author = ClientInfo.ProjectName`, same source as
  the `Origin` server-enforcement fix (H-04). Lossy but consistent.
- **Registry MVP**: `Author = users.json` row's `user_id` (e.g.,
  `alice@acme.com`). Precise per-human attribution.
- **PKI stretch**: `Author = signed claim's sub field`. Cryptographic identity.

**Rationale**: Dropping the field is wrong because the registry MVP will
already give us a per-user identity to stamp — removing Author just to re-add
it later is churn. "Override" and "promote" are cosmetically different forms
of the same decision (server fills from auth context); "promote" is what
happens naturally once the registry MVP types the field as `UserID`.
Client-sourced Author is indefensible because it replicates the Origin
spoofing vector in a second field.

**Consequence**:

- The Author field stays on the wire and in `Entry{}`.
- Client-side code that populates `pe.Author` from local config becomes a
  no-op. Audit `ctx connect publish` and `ctx add --share` for any such
  code paths before the server-enforcement fix lands.
- `handler.go publish()` fills Author from the authenticated context (the
  same `ClientInfo` that H-04 pulls for Origin). Single unified
  auth-to-handler pipe.
- `docs/security/hub.md` "Compromised client token" section gets rewritten:
  attribution becomes **wrong** on compromise (attacker's token maps to
  attacker's identity), not **forgeable** (attacker cannot stamp someone
  else's name).
- The sysadmin-registry spec (`specs/hub-identity-registry.md`, tasked)
  MUST include a `user_id` field per row — it's the stamping source.
- Three open tasks collapse into one: H-22 resolves to "implement
  server-authoritative Author" instead of "decide Author fate." TASKS.md
  updated.

**Alternatives considered**:

- **Keep client-authoritative**: rejected. Same spoofing vector as Origin;
  trivially defeats any downstream attribution check.
- **Drop the field**: rejected. The registry MVP will need per-human
  attribution anyway. Dropping today is churn that gets undone
  immediately.
- **Override at client-side before publish**: rejected. Puts the security
  boundary on the wrong side of the trust zone. Must be server-side.

**Follow-up — client-advisory metadata**: the client still has useful
information to share that isn't an identity claim: a human-friendly
display name, the machine that made the publish, the tool version, a
CI system label, a team/role handle. This lives on a **new sibling
field `Meta`** (a `ClientMetadata` sub-struct), not on `Author`. The
separation of types is what protects the security property: `Author`
is reserved for server-authoritative identity, `Meta` is
client-advisory and explicitly labeled as such in any rendered
surface. `Meta` fields are size-capped individually (256 bytes) and
in aggregate (2 KB), validated for plain-string content (no
newlines, no control characters), and never claimed as attribution
in any API response. The renderer MUST label `Meta`-sourced values
with prose like "client label" or "client-reported" so readers
cannot mistake them for authoritative identity. See TASKS.md for
the implementation task.

---

## [2026-03-05-205424] Gitignore .context/memory/ for this project

**Status**: Accepted

**Context**: Memory mirror contains copies of MEMORY.md which holds strategic
analysis and session notes

**Decision**: Gitignore .context/memory/ for this project

**Rationale**: Strategic content should not be in git history. Docs updated to
say 'often git-tracked' for the general recommendation — this project is the
exception.

**Consequence**: Mirror and archives are local-only for this project. Other
projects can still track them. Sync and drift detection work the same way
regardless.

---



## [2026-03-01-222733] PersistentPreRunE init guard with three-level exemption

**Status**: Accepted

**Context**: ctx commands handled missing .context/ inconsistently — some
caught errors, some got confusing file-not-found messages, some produced empty
output

**Decision**: PersistentPreRunE init guard with three-level exemption

**Rationale**: Single PersistentPreRunE on root command gives one clear error.
Three-level exemption (hidden commands, annotated commands, grouping commands)
covers all edge cases without per-command boilerplate

**Consequence**: Boundary violation now returns an error instead of os.Exit(1),
making it testable. The subprocess-based boundary test was simplified to a
direct error assertion

---

---

## [2026-02-26-100002] ctx init and CLAUDE.md handling (consolidated)

**Status**: Accepted

**Consolidated from**: 3 decisions (2026-01-20)

- `ctx init` handles CLAUDE.md intelligently: creates if missing, backs up and
  offers merge if existing, uses marker comment for idempotency. The `--merge`
  flag enables non-interactive append.
- `ctx init` always generates `.claude/hooks/` alongside `.context/` with no
  flag needed. Other AI tools ignore `.claude/`; Claude Code users get seamless
  zero-config experience.
- Core tool stays generic and tool-agnostic, with optional Claude Code
  enhancements via `.claude/hooks/`. Other AI tools can be supported similarly
  (`ctx hook cursor`, etc.).

---

