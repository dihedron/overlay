package text

import (
	"image"
	"image/color"
	"image/draw"
	"log/slog"
	"os"

	"github.com/dihedron/overlay/command/base"
	"github.com/jessevdk/go-flags"
	"golang.org/x/image/font"
	"golang.org/x/image/font/gofont/goregular"
	"golang.org/x/image/font/opentype"
	"golang.org/x/image/math/fixed"
)

// Text is the command that adds text as an overlay to an image.
type Text struct {
	base.InputCommand
	base.OutputCommand
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

	// write to output
	err = cmd.WriteOutput(dst)
	if err != nil {
		slog.Error("error writing output stream", "name", cmd.Output, "error", err)
		return err
	}
	slog.Debug("image correctly encoded", "filename", cmd.Output, "format", cmd.Format)

	return nil
}
