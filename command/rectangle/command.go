package rectangle

import (
	"fmt"
	"log/slog"

	"github.com/dihedron/overlay/command/base"
	"github.com/gogpu/gg"
)

// Rectangle is the command that adds a rectangle as an overlay to an image.
type Rectangle struct {
	base.InputCommand
	base.OutputCommand
	// Point is the position in the image where the rectangle will start.
	Point base.PointF `short:"p" long:"point" description:"The coordinates where the rectangle will be written, as an (x,y) point" optional:"true"`
	// Size is the size of the rectangle to be written to the image.
	Size base.PointF `short:"s" long:"size" description:"The size of the rectangle to be written to the image, as an (width,height) point" optional:"true"`
	// Colour is the colour of the rectangle to be written to the image.
	Colour base.Colour `short:"c" long:"colour" description:"The colour of the rectangle to be written to the image" optional:"true" default:"#000000"`
	// Fill is whether the rectangle should be filled with the given colour.
	Fill bool `short:"f" long:"fill" description:"Whether the rectangle should be filled with the given colour, by default it is not" optional:"true"`
	// Stroke is the width of the rectangle stroke, when fill is false.
	Stroke float64 `short:"w" long:"stroke" description:"The width of the rectangle stroke, when fill is false" optional:"true" default:"1"`
	// Radius defines a rounded rectangle by rounding the corners of the rectangle
	Radius float64 `short:"r" long:"radius" description:"The radius of the rectangle corners" optional:"true" default:"0"`
}

// Execute is the real implementation of the Rectangle command.
func (cmd *Rectangle) Execute(args []string) error {
	slog.Debug("running rectangle command")

	// open the input and output streams
	underlay, err := cmd.ReadInput()
	if err != nil {
		slog.Error("error reading input stream", "name", cmd.Input, "error", err)
		return err
	}

	// create the device context with the underlay image
	dc := gg.NewContextForImage(underlay)
	defer dc.Close()

	// set the colour
	dc.SetRGBA(float64(cmd.Colour.R), float64(cmd.Colour.G), float64(cmd.Colour.B), float64(cmd.Colour.A))

	if cmd.Radius > 0 {
		slog.Debug("drawing rounded rectangle", "point", cmd.Point, "size", cmd.Size, "radius", cmd.Radius)
		// define the rounded rectangle area: (x1, y1) to (x2, y2)
		// rounded rectangle is defined by the top-left corner and the size
		dc.DrawRoundedRectangle(cmd.Point.X, cmd.Point.Y, cmd.Size.X, cmd.Size.Y, cmd.Radius)
	} else {
		slog.Debug("drawing rectangle", "point", cmd.Point, "size", cmd.Size)
		// define the rectangle area: (x1, y1) to (x2, y2)
		// rectangle is defined by the top-left corner and the size
		dc.DrawRectangle(cmd.Point.X, cmd.Point.Y, cmd.Size.X, cmd.Size.Y)
	}

	if cmd.Fill {
		slog.Debug("drawing rectangle as fill", "colour", cmd.Colour)
		dc.Fill()
	} else if cmd.Stroke > 0 {
		slog.Debug("drawing rectangle as stroke", "width", cmd.Stroke)
		dc.SetLineWidth(cmd.Stroke)
		dc.Stroke()
	} else {
		slog.Error("either --fill or --stroke must be specified")
		return fmt.Errorf("either --fill or --stroke must be specified")
	}

	slog.Debug("rectangle overlaid on the image", "point", cmd.Point, "size", cmd.Size, "colour", cmd.Colour)

	// write the image to the output stream
	img := dc.Image()
	if err := cmd.WriteOutput(img); err != nil {
		slog.Error("error writing output stream", "name", cmd.Output, "error", err)
		return err
	}
	slog.Debug("image correctly encoded", "filename", cmd.Output, "format", cmd.Format)

	return nil
}
