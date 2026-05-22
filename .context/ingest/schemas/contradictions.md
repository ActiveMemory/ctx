# Contradictions Schema

Shape for `.context/kb/contradictions.md`. Each entry records
two or more `EV-###` rows that disagree, the demotion applied
under the demotion policy, and a resolution status. Resolved
contradictions stay in the file as audit trail.

## Fields

| Field               | Description                                                       |
|---------------------|-------------------------------------------------------------------|
| `id`                | Stable `C-###` identifier; zero-padded; never renumbered.         |
| `opened`            | ISO date the contradiction was filed.                             |
| `evidence`          | Two or more `EV-###` IDs that disagree.                           |
| `summary`           | One-line statement of what the rows disagree about.               |
| `demotion-applied`  | Which claim was demoted to which band, with reason.               |
| `resolution-status` | `open` or `resolved`; resolved entries name the winning claim.    |

## Example

```markdown
## C-007: Widget Scope (Per-Folder vs Per-File)

- **Opened:** 2026-05-16.
- **Evidence:** EV-042, EV-051.
- **Summary.** EV-042 says a widget is a folder; EV-051 says a
  widget is a single `.md` file with frontmatter.
- **Demotion applied.** EV-051 demoted from `medium` to `low`;
  EV-042 corroborated by EV-043 holds at `high`.
- **Resolution status:** resolved; folder-shaped (2026-05-16).
```

## Rules

- Every contradiction applies the demotion policy at filing
  time; recording the disagreement without demoting an older
  claim is a hard anti-pattern.
- Cite at least two `EV-###` IDs; a contradiction is between
  evidence rows, not between prose paragraphs.
- Never renumber or delete `EV-###` rows when reconciling;
  demote in place and let the row stand as audit trail.
- Resolved entries are kept; deletion erases kb history.
- Open contradictions surface in `/ctx-kb-site-review` and
  drive a paired `outstanding-questions.md` row naming both
  sides and what evidence would resolve.
