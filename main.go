package main

import (
	"os"

	"github.com/rotisserie/eris"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
	"github.com/zrs01/pst/internal/docb"
)

var version = "development"

func main() {
	cliapp := cli.NewApp()
	cliapp.Name = "pst"
	cliapp.Usage = "Program specfication tool"
	cliapp.Version = version
	cliapp.Commands = []*cli.Command{}

	debug := false
	var ifile, ofile, cfile, dfile string

	cliapp.Flags = []cli.Flag{
		&cli.BoolFlag{
			Name:        "debug",
			Aliases:     []string{"d"},
			Usage:       "debug mode",
			Required:    false,
			Destination: &debug,
		},
		&cli.StringFlag{
			Name:        "config",
			Aliases:     []string{"c"},
			Usage:       "config file",
			Required:    false,
			Destination: &cfile,
		},
		&cli.StringFlag{
			Name:        "input",
			Aliases:     []string{"i"},
			Usage:       "input file",
			Required:    true,
			Destination: &ifile,
		},
		&cli.StringFlag{
			Name:        "output",
			Aliases:     []string{"o"},
			Usage:       "output file",
			Required:    true,
			Destination: &ofile,
		},
		&cli.StringFlag{
			Name:        "document",
			Aliases:     []string{"m"},
			Usage:       "existing .docx file",
			Required:    false,
			Destination: &dfile,
		},
	}
	cliapp.Action = func(ctx *cli.Context) error {
		// return converter.Build(cfile, ifile, ofile, dfile)
		return docb.Build(cfile, ifile, ofile, dfile)
	}

	if err := cliapp.Run(os.Args); err != nil {
		logrus.Error(eris.ToString(err, debug))
	}
}
