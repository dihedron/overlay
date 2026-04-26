package height

import (
	"fmt"
	"log/slog"

	"github.com/dihedron/overlay/command/base"
)

type Height struct {
	base.InputCommand
}

// Execute is the implementation of the height command.
func (cmd *Height) Execute(args []string) error {
	slog.Debug("running height command")

	img, err := cmd.ReadInput()
	if err != nil {
		slog.Error("error reading input stream", "name", cmd.Input, "error", err)
		return err
	}

	fmt.Printf("%d", img.Bounds().Dy())

	return nil
}
