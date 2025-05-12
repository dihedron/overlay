package main

import (
	"errors"
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
	"runtime"
	"runtime/pprof"
	"strconv"
	"strings"

	"github.com/jessevdk/go-flags"
	"golang.org/x/image/bmp"
	"golang.org/x/image/font"
	"golang.org/x/image/font/gofont/goregular"
	"golang.org/x/image/font/opentype"
	"golang.org/x/image/math/fixed"
)

type Point struct {
	X, Y int
}

func (p *Point) UnmarshalFlag(value string) error {
	parts := strings.Split(value, ",")

	if len(parts) != 2 {
		return errors.New("invalid format: expected two numbers separated by a ,")
	}
	x, err := strconv.ParseInt(parts[0], 10, 32)
	if err != nil {
		return err
	}
	y, err := strconv.ParseInt(parts[1], 10, 32)
	if err != nil {
		return err
	}
	p.X = int(x)
	p.Y = int(y)
	return nil
}

func (p Point) MarshalFlag() (string, error) {
	return fmt.Sprintf("%d,%d", p.X, p.Y), nil
}

type Color color.RGBA

func (c *Color) UnmarshalFlag(value string) error {
	if len(value) > 0 && value[0] != '#' {
		return fmt.Errorf("invalid color string format")
	}

	switch len(value) {
	case 4: // #RGB

		if r, err := strconv.ParseUint(value[1:2], 16, 8); err != nil {
			return err
		} else {
			c.R = uint8(r)
		}

		if g, err := strconv.ParseUint(value[2:3], 16, 8); err != nil {
			return err
		} else {
			c.G = uint8(g)
		}

		if b, err := strconv.ParseUint(value[3:4], 16, 8); err != nil {
			return err
		} else {
			c.B = uint8(b)
		}

		c.A = 255

	case 5: // #RGBA

		if r, err := strconv.ParseUint(value[1:2], 16, 8); err != nil {
			return err
		} else {
			c.R = uint8(r)
		}

		if g, err := strconv.ParseUint(value[2:3], 16, 8); err != nil {
			return err
		} else {
			c.G = uint8(g)
		}

		if b, err := strconv.ParseUint(value[3:4], 16, 8); err != nil {
			return err
		} else {
			c.B = uint8(b)
		}

		if a, err := strconv.ParseUint(value[4:5], 16, 8); err != nil {
			return err
		} else {
			c.A = uint8(a)
		}

	case 7: // #RRGGBB

		if r, err := strconv.ParseUint(value[1:3], 16, 8); err != nil {
			return err
		} else {
			c.R = uint8(r)
		}

		if g, err := strconv.ParseUint(value[3:5], 16, 8); err != nil {
			return err
		} else {
			c.G = uint8(g)
		}

		if b, err := strconv.ParseUint(value[5:7], 16, 8); err != nil {
			return err
		} else {
			c.B = uint8(b)
		}

		c.A = 255

	case 9: // #RRGGBBAA

		if r, err := strconv.ParseUint(value[1:3], 16, 8); err != nil {
			return err
		} else {
			c.R = uint8(r)
		}

		if g, err := strconv.ParseUint(value[3:5], 16, 8); err != nil {
			return err
		} else {
			c.G = uint8(g)
		}

		if b, err := strconv.ParseUint(value[5:7], 16, 8); err != nil {
			return err
		} else {
			c.B = uint8(b)
		}

		if a, err := strconv.ParseUint(value[7:9], 16, 8); err != nil {
			return err
		} else {
			c.A = uint8(a)
		}
	default:
		return fmt.Errorf("invalid color string format ")
	}
	slog.Debug("parsed color", "red", c.R, "green", c.G, "blue", c.B, "alpha", c.A)
	return nil
}

func (c Color) MarshalFlag() (string, error) {
	return fmt.Sprintf("#%02X%02X%02X", c.R, c.G, c.B), nil
}

