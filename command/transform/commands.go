package transform

import (
	"github.com/dihedron/overlay/command/transform/rotate"
	"github.com/dihedron/overlay/command/transform/zoom"
)

// Commands is the set of root command groups.
type Commands struct {
	// Rotate rotates an image.
	Rotate rotate.Rotate `command:"rotate" alias:"r" description:"Rotate an image."`
	// Zoom zooms an image.
	Zoom zoom.Zoom `command:"zoom" alias:"z" description:"Zoom an image."`
}
