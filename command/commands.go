package command

import (
	"github.com/dihedron/overlay/command/draw"
	"github.com/dihedron/overlay/command/transform"
	"github.com/dihedron/overlay/command/version"
)

// Commands is the set of root command groups.
type Commands struct {
	// Draw is the set of subcommands to create a canvas and to paint images, shapes and text over it.
	Draw draw.Commands `command:"draw" alias:"d" description:"Paint images, shapes and text over a canvas"`
	// Transform is the set of subcommands to transform images.
	Transform transform.Commands `command:"transform" alias:"x" description:"Transform images"`
	// Version prints overlay version information and exits.
	Version version.Version `command:"version" alias:"v" description:"Show the command version and exit."`
}
