package zoom

import (
	"image"
	"log/slog"

	"github.com/anthonynsimon/bild/transform"
	"github.com/dihedron/overlay/command/base"
)

// Zoom zooms an image.
type Zoom struct {
	base.InputCommand
	base.OutputCommand
	// Pivot is the point around which the image will be zoomed.
	Pivot base.Size `short:"p" long:"pivot" description:"The point around which the image will be zoomed, as an (x,y) point" optional:"true"`
	// Factor determines the factor by which the image will be zoomed (>1.0 zooms in, <1.0 zooms out).
	Factor float64 `short:"f" long:"factor" description:"The factor by which the image will be zoomed (>1.0 zooms in, <1.0 zooms out)" optional:"true" default:"2.0"`
}

// Execute is the real implementation of the Zoom command.
func (cmd *Zoom) Execute(args []string) error {
	slog.Debug("running zoom command")

	img, err := cmd.ReadInput()
	if err != nil {
		slog.Error("error reading input stream", "name", cmd.Input, "error", err)
		return err
	}

	opts := &transform.ZoomOptions{
		Pivot: &image.Point{X: cmd.Pivot.X, Y: cmd.Pivot.Y},
	}

	result := transform.Zoom(img, cmd.Factor, opts)
	slog.Debug("image zoomed", "factor", cmd.Factor, "pivot", cmd.Pivot)

	if err := cmd.WriteOutput(result); err != nil {
		slog.Error("error writing output stream", "name", cmd.Output, "error", err)
		return err
	}
	slog.Debug("image correctly encoded", "filename", cmd.Output, "format", cmd.Format)

	return nil
}
