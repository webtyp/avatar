package avatar

import "webtyp.com/svg"

type Layer uint8

const (
	LayerTorso Layer = iota
	LayerHead
	LayerHair
	LayerBrows
	LayerEyes
	LayerNose
	LayerMouth
	LayerPoses
	LayerExtras
)

type Catalog struct {
	Prefix string
	Hair   svg.Icon
	Torso  svg.Icon
	Head   svg.Icon
}

// NewCatalog constructs a Catalog using the standard Inkscape layer/variant ID conventions:
// e.g. <prefix>-torso-normal, <prefix>-hair-normal, <prefix>-eyes-<state>, etc.
func NewCatalog(prefix string) *Catalog {
	return &Catalog{Prefix: prefix}
}

func (c *Catalog) WithHair(icon svg.Icon) *Catalog {
	c.Hair = icon
	return c
}

func (c *Catalog) WithTorso(icon svg.Icon) *Catalog {
	c.Torso = icon
	return c
}

func (c *Catalog) WithHead(icon svg.Icon) *Catalog {
	c.Head = icon
	return c
}
