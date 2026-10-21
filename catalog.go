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
}

func NewCatalog(prefix string) Catalog {
	c := Catalog{
		prefix: prefix,
		torso:  "#" + prefix + "-torso-normal",
		head:   "#" + prefix + "-head-normal",
		hair:   "#" + prefix + "-hair-normal",
	}

	layers := [3]string{"brows", "eyes", "mouth"}
	for l, layerName := range layers {
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
