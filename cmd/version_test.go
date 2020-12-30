package cmd

import (
	"testing"
)

func TestBuildEnvsString(t *testing.T) {
	want := `TSS_BUILD_VERSION="0.1.2-development"
TSS_BUILD_COMMIT="aaeb6725631dcff02055855ee263ef5f45ed1eea-development"
TSS_BUILD_DATE="2020-12-28T11:01:32Z-development"
TSS_BUILD_GOOS="darwin"
TSS_BUILD_GOARCH="amd64"
TSS_BUILD_GOVERSION="go1.15.6"
`
	out := GetBuildEnvsString("darwin", "amd64", "go1.15.6")
	if out != want {
		t.Errorf("BuildEnvString(); want \n%s", want)
	}
}
