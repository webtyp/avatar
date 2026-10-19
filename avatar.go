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

func New(cat *Catalog) *Avatar {
	return &Avatar{
		catalog:       cat,
		stateSig:      NewString(""),
		torsoSig:      NewString(""),
		headSig:       NewString(""),
		hairSig:       NewString(""),
		browsSig:      NewString(""),
		eyesSig:       NewString(""),
		noseSig:       NewString(""),
		mouthSig:      NewString(""),
		poseSig:       NewString(""),
		extraSig:      NewString(""),
		visibleLayers: [9]bool{true, true, true, true, true, true, true, true, true},
	}
}

func (a *Avatar) WidgetName() widget.Name {
	return NameAvatar
}

func (a *Avatar) WidgetKind() widget.Kind {
	return widget.Region
}

func (a *Avatar) SetState(s State) {
	a.state = s
	a.stateSig.Set(s.String())
}

func (a *Avatar) Speak(text string) {
	a.speaking = true
}

func (a *Avatar) Stop() {
	a.speaking = false
}

func (a *Avatar) SetLayerVisible(l Layer, visible bool) {
	a.visibleLayers[l] = visible
}

func (a *Avatar) Init(ctx Ctx) {
	// engine loop
	go engine(ctx, a)
}

func (a *Avatar) Render() *Element {
	prefix := a.catalog.Prefix

	renderLayer := func(l Layer, part string, sig *SignalString) *Element {
		use := NewElement("use").
			BindAttrFunc("href", func() string {
				val := sig.Get()
				if val == "" || val == "none" {
					return ""
				}
				if part == "pose" {
					return "#" + prefix + "-pose-" + val
				}
				return "#" + prefix + "-" + part + "-" + val
			})
		g := NewElement("g").Child(use)
		if !a.visibleLayers[l] {
			g.Attr("style", "display:none")
		}
		return g
	}

	return NewElement("svg").
		Attr("viewBox", "0 0 512 512").
		Attr("class", "avatar").
		Child(
			renderLayer(LayerTorso, "torso", a.torsoSig),
			renderLayer(LayerHead, "head", a.headSig),
			renderLayer(LayerHair, "hair", a.hairSig),
			renderLayer(LayerBrows, "brows", a.browsSig),
			renderLayer(LayerEyes, "eyes", a.eyesSig),
			renderLayer(LayerNose, "nose", a.noseSig),
			renderLayer(LayerMouth, "mouth", a.mouthSig),
			renderLayer(LayerPoses, "pose", a.poseSig),
			renderLayer(LayerExtras, "extras", a.extraSig),
		)
}
