package base

import (
	"github.com/jessevdk/go-flags"
)

type Command struct {
	// Output is the name of the output file.
	Output flags.Filename `short:"o" long:"output" description:"The name of the output file or - for STDOUT" optional:"true" default:"-"`
	// Format is the output format, if an output filename is not specified; it is used for chaining.
	Format string `short:"x" long:"format" description:"Format of the output image" optional:"true" choice:"jpeg" choice:"jpg" choice:"png" choice:"gif" choice:"bmp" default:"png"`
	// DPI is the image resolution in Dots Per Inch.
	DPI float64 `short:"d" long:"dpi" description:"The image resolution in DPI - Dots Per Inch" optional:"true" default:"72"`
}
