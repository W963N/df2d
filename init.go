package main

import (
	"flag"
	"fmt"
	"io"

	"github.com/google/goterm/term"
)

const (
	TITLE = "df2d"
)

func clyFprintf(w io.Writer, format string) {
	fmt.Fprintf(w, term.Yellowf(format))
}

func clrFprintf(w io.Writer, format string) {
	fmt.Fprintf(w, term.Redf(format))
}

func clbFprintf(w io.Writer, format string) {
	fmt.Fprintf(w, term.Bluef(format))
}

func init() {
	flag.CommandLine.Init(TITLE, flag.ContinueOnError)
	flag.CommandLine.Usage = func() {
		o := flag.CommandLine.Output()
		clyFprintf(o, "\nW963N Memory ~~ "+flag.CommandLine.Name()+"\n")
		clbFprintf(o, "\n@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@\n")
		clbFprintf(o, "@@@@@@@@@@@@@  .@@@@@@@@@.  @@@@@@@@@@@@\n")
		clbFprintf(o, "@@@@@@@@@ (@@@@@@@@@@@@@@@@@@@) @@@@@@@@\n")
		clbFprintf(o, "@@@@@@ @@ @@@@             @@@@ @@ @@@@@\n")
		clbFprintf(o, "@@@@ @@@@@ @ @@@         @@@ @ @@@@ @@@@\n")
		clbFprintf(o, "@@@ @@@@@@ @@@@@@       @@@@@@ @@@@@ @@@\n")
		clbFprintf(o, "@@ @@@@@@,@@@@@@@       @@@@@@@,@@@@@ @@\n")
		clbFprintf(o, "@@ @@@@@@ @@@@@@@@     @@@@@@@@ @@@@@@ @\n")
		clbFprintf(o, "@       @@@      @     @      @@@      @\n")
		clbFprintf(o, "@@ @@@ @@@@@@   @@@   @@@   @@@@@@ @@ @@\n")
		clbFprintf(o, "@@@ @@ @@@@@@@@@@@@   @@@@@@@@@@@@ @@ @@\n")
		clbFprintf(o, "@@@@  @@@@@@@@@@@@@@ @@@@@@@@@@@@@ @ @@@\n")
		clbFprintf(o, "@@@@@  @@@@@@@@@@@@@ @@@@@@@@@@@@   @@@@\n")
		clbFprintf(o, "@@@@@@@   @@@@@@@  @@@  @@@@@@@   @@@@@@\n")
		clbFprintf(o, "@@@@@@@@@@ @@@@@@@@@@@@@@@@@@@ @@@@@@@@@\n")
		clbFprintf(o, "@@@@@@@@@@@@@@@     @     @@@@@@@@@@@@@@\n\n")
		clrFprintf(o, "\nUsage: \n")
		fmt.Fprintf(o, "  df2d "+term.Greenf("-e")+" [Env File]\n")
		clrFprintf(o, "\nOptions: \n")
		flag.PrintDefaults()
		clbFprintf(o, "\nHobbyright 2023 walnut 🐿🐿🐿 .\n\n")
	}
	flag.StringVar(&env, "e", "./env.toml", "path of env.toml.")
}

var (
	env string
)
