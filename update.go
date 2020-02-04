package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"

	"github.com/spf13/viper"
	"golang.org/x/mod/modfile"
	"golang.org/x/mod/semver"
)

func runCmd(cmd string, args ...string) string {
	p := viper.GetString("path")
	if p == "" {
		p = "./"
	}

	c := exec.Command(cmd, args...)
	c.Stderr = os.Stderr
	c.Dir = p

	out, err := c.Output()
	if err != nil {
		fmt.Println(err)
		os.Exit(2)
	}

	out = bytes.TrimSpace(out)
	return string(out)
}

func updateCommand() {
	goModPath := runCmd("go", "env", "GOMOD")

	data, err := ioutil.ReadFile(goModPath)
	if err != nil {
		fmt.Println(err)
		os.Exit(2)
	}

	mod, err := modfile.Parse(goModPath, data, nil)
	if err != nil {
		fmt.Println(err)
		os.Exit(2)
	}

	for _, r := range mod.Require {
		if r.Indirect {
			fmt.Printf("%s %s // indirect\n", r.Mod.Path, r.Mod.Version)
			continue
		}

		max := r.Mod.Version
		out := runCmd("go", "list", "-m", "-versions", r.Mod.Path)

		if items := strings.Fields(out); len(items) > 0 {
			items = items[1:]

			for i := range items {
				if !semver.IsValid(items[i]) || semver.Canonical(items[i]) == "" || semver.Prerelease(items[i]) != "" {
					continue
				}

				if semver.Compare(items[i], max) == 1 {
					max = items[i]
				}
			}
		}

		fmt.Printf("%s %s", r.Mod.Path, max)
		if max != r.Mod.Version {
			fmt.Printf(" // %s => %s\n", r.Mod.Version, max)
		} else {
			fmt.Println()
		}
	}
}
