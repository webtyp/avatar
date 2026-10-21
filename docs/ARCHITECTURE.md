# avatar — architecture

STATUS (remove this line when AVATAR_MASTER_PLAN closes): v0.1.0 is being redesigned in this same
repository (activity/emotion/autonomous channels, CSS-driven loops, `cmd/avatarc`).

## Design of v0.1.0 (record)

| Decision | Why |
|---|---|
| The `.svg` stays 100 % Inkscape-pure: native layers, `display:none` for hidden variants, `fill="#ffffff"` on `*-bg` occlusion paths | The designer edits one file in Inkscape; the Go build cleans it |
| Modular anatomy: `torso`, `head`, `hair` never deform; `poses` holds only the hands (`*-bg` + `*-ink`) | Another character swaps a slot without touching eyes, brows, mouths or hands |
| Each variant becomes its own `<symbol id="<prefix>-<layer>-<variant>">` through `(*sprite.Sprite).AddLayers` (`webtyp/svg` v0.3.28) | `<use>` content is a closed shadow tree: CSS cannot show or hide nodes inside one symbol |
| Geometry only in `!wasm` (`svg.go`, `IconSvg()`); WASM carries only symbol ids | `sitec` inlines the sprite in `<body>`; the 115 KB of paths never reach the binary |

## Prior art

| Product | What it separates | What we take |
|---|---|---|
| [Adobe Character Animator](https://helpx.adobe.com/adobe-character-animator/desktop/creating-and-controlling-puppets/prepare-artwork.html) | A layered Illustrator/Photoshop file whose layer **names** auto-rig the puppet: `Left Blink`, mouth visemes `Ah, D, Ee, F, L, M, Oh, R, S, Uh, W-Oo`, plus `Neutral, Smile, Surprised` | Closest analogue: naming convention in a vector file → live rig. Blink is its own layer, independent of mouth and expression |
| [TalkingHead (met4citizen)](https://github.com/met4citizen/TalkingHead) | `mood` (a baseline of expression), lip-sync from the TTS, gestures and poses as separate calls; LLM demos expose `set_mood`, `make_hand_gesture`, `make_facial_expression` as tools | Mood is a baseline that speech plays over, never a state that replaces it |
| [Inworld AI](https://docs.inworld.ai/docs/tutorial-integrations/Unity/legacy/v2/api-reference/BodyAnimation) | The character engine sends **emotion** and **gesture** events next to the text; the client maps them to its own animations (`onEmotionChanged`, `HandleEmotion`, `HandleGesture`) | The agent emits data; the page decides how to draw it |
| Unity Mecanim layers + Avatar Mask; ARKit blendshapes | One layer per body region; visemes apart from emotion blendshapes | Each layer of the SVG has one owner channel |
| Rive, Lottie, Avataaars | Runtime-heavy (Rive), not editable in Inkscape (Lottie), geometry in the JS bundle (Avataaars) | Not adopted |

## Measured on v0.1.0 (2026-10-20, `mjosefa-cote`, TinyGo build, webtyp MCP)

Time from a state button click to the mouth `<use>` changing: Normal→Sad 1.7 s, Normal→Think
3.7 s, Fatigue→Greeting **20.3 s**, any other 505 ms. No long task was recorded: the page never
froze; the engine ignores the order until its pending timer fires.

## Known defects in v0.1.0

1. `SetState` and `Speak` only write a field; the engine reads it when its current timer fires
   (`engine.go`: 500 ms poll, 3.5–5 s in `Normal`, 10–30 s in `Fatigue`).
2. `State` `Speak` has no case in `engine.go`: its button does nothing. `Speak(text)` ignores the
   text and never ends by itself.
3. Blink exists only in `Normal`; while speaking the eyes cycle `talk-a`/`talk-o`, so there is
   no blink.
4. `Catalog.WithHair`/`WithTorso`/`WithHead` are never read: `Render` builds every href from
   `Prefix`. A second character cannot be built today.
5. `SetLayerVisible` is read only in `Render`: toggles do nothing live.
6. Allocations per frame: `Render` concatenates the href string on every signal change; every
   `time.AfterFunc` allocates a timer and a `js.FuncOf` (≈ 7 per second while speaking).
7. `stateSig` is never bound; `go engine(...)` starts a goroutine that only schedules a timeout.
8. Consumer: `mjosefa-cote/svg.go` discards the `AddLayers` error (`_ =`); `web/client.go`
   writes its CSS by hand; `docs/avatar/preview.html` was not deleted.
