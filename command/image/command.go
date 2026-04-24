package image

import (
	"errors"
	"image"
	"log/slog"
	"os"

	"github.com/dihedron/overlay/command/base"
	"github.com/gogpu/gg"
	"github.com/jessevdk/go-flags"
)

// Image is the command that superimposes an image as an overlay to the given image.
type Image struct {
	base.OutputCommand
	base.InputCommand
	// Image is the image to superimpose as an overlay to the image.
	Image flags.Filename `short:"y" long:"image" description:"The image to superimpose as an overlay to the given image" optional:"true"`
	// Point is the position in the image where the image will be superimposed.
	Point base.PointF `short:"p" long:"point" description:"The coordinates where the image will be superimposed, as an (x,y) point" optional:"true"`
}

// Execute is the real implementation of the Image command.
func (cmd *Image) Execute(args []string) error {
	slog.Debug("running image command")

	// read the underlay image
	underlay, err := cmd.ReadInput()
	if err != nil {
		slog.Error("error reading input stream", "name", cmd.Input, "error", err)
		return err
	}
	slog.Debug("underlay image decoded", "name", cmd.Input, "width", underlay.Bounds().Dx(), "height", underlay.Bounds().Dy())

	// open the overlay image file
	slog.Debug("reading overlay from file", "name", cmd.Image)
	var f *os.File
	if f, err = os.Open(string(cmd.Image)); err != nil {
		slog.Error("error opening overlay image file", "name", cmd.Image, "error", err)
		return err
	}
	defer f.Close()

	// decode the overlay image
	overlay, _, err := image.Decode(f)
	if err != nil {
		slog.Error("error decoding input data for overlay image", "name", cmd.Image, "error", err)
		os.Exit(1)
	}
	slog.Debug("overlay image decoded", "name", cmd.Image, "width", overlay.Bounds().Dx(), "height", overlay.Bounds().Dy())

	// check if the overlay image is larger than the underlay image
	if overlay.Bounds().Dx() > underlay.Bounds().Dx() || overlay.Bounds().Dy() > underlay.Bounds().Dy() {
		slog.Error("overlay image is larger than the underlay image", "name", cmd.Image)
		return errors.New("overlay image is larger than the underlay image")
	}
	slog.Debug("overlay image is smaller than the underlay image", "name", cmd.Image)

	// create the device context with the underlay image
	dc := gg.NewContextForImage(underlay)
	defer dc.Close()

	// copy the overlay image on the underlay image at the given point
	dc.DrawImage(gg.ImageBufFromImage(overlay), cmd.Point.X, cmd.Point.Y)

	// write the result to the output stream
	img := dc.Image()
	if err := cmd.WriteOutput(img); err != nil {
		slog.Error("error writing output stream", "name", cmd.Output, "error", err)
		return err
	}

	slog.Debug("command done")
	return nil
}
