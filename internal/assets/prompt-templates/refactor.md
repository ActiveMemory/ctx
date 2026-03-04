# Refactor

Refactor the specified code following these rules:

1. **Write or verify tests first** — confirm existing behavior is captured before changing structure.
2. **Preserve all existing behavior** — refactoring changes structure, not outcomes.
3. **Make one structural change at a time** — keep each step reviewable and revertible.
4. **Run tests after each step** — catch regressions immediately, not at the end.
5. **Check project conventions** — consult `.context/CONVENTIONS.md` to ensure the refactored code follows established patterns.

If a refactoring step would change observable behavior, stop and flag it as a separate task.
