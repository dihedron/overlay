package size

import (
	"fmt"
	"log/slog"

	"github.com/dihedron/overlay/command/base"
)

type Size struct {
	base.InputCommand
}

// Execute is the implementation of the size command.
func (cmd *Size) Execute(args []string) error {
	slog.Debug("running size command")

	img, err := cmd.ReadInput()
	if err != nil {
		slog.Error("error reading input stream", "name", cmd.Input, "error", err)
		return err
	}

	fmt.Printf("%dx%d", img.Bounds().Dx(), img.Bounds().Dy())

	return nil
}
