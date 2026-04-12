---
#   /    ctx:                         https://ctx.ist
# ,'`./    do you remember?
# `.,'\
#   \    Copyright 2026-present Context contributors.
#                 SPDX-License-Identifier: Apache-2.0

title: Hub Security Audit
status: Draft
author: ctx team (post-PR #60 audit)
date: 2026-04-11
commit: ed8158b0
---

# Hub Security Audit

## TL;DR

**The `ctx` Hub is safe for the declared "trusted team on a
LAN" story and unsafe for anything beyond it.** The audit
found 30 issues ranging from critical (no transport
security, no cluster authentication) to informational
(minor error-message leaks). The critical and high-severity
findings cluster around three themes:

1. **No transport security anywhere**. Both the hub server
   and every client hard-code plaintext gRPC
   (`insecure.NewCredentials()`). There is no TLS option
   in the code. An operator who puts nginx in front of
   the hub cannot actually connect a client to the
   TLS-terminated endpoint because the client refuses to
   speak TLS.
2. **No identity layer**. Authentication is bearer-token
   only; tokens identify projects, not humans. Project
   names are user-asserted at registration and
   user-asserted again on every publish. `Origin` on
   entries is not cross-checked against the authenticated
   client's project name. The server knows which client
   is talking to it but never uses that information to
   gate publishes.
3. **Raft is fully unauthenticated**. The cluster
   transport is plaintext TCP with no authentication; any
   peer that can reach the Raft bind port can join the
   cluster, partition it, or DoS it.

The hub is **acceptable for the "personal cross-project
brain" and "small trusted team on a LAN" stories** — the
trust model assumes every participant is friendly and
every network segment is private. It is **unsafe for any
public-internet deployment**, and shipping one today
would be the "keeps you awake at night in a bad way"
scenario that motivated this audit.

The sysadmin-registry MVP tasked earlier (TASKS.md Hub
identity layer phase) addresses a significant subset of
these findings and is the right near-term work item. A
subsequent **signed-claim / PKI auth** phase is required
before any real public-internet deployment. Both tracks
are tasked; this spec makes the findings that motivate
them concrete and auditable.

## Metadata

| Field      | Value                               |
|------------|-------------------------------------|
| Audit date | 2026-04-11                          |
| Commit     | `ed8158b0` (main)                   |
| Scope      | `internal/hub/`, `internal/cli/connect/core/config/`, relevant callers |
| Method     | Static review of all `.go` files in scope; cross-reference against threat model below |
| Auditor    | ctx core team, during PR #60 follow-up |

## Scope

**In scope**:

- Hub server: `internal/hub/` — gRPC service, auth,
  storage, cluster, fan-out, replication.
- Hub client: `internal/cli/connect/` —
  registration, publish, subscribe, sync, listen.
- Client-side credential storage:
  `internal/cli/connect/core/config/` +
  `internal/crypto/` (reused).
- CLI entry points: `internal/cli/hub/`,
  `internal/cli/serve/` (the old static-site serve path,
  no longer hosting the hub but still relevant for
  history), `internal/cli/add/` (the `--share` flag).

**Out of scope**:

- Claude Code plugin auth (belongs to Claude Code).
- The ctx MCP server's non-hub tools (`ctx_search`,
  `ctx_steering_get`): inherit the client's local trust
  domain and are covered by MCP's own auth model.
- The encrypted scratchpad (`internal/crypto/` in
  isolation): already audited as a separate concern.
