package info

import (
	"github.com/dihedron/overlay/command/info/height"
	"github.com/dihedron/overlay/command/info/sample"
	"github.com/dihedron/overlay/command/info/size"
	"github.com/dihedron/overlay/command/info/width"
)

// Commands is the set of root info command groups.
type Commands struct {
	// Height gets the height of an image.
	Height height.Height `command:"height" alias:"h" description:"Get the height of an image."`
	// Width gets the width of an image.
	Width width.Width `command:"width" alias:"w" description:"Get the width of an image."`
	// Size gets the size of an image.
	Size size.Size `command:"size" alias:"s" description:"Get the size of an image."`
	// Sample gets the color of a pixel in an image.
	Sample sample.Sample `command:"sample" alias:"p" description:"Get the color of a pixel in an image."`
}
