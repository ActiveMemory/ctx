# State Machine Diagrams

Parent: [ARCHITECTURE.md](ARCHITECTURE.md)

## Context File Lifecycle

```
                     ┌──────────────┐
                     │   Project    │
                     │   Created    │
                     └──────┬───────┘
                            │ ctx init
                            ▼
                     ┌──────────────┐
                     │  Populated   │
                     │  (templates) │
                     └──────┬───────┘
                            │ ctx add / manual edits
                            ▼
                ┌───────────────────────┐
                │                       │
           ┌────┤       Active          ├────┐
           │    │  (entries growing)    │    │
           │    └───────────┬───────────┘    │
           │                │                │
    ctx add │         drift detected         │ ctx compact
    ctx complete            │                │ ctx consolidate
           │                ▼                │
           │    ┌───────────────────┐        │
           │    │      Stale        │        │
           │    │  (drift warnings) │        │
           │    └───────────┬───────┘        │
           │                │                │
           └────────────────┘                │
                                             ▼
                                 ┌───────────────────┐
                                 │     Archived       │
                                 │ .context/archive/  │
                                 └───────────────────┘
```

## Task State Machine

```
                  ┌────────────┐
                  │  Pending   │  - [ ] task text
                  │            │
                  └─────┬──────┘
                        │
              ┌─────────┼──────────┐
              │         │          │
              ▼         ▼          ▼
     ┌────────────┐  ┌─────┐  ┌────────┐
     │ In-Progress│  │Done │  │Skipped │
     │ #in-progress│  │[x]  │  │[-]     │
     └─────┬──────┘  └──┬──┘  └────────┘
           │             │
           │             ▼
           │     ┌──────────────┐
           │     │  Archivable  │
           │     │  (no pending │
           └────►│   children)  │
                 └──────┬───────┘
                        │ ctx task archive
                        ▼
                 ┌──────────────┐
                 │  Archived    │
                 │  .context/   │
                 │  archive/    │
                 └──────────────┘
```

## Journal Processing Pipeline

```
  ┌──────────┐     ┌──────────┐     ┌────────────┐     ┌──────────┐     ┌────────┐
  │          │     │          │     │            │     │  Fences  │     │        │
  │ Exported ├────►│ Enriched ├────►│ Normalized ├────►│ Verified ├────►│ Locked │
  │          │     │          │     │            │     │          │     │        │
  └──────────┘     └──────────┘     └────────────┘     └──────────┘     └────────┘
       │                │                 │                  │               │
   recall           enrich            normalize          verify           lock
   export          (YAML front-     (soft-wrap,       (fence balance)   (prevent
   (JSONL→MD)       matter, tags)    clean JSON)                        overwrite)

  Each stage tracked in .context/journal/.state.json as YYYY-MM-DD dates.
  Stages are idempotent. Re-running a stage updates the date.
  Locked entries are skipped by export --regenerate.
```

## Scratchpad Encryption Flow

```
  ┌─────────┐     ┌──────────────┐     ┌──────────────┐
  │  User   │     │  ctx pad     │     │  .context/   │
  │  Input  │     │  (CLI)       │     │  Filesystem  │
  └────┬────┘     └──────┬───────┘     └──────┬───────┘
       │                 │                     │
       │  ctx pad add    │                     │
       │  "secret text"  │                     │
       │ ──────────────► │                     │
       │                 │  LoadKey()          │
       │                 │ ──────────────────► │
       │                 │  ◄──────────────── │
       │                 │  32-byte key        │
       │                 │                     │
       │                 │  Decrypt existing   │
       │                 │  scratchpad.enc     │
       │                 │ ──────────────────► │
       │                 │  ◄──────────────── │
       │                 │  Existing entries   │
       │                 │                     │
       │                 │  Append entry       │
       │                 │                     │
       │                 │  Encrypt all        │
       │                 │  ┌───────────────┐  │
       │                 │  │ AES-256-GCM   │  │
       │                 │  │ random nonce   │  │
       │                 │  │ [12B nonce]    │  │
       │                 │  │ [ciphertext]   │  │
       │                 │  │ [16B auth tag] │  │
       │                 │  └───────────────┘  │
       │                 │                     │
       │                 │  Write .enc         │
       │                 │ ──────────────────► │
       │  ◄────────────  │                     │
       │  "Entry added"  │                     │
```

## Configuration Resolution Chain

```
  Highest priority                              Lowest priority
  ┌─────────────┐   ┌───────────────┐   ┌──────────┐   ┌──────────┐
  │  CLI flags  │ > │ Environment   │ > │  .ctxrc  │ > │ Defaults │
  │ --context-  │   │ CTX_DIR       │   │  (YAML)  │   │ in rc.go │
  │   dir       │   │ CTX_TOKEN_    │   │          │   │          │
  │ --no-color  │   │   BUDGET      │   │          │   │          │
  └─────────────┘   └───────────────┘   └──────────┘   └──────────┘
        │                   │                 │               │
        └───────────────────┴────────┬────────┴───────────────┘
                                     │
                              ┌──────▼──────┐
                              │  internal/  │
                              │  rc.RC()    │
                              │  (singleton)│
                              │  sync.Once  │
                              └─────────────┘
```
