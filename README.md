# avatar
<img src="docs/img/badges.svg">

Interactive modular SVG avatar component and state engine.

## Build a character

Use `avatar.NewCatalog` to construct a catalog and swap character parts.
Every element is composed based on the Inkscape layer ID conventions:

```go
prefix := "cote" // the character prefix in the SVG ids
cat := avatar.NewCatalog(prefix).
    WithHair("cote-hair-short").
    WithTorso("cote-torso-formal").
    WithHead("cote-head-normal")
    // WithShirt, WithEarrings, WithGlasses

av := avatar.New(cat)
```

## Animate it

Avatars possess both an Activity and an Emotion. Modifying either applies changes immediately.
To reset to an Idle state, call `Stop()`.

```go
av.SetEmotion(avatar.Sad)
av.SetActivity(avatar.Think)
av.SetEmotion(avatar.Angry)
av.SetActivity(avatar.Speak) // shows speaking loop with rest on angry

av.Speak("Hello!") // Automatically stops returning to idle after text duration
av.Stop()
```

Examples of layer variants applied:
- `New(NewCatalog("t"))` -> `#t-brows-normal`, `#t-eyes-normal`, `#t-mouth-normal`
- `SetEmotion(Sad)` -> `#t-brows-sad`, `#t-eyes-sad`, `#t-mouth-sad`, `#t-nose-sad`, `#t-extras-drop`
- `SetEmotion(Angry); SetActivity(Speak)` -> `#t-brows-anger`, `#t-mouth-anger` + speech loop
- `SetActivity(Think)` -> `#t-brows-think`, `#t-mouth-think`, think loop

## Reference character

The reference character Cote lives in the `reference` subpackage.

```go
import "webtyp.com/avatar/reference"

cat := reference.Catalog()

// To extract the sprite during build time:
sprite := reference.IconSvg()
```
