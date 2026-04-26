package crop

import (
	"image"
	"log/slog"

	"github.com/anthonynsimon/bild/transform"
	"github.com/dihedron/overlay/command/base"
)

// Crop crops an image to the given rectangle.
type Crop struct {
	base.InputCommand
	base.OutputCommand
	// Rectangle is the rectangle to crop the image to.
	Rectangle base.Rectangle `short:"r" long:"rectangle" description:"The rectangle to crop the image to" optional:"true" default:"0,0,0,0"`
}

// Execute is the real implementation of the Rotate command.
func (cmd *Crop) Execute(args []string) error {
	slog.Debug("running rotate command")

	img, err := cmd.ReadInput()
	if err != nil {
		slog.Error("error reading input stream", "name", cmd.Input, "error", err)
		return err
	}

	result := transform.Crop(img, image.Rectangle{
		Min: image.Point{X: cmd.Rectangle.TopLeft.X, Y: cmd.Rectangle.TopLeft.Y},
		Max: image.Point{X: cmd.Rectangle.BottomRight.X, Y: cmd.Rectangle.BottomRight.Y},
	})
	slog.Debug("image cropped", "rectangle", cmd.Rectangle)

	if err := cmd.WriteOutput(result); err != nil {
		slog.Error("error writing output stream", "name", cmd.Output, "error", err)
		return err
	}
	slog.Debug("image correctly encoded", "filename", cmd.Output, "format", cmd.Format)

	return nil
}
