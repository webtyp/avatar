package reference

import "webtyp.com/avatar"

//go:generate go run webtyp.com/avatar/cmd/avatarc -prefix cote -svg cote.svg -pkg reference -out symbols_gen.go

// Catalog is the reference character: Cote, drawn following the Inkscape layer convention.
func Catalog() avatar.Catalog {
	return avatar.NewCatalog(Prefix)
}
