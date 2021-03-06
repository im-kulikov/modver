package main

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/pkg/errors"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

type empty int

const (
	devNull = empty(0)
	about   = "\nSimple program to receive latest version (for go.mod) of passed package.\nVersion %s (%s).\n"
)

// Read returns nothing
func (empty) Read([]byte) (int, error) { return 0, io.EOF }

func initSettings() {
	var (
		h, v  bool
		u, dr bool
	)

	// flags setup:
	flags := pflag.NewFlagSet("commandline", pflag.ExitOnError)
	flags.SortFlags = false

	flags.Bool("verbose", false, "verbose")

	flags.BoolVarP(&h, "help", "h", false, "show help")
	flags.BoolVarP(&v, "version", "v", false, "show version")

	flags.BoolVarP(&u, "update", "u", false, "fetch and update current modules")
	flags.BoolVar(&dr, "dry-run", false, "only check updates for current modules (use with -u / --update)")

	flags.Bool(commitFlag, false, "display latest commit version (for example v0.0.0-<hash>-<date>)")
	flags.String(branchFlag, "master", "use passed branch to receive version (for remote repos only)")

	// set prefers:
	viper.Set("app.name", "modver")
	viper.Set("app.author", "Evgeniy Kulikov <im@kulikov.im>")
	viper.Set("app.version", version+"("+build+")")

	if err := viper.BindPFlags(flags); err != nil {
		panic(err)
	}

	if err := viper.ReadConfig(devNull); err != nil {
		panic(err)
	}

	if err := flags.Parse(os.Args); err != nil {
		panic(err)
	}

	switch {
	case h:
		fmt.Printf(about, version, build)
		fmt.Println("modver [global options] {repo-path or url}")
		flags.PrintDefaults()
		os.Exit(0)
	case v:
		fmt.Printf(about, version, build)
		os.Exit(0)
	case u:
		updateCommand(dr)

		os.Exit(0)
	}

	if args := flags.Args(); len(args) >= 2 {
		viper.SetDefault("path", args[1])
	} else {
		viper.SetDefault("path", "./")
	}
}

func main() {
	initSettings()

	if err := latestRevision(); err == nil || errors.Cause(err) == context.Canceled {
		os.Exit(0)
	} else if !viper.GetBool("verbose") {
		fmt.Println(err)
		os.Exit(2)
	} else {
		panic(err)
	}
}
