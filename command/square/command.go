package square

import (
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
	"strings"

	"github.com/dihedron/overlay/command/base"
	"github.com/jessevdk/go-flags"
	"golang.org/x/image/bmp"
)

// Square is the command that adds text as an overlay to an image.
type Square struct {
	base.Command
	// Input is the name of the input file.
	Input flags.Filename `short:"i" long:"input" description:"The name of the input file or - for STDIN" optional:"true" default:"-"`
	// Point is the position in the image where the square will start.
	Point base.Point `short:"p" long:"point" description:"The coordinates where the square will be written, as an (x,y) point" optional:"true"`
	// Size is the size of the square to be written to the image.
	Size base.Point `short:"s" long:"size" description:"The size of the square to be written to the image, as an (width,height) point" optional:"true"`
	// Colour is the colour of the square to be written to the image.
	Colour base.Colour `short:"c" long:"colour" description:"The colour of the square to be written to the image" optional:"true" default:"#000000"`
	// Fill is whether the square should be filled with the given colour.
	Fill bool `short:"f" long:"fill" description:"Whether the square should be filled with the given colour" optional:"true" default:"false"`
}

// Execute is the real implementation of the Square command.
func (cmd *Square) Execute(args []string) error {
	slog.Debug("running text command")

	// open the input and output streams
	var (
		input  io.Reader
		output io.Writer
	)
	if cmd.Input == "-" {
		slog.Debug("getting image from STDIN")
		input = os.Stdin
	} else {
		// open the underlay image file
		slog.Debug("reading input from file", "name", cmd.Input)
		var err error
		if input, err = os.Open(string(cmd.Input)); err != nil {
			slog.Error("error opening input file", "name", cmd.Input, "error", err)
			os.Exit(1)
		}
		if input, ok := input.(io.ReadCloser); ok {
			slog.Debug("input needs to be closed at application shutdown", "name", cmd.Input)
			defer input.Close()
		}
	}

	if cmd.Output == "-" {
		slog.Debug("writing image to STDOUT", "format", cmd.Format)
		output = os.Stdout
	} else {
		switch strings.ToLower(filepath.Ext(string(cmd.Output))) {
		case ".jpg", ".jpeg":
			cmd.Format = "jpg"
		case ".png":
			cmd.Format = "png"
		case ".gif":
			cmd.Format = "gif"
		case ".bmp":
			cmd.Format = "bmp"
		default:
			fmt.Fprintf(os.Stderr, "Unsupported output file type: %s\n", filepath.Ext(string(cmd.Output)))
			slog.Error("unsupported output image type", "name", cmd.Output)
			os.Exit(1)
		}
		slog.Debug("writing output to file", "name", cmd.Output, "format", cmd.Format)

		// open the output file
		var err error
		if output, err = os.Create(string(cmd.Output)); err != nil {
			slog.Error("error opening output file", "name", cmd.Output, "error", err)
			os.Exit(1)
		}
		if output, ok := output.(io.WriteCloser); ok {
			slog.Debug("output needs to be closed at application shutdown", "name", cmd.Output)
			defer output.Close()
		}
	}

	slog.Debug("streams ready")

	// decode the underlay image
	underlay, _, err := image.Decode(input)
	if err != nil {
		slog.Error("error decoding input data for underlay image", "name", cmd.Input, "error", err)
		os.Exit(1)
	}
	slog.Debug("underlay image decoded", "name", cmd.Input, "width", underlay.Bounds().Dx(), "height", underlay.Bounds().Dy())

	// create a new image with the same dimensions as the original
	dst := image.NewRGBA(underlay.Bounds())
	draw.Draw(dst, dst.Bounds(), underlay, image.Point{0, 0}, draw.Src)

	slog.Debug("image copied to destination context", "width", dst.Bounds().Dx(), "height", dst.Bounds().Dy())

	slog.Debug("overlaying square on the image", "point", cmd.Point, "size", cmd.Size, "colour", cmd.Colour)

	// set the colour
	background := color.RGBA{R: cmd.Colour.R, G: cmd.Colour.G, B: cmd.Colour.B, A: cmd.Colour.A}

	// define the rectangle area: (x1, y1) to (x2, y2)
	squareRect := image.Rect(cmd.Point.X, cmd.Point.Y, cmd.Point.X+cmd.Size.X, cmd.Point.Y+cmd.Size.Y)

	// draw the square on top of the image
	if cmd.Fill {
		// draw the filled square
		draw.Draw(dst, squareRect, &image.Uniform{background}, image.Point{}, draw.Src)
	} else {
		// draw the outline of the square
		draw.Draw(dst, squareRect, &image.Uniform{background}, image.Point{}, draw.Over)
	}

	slog.Debug("square overlaid on the image", "point", cmd.Point, "size", cmd.Size, "colour", cmd.Colour)

	// encode the output image
	slog.Debug("encoding output image", "name", cmd.Output, "format", cmd.Format)
	switch cmd.Format {
	case "jpg", "jpeg":
		slog.Debug("encoding output file as JPEG", "name", cmd.Output)
		if err = jpeg.Encode(output, dst, nil); err != nil {
			slog.Error("error encoding output file", "name", cmd.Output, "error", err, "format", cmd.Format)
			os.Exit(1)
		}
	case "png":
		slog.Debug("encoding output file as PNG", "name", cmd.Output)
		if err = png.Encode(output, dst); err != nil {
			slog.Error("error encoding output file", "name", cmd.Output, "error", err, "format", cmd.Format)
			os.Exit(1)
		}
	case "gif":
		slog.Debug("encoding output file as GIF", "name", cmd.Output)
		if err = gif.Encode(output, dst, nil); err != nil {
			slog.Error("error encoding output file", "name", cmd.Output, "error", err, "format", cmd.Format)
			os.Exit(1)
		}
	case "bmp":
		slog.Debug("encoding output file as BMP", "name", cmd.Output)
		if err = bmp.Encode(output, dst); err != nil {
			slog.Error("error encoding output file", "name", cmd.Output, "error", err, "format", cmd.Format)
			os.Exit(1)
		}
	default:
		fmt.Fprintf(os.Stderr, "Unsupported output format: %s\n", cmd.Format)
		slog.Error("unsupported output format", "name", cmd.Output, "format", cmd.Format)
		os.Exit(1)
	}
	slog.Debug("image correctly encoded", "filename", cmd.Output, "format", cmd.Format)

	return nil
}
