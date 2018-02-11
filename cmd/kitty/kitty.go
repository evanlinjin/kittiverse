package main

import (
	"encoding/json"
	"errors"
	"github.com/kittycash/kittiverse/src/kitty/generator"
	"github.com/kittycash/kittiverse/src/kitty/generator/container/v0"
	"github.com/kittycash/kittiverse/src/kitty/graphics"
	"gopkg.in/urfave/cli.v1"
	"image"
	"image/png"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"github.com/kittycash/kittiverse/src/kitty/genetics"
)

var app = cli.NewApp()

func init() {
	app.Name = "kitty"
	app.Usage = "managing the kittiverse"
	app.Commands = cli.Commands{
		cli.Command{
			Name:  "layer",
			Usage: "tools for modifying and managing kitty layers",
			Subcommands: cli.Commands{
				cli.Command{
					Name:  "rotate",
					Usage: "rotates a layer",
					Flags: cli.FlagsByName{
						cli.StringFlag{
							Name:  "source, s",
							Usage: "image source",
						},
						cli.StringFlag{
							Name:  "destination, d",
							Usage: "image destination",
						},
						cli.Float64Flag{
							Name:  "angle, a",
							Usage: "angle in clockwise radians of rotation",
						},
					},
					Action: func(ctx *cli.Context) error {
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
					},
				},
				cli.Command{
					Name:  "scale",
					Usage: "scales a layer by multiplication x & y",
					Flags: cli.FlagsByName{
						cli.StringFlag{
							Name:  "source, s",
							Usage: "image source",
						},
						cli.StringFlag{
							Name:  "destination, d",
							Usage: "image destination",
						},
						cli.Float64Flag{
							Name:  "scaleX, x",
							Usage: "factor of scale in x direction",
						},
						cli.Float64Flag{
							Name:  "scaleY, y",
							Usage: "factor of scale in y direction",
						},
					},
					Action: func(ctx *cli.Context) error {
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
					},
				},
				cli.Command{
					Name:  "remove_whitespace",
					Usage: "removes whitespace of image and spits into smaller image and config file",
					Flags: cli.FlagsByName{
						cli.StringFlag{
							Name:  "source, s",
							Usage: "image source",
						},
						cli.StringFlag{
							Name:  "destination, d",
							Usage: "image destination",
						},
					},
					Action: func(ctx *cli.Context) error {
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
					},
				},
				cli.Command{
					Name:  "include_whitespace",
					Usage: "adds whitespace to an image with associated config file",
					Flags: cli.FlagsByName{
						cli.StringFlag{
							Name:  "source, s",
							Usage: "image source",
						},
						cli.StringFlag{
							Name:  "destination, d",
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
					Action: func(ctx *cli.Context) error {
						var (
							srcName   = ctx.String("source")
							dstName   = ctx.String("destination")
							dstWidth  = ctx.Int("width")
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
					},
				},
			},
		},
		cli.Command{
			Name:  "admin",
			Usage: "tools for compiling and managing kitty generation files",
			Subcommands: cli.Commands{
				cli.Command{
					Name:  "compile",
					Usage: "compiles a kitty generation file",
					Flags: cli.FlagsByName{
						cli.StringFlag{
							Name:  "dir, d",
							Usage: "path of loose files to compile from",
							Value: "kitty_layers",
						},
						cli.StringFlag{
							Name:  "output, o",
							Usage: "path of output file",
							Value: "file.kcg",
						},
					},
					Action: func(ctx *cli.Context) error {
						gen := generator.NewInstance(
							v0.NewImagesContainer(),
							v0.NewLayersContainer(),
						)
						if e := gen.Compile(ctx.String("dir")); e != nil {
							return e
						}
						log.Println("[ALLELE_RANGES]", gen.GetAlleleRanges().String(true))
						f, e := os.Create(ctx.String("output"))
						if e != nil {
							return e
						}
						return gen.Export(f)
					},
				},
			},
		},
		cli.Command{
			Name: "dna",
			Usage: "tools for generating and managing kitty DNA",
			Subcommands: cli.Commands{
				cli.Command{
					Name: "random",
					Usage: "generates a random DNA",
					Flags: cli.FlagsByName{
						cli.StringFlag{
							Name: "file, f",
							Usage: "path of '.kcg' file to use",
							Value: "file.kcg",
						},
					},
					Action: func(ctx *cli.Context) error {
						gen := generator.NewInstance(
							v0.NewImagesContainer(),
							v0.NewLayersContainer(),
						)
						f, e := os.Open(ctx.String("file"))
						if e != nil {
							return e
						}
						s, _ := f.Stat()
						if e := gen.Import(f, int(s.Size())); e != nil {
							return e
						}
						if e := f.Close(); e != nil {
							return e
						}
						out, e := json.MarshalIndent(struct {
							DNA string `json:"dna"`
						}{
							DNA: gen.GetAlleleRanges().RandomDNA().Hex(),
						}, "", "    ")
						if e != nil {
							return e
						} else {
							log.Println(string(out))
						}
						return nil
					},
				},
				cli.Command{
					Name: "image",
					Usage: "generates a kitty image from DNA",
					Flags: cli.FlagsByName{
						cli.StringFlag{
							Name: "dna, d",
							Usage: "hex representation of DNA",
							Value: genetics.DNA{}.Hex(),
						},
						cli.StringFlag{
							Name: "file, f",
							Usage: "path of '.kcg' file to use",
							Value: "file.kcg",
						},
						cli.StringFlag{
							Name: "output, o",
							Usage: "path of the output file of the kitty generated from DNA",
							Value: "kitty.png",
						},
					},
					Action: func(ctx *cli.Context) error {
						gen := generator.NewInstance(
							v0.NewImagesContainer(),
							v0.NewLayersContainer(),
						)
						f, e := os.Open(ctx.String("file"))
						if e != nil {
							return e
						}
						s, _ := f.Stat()
						if e := gen.Import(f, int(s.Size())); e != nil {
							return e
						}
						if e := f.Close(); e != nil {
							return e
						}
						f, e = os.Create(ctx.String("output"))
						if e != nil {
							return e
						}
						defer f.Close()
						dna, e := genetics.NewDNAFromHex(ctx.String("dna"))
						if e != nil {
							return e
						}
						img, e := gen.GenerateKitty(dna)
						if e != nil {
							return e
						}
						return png.Encode(f, img)
					},
				},
			},
		},
	}

}

func main() {
	if e := app.Run(os.Args); e != nil {
		log.Println(e)
	}
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
