package flip

import (
	"log/slog"

	"github.com/anthonynsimon/bild/transform"
	"github.com/dihedron/overlay/command/base"
)

// Flips an image vertically.
type FlipVertically struct {
	base.InputCommand
	base.OutputCommand
}

// Execute is the real implementation of the FlipVertically command.
func (cmd *FlipVertically) Execute(args []string) error {
	slog.Debug("running flip vertically command")

	img, err := cmd.ReadInput()
	if err != nil {
		slog.Error("error reading input stream", "name", cmd.Input, "error", err)
		return err
	}

	result := transform.FlipV(img)
	slog.Debug("image flipped vertically")

	if err := cmd.WriteOutput(result); err != nil {
		slog.Error("error writing output stream", "name", cmd.Output, "error", err)
		return err
	}
	slog.Debug("image correctly encoded", "filename", cmd.Output, "format", cmd.Format)

	return nil
}
