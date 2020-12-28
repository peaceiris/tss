// Command tss prints timestamps relative to the program start, and the previous
// line of input.
package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"runtime"

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
var buildVersion string
var buildCommit string
var buildDate string

func BuildVersionString() string {
	version := buildVersion
	commit := buildCommit

	date := buildDate
	if date == "" {
		date = "unknown"
	}

	return fmt.Sprintf(`TSS_BUILD_VERSION="%s"
TSS_BUILD_COMMIT="%s"
TSS_BUILD_DATE="%s"
TSS_BUILD_GOOS="%s"
TSS_BUILD_GOARCH="%s"
TSS_BUILD_GOVERSION="%s"`, version, commit, date, runtime.GOOS, runtime.GOARCH, runtime.Version())
}

func main() {
	version := flag.Bool("version", false, "Print the version string")
	v := flag.Bool("v", false, "Print the version string")
	flag.Parse()
	if *version || *v {
		fmt.Println(BuildVersionString())
		os.Exit(0)
	}
	if _, err := tss.Copy(os.Stdout, os.Stdin); err != nil {
		log.Fatal(err)
	}
}
