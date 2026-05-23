# Outstanding Questions Schema

Shape for `.context/kb/outstanding-questions.md`. Each entry
is an open gap the kb has not yet resolved. Stable `Q-###`
IDs; newest first; never renumber.

Questions block the promotion of any claim that depends on
them: a topic page with an open `Q-###` cited in its
`## Open questions` section cannot ship at `confidence: high`.

## Fields

| Field                      | Description                                                |
|----------------------------|------------------------------------------------------------|
| `id`                       | Stable `Q-###` identifier; zero-padded; never renumbered.  |
| `opened`                   | ISO date the question was opened.                          |
| `question`                 | One-sentence open question, phrased as a question.         |
| `why-it-matters`           | Why answering this changes a topic page or a band.         |
| `what-evidence-would-resolve` | The shape of evidence that would close this entry.      |
| `related-ev`               | Optional comma-separated `EV-###` references.              |

## Example

```markdown
## Q-009: Does the Widget Contract Permit Nested Skills?

- **Opened:** 2026-05-16.
- **Question.** Does the widget folder shape allow a nested
  `SKILL.md` inside a subfolder, or is one skill per folder
  the hard limit?
- **Why it matters.** The `widget-composition` topic page
  cannot promote above `low` until this is answered; the page
  currently hedges with both readings.
- **What evidence would resolve.** A statement in the upstream
  spec, or a worked example from the vendor's own repo.
- **Related EV:** EV-042, EV-043.
```

## Rules

- Opening a question without stating what evidence would
  resolve it is a hard anti-pattern; the field is required.
- Phrase the question as an interrogative; declarative
  statements belong in `evidence-index.md` with a band.
- `low`-confidence topic-page content ships only when paired
  with a matching `outstanding-questions.md` entry; the entry
  is the license to ship the hedge.
- Closed questions are not deleted; they record what was once
  open and how it resolved.
- Do not reuse a retired `Q-###` ID for a new question; the
  ID space is monotonic.
