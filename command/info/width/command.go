package width

import (
	"fmt"
	"log/slog"

	"github.com/dihedron/overlay/command/base"
)

type Width struct {
	base.InputCommand
}

// Execute is the implementation of the width command.
func (cmd *Width) Execute(args []string) error {
	slog.Debug("running width command")

	img, err := cmd.ReadInput()
	if err != nil {
		slog.Error("error reading input stream", "name", cmd.Input, "error", err)
		return err
	}

	fmt.Printf("%d", img.Bounds().Dx())

	return nil
}
