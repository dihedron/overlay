package base

import (
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"strings"

	"github.com/jessevdk/go-flags"
)

type Command struct {
	// Output is the name of the output file.
	Output flags.Filename `short:"o" long:"output" description:"The name of the output file or - for STDOUT" optional:"true" default:"-"`
	// Format is the output format, if an output filename is not specified; it is used for chaining.
	Format string `short:"f" long:"format" description:"Format of the output image" optional:"true" choice:"jpeg" choice:"jpg" choice:"png" choice:"gif" choice:"bmp" default:"png"`
	// DPI is the image resolution in Dots Per Inch.
	DPI float64 `short:"d" long:"dpi" description:"The image resolution in DPI - Dots Per Inch" optional:"true" default:"72"`
}

func (cmd *Command) OutputStream() (io.Writer, error) {
	// open the output stream
	var output io.Writer

	if cmd.Output == "-" {
		slog.Debug("writing image to STDOUT", "format", cmd.Format)
		return os.Stdout, nil
	}

	// check the output format
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
		err := fmt.Errorf("unsupported output file type: %s", filepath.Ext(string(cmd.Output)))
		slog.Error("unsupported output image type", "name", cmd.Output)
		return nil, err
	}
	slog.Debug("writing output to file", "name", cmd.Output, "format", cmd.Format)

	// open the output file
	var err error
	if output, err = os.Create(string(cmd.Output)); err != nil {
		slog.Error("error opening output file", "name", cmd.Output, "error", err)
		return nil, err
	}
	// if output, ok := output.(io.WriteCloser); ok {
	// 	slog.Debug("output needs to be closed at application shutdown", "name", cmd.Output)
	// 	defer output.Close()
	// }

	slog.Debug("output stream ready")
	return output, nil
}
