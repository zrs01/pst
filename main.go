package main

import (
	"os"

	"github.com/rotisserie/eris"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
	"github.com/zrs01/pst/internal/converter"
)

func main() {
	cliapp := cli.NewApp()
	cliapp.Name = "dst"
	cliapp.Usage = "Database schema tool"
	cliapp.Version = "0.0.1-202206"
	cliapp.Commands = []*cli.Command{}

	debug := false
	var ifile, ofile string

	cliapp.Flags = []cli.Flag{
		&cli.BoolFlag{
			Name:        "debug",
			Aliases:     []string{"d"},
			Usage:       "Debug mode",
			Required:    false,
			Destination: &debug,
		},
		&cli.StringFlag{
			Name:        "input",
			Aliases:     []string{"i"},
			Usage:       "Input file",
			Required:    true,
			Destination: &ifile,
		},
		&cli.StringFlag{
			Name:        "output",
			Aliases:     []string{"o"},
			Usage:       "Output file",
			Required:    true,
			Destination: &ofile,
		},
	}
	cliapp.Action = func(ctx *cli.Context) error {
		return converter.Build(ifile, ofile)
	}

	if err := cliapp.Run(os.Args); err != nil {
		logrus.Error(eris.ToString(err, debug))
	}
}
