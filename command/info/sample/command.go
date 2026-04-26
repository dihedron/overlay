package sample

import (
	"fmt"
	"log/slog"

	"github.com/dihedron/overlay/command/base"
)

type Sample struct {
	base.InputCommand
	// Point is the point to sample at.
	Point base.Size `short:"p" long:"point" description:"The coordinates where the pixel will be sampled, as an (x,y) point" required:"true"`
}

// Execute is the implementation of the sample command.
func (cmd *Sample) Execute(args []string) error {
	slog.Debug("running sample command")

	img, err := cmd.ReadInput()
	if err != nil {
		slog.Error("error reading input stream", "name", cmd.Input, "error", err)
		return err
	}

	r, g, b, a := img.At(cmd.Point.X, cmd.Point.Y).RGBA()
	r >>= 8
	g >>= 8
	b >>= 8
	a >>= 8
	fmt.Printf("#%02X%02X%02X%02X", r, g, b, a)

	return nil
}
