package avatar

import (
	"math/rand"
	"webtyp.com/dom"
	"webtyp.com/time"
)

func engine(ctx dom.Ctx, a *Avatar) {
	// Initialize default signals
	a.torsoSig.Set("normal")
	a.headSig.Set("normal")
	a.hairSig.Set("normal")
	a.browsSig.Set("normal")
	a.eyesSig.Set("normal")
	a.noseSig.Set("none")
	a.mouthSig.Set("normal")
	a.poseSig.Set("none")
	a.extraSig.Set("none")

	var timer time.Timer
	running := true
	ctx.OnCleanup(func() {
		running = false
		if timer != nil {
			timer.Stop()
		}
	})

	var step func()
	step = func() {
		if !running {
			return
		}

		if a.speaking {
			a.mouthSig.Set("talk-a")
			a.poseSig.Set("talk-a")
			a.extraSig.Set("speech-a")
			a.browsSig.Set("talk-a")
			a.eyesSig.Set("talk-a")
			timer = time.AfterFunc(150, func() {
				if !running || !a.speaking {
					step()
					return
				}
				a.mouthSig.Set("talk-o")
				a.poseSig.Set("talk-o")
				a.extraSig.Set("speech-o")
				a.browsSig.Set("talk-o")
				a.eyesSig.Set("talk-o")
				timer = time.AfterFunc(150, func() {
					if !running || !a.speaking {
						step()
						return
					}
					a.mouthSig.Set("normal")
					a.poseSig.Set("none")
					a.extraSig.Set("none")
					a.browsSig.Set("normal")
					a.eyesSig.Set("normal")
					timer = time.AfterFunc(120, step)
				})
			})
			return
		}

		switch a.state {
		case Normal:
			a.browsSig.Set("normal")
			a.eyesSig.Set("normal")
			a.noseSig.Set("none")
			a.mouthSig.Set("normal")
			a.poseSig.Set("none")
			a.extraSig.Set("none")
			sleepTime := 3500 + rand.Intn(1500)
			timer = time.AfterFunc(sleepTime, func() {
				if !running || a.state != Normal || a.speaking {
					step()
					return
				}
				a.eyesSig.Set("closed")
				timer = time.AfterFunc(150, step)
			})

		case Greeting:
			a.browsSig.Set("greeting")
			a.eyesSig.Set("greeting")
			a.noseSig.Set("none")
			a.mouthSig.Set("greeting")
			a.poseSig.Set("greeting")
			a.extraSig.Set("wave")
			timer = time.AfterFunc(500, step)

		case Think:
			a.browsSig.Set("think")
			a.eyesSig.Set("think")
			a.noseSig.Set("sad")
			a.mouthSig.Set("think")
			a.poseSig.Set("think")
			a.extraSig.Set("thought")
			timer = time.AfterFunc(500, step)

		case Listen:
			a.browsSig.Set("listen")
			a.eyesSig.Set("listen")
			a.noseSig.Set("none")
			a.mouthSig.Set("listen")
			a.poseSig.Set("none")
			a.extraSig.Set("sound")
			timer = time.AfterFunc(500, step)

		case Fatigue:
			a.browsSig.Set("fatigue")
			a.eyesSig.Set("fatigue")
			a.noseSig.Set("none")
			a.mouthSig.Set("fatigue")
			a.poseSig.Set("none")
			a.extraSig.Set("zzz")
			sleepTime := 10000 + rand.Intn(20000)
			timer = time.AfterFunc(sleepTime, func() {
				if !running || a.state != Fatigue || a.speaking {
					step()
					return
				}
				// Deep yawn sequence (~5.5s)
				a.mouthSig.Set("talk-a")
				timer = time.AfterFunc(1500, func() {
					if !running || a.state != Fatigue || a.speaking {
						step()
						return
					}
					a.mouthSig.Set("astonishment")
					a.eyesSig.Set("closed")
					timer = time.AfterFunc(2000, func() {
						if !running || a.state != Fatigue || a.speaking {
							step()
							return
						}
						a.mouthSig.Set("talk-o")
						timer = time.AfterFunc(2000, func() {
							if !running || a.state != Fatigue || a.speaking {
								step()
								return
							}
							a.mouthSig.Set("fatigue")
							a.eyesSig.Set("fatigue")
							step()
						})
					})
				})
			})

		case Sad:
			a.browsSig.Set("sad")
			a.eyesSig.Set("sad")
			a.noseSig.Set("sad")
			a.mouthSig.Set("sad")
			a.poseSig.Set("none")
			a.extraSig.Set("drop")
			timer = time.AfterFunc(500, step)

		case Astonishment:
			a.browsSig.Set("astonishment")
			a.eyesSig.Set("astonishment")
			a.noseSig.Set("none")
			a.mouthSig.Set("astonishment")
			a.poseSig.Set("none")
			a.extraSig.Set("marks")
			timer = time.AfterFunc(500, step)

		case Anger:
			a.browsSig.Set("anger")
			a.eyesSig.Set("anger")
			a.noseSig.Set("none")
			a.mouthSig.Set("anger")
			a.poseSig.Set("none")
			a.extraSig.Set("none")
			timer = time.AfterFunc(500, step)

		default:
			timer = time.AfterFunc(500, step)
		}
	}

	step()
}
