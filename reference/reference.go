package reference

import "webtyp.com/avatar"

const Prefix = "cote" // symbol prefix of cote.svg

// Catalog is the reference character: Cote, drawn following the Inkscape layer convention.
func Catalog() avatar.Catalog {
	return avatar.NewCatalog(Prefix)
}
