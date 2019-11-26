package main

import (
	"fmt"
	"os"

	"github.com/urfave/cli/v2"
)

func main() {
	app := cli.NewApp()
	app.Authors = []*cli.Author{
		{
			Name:  "Evgeniy Kulikov",
			Email: "im@kulikov.im",
		},
	}

	app.Flags = []cli.Flag{
		&cli.BoolFlag{
			Name:  commitFlag,
			Usage: "display latest commit version (for example v0.0.0-<hash>-<date>)",
		},

		&cli.StringFlag{
			Value: "master",
			Name:  branchFlag,
			Usage: "use passed branch to receive version (for remote repos only)",
		},
	}

	app.Usage = "Simple program to receive latest version (for go.mod) of passed package."
	app.UsageText = "modver [global options] {repo-path or url}"
	app.Version = fmt.Sprintf("%s (%s)", version, build)
	app.Action = latestRevision
	if err := app.Run(os.Args); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		cli.OsExiter(1)
	}
}
