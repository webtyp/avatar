---
PLAN: "feat!: avatar v0.2 — independent activity and emotion, CSS-driven loops, working catalog slots"
TAG: v0.2.0
EXECUTOR: jules
REVIEWER: none
STATUS: review
SESSION: 3958202475294254305
PR: https://github.com/webtyp/avatar/pull/2
REVIEWED_COMMIT: 95d1bc0c3824790a7aabd9aa758dbdd0daa1688c
---

> Read AGENTS.md at the repo root before any change.

# Plan — `avatar` v0.2.0: activity × emotion, setters that apply at once, motion in CSS

v0.1.0 has one `State` enum mixing what the avatar *does* with how it *feels*, applies a state only
when its pending timer fires (measured: up to 20 s), never blinks while speaking, ignores
`WithHair/WithTorso/WithHead`, and allocates a string per frame. The record and the measurement are
in [`docs/ARCHITECTURE.md`](docs/ARCHITECTURE.md). This plan replaces the engine.

## Tools

```bash
go install webtyp.com/devflow/cmd/gotest@latest
```

## Ecosystem search

Reuses `webtyp/dom` (elements, `SignalString`, `BindAttr`, `Ctx`), `webtyp/widget` + `widget/style`
(parts, `Wave`, and the new `Blink`/`Breathe`/`Frame`), `webtyp/svg` (`Icon`), `webtyp/time`
(`AfterFunc`). Nothing new is created outside this repo.

## Design gate

1. **Prior art.** Unity Mecanim animation layers with avatar masks and ARKit (visemes apart from
   emotion blendshapes): one owner per body region. TalkingHead (met4citizen): a *mood* is a baseline
   speech plays over. Inworld: emotion arrives as data next to the text; the client draws it.
   Adobe Character Animator: blink is its own layer, independent of mouth and expression.
2. **Novice names.** `av.SetActivity(avatar.Think)`, `av.SetEmotion(avatar.Sad)`,
   `av.Speak("…")`, `av.Stop()`, `cat.WithGlasses(icon)`.
3. **Ledger.** Concepts +2 (`Activity`, `Emotion`) / −1 (`State`). Ways to set what the avatar
   shows: 1 → 1 (`SetState` is deleted). Call site unchanged in size.
