// Command tss prints timestamps relative to the program start, and the previous
// line of input.
package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	tss "github.com/peaceiris/tss/lib"
)

func init() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, `tss [-v] [-h]

Annotate stdin with timestamps per line.
`)
	}
}

// Version tss version
var Version string
var Commit string
var Date string

func main() {
	version := flag.Bool("version", false, "Print the version string")
	v := flag.Bool("v", false, "Print the version string")
	flag.Parse()
	if *version || *v {
		envString := fmt.Sprintf(`TSS_VERSION=\"%s\"
TSS_BUILD_COMMIT=\"%s\"
TSS_BUILD_DATE=\"%s\"
`, Version, Commit, Date)
		fmt.Println(envString)
		os.Exit(0)
	}
	if _, err := tss.Copy(os.Stdout, os.Stdin); err != nil {
		log.Fatal(err)
	}
}
