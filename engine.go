package avatar

import (
	"math/rand"
	"webtyp.com/time"
)

func (a *Avatar) Speak(text string) {
	a.SetActivity(Speak)

	runesCount := 0
	for range text {
		runesCount++
	}
	millis := runesCount * 60
	if millis < 1200 {
		millis = 1200
	}

	a.timer = time.AfterFunc(millis, func() {
		if a.activity == Speak {
			a.Stop()
		}
	})
}

func (a *Avatar) yawn() {
	if a.activity != Fatigue {
		return
	}

	// Yawn sequence
	// step 1: mouth talk-o
	a.yawnStep = 1
	a.apply()

	a.timer = time.AfterFunc(1500, func() {
		if a.activity != Fatigue {
			return
		}
		// step 2: mouth astonishment + eyes closed + nose closed
		a.yawnStep = 2
		a.apply()

		a.timer = time.AfterFunc(2000, func() {
			if a.activity != Fatigue {
				return
			}
			// step 3: mouth talk-o
			a.yawnStep = 3
			a.apply()

			a.timer = time.AfterFunc(2000, func() {
				if a.activity != Fatigue {
					return
				}
				// step 0: back to normal fatigue
				a.yawnStep = 0
				a.apply()
				a.scheduleYawn()
			})
		})
	})
}

func (a *Avatar) scheduleYawn() {
	a.timer = time.AfterFunc(10000+rand.Intn(20001), a.yawn)
}

func (a *Avatar) apply() {
	if a.catalog.prefix == "" {
		return
	}

	act := a.activity
	emo := a.emotion

	// If transitioning to Fatigue, we need to schedule the yawn
	if act == Fatigue && a.timer == nil {
		a.scheduleYawn()
	}

	emoVariant := variantNormal
	switch emo {
	case Happy:
		emoVariant = variantGreeting
	case Sad:
		emoVariant = variantSad
	case Astonished:
		emoVariant = variantAstonishment
	case Angry:
		emoVariant = variantAnger
	}

	actVariant := variantNormal
	switch act {
	case Greeting:
		actVariant = variantGreeting
	case Listen:
		actVariant = variantListen
	case Think:
		actVariant = variantThink
	case Speak:
		actVariant = variantNormal
	case Fatigue:
		actVariant = variantFatigue
	}

	// brows, eyes
	beVariant := emoVariant
	if emo == Neutral {
		beVariant = actVariant
	}
	if a.yawnStep == 2 {
		beVariant = variantClosed
	}
	a.browsSig.Set(a.catalog.face[0][beVariant])
	a.eyesSig.Set(a.catalog.face[1][beVariant])

	// mouth (static)
	mouthVariant := emoVariant
	switch act {
	case Greeting, Listen, Think, Fatigue:
		mouthVariant = actVariant
	}
	if a.yawnStep == 1 || a.yawnStep == 3 {
		mouthVariant = variantTalkO
	} else if a.yawnStep == 2 {
		mouthVariant = variantAstonishment
	}

	a.mouthSig.Set(a.catalog.face[2][mouthVariant])

	// speech frames
	a.mouthRestSig.Set(a.catalog.face[2][emoVariant]) // rest frame = emotion's mouth

	// nose
	if emo == Sad {
		a.noseSig.Set(a.catalog.noseSad)
	} else if a.yawnStep == 2 {
		a.noseSig.Set(a.catalog.noseClosed)
	} else {
		a.noseSig.Set("")
	}

	// extras static
	if act == Idle {
		if emo == Sad {
			a.extraSig.Set(a.catalog.extrasDrop)
		} else if emo == Astonished {
			a.extraSig.Set(a.catalog.extrasMarks)
		} else {
			a.extraSig.Set("")
		}
	} else {
		a.extraSig.Set("")
	}

	// groups display based on activity
	a.speakSig.Set(displayHidden)
	a.mouthStaticSig.Set(displayShown)
	a.greetSig.Set(displayHidden)
	a.thinkSig.Set(displayHidden)
	a.extraStaticSig.Set(displayHidden)
	a.fatigueSig.Set(displayHidden)
	a.listenSig.Set(displayHidden)

	switch act {
	case Speak:
		a.speakSig.Set(displayShown)
		a.mouthStaticSig.Set(displayHidden)
	case Greeting:
		a.greetSig.Set(displayShown)
	case Think:
		a.thinkSig.Set(displayShown)
	case Fatigue:
		a.fatigueSig.Set(displayShown)
	case Listen:
		a.listenSig.Set(displayShown)
	case Idle:
		if emo == Sad || emo == Astonished {
			a.extraStaticSig.Set(displayShown)
		}
	}
}
