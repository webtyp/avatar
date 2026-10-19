//go:build !wasm

package avatar

import (
	"webtyp.com/css"
	"webtyp.com/widget/style"
)

func (a *Avatar) RenderCSS() *css.Stylesheet {
	return style.For(a).
		Root(
			style.CenterContent(),
			style.KeepSize(),
		).
		Stylesheet()
}
