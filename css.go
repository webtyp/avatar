//go:build !wasm

package avatar

import (
	"webtyp.com/widget/style"
)

func (a *Avatar) RenderCSS() *style.Sheet {
	return style.For(a).
		Root(
			style.CenterContent(),
			style.KeepSize(),
		).
		Part(PartBody, style.Breathe()).
		Part(PartEyes, style.Blink()).
		Part(PartGreetHand, style.Wave()).
		Part(PartTalkHands, style.Wave()).
		Part(PartMouthA, style.Frame(0, 3, style.CadenceFast)).
		Part(PartMouthO, style.Frame(1, 3, style.CadenceFast)).
		Part(PartMouthRest, style.Frame(2, 3, style.CadenceFast)).
		Part(PartHandA, style.Frame(0, 2, style.CadenceFast)).
		Part(PartHandO, style.Frame(1, 2, style.CadenceFast)).
		Part(PartSpeechA, style.Frame(0, 2, style.CadenceFast)).
		Part(PartSpeechO, style.Frame(1, 2, style.CadenceFast)).
		Part(PartDot1, style.Frame(0, 3, style.CadenceSlow)).
		Part(PartDot2, style.Frame(1, 3, style.CadenceSlow)).
		Part(PartDot3, style.Frame(2, 3, style.CadenceSlow)).
		Part(PartZzz1, style.Frame(0, 3, style.CadenceSlow)).
		Part(PartZzz2, style.Frame(1, 3, style.CadenceSlow)).
		Part(PartZzz3, style.Frame(2, 3, style.CadenceSlow)).
		Part(PartSound1, style.Frame(0, 3, style.CadenceSlow)).
		Part(PartSound2, style.Frame(1, 3, style.CadenceSlow)).
		Part(PartSound3, style.Frame(2, 3, style.CadenceSlow))
}
