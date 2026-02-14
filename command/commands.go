package command

import (
	"github.com/dihedron/overlay/command/canvas"
	"github.com/dihedron/overlay/command/image"
	"github.com/dihedron/overlay/command/square"
	"github.com/dihedron/overlay/command/text"
	"github.com/dihedron/overlay/command/version"
)

// Commands is the set of root command groups.
type Commands struct {
	// Canvas creates a new image with the given size and colour.
	Canvas canvas.Canvas `command:"canvas" alias:"c" description:"Create a new image with the given size and colour." `
	// Square adds a square as an overlay to an image.
	Square square.Square `command:"square" alias:"q" description:"Add a square as an overlay to an image." `
	// Image superimposes an image as an overlay to the given image.
	Image image.Image `command:"image" alias:"i" description:"Superimposes an image as an overlay to the given image." `
	// Text adds text as an overlay to an image.
	Text text.Text `command:"text" alias:"t" description:"Add text as an overlay to an image." `
	// Version prints overlay version information and exits.
	Version version.Version `command:"version" alias:"v" description:"Show the command version and exit."`
}
