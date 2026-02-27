# Package Dependencies (Diagrams)

Parent: [ARCHITECTURE.md](ARCHITECTURE.md)

## ASCII Dependency Tree

```
                          ┌─────────────┐
                          │  cmd/ctx    │
                          │  main.go    │
                          └──────┬──────┘
                                 │
                          ┌──────▼──────┐
                          │  bootstrap  │
                          │  (root cmd) │
                          └──────┬──────┘
                                 │
           ┌─────────────────────┼─────────────────────┐
           │                     │                     │
    ┌──────▼──────┐       ┌──────▼──────┐       ┌──────▼──────┐
    │   cli/add   │       │  cli/agent  │  ...  │ cli/system  │
    │  cli/drift  │       │  cli/recall │       │  cli/watch  │
    └──────┬──────┘       └──────┬──────┘       └──────┬──────┘
           │                     │                     │
    ┌──────┴─────────────────────┴─────────────────────┘
    │  Shared dependencies (selected per command)
    │
    ├──► context ──► rc ──► config          (leaf)
    ├──► drift ──► context, index, rc
    ├──► index ──► config
    ├──► task ──► config
    ├──► validation ──► config
    ├──► recall/parser ──► config
    ├──► claude ──► assets
    ├──► notify ──► crypto, rc, config
    ├──► journal/state ──► config
    ├──► crypto                             (leaf, stdlib only)
    └──► sysinfo                            (leaf, stdlib only)
```

## Dependency Matrix

```
                    config  assets  rc  context  crypto  sysinfo  drift  index  task  validation  recall/parser  claude  notify  journal/state
config                 -
assets                 -      -
rc                     ✓             -
context                ✓             ✓     -
crypto                                           -
sysinfo                                                    -
drift                  ✓             ✓     ✓                        -      ✓
index                  ✓                                                   -
task                   ✓                                                          -
validation             ✓                                                                 -
recall/parser          ✓                                                                            -
claude                        ✓                                                                            -
notify                 ✓             ✓                  ✓                                                          -
journal/state          ✓                                                                                                  -
```

## Mermaid Graph

```mermaid
graph TD
    config["config<br/>(constants, regex, file names)"]
    assets["assets<br/>(embedded templates)"]

    rc["rc<br/>(runtime config)"] --> config
    context["context<br/>(loader)"] --> rc
    context --> config
    drift["drift<br/>(detector)"] --> config
    drift --> context
    drift --> index
    index["index<br/>(reindexing)"] --> config
    task["task<br/>(parsing)"] --> config
    validation["validation<br/>(sanitize)"] --> config
    recall_parser["recall/parser<br/>(session parsing)"] --> config
    claude["claude<br/>(hooks, skills)"] --> assets
    crypto["crypto<br/>(AES-256-GCM)"]
    sysinfo["sysinfo<br/>(OS metrics)"]
    notify["notify<br/>(webhooks)"] --> config
    notify --> crypto
    notify --> rc
    journal_state["journal/state<br/>(pipeline state)"] --> config

    bootstrap["bootstrap<br/>(CLI entry)"] --> rc
    bootstrap --> cli_all["cli/* (23 commands)"]
    cli_all --> config
    cli_all --> rc
    cli_all --> context
    cli_all --> drift
    cli_all --> index
    cli_all --> task
    cli_all --> validation
    cli_all --> recall_parser
    cli_all --> claude
    cli_all --> assets
    cli_all --> crypto
    cli_all --> sysinfo
    cli_all --> notify
    cli_all --> journal_state
```
