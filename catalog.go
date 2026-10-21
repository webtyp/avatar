package avatar

import "webtyp.com/svg"

type variant uint8

const (
	variantNormal variant = iota
	variantGreeting
	variantSad
	variantAstonishment
	variantAnger
	variantFatigue
	variantListen
	variantThink
	variantTalkA
	variantTalkO
	variantClosed
	variantCount
)

var variantNames = [variantCount]string{
	"normal",
	"greeting",
	"sad",
	"astonishment",
	"anger",
	"fatigue",
	"listen",
	"think",
	"talk-a",
	"talk-o",
	"closed",
}

// Catalog is a value type. Every With* returns a copy.
type Catalog struct {
	prefix   string
	torso    string
	head     string
	hair     string
	shirt    string
	earrings string
	glasses  string

	// Precomputed face variant hrefs for brows, eyes, mouth
	face [3][variantCount]string

	noseSad      string
	noseClosed   string
	extrasDrop   string
	extrasMarks  string
	mouthTalkA   string
	mouthTalkO   string
	poseGreeting string
	poseThink    string
	poseTalkA    string
	poseTalkO    string
	speechA      string
	speechO      string
	bubble1      string
	bubble2      string
	bubble3      string
	bubble3Ring  string
	dot1         string
	dot2         string
	dot3         string
	soundInner   string
	soundMid     string
	soundOuter   string
	zzz1         string
	zzz2         string
	zzz3         string
}

func NewCatalog(prefix string) Catalog {
	c := Catalog{
		prefix: prefix,
		torso:  "#" + prefix + "-torso-normal",
		head:   "#" + prefix + "-head-normal",
		hair:   "#" + prefix + "-hair-normal",

		noseSad:      "#" + prefix + "-" + symNoseSad,
		noseClosed:   "#" + prefix + "-" + symNoseClosed,
		extrasDrop:   "#" + prefix + "-" + symExtrasDrop,
		extrasMarks:  "#" + prefix + "-" + symExtrasMarks,
		mouthTalkA:   "#" + prefix + "-" + symMouthTalkA,
		mouthTalkO:   "#" + prefix + "-" + symMouthTalkO,
		poseGreeting: "#" + prefix + "-" + symPoseGreeting,
		poseThink:    "#" + prefix + "-" + symPoseThink,
		poseTalkA:    "#" + prefix + "-" + symPoseTalkA,
		poseTalkO:    "#" + prefix + "-" + symPoseTalkO,
		speechA:      "#" + prefix + "-" + symSpeechA,
		speechO:      "#" + prefix + "-" + symSpeechO,
		bubble1:      "#" + prefix + "-" + symBubble1,
		bubble2:      "#" + prefix + "-" + symBubble2,
		bubble3:      "#" + prefix + "-" + symBubble3,
		bubble3Ring:  "#" + prefix + "-" + symBubble3Ring,
		dot1:         "#" + prefix + "-" + symDot1,
		dot2:         "#" + prefix + "-" + symDot2,
		dot3:         "#" + prefix + "-" + symDot3,
		soundInner:   "#" + prefix + "-" + symSoundInner,
		soundMid:     "#" + prefix + "-" + symSoundMid,
		soundOuter:   "#" + prefix + "-" + symSoundOuter,
		zzz1:         "#" + prefix + "-" + symZzz1,
		zzz2:         "#" + prefix + "-" + symZzz2,
		zzz3:         "#" + prefix + "-" + symZzz3,
	}

	for l, layerName := range faceLayers {
		for v, varName := range variantNames {
			c.face[l][v] = "#" + prefix + "-" + layerName + "-" + varName
		}
	}

	return c
}

func (c Catalog) WithTorso(icon svg.Icon) Catalog {
	c.torso = "#" + string(icon)
	return c
}

func (c Catalog) WithHead(icon svg.Icon) Catalog {
	c.head = "#" + string(icon)
	return c
}

func (c Catalog) WithHair(icon svg.Icon) Catalog {
	c.hair = "#" + string(icon)
	return c
}

func (c Catalog) WithShirt(icon svg.Icon) Catalog {
	c.shirt = "#" + string(icon)
	return c
}

func (c Catalog) WithEarrings(icon svg.Icon) Catalog {
	c.earrings = "#" + string(icon)
	return c
}

func (c Catalog) WithGlasses(icon svg.Icon) Catalog {
	c.glasses = "#" + string(icon)
	return c
}
