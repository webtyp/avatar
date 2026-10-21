package avatar

import (
	"webtyp.com/dom"
	"webtyp.com/time"
	"webtyp.com/widget"
)

const NameAvatar = widget.Name("avatar")

const (
	displayShown  = "inline"
	displayHidden = "none"
)

type Avatar struct {
	dom.Element

	catalog  Catalog
	activity Activity
	emotion  Emotion

	timer    time.Timer
	yawnStep int

	visibleLayers [layerCount]bool
	layerSig      [layerCount]*dom.SignalString

	// Specific part signals
	browsSig       *dom.SignalString
	eyesSig        *dom.SignalString
	noseSig        *dom.SignalString
	mouthSig       *dom.SignalString
	mouthRestSig   *dom.SignalString
	mouthStaticSig *dom.SignalString
	speakSig       *dom.SignalString
	greetSig       *dom.SignalString
	thinkSig       *dom.SignalString
	extraSig       *dom.SignalString
	extraStaticSig *dom.SignalString
	fatigueSig     *dom.SignalString
	listenSig      *dom.SignalString
}

func New(cat Catalog) *Avatar {
	a := &Avatar{
		catalog:  cat,
		activity: Idle,
		emotion:  Neutral,

		browsSig:       dom.NewString(""),
		eyesSig:        dom.NewString(""),
		noseSig:        dom.NewString(""),
		mouthSig:       dom.NewString(""),
		mouthRestSig:   dom.NewString(""),
		mouthStaticSig: dom.NewString(""),
		speakSig:       dom.NewString(displayHidden),
		greetSig:       dom.NewString(displayHidden),
		thinkSig:       dom.NewString(displayHidden),
		extraSig:       dom.NewString(""),
		extraStaticSig: dom.NewString(displayHidden),
		fatigueSig:     dom.NewString(displayHidden),
		listenSig:      dom.NewString(displayHidden),
	}
	for i := range a.visibleLayers {
		a.visibleLayers[i] = true
		a.layerSig[i] = dom.NewString(displayShown)
	}

	a.apply()
	return a
}

func (a *Avatar) SetActivity(act Activity) {
	if a.timer != nil {
		a.timer.Stop()
		a.timer = nil
	}
	a.yawnStep = 0
	a.activity = act
	a.apply()
}

func (a *Avatar) SetEmotion(e Emotion) {
	a.emotion = e
	a.apply()
}

func (a *Avatar) Stop() {
	a.SetActivity(Idle)
}

func (a *Avatar) SetLayerVisible(l Layer, visible bool) {
	if l >= layerCount {
		return
	}
	a.visibleLayers[l] = visible
	if visible {
		a.layerSig[l].Set(displayShown)
	} else {
		a.layerSig[l].Set(displayHidden)
	}
}

func (a *Avatar) Init(ctx dom.Ctx) {
	ctx.OnCleanup(func() {
		if a.timer != nil {
			a.timer.Stop()
		}
	})
}

func (a *Avatar) WidgetName() widget.Name {
	return NameAvatar
}

func (a *Avatar) WidgetKind() widget.Kind {
	return widget.Region
}

