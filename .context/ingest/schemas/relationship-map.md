# Relationship Map Schema

Shape for `.context/kb/relationship-map.md`. Cross-topic and
cross-source relationships recorded as rows keyed on topic
slug or `EV-###` ID. The relationship map is the kb's edge
list; nodes are topic pages and evidence rows.

Relationships key on *slug + EV-### ID*, never on file path,
so folder reorganisations and lazy sub-page splits do not
invalidate the graph.

## Fields

| Column    | Description                                                              |
|-----------|--------------------------------------------------------------------------|
| `from`    | Originating topic slug or `EV-###` ID.                                   |
| `to`      | Destination topic slug or `EV-###` ID.                                   |
| `kind`    | Relationship kind (`depends-on`, `refines`, `contradicts`, `supersedes`).|
| `summary` | One-line description of the relationship.                                |

## Example

| from        | to        | kind        | summary                                                  |
|-------------|-----------|-------------|----------------------------------------------------------|
| widgets     | EV-042    | depends-on  | Widget topic page leans on EV-042 for its scope claim.   |

## Rules

- Key on slug and `EV-###`, never on file path; lazy
  sub-page splits and folder moves must not break edges.
- `kind` is drawn from a controlled vocabulary; introducing
  a new kind requires a `domain-decisions.md` entry naming
  the rationale.
- `contradicts` edges pair with a row in
  `contradictions.md`; the relationship map is the index,
  the contradictions file is the body.
- `supersedes` edges pair with the superseded row's demotion
  in `evidence-index.md`; never delete the superseded row.
- Self-edges (a slug pointing to itself, or an `EV-###`
  pointing to itself) are not meaningful and are rejected
  on write.
