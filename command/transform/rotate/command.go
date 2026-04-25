package rotate

import (
	"image"
	"log/slog"

	"github.com/anthonynsimon/bild/transform"
	"github.com/dihedron/overlay/command/base"
)

// Rotate rotates an image.
type Rotate struct {
	base.InputCommand
	base.OutputCommand
	// Angle is the angle in degrees to rotate the image.
	Angle float64 `short:"a" long:"angle" description:"The angle in degrees to rotate the image" optional:"true" default:"90"`
	// Pivot is the point around which the image will be rotated.
	Pivot base.Size `short:"p" long:"pivot" description:"The point around which the image will be rotated, as an (x,y) point" optional:"true"`
	// Resize determines whether the output image should be resized to fit the rotated image.
	Resize bool `short:"r" long:"resize" description:"Whether the output image should be resized to fit the rotated image" optional:"true"`
}

// Execute is the real implementation of the Rotate command.
func (cmd *Rotate) Execute(args []string) error {
	slog.Debug("running rotate command")

	img, err := cmd.ReadInput()
	if err != nil {
		slog.Error("error reading input stream", "name", cmd.Input, "error", err)
		return err
	}

	opts := &transform.RotationOptions{
		ResizeBounds: cmd.Resize,
		Pivot:        &image.Point{X: cmd.Pivot.X, Y: cmd.Pivot.Y},
	}

	result := transform.Rotate(img, cmd.Angle, opts)
	slog.Debug("image rotated", "angle", cmd.Angle, "pivot", cmd.Pivot, "resize", cmd.Resize)

	if err := cmd.WriteOutput(result); err != nil {
		slog.Error("error writing output stream", "name", cmd.Output, "error", err)
		return err
	}
	slog.Debug("image correctly encoded", "filename", cmd.Output, "format", cmd.Format)

	return nil
}
