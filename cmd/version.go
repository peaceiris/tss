package cmd

import (
	"fmt"
	"runtime"

	"github.com/spf13/cobra"
)

// tss build env
var (
	buildVersion string = "0.1.2-development"
	buildCommit  string = "aaeb6725631dcff02055855ee263ef5f45ed1eea-development"
	buildDate    string = "2020-12-28T11:01:32Z-development"
)

// NewCmdVersion creates version command
func NewCmdVersion() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "version",
		Short: "Print the version number of tss",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Printf(GetBuildEnvsString(runtime.GOOS, runtime.GOARCH, runtime.Version()))
		},
	}
	return cmd
}

// BuildEnv has shell environment name and value
type BuildEnv struct {
	name  string
	value string
}

// GetBuildEnvsString returns tss build environment variables
func GetBuildEnvsString(goOS string, goArch string, goVersion string) string {
	version := buildVersion
	commit := buildCommit

	buildEnvs := [6]BuildEnv{
		{name: "VERSION", value: version},
		{name: "COMMIT", value: commit},
		{name: "DATE", value: buildDate},
		{name: "GOOS", value: goOS},
		{name: "GOARCH", value: goArch},
		{name: "GOVERSION", value: goVersion},
	}

	out := ""
	for _, e := range buildEnvs {
		out += fmt.Sprintf("TSS_BUILD_%s=\"%s\"\n", e.name, e.value)
	}
	return out
}
