package cmd

import (
	"bytes"
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBuildEnvsString(t *testing.T) {
	command := "tss version"
	want := `TSS_BUILD_VERSION="0.1.2-development"
TSS_BUILD_COMMIT="aaeb6725631dcff02055855ee263ef5f45ed1eea-development"
TSS_BUILD_DATE="2020-12-28T11:01:32Z-development"
TSS_BUILD_GOOS="darwin"
TSS_BUILD_GOARCH="amd64"
TSS_BUILD_GOVERSION="go1.15.6"
`

	buf := new(bytes.Buffer)
	cmd := NewCmdVersion()
	cmd.SetOutput(buf)
	cmdArgs := strings.Split(command, " ")
	fmt.Printf("cmdArgs %+v\n", cmdArgs)
	cmd.SetArgs(cmdArgs[1:])
	cmd.Execute()
	out := buf.String()

	assert.Equal(t, out, want, "they should be equal")
}
