package transform

import "github.com/dihedron/overlay/command/transform/rotate"

// Commands is the set of root command groups.
type Commands struct {
	// Rotate rotates an image.
	Rotate rotate.Rotate `command:"rotate" alias:"r" description:"Rotate an image."`
}
