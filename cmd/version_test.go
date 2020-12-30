package cmd

import (
	"bytes"
	"fmt"
	"runtime"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBuildEnvsString(t *testing.T) {
	command := "tss version"

	buildEnvs := [6]BuildEnv{
		{name: "VERSION", value: "0.1.2-development"},
		{name: "COMMIT", value: "aaeb6725631dcff02055855ee263ef5f45ed1eea-development"},
		{name: "DATE", value: "2020-12-28T11:01:32Z-development"},
		{name: "GOOS", value: runtime.GOOS},
		{name: "GOARCH", value: runtime.GOARCH},
		{name: "GOVERSION", value: runtime.Version()},
	}

	want := ""
	for _, e := range buildEnvs {
		want += fmt.Sprintf("TSS_BUILD_%s=\"%s\"\n", e.name, e.value)
	}

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