func (a *Avatar) Render() *dom.Element {
	prefix := a.catalog.prefix

	renderLayer := func(l Layer, hrefFunc func() string, id string) *dom.Element {
		use := dom.NewElement("use").BindAttrFunc("href", hrefFunc)
		g := dom.NewElement("g").BindAttr("display", a.layerSig[l])
		if id != "" {
			g.Class(id)
		}
		return g.Child(use)
	}

	optLayer := func(l Layer, href string) *dom.Element {
		if href == "" {
			return dom.NewElement("g").Attr("display", "none") // don't render empty catalog layers
		}
		return renderLayer(l, func() string { return href }, "")
	}

	return dom.NewElement("svg").
		Attr("viewBox", "0 0 512 512").
		Class("avatar").
		Child(
			dom.NewElement("g").Class(NameAvatar.Class(PartBody).String()).Child(
				optLayer(LayerTorso, a.catalog.torso),
				optLayer(LayerShirt, a.catalog.shirt),
				optLayer(LayerHead, a.catalog.head),
				optLayer(LayerHair, a.catalog.hair),
				optLayer(LayerEarrings, a.catalog.earrings),

				renderLayer(LayerBrows, func() string { return a.browsSig.Get() }, ""),

				dom.NewElement("g").Class(NameAvatar.Class(PartEyes).String()).BindAttr("display", a.layerSig[LayerEyes]).Child(
					dom.NewElement("use").BindAttrFunc("href", func() string { return a.eyesSig.Get() }),
				),

				renderLayer(LayerNose, func() string { return a.noseSig.Get() }, ""),
				optLayer(LayerGlasses, a.catalog.glasses),

				dom.NewElement("g").BindAttr("display", a.layerSig[LayerMouth]).Child(
					dom.NewElement("g").BindAttr("display", a.mouthStaticSig).Child(
						dom.NewElement("use").BindAttrFunc("href", func() string { return a.mouthSig.Get() }),
					),
					dom.NewElement("g").BindAttr("display", a.speakSig).Child(
						dom.NewElement("use").Class(NameAvatar.Class(PartMouthA).String()).Attr("href", "#"+prefix+"-mouth-talk-a"),
						dom.NewElement("use").Class(NameAvatar.Class(PartMouthO).String()).Attr("href", "#"+prefix+"-mouth-talk-o"),
						dom.NewElement("use").Class(NameAvatar.Class(PartMouthRest).String()).BindAttrFunc("href", func() string { return a.mouthRestSig.Get() }),
					),
				),
			),

			dom.NewElement("g").BindAttr("display", a.layerSig[LayerPoses]).Child(
				dom.NewElement("g").Class(NameAvatar.Class(PartGreetHand).String()).BindAttr("display", a.greetSig).Child(
					dom.NewElement("use").Attr("href", "#"+prefix+"-pose-greeting"),
				),
				dom.NewElement("g").BindAttr("display", a.thinkSig).Child(
					dom.NewElement("use").Attr("href", "#"+prefix+"-pose-think"),
				),
				dom.NewElement("g").Class(NameAvatar.Class(PartTalkHands).String()).BindAttr("display", a.speakSig).Child(
					dom.NewElement("use").Class(NameAvatar.Class(PartHandA).String()).Attr("href", "#"+prefix+"-pose-talk-a"),
					dom.NewElement("use").Class(NameAvatar.Class(PartHandO).String()).Attr("href", "#"+prefix+"-pose-talk-o"),
				),
			),

			dom.NewElement("g").BindAttr("display", a.layerSig[LayerExtras]).Child(
				dom.NewElement("g").BindAttr("display", a.extraStaticSig).Child(
					dom.NewElement("use").BindAttrFunc("href", func() string { return a.extraSig.Get() }),
				),
				dom.NewElement("g").BindAttr("display", a.speakSig).Child(
					dom.NewElement("use").Class(NameAvatar.Class(PartSpeechA).String()).Attr("href", "#"+prefix+"-extras-speech-a"),
					dom.NewElement("use").Class(NameAvatar.Class(PartSpeechO).String()).Attr("href", "#"+prefix+"-extras-speech-o"),
				),
				dom.NewElement("g").BindAttr("display", a.thinkSig).Child(
					dom.NewElement("use").Attr("href", "#"+prefix+"-thought-bubble-1"),
					dom.NewElement("use").Attr("href", "#"+prefix+"-thought-bubble-2"),
					dom.NewElement("use").Attr("href", "#"+prefix+"-thought-bubble-3"),
					dom.NewElement("use").Attr("href", "#"+prefix+"-thought-bubble-3-ring"),
					dom.NewElement("use").Class(NameAvatar.Class(PartDot1).String()).Attr("href", "#"+prefix+"-thought-dot-1"),
					dom.NewElement("use").Class(NameAvatar.Class(PartDot2).String()).Attr("href", "#"+prefix+"-thought-dot-2"),
					dom.NewElement("use").Class(NameAvatar.Class(PartDot3).String()).Attr("href", "#"+prefix+"-thought-dot-3"),
				),
				dom.NewElement("g").BindAttr("display", a.listenSig).Child(
					dom.NewElement("use").Class(NameAvatar.Class(PartSound1).String()).Attr("href", "#"+prefix+"-sound-inner"),
					dom.NewElement("use").Class(NameAvatar.Class(PartSound2).String()).Attr("href", "#"+prefix+"-sound-mid"),
					dom.NewElement("use").Class(NameAvatar.Class(PartSound3).String()).Attr("href", "#"+prefix+"-sound-outer"),
				),
				dom.NewElement("g").BindAttr("display", a.fatigueSig).Child(
					dom.NewElement("use").Class(NameAvatar.Class(PartZzz1).String()).Attr("href", "#"+prefix+"-zzz-1"),
					dom.NewElement("use").Class(NameAvatar.Class(PartZzz2).String()).Attr("href", "#"+prefix+"-zzz-2"),
					dom.NewElement("use").Class(NameAvatar.Class(PartZzz3).String()).Attr("href", "#"+prefix+"-zzz-3"),
				),
			),
		)
}
