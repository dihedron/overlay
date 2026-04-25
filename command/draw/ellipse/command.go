package ellipse

import (
	"fmt"
	"log/slog"

	"github.com/dihedron/overlay/command/base"
	"github.com/gogpu/gg"
)

// Ellipse is the command that adds an ellipse as an overlay to an image.
type Ellipse struct {
	base.InputCommand
	base.OutputCommand
	// Point is the position in the image where the ellipse will start.
	Point base.Point `short:"p" long:"point" description:"The coordinates where the ellipse will be written, as an (x,y) point" optional:"true"`
	// Colour is the colour of the ellipse to be written to the image.
	Colour base.Colour `short:"c" long:"colour" description:"The colour of the ellipse to be written to the image" optional:"true" default:"#000000"`
	// Fill is whether the ellipse should be filled with the given colour.
	Fill bool `short:"f" long:"fill" description:"Whether the ellipse should be filled with the given colour, by default it is not" optional:"true"`
	// Stroke is the width of the ellipse stroke, when fill is false.
	Stroke float64 `short:"w" long:"stroke" description:"The width of the ellipse stroke, when fill is false" optional:"true" default:"1"`
	// Radius defines the radii (rx and ry) of the ellipse
	Radius base.Point `short:"r" long:"radius" description:"The radii of the ellipse" optional:"true" default:"10,10"`
}

// Execute is the real implementation of the Rectangle command.
func (cmd *Ellipse) Execute(args []string) error {
	slog.Debug("running ellipse command")

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

	dc.DrawEllipse(cmd.Point.X, cmd.Point.Y, cmd.Radius.X, cmd.Radius.Y)

	if cmd.Fill {
		slog.Debug("drawing ellipse as fill", "colour", cmd.Colour)
		dc.Fill()
	} else if cmd.Stroke > 0 {
		slog.Debug("drawing ellipse as stroke", "width", cmd.Stroke)
		dc.SetLineWidth(cmd.Stroke)
		dc.Stroke()
	} else {
		slog.Error("either --fill or --stroke must be specified")
		return fmt.Errorf("either --fill or --stroke must be specified")
	}

	slog.Debug("ellipse overlaid on the image", "point", cmd.Point, "radii", cmd.Radius, "colour", cmd.Colour)

	// write the image to the output stream
	img := dc.Image()
	if err := cmd.WriteOutput(img); err != nil {
		slog.Error("error writing output stream", "name", cmd.Output, "error", err)
		return err
	}
	slog.Debug("image correctly encoded", "filename", cmd.Output, "format", cmd.Format)

	return nil
}
