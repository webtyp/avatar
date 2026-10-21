package main

import (
	"flag"
	"fmt"
	"os"
	"webtyp.com/avatar"
)

func main() {
	prefix := flag.String("prefix", "", "prefix")
	svgFile := flag.String("svg", "", "svg file")
	pkg := flag.String("pkg", "", "package name")
	out := flag.String("out", "", "output file")
	flag.Parse()

	if *prefix == "" || *svgFile == "" || *pkg == "" || *out == "" {
		fmt.Fprintln(os.Stderr, "missing required flags")
		os.Exit(2)
	}

	ids, missing, err := avatar.Check(*prefix, *svgFile)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(2)
	}

	if len(missing) > 0 {
		for _, m := range missing {
			fmt.Fprintf(os.Stderr, "missing symbol: %s\n", m)
		}
		fmt.Fprintf(os.Stderr, "%d missing symbols in %s\n", len(missing), *svgFile)
		os.Exit(1)
	}

	res := avatar.GenerateConstants(*pkg, *prefix, ids)
	if err := os.WriteFile(*out, res, 0644); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(2)
	}
}
