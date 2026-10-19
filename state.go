package avatar

type State uint8

const (
	Normal State = iota
	Greeting
	Speak
	Think
	Listen
	Fatigue
	Sad
	Astonishment
	Anger
)

func (s State) String() string {
	switch s {
	case Normal:
		return "normal"
	case Greeting:
		return "greeting"
	case Speak:
		return "speak"
	case Think:
		return "think"
	case Listen:
		return "listen"
	case Fatigue:
		return "fatigue"
	case Sad:
		return "sad"
	case Astonishment:
		return "astonishment"
	case Anger:
		return "anger"
	default:
		return "unknown"
	}
}
