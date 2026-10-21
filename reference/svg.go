//go:build !wasm

package reference

import (
	_ "embed"
	"webtyp.com/svg/sprite"
)

//go:embed cote.svg
var coteSVG string

// IconSvg is the reference character's sprite, extracted by sitec at build time.
func IconSvg() *sprite.Sprite {
	s := sprite.NewSprite()
	if err := s.AddLayers(Prefix, coteSVG); err != nil {
		panic(err) // build time only: a broken reference drawing must stop the build
	}
	return s
}
