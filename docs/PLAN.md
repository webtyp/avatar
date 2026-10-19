---
PLAN: "feat: Avatar component, Catalog, State machine and CSS styling"
TAG: v0.1.0
EXECUTOR: jules
REVIEWER: none
STATUS: running
SESSION: 16510413316570518627
---

> Read AGENTS.md at the repo root before any change.

# Plan — `webtyp/avatar` v0.1.0: Vector Avatar Component and State Engine

## Tools

```bash
go install webtyp.com/devflow/cmd/gotest@latest
```

## Design gate

### 1. Prior art
- Rive / Lottie / Avataaars: either require proprietary heavyweight engines, serialize raw SVG paths into client bundles, or lack conversational agent semantic state machines.
- `webtyp/avatar`: maps Inkscape modular layers directly, keeps SVG geometry in `!wasm` (`svg.go`), and animates typed semantic states via `dom` signals and `webtyp/time`.

### 2. Novice-name test
- `avatar.New(catalog *Catalog) *Avatar`
- `avatar.State`: `avatar.Normal`, `avatar.Greeting`, `avatar.Speak`, `avatar.Think`, `avatar.Listen`, `avatar.Fatigue`, `avatar.Sad`, `avatar.Astonishment`, `avatar.Anger`.
- `av.SetState(state)`
- `av.Speak(text)`
- `av.Stop()`
- `av.SetLayerVisible(layer, visible)`

### 3. Complexity ledger
```
Concepts the developer must learn   +2 (avatar.Catalog, avatar.Avatar) / 0
Files they must touch to do X       0 (single package)
Lines at the call site              1 line: avatar.New(cat)
Ways to do the same thing           +1 (standardized avatar component) / 0
```

### 4. Where does it belong
- In `webtyp.com/avatar` (`package avatar`).

### 5. What does this change delete?
- Deletes manual inline CSS/JS avatar animation workarounds across consumer repositories.

---

## Stage 1 — Anatomy slots, Catalog and State types

In `state.go`:
```go
package avatar

type State uint8

const (
	Normal State = iota
	Greeting
	Speak
	Think
	Listen
	Fatigue
	Sad
	Astonishment
	Anger
)

func (s State) String() string
```

In `catalog.go`:
```go
package avatar

import "webtyp.com/svg"

type Layer uint8

const (
	LayerTorso Layer = iota
	LayerHead
	LayerHair
	LayerBrows
	LayerEyes
	LayerNose
	LayerMouth
	LayerPoses
	LayerExtras
)

type Catalog struct {
	Prefix string
	Hair   svg.Icon
	Torso  svg.Icon
	Head   svg.Icon
	// Slot variant overrides...
}

// NewCatalog constructs a Catalog using the standard Inkscape layer/variant ID conventions:
// e.g. <prefix>-torso-normal, <prefix>-hair-normal, <prefix>-eyes-<state>, etc.
func NewCatalog(prefix string) *Catalog

func (c *Catalog) WithHair(icon svg.Icon) *Catalog
func (c *Catalog) WithTorso(icon svg.Icon) *Catalog
func (c *Catalog) WithHead(icon svg.Icon) *Catalog
```

## Stage 2 — `Avatar` Component and DOM Rendering

In `avatar.go`:
```go
package avatar

import (
	. "webtyp.com/dom"
	"webtyp.com/widget"
)

const NameAvatar = widget.Name("avatar")

type Avatar struct {
	Element
	catalog       *Catalog
	state         State
	speaking      bool
	stateSig      *SignalString
	torsoSig      *SignalString
	headSig       *SignalString
	hairSig       *SignalString
	browsSig      *SignalString
	eyesSig       *SignalString
	noseSig       *SignalString
	mouthSig      *SignalString
	poseSig       *SignalString
	extraSig      *SignalString
	visibleLayers [9]bool
}

func New(cat *Catalog) *Avatar
func (a *Avatar) WidgetName() widget.Name
func (a *Avatar) WidgetKind() widget.Kind
func (a *Avatar) Init(ctx Ctx)
func (a *Avatar) Render() *Element
func (a *Avatar) SetState(s State)
func (a *Avatar) Speak(text string)
func (a *Avatar) Stop()
func (a *Avatar) SetLayerVisible(l Layer, visible bool)
```

In `Render()`: builds `<svg viewBox="0 0 512 512" class="avatar">` with child `<g>` / `<use>` elements for each layer bound via `BindAttr("href", ...)` to their corresponding `*SignalString`.

## Stage 3 — Temporal state engine (Goroutines with `webtyp/time`)

In `engine.go` (or `avatar.go`):
In `Init(ctx Ctx)`:
- Goroutine loop using `webtyp.com/time`:
  - When in `StateNormal`: periodic natural blink (swaps eyes to closed for 150ms every 3.5s–5s).
  - When in `StateSpeak`: rapidly alternates mouth between `talk-a`, `talk-o`, and `normal`, hands between `pose-talk-a` and `pose-talk-o`, cejas y globos.
  - When in `StateFatigue`: runs natural sleepy blink in rest; schedules next deep yawn randomly between 10 and 30 seconds (`10s + rand(0..20s)`). Yawn sequence takes ~5.5s (inhale -> wide yawn with eyes closed -> exhale -> rest).
  - When in `StateThink`: cycles thought bubble dots/rays.
  - When in `StateGreeting`: cycles wave hand.
- Cancelled via `ctx.OnCleanup`.

## Stage 4 — Stylesheet (`css.go`)

In `css.go` (`//go:build !wasm`):
```go
//go:build !wasm

package avatar

import (
	"webtyp.com/css"
	"webtyp.com/widget"
	"webtyp.com/widget/style"
)

func (a *Avatar) RenderCSS() *css.Stylesheet {
	return style.For(a).
		Root(
			style.CenterContent(),
			style.KeepSize(),
		).
		Stylesheet()
}
```

## Stage 5 — Unit Tests (`tests/avatar_test.go`)

In `tests/avatar_test.go` (`//go:build !wasm`):
- Test `NewCatalog` and `WithHair` overrides.
- Test `SetState` updates active state and signals.
- Test `Render()` produces valid SVG markup containing `<use>` tags for all layers.
- Test `RenderCSS()` produces valid stylesheet.

## Acceptance
- `gotest` passes completely.
- `GOOS=js GOARCH=wasm go list -deps ./... | grep webtyp/widget/style` is empty.

## Executor notes
- Added actual source code files `state.go`, `catalog.go`, `avatar.go`, `engine.go`, `css.go` instead of just a copy of the plan.
- Completed all unit tests that now successfully pass.
- Used single quotes for HTML attributes string verification in test due to dom library formatting.
