package canvas

import (
	"log/slog"

	"github.com/dihedron/overlay/command/base"
	"github.com/gogpu/gg"
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

	dc := gg.NewContext(cmd.Size.X, cmd.Size.Y)
	defer dc.Close()

	// clear background with uniform colour
	dc.ClearWithColor(gg.RGBA2(float64(cmd.Colour.R), float64(cmd.Colour.G), float64(cmd.Colour.B), float64(cmd.Colour.A)))

	// write the image to the output stream
	img := dc.Image()
	if err := cmd.WriteOutput(img); err != nil {
		slog.Error("error writing output stream", "name", cmd.Output, "error", err)
		return err
	}
	slog.Debug("image correctly encoded", "filename", cmd.Output, "format", cmd.Format)

	return nil
}
