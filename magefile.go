// +build mage

package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
)

const (
	packageName = "github.com/peaceiris/tss"
)

var ldflags = `-X $PACKAGE/cmd.buildVersion=$VERSION -X $PACKAGE/cmd.buildCommit=$COMMIT_HASH -X $PACKAGE/cmd.buildDate=$BUILD_DATE`

// allow user to override go executable by running as GOEXE=xxx make ... on unix-like systems
var goexe = "go"

func init() {
	if exe := os.Getenv("GOEXE"); exe != "" {
		goexe = exe
	}

	// We want to use Go 1.11 modules even if the source lives inside GOPATH.
	// The default is "auto".
	os.Setenv("GO111MODULE", "on")
}

func runWith(env map[string]string, cmd string, inArgs ...interface{}) error {
	s := argsToStrings(inArgs...)
	return sh.RunWith(env, cmd, s...)
}

// Install development tools
func Setup() error {
	if err := sh.Run(goexe, "get", "-u", "golang.org/x/lint/golint"); err != nil {
		return err
	}
	if err := sh.Run(goexe, "get", "-u", "github.com/stretchr/testify"); err != nil {
		return err
	}
	if err := sh.Run(goexe, "mod", "tidy"); err != nil {
		return err
	}
	if err := sh.Run(goexe, "mod", "verify"); err != nil {
		return err
	}
	return nil
}

// Build binary
func Build() error {
	return runWith(flagEnv(), goexe, "build", "-ldflags", ldflags, buildFlags(), packageName)
}

// Build binary with race detector enabled
func BuildRace() error {
	return runWith(flagEnv(), goexe, "build", "-race", "-ldflags", ldflags, buildFlags(), packageName)
}

// Install binary
func Install() error {
	return runWith(flagEnv(), goexe, "install", "-ldflags", ldflags, buildFlags(), packageName)
}

// Uninstall binary
func Uninstall() error {
	return sh.Run(goexe, "clean", "-i", packageName)
}

func flagEnv() map[string]string {
	hash, _ := sh.Output("git", "rev-parse", "--short", "HEAD")
	version, _ := sh.Output("git", "describe", "--tags")
	return map[string]string{
		"VERSION":     version,
		"PACKAGE":     packageName,
		"COMMIT_HASH": hash,
		"BUILD_DATE":  time.Now().Format("2006-01-02T15:04:05Z0700"),
	}
}

// Run tests and linters
func Check() {

	mg.Deps(Fmt, Vet)

	// don't run two tests in parallel, they saturate the CPUs anyway, and running two
	// causes memory issues in CI.
	mg.Deps(TestRace)
}

func testGoFlags() string {
	if isCI() {
		return ""
	}

	return "-test.short"
}

// Run tests
func Test() error {
	env := map[string]string{"GOFLAGS": testGoFlags()}
	return runCmd(env, goexe, "test", "-coverpkg", "./...", "-covermode", "atomic", "-coverprofile", "coverage.txt", "./...", buildFlags())
}

// Run tests with race detector
func TestRace() error {
	env := map[string]string{"GOFLAGS": testGoFlags()}
	return runCmd(env, goexe, "test", "-race", "-coverpkg", "./...", "-covermode", "atomic", "-coverprofile", "coverage.txt", "./...", buildFlags())
}

// Run gofmt -w
func Fmtw() error {
	if !isGoLatest() {
		return nil
	}
	pkgs, err := tssPackages()
	if err != nil {
		return err
	}
	for _, pkg := range pkgs {
		files, err := filepath.Glob(filepath.Join(pkg, "*.go"))
		if err != nil {
			return nil
		}
		for _, f := range files {
			_, err := sh.Output("gofmt", "-w", f)
			if err != nil {
				fmt.Printf("ERROR: running gofmt on %q: %v\n", f, err)
			}
		}
	}
	return nil
}

// Run gofmt linter
func Fmt() error {
	if !isGoLatest() {
		return nil
	}
	pkgs, err := tssPackages()
	if err != nil {
		return err
	}
	failed := false
	first := true
	for _, pkg := range pkgs {
		files, err := filepath.Glob(filepath.Join(pkg, "*.go"))
		if err != nil {
			return nil
		}
		for _, f := range files {
			// gofmt doesn't exit with non-zero when it finds unformatted code
			// so we have to explicitly look for output, and if we find any, we
			// should fail this target.
			s, err := sh.Output("gofmt", "-l", f)
			if err != nil {
				fmt.Printf("ERROR: running gofmt on %q: %v\n", f, err)
				failed = true
			}
			if s != "" {
				if first {
					fmt.Println("The following files are not gofmt'ed:")
					first = false
				}
				failed = true
				fmt.Println(f)
			}
		}
	}
	if failed {
		return errors.New("improperly formatted go files")
	}
	return nil
}