- ctx's UX affordances against accidental secret
  publishing: covered elsewhere (TASKS.md "secret-leak
  runbook" task).

## Threat model

### Trusted roles

- **Hub operator (sysadmin)**: has root on the hub host;
  holds the admin token; can read/write all hub files;
  can stop/start the hub; can hand-edit `clients.json`
  for revocation. Trusted without bound.
- **Client-side developer**: has local root on their
  workstation; holds one client token per registered
  project; publishes entries as themselves in good faith.
  Trusted as long as they don't deliberately abuse the
  system.

### Untrusted roles

- **Any network eavesdropper** on the path between a
  client and the hub (today: arbitrary, because there's
  no TLS; under a reasonable deployment: anyone on the
  same LAN or VPN).
- **Any attacker** outside the trust zone who can reach
  the hub's gRPC port.
- **Any attacker** who gains read access to the hub's
  data directory without root (e.g., via a misconfigured
  backup, a shared mount, or a container escape).
- **A compromised client workstation** (malware,
  unattended unlock, stolen laptop): the attacker can
  decrypt `.connect.enc` as the legitimate user and gain
  that user's hub token.

### Assets

- **Bearer tokens** (admin, client): grant hub access.
- **`entries.jsonl`**: append-only content; considered
  readable by every legitimate client that subscribes.
- **`clients.json`**: token registry; must remain private.
- **`admin.token`**: the operator's root credential.
- **Cluster state** (`raft/`): leadership metadata; not
  sensitive content, but integrity matters.
- **`.connect.enc`** (client-side): encrypted client
  token + hub address. Compromised client credentials
  let an attacker impersonate the legitimate user.

### Assumptions

- The LAN or VPN between client and hub is considered
  private by the declared "trusted team" story.
  Findings that depend on a hostile network
  (eavesdropping, man-in-the-middle) are **critical under
  a public-internet threat model** and **acceptable** for
  the LAN story only.
- The hub host is hardened against general workstation
  compromise (updates applied, unprivileged user, etc.).
  This audit does not assess host hardening.

## Per-story verdicts

### Story 1 — Personal cross-project brain (localhost hub)

**Verdict: Acceptable.** The network is `127.0.0.1`, the
same user owns the client and the server, and every
identity collapses to "me." The absence of TLS, identity
layer, and rate limiting are non-issues because there is
no adversary in the model.

One caveat: a compromised user account (malware) has
direct filesystem access and bypasses everything the hub
could protect. The hub is not a security boundary
against local-user compromise, and it isn't trying to
be.

### Story 2 — Small trusted team on a LAN

**Verdict: Acceptable with documented caveats.** The
trust-within-cluster model holds: everyone holding a
client token is friendly, the LAN is private, the hub
host is hardened. The caveats we should document:

- **Never expose the gRPC port to the public internet
  or to non-LAN subnets**, because there's no TLS and
  no rate limiting.
- **Treat `clients.json` on the hub host like
  `/etc/shadow`** — it holds every client token in
  plaintext. File-level protection (chmod 700 on the
  data dir, unprivileged user via systemd hardening) is
  the only defense.
- **Rotate the admin token periodically** and after any
  team member with sysadmin access leaves.
- **Accept that attribution is soft** — `Origin` is
  self-asserted and cannot be trusted for audit
  purposes.

The `hub-team.md` recipe already names these caveats in
prose; they're not surprising under the declared trust
model.

### Story 3 — Public-internet / multi-user

**Verdict: UNSAFE. Do not deploy.** Every critical
finding in this document applies, and several
high-severity findings compound catastrophically without
transport security. Shipping the current hub on the
public internet, even behind nginx, would leak
credentials (client TLS is impossible), expose the Raft
cluster to remote disruption, and allow unattributable
impersonation.

The path to "safe enough for Story 3" is:

1. Native TLS on both client and server (H-01, H-02).
2. Sysadmin-registry identity layer (TASKS.md Hub
   identity phase, resolves H-04, H-05, H-06, partial
   resolution of H-03).
3. Per-token rate limiting (H-08) and per-token listener
   caps (H-09).
4. Raft transport authentication or replacement
   (H-10, H-11).
5. Signed-claim / PKI auth (stretch task in Hub identity
   phase).

Until all five land, the declared Story 3 posture stays
"not supported."

## Findings

Findings are numbered H-NN and ordered by severity.
Within each severity bucket, ordering is by structural
importance.

### Critical (transport and network compromise)

#### H-01 — No server-side TLS

**Severity**: Critical (public-internet), Medium (LAN).

**Location**: `internal/hub/server.go:30`:

```go
gs := grpc.NewServer()
```

`grpc.NewServer()` is called with no TLS credentials and
no `grpc.Creds(...)` option. The server accepts only
plaintext gRPC connections.

**Impact**: All RPC traffic — including bearer tokens in
the `authorization` header and entry content — is sent
in the clear. Any network eavesdropper on the path
between client and server can capture tokens and entry
bodies. The declared "trusted team on a LAN" story
papers over this by assuming the LAN is private, which
is a real but fragile assumption.

**Recommendation**: Add a `--tls-cert` and `--tls-key`
flag pair to `ctx hub start`. When both are set, build
`credentials.NewTLS(...)` and pass it to
`grpc.NewServer` via `grpc.Creds`. Keep the plaintext
default for `localhost`-only Story 1 deployments.
Document the required upgrade for any non-localhost
deployment.

**Fix complexity**: Small (a few hours).

**Existing task coverage**: None. New task needed.

---

#### H-02 — No client-side TLS (plaintext hard-wired)

**Severity**: Critical (public-internet), High (LAN).

**Location**: Three sites, all identical:

- `internal/hub/client.go:30` (`NewClient`)
- `internal/hub/sync_helper.go:28` (`replicateOnce`)
- `internal/hub/failover.go:34` (`NewFailoverClient`)

All three use:

```go
grpc.WithTransportCredentials(insecure.NewCredentials())
```

**Impact**: The client **cannot** speak TLS even if the
server supports it. An operator who puts nginx in front
of the hub for TLS termination will find that clients
fail to connect (the client handshake doesn't negotiate
TLS, and the TLS terminator will reject or misrouter the
plaintext payload). This directly contradicts the
guidance in `docs/recipes/hub-multi-machine.md` which
recommends an nginx reverse proxy — **that recommendation
is currently un-implementable**.

**Recommendation**:

1. Introduce a `hub_addr` scheme discriminator:
   `grpc://host:port` for plaintext (current behavior),
   `grpcs://host:port` for TLS.
2. In `NewClient` and friends, parse the scheme and
   build the appropriate credential bundle
   (`insecure.NewCredentials()` vs
   `credentials.NewClientTLSFromCert(...)` with the
   system trust store).
3. Add an optional `--ca-cert` flag for self-signed
   deployments.
4. Update `hub-multi-machine.md` to show both forms
   (plain LAN and TLS-terminated via `grpcs://`).

**Fix complexity**: Medium (parsing + flag wiring + doc
updates).

**Existing task coverage**: None. New task needed
(paired with H-01).

---

#### H-10 — Raft transport unauthenticated

**Severity**: Critical (public-internet and LAN).

**Location**: `internal/hub/cluster.go:61-67`:

```go
transport, transErr := raft.NewTCPTransport(
    bindAddr, addr, 3,
    10*time.Second, os.Stderr,
)
```

`raft.NewTCPTransport` accepts plaintext TCP with no
peer authentication. Any host that can open a TCP
connection to the Raft bind port is treated as a
potential cluster member.

**Impact**:

- An attacker on the same network can impersonate a peer
  and cause leadership churn, partitions, or a DoS by
  flooding Raft RPCs.
- On the public internet, this is a remote unauthenticated
  DoS vector. The Raft port is typically `gRPC-port +
  1`, which is trivial to discover.
- Because the Raft FSM is a no-op (`leaderFSM` in
  `fsm.go`), a malicious peer cannot directly mutate
  data — but it **can** hijack leadership and then use
  its legitimate gRPC bearer credentials to publish
  entries that followers will replicate blindly (see
  H-13).

**Recommendation**:

1. Replace `raft.NewTCPTransport` with a TLS-wrapped
   transport using mutual-TLS between cluster peers.
   Peers authenticate via certificates issued from a
   cluster CA managed by the sysadmin.
2. Alternatively (simpler but weaker): gate the Raft
   port with a pre-shared secret that every peer
   presents in the Raft handshake. Not ideal but
   dramatically better than nothing.
3. Bind the Raft port to `127.0.0.1` by default and
   require an explicit `--raft-bind` flag to expose it
   on the network. Makes misconfiguration harder.

**Fix complexity**: Medium-High (Raft transport is
library-imposed; depends on what hashicorp/raft
supports).

**Existing task coverage**: None. New task needed.

---

#### H-11 — Raft transport unencrypted

**Severity**: Critical (public-internet), Medium (LAN).

**Location**: Same as H-10.

**Impact**: All Raft log operations traverse the network
in plaintext. Because the FSM is a no-op, the Raft log
carries no sensitive application data — but it does
carry cluster membership changes, leadership transfers,
and peer identity. An attacker sniffing inter-node
traffic can enumerate the cluster topology and
potentially craft malicious Raft messages based on the
observed state.

**Recommendation**: Addressed by H-10 fix (mTLS between
peers covers both authentication and encryption).

**Existing task coverage**: Merge with H-10.

---

#### H-13 — Replication trusts master without re-validation

**Severity**: Critical (if a master is ever compromised).

**Location**: `internal/hub/sync_helper.go:61-77`:

```go
for {
    msg := &EntryMsg{}
    if recvErr := stream.RecvMsg(msg); recvErr != nil {
        return
    }
    entry := Entry{
        ID:        msg.ID,
        Type:      msg.Type,
        Content:   msg.Content,
        Origin:    msg.Origin,
        ...
    }
    _, _ = store.Append([]Entry{entry})
}
```

The replication loop receives entries from the master
and calls `store.Append` directly. It does not:

- Re-validate the entry against `validateEntry` (type
  allowlist, size cap, required fields).
- Check that `Origin` matches any known project on the
  follower.
- Verify any signature or MAC over the entry.

**Impact**: A compromised master — which, thanks to H-10,
is trivially achievable by any attacker on the Raft
network — can inject arbitrary entries into every
follower's store. The attack surface is enormous:

- Publish entries with `Type` outside the allowlist that
  slip past client-side renderers.
- Publish entries with `Content` larger than 1MB,
  bypassing the per-RPC validation.
- Forge `Origin` to impersonate any project on any
  follower.
- Publish entries with back-dated timestamps to
  invalidate any ordering assumptions clients make.

**Recommendation**:

1. Call `validateEntry` on every replicated entry
   before appending.
2. Introduce a follower-side content sanity check (at
   minimum: non-empty ID, known type, content under
   limit).
3. Once H-10 is fixed, the "compromised master"
   premise is materially reduced, but defense-in-depth
   validation is still cheap.
4. Long-term: require the publishing client to sign
   entries before publish, and have followers verify
   the signature on replication. This eliminates the
   "trust the master" premise entirely at the cost of
   a per-client signing key.

**Fix complexity**: Small for steps 1-2; large for
step 4 (requires the PKI work in the identity-layer
stretch task).

**Existing task coverage**: None. New task needed.

---

### High (authentication, authorization, attribution)

#### H-03 — Plaintext token storage in `clients.json`

**Severity**: High.

**Location**: `internal/hub/types.go:49-53`, `internal/hub/store.go:146-163`.

```go
type ClientInfo struct {
    ID          string `json:"id"`
    ProjectName string `json:"project_name"`
    Token       string `json:"token"`   // ← stored in plaintext
}
```

`RegisterClient` calls `saveJSON(clientsPath, s.clients)`
which writes the struct verbatim, including the
plaintext `Token` field.

**Impact**: File-level compromise of
`<data-dir>/clients.json` yields every client token in
the clear. This turns a **host compromise** into a
**total hub compromise**, because the attacker can now
impersonate every registered project.

**Recommendation**:

1. Hash tokens with a strong password hash (argon2id or
   bcrypt) before persisting. Only plaintext token
   leaves the server at registration time; only hash
   ever touches disk.
2. Validate tokens by hashing the presented value and
   comparing to the stored hash.
3. Migrate existing `clients.json` deployments via a
   one-shot migration: read old file, hash each token,
   rewrite. Print a warning to the operator that all
   existing tokens will continue working but cannot be
   recovered from the file (so the sysadmin should note
   them before migration).

**Fix complexity**: Small (argon2 is in `golang.org/x/crypto`,
the validation loop is a one-line change).

**Why hashing and not encryption with the pad key?**
A reasonable-sounding alternative is to reuse the existing
global encryption key at `~/.ctx/.ctx.key` (the one that
protects `ctx pad` entries) to encrypt the `Token` field
in `clients.json` at rest. Rejected, for two reasons:

1. **Same-host compromise defeats it.** The pad key lives
   on the same machine as `<data-dir>/clients.json`. The
   threat model H-03 defends against is *file-level host
   compromise yields every token in the clear* — and an
   attacker who can read the data dir almost certainly
   can also read `~/.ctx/`. Co-located key + ciphertext
   is not meaningfully better than plaintext against that
   threat. Hashing with argon2id is, because re-deriving
   each token still costs argon2id work per attempt with
   no shortcut.
2. **Wrong primitive for the job.** Tokens are
   server-side credentials the server only needs to
   *verify* (compare a presented value to a stored
   record). They never need to be recovered to plaintext
   on the server. Encryption is the right primitive when
   the server must read the plaintext back (which is why
   `ctx pad` uses it — the user wants to see their note).
   For verify-only, hashing is strictly better: it
   removes a recovery path that doesn't need to exist and
   sidesteps a future key-rotation problem.

Encryption-of-tokens would help if and only if the key
lived **off-host** (OS keyring, HSM, external secrets
manager). That's the "behind the local keyring" half of
the brainstorm follow-up note below — a separate, larger
piece of work than H-03, tracked but not blocking.

**Existing task coverage**:
`#### Design follow-ups surfaced by the brainstorm
(2026-04-11)` → "Hash `clients.json` tokens or move them
behind the local keyring". Re-affirmed here with
concrete remediation.

---

#### H-04 — `Origin` not server-enforced on publish

**Severity**: High.

**Location**: `internal/hub/handler.go:62-98`,
`internal/hub/validate.go:23-47`,
`internal/hub/store.go` (`ValidateToken`).

The publish handler copies `pe.Origin` verbatim from the
client's PublishRequest into the stored Entry. The
authenticated client's `ClientInfo.ProjectName` is
available inside `validateBearer` but is never attached
to the request context, so the handler has no way to
cross-check `pe.Origin` against the authenticated
identity.

**Impact**: A client with a valid token for project
`alpha` can publish entries with `Origin: beta` and the
hub will store them as if they came from `beta`. This
allows:

- **Accidental mislabeling**: a client config bug (typo,
  stale script) silently poisons the attribution stream.
- **Deliberate impersonation**: any client with a valid
  token can impersonate any other project, including
  projects they are not registered for.
- **Attribution washing**: entries published via a
  compromised low-privilege client can appear to come
  from a high-privilege one, defeating any
  after-the-fact "who published this?" investigation.

**Current (broken) user-to-server trace**:

1. User runs `ctx connect publish --type learning "..."`.
2. Client reads `.context/.connect.enc`, pulls out the
   bearer token and local project name, builds a
   `PublishRequest{Origin: "<local project>", ...}`,
   attaches `authorization: Bearer <token>` to gRPC
   metadata, sends it.
3. Auth interceptor calls `validateBearer`. It reads
   the header, strips `Bearer `, and calls
   `store.ValidateToken(token)`. That returns
   `error` only — **the matched `ClientInfo` is
   discarded**.
4. `handler.publish()` copies `pe.Origin` verbatim into
   the stored `Entry` (`handler.go:81`). No cross-check.
5. Entry is persisted and replicated under whatever
   `Origin` the client claimed.

The break is between steps 3 and 4: the auth layer knows
which project the token maps to but throws that
knowledge away before the handler runs.

**Recommendation (fix shape)**:

1. **`store.go` — promote `ValidateToken` to
   `LookupToken`.** Today it returns `error`; change it
   to `LookupToken(token string) (*ClientInfo, error)`.
   There is only one caller (`validateBearer`), so the
   rename is cheap; prefer promoting the existing method
   over adding a sibling.

2. **`validate.go` — return the enriched context.**
   `validateBearer` changes signature from
   `(ctx, store) error` to
   `(ctx, store) (context.Context, error)`. On success,
   it attaches the looked-up `*ClientInfo` to the
   context via `context.WithValue` under a private
   package key:

   ```go
   type ctxKey int
   const clientInfoKey ctxKey = 1

   func validateBearer(
       ctx context.Context, store *Store,
   ) (context.Context, error) {
       // ... same metadata parsing ...
       info, err := store.LookupToken(token)
       if err != nil {
           return ctx, status.Error(
               codes.Unauthenticated, "invalid token",
           )
       }
       return context.WithValue(
           ctx, clientInfoKey, info,
       ), nil
   }

   func clientFromContext(
       ctx context.Context,
   ) *ClientInfo {
       v, _ := ctx.Value(clientInfoKey).(*ClientInfo)
       return v
   }
   ```

   The gRPC auth interceptor that calls `validateBearer`
   must propagate the **returned** context to the
   downstream handler — that's how metadata rides
   through the middleware chain.

3. **`handler.go` — overwrite, don't trust.** In
   `publish()`, replace `Origin: pe.Origin` with a
   value pulled from the context:

   ```go
   client := clientFromContext(ctx)
   for i, pe := range req.Entries {
       entries[i] = Entry{
           ID:        pe.ID,
           Type:      pe.Type,
           Content:   pe.Content,
           Origin:    client.ProjectName, // server-enforced
           Meta:      pe.Meta,
           Timestamp: time.Unix(pe.Timestamp, 0),
       }
   }
   ```

   `client` is guaranteed non-nil here because the
   interceptor rejects unauthenticated requests before
   the handler runs. If a future refactor ever allows
   unauthenticated publish, this becomes a panic-on-nil
   bug — add a defensive nil-check in that future diff,
   not this one.

4. **Silent rewrite vs. loud reject.** The client's
   `pe.Origin` field is now advisory. Three options:

   - **(a) Silent overwrite**: ignore `pe.Origin`
     entirely. Simplest. Audit-recommended default.
   - **(b) Reject mismatch**: if
     `pe.Origin != "" && pe.Origin != client.ProjectName`,
     return `codes.InvalidArgument`. Surfaces client
     config bugs loudly (typo fails fast instead of
     getting silently rewritten).
   - **(c) Warn**: log a server-side warning on
     mismatch but persist the corrected value.

   Prefer **(b)** for the initial implementation: it
   catches client-side configuration drift that (a)
   would hide. The cost is a handful of extra lines; the
   benefit is that a broken CI script fails visibly
   instead of silently corrupting attribution. Drop to
   (a) only if real-world friction shows it.

5. **Test — the whole point.** Regression test in
   `internal/hub/handler_test.go` (new file or existing):

   ```go
   func TestPublish_OriginServerEnforced(t *testing.T) {
       // Register project "alpha", get token T_alpha.
       // Register project "beta",  get token T_beta.
       // Publish with T_alpha but Origin: "beta".
       // Assert: under policy (b) the RPC returns
       //         InvalidArgument; under (a) the stored
       //         entry has Origin == "alpha".
   }
   ```

   Without this test the fix is a promise. With it, any
   future refactor that breaks the wiring fails CI.

**Fixed user-to-server trace (post-fix)**:

1. Client sends `PublishRequest{Origin: "beta"}` with
   `T_alpha` in the header.
2. Auth interceptor calls `validateBearer`, which
   looks up `T_alpha` → `ClientInfo{ProjectName: "alpha", ...}`
   and attaches it to the context.
3. `handler.publish()` reads
   `clientFromContext(ctx) → {ProjectName: "alpha"}`.
   Under policy (b), rejects with `InvalidArgument`
   because `"beta" != "alpha"`. Under policy (a),
   overwrites every entry's `Origin` to `"alpha"` and
   persists.
4. Stored entry (if persisted): `Origin: "alpha"`. The
   client's `"beta"` claim is either refused or
   discarded; never persisted.

**Threat-model caveat the fix does NOT address**:

Under today's single-admin-token deployment model,
server-enforced Origin does **not** prevent a legitimate
`alpha` client from publishing garbage *as alpha*. It
only prevents `alpha` from publishing *as beta*. This
collapses the "any valid token can impersonate any
project" class of attack into "a compromised token can
only impersonate its own project" — a large improvement,
but not a complete fix for attribution integrity. The
remaining holes:

- **H-05** (project squatting at registration): the
  admin-token holder can pre-register `production` and
  then legitimately publish as `production`.
- **H-06** (no user identity layer): attribution is
  project-scoped, not human-scoped — two developers on
  the same project are indistinguishable.

Both close when the sysadmin-registry MVP lands (Hub
identity layer phase, TASKS.md).

**Pairing with H-22a — share the commit**:

H-04 (Origin) and H-22a (Author) are the **same
plumbing**: both need `ClientInfo` on the context, both
need `handler.go publish()` to stamp a field from that
context instead of from the request. Land them in a
single commit. The commit message references both
`specs/hub-security-audit.md` H-04/H-22a and
`.context/DECISIONS.md [2026-04-11-180000]`. Under the
pre-registry model, both fields stamp from
`client.ProjectName`; under the registry MVP, Author
upgrades to `client.UserID` while Origin stays on
`client.ProjectName`.

**Fix complexity**: Small. Concrete line counts:

- `store.go`: ~3 lines (signature change, return
  `*ClientInfo`).
- `validate.go`: ~10 lines (signature change,
  context-value helper, `clientFromContext` accessor).
- Auth interceptor site: 1-line change (use returned
  context for downstream handler).
- `handler.go publish()`: 2 lines (pull client, swap
  `Origin` and `Author` sources).
- Test: ~40 lines (new regression test).

Total: <60 lines, single commit, single spec reference.

**Existing task coverage**:
"Server-enforce `Origin` on publish". Already captured
in the PR #60 follow-up section. This audit upgrades its
framing from "small cheap fix" to "High-severity
finding" and adds the H-22a pairing requirement.

---

#### H-05 — Project squatting at registration (no allowlist)

**Severity**: High.

**Location**: `internal/hub/handler.go:20-59` (`register`),
`internal/hub/store.go:146-163` (`RegisterClient`).

The register handler validates the admin token and
non-empty `ProjectName`, then calls `RegisterClient`
which only rejects **duplicate** names. Any non-duplicate
project name is accepted and persisted.

**Impact**: A user who holds the admin token can register
any non-duplicate project name — including names they
have no legitimate claim to. Under the current deployment
model (single shared admin token), an insider attacker
can:

1. Register `--project production-platform` before the
   real production team does.
2. Receive a client token for `production-platform`.
3. Publish entries tagged as `production-platform`.
4. Every subscribed client receives the forgeries as if
   they came from production.

This is especially dangerous in combination with H-04:
first squat the name, then publish entries that are
cryptographically indistinguishable from legitimate
production entries because there's no server-side
attribution check.

**Recommendation**: Resolved by the sysadmin-registry
MVP in the Hub identity layer phase. Pre-seeded
`users.json` rejects registrations that don't match a
pre-declared `{user, project}` pair.

Short-term mitigation (before the registry ships):
document in `hub-multi-machine.md` that the admin token
MUST be held only by the sysadmin, never shared with
team members.

**Existing task coverage**: Sysadmin-registry MVP
(Hub identity layer phase, TASKS.md).

---

#### H-06 — No user identity layer

**Severity**: High.

**Location**: Pervasive — there is no concept of "user"
anywhere in the hub data model. `ClientInfo` identifies
projects, not humans.

**Impact**: Consequences include:

- Attribution cannot distinguish two developers on the
  same team.
- Revocation is project-scoped, not user-scoped:
  removing a departed employee requires manual `ls`
  through `clients.json` to find all rows where they
  might have registered, with no way to know for sure.
- No audit trail can answer "who published this?"
  beyond "the holder of token X at the time."

**Recommendation**: Sysadmin-registry MVP adds
per-user rows to `users.json`. Subsequent PKI phase adds
cryptographic user identity.

**Existing task coverage**: Hub identity layer phase
(both MVP and PKI stretch tasks).

---

#### H-08 — No per-token rate limiting

**Severity**: High.

**Location**: `internal/hub/grpc.go` (service descriptor
and handler wiring). No rate-limiting middleware is
registered.

**Impact**: A valid client token can call Publish at
maximum TCP throughput. With `maxContentLen = 1 MB`, a
sustained publisher can commit 1 GB/s of storage
pressure (modulo gRPC max message size, which defaults
to 4 MB per message). A misbehaving client — not even a
malicious one, just a buggy retry loop — can fill the
disk in minutes.

**Impact is amplified by H-15 (append-all-rewrite)**:
every publish rewrites the entire `entries.jsonl` file,
so the effective cost of one publish scales with the
total store size. A publish flood becomes quadratic.

**Recommendation**:

1. Add a token-keyed rate limiter on Publish. A
   reasonable starting point: 10 entries/sec per token,
   100 entries/sec burst.
2. Per-token cumulative byte quota per day (e.g., 100
   MB/day) as a coarser secondary limit.
3. Rate-limit Listen stream opens per token (cap
   concurrent streams per token at, say, 4).
4. Return `codes.ResourceExhausted` with a Retry-After
   hint so clients back off gracefully.

**Fix complexity**: Small (there are good Go
rate-limiting libraries; `golang.org/x/time/rate` is in
the stdlib extended tree).

**Existing task coverage**: None. New task needed.

---

#### H-09 — No per-client listener limit

**Severity**: High.

**Location**: `internal/hub/fanout.go` (`subscribe`,
`broadcast`).

`subscribe` allocates a 64-entry buffered channel for
every caller. `broadcast` sends to every subscribed
channel. There is no cap on the number of concurrent
subscribers per token or in total.

**Impact**: A single client with a valid token can open
N Listen streams. Each stream holds a 64-entry buffer of
`[]Entry`. Memory cost is bounded per listener but
unbounded in N. An attacker can exhaust hub memory by
opening thousands of Listen streams.

Secondary impact: every Listen stream is a goroutine,
so CPU scheduling overhead grows linearly in N too.

**Recommendation**: Cap concurrent listeners per token
(e.g., 4) and total concurrent listeners (e.g., 256).
Reject further subscribe attempts with
`codes.ResourceExhausted`. Track per-token counts in the
`fanOut` struct.

**Fix complexity**: Small.

**Existing task coverage**: None. New task needed
(paired with H-08).

---

#### H-17 — No batch size cap on PublishRequest

**Severity**: High.

**Location**: `internal/hub/types.go:142-144`
(`PublishRequest.Entries []PublishEntry`),
`internal/hub/entry_validate.go:28-53`.

`validateEntry` checks size per entry (1 MB) but nothing
checks the number of entries in a single PublishRequest.
The gRPC default `MaxRecvMsgSize` of 4 MB caps the total
wire size, but:

- A batch of 4000 entries at 1 KB each fits in 4 MB
  and creates 4000 disk writes (amplified by H-15).
- Increasing MaxRecvMsgSize (common for high-throughput
  deployments) makes this worse.

**Impact**: A single Publish call can trigger
disproportionate server-side work. Combined with H-08
(no rate limiting), a client can cause sustained
server-side load with one RPC per second.

**Recommendation**: Cap `PublishRequest.Entries` at a
reasonable number (e.g., 32 per request) and reject
larger batches with `codes.InvalidArgument`. Document the
limit.

**Fix complexity**: Small.

**Existing task coverage**: None. New task needed.

---

#### H-18 — No audit log of RPC operations

**Severity**: High.

**Location**: Pervasive — no structured logging of RPC
invocations exists beyond gRPC's default error
logging.

**Impact**: After a security incident, there is no way
to answer:

- Who called Publish on entry seq 12345?
- When was client X's token last used?
- How many Register attempts failed between 02:00 and
  03:00?
- What RPCs did a now-revoked token issue before
  revocation?

`entries.jsonl` records the result of successful
Publish calls, but nothing records failed attempts,
authentication failures, Sync cursor positions, or any
other evidence of activity.

**Recommendation**: Create `audits.jsonl` next to
`entries.jsonl`. Every RPC appends one line containing:

```json
{"ts":"2026-04-11T18:22:01Z","method":"Publish","user":"alice@acme.com","project":"alpha","status":"ok","entry_count":3}
```

- Authentication failures logged with `status: "deny"`
  and the reason.
- Retention independent of `entries.jsonl` (sysadmin
  can rotate audit logs on a shorter schedule).
- Available to `ctx hub status --audit` for
  operator-side inspection.

**Fix complexity**: Small-Medium (new file handle, log
rotation strategy).

**Existing task coverage**: Listed in the Hub identity
layer phase ("Add per-user audit log"). Re-affirmed here.

---

#### H-19 — No revocation mechanism

**Severity**: High.

**Location**: No code exists for token revocation. The
only path is manual JSON editing: stop the hub, edit
`clients.json`, restart.

**Impact**: When an employee leaves, when a workstation
is stolen, when a token is suspected compromised — the
operator has to hand-edit a file and restart the hub.
This is:

- **Error-prone** (wrong project name, JSON syntax
  error, forgot to restart).
- **Slow** (seconds to minutes of downtime).
- **Invisible** (no audit trail of what was revoked or
  when).

**Recommendation**:

1. Add `ctx hub users remove <id>` subcommand that
   edits `users.json` (under the identity-layer MVP) or
   `clients.json` (under the current model) and signals
   the running hub to reload.
2. Hub watches the registry file via `fsnotify` or
   polls it every N seconds. Reloading is non-fatal:
   any token no longer in the registry starts failing
   on the next RPC.
3. Revocation events go to `audits.jsonl` (H-18).

**Fix complexity**: Small-Medium (file watching + reload
state machine).

**Existing task coverage**: Partially covered by the
Hub identity layer phase ("Add `ctx hub users`
subcommand group"). Re-affirmed with concrete
remediation.

---

#### H-22 — `Entry.Author` is unauthenticated freeform

**Severity**: High.

**Location**: `internal/hub/types.go:33-41`
(`Entry.Author`), `internal/hub/handler.go:83`
(publish copies it verbatim).

The `Author` field exists on the Entry wire format and
is optional. The publish handler copies it verbatim
from client input with no validation.

**Impact**: Impersonation on display: a client can
publish entries claiming `Author: "alice@acme.com"`
regardless of who actually authenticated. Any downstream
system that reads `.context/hub/*.md` and trusts the
Author field for attribution is trivially defeated.

**Recommendation**: Three options, decision required
(already tasked in the design follow-ups):

1. **Drop** — remove the field entirely. Simplest, least
   surprising. Recommended unless there's a concrete
   use case for per-entry attribution.
2. **Override** — like H-04, the server overwrites
   `Author` with the authenticated client's identity
   (user_id under the registry MVP; ProjectName under
   the current model).
3. **Promote** — make `Author` a first-class identity
   field backed by the identity layer. Requires H-06
   resolution.

**Fix complexity**: Small (drop or override), Medium
(promote).

**Existing task coverage**:
"Decide the fate of `Entry.Author`". Re-affirmed with
explicit options.

---

### Medium (correctness, operational integrity)

#### H-07 — No token TTL / expiry

**Severity**: Medium.

**Location**: `internal/hub/types.go` (`ClientInfo`
struct has no expiration field), `internal/hub/store.go
ValidateToken` (no expiry check).

Tokens are valid forever until manually removed from
`clients.json`.

**Impact**: Dormant tokens accumulate. A leaked token
from three years ago is still valid today. There is no
hygiene pressure to rotate.

**Recommendation**: Optional `expires_at` field per
row. Tokens without it keep current semantics; tokens
with it are rejected after the timestamp. See the TTL
decision in the identity-layer spec task.

**Existing task coverage**: TTL decision documented in
the Hub identity layer phase spec task.

---

#### H-12 — Raft bootstrap race condition

**Severity**: Medium.

**Location**: `internal/hub/cluster.go:88-117`:

```go
if len(peers) == 0 {
    ...
    r.BootstrapCluster(config)
} else {
    ...
    r.BootstrapCluster(config)
}
```

`BootstrapCluster` is called unconditionally on every
node, every startup. In a multi-node cluster, if all
three nodes start simultaneously with `--peers`, each
calls BootstrapCluster with its own initial config.
This produces racing bootstraps.

**Impact**: Cluster might enter a split-brain state on
first startup, where two nodes each think they're the
bootstrap authority. The hashicorp/raft library detects
this on subsequent log operations but the initial state
is undefined.

**Recommendation**:

1. Mark a single node as the bootstrap leader via a
   `--bootstrap` flag, and only that node calls
   BootstrapCluster. Other nodes use
   `AddVoter`/`AddNonvoter` to join an existing cluster.
2. After first successful bootstrap, persist a
   `bootstrapped` flag in the raft dir and skip
   BootstrapCluster on subsequent starts.

**Fix complexity**: Small-Medium (state tracking in the
cluster dir).

**Existing task coverage**: None. New task needed.

---

#### H-14 — Sequence re-assignment on replication

**Severity**: Medium.

**Location**: `internal/hub/sync_helper.go:66-76`
(builds an Entry from the replicated message),
`internal/hub/store.go:Append` (overwrites Sequence).

Followers replicate master entries through
`store.Append`, which unconditionally assigns a new
sequence from the follower's local counter, discarding
the master's sequence.

**Impact**: The same entry has different sequence
numbers on different nodes. A client that fails over
from node A (cursor at seq 42) to node B (same entry at
seq 17) resumes from a meaningless cursor position and
ends up re-replicating entries it has already seen.

This is a **correctness** finding with security
implications: it also means the append-only log is not
a reliable cross-node reference for "which entries does
the cluster consider canonical."

**Recommendation**:

1. Add a `masterSequence` field to Entry that preserves
   the master-assigned sequence across replication.
2. Clients track both local and master sequence in
   their `.sync-state.json` and cursor by master
   sequence when talking to any node.
3. Alternative: derive deterministic sequences from the
   entry ID (e.g., `murmur3(ID) % 2^32`), so all nodes
   agree on the same sequence for the same entry. More
   complex but cleaner.

**Fix complexity**: Medium.

**Existing task coverage**: None. New task needed.

---

#### H-15 — `appendFile` is not atomic and not actually appending

**Severity**: Medium (data integrity), Low (security).

**Location**: `internal/hub/persist.go:73-82`:

```go
func appendFile(path string, data []byte) error {
    existing, readErr := io.SafeReadUserFile(path)
    if readErr != nil && !os.IsNotExist(readErr) {
        return readErr
    }
    return io.SafeWriteFile(
        path, append(existing, data...), fs.PermFile,
    )
}
```

The function:

1. Reads the **entire** existing file into memory.
2. Concatenates the new data in Go memory.
3. Writes the whole thing back via `SafeWriteFile`
   (which does atomic temp+rename).

**Impact**:

- **Performance**: every publish is O(N) in total store
  size. At 10k entries, 1 KB each, every publish reads
  and writes 10 MB. Per-publish throughput collapses as
  the store grows.
- **Correctness under pressure**: SafeWriteFile is
  atomic at the filesystem level, but the read-all-then-
  rewrite pattern defeats the "append-only" semantics:
  if a read fails due to memory pressure, the new write
  can overwrite a short version of the file and
  **lose all historical entries**.
- **Disk usage**: temp+rename doubles peak disk usage.

**Recommendation**: Replace with a true append using
`os.OpenFile(path, os.O_APPEND|os.O_WRONLY|os.O_CREATE,
0600)`. Atomic append on POSIX file systems for writes
shorter than `PIPE_BUF` (4096 bytes). For longer writes
(multi-entry batches), use a temp file + rename of the
appended bytes only, not the whole file.

**Fix complexity**: Small (standard Go pattern).

**Existing task coverage**: None. New task needed.

---

#### H-24 — No redaction mechanism for published entries

**Severity**: Medium.

**Location**: No code exists for entry deletion. The
store is append-only by design.

**Impact**: Accidental publication of a secret, PII, or
legally problematic content (e.g., a name that must be
removed under GDPR right-to-be-forgotten) has no
supported remediation. The operator has to stop the hub,
hand-edit `entries.jsonl`, manually update
`meta.json` sequence counter, and restart — all manual,
all error-prone. Every client still has the entry in
its local `.context/hub/*.md` mirror; there's no signal
to re-sync.

**Recommendation**:

1. Add `ctx hub redact <seq>` subcommand that:
   - Marks the entry as redacted in an `entries_redacted.jsonl`
     sidecar (append-only).
   - Broadcasts a redaction notification via the Listen
     stream so subscribed clients drop the entry from
     their local mirrors.
2. On query, filter out redacted entries.
3. Log the redaction to `audits.jsonl` (H-18) with the
   operator's user_id.

**Fix complexity**: Medium.

**Existing task coverage**: None. New task needed.

---

#### H-29 — No bound on in-memory entry cache

**Severity**: Medium.

**Location**: `internal/hub/types.go:78-85`
(`Store.entries []Entry`),
`internal/hub/store.go:loadEntries`.

The store loads the entire `entries.jsonl` into an
in-memory slice on startup and grows it unbounded on
every Publish.

**Impact**: For a long-lived hub, memory grows linearly
with total entries published. At 10M entries × 1 KB
average, the hub holds 10 GB resident. OOM risk on
memory-constrained hosts; no way to bound it.

**Recommendation**:

1. Implement an LRU cache over entries.jsonl instead of
   a full in-memory slice. Query pages through the file
   with an offset index.
2. Persistent index file
   (`entries.idx` → sequence → byte offset) for O(log N)
   seeks without reading the whole file.
3. Secondary: entries.jsonl rotation at a threshold.
   Oldest entries move to `entries-YYYY-MM.jsonl.gz`.

**Fix complexity**: Medium-Large.

**Existing task coverage**: None. New task needed.

---

#### H-30 — No connection-level backpressure

**Severity**: Medium.

**Location**: `internal/hub/server.go:30` (`grpc.NewServer`
called with no options).

Default gRPC server has no connection limit, no
per-connection stream limit, no idle timeout, no
keepalive enforcement.

**Impact**: A slow-loris style attack (open many
connections, send partial messages) can exhaust
goroutines and file descriptors. Idle clients that never
send anything keep connections open indefinitely.

**Recommendation**:

1. Set `grpc.KeepaliveEnforcementPolicy` requiring
   clients to ping every N seconds or be dropped.
2. Set `grpc.KeepaliveParams` to close idle connections.
3. Set `grpc.MaxConcurrentStreams` per connection.
4. Cap total concurrent connections at the process
   level (not gRPC itself — a listener wrapper).

**Fix complexity**: Small.

**Existing task coverage**: None. New task needed
(paired with H-08 and H-09 as the "rate limiting and
backpressure" track).

---

### Low (defense-in-depth, hardening)

#### H-16 — Content not sanitized for markdown injection

**Severity**: Low.

**Location**: `internal/hub/entry_validate.go:28-53`
(validation), `internal/cli/connect/core/render/render.go`
(client-side renderer).

`validateEntry` does not inspect `Content` beyond size.
When the client-side renderer writes entries to
`.context/hub/decisions.md` (and siblings), the content
is concatenated into a markdown document. A malicious
Content field like `\n---\ntitle: Fake\n---\n# Impostor`
can inject what looks like a new frontmatter block and
a new entry.

**Impact**: A publisher with legitimate credentials can
confuse client-side parsers. Not a server-side escalation,
but a way to manipulate what readers see in their
`.context/hub/*.md` files.

**Recommendation**: Escape or fence Content when
rendering on the client. Either wrap every entry's body
in explicit markers (`<!-- BEGIN ENTRY seq=42 -->` /
`<!-- END ENTRY seq=42 -->`) or encode known-disruptive
patterns (triple-dash lines at the start of a new line).

**Fix complexity**: Small.

**Existing task coverage**: None. New task needed.

---

#### H-20 — Token validation is not fully constant-time

**Severity**: Low (theoretical timing side-channel).

**Location**: `internal/hub/store.go:173-188`:

```go
func (s *Store) ValidateToken(bearerToken string) *ClientInfo {
    idx, ok := s.tokenIdx[bearerToken]  // ← map lookup (not constant time)
    ...
    if subtle.ConstantTimeCompare(
        []byte(stored), []byte(bearerToken),
    ) != 1 {
        return nil
    }
    ...
}
```

Go maps are hash-table backed with randomized hashing.
Map lookup is average O(1) but not strictly
constant-time per key. An attacker timing many
`Register`/RPC attempts could in principle extract
information about bucket collisions between their
guesses and valid tokens.

**Impact**: Largely theoretical. Go's hash randomization
makes reliable exploitation unlikely. Still, the
pattern "fast path via map, slow path via constant-time
compare" is not best practice.

**Recommendation**: For each incoming token:

1. Iterate all `ClientInfo` entries.
2. Do `subtle.ConstantTimeCompare` against every stored
   token.
3. OR the results together. Return the matching
   `ClientInfo` if exactly one matched.

This gives strict constant-time validation at the cost
of O(N) per RPC where N is the total number of clients.
For typical deployments (N < 1000), the cost is
negligible compared to gRPC overhead.

**Fix complexity**: Small.

**Existing task coverage**: None. Low enough to be a
stretch item within the token-hashing work (H-03).

---

#### H-21 — Token header parser tolerates missing "Bearer " prefix

**Severity**: Low.

**Location**: `internal/hub/validate.go:40`:

```go
token := strings.TrimPrefix(vals[0], bearerPrefix)
```

`strings.TrimPrefix` silently returns the original string
if the prefix is absent. So `authorization: ctx_cli_abc`
is treated the same as
`authorization: Bearer ctx_cli_abc`.

**Impact**: Minor protocol violation accepted.
Non-exploitable today, but it hides client bugs (a
client that forgets to prepend "Bearer " still works).

**Recommendation**: Require the exact `Bearer ` prefix
and reject otherwise with `codes.Unauthenticated`. Log
the event at INFO level for operator visibility.

**Fix complexity**: Trivial.

**Existing task coverage**: None. Low priority.

---

#### H-23 — Admin token stored plaintext on disk

**Severity**: Low (in the trusted-team model), High
(under a hostile-host model).

**Location**: `internal/cli/hub/core/server/setup.go:72-74`:

```go
if writeErr := io.SafeWriteFile(
    tokenPath, []byte(adminToken), fs.PermSecret,
); writeErr != nil {
```

The admin token is persisted to `<data-dir>/admin.token`
as plaintext (mode 0600).

**Impact**: Anyone with file read access to the hub data
directory gets the admin token. Root on the hub host is
already total compromise, so this is only a concern if
the data directory is exposed via other mechanisms
(shared volume, misconfigured backup, accidental
`chmod -R`). Under the declared trust model this is
acceptable.

**Recommendation**:

1. Prompt the operator to store the admin token in a
   secrets manager and delete the file after first run.
2. Alternatively, derive the admin token from a
   passphrase via argon2id so the file can store a
   hash only.
3. Document the file's sensitivity in
   `operations/hub.md` (partially done).

**Fix complexity**: Small.

**Existing task coverage**: Partially covered by the
H-03 clients.json hashing task (same pattern).

---

#### H-25 — Auth error messages distinguish missing vs invalid

**Severity**: Low (info leak).

**Location**: `internal/hub/validate.go:26-45`:

```go
if !ok { return ... "missing metadata" }
if len(vals) == 0 { return ... "missing token" }
if store.ValidateToken(token) == nil { return ... "invalid token" }
```

Three distinct error messages give an attacker
probing information: did they send metadata? Did they
send a token at all? Was the token in the right format?

**Impact**: Aids attacker reconnaissance. Low severity
because the attacker can't easily enumerate valid tokens
from the error messages, but best practice is to collapse
auth failures into a single generic code.

**Recommendation**: Return a single `codes.Unauthenticated`
with a fixed message ("authentication required") for all
three cases. Log the specific reason server-side only
(for operator debugging).

**Fix complexity**: Trivial.

**Existing task coverage**: None.

---

#### H-28 — Raft bind port predictable (`gRPC port + 1`)

**Severity**: Low.

**Location**: `internal/cli/hub/core/server/run.go:83`:

```go
bindAddr := fmt.Sprintf(":%d", port+1)
```

If a client knows the gRPC port (public via the
`Status` RPC), they can trivially compute the Raft
port. Combined with H-10 (unauthenticated Raft), this
makes the cluster a one-scan attack surface.

**Impact**: Reduces the cost of finding the Raft port
from "port scan" to "add 1." Reinforces H-10.

**Recommendation**: Accept a separate `--raft-port` or
`--raft-bind` flag instead of deriving from
`--port`. Default to a random high port or refuse to
start without an explicit value.

**Fix complexity**: Trivial.

**Existing task coverage**: None. Merge with H-10.

---

### Informational

#### H-26 — Daemon re-exec now uses the correct flag

**Severity**: Info (historical note, already fixed).

**Location**: `internal/cli/hub/core/server/daemon.go:56`
used to re-exec `ctx serve --hub`, which became
incorrect after the `ctx hub start` split. Fixed in this
session. Recorded here so the audit trail captures the
pre-fix state.

**Status**: Closed.

---

#### H-27 — mTLS / asymmetric auth not considered

**Severity**: Info.

**Discussion**: The current design uses bearer tokens
exclusively. An mTLS design would give cryptographic
client identity without shared secrets on the wire, at
the cost of a client-certificate management layer. This
is the "signed-claim / PKI" stretch task in the Hub
identity layer phase — called out here for completeness.

**Existing task coverage**: Hub identity layer phase
stretch task (PKI).

---

## Recommendations by timeline

### Do now (before any non-localhost deployment)

These are prerequisites for the Story 2 (trusted team)
deployment we're actively documenting:

- **H-01** + **H-02**: Add TLS support to both server
  and client. Without this, every LAN deployment is one
  sniffed packet away from token compromise.
- **H-04**: Server-enforce `Origin` on publish.
  Five-line fix, eliminates attribution forgery.
- **H-15**: Fix `appendFile` to actually append.
  Correctness plus availability.

### Short term (next sprint / next minor release)

These are the Story 2 hardening track:

- **H-03**: Hash `clients.json` tokens.
- **H-08** + **H-09** + **H-30**: Rate limiting, listener
  caps, backpressure — the DoS hardening bundle.
- **H-17**: Cap batch size on PublishRequest.
- **H-18**: Audit log (`audits.jsonl`).
- **H-19**: Revocation command.
- **H-22**: Decide `Entry.Author` fate.

### Medium term (next quarter / next major)

These unlock the sysadmin-registry MVP (Story 2 → Story 2.5):

- Full **Hub identity layer phase** from TASKS.md.
- **H-13**: Follower-side replication validation.
- **H-14**: Preserve master sequence on replication.
- **H-12**: Deterministic Raft bootstrap.
- **H-24**: Entry redaction.
- **H-29**: Bounded in-memory entry cache.

### Long term (Story 3 enablement)

These are prerequisites for any real public-internet
deployment:

- **H-10** + **H-11**: Authenticated and encrypted Raft
  transport (mTLS between peers).
- **PKI stretch task** from the Hub identity layer phase:
  signed short-lived claims replacing bearer tokens.
- Content sanitation (**H-16**) with cryptographically
  signed entries.
- Per-project ACLs on reads (currently out of scope even
  in the identity-layer phase).

## Defense-in-depth posture

Even after all findings are addressed, the hub should
treat the following as "always-on" hardening:

- **Run as an unprivileged user** via systemd with
  `ProtectSystem=strict`, `NoNewPrivileges=true`,
  `PrivateTmp=true`, `ReadWritePaths=/var/lib/ctx-hub`.
- **Dedicated filesystem or quota** for the data
  directory so a runaway publish fills one volume, not
  `/`.
- **Off-host backups** of `entries.jsonl` on an
  independent schedule from the hub binary lifecycle.
- **Monitoring**: `ctx hub status --exit-code` as a
  liveness probe; `entries.jsonl` growth rate alert;
  Raft leader flap alert.
- **Patch cadence**: Go runtime updates, gRPC updates,
  hashicorp/raft updates tracked against CVE feeds.

## Appendix: files reviewed

| File                                              | Findings drawn from |
|---------------------------------------------------|---------------------|
| `internal/hub/server.go`                          | H-01, H-30          |
| `internal/hub/client.go`                          | H-02                |
| `internal/hub/sync_helper.go`                     | H-02, H-13, H-14    |
| `internal/hub/failover.go`                        | H-02                |
| `internal/hub/handler.go`                         | H-04, H-05, H-22    |
| `internal/hub/validate.go`                        | H-04 (root), H-20, H-21, H-25 |
| `internal/hub/entry_validate.go`                  | H-04, H-16, H-17    |
| `internal/hub/store.go`                           | H-03, H-20, H-29    |
| `internal/hub/types.go`                           | H-03, H-06, H-07, H-22, H-29 |
| `internal/hub/persist.go`                         | H-15                |
| `internal/hub/fanout.go`                          | H-09                |
| `internal/hub/cluster.go`                         | H-10, H-11, H-12, H-28 |
| `internal/hub/fsm.go`                             | (no-op, confirmed)  |
| `internal/hub/auth.go`, `token.go`                | H-23 (primitives sound) |
| `internal/hub/grpc.go`                            | H-08, H-09, H-30    |
| `internal/cli/connect/core/config/config.go`      | (sound)             |
| `internal/crypto/crypto.go`                       | (sound)             |

## Appendix: threat-model dimensions NOT covered

For completeness, the following threat-model dimensions
are out of scope for this audit and should be tracked
separately:

- **Supply chain**: Go module pinning, dependency CVE
  monitoring, reproducible builds.
- **Build integrity**: signed binaries, release checksums,
  transparency log.
- **Operational runbooks**: incident response, tabletop
  exercises, key rotation drills.
- **Third-party library CVEs**: `hashicorp/raft`,
  `hashicorp/raft-boltdb`, `google.golang.org/grpc` all
  have their own attack surfaces we did not analyze.
- **Client-side workstation hardening**: disk encryption,
  screen lock, malware prevention.
- **AI-agent misbehavior**: a local agent writing harmful
  `ctx add --share` commands as a user mistake. Covered
  by `ctx`'s per-command confirmation UX (task in the
  TASKS.md secret-leak runbook item).

---

**End of audit.** Findings total: **30** (5 Critical,
12 High, 7 Medium, 4 Low, 2 Info). Critical and
High-severity findings cluster around transport
security, identity, and attribution — exactly the axes
that separate a "trusted LAN deployment" from a "real
public service," and exactly the axes the Hub identity
layer phase begins to address.
