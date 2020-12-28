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

// tss build env
var buildVersion string = "0.1.2-development"
var buildCommit string = "aaeb6725631dcff02055855ee263ef5f45ed1eea-development"
var buildDate string = "2020-12-28T11:01:32Z-development"

// BuildEnvString returns tss build environment variables
func BuildEnvString(goOS string, goArch string, goVersion string) string {
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
TSS_BUILD_GOVERSION="%s"`, version, commit, date, goOS, goArch, goVersion)
}

func main() {
	version := flag.Bool("version", false, "Print the version string")
	v := flag.Bool("v", false, "Print the version string")
	flag.Parse()
	if *version || *v {
		fmt.Println(BuildEnvString(runtime.GOOS, runtime.GOARCH, runtime.Version()))
		os.Exit(0)
	}
	if _, err := tss.Copy(os.Stdout, os.Stdin); err != nil {
		log.Fatal(err)
	}
}
