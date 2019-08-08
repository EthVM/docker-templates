package main

import (
	"os"
	"strings"

	"github.com/One-com/gonelog/log"
	"gopkg.in/urfave/cli.v1"
)

const (
	defaultDelims = "{{:}}"
)

var (
	// Version generated from Makefile
	VERSION string

	stdOutFlag = cli.BoolFlag{
		Name:        "stdout",
		Usage:       "forces output to be written to stdout",
		Destination: &forceStdOutFlag,
	}

	logLevelFlag = cli.IntFlag{
		Name:        "log-level",
		Value:       4,
		Usage:       "log level to emit to the screen",
		Destination: &logLevel,
	}

	forceFlag = cli.BoolFlag{
		Name:        "no-overwrite",
		Usage:       "do not overwrite destination file if it already exists",
		Destination: &noOverwriteFlag,
	}

	delimsFlag = cli.StringFlag{
		Name:        "delims",
		Usage:       `template tag delimiters. Default "{{":"}}"`,
		Value:       defaultDelims,
		Destination: &rawDelims,
	}

	renderCommand = cli.Command{
		Action:    renderCmd,
		Name:      "render",
		Usage:     "renders the specified definition file(s)",
		ArgsUsage: "[definition]",
	}

	app             = cli.NewApp()
	rawDelims       string
	delims          = strings.Split(defaultDelims, ":")
	noOverwriteFlag bool
	forceStdOutFlag bool
	logLevel        int

	l = log.New(os.Stdout, "", log.Lcolor|log.Ldate|log.Ltime|log.Llevel)
)

func init() {
	app.Name = "docker-templates"
	app.Usage = "render Docker Compose / Stack file templates with the power of go templates"
	app.Version = VERSION
	app.Flags = []cli.Flag{
		stdOutFlag,
		delimsFlag,
		logLevelFlag,
	}
	app.Commands = []cli.Command{renderCommand}

	if rawDelims != "" {
		delims := strings.Split(rawDelims, ":")
		if len(delims) != 2 {
			l.Panicf("Bad delimiters argument: %s. Expected \"left:right\"", delims)
		}
	}
}

func main() {
	if err := app.Run(os.Args); err != nil {
		l.Panicf("Error found! Aborting execution: %s", err)
	}
}
