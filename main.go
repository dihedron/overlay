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
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/dihedron/overlay/version"
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
	Input flags.Filename `short:"i" long:"input" description:"The Name of the input file (default: STDIN)." optional:"true" default:"-"`
	// Output is the name of the output file.
	Output flags.Filename `short:"o" long:"output" description:"Name of the output file (default: STDOUT)." optional:"true" default:"-"`
	// Format is the output format, if an output filename is not specified; it is used for chaining.
	Format string `short:"F" long:"format" description:"Format of the output image." optional:"true" choice:"jpeg" choice:"jpg" choice:"png" choice:"gif" choice:"bmp" default:"png"`
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

	if len(os.Args) == 2 && os.Args[1] == "--version" {
		version.Print(os.Stdout)
		os.Exit(0)
	}

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
			fmt.Fprintf(os.Stderr, "%v\n", err)
			os.Exit(1)
		}
	}

	slog.Debug("start running...", "options", options)

	var (
		input  io.Reader
		output io.Writer
	)
	if options.Input == "-" {
		slog.Debug("getting image from STDIN")
		input = os.Stdin
	} else {
		// open the image file
		slog.Debug("reading input from file", "name", options.Input)
		var err error
		if input, err = os.Open(string(options.Input)); err != nil {
			slog.Error("error opening input file", "name", options.Input, "error", err)
			os.Exit(1)
		}
		if input, ok := input.(io.ReadCloser); ok {
			slog.Debug("input needs to be closed at application shutdown", "name", options.Input)
			defer input.Close()
		}
	}

	if options.Output == "-" {
		slog.Debug("writing image to STDOUT", "format", options.Format)
		output = os.Stdout
	} else {
		switch strings.ToLower(filepath.Ext(string(options.Output))) {
		case ".jpg", ".jpeg":
			options.Format = "jpg"
		case ".png":
			options.Format = "png"
		case ".gif":
			options.Format = "gif"
		case ".bmp":
			options.Format = "bmp"
		default:
			fmt.Fprintf(os.Stderr, "Unsupported output file type: %s\n", filepath.Ext(string(options.Output)))
			slog.Error("unsupported output image type", "name", options.Output)
			os.Exit(1)
		}
		slog.Debug("writing output to file", "name", options.Output, "format", options.Format)

		// open the output file
		var err error
		if output, err = os.Create(string(options.Output)); err != nil {
			slog.Error("error opening output file", "name", options.Output, "error", err)
			os.Exit(1)
		}
		if output, ok := output.(io.WriteCloser); ok {
			slog.Debug("output needs to be closed at application shutdown", "name", options.Output)
			defer output.Close()
		}
	}

	slog.Debug("streams ready")

	// decode the image
	img, _, err := image.Decode(input)
	if err != nil {
		slog.Error("error decoding input data", "name", options.Input, "error", err)
		os.Exit(1)
	}

	slog.Debug("image decoded", "name", options.Input, "width", img.Bounds().Dx(), "height", img.Bounds().Dy())

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

	switch options.Format {
	case "jpg", "jpeg":
		slog.Debug("encoding output file as JPEG", "name", options.Output)
		if err = jpeg.Encode(output, dst, nil); err != nil {
			slog.Error("error encoding output file", "name", options.Output, "error", err, "format", options.Format)
			os.Exit(1)
		}
	case "png":
		slog.Debug("encoding output file as PNG", "name", options.Output)
		if err = png.Encode(output, dst); err != nil {
			slog.Error("error encoding output file", "name", options.Output, "error", err, "format", options.Format)
			os.Exit(1)
		}
	case "gif":
		slog.Debug("encoding output file as GIF", "name", options.Output)
		if err = gif.Encode(output, dst, nil); err != nil {
			slog.Error("error encoding output file", "name", options.Output, "error", err, "format", options.Format)
			os.Exit(1)
		}
	case "bmp":
		slog.Debug("encoding output file as BMP", "name", options.Output)
		if err = bmp.Encode(output, dst); err != nil {
			slog.Error("error encoding output file", "name", options.Output, "error", err, "format", options.Format)
			os.Exit(1)
		}
	default:
		fmt.Fprintf(os.Stderr, "Unsupported output format: %s\n", options.Format)
		slog.Error("unsupported output format", "name", options.Output, "format", options.Format)
		os.Exit(1)
	}
	slog.Debug("image correctly encoded", "filename", options.Output, "format", options.Format)
}
