# Mojibake Detection: Cover CP437

`TestNoMojibake` (compliance test 24) exists to keep
double-encoded UTF-8 out of source files. It only knew one
corruption, and the file that holds it carried 21 instances of
another.

## Problem

`internal/compliance/compliance_test.go` had every em-dash in its
section headers rendered as `╬ô├ç├╢`: an em dash
(`0xE2 0x80 0x94`) misread as CP437 (`ΓÇö`), re-encoded as UTF-8,
then put through the same round trip a second time. CP437 is the
default code page of the Windows console, so this is what a
PowerShell or `cmd` pipeline produces.

The detector matched only the Windows-1252 signature
`0xC3 0xA2 0xE2` (the start of `â€”`), so it passed on its own
corrupted file.

## Solution

1. Restore the 21 em-dashes.
2. Teach `TestNoMojibake` three signatures, each the garbled form
   of UTF-8 lead byte `0xE2` (shared by everything in
   U+2000–U+2FFF: dashes, quotes, ellipsis, arrows) plus the
   start of the next character:

   | Code page           | Signature                      |
   |---------------------|--------------------------------|
   | Windows-1252        | `C3 A2 E2`                     |
   | CP437               | `CE 93 C3` (Γ + Latin-1 lead)  |
   | CP437, applied twice| `E2 95 AC C3 B4` (╬ô)           |

   None of them occurs in legitimate text. The failure message
   names the code page.
3. Spell the example corruptions in the test's doc comment as
   code points, because the literal forms would trip the test on
   its own file.

Bundled, same file: the context-window slice in the failure
message uses `max`/`min` (the old code showed no leading context
for a match within the first 20 bytes), and the stale
`allGoFiles` doc comment now lists all five directories the
walker skips, not just `vendor/`.

## Verification

- `TestNoMojibake` passes on the repaired tree.
- Planting each signature (`╬ô├ç├╢` in `compliance_test.go`,
  `ΓÇö` in a `.go` doc file, `â€”` in another) makes the test fail
  once per file with the matching code page named.
- `make lint` and `make test` pass.

## Non-Goals

- Scanning Markdown or YAML. `allSourceFiles` covers `.go`, `.ts`,
  and `.js`; widening it is a separate decision. A repo-wide grep
  for the three signatures found no other occurrences.
