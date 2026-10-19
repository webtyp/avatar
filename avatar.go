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
	// e.g. <prefix>-torso-normal
	prefix := a.catalog.Prefix

	return NewElement("svg").
		Attr("viewBox", "0 0 512 512").
		Attr("class", "avatar").
		Child(
			Show(NewBool(a.visibleLayers[LayerTorso]), func() *Element {
				return NewElement("g").Child(
					NewElement("use").BindAttrFunc("href", func() string { return "#" + prefix + "-torso-" + a.torsoSig.Get() }),
				)
			}),
			Show(NewBool(a.visibleLayers[LayerHead]), func() *Element {
				return NewElement("g").Child(
					NewElement("use").BindAttrFunc("href", func() string { return "#" + prefix + "-head-" + a.headSig.Get() }),
				)
			}),
			Show(NewBool(a.visibleLayers[LayerHair]), func() *Element {
				return NewElement("g").Child(
					NewElement("use").BindAttrFunc("href", func() string { return "#" + prefix + "-hair-" + a.hairSig.Get() }),
				)
			}),
			Show(NewBool(a.visibleLayers[LayerBrows]), func() *Element {
				return NewElement("g").Child(
					NewElement("use").BindAttrFunc("href", func() string { return "#" + prefix + "-brows-" + a.browsSig.Get() }),
				)
			}),
			Show(NewBool(a.visibleLayers[LayerEyes]), func() *Element {
				return NewElement("g").Child(
					NewElement("use").BindAttrFunc("href", func() string { return "#" + prefix + "-eyes-" + a.eyesSig.Get() }),
				)
			}),
			Show(NewBool(a.visibleLayers[LayerNose]), func() *Element {
				return NewElement("g").Child(
					NewElement("use").BindAttrFunc("href", func() string { return "#" + prefix + "-nose-" + a.noseSig.Get() }),
				)
			}),
			Show(NewBool(a.visibleLayers[LayerMouth]), func() *Element {
				return NewElement("g").Child(
					NewElement("use").BindAttrFunc("href", func() string { return "#" + prefix + "-mouth-" + a.mouthSig.Get() }),
				)
			}),
			Show(NewBool(a.visibleLayers[LayerPoses]), func() *Element {
				return NewElement("g").Child(
					NewElement("use").BindAttrFunc("href", func() string { return "#" + prefix + "-poses-" + a.poseSig.Get() }),
				)
			}),
			Show(NewBool(a.visibleLayers[LayerExtras]), func() *Element {
				return NewElement("g").Child(
					NewElement("use").BindAttrFunc("href", func() string { return "#" + prefix + "-extras-" + a.extraSig.Get() }),
				)
			}),
		)
}
