// Command tss prints timestamps relative to the program start, and the previous
// line of input.
package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"runtime"

	tss "github.com/peaceiris/tss/cmd"
)

func init() {
	flag.Usage = func() {
		fmt.Println(`tss [-v] [-h]

Annotate stdin with timestamps per line.`)
	}
}

func main() {
	version := flag.Bool("version", false, "Print the version string")
	v := flag.Bool("v", false, "Print the version string")
	flag.Parse()
	if *version || *v {
		fmt.Printf(tss.GetBuildEnvsString(runtime.GOOS, runtime.GOARCH, runtime.Version()))
		os.Exit(0)
	}
	if _, err := tss.Copy(os.Stdout, os.Stdin); err != nil {
		log.Fatal(err)
	}
}
