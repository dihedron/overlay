package circle

import (
	"fmt"
	"log/slog"

	"github.com/dihedron/overlay/command/base"
	"github.com/gogpu/gg"
)

// Circle is the command that adds a circle as an overlay to an image.
type Circle struct {
	base.InputCommand
	base.OutputCommand
	// Point is the position in the image where the circle will start.
	Point base.Point `short:"p" long:"point" description:"The coordinates where the circle will be written, as an (x,y) point" optional:"true"`
	// // Size is the size of the square to be written to the image.
	// Size base.Point `short:"s" long:"size" description:"The size of the square to be written to the image, as an (width,height) point" optional:"true"`
	// Colour is the colour of the circle to be written to the image.
	Colour base.Colour `short:"c" long:"colour" description:"The colour of the circle to be written to the image" optional:"true" default:"#000000"`
	// Fill is whether the circle should be filled with the given colour.
	Fill bool `short:"f" long:"fill" description:"Whether the circle should be filled with the given colour, by default it is not" optional:"true"`
	// Stroke is the width of the circle stroke, when fill is false.
	Stroke float64 `short:"w" long:"stroke" description:"The width of the circle stroke, when fill is false" optional:"true" default:"1"`
	// Radius defines the radius of the circle.
	Radius float64 `short:"r" long:"radius" description:"The radius of the circle" optional:"true" default:"10"`
}

// Execute is the real implementation of the Circle command.
func (cmd *Circle) Execute(args []string) error {
	slog.Debug("running circle command")

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

	slog.Debug("drawing circle", "point", cmd.Point, "radius", cmd.Radius)
	dc.DrawCircle(cmd.Point.X, cmd.Point.Y, cmd.Radius)

	if cmd.Fill {
		slog.Debug("drawing circle as fill", "colour", cmd.Colour)
		dc.Fill()
	} else if cmd.Stroke > 0 {
		slog.Debug("drawing circle as stroke", "width", cmd.Stroke)
		dc.SetLineWidth(cmd.Stroke)
		dc.Stroke()
	} else {
		slog.Error("either --fill or --stroke must be specified")
		return fmt.Errorf("either --fill or --stroke must be specified")
	}

	// if cmd.Radius > 0 {
	// 	slog.Debug("drawing rounded rectangle", "point", cmd.Point, "size", cmd.Size, "radius", cmd.Radius)
	// 	// define the rounded rectangle area: (x1, y1) to (x2, y2)
	// 	// rounded rectangle is defined by the top-left corner and the size
	// 	dc.DrawRoundedRectangle(float64(cmd.Point.X), float64(cmd.Point.Y), float64(cmd.Size.X), float64(cmd.Size.Y), float64(cmd.Radius))
	// } else {
	// 	slog.Debug("drawing rectangle", "point", cmd.Point, "size", cmd.Size)
	// 	// define the rectangle area: (x1, y1) to (x2, y2)
	// 	// rectangle is defined by the top-left corner and the size
	// 	dc.DrawRectangle(float64(cmd.Point.X), float64(cmd.Point.Y), float64(cmd.Size.X), float64(cmd.Size.Y))
	// }

	// if cmd.Fill {
	// 	slog.Debug("drawing rectangle as fill", "colour", cmd.Colour)
	// 	dc.Fill()
	// } else if cmd.Stroke > 0 {
	// 	slog.Debug("drawing rectangle as stroke", "width", cmd.Stroke)
	// 	dc.SetLineWidth(cmd.Stroke)
	// 	dc.Stroke()
	// } else {
	// 	slog.Error("either --fill or --stroke must be specified")
	// 	return fmt.Errorf("either --fill or --stroke must be specified")
	// }

	slog.Debug("circle overlaid on the image", "point", cmd.Point, "radius", cmd.Radius, "colour", cmd.Colour)

	// write the image to the output stream
	img := dc.Image()
	if err := cmd.WriteOutput(img); err != nil {
		slog.Error("error writing output stream", "name", cmd.Output, "error", err)
		return err
	}
	slog.Debug("image correctly encoded", "filename", cmd.Output, "format", cmd.Format)

	return nil
}
