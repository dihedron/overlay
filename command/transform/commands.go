package transform

import (
	"github.com/dihedron/overlay/command/transform/crop"
	"github.com/dihedron/overlay/command/transform/flip"
	"github.com/dihedron/overlay/command/transform/rotate"
	"github.com/dihedron/overlay/command/transform/zoom"
)

// Commands is the set of root transform command groups.
type Commands struct {
	// FlipHorizontally flips an image horizontally.
	FlipHorizontally flip.FlipHorizontally `command:"fliph" alias:"h" description:"Flip an image horizontally."`
	// FlipVertically flips an image vertically.
	FlipVertically flip.FlipVertically `command:"flipv" alias:"v" description:"Flip an image vertically."`
	// Rotate rotates an image.
	Rotate rotate.Rotate `command:"rotate" alias:"r" description:"Rotate an image."`
	// Crop crops an image.
	Crop crop.Crop `command:"crop" alias:"c" description:"Crop an image."`
	// Zoom zooms an image.
	Zoom zoom.Zoom `command:"zoom" alias:"z" description:"Zoom an image."`
}
