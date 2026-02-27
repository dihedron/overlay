package canvas

import (
	"image"
	"image/color"
	"image/draw"
	"log/slog"

	"github.com/dihedron/overlay/command/base"
)

// Canvas is the command that creates a new image with the given size and colour.
type Canvas struct {
	base.OutputCommand
	// Size is the size of the canvas.
	Size base.Point `short:"s" long:"size" description:"The size of the canvas, as an (width,height) pair" required:"true"`
	// Colour is the colour used to fill the canvas.
	Colour base.Colour `short:"c" long:"colour" description:"The colour used to fill the canvas" optional:"true" default:"#FFFFFF"`
}

// Execute is the real implementation of the Canvas command.
func (cmd *Canvas) Execute(args []string) error {
	slog.Debug("running canvas command")

	// create a blank image with the given size
	img := image.NewRGBA(image.Rect(0, 0, cmd.Size.X, cmd.Size.Y))

	// set the colour
	background := color.RGBA{R: cmd.Colour.R, G: cmd.Colour.G, B: cmd.Colour.B, A: cmd.Colour.A}

	// fill the image with the colour
	draw.Draw(img, img.Bounds(), &image.Uniform{background}, image.Point{}, draw.Src)

	// write the image to the output stream
	if err := cmd.WriteOutput(img); err != nil {
		slog.Error("error writing output stream", "name", cmd.Output, "error", err)
		return err
	}
	slog.Debug("image correctly encoded", "filename", cmd.Output, "format", cmd.Format)

	return nil
}
