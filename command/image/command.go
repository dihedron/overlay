package image

import (
	"fmt"
	"image"
	"image/draw"
	"log/slog"
	"os"

	"github.com/dihedron/overlay/command/base"
	"github.com/jessevdk/go-flags"
)

// Image is the command that superimposes an image as an overlay to the given image.
type Image struct {
	base.OutputCommand
	base.InputCommand
	// Image is the image to superimpose as an overlay to the image.
	Image flags.Filename `short:"y" long:"image" description:"The image to superimpose as an overlay to the given image" optional:"true"`
	// Point is the position in the image where the image will be superimposed.
	Point base.Point `short:"p" long:"point" description:"The coordinates where the image will be superimposed, as an (x,y) point" optional:"true"`
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

	// create a new image with the same dimensions as the original
	dst := image.NewRGBA(underlay.Bounds())
	draw.Draw(dst, dst.Bounds(), underlay, image.Point{0, 0}, draw.Src)
	slog.Debug("image copied to destination context", "width", dst.Bounds().Dx(), "height", dst.Bounds().Dy())

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
		fmt.Fprintf(os.Stderr, "Overlay image is larger than the underlay image\n")
		slog.Error("overlay image is larger than the underlay image", "name", cmd.Image)
		os.Exit(1)
	}
	slog.Debug("overlay image is smaller than the underlay image", "name", cmd.Image)

	//combine the image
	draw.Draw(dst, overlay.Bounds().Add(image.Point(cmd.Point)), overlay, image.Point{0, 0}, draw.Over)

	if err := cmd.WriteOutput(dst); err != nil {
		slog.Error("error writing output stream", "name", cmd.Output, "error", err)
		return err
	}

	slog.Debug("command done")
	return nil
}
