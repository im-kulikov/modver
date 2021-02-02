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
	buf := new(bytes.Buffer)

	p := viper.GetString("path")
	if p == "" {
		p = "./"
	}

	c := exec.Command(cmd, args...)
	c.Stdout = buf
	c.Stderr = os.Stderr
	c.Dir = p

	if err := c.Run(); err != nil {
		fmt.Printf("Something went wrong: %s\n", err)
		os.Exit(2)
	}

	out := buf.Bytes()
	out = bytes.TrimSpace(out)
	return string(out)
}

func updateCommand(dryRun bool) {
	if dryRun {
		fmt.Printf("--dry-run option enabled, check only for updates\n\n")
	}

	goModPath := runCmd("go", "env", "GOMOD")
	if goModPath == "" {
		fmt.Println("Could not find modfile")
		os.Exit(2)
	}

	info, err := os.Stat(goModPath)
	if err != nil {
		fmt.Printf("Could not check file: %s\n", err)
		os.Exit(2)
	}

	data, err := ioutil.ReadFile(goModPath)
	if err != nil {
		fmt.Printf("Could not read modfile: %s\n", err)
		os.Exit(2)
	}

	mod, err := modfile.Parse(goModPath, data, nil)
	if err != nil {
		fmt.Printf("Could not parse modfile: %s\n", err)
		os.Exit(2)
	}

	changed := false
	result := new(strings.Builder)
	for _, r := range mod.Require {
		if r.Indirect {
			continue
		}

		now := r.Mod.Version
		max := r.Mod.Version
		out := runCmd("go", "list", "-m", "-mod=mod", "-versions", r.Mod.Path)

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

		fmt.Printf(".")

		if max != now {
			changed = true
			mod.AddRequire(r.Mod.Path, max)
			_, _ = fmt.Fprintf(result, "%s %s // %s => %s\n", r.Mod.Path, max, now, max)
		}
	}

	fmt.Println()

	if !changed {
		fmt.Println("Everything is up to date")
		return
	}

	_, _ = os.Stdout.WriteString(result.String())

	if dryRun {
		return
	}

	out, err := mod.Format()
	if err != nil {
		fmt.Printf("Could not format: %s\n", err)
		os.Exit(2)
	} else if err = ioutil.WriteFile(goModPath, out, info.Mode()); err != nil {
		fmt.Printf("Could not write file: %s\n", err)
		os.Exit(2)
	}
}