type Command struct {
	// Input is the name of the input file.
	Input flags.Filename `short:"i" long:"input" description:"The name of the input file or - for STDIN" optional:"true" default:"-"`
	// Output is the name of the output file.
	Output flags.Filename `short:"o" long:"output" description:"The name of the output file or - for STDOUT" optional:"true" default:"-"`
	// Format is the output format, if an output filename is not specified; it is used for chaining.
	Format string `short:"x" long:"format" description:"Format of the output image" optional:"true" choice:"jpeg" choice:"jpg" choice:"png" choice:"gif" choice:"bmp" default:"png"`
	// Text is the text to write as an overlay to the image.
	Text string `short:"t" long:"text" description:"The text to add as an overlay to the given image" optional:"true"`
	// Image is the image to superimpose as an overlay to the image.
	Image string `short:"y" long:"image" description:"The image to superimpose as an overlay to the given image" optional:"true"`
	// Point is the position in the image where the text will start.
	Point Point `short:"p" long:"point" description:"The coordinates where the text/image will be written, as an (x,y) point" optional:"true"`
	// Font is the font to use for writing to the image.
	Font flags.Filename `short:"f" long:"font" description:"The name of the font to be used for writing" optional:"true"`
	// Font is the font to use for writing to the image.
	Color Color `short:"c" long:"color" description:"The color of the font to be used for writing" optional:"true" default:"#000000"`
	// Size is the size of font to use for writing to the image.
	Size float64 `short:"s" long:"size" description:"The size of the font to be used for writing" optional:"true" default:"12"`
	// DPI is the image resolution in Dots Per Inch.
	DPI float64 `short:"d" long:"dpi" description:"The image resolution in DPI - Dots Per Inch" optional:"true" default:"72"`
	// CPUProfile sets the (optional) path of the file for CPU profiling info.
	CPUProfile string `short:"C" long:"cpu-profile" description:"The (optional) path where the CPU profiler will store its data." optional:"yes"`
	// MemProfile sets the (optional) path of the file for memory profiling info.
	MemProfile string `short:"M" long:"mem-profile" description:"The (optional) path where the memory profiler will store its data." optional:"yes"`
}

func (cmd *Command) Execute(args []string) error {
	slog.Debug("start running...", "options", *cmd)

	if cmd.CPUProfile != "" {
		slog.Debug("CPU profiling enabled", "file", cmd.CPUProfile)
		c := cmd.ProfileCPU()
		if c != nil {
			defer c.Close()
		}
	}
	if cmd.MemProfile != "" {
		slog.Debug("memory profiling enabled", "file", cmd.MemProfile)
		cmd.ProfileMemory()
	}

	// validate overlay mode: ony one of image and text is allowed
	if cmd.Text == "" && cmd.Image == "" {
		fmt.Fprintf(os.Stderr, "No text or image specified to overlay on the image\n")
		slog.Error("no text or image specified to overlay on the image")
		os.Exit(1)
	}
	if cmd.Text != "" && cmd.Image != "" {
		fmt.Fprintf(os.Stderr, "Both text and image specified to overlay on the image\n")
		slog.Error("both text and image specified to overlay on the image")
		os.Exit(1)
	}

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

	// now, depending on the mode, overlay the text or the image
	var fnt *opentype.Font
	if cmd.Text != "" {
		slog.Debug("overlaying text on the image", "text", cmd.Text)

		if cmd.Font != "" {

			// read the font data
			fontData, err := os.ReadFile(string(cmd.Font))
			if err != nil {
				slog.Error("error reading font file", "name", cmd.Font, "error", err)
				os.Exit(1)
			}
			slog.Debug("font data read", "filename", cmd.Font)

			// parse the font data into a font
			fnt, err = opentype.Parse(fontData)
			if err != nil {
				slog.Error("error parsing font data", "name", cmd.Font, "error", err)
				os.Exit(1)
			}

		} else {
			slog.Debug("using default font")
			fnt, err = opentype.Parse(goregular.TTF)
			if err != nil {
				slog.Error("error parsing default font data", "name", cmd.Font, "error", err)
				os.Exit(1)
			}
		}
		slog.Debug("font parsed")

		// create the font face
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
			Src:  image.NewUniform(color.RGBA(cmd.Color)),
			Face: fontFace,
			Dot:  point,
		}
		d.DrawString(cmd.Text)
		slog.Debug("text overlaid on the image", "text", cmd.Text, "point", cmd.Point)
	} else {
		slog.Debug("overlaying image on the image", "image", cmd.Image)

		// open the overlay image file
		slog.Debug("reading overlay from file", "name", cmd.Image)
		var (
			err error
			f   io.Reader
		)
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
	}

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

func (cmd *Command) ProfileCPU() *Closer {
	var f *os.File
	if cmd.CPUProfile != "" {
		var err error
		f, err = os.Create(cmd.CPUProfile)
		if err != nil {
			slog.Error("could not create CPU profile", "file", cmd.CPUProfile, "error", err)
		}
		if err := pprof.StartCPUProfile(f); err != nil {
			slog.Error("could not start CPU profiler", "error", err)
		}
	}
	return &Closer{
		file: f,
	}
}

func (cmd *Command) ProfileMemory() {
	if cmd.MemProfile != "" {
		f, err := os.Create(cmd.MemProfile)
		if err != nil {
			slog.Error("could not create memory profile", "file", cmd.MemProfile, "error", err)
		}
		defer f.Close()
		runtime.GC() // get up-to-date statistics
		if err := pprof.WriteHeapProfile(f); err != nil {
			slog.Error("could not write memory profile", "error", err)
		}
	}
}

type Closer struct {
	file *os.File
}

func (c *Closer) Close() {
	if c.file != nil {
		pprof.StopCPUProfile()
		c.file.Close()
	}
}
