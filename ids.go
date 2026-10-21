package avatar

// Symbol names (the id is "<prefix>-<name>"). The engine uses only these.
const (
	symTorso = "torso-normal"
	symHead  = "head-normal"
	symHair  = "hair-normal"

	symNoseSad     = "nose-sad"
	symNoseClosed  = "nose-closed"
	symExtrasDrop  = "extras-drop"
	symExtrasMarks = "extras-marks"

	symMouthTalkA   = "mouth-talk-a"
	symMouthTalkO   = "mouth-talk-o"
	symPoseGreeting = "pose-greeting"
	symPoseThink    = "pose-think"
	symPoseTalkA    = "pose-talk-a"
	symPoseTalkO    = "pose-talk-o"
	symSpeechA      = "extras-speech-a"
	symSpeechO      = "extras-speech-o"
	symBubble1      = "thought-bubble-1"
	symBubble2      = "thought-bubble-2"
	symBubble3      = "thought-bubble-3"
	symBubble3Ring  = "thought-bubble-3-ring"
	symDot1         = "thought-dot-1"
	symDot2         = "thought-dot-2"
	symDot3         = "thought-dot-3"
	symSoundInner   = "sound-inner"
	symSoundMid     = "sound-mid"
	symSoundOuter   = "sound-outer"
	symZzz1         = "zzz-1"
	symZzz2         = "zzz-2"
	symZzz3         = "zzz-3"
)

// faceLayers are the three layers that have one symbol per variant ("<layer>-<variant>").
var faceLayers = [3]string{"brows", "eyes", "mouth"}

// RequiredSymbols returns every symbol id a character drawn with this prefix must define, in a
// stable order: torso, head, hair; then brows/eyes/mouth × every variant (variantNames order);
// then the constants above in declaration order. The optional slots (shirt, earrings, glasses)
// are not required.
func RequiredSymbols(prefix string) []string {
	var ids []string
	seen := make(map[string]bool)

	add := func(id string) {
		if !seen[id] {
			seen[id] = true
			ids = append(ids, id)
		}
	}

	p := prefix + "-"
	add(p + symTorso)
	add(p + symHead)
	add(p + symHair)

	for _, layerName := range faceLayers {
		for _, varName := range variantNames {
			add(p + layerName + "-" + varName)
		}
	}

	for _, sym := range []string{
		symNoseSad,
		symNoseClosed,
		symExtrasDrop,
		symExtrasMarks,
		symMouthTalkA,
		symMouthTalkO,
		symPoseGreeting,
		symPoseThink,
		symPoseTalkA,
		symPoseTalkO,
		symSpeechA,
		symSpeechO,
		symBubble1,
		symBubble2,
		symBubble3,
		symBubble3Ring,
		symDot1,
		symDot2,
		symDot3,
		symSoundInner,
		symSoundMid,
		symSoundOuter,
		symZzz1,
		symZzz2,
		symZzz3,
	} {
		add(p + sym)
	}

	return ids
}
