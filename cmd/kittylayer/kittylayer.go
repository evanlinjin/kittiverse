package main

import (
	"gopkg.in/urfave/cli.v1"
	"github.com/kittycash/kittiverse/src/imager/layer"
	"image/png"
	"os"
	"errors"
	"strings"
	"log"
	"image"
	"encoding/json"
	"io/ioutil"
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
			Usage: "removes whitespace of image and spits into smaller image and config file",
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
		cli.Command{
			Name: "include_whitespace",
			Usage: "adds whitespace to an image with associated config file",
			Flags: cli.FlagsByName{
				cli.StringFlag{
					Name: "source, s",
					Usage: "image source",
				},
				cli.StringFlag{
					Name: "destination, d",
					Usage: "image destination",
				},
				cli.IntFlag{
					Name:  "width, x",
					Usage: "width of destination image in pixels",
					Value: 1200,
				},
				cli.IntFlag{
					Name:  "height, y",
					Usage: "height of destination image in pixels",
					Value: 1200,
				},
			},
			Action: cli.ActionFunc(includeWhitespace),
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

func includeWhitespace(ctx *cli.Context) error {
	var (
		srcName = ctx.String("source")
		dstName = ctx.String("destination")
		dstWidth = ctx.Int("width")
		dstHeight = ctx.Int("height")
		placement = new(layer.Placement)
	)

	src, e := openImage(srcName, func(fn string) error {
		fn = strings.TrimSuffix(fn, ".png") + ".json"

		data, e := ioutil.ReadFile(fn)
		if e != nil {
			return e
		}

		return json.Unmarshal(data, placement)
	})
	if e != nil {
		return e
	}

	dstBounds := image.Rect(0, 0, dstWidth, dstHeight)
	dst, e := layer.IncludeWhitespace(src, dstBounds, placement)
	if e != nil {
		return e
	}

	return createImage(dstName, dst)
}

/*
	<<< HELPER FUNCTIONS >>>
*/

func openImage(srcName string, fnActions ...fnAction) (image.Image, error) {
	for _, action := range fnActions {
		if e := action(srcName); e != nil {
			return nil, e
		}
	}

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