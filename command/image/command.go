package image

import (
	"fmt"
	"image"
	"image/draw"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io"
	"log/slog"
	"os"

	"github.com/dihedron/overlay/command/base"
	"github.com/jessevdk/go-flags"
	"golang.org/x/image/bmp"
)

// Image is the command that superimposes an image as an overlay to the given image.
type Image struct {
	base.Command
	// Input is the name of the input file.
	Input flags.Filename `short:"i" long:"input" description:"The name of the input file or - for STDIN" optional:"true" default:"-"`
	// Image is the image to superimpose as an overlay to the image.
	Image string `short:"y" long:"image" description:"The image to superimpose as an overlay to the given image" optional:"true"`
	// Point is the position in the image where the image will be superimposed.
	Point base.Point `short:"p" long:"point" description:"The coordinates where the image will be superimposed, as an (x,y) point" optional:"true"`
}

// Execute is the real implementation of the Image command.
func (cmd *Image) Execute(args []string) error {
	slog.Debug("running image command")

	// open the input and output streams
	var (
		input  io.Reader
		output io.Writer
		err    error
	)

	// open the output stream
	if output, err = cmd.OutputStream(); err != nil {
		slog.Error("error opening output stream", "name", cmd.Output, "error", err)
		return err
	}

	// ensure the output stream is closed at application shutdown
	if output, ok := output.(io.WriteCloser); ok {
		slog.Debug("output needs to be closed at application shutdown", "name", cmd.Output)
		defer output.Close()
	}

	// open the input stream
	if input, err = cmd.InputStream(); err != nil {
		slog.Error("error opening input stream", "name", cmd.Input, "error", err)
		return err
	}

	// ensure the input stream is closed at application shutdown
	if input, ok := input.(io.ReadCloser); ok {
		slog.Debug("input needs to be closed at application shutdown", "name", cmd.Input)
		defer input.Close()
	}

	slog.Debug("streams ready")

	// decode the underlay image
	underlay, _, err := image.Decode(input)
	if err != nil {
		slog.Error("error decoding input data for underlay image", "name", cmd.Input, "error", err)
		os.Exit(1)
	}
	slog.Debug("underlay image decoded", "name", cmd.Input, "width", underlay.Bounds().Dx(), "height", underlay.Bounds().Dy())

	// create a new image with the same dimensions as the original
	dst := image.NewRGBA(underlay.Bounds())
	draw.Draw(dst, dst.Bounds(), underlay, image.Point{0, 0}, draw.Src)

	slog.Debug("image copied to destination context", "width", dst.Bounds().Dx(), "height", dst.Bounds().Dy())

	slog.Debug("overlaying image on the image", "image", cmd.Image)

	// open the overlay image file
	slog.Debug("reading overlay from file", "name", cmd.Image)
	var f io.Reader
	if f, err = os.Open(cmd.Image); err != nil {
		slog.Error("error opening overlay image file", "name", cmd.Image, "error", err)
		os.Exit(1)
	}
	if f, ok := f.(io.ReadCloser); ok {
		slog.Debug("input needs to be closed at application shutdown", "name", cmd.Image)
		defer f.Close()
	}

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

	// encode the output image
	slog.Debug("encoding output image", "name", cmd.Output, "format", cmd.Format)
	switch cmd.Format {
	case "jpg", "jpeg":
		slog.Debug("encoding output file as JPEG", "name", cmd.Output)
		if err = jpeg.Encode(output, dst, nil); err != nil {
			slog.Error("error encoding output file", "name", cmd.Output, "error", err, "format", cmd.Format)
			os.Exit(1)
		}
	case "png":
		slog.Debug("encoding output file as PNG", "name", cmd.Output)
		if err = png.Encode(output, dst); err != nil {
			slog.Error("error encoding output file", "name", cmd.Output, "error", err, "format", cmd.Format)
			os.Exit(1)
		}
	case "gif":
		slog.Debug("encoding output file as GIF", "name", cmd.Output)
		if err = gif.Encode(output, dst, nil); err != nil {
			slog.Error("error encoding output file", "name", cmd.Output, "error", err, "format", cmd.Format)
			os.Exit(1)
		}
	case "bmp":
		slog.Debug("encoding output file as BMP", "name", cmd.Output)
		if err = bmp.Encode(output, dst); err != nil {
			slog.Error("error encoding output file", "name", cmd.Output, "error", err, "format", cmd.Format)
			os.Exit(1)
		}
	default:
		fmt.Fprintf(os.Stderr, "Unsupported output format: %s\n", cmd.Format)
		slog.Error("unsupported output format", "name", cmd.Output, "format", cmd.Format)
		os.Exit(1)
	}
	slog.Debug("image correctly encoded", "filename", cmd.Output, "format", cmd.Format)

	slog.Debug("command done")
	return nil
}

func (cmd *Image) InputStream() (io.Reader, error) {
	// open the input stream
	var (
		input io.Reader
		err   error
	)

	if cmd.Input == "-" {
		slog.Debug("getting image from STDIN")
		return os.Stdin, nil
	}

	// open the underlay image file
	slog.Debug("reading input from file", "name", cmd.Input)
	if input, err = os.Open(string(cmd.Input)); err != nil {
		slog.Error("error opening input file", "name", cmd.Input, "error", err)
		return nil, err
	}

	return input, nil
}
