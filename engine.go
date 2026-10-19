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
	a.noseSig.Set("normal")
	a.mouthSig.Set("normal")
	a.poseSig.Set("normal")
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
			timer = time.AfterFunc(100, func() {
				if !running || !a.speaking {
					step()
					return
				}
				a.mouthSig.Set("talk-o")
				a.poseSig.Set("talk-o")
				timer = time.AfterFunc(100, func() {
					if !running || !a.speaking {
						step()
						return
					}
					a.mouthSig.Set("normal")
					a.poseSig.Set("normal")
					timer = time.AfterFunc(100, step)
				})
			})
			return
		}

		switch a.state {
		case Normal:
			a.eyesSig.Set("normal")
			sleepTime := 3500 + rand.Intn(1500)
			timer = time.AfterFunc(sleepTime, func() {
				if !running || a.state != Normal {
					step()
					return
				}
				a.eyesSig.Set("closed")
				timer = time.AfterFunc(150, step)
			})

		case Fatigue:
			a.eyesSig.Set("closed")
			sleepTime := 10000 + rand.Intn(20000)
			timer = time.AfterFunc(sleepTime, func() {
				if !running || a.state != Fatigue {
					step()
					return
				}
				a.mouthSig.Set("yawn-inhale")
				timer = time.AfterFunc(1500, func() {
					if !running || a.state != Fatigue {
						step()
						return
					}
					a.mouthSig.Set("yawn-wide")
					timer = time.AfterFunc(2000, func() {
						if !running || a.state != Fatigue {
							step()
							return
						}
						a.mouthSig.Set("yawn-exhale")
						timer = time.AfterFunc(2000, func() {
							if !running || a.state != Fatigue {
								step()
								return
							}
							a.mouthSig.Set("normal")
							step()
						})
					})
				})
			})

		case Think:
			a.extraSig.Set("think-1")
			timer = time.AfterFunc(500, func() {
				if !running || a.state != Think {
					step()
					return
				}
				a.extraSig.Set("think-2")
				timer = time.AfterFunc(500, func() {
					if !running || a.state != Think {
						step()
						return
					}
					a.extraSig.Set("think-3")
					timer = time.AfterFunc(500, step)
				})
			})

		case Greeting:
			a.poseSig.Set("wave-1")
			timer = time.AfterFunc(250, func() {
				if !running || a.state != Greeting {
					step()
					return
				}
				a.poseSig.Set("wave-2")
				timer = time.AfterFunc(250, step)
			})

		default:
			timer = time.AfterFunc(100, step)
		}
	}

	step()
}
