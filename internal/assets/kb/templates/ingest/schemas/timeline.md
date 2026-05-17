# Timeline Schema

Shape for `.context/kb/timeline.md`. Dated events worth
recording in the kb, oldest first. Each entry pins a real
occurrence in the domain to `EV-###` rows so the timeline
can be cross-checked against its sources.

The timeline is not a changelog of kb activity; that lives
in `SESSION_LOG.md` and the closeouts under
`.context/ingest/closeouts/`. The timeline records events
the kb is studying, not events the kb performed.

## Fields

| Field            | Description                                                    |
|------------------|----------------------------------------------------------------|
| `date`           | ISO date the event occurred (not the date it was ingested).    |
| `event`          | One-paragraph description of what happened.                    |
| `source-ev`      | Comma-separated `EV-###` references that ground the event.     |
| `related-topics` | Optional topic slugs touched by the event.                     |

## Example

```markdown
## 2026-05-16: Widget Contract v2 Published

- **Event.** The upstream maintainer published v2 of the
  widget contract, adding the `compatibility` frontmatter
  field and tightening the `description` length budget.
- **Source EV:** EV-042, EV-051.
- **Related topics:** widgets, widget-composition.
```

## Rules

- Pin the event to its `occurred:` date, not the date it was
  ingested; the `extracted:` field on the evidence row
  records the latter.
- Cite at least one `EV-###`; an entry without source
  backing is speculation, not a timeline event.
- Speculation about future events belongs in
  `outstanding-questions.md`, not here; filing a dated event
  for something the kb merely expects to happen is a hard
  anti-pattern.
- Newer events are appended below older ones; do not
  reorder once published.
- Cross-reference related topic slugs so the timeline can be
  scanned for blast radius when a single event touches
  multiple pages.
