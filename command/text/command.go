package text

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"strings"

	"github.com/dihedron/overlay/command/base"
	"github.com/jessevdk/go-flags"
	"golang.org/x/image/bmp"
	"golang.org/x/image/font"
	"golang.org/x/image/font/gofont/goregular"
	"golang.org/x/image/font/opentype"
	"golang.org/x/image/math/fixed"
)

// Text is the command that adds text as an overlay to an image.
type Text struct {
	base.Command
	// Input is the name of the input file.
	Input flags.Filename `short:"i" long:"input" description:"The name of the input file or - for STDIN" optional:"true" default:"-"`
	// Text is the text to write as an overlay to the image.
	Text string `short:"t" long:"text" description:"The text to add as an overlay to the given image" optional:"true"`
	// Point is the position in the image where the text will start.
	Point base.Point `short:"p" long:"point" description:"The coordinates where the text/image will be written, as an (x,y) point" optional:"true"`
	// Font is the font to use for writing to the image.
	Font flags.Filename `short:"f" long:"font" description:"The name of the font to be used for writing" optional:"true"`
	// Colour is the colour of the font to be used for writing to the image.
	Colour base.Colour `short:"c" long:"colour" description:"The colour of the font to be used for writing" optional:"true" default:"#000000"`
	// Size is the size of font to use for writing to the image.
	Size float64 `short:"s" long:"size" description:"The size of the font to be used for writing" optional:"true" default:"12"`
}

// Execute is the real implementation of the Text command.
func (cmd *Text) Execute(args []string) error {
	slog.Debug("running text command")

	// open the input and output streams
	var (
		input  io.Reader
		output io.Writer
	)
	if cmd.Input == "-" {
		slog.Debug("getting image from STDIN")
		input = os.Stdin
	} else {
		// open the underlay image file
		slog.Debug("reading input from file", "name", cmd.Input)
		var err error
		if input, err = os.Open(string(cmd.Input)); err != nil {
			slog.Error("error opening input file", "name", cmd.Input, "error", err)
			os.Exit(1)
		}
		if input, ok := input.(io.ReadCloser); ok {
			slog.Debug("input needs to be closed at application shutdown", "name", cmd.Input)
			defer input.Close()
		}
	}

	if cmd.Output == "-" {
		slog.Debug("writing image to STDOUT", "format", cmd.Format)
		output = os.Stdout
	} else {
		switch strings.ToLower(filepath.Ext(string(cmd.Output))) {
		case ".jpg", ".jpeg":
			cmd.Format = "jpg"
		case ".png":
			cmd.Format = "png"
		case ".gif":
			cmd.Format = "gif"
		case ".bmp":
			cmd.Format = "bmp"
		default:
			fmt.Fprintf(os.Stderr, "Unsupported output file type: %s\n", filepath.Ext(string(cmd.Output)))
			slog.Error("unsupported output image type", "name", cmd.Output)
			os.Exit(1)
		}
		slog.Debug("writing output to file", "name", cmd.Output, "format", cmd.Format)

		// open the output file
		var err error
		if output, err = os.Create(string(cmd.Output)); err != nil {
			slog.Error("error opening output file", "name", cmd.Output, "error", err)
			os.Exit(1)
		}
		if output, ok := output.(io.WriteCloser); ok {
			slog.Debug("output needs to be closed at application shutdown", "name", cmd.Output)
			defer output.Close()
		}
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

	slog.Debug("overlaying text on the image", "text", cmd.Text)

	// create the font face
	var fnt *opentype.Font
	if cmd.Font != "" {
		// read the font data
		fontData, err := os.ReadFile(string(cmd.Font))
		if err != nil {
			slog.Error("error reading font file", "name", cmd.Font, "error", err)
			os.Exit(1)
		}
		slog.Debug("font data read", "filename", cmd.Font)

		// parse the font data into a font
		if fnt, err = opentype.Parse(fontData); err != nil {
			slog.Error("error parsing font data", "name", cmd.Font, "error", err)
			os.Exit(1)
		}
	} else {
		slog.Debug("using default font")
		if fnt, err = opentype.Parse(goregular.TTF); err != nil {
			slog.Error("error parsing default font data", "name", cmd.Font, "error", err)
			os.Exit(1)
		}
	}
	slog.Debug("font parsed")

	fontFace, err := opentype.NewFace(fnt, &opentype.FaceOptions{
		Size:    float64(cmd.Size),
		DPI:     cmd.DPI,
		Hinting: font.HintingNone,
	})
	if err != nil {
		slog.Error("error creating font face", "name", cmd.Font, "error", err)
		os.Exit(1)
	}

	point := fixed.Point26_6{
		X: fixed.I(cmd.Point.X),
		Y: fixed.I(cmd.Point.Y),
	}

	d := &font.Drawer{
		Dst:  dst,
		Src:  image.NewUniform(color.RGBA(cmd.Colour)),
		Face: fontFace,
		Dot:  point,
	}
	d.DrawString(cmd.Text)
	slog.Debug("text overlaid on the image", "text", cmd.Text, "point", cmd.Point)

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

	return nil
}
