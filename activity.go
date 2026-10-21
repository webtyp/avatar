package avatar

type Activity uint8

const (
	Idle Activity = iota
	Greeting
	Listen
	Think
	Speak
	Fatigue
)
