package test

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/likexian/gokit/assert"
)

var (
	binName = "go-ucce"
	dir     = "../"
)

func TestMain(m *testing.M) {

	fmt.Println("Building tool...")

	if runtime.GOOS == "windows" {
		binName += ".exe"
	}

	build := exec.Command("cd", dir, "&&", "go", "build", "-o", binName)

	if err := build.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Cannot build tool %s: %s", binName, err)
		os.Exit(1)
	}

	fmt.Println("Running tests....")
	result := m.Run()

	fmt.Println("Cleaning up...")
	os.Remove(binName)

	os.Exit(result)
}

func TestLicence(t *testing.T) {
	cmdPath := filepath.Join("../", "go-uuce")

	t.Run("Launch binary", func(t *testing.T) {

		cmd := exec.Command(cmdPath, "licence", "-a", "localhost", "-u", "user", "-x", "password", "-t", "finesse")
		out, _ := cmd.CombinedOutput()

		expected := "this ucce instance type doesn't exist finesse" + "\n"

		assert.Contains(t, string(out), expected)
	})
}

type Cmd struct {
	args interface{}
}

func TestFlagsRequired(t *testing.T) {

	cmdPath := filepath.Join("../", "go-uuce")

	var testCases = []struct {
		name string
		cmd  string
		out  string
		//err  error
	}{
		{name: "none flag", cmd: cmdPath + "status", out: "Error: required flag(s) \"host\", \"password\", \"typeOf\", \"user\" not set"},
		{name: "one flag", cmd: cmdPath + "status" + "-u" + "user", out: "Error: required flag(s) \"password\", \"typeOf\", \"user\" not set"},
		//{name: "two flags", cmd: "", out: "Error: required flag(s) \"password\", \"typeOf\" not set"},
		//{name: "three flags", cmd: "", out: "Error: required flag(s) \"typeOf\" not set"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {

			cmd := exec.Command(tc.cmd)
			out, err := cmd.CombinedOutput()
			if err != nil {
				fmt.Print(err)
			}

			assert.Contains(t, string(out), tc.out)
		})
	}

	// t.Run("Launch binary", func(t *testing.T) {

	// 	cmd := exec.Command(cmdPath, "status")
	// 	out, _ := cmd.CombinedOutput()

	// 	expected := "Error: required flag(s) \"host\", \"password\", \"typeOf\", \"user\" not set" + "\n"

	// 	assert.Contains(t, string(out), expected)
	// })
}