var (
	pkgPrefixLen = len("github.com/peaceiris/tss")
	pkgs         []string
	pkgsInit     sync.Once
)

func tssPackages() ([]string, error) {
	var err error
	pkgsInit.Do(func() {
		var s string
		s, err = sh.Output(goexe, "list", "./...")
		if err != nil {
			return
		}
		pkgs = strings.Split(s, "\n")
		for i := range pkgs {
			pkgs[i] = "." + pkgs[i][pkgPrefixLen:]
		}
	})
	return pkgs, err
}

// Run golint linter
func Lint() error {
	pkgs, err := tssPackages()
	if err != nil {
		return err
	}
	failed := false
	for _, pkg := range pkgs {
		// We don't actually want to fail this target if we find golint errors,
		// so we don't pass -set_exit_status, but we still print out any failures.
		if _, err := sh.Exec(nil, os.Stderr, nil, "golint", pkg); err != nil {
			fmt.Printf("ERROR: running go lint on %q: %v\n", pkg, err)
			failed = true
		}
	}
	if failed {
		return errors.New("errors running golint")
	}
	return nil
}

// Run go vet linter
func Vet() error {
	if err := sh.Run(goexe, "vet", "./..."); err != nil {
		return fmt.Errorf("error running go vet: %v", err)
	}
	return nil
}

func runCmd(env map[string]string, cmd string, args ...interface{}) error {
	if mg.Verbose() {
		return runWith(env, cmd, args...)
	}
	output, err := sh.OutputWith(env, cmd, argsToStrings(args...)...)
	if err != nil {
		fmt.Fprint(os.Stderr, output)
	}

	return err
}

func isGoLatest() bool {
	return strings.Contains(runtime.Version(), "1.15")
}

func isCI() bool {
	return os.Getenv("CI") != ""
}

func buildFlags() []string {
	if runtime.GOOS == "windows" {
		return []string{"-buildmode", "exe"}
	}
	return nil
}

func argsToStrings(v ...interface{}) []string {
	var args []string
	for _, arg := range v {
		switch v := arg.(type) {
		case string:
			if v != "" {
				args = append(args, v)
			}
		case []string:
			if v != nil {
				args = append(args, v...)
			}
		default:
			panic("invalid type")
		}
	}

	return args
}

func BumpVersion(releaseType string) error {
	fmt.Printf("\nnpx standard-version --release-as %s --preset eslint\n", releaseType)
	out, err := sh.Output("npx", "standard-version", "--release-as", releaseType, "--preset", "eslint")
	fmt.Println(out)
	if err != nil {
		return err
	}
	if err := sh.Run("git", "push", "origin", "main"); err != nil {
		return err
	}
	if err := sh.Run("git", "push", "origin", "--tags"); err != nil {
		return err
	}
	return nil
}

// Run npx standard-version --release-as patch --preset eslint
func BumpPatchVersion() error {
	return BumpVersion("patch")
}

// Run npx standard-version --release-as minor --preset eslint
func BumpMinorVersion() error {
	return BumpVersion("minor")
}

// Run npx standard-version --release-as patch --preset eslint --dry-run
func BumpVersionTest() error {
	fmt.Printf("\nnpx standard-version --release-as patch --preset eslint --dry-run\n")
	out, err := sh.Output("npx", "standard-version", "--release-as", "patch", "--preset", "eslint", "--dry-run")
	fmt.Println(out)
	return err
}

// Run goreleaser check
func GoreleaserCheck() error {
	fmt.Printf("\ngoreleaser check\n")
	return sh.Run("goreleaser", "check")
}

// Run goreleaser --snapshot --skip-publish --rm-dist
func GoreleaserTest() error {
	fmt.Printf("\ngoreleaser --snapshot --skip-publish --rm-dist\n")
	return sh.Run("goreleaser", "--snapshot", "--skip-publish", "--rm-dist")
}

// Run GoreleaserCheck, GoreleaserTest, BumpVersionTest
func ReleaseTest() {
	mg.SerialDeps(GoreleaserCheck, GoreleaserTest, BumpVersionTest)
}
