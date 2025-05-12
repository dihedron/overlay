package main

import (
	"errors"
	"fmt"
	"image/color"
	"log/slog"
	"strconv"
	"strings"

	"github.com/jessevdk/go-flags"
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
	if len(value) > 0 && value[0] != '#' {
		return fmt.Errorf("invalid color string format")
	}

	switch len(value) {
	case 4: // #RGB

		if r, err := strconv.ParseUint(value[1:2], 16, 8); err != nil {
			return err
		} else {
			c.R = uint8(r)
		}

		if g, err := strconv.ParseUint(value[2:3], 16, 8); err != nil {
			return err
		} else {
			c.G = uint8(g)
		}

		if b, err := strconv.ParseUint(value[3:4], 16, 8); err != nil {
			return err
		} else {
			c.B = uint8(b)
		}

		c.A = 255

	case 5: // #RGBA

		if r, err := strconv.ParseUint(value[1:2], 16, 8); err != nil {
			return err
		} else {
			c.R = uint8(r)
		}

		if g, err := strconv.ParseUint(value[2:3], 16, 8); err != nil {
			return err
		} else {
			c.G = uint8(g)
		}

		if b, err := strconv.ParseUint(value[3:4], 16, 8); err != nil {
			return err
		} else {
			c.B = uint8(b)
		}

		if a, err := strconv.ParseUint(value[4:5], 16, 8); err != nil {
			return err
		} else {
			c.A = uint8(a)
		}

	case 7: // #RRGGBB

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

		c.A = 255

	case 9: // #RRGGBBAA

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

		if a, err := strconv.ParseUint(value[7:9], 16, 8); err != nil {
			return err
		} else {
			c.A = uint8(a)
		}
	default:
		return fmt.Errorf("invalid color string format ")
	}
	slog.Debug("parsed color", "red", c.R, "green", c.G, "blue", c.B, "alpha", c.A)
	return nil
}

func (c Color) MarshalFlag() (string, error) {
	return fmt.Sprintf("#%02X%02X%02X", c.R, c.G, c.B), nil
}

type Options struct {
	// Input is the name of the input file.
	Input flags.Filename `short:"i" long:"input" description:"The name of the input file or - for STDIN" optional:"true" default:"-"`
	// Output is the name of the output file.
	Output flags.Filename `short:"o" long:"output" description:"The name of the output file or - for STDOUT" optional:"true" default:"-"`
	// Format is the output format, if an output filename is not specified; it is used for chaining.
	Format string `short:"x" long:"format" description:"Format of the output image" optional:"true" choice:"jpeg" choice:"jpg" choice:"png" choice:"gif" choice:"bmp" default:"png"`
	// Text is the text to write as an overlay to the image.
	Text string `short:"t" long:"text" description:"The text to add as an overlay to the given image" optional:"true"`
	// Image is the image to superimpose as an overlay to the image.
	Image string `short:"y" long:"image" description:"The image to superimpose as an overlay to the given image" optional:"true"`
	// Point is the position in the image where the text will start.
	Point Point `short:"p" long:"point" description:"The coordinates where the text/image will be written, as an (x,y) point" optional:"true"`
	// Font is the font to use for writing to the image.
	Font flags.Filename `short:"f" long:"font" description:"The name of the font to be used for writing" optional:"true"`
	// Font is the font to use for writing to the image.
	Color Color `short:"c" long:"color" description:"The color of the font to be used for writing" optional:"true" default:"#000000"`
	// Size is the size of font to use for writing to the image.
	Size float64 `short:"s" long:"size" description:"The size of the font to be used for writing" optional:"true" default:"12"`
	// DPI is the image resolution in Dots Per Inch.
	DPI float64 `short:"d" long:"dpi" description:"The image resolution in DPI - Dots Per Inch" optional:"true" default:"72"`
}
