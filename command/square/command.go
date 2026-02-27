package square

import (
	"image"
	"image/color"
	"image/draw"
	"log/slog"

	"github.com/dihedron/overlay/command/base"
)

// Square is the command that adds text as an overlay to an image.
type Square struct {
	base.InputCommand
	base.OutputCommand
	// Point is the position in the image where the square will start.
	Point base.Point `short:"p" long:"point" description:"The coordinates where the square will be written, as an (x,y) point" optional:"true"`
	// Size is the size of the square to be written to the image.
	Size base.Point `short:"s" long:"size" description:"The size of the square to be written to the image, as an (width,height) point" optional:"true"`
	// Colour is the colour of the square to be written to the image.
	Colour base.Colour `short:"c" long:"colour" description:"The colour of the square to be written to the image" optional:"true" default:"#000000"`
	// Fill is whether the square should be filled with the given colour.
	Fill bool `short:"f" long:"fill" description:"Whether the square should be filled with the given colour" optional:"true"`
}

// Execute is the real implementation of the Square command.
func (cmd *Square) Execute(args []string) error {
	slog.Debug("running text command")

	// open the input and output streams
	underlay, err := cmd.ReadInput()
	if err != nil {
		slog.Error("error reading input stream", "name", cmd.Input, "error", err)
		return err
	}

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
	if err := cmd.WriteOutput(dst); err != nil {
		slog.Error("error writing output stream", "name", cmd.Output, "error", err)
		return err
	}
	slog.Debug("image correctly encoded", "filename", cmd.Output, "format", cmd.Format)

	return nil
}
