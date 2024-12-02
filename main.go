package main

import (
	"errors"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/gif"
	"image/jpeg"
	"image/png"
	"log/slog"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/jessevdk/go-flags"
	"golang.org/x/image/bmp"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
	"golang.org/x/image/math/fixed"
)

type Point struct {
	X, Y int
}

func (p *Point) UnmarshalFlag(value string) error {
	parts := strings.Split(value, ",")

	if len(parts) != 2 {
		return errors.New("invalid format: expected two numbers separated by a ,")
	}
	x, err := strconv.ParseInt(parts[0], 10, 32)
	if err != nil {
		return err
	}
	y, err := strconv.ParseInt(parts[1], 10, 32)
	if err != nil {
		return err
	}
	p.X = int(x)
	p.Y = int(y)
	return nil
}

func (p Point) MarshalFlag() (string, error) {
	return fmt.Sprintf("%d,%d", p.X, p.Y), nil
}

type Color color.RGBA

func (c *Color) UnmarshalFlag(value string) error {
	if len(value) != 7 || value[0] != '#' {
		c.R = 0
		c.G = 0
		c.B = 0
		c.A = 0
		return fmt.Errorf("invalid color string format")
	}

	//c.A = 127
	c.A = 255

	if r, err := strconv.ParseUint(value[1:3], 16, 8); err != nil {
		return err
	} else {
		c.R = uint8(r)
	}

	if g, err := strconv.ParseUint(value[3:5], 16, 8); err != nil {
		return err
	} else {
		c.G = uint8(g)
	}

	if b, err := strconv.ParseUint(value[5:7], 16, 8); err != nil {
		return err
	} else {
		c.B = uint8(b)
	}

	// if a, err := strconv.ParseUint(value[7:], 16, 8); err != nil {
	// 	return err
	// } else {
	// 	c.A = uint8(a)
	// }
	return nil
}

func (c Color) MarshalFlag() (string, error) {
	return fmt.Sprintf("#%02X%02X%02X", c.R, c.G, c.B), nil
}

type Options struct {
	// Input is the name of the input file.
	Input flags.Filename `short:"i" long:"input" description:"Name of the input file." required:"true"`
	// Output is the name of the output file.
	Output flags.Filename `short:"o" long:"output" description:"Name of the output file." required:"true"`
	// Text is the text to write as an overlay to the image.
	Text string `short:"t" long:"text" description:"The text to add as an overly to the given image." default:"hallo, world!"`
	// Point is the position in the image where the text will start.
	Point Point `short:"p" long:"point" description:"The coordinates where the text will be written, as an (x,y) point." required:"true"`
	// Font is the font to use for writing to the image.
	Font flags.Filename `short:"f" long:"font" description:"The name of the font to be used for writing." required:"true"`
	// Font is the font to use for writing to the image.
	Color Color `short:"c" long:"color" description:"The color of the font to be used for writing." optional:"true" default:"#000000"`
	// Size is the size of font to use for writing to the image.
	Size float64 `short:"s" long:"size" description:"The size of the font to be used for writing." required:"true"`
	// DPI is the image resolution in Dots Per Inch.
	DPI float64 `short:"d" long:"dpi" description:"The image resolution in DPI (Dots Per Inch)." optional:"true" default:"72"`
}

func main() {

	var options Options

	var parser = flags.NewParser(&options, flags.Default)

	if _, err := parser.Parse(); err != nil {
		switch flagsErr := err.(type) {
		case flags.ErrorType:
			if flagsErr == flags.ErrHelp {
				os.Exit(0)
			}
			os.Exit(1)
		default:
			os.Exit(1)
		}
	}

	slog.Debug("running with options", "options", options)

	slog.Debug("opening input file", "filename", options.Input)

	// open the image file
	imgFile, err := os.Open(string(options.Input))
	if err != nil {
		slog.Error("error opening input file", "name", options.Input, "error", err)
		os.Exit(1)
	}
	defer imgFile.Close()

	slog.Debug("input file open", "filename", options.Input)

	// decode the image
	img, _, err := image.Decode(imgFile)
	if err != nil {
		slog.Error("error decoding input file", "name", options.Input, "error", err)
		os.Exit(1)
	}

	slog.Debug("image decoded", "filename", options.Input, "width", img.Bounds().Dx(), "height", img.Bounds().Dy())

	// create a new image with the same dimensions as the original
	dst := image.NewRGBA(img.Bounds())
	draw.Draw(dst, dst.Bounds(), img, image.Point{0, 0}, draw.Src)

	slog.Debug("image copied to destination context", "width", dst.Bounds().Dx(), "height", dst.Bounds().Dy())

	// read the font data
	fontData, err := os.ReadFile(string(options.Font))
	if err != nil {
		slog.Error("error reading font file", "name", options.Font, "error", err)
		os.Exit(1)
	}

	slog.Debug("font data read", "filename", options.Font)

	// parse the font data into a font
	f, err := opentype.Parse(fontData)
	if err != nil {
		slog.Error("error parsing font data", "name", options.Font, "error", err)
		os.Exit(1)
	}

	slog.Debug("font parsed")

	// create the font face
	fontFace, err := opentype.NewFace(f, &opentype.FaceOptions{
		Size:    float64(options.Size),
		DPI:     options.DPI,
		Hinting: font.HintingNone,
	})
	if err != nil {
		slog.Error("error creating font face", "name", options.Font, "error", err)
		os.Exit(1)
	}

	point := fixed.Point26_6{
		X: fixed.I(options.Point.X),
		Y: fixed.I(options.Point.Y),
	}

	d := &font.Drawer{
		Dst:  dst,
		Src:  image.NewUniform(color.RGBA(options.Color)),
		Face: fontFace,
		Dot:  point,
	}
	d.DrawString(options.Text)

	// encode the image as JPEG
	out, err := os.Create(string(options.Output))
	if err != nil {
		slog.Error("error opening output file", "name", options.Output, "error", err)
		os.Exit(1)

	}
	defer out.Close()

	switch strings.ToLower(filepath.Ext(string(options.Input))) {
	case ".jpg", ".jpeg":
		slog.Debug("encoding output file as JPEG", "name", options.Output)
		if err = jpeg.Encode(out, dst, nil); err != nil {
			slog.Error("error encoding output file", "name", options.Output, "error", err)
			os.Exit(1)
		}
	case ".png":
		slog.Debug("encoding output file as PNG", "name", options.Output)
		if err = png.Encode(out, dst); err != nil {
			slog.Error("error encoding output file", "name", options.Output, "error", err)
			os.Exit(1)
		}
	case ".gif":
		slog.Debug("encoding output file as GIF", "name", options.Output)
		if err = gif.Encode(out, dst, nil); err != nil {
			slog.Error("error encoding output file", "name", options.Output, "error", err)
			os.Exit(1)
		}
	case ".bmp":
		slog.Debug("encoding output file as BMP", "name", options.Output)
		if err = bmp.Encode(out, dst); err != nil {
			slog.Error("error encoding output file", "name", options.Output, "error", err)
			os.Exit(1)
		}
	default:
		fmt.Fprintf(os.Stderr, "Unsupported input file type: %s\n", filepath.Ext(string(options.Input)))
		slog.Error("unsupported input image type", "name", options.Input)
		os.Exit(1)
	}

	slog.Debug("image correctly encoded", "filename", options.Output)
}
