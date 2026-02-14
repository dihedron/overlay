package canvas

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
	"golang.org/x/image/bmp"
)

// Canvas is the command that creates a new image with the given size and colour.
type Canvas struct {
	base.Command
	// Size is the size of the canvas.
	Size base.Point `short:"s" long:"size" description:"The size of the canvas, as an (width,height) pair" required:"true"`
	// Colour is the colour used to fill the canvas.
	Colour base.Colour `short:"c" long:"colour" description:"The colour used to fill the canvas" optional:"true" default:"#FFFFFF"`
}

// Execute is the real implementation of the Canvas command.
func (cmd *Canvas) Execute(args []string) error {
	slog.Debug("running canvas command")

	// open the output stream
	var output io.Writer

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

	slog.Debug("output stream ready")

	// create a blank 600x400 RGBA image
	img := image.NewRGBA(image.Rect(0, 0, cmd.Size.X, cmd.Size.Y))

	// set the colour
	background := color.RGBA{R: cmd.Colour.R, G: cmd.Colour.G, B: cmd.Colour.B, A: cmd.Colour.A}

	// fill the image with the colour
	draw.Draw(img, img.Bounds(), &image.Uniform{background}, image.Point{}, draw.Src)

	// encode the output image
	var err error
	slog.Debug("encoding output image", "name", cmd.Output, "format", cmd.Format)
	switch cmd.Format {
	case "jpg", "jpeg":
		slog.Debug("encoding output file as JPEG", "name", cmd.Output)
		if err = jpeg.Encode(output, img, nil); err != nil {
			slog.Error("error encoding output file", "name", cmd.Output, "error", err, "format", cmd.Format)
			os.Exit(1)
		}
	case "png":
		slog.Debug("encoding output file as PNG", "name", cmd.Output)
		if err = png.Encode(output, img); err != nil {
			slog.Error("error encoding output file", "name", cmd.Output, "error", err, "format", cmd.Format)
			os.Exit(1)
		}
	case "gif":
		slog.Debug("encoding output file as GIF", "name", cmd.Output)
		if err = gif.Encode(output, img, nil); err != nil {
			slog.Error("error encoding output file", "name", cmd.Output, "error", err, "format", cmd.Format)
			os.Exit(1)
		}
	case "bmp":
		slog.Debug("encoding output file as BMP", "name", cmd.Output)
		if err = bmp.Encode(output, img); err != nil {
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
