package avatar

type Emotion uint8

const (
	Neutral Emotion = iota
	Happy
	Sad
	Astonished
	Angry
)

func (e Emotion) String() string {
	switch e {
	case Happy:
		return "happy"
	case Sad:
		return "sad"
	case Astonished:
		return "astonished"
	case Angry:
		return "angry"
	case Neutral:
		fallthrough
	default:
		return "neutral"
	}
}
