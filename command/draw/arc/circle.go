package arc

import (
	"fmt"
	"log/slog"
	"math"

	"github.com/dihedron/overlay/command/base"
	"github.com/gogpu/gg"
)

// CircularArc is the command that adds an arc as an overlay to an image.
type CircularArc struct {
	base.InputCommand
	base.OutputCommand
	// Point is the position in the image where the arc will start.
	Point base.Point `short:"p" long:"point" description:"The coordinates where the arc will be written, as an (x,y) point" optional:"true"`
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
	// Angle defines the angle of the arc.
	Angle base.Point `short:"a" long:"angle" description:"The angle of the arc, as an (start,end) angles in degrees" optional:"true" default:"0,90"`
}

// Execute is the real implementation of the CircularArc command.
func (cmd *CircularArc) Execute(args []string) error {
	slog.Debug("running circular arc command")

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

	slog.Debug("drawing circular arc", "point", cmd.Point, "radius", cmd.Radius, "angle", cmd.Angle)
	dc.DrawArc(cmd.Point.X, cmd.Point.Y, cmd.Radius, cmd.Angle.X/180*math.Pi, cmd.Angle.Y/180*math.Pi)

	if cmd.Fill {
		// TODO: implelemn path closing to create a sector from a chord instead of a sector
		slog.Debug("drawing circular arc as fill", "colour", cmd.Colour)
		dc.Fill()
	} else if cmd.Stroke > 0 {
		slog.Debug("drawing circular arc as stroke", "width", cmd.Stroke)
		dc.SetLineWidth(cmd.Stroke)
		dc.Stroke()
	} else {
		slog.Error("either --fill or --stroke must be specified")
		return fmt.Errorf("either --fill or --stroke must be specified")
	}

	slog.Debug("circular arc overlaid on the image", "point", cmd.Point, "radius", cmd.Radius, "angle", cmd.Angle, "colour", cmd.Colour)

	// write the image to the output stream
	img := dc.Image()
	if err := cmd.WriteOutput(img); err != nil {
		slog.Error("error writing output stream", "name", cmd.Output, "error", err)
		return err
	}
	slog.Debug("image correctly encoded", "filename", cmd.Output, "format", cmd.Format)

	return nil
}
