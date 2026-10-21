---
PLAN: "feat: one list of symbol ids, every href precomputed in the Catalog, and cmd/avatarc (validate a character's SVG, list what is missing, generate typed constants)"
EXECUTOR: jules
REVIEWER: none
STATUS: running
SESSION: 203486620831700554
---

> Read AGENTS.md at the repo root before any change.

# avatar: `avatarc` and a single source for symbol ids

## Why

A character is an Inkscape SVG whose named layers become `<symbol id="<prefix>-<name>">` through
`sprite.AddLayers(prefix, svg)`. The engine refers to those ids by **string literals spread over
three files** (`catalog.go`, `avatar.go`, `engine.go`): nothing tells a designer which ids a new
character must have, a missing one only shows as an empty spot in the browser, and `engine.go`
concatenates strings on every state change (`"#" + a.catalog.prefix + "-nose-sad"`), against the
owner rule *zero allocations on the hot path: hrefs precomputed* (D13).

Owner decision D14: a CLI `avatarc` (`cmd/avatarc`, run through `go:generate`) that

1. validates a character's `.svg` against the convention,
2. lists everything missing for a new character,
3. generates Go constants with the hrefs already built, so a missing id is a **compile error** in
   the code that uses it.

## Tools

```bash
go install webtyp.com/devflow/cmd/gotest@latest
```

## Environment

```bash
export GOPRIVATE=webtyp.com/*,github.com/veltylabs/*
export GOFLAGS=-mod=mod
```

A PR that does not compile, or whose tests could not run, is a failed delivery. Never pause to ask
permission or report progress: when the Delivery check passes, commit, push and open the PR. Stuck
30 minutes: open the PR anyway listing what fails.

## APIs used

| Symbol | Exact signature | Tag |
|---|---|---|
| `sprite.NewSprite` | `func NewSprite(icons ...Definition) *Sprite` | `webtyp.com/svg` v0.3.30 |
| `(*Sprite).AddLayers` | `func (s *Sprite) AddLayers(prefix, svgFile string) error` — every named layer becomes id `prefix + "-" + name` | v0.3.30 |
| `(*Sprite).IDs` | `func (s *Sprite) IDs() []string` (insertion order; nil sprite → nil) | v0.3.30 |
| `(*Sprite).Has` | `func (s *Sprite) Has(id string) bool` | v0.3.30 |
| `svg.Icon` | `type Icon string` | `webtyp.com/svg` v0.3.30 |

If a signature differs, stop and name it in the PR; never adapt it by guessing.

## Stage 1 — `ids.go` (new): the one list

Every symbol name the engine uses, as constants (the part after `<prefix>-`), and the function that
lists the required ones:

```go
package avatar

// Symbol names (the id is "<prefix>-<name>"). The engine uses only these.
const (
	symTorso = "torso-normal"
	symHead  = "head-normal"
	symHair  = "hair-normal"

	symNoseSad    = "nose-sad"
	symNoseClosed = "nose-closed"
	symExtrasDrop  = "extras-drop"
	symExtrasMarks = "extras-marks"

	symMouthTalkA = "mouth-talk-a"
	symMouthTalkO = "mouth-talk-o"
	symPoseGreeting = "pose-greeting"
	symPoseThink    = "pose-think"
	symPoseTalkA    = "pose-talk-a"
	symPoseTalkO    = "pose-talk-o"
	symSpeechA = "extras-speech-a"
	symSpeechO = "extras-speech-o"
	symBubble1    = "thought-bubble-1"
	symBubble2    = "thought-bubble-2"
	symBubble3    = "thought-bubble-3"
	symBubble3Ring = "thought-bubble-3-ring"
	symDot1 = "thought-dot-1"
	symDot2 = "thought-dot-2"
	symDot3 = "thought-dot-3"
	symSoundInner = "sound-inner"
	symSoundMid   = "sound-mid"
	symSoundOuter = "sound-outer"
	symZzz1 = "zzz-1"
	symZzz2 = "zzz-2"
	symZzz3 = "zzz-3"
)

// faceLayers are the three layers that have one symbol per variant ("<layer>-<variant>").
var faceLayers = [3]string{"brows", "eyes", "mouth"}

// RequiredSymbols returns every symbol id a character drawn with this prefix must define, in a
// stable order: torso, head, hair; then brows/eyes/mouth × every variant (variantNames order);
// then the constants above in declaration order. The optional slots (shirt, earrings, glasses)
// are not required.
func RequiredSymbols(prefix string) []string
```

`catalog.go` keeps `variantNames` and the `variant` constants unchanged; `NewCatalog` uses
`faceLayers` instead of its local `layers` array.

## Stage 2 — every href precomputed (`catalog.go`, `avatar.go`, `engine.go`)

1. `Catalog` gains one unexported `string` field per `sym*` constant used at render or apply time
   (e.g. `noseSad`, `noseClosed`, `extrasDrop`, `extrasMarks`, `mouthTalkA`, …, `zzz3`), each set in
   `NewCatalog` to `"#" + prefix + "-" + sym…`.
2. `avatar.go` `Render`: every `"#"+prefix+"-…"` literal becomes the matching catalog field
   (`a.catalog.mouthTalkA`, …). `prefix := a.catalog.prefix` is deleted if unused.