4. **Where.** One repo: the user keeps the name `avatar`; the contract types and the SVG
   implementation stay together until a second implementation exists. The reference character
   (Cote's drawing, made public by the owner) lives in the sub-package `avatar/reference`: the
   laboratory (`agenteval`), this repo's tests and `mjosefa-cote` import one copy.
5. **Deletes.** `State`, `SetState`, `state.go`, the timer chain in `engine.go`, `go engine(...)`,
   the exported `Catalog` fields.

## APIs used

| Symbol | Signature | Tag | File |
|---|---|---|---|
| `dom.NewElement` | `func NewElement(tag string) *Element` | dom v0.13.26 | `element.go:50` |
| `(*Element).Attr` / `Class` / `Child` | `Attr(key, val string) *Element`, `Class(class ...string) *Element`, `Child(c ...Component) *Element` | dom v0.13.26 | `element.go:108-248` |
| `(*Element).BindAttr` | `BindAttr(name string, s *SignalString) *Element` | dom v0.13.26 | `element.go:301` |
| `(*Element).String` | `String() string` — serializes the CURRENT value of every bound attribute | dom v0.13.26 | `element.go:423` |
| `dom.NewString`, `(*SignalString).Set` | `NewString(v string) *SignalString`; `Set(v string)` is a no-op when `v` equals the current value | dom v0.13.26 | `signal.go:59-77` |
| `dom.Ctx` | `interface { OnCleanup(fn func()); … }` | dom v0.13.26 | `interface.dom.go:83` |
| `widget.Name.Class` | `func (n Name) Class(p Part) Class` → `"avatar__eyes"`; `Class.String()` | widget | `widget.go:31` |
| `style.For` / `Root` / `Part` | `For(w widget.Widget) *Sheet`; `Root(opts ...Option) *Sheet`; `Part(p widget.Part, opts ...Option) *Sheet`; `Stylesheet() *css.Stylesheet` | widget | `style/sheet.go:266`, `sheet_builder.go` |
| `style.Wave`, `Blink`, `Breathe` | `func() Option` | widget (Stage 1 tag) | `style/except_visual.go` |
| `style.Frame` | `func Frame(index, count int, c Cadence) Option`; `CadenceFast` (150ms/frame), `CadenceSlow` (400ms/frame) | widget (Stage 1 tag) | `style/except_visual.go`, `style/scale.go` |
| `svg.Icon` | `type Icon string` | svg v0.3.30 | `icon.go` |
| `time.AfterFunc` | `func AfterFunc(milliseconds int, f func()) Timer`; `Timer` is `interface{ Stop() bool }` | time v0.6.2 | `api.go:29,131` |

If a signature in the published tag differs from this table, stop and name the symbol in the PR;
never adapt it by guessing.

## The public API (exactly this; nothing else exported)

```go
// activity.go
type Activity uint8
const ( Idle Activity = iota; Greeting; Listen; Think; Speak; Fatigue )

// emotion.go
type Emotion uint8
const ( Neutral Emotion = iota; Happy; Sad; Astonished; Angry )

// catalog.go — value type, every With* returns a copy (deriving a character never mutates the base)
type Catalog struct { /* unexported */ }
func NewCatalog(prefix string) Catalog
func (c Catalog) WithTorso(icon svg.Icon) Catalog
func (c Catalog) WithHead(icon svg.Icon) Catalog
func (c Catalog) WithHair(icon svg.Icon) Catalog
func (c Catalog) WithShirt(icon svg.Icon) Catalog     // empty by default: not drawn
func (c Catalog) WithEarrings(icon svg.Icon) Catalog  // empty by default: not drawn
func (c Catalog) WithGlasses(icon svg.Icon) Catalog   // empty by default: not drawn

// layer.go — z-order, bottom to top
type Layer uint8
const ( LayerTorso Layer = iota; LayerShirt; LayerHead; LayerHair; LayerEarrings; LayerBrows;
	LayerEyes; LayerNose; LayerGlasses; LayerMouth; LayerPoses; LayerExtras; layerCount )

// avatar.go
const NameAvatar = widget.Name("avatar")
func New(cat Catalog) *Avatar
func (a *Avatar) SetActivity(act Activity)   // applied before it returns
func (a *Avatar) SetEmotion(e Emotion)       // applied before it returns; kept across activities
func (a *Avatar) Speak(text string)          // SetActivity(Speak), then Idle after speechMillis(text)
func (a *Avatar) Stop()                      // SetActivity(Idle)
func (a *Avatar) SetLayerVisible(l Layer, visible bool)
func (a *Avatar) Init(ctx dom.Ctx)           // registers the cleanup that stops pending timers
func (a *Avatar) Render() *dom.Element
func (a *Avatar) WidgetName() widget.Name; func (a *Avatar) WidgetKind() widget.Kind // unchanged
```

`NewCatalog(prefix)` defaults: torso `<prefix>-torso-normal`, head `<prefix>-head-normal`, hair
`<prefix>-hair-normal`. **Every href string is built once in `NewCatalog`/`With*`** (`"#" + icon`),
including a table per face layer, so no setter and no render concatenates strings.

## Variants — the rules (one owner per layer)

Face variant names (`<prefix>-<layer>-<name>`, layers `brows`, `eyes`, `mouth`): `normal`,
`greeting`, `sad`, `astonishment`, `anger`, `fatigue`, `listen`, `think`, `talk-a`, `talk-o`,
`closed`. Keep them as an unexported `variant uint8` enum + `[variantCount]string` names table; the
catalog precomputes `face [3][variantCount]string` hrefs.

Emotion → its own variant: Neutral→`normal`, Happy→`greeting`, Sad→`sad`,
Astonished→`astonishment`, Angry→`anger`.
Activity → its own variant: Idle→`normal`, Greeting→`greeting`, Listen→`listen`, Think→`think`,
Speak→`normal`, Fatigue→`fatigue`.

| Layer | Rule |
|---|---|
| brows, eyes | emotion's variant; when emotion is `Neutral`, the activity's variant |
| mouth (static `<use>`) | Greeting/Listen/Think/Fatigue → activity's variant; Idle → emotion's variant; Speak → static mouth hidden, the speech frames shown (below) |
| speech frames | three `<use>`: `mouth-talk-a`, `mouth-talk-o`, and a rest frame = **emotion's** mouth (so speaking angry rests on the angry mouth) |
| nose | `<prefix>-nose-sad` when emotion is Sad; `<prefix>-nose-closed` during yawn step 2; else empty href |
| poses | Greeting → `pose-greeting`; Think → `pose-think`; Speak → `pose-talk-a`/`pose-talk-o` frames; else none |
| extras | Speak → `extras-speech-a`/`extras-speech-o` frames; Think → `thought-bubble-1`, `-2`, `-3`, `-3-ring` static + `thought-dot-1..3` frames; Listen → `sound-inner`/`-mid`/`-outer` frames; Fatigue → `zzz-1..3` frames; Greeting → `extras-wave`; Idle → Sad `extras-drop`, Astonished `extras-marks`, else none |

Worked examples (what `av.Render().String()` contains right after the calls, prefix `t`):

| Calls | Contains |
|---|---|
| `New(NewCatalog("t"))` | `href="#t-brows-normal"`, `href="#t-eyes-normal"`, `href="#t-mouth-normal"` |
| `SetEmotion(Sad)` | `#t-brows-sad`, `#t-eyes-sad`, `#t-mouth-sad`, `#t-nose-sad`, `#t-extras-drop` |
| `SetEmotion(Angry); SetActivity(Speak)` | `#t-brows-anger`; speech group `display="inline"`; rest frame `#t-mouth-anger`; static mouth group `display="none"` |
| `SetActivity(Think)` (Neutral) | `#t-brows-think`, `#t-mouth-think`, think group `display="inline"` |
| `SetEmotion(Happy); SetActivity(Think)` | `#t-brows-greeting` (emotion wins), `#t-mouth-think` (activity wins) |
| `SetActivity(Fatigue); SetActivity(Greeting)` | fatigue group `display="none"`, `#t-mouth-greeting` — immediately |
| `NewCatalog("t").WithHair("u-hair-short")` | `#u-hair-short`; and `NewCatalog("t")` alone still `#t-hair-normal` |
| `SetLayerVisible(LayerHair, false)` | the hair layer group `display="none"` |

## Stage 1 — dependencies

`go get webtyp.com/widget@v0.6.44` and `go get webtyp.com/svg@v0.3.30`, then
`go mod tidy`.

## Stage 2 — types and catalog

Files: `activity.go`, `emotion.go` (with `String()` like today's `state.go`), `layer.go`,
`catalog.go` (rewrite). Delete `state.go` (`grep -rn "avatar.State\|SetState" .` → empty).

## Stage 3 — render tree (`avatar.go`)

`Render()` builds once:

```
<svg viewBox="0 0 512 512" class="avatar">
  <g class="avatar__body">                         ← Breathe
    one <g display=layerSig[l]> per layer torso…mouth, each with one <use href=…>
    eyes: <g class="avatar__eyes" …>              ← Blink (always on, never paused)
    mouth: <g display=mouthStaticSig><use href=mouthSig/></g>
           <g display=speakSig> <use class="avatar__mouth-a"/> <use class="avatar__mouth-o"/> <use class="avatar__mouth-rest" href=mouthRestSig/> </g>
  </g>
  poses:  <g class="avatar__greet-hand" display=greetSig>  ← Wave
          <g display=thinkSig><use pose-think/></g>
          <g class="avatar__talk-hands" display=speakSig> <use avatar__hand-a/> <use avatar__hand-o/> </g>  ← Wave on the group
  extras: <g display=extraStaticSig><use href=extraSig/></g>
          speech (2 frames), thought (4 static + 3 dot frames), sound (3 frames), zzz (3 frames), each group display bound to its activity signal
</svg>
```

Layers with an empty catalog href (shirt, earrings, glasses by default) are not rendered at all.
Display values are the constants `displayShown = "inline"` and `displayHidden = "none"`. Parts are
`widget.Part` constants (`PartBody`, `PartEyes`, `PartMouthA`, `PartMouthO`, `PartMouthRest`,
`PartGreetHand`, `PartTalkHands`, `PartHandA`, `PartHandO`, `PartSpeechA`, `PartSpeechO`,
`PartDot1..3`, `PartZzz1..3`, `PartSound1..3`); classes come from `NameAvatar.Class(part).String()`.

## Stage 4 — engine (`engine.go`, rewrite)

- `apply()` computes every signal from `(activity, emotion, yawnStep)` with the table above and
  calls `Set` on each. `SetActivity`/`SetEmotion` stop any pending timer, store the field, call
  `apply()`. No goroutine (`grep -rn "go " --include=*.go . | grep -v _test` shows no `go func`/`go engine`).
- `Speak(text)`: `SetActivity(Speak)`, then one `time.AfterFunc(speechMillis(text), a.Stop)`.
  `speechMillis` = 60 × number of runes (count with `for range text`), minimum 1200.
- Fatigue: on entering, `time.AfterFunc(10000 + rand.Intn(20001), a.yawn)`; the yawn steps
  1500 ms mouth `talk-o` → 2000 ms mouth `astonishment` + eyes `closed` + nose closed → 2000 ms
  mouth `talk-o` → back to the table, then schedule the next yawn. One `timer time.Timer` field
  holds whatever is pending; `Init` registers `ctx.OnCleanup` that stops it.
- `css.go`: `Root(CenterContent(), KeepSize())` as today, plus `Part(PartBody, Breathe())`,
  `Part(PartEyes, Blink())`, `Part(PartGreetHand, Wave())`, `Part(PartTalkHands, Wave())`, mouth
  `Frame(0..2, 3, CadenceFast)`, hands and speech `Frame(0..1, 2, CadenceFast)`, dots, zzz and
  sound `Frame(0..2, 3, CadenceSlow)`.

## Stage 5 — tests (`tests/avatar_test.go`, rewrite; `//go:build !wasm`)

One test per row of the worked-examples table: build, call, `out := av.Render().String()`,
`strings.Contains(out, want)`. Plus: `Speak("hola")` then `Stop()` → speech group `display="none"`;
`speechMillis` is unexported — test it only through `Speak` + `Stop`. `css.go`: `RenderCSS()`
does not panic and `Validate()` of the sheet returns no errors. Delete the v0.1 `TestAvatar_State`.

## Stage 6 — the reference character (`reference/`, new package `webtyp.com/avatar/reference`)

`reference/cote.svg` is already in the repo (the Inkscape file of Cote, 115 KB): do not edit it.
A sub-package, so the drawing reaches only the binaries that import it (`sitec` extracts the
`IconSvg()` of every package a build reaches; in the root package it would land in every
consumer's page).

```go
// reference/reference.go — untagged
package reference

import "webtyp.com/avatar"

const Prefix = "cote" // symbol prefix of cote.svg

// Catalog is the reference character: Cote, drawn following the Inkscape layer convention.
func Catalog() avatar.Catalog { return avatar.NewCatalog(Prefix) }
```

```go
// reference/svg.go
//go:build !wasm

package reference

import (
	_ "embed"
	"webtyp.com/svg/sprite"
)

//go:embed cote.svg
var coteSVG string

// IconSvg is the reference character's sprite, extracted by sitec at build time.
func IconSvg() *sprite.Sprite {
	s := sprite.NewSprite()
	if err := s.AddLayers(Prefix, coteSVG); err != nil {
		panic(err) // build time only: a broken reference drawing must stop the build
	}
	return s
}
```

`(*sprite.Sprite).AddLayers(prefix, svgContent string) error` (svg v0.3.30, `sprite/layers.go:16`;
the second argument is the file's **content**); `(*Sprite).String() string` renders the sprite.

Test `tests/reference_test.go` (`//go:build !wasm`): `out := reference.IconSvg().String()`; for every
symbol id the engine can reference (all face variants of the table above for `brows`, `eyes`,
`mouth`; `nose-sad`, `nose-closed`; every `pose-*` and `extras-*` id of the extras rule; the frame
pieces `thought-dot-1..3`, `thought-bubble-1..3`, `thought-bubble-3-ring`, `zzz-1..3`,
`sound-inner`, `sound-mid`, `sound-outer`, `extras-speech-a`, `extras-speech-o`) assert
`strings.Contains(out, `id="cote-<id>"`)`. Exception, verbatim in a comment: the drawing has no
`happy` variants; `Happy` uses `greeting`. This test is the guard that the engine and the drawing
agree.

## Stage 7 — docs

`README.md`: three sections, "Build a character" (catalog, slots, Inkscape ID convention above),
"Animate it" (`SetActivity`, `SetEmotion`, `Speak`), with the worked examples, and "Reference
character" (`reference.Catalog()` + importing the package so its sprite is extracted).
`docs/ARCHITECTURE.md`: delete items 1–7 of "Known defects in v0.1.0" (fixed here; keep item 8)
and add a "Channels" section with the variant rules table. `AGENTS.md` "This repo": replace the
`StateFatigue`/`StateSpeak` timing line with "Timers exist only for `Speak` end and the Fatigue
yawn; repeating motion is CSS (`widget/style`)".

## Delivery check

```bash
GOFLAGS=-mod=mod GOPRIVATE=webtyp.com/* go build ./...
gotest
grep -rn "SetState\|avatar.State\b" .     # empty
test ! -e state.go
GOOS=js GOARCH=wasm go list -deps ./... | grep webtyp/svg/sprite   # empty: sprite never reaches WASM
```

A PR that does not compile, or whose tests could not run, is a failed delivery. Never pause to ask
permission or report progress: when the check passes, commit, push and open the PR. Stuck for 30
minutes → open the PR anyway, listing what fails.

## Stages

| # | What | Files |
|---|---|---|
| 1 | dependencies | `go.mod`, `go.sum` |
| 2 | types, catalog | `activity.go`, `emotion.go`, `layer.go`, `catalog.go`; delete `state.go` |
| 3 | render tree | `avatar.go` |
| 4 | engine, styles | `engine.go`, `css.go` |
| 5 | tests | `tests/avatar_test.go` |
| 6 | reference character | `reference/reference.go`, `reference/svg.go` (new), `tests/reference_test.go` (new); `reference/cote.svg` exists |
| 7 | docs | `README.md`, `docs/ARCHITECTURE.md`, `AGENTS.md` |
