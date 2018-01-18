package main

import (
	"github.com/kittycash/kittiverse/src/incubator"
	"gopkg.in/urfave/cli.v1"
	"os"
	"log"
)

var (
	ImagesDir = "Kitties"
)

func main() {

	app := cli.NewApp()
	app.Name = "generatekitty"
	app.Usage = "generates kitties"
	app.Flags = cli.FlagsByName{
		cli.StringFlag{
			Name: "images-dir, i",
			Usage: "directory where we store all the modular image files",
			Value: ImagesDir,
			Destination: &ImagesDir,
		},
	}
	app.Action = run
	if e := app.Run(os.Args); e != nil {
		log.Println(e)
	}
}

func run(_ *cli.Context) error {
	if e := incubator.SetRootDir(ImagesDir); e != nil {
		return e
	}
	config := &incubator.KittyGenSpecs{
		Version: 0,
		DNA:     incubator.DNAGenSpecs{
			Group: 0,
			Color: 38,
			Pattern: 3,
			Body: 0,
			Brows: -1,
			Ears: 0,
			Eyes: 0,
			Head: 0,
			Nose: 2,
			Tail: 1,
		},
		Accessories: incubator.AccessoriesGenSpecs{
			Collar: &incubator.ItemGenSpecs{
				ID:    0,
				Color: 13,
			},
		},
	}
	_, e := incubator.GenerateKitty(config, true, "kitty.png")
	if e != nil {
		return e
	}
	return nil
}