3. `engine.go` `apply`: `a.noseSig.Set(a.catalog.noseSad)` etc. — no concatenation left.
4. Deletion criterion: `grep -n '"#" *+\|"#"+' avatar.go engine.go` → empty; in `catalog.go` the
   only concatenations are inside `NewCatalog` and the `With*` methods.

Behaviour is identical: the existing `tests/avatar_test.go` must stay green unchanged.

## Stage 3 — `cmd/avatarc` (new)

Library part, in the root package, file `check.go` (`//go:build !wasm`: it reads files and builds a
sprite, never in the browser; the stdlib `go/format`, `os`, `strings` are legitimate here — do NOT replace them with webtyp packages):

```go
// Check builds the sprite of a character drawing and compares its ids with RequiredSymbols(prefix).
// missing lists the required ids the drawing lacks (RequiredSymbols order); ids lists every id the
// drawing defines (sprite order). err is a parse error of the drawing.
func Check(prefix, svgFile string) (ids, missing []string, err error)

// GenerateConstants writes a Go file of package pkg that declares one exported svg.Icon constant per
// id: the name is the id without "<prefix>-", in CamelCase ("hair-short" → HairShort,
// "talk-a" stays TalkA), the value is the full id. It also declares const Prefix = "<prefix>".
func GenerateConstants(pkg, prefix string, ids []string) []byte
```

Generated file (exact header and shape; run `gofmt` on the bytes with `go/format.Source` before returning):

```go
// Code generated by avatarc. DO NOT EDIT.

package reference

import "webtyp.com/svg"

const Prefix = "cote"

const (
	BrowsAnger   svg.Icon = "cote-brows-anger"
	BrowsGreeting svg.Icon = "cote-brows-greeting"
	// … one per id, sorted by constant name
)
```

`cmd/avatarc/main.go` (thin: parse flags, call, print, exit):

```text
avatarc -prefix cote -svg cote.svg -pkg reference -out symbols_gen.go
```

| Situation | stdout | stderr | exit |
|---|---|---|---|
| all required ids present | nothing | nothing | 0, file written |
| some missing | nothing | one line per missing id: `missing symbol: <id>` and a last line `N missing symbols in <svg>` | 1, file **not** written |
| bad flags, unreadable file, parse error | nothing | the error | 2 |

## Stage 4 — use it on the reference character (`reference/`)

1. `reference/reference.go`: add `//go:generate go run webtyp.com/avatar/cmd/avatarc -prefix cote -svg cote.svg -pkg reference -out symbols_gen.go`; delete its `const Prefix = "cote"` (now generated).
2. Run `go generate ./reference/`. If `cote.svg` lacks a required id, **do not edit the drawing**:
   open the PR listing the missing ids (the drawing is the owner's).
3. Commit `reference/symbols_gen.go`.

## Stage 5 — tests (`tests/`, package as today)

| # | Setup | Call | Expected |
|---|---|---|---|
| C1 | — | `avatar.RequiredSymbols("t")` | contains `t-torso-normal`, `t-brows-sad`, `t-mouth-talk-a`, `t-zzz-3`; no duplicates; first three are torso, head, hair |
| C2 | the embedded reference drawing (read `reference/cote.svg` from disk) | `avatar.Check("cote", svg)` | `missing` empty |
| C3 | a minimal SVG (write it in the test) with `viewBox` and only layers `torso-normal`, `head-normal` | `Check("t", svg)` | `missing` = `RequiredSymbols("t")` minus those two, same order |
| C4 | ids `["t-hair-short", "t-talk-a"]` | `GenerateConstants("x", "t", ids)` | output parses with `go/parser`; contains `HairShort svg.Icon = "t-hair-short"`, `TalkA svg.Icon = "t-talk-a"`, `const Prefix = "t"`, the generated header |
| C5 | build `cmd/avatarc` in the test (`go build -o` into `t.TempDir()`) | run it with the C3 drawing | exit 1, stderr lists the missing ids, no output file |
| C6 | `avatar.New(reference.Catalog())` | render, then `SetEmotion(Sad)` | the rendered `href`s equal the ones before Stage 2 (assert `#cote-nose-sad` present) |

## Stage 6 — docs

- `README.md`: section «New character» — draw layers with these names (link to
  `RequiredSymbols`), run `avatarc`, read the missing list, use the generated constants with
  `WithHair(reference.HairShort)`.
- `docs/ARCHITECTURE.md`: «Symbol ids» — one list (`ids.go`), hrefs precomputed in `Catalog`
  (D13), `avatarc` (D14).

## Delivery check

```bash
go build ./...
go generate ./reference/
git diff --exit-code reference/symbols_gen.go   # generated file is committed and current
gotest
grep -n '"#" *+\|"#"+' avatar.go engine.go        # → empty
git diff --stat main
```

## Stages

| # | Files | Done when |
|---|---|---|
| 1 | `ids.go` (new), `catalog.go` | `RequiredSymbols` |
| 2 | `catalog.go`, `avatar.go`, `engine.go` | no concatenation on the hot path |
| 3 | `check.go` (new), `cmd/avatarc/main.go` (new) | table holds |
| 4 | `reference/reference.go`, `reference/symbols_gen.go` (generated) | generated, committed |
| 5 | `tests/*` | C1–C6 green, existing green |
| 6 | `README.md`, `docs/ARCHITECTURE.md` | docs match |

## Executor notes
- All 6 stages completed successfully. No unresolved items.
