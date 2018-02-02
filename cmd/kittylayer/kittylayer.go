package main

import (
	"gopkg.in/urfave/cli.v1"
	"github.com/kittycash/kittiverse/src/tools/layer"
	"image/png"
	"os"
	"errors"
	"strings"
	"log"
	"image"
	"encoding/json"
)

func main() {
	app := cli.NewApp()
	app.Name = "kittylayer"
	app.Usage = "tools for modifying and managing kitty layers"
	app.Commands = cli.Commands{
		cli.Command{
			Name: "rotate",
			Usage: "rotates a layer",
			Flags: cli.FlagsByName{
				cli.StringFlag{
					Name: "source, s",
					Usage: "image source",
				},
				cli.StringFlag{
					Name: "destination, d",
					Usage: "image destination",
				},
				cli.Float64Flag{
					Name: "angle, a",
					Usage: "angle in clockwise radians of rotation",
				},
			},
			Action: cli.ActionFunc(rotate),
		},
		cli.Command{
			Name: "scale",
			Usage: "scales a layer by multiplication x & y",
			Flags: cli.FlagsByName{
				cli.StringFlag{
					Name: "source, s",
					Usage: "image source",
				},
				cli.StringFlag{
					Name: "destination, d",
					Usage: "image destination",
				},
				cli.Float64Flag{
					Name: "scaleX, x",
					Usage: "factor of scale in x direction",
				},
				cli.Float64Flag{
					Name: "scaleY, y",
					Usage: "factor of scale in y direction",
				},
			},
			Action: cli.ActionFunc(scale),
		},
		cli.Command{
			Name: "remove_whitespace",
			Usage: "Removes whitespace of image and spits into smaller image and config file",
			Flags: cli.FlagsByName{
				cli.StringFlag{
					Name: "source, s",
					Usage: "image source",
				},
				cli.StringFlag{
					Name: "destination, d",
					Usage: "image destination",
				},
			},
			Action: cli.ActionFunc(removeWhitespace),
		},
	}
	if e := app.Run(os.Args); e != nil {
		log.Println(e)
	}
}

func rotate(ctx *cli.Context) error {
	var (
		srcName = ctx.String("source")
		dstName = ctx.String("destination")
		angle   = ctx.Float64("angle")
	)

	src, e := openImage(srcName)
	if e != nil {
		return e
	}

	dst, e := layer.Rotate(src, angle)
	if e != nil {
		return e
	}

	return createImage(dstName, dst)
}

func scale(ctx *cli.Context) error {
	var (
		srcName = ctx.String("source")
		dstName = ctx.String("destination")
		scaleX  = ctx.Float64("scaleX")
		scaleY  = ctx.Float64("scaleY")
	)

	src, e := openImage(srcName)
	if e != nil {
		return e
	}

	dst, e := layer.Scale(src, scaleX, scaleY)
	if e != nil {
		return e
	}

	return createImage(dstName, dst)
}

func removeWhitespace(ctx *cli.Context) error {
	var (
		srcName = ctx.String("source")
		dstName = ctx.String("destination")
	)

	src, e := openImage(srcName)
	if e != nil {
		return e
	}

	dst, dstConfig := layer.RemoveWhitespace(src)

	return createImage(dstName, dst, func(fn string) error {
		cName := strings.TrimSuffix(fn, ".png") + ".json"

		cData, e := json.MarshalIndent(dstConfig, "", "    ")
		if e != nil {
			return e
		}

		cf, e := os.Create(cName)
		if e != nil {
			return e
		}

		_, e = cf.Write(cData)
		return e
	})
}

/*
	<<< HELPER FUNCTIONS >>>
*/

func openImage(srcName string) (image.Image, error) {
	sf, e := os.Open(srcName)
	if e != nil {
		return nil, errors.New("failed to open source: " + e.Error())
	}
	defer sf.Close()

	src, e := png.Decode(sf)
	if e != nil {
		return nil, errors.New("failed to decode source: " + e.Error())
	}
	return src, nil
}

type fnAction func(fn string) error

func createImage(dstName string, dst image.Image, fnActions ...fnAction) error {
	if strings.HasSuffix(dstName, ".png") == false {
		dstName += ".png"
	}

	df, e := os.Create(dstName)
	if e != nil {
		return errors.New("failed to create image: " + e.Error())
	}

	for _, action := range fnActions {
		if e := action(dstName); e != nil {
			return e
		}
	}

	return png.Encode(df, dst)
}