package main

import (
	"github.com/urfave/cli/v2"
	"github.com/wii-tools/archon/cmd/arc"
	"github.com/wii-tools/archon/cmd/wad"
	"log"
	"os"
)

func main() {
	app := &cli.App{
		Usage: "your swiss army knife for the Wii",
		UsageText: "Archon is a multi-purpose tool designed to assist with the manipulation of various formats. " +
			"Commands below contain subcommands. Run them individually to obtain their help.",
		Commands: []*cli.Command{
			{
				Name:  "wad",
				Usage: "WAD-related operations",
				Subcommands: []*cli.Command{
					{
						Name:   "inspect",
						Action: wad.Inspect,
						Flags: []cli.Flag{
							&cli.StringFlag{Name: "in", Usage: "path to the WAD file", TakesFile: true, Required: true},
						},
					},
					{
						Name:   "unpack",
						Action: wad.Unpack,
						Flags: []cli.Flag{
							&cli.StringFlag{Name: "in", Usage: "path to the WAD file", TakesFile: true, Required: true},
							&cli.StringFlag{Name: "out", Usage: "directory to extract contents to", Required: true},
							&cli.BoolFlag{Name: "no-decrypt", Usage: "if specified, does not decrypted WAD contents"},
							&cli.BoolFlag{Name: "id", Usage: "name contents by ID, instead of index"},
						},
					},
					{
						Name:   "pack",
						Action: wad.Pack,
						Flags: []cli.Flag{
							&cli.StringFlag{Name: "in", Usage: "directory to extract contents to", Required: true},
							&cli.StringFlag{Name: "out", Usage: "path to output the WAD", Required: true},
							&cli.BoolFlag{Name: "id", Usage: "read contents by ID, instead of index"},
						},
					},
				},
			},
			{
				Name: "u8",
				Usage: "U8 archive related operations",
				Subcommands: []*cli.Command{
					{
						Name: "inspect",
						Action: arc.Inspect,
						Flags: []cli.Flag{
							&cli.StringFlag{Name: "in", Usage: "path to the U8 archive", TakesFile: true, Required: true},
						},
					},
					{
						Name: "extract",
						Action: arc.Extract,
						Flags: []cli.Flag{
							&cli.StringFlag{Name: "in", Usage: "path to the U8 archive", TakesFile: true, Required: true},
							&cli.StringFlag{Name: "out", Usage: "path to where you want to extract", Required: true},
						},
					},
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
