package main

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

	"github.com/dihedron/overlay/metadata"
	"github.com/jessevdk/go-flags"
	"golang.org/x/image/bmp"
	"golang.org/x/image/font"
	"golang.org/x/image/font/gofont/goregular"
	"golang.org/x/image/font/opentype"
	"golang.org/x/image/math/fixed"
)

func main() {

	if len(os.Args) == 2 && (os.Args[1] == "version" || os.Args[1] == "--version") {
		metadata.Print(os.Stdout)
		os.Exit(0)
	} else if len(os.Args) == 3 && os.Args[1] == "version" && (os.Args[2] == "--verbose" || os.Args[2] == "-v") {
		metadata.PrintFull(os.Stdout)
		os.Exit(0)
	}

	var options Options

	var parser = flags.NewParser(&options, flags.Default)

	if _, err := parser.Parse(); err != nil {
		switch flagsErr := err.(type) {
		case flags.ErrorType:
			if flagsErr == flags.ErrHelp {
				os.Exit(0)
			}
			os.Exit(1)
		default:
			fmt.Fprintf(os.Stderr, "%v\n", err)
			os.Exit(1)
		}
	}

	slog.Debug("start running...", "options", options)

	// validate overlay mode: ony one of image and text is allowed
	if options.Text == "" && options.Image == "" {
		fmt.Fprintf(os.Stderr, "No text or image specified to overlay on the image\n")
		slog.Error("no text or image specified to overlay on the image")
		os.Exit(1)
	}
	if options.Text != "" && options.Image != "" {
		fmt.Fprintf(os.Stderr, "Both text and image specified to overlay on the image\n")
		slog.Error("both text and image specified to overlay on the image")
		os.Exit(1)
	}

	// open the input and output streams
	var (
		input  io.Reader
		output io.Writer
	)
	if options.Input == "-" {
		slog.Debug("getting image from STDIN")
		input = os.Stdin
	} else {
		// open the underlay image file
		slog.Debug("reading input from file", "name", options.Input)
		var err error
		if input, err = os.Open(string(options.Input)); err != nil {
			slog.Error("error opening input file", "name", options.Input, "error", err)
			os.Exit(1)
		}
		if input, ok := input.(io.ReadCloser); ok {
			slog.Debug("input needs to be closed at application shutdown", "name", options.Input)
			defer input.Close()
		}
	}

	if options.Output == "-" {
		slog.Debug("writing image to STDOUT", "format", options.Format)
		output = os.Stdout
	} else {
		switch strings.ToLower(filepath.Ext(string(options.Output))) {
		case ".jpg", ".jpeg":
			options.Format = "jpg"
		case ".png":
			options.Format = "png"
		case ".gif":
			options.Format = "gif"
		case ".bmp":
			options.Format = "bmp"
		default:
			fmt.Fprintf(os.Stderr, "Unsupported output file type: %s\n", filepath.Ext(string(options.Output)))
			slog.Error("unsupported output image type", "name", options.Output)
			os.Exit(1)
		}
		slog.Debug("writing output to file", "name", options.Output, "format", options.Format)

		// open the output file
		var err error
		if output, err = os.Create(string(options.Output)); err != nil {
			slog.Error("error opening output file", "name", options.Output, "error", err)
			os.Exit(1)
		}
		if output, ok := output.(io.WriteCloser); ok {
			slog.Debug("output needs to be closed at application shutdown", "name", options.Output)
			defer output.Close()
		}
	}

	slog.Debug("streams ready")

	// decode the underlay image
	underlay, _, err := image.Decode(input)
	if err != nil {
		slog.Error("error decoding input data for underlay image", "name", options.Input, "error", err)
		os.Exit(1)
	}
	slog.Debug("underlay image decoded", "name", options.Input, "width", underlay.Bounds().Dx(), "height", underlay.Bounds().Dy())

	// create a new image with the same dimensions as the original
	dst := image.NewRGBA(underlay.Bounds())
	draw.Draw(dst, dst.Bounds(), underlay, image.Point{0, 0}, draw.Src)

	slog.Debug("image copied to destination context", "width", dst.Bounds().Dx(), "height", dst.Bounds().Dy())

	// now, depending on the mode, overlay the text or the image
	var f *opentype.Font
	if options.Text != "" {
		slog.Debug("overlaying text on the image", "text", options.Text)

		if options.Font != "" {

			// read the font data
			fontData, err := os.ReadFile(string(options.Font))
			if err != nil {
				slog.Error("error reading font file", "name", options.Font, "error", err)
				os.Exit(1)
			}
			slog.Debug("font data read", "filename", options.Font)

			// parse the font data into a font
			f, err = opentype.Parse(fontData)
			if err != nil {
				slog.Error("error parsing font data", "name", options.Font, "error", err)
				os.Exit(1)
			}

		} else {
			slog.Debug("using default font")
			f, err = opentype.Parse(goregular.TTF)
			if err != nil {
				slog.Error("error parsing default font data", "name", options.Font, "error", err)
				os.Exit(1)
			}
		}
		slog.Debug("font parsed")

		// create the font face
		fontFace, err := opentype.NewFace(f, &opentype.FaceOptions{
			Size:    float64(options.Size),
			DPI:     options.DPI,
			Hinting: font.HintingNone,
		})
		if err != nil {
			slog.Error("error creating font face", "name", options.Font, "error", err)
			os.Exit(1)
		}

		point := fixed.Point26_6{
			X: fixed.I(options.Point.X),
			Y: fixed.I(options.Point.Y),
		}

		d := &font.Drawer{
			Dst:  dst,
			Src:  image.NewUniform(color.RGBA(options.Color)),
			Face: fontFace,
			Dot:  point,
		}
		d.DrawString(options.Text)
		slog.Debug("text overlayed on the image", "text", options.Text, "point", options.Point)
	} else {
		slog.Debug("overlaying image on the image", "image", options.Image)

		// open the overlay image file
		slog.Debug("reading overlay from file", "name", options.Image)
		var (
			err error
			f   io.Reader
		)
		if f, err = os.Open(options.Image); err != nil {
			slog.Error("error opening overlay image file", "name", options.Image, "error", err)
			os.Exit(1)
		}
		if f, ok := f.(io.ReadCloser); ok {
			slog.Debug("input needs to be closed at application shutdown", "name", options.Image)
			defer f.Close()
		}

		// decode the overlay image
		overlay, _, err := image.Decode(f)
		if err != nil {
			slog.Error("error decoding input data for overlay image", "name", options.Image, "error", err)
			os.Exit(1)
		}
		slog.Debug("overlay image decoded", "name", options.Image, "width", overlay.Bounds().Dx(), "height", overlay.Bounds().Dy())

		//offset := image.Pt(overlay.XPos, overlay.YPos)
		//combine the image
		draw.Draw(dst, overlay.Bounds().Add(image.Point(options.Point)), overlay, image.Point{0, 0}, draw.Over)
	}

	switch options.Format {
	case "jpg", "jpeg":
		slog.Debug("encoding output file as JPEG", "name", options.Output)
		if err = jpeg.Encode(output, dst, nil); err != nil {
			slog.Error("error encoding output file", "name", options.Output, "error", err, "format", options.Format)
			os.Exit(1)
		}
	case "png":
		slog.Debug("encoding output file as PNG", "name", options.Output)
		if err = png.Encode(output, dst); err != nil {
			slog.Error("error encoding output file", "name", options.Output, "error", err, "format", options.Format)
			os.Exit(1)
		}
	case "gif":
		slog.Debug("encoding output file as GIF", "name", options.Output)
		if err = gif.Encode(output, dst, nil); err != nil {
			slog.Error("error encoding output file", "name", options.Output, "error", err, "format", options.Format)
			os.Exit(1)
		}
	case "bmp":
		slog.Debug("encoding output file as BMP", "name", options.Output)
		if err = bmp.Encode(output, dst); err != nil {
			slog.Error("error encoding output file", "name", options.Output, "error", err, "format", options.Format)
			os.Exit(1)
		}
	default:
		fmt.Fprintf(os.Stderr, "Unsupported output format: %s\n", options.Format)
		slog.Error("unsupported output format", "name", options.Output, "format", options.Format)
		os.Exit(1)
	}
	slog.Debug("image correctly encoded", "filename", options.Output, "format", options.Format)
}
