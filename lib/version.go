package tss

import "fmt"

// tss build env
var (
	buildVersion string = "0.1.2-development"
	buildCommit  string = "aaeb6725631dcff02055855ee263ef5f45ed1eea-development"
	buildDate    string = "2020-12-28T11:01:32Z-development"
)

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
