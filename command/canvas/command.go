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
	var (
		output io.Writer
		err    error
	)

	slog.Debug("running canvas command")

	// open the output stream
	if output, err = cmd.OutputStream(); err != nil {
		slog.Error("error opening output stream", "name", cmd.Output, "error", err)
		return err
	}

	// ensure the output stream is closed at application shutdown
	if output, ok := output.(io.WriteCloser); ok {
		slog.Debug("output needs to be closed at application shutdown", "name", cmd.Output)
		defer output.Close()
	}

	// create a blank 600x400 RGBA image
	img := image.NewRGBA(image.Rect(0, 0, cmd.Size.X, cmd.Size.Y))

	// set the colour
	background := color.RGBA{R: cmd.Colour.R, G: cmd.Colour.G, B: cmd.Colour.B, A: cmd.Colour.A}

	// fill the image with the colour
	draw.Draw(img, img.Bounds(), &image.Uniform{background}, image.Point{}, draw.Src)

	// encode the output image
	slog.Debug("encoding output image", "name", cmd.Output, "format", cmd.Format)
	switch cmd.Format {
	case "jpg", "jpeg":
		slog.Debug("encoding output file as JPEG", "name", cmd.Output)
		if err = jpeg.Encode(output, img, nil); err != nil {
			slog.Error("error encoding output file", "name", cmd.Output, "error", err, "format", cmd.Format)
			return err
		}
	case "png":
		slog.Debug("encoding output file as PNG", "name", cmd.Output)
		if err = png.Encode(output, img); err != nil {
			slog.Error("error encoding output file", "name", cmd.Output, "error", err, "format", cmd.Format)
			return err
		}
	case "gif":
		slog.Debug("encoding output file as GIF", "name", cmd.Output)
		if err = gif.Encode(output, img, nil); err != nil {
			slog.Error("error encoding output file", "name", cmd.Output, "error", err, "format", cmd.Format)
			return err
		}
	case "bmp":
		slog.Debug("encoding output file as BMP", "name", cmd.Output)
		if err = bmp.Encode(output, img); err != nil {
			slog.Error("error encoding output file", "name", cmd.Output, "error", err, "format", cmd.Format)
			return err
		}
	default:
		fmt.Fprintf(os.Stderr, "Unsupported output format: %s\n", cmd.Format)
		slog.Error("unsupported output format", "name", cmd.Output, "format", cmd.Format)
		return fmt.Errorf("unsupported output format: %s", cmd.Format)
	}
	slog.Debug("image correctly encoded", "filename", cmd.Output, "format", cmd.Format)

	return nil
}